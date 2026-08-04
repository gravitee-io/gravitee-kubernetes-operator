# Rules: writing e2e tests

Audience: any agent **adding or changing a test** under `test/platform-test/e2e/`.

The suite drives a **real** Kubernetes cluster running APIM + Gateway + the GKO
operator, plus the Terraform APIM provider. There are no mocks: every test
mutates live cluster and APIM state, serially, against one shared APIM.

Procedure lives in the `write-e2e-test` skill. These are the invariants that
hold whether or not you invoke it.

---

## Golden rules

> These are the mistakes that have actually broken CI.

1. **Every test cleans up everything it creates, with a safety net.** Inline
   cleanup does not run if the test times out, and a leaked APIM resource
   poisons every later test that reuses the same name. `forEachProvisioner`
   gives you this for free; a hand-rolled test needs its own
   `afterEach`/`afterAll`.
2. **Never fix a failure by raising a timeout or skipping the test.** A 30s
   timeout that needs 31s is hiding a real reconcile, apply, or consistency
   problem.
3. **Run the test you changed before reporting done.**
   `npm run e2e -- --grep @GKO-xxxx`. Do not claim green without a run.
4. **Use `npm run e2e`, never bare `npx playwright test`.** The bare command
   skips `globalSetup` (the infra pre-flight checks) and gives misleading
   failures.
5. **Clean up in reverse dependency order:** subscriptions, then applications,
   then APIs. The GKO admission webhook blocks deleting an application that
   still has subscriptions. Never delete shared preconditions: the `dev-ctx`
   ManagementContext, or the APIM org/environment Terraform authenticates
   against.
6. **Every test carries its provisioner tag.** Lane selection is by title tag
   alone. `npm run check:lanes` enforces it.

## Polling and eventual consistency

Both `kubectl apply` and `terraform apply` return before APIM/Gateway have
converged. Never assert immediately after an apply.

- Use `mapi.waitForApiMatches()` / `expect.poll()` / the `poll()` util, not a
  single-shot assertion. This is the convergence check that matters for **both**
  drivers, since both write to APIM via the Automation API.
- **Combine polled checks atomically:** `expect.poll(() => fetch…).toMatchObject({…})`
  rather than polling one field then re-fetching for the rest.
- **To trigger a GKO reconcile, re-`kubectl apply -f` a modified CR file.** Not
  `kubectl patch` or `kubectl annotate`. For Terraform, edit the vars and
  re-apply.

## Resource isolation

The suite runs **serially with a single worker** against **one shared APIM**:

- **API/App names are a shared global namespace.** If one test leaks a name, the
  next file's apply collides with stuck state and times out, so one root failure
  cascades. Prefer a unique, test-scoped name over reusing an existing fixture's.
- **APIM/MongoDB state persists across cluster restarts.** Only `kind delete
  cluster` or a full Helm uninstall plus PV delete wipes it.

### Safety-net cleanup (hand-rolled tests only)

`forEachProvisioner` handles this for you. A test that provisions directly needs:

```ts
import * as kubectl from "../../../helpers/kubectl.js";

test.describe(`… ${PROVISIONER.GKO}`, () => {
  test.afterEach(async () => {
    for (const f of ["<sub>/crd.yaml", "<app>/crd.yaml", "<api>/crd.yaml"]) {
      await kubectl.del(fixture(f)).catch(() => {});   // reverse dependency order
    }
  });
});
```

For Terraform, track the workspace and `destroyWorkspace(ws)` in `afterAll`; it
re-runs `destroy` as a no-op if a test already destroyed inline, so it is always
safe to call.

## `view` vs `checks`

- **`provisioned.view`** is the agnostic readout: *"did MY layer land this
  role?"*, answered from the provisioner's own record (GKO's CR status,
  Terraform's state), never from APIM. Callable from shared bodies with no
  narrowing, via `assertProvisioner(provisioned, role, "applied" | "failed" | "gone")`.

  Use it **after `remove()`** (expect `"gone"`) and on **failure paths**.
  Do **not** use it as a convergence wait after `update()`: GKO's condition
  `observedGeneration` can lag indefinitely after a re-apply (GKO-2940), so a
  post-update read can return the *pre*-update `Accepted=True` and pass without
  asserting anything. `mapi` is the convergence signal after an update.

  Terraform's `read()` needs `addresses: { role: "apim_x.name" }` in the scenario
  spec (`[0]`-suffixed for a `count`-gated resource) and throws without it.

- **`provisioned.checks`** is the escape hatch for assertions only one
  provisioner can make, narrowed with `isTerraform()`. Terraform-only today:
  drift, redaction, plan exit codes, taint. **GKO declares none**: its
  control-plane readouts *are* the `view` question, and the remaining Kubernetes
  primitives (admission rejection, events) are reached through the kubectl engine
  in tests where nothing is provisioned.

- **Gaps without noise:** a planned-but-unimplemented provisioner goes in
  `pending: { terraform: "<reason>" }` and renders as a visible `test.fixme`,
  never a silent skip. A provisioner absent from both is N/A by design.

## Xray tagging

- Test ids live **only** in [`e2e/helpers/tags.ts`](../../e2e/helpers/tags.ts) as
  `XRAY.*` constants, and are interpolated into test titles:
  ``test(`Description ${XRAY.CATEGORY.TEST_ID}`, …)``.
- A new test with no ticket yet gets a `@GKO-TBD-*` placeholder in `tags.ts`;
  run `/xray-sync-tests` to file the real Test tickets and rewrite the file.
- Each provisioner arm carries **its own** id: GKO and Terraform are different
  Jira tickets.
- E2E files reference **Xray Test ids only**. Never a Story or bug key in a
  describe title, docstring, comment or fixture name.
- `TAGS.REGRESSION` marks the regression pack. `since("4.12")` declares the
  oldest APIM version a test needs, so `--run-up-to-version` skips it on older
  clusters instead of failing.

## APIM behaviours worth knowing

Quirks of the **APIM backend / Automation API**, not the operator:

- **Origin labels:** `origin: MANAGEMENT` = written via mAPI; `origin: KUBERNETES`
  = written via the Automation API, the write path for **both** GKO and
  Terraform, so origin alone does not tell you which driver created a resource.
- **API-key listing returns revoked/expired keys:** no server-side filter. Filter
  client-side on `revoked`/`expired`.
- **API-key values are unique per API**, including already-revoked entries.
  Custom-key tests must generate a per-run unique value.
- **`syncFrom: MANAGEMENT`** is the default for almost all V4 fixtures; not a
  discriminator when triaging.

GKO-specific: **HRID to ID** is `namespace + "-" + name`, from which APIM derives
a deterministic UUIDv3. Use it to correlate CR to APIM API.

## When a test fails

Always consider that the root cause is a **bug in the component under test**, not
the test. Flag it explicitly rather than working around it. The
`investigate-e2e-failure` skill carries the full decision tree.

## Committing

> **No AI attribution on commits or PRs.** Whatever agent you are, do not add an
> AI co-author or attribution trailer: no `Co-Authored-By: …`, no "Generated
> with …" footer.

- Conventional Commits, `test:` / `docs:` / `fix:` prefix, plain body.
- Wrap commit bodies at ~72 columns: CI fails on `body-max-line-length` > 120.
- Ask before committing or pushing.

Attribution is **enforced by committed config**, but verify your trailers if your
tool ignores it: [`.claude/settings.json`](../../.claude/settings.json) sets
`attribution.commit`/`attribution.pr` to empty strings;
[`.cursor/cli-config.json`](../../.cursor/cli-config.json) sets
`attributeCommitsToAgent`/`attributePRsToAgent` to `false`. Adding a new agent?
Drop its equivalent config under `test/platform-test/` (this suite is
self-contained and may move to its own repo) rather than relying on this prose.
