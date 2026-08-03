# AGENTS.md: platform-test (E2E suite)

Guidance for AI coding agents working under `test/platform-test/`. This is the
Playwright (TypeScript) end-to-end suite that drives a **real** local Kubernetes
cluster running APIM + Gateway + the GKO operator, plus the Terraform APIM
provider. There are **no mocks**: every test mutates live cluster + APIM state.

> Read this before editing or adding tests. For environment bootstrap (`gck`,
> Helm, pre-flight checks) read [`e2e/README.md`](e2e/README.md); for the
> assertion library API read [`README.md`](README.md); for coverage status and
> what is left to migrate read [`PARITY.md`](PARITY.md). Do not duplicate those
> here.

---

## Golden rules

> **Critical: these are the mistakes that have actually broken CI.**

1. **Every test cleans up everything it creates, with a safety net.** Inline
   cleanup does **not** run if the test times out, and a leaked APIM resource
   poisons every later test that reuses the same name. `forEachProvisioner` gives
   you this for free; a hand-rolled test needs its own `afterEach`/`afterAll`.
2. **Never fix a failure by raising a timeout or skipping the test.** A 30s
   timeout that needs 31s is hiding a real reconcile, apply, or consistency
   problem.
3. **Run the test you changed before reporting done.** `npm run e2e -- --grep
   @GKO-xxxx`. Don't claim green without a run.
4. **Use `npm run e2e`, never bare `npx playwright test`.** The bare command
   skips `globalSetup` (the infra pre-flight checks) and gives misleading
   failures.
5. **Clean up in reverse dependency order:** subscriptions → applications →
   APIs. The GKO admission webhook blocks deleting an application that still has
   subscriptions. Never delete shared preconditions: the `dev-ctx`
   ManagementContext, or the APIM org/environment Terraform authenticates against.
6. **Every test carries its provisioner tag.** Lane selection is by title tag
   alone (below). `npm run check:lanes` enforces it.

---

## One authoring model

Everything is provisioned through the **provisioner layer**. A test declares
*intent* once; a provisioner applies it through one path (GKO, Terraform, later
UI/HTTP); the same `mapi` / `gateway` assertions run against the resulting
platform state. A GKO-only concern is simply a scenario that declares one
provisioner.

```
test/platform-test/
  src/                          # @gravitee/platform-test (runner-agnostic; `npm run typecheck`)
    assertions/apim/            # mapi, gateway — shared, provisioner-agnostic assertions
    config/                     # config.yaml loading
    utils/match/                # poll, deepPartialMatch
    provisioners/
      registry.ts               #   the provisioner set — everything else derives from it
      types.ts  view.ts  base.ts
      engines/kubectl.ts  engines/terraform-core.ts
      gko/  terraform/          #   one folder per provisioner
  e2e/
    playwright.config.ts        # serial: workers=1, retries=0, 30s; *.test.ts + *.scenario.ts
    playwright.upgrade.config.ts
    global-setup.ts             # pre-flight: APIM, Gateway, K8s, GKO reachable
    setup.ts                    # fixtures: { mapi, gateway, kubectl, mtlsGatewayBaseUrl } + fixture()
    helpers/
      for-each-provisioner.ts   #   expand one scenario into one tagged test per provisioner
      provisioner-env.ts        #   gkoScenario()/tfScenario(): bind engines + config + paths
      tags.ts                   #   XRAY.* ids, TAGS.REGRESSION, PROVISIONER.*, since()
    tests/
      user-journeys/<persona>/<journey>/   # shared intent, every provisioner
      gko/<area>/                          # the OPERATOR under test
      terraform/                           # the PROVIDER under test
      upgrade/                             # survival across a version change
    fixtures/<domain>/<scenario>/          # legacy central fixtures, used by tests/gko
```

### Where a test goes

| The system under test is… | Goes in | Shape |
|---|---|---|
| the platform (APIM/gateway behaviour a customer sees) | `tests/user-journeys/<persona>/<journey>/` | `*.scenario.ts` via `forEachProvisioner` |
| the operator itself (admission, `.status`, templating, V2) | `tests/gko/<area>/` | `*.test.ts` tagged `PROVISIONER.GKO` |
| the provider itself (drift, plan exit codes, redaction, import) | `tests/terraform/` | `*.test.ts` tagged `PROVISIONER.TERRAFORM` |
| one journey's provisioner-specific slice | the journey folder | `<journey>-gko-only.test.ts` / `-tf-only.test.ts` |

`tests/gko/` is **not a lane** — it is the operator under test. Anything a
customer could also do through Terraform belongs in a journey.

## Adding a user journey

1. **Pick the persona folder** — who performs this journey? `api-producer`,
   `api-consumer`, `platform-admin`. See the
   [catalog](./e2e/tests/user-journeys/README.md).
2. **A new folder only when the _story_ differs, not when the config does.**
   Variants of one story (message-API entrypoints, plan security types, api-key
   modes) are a **named variant table + loop** inside the journey's own
   `*.scenario.ts`, sharing one assertion body.
3. **Confirm the Terraform path first.** Read the provider schema on
   `origin/main` (the local checkout drifts) — see
   [PARITY.md → Regenerating](./PARITY.md#regenerating-this-document). A missing
   attribute is an Automation-API backlog item: file it and mark the arm
   `pending`, do not silently drop the journey.
4. **Author one shared intent**, not two tests. The body must not branch on
   `provisionerId`. Assert the fields that carry regression value
   (`visibility`, `lifecycleState`, `security.type`) so the use-case framing does
   not lose granularity.
5. **De-dup.** Once the journey covers what a standalone GKO/TF test did, **delete
   that test** and carry its Xray id onto the journey arm.
6. **Tag & sync.** Each arm carries its own Xray id; add a `@GKO-TBD-*`
   placeholder for a new arm in `helpers/tags.ts`, then run `/xray-sync-tests`.
   Update the catalog and PARITY.md.

```ts
forEachProvisioner<MyParams>(
  {
    title: "API is started and reachable",
    provisioners: {
      gko: gkoScenario<MyParams>({
        manifests: [path.join(here, "gko/api.yaml")],  // co-located, absolute
        roles: { api: "e2e-v4-keyless" },              // role -> CR name (kind by convention)
        contextPath: "/e2e-v4-keyless",
      }),
      terraform: tfScenario<MyParams>({ fixture: path.join(here, "terraform") }),
    },
    xray: { gko: XRAY.X.GKO_ID, terraform: XRAY.X.TF_ID },  // a list is allowed per arm
    tags: [TAGS.REGRESSION],
    timeoutMs: { gko: 60_000 },                       // TF defaults to TF_WORKSPACE_TIMEOUT_MS
  },
  async ({ provisioned, mapi, gateway }) => {
    await mapi.waitForApiStarted(await provisioned.apiId());
    await gateway.assertResponds(await provisioned.contextPath(), { status: 200 });
  },
  {} as MyParams,
);
```

### Handle surface

`provisioned.apiId()` / `subscriptionId()` / `applicationId()` / `groupId()`
return the APIM UUID (pass a label like `apiId("two-plans")` only when a scenario
has two of the same kind). Plus `contextPath()`, `update(params)`,
`remove(role)`, `destroy()`. Ids are resolved once then cached. The generator
destroys the handle for you, with an `afterEach` safety net that survives a
timeout.

Adding a kind: the getters live once in `BaseProvisioned` and delegate to each
provisioner's `resolveId(role)`. GKO reads `.status.id` of the role's CR (kind by
convention: `api`→apiv4definition, `application`→application,
`subscription`→subscription, `group`→group; use `{ kind, name }` otherwise).
Terraform reads `terraform output` (`api`→`api_id`, `subscription`→`sub_id`,
`application`→`app_id`, `group`→`group_id`; override via `outputs`).

Parameterization that differs structurally per provisioner lives in a co-located
`params.ts` exposing one shared param type plus the per-provisioner closures (GKO
`applyParams`, TF `toVars`). See `api-consumer/subscribe-and-call/`.

### `view` vs `checks`

- **`provisioned.view`** — the agnostic readout: *"did MY layer land this role?"*,
  answered from the provisioner's own record (GKO's CR status, Terraform's state),
  never from APIM. Callable from shared bodies with no narrowing, via
  `assertProvisioner(provisioned, role, "applied" | "failed" | "gone")`.

  Use it **after `remove()`** (expect `"gone"`) and on **failure paths**.
  Do **not** use it as a convergence wait after `update()`: GKO's condition
  `observedGeneration` can lag indefinitely after a re-apply (GKO-2940), so a
  post-update read can return the *pre*-update `Accepted=True` and pass without
  asserting anything. `mapi` is the convergence signal after an update.

  Terraform's `read()` needs `addresses: { role: "apim_x.name" }` in the scenario
  spec (`[0]`-suffixed for a `count`-gated resource) and throws without it.

- **`provisioned.checks`** — the escape hatch for assertions only one provisioner
  can make, narrowed with `isTerraform()`. Terraform-only today: drift,
  redaction, plan exit codes, taint. **GKO declares none** — its control-plane
  readouts *are* the `view` question, and the remaining Kubernetes primitives
  (admission rejection, events) are reached through the kubectl engine in tests
  where nothing is provisioned.

- **Gaps without noise:** a planned-but-unimplemented provisioner goes in
  `pending: { terraform: "<reason>" }` and renders as a visible `test.fixme`,
  never a silent skip. A provisioner absent from both is N/A by design.

## Selecting what to run

`npm run e2e` is the single entry point for every suite; the other scripts are
aliases for a flag on it.

```bash
npm run e2e                            # everything
npm run e2e:regression                 # @regression only
npm run e2e -- --provision-with gko    # one provisioner lane (alias: npm run e2e:gko)
npm run e2e -- --grep @GKO-176         # one test by Xray tag
npm run e2e -- --run-up-to-version 4.11    # skip tests tagged @since-<newer>
npm run e2e -- --suite upgrade --phase before   # alias: npm run e2e:upgrade:before
npm run check:lanes                    # lane partition guard (no cluster needed)
npm run typecheck                      # src; also typecheck:e2e, typecheck:examples
```

**Lane selection is by title tag only.** Every provisioner-specific test carries
`PROVISIONER.GKO` / `PROVISIONER.TERRAFORM` in its `test.describe` title, and
`forEachProvisioner` appends the tag to each generated arm. Nothing keys off the
folder, so the tree can be reorganised freely — but a test with a missing or
wrong tag silently runs in every lane, or none. `npm run check:lanes` asserts the
lanes partition the suite exactly; run it after adding or moving tests.

Do **not** use `--grep @gko`: Playwright's CLI `--grep` is case-insensitive, so
`@gko` also matches every `@GKO-NNNN` Xray tag and selects the whole suite. The
config's own `grepInvert` is case-sensitive for exactly this reason.

The flags are orthogonal and combine: `--provision-with gko --run-up-to-version 4.11`.

## Polling & eventual consistency

Both `kubectl apply` and `terraform apply` return before APIM/Gateway have
converged. Never assert immediately after an apply.

- Use `mapi.waitForApiMatches()` / `expect.poll()` / the `poll()` util, not a
  single-shot assertion. This is the convergence check that matters for **both**
  drivers, since both write to APIM via the Automation API.
- **Combine polled checks atomically:** `expect.poll(() => fetch…).toMatchObject({…})`
  rather than polling one field then re-fetching for the rest.
- **To trigger a GKO reconcile, re-`kubectl apply -f` a modified CR file.** Not
  `kubectl patch`/`annotate`. For Terraform, edit the vars and re-apply.

## Resource isolation

The suite runs **serially with a single worker** against **one shared APIM**:

- **API/App names are a shared global namespace.** If one test leaks a name, the
  next file's apply collides with stuck state and times out, so one root failure
  cascades. Prefer a unique, test-scoped name over reusing an existing fixture's.
- **APIM/MongoDB state persists across cluster restarts.** Only `kind delete
  cluster` or a full Helm uninstall + PV delete wipes it.

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

## APIM behaviours worth knowing

Quirks of the **APIM backend / Automation API**, not the operator:

- **Origin labels:** `origin: MANAGEMENT` = written via mAPI; `origin: KUBERNETES`
  = written via the Automation API — the write path for **both** GKO and
  Terraform, so origin alone does not tell you which driver created a resource.
- **API-key listing returns revoked/expired keys:** no server-side filter. Filter
  client-side on `revoked`/`expired`.
- **API-key values are unique per API**, including already-revoked entries. Custom-key
  tests must generate a per-run unique value.
- **`syncFrom: MANAGEMENT`** is the default for almost all V4 fixtures; not a
  discriminator when triaging.

GKO-specific: **HRID → ID** is `namespace + "-" + name`, from which APIM derives a
deterministic UUIDv3. Use it to correlate CR ↔ APIM API.

## When a test fails

1. **Leaked resource / cascade?** A wave of generic 30s timeouts across unrelated
   suites usually means one earlier test leaked a shared-named resource. Find the
   *first* failure and the missing safety-net cleanup.
2. **Eventual consistency?** Convert a flaky single-shot assertion to
   `poll()` / `expect.poll()`. Don't bump the timeout.
3. **APIM image too old?** Some tests need a fix not yet in the pinned APIM (the
   version comes from the `gravitee-io/gravitee` CircleCI orb, not this repo).

Always consider that the root cause is a **bug in the component under test**, not
the test. Flag it explicitly rather than working around it.

## Committing

> **Critical: no AI attribution on commits or PRs.** Whatever agent you are, do
> **not** add an AI co-author or attribution trailer: no `Co-Authored-By: …`, no
> "Generated with …" footer. Match the repo's style: a `test:` / `docs:` / `fix:`
> prefixed subject and a plain body.

This is **enforced by committed config**, but verify your trailers if your tool
ignores it: [`.claude/settings.json`](.claude/settings.json) sets
`attribution.commit`/`attribution.pr` to empty strings;
[`.cursor/cli-config.json`](.cursor/cli-config.json) sets
`attributeCommitsToAgent`/`attributePRsToAgent` to `false`. Adding a new agent?
Drop its equivalent config under `test/platform-test/` (this suite is
self-contained and may move to its own repo) rather than relying on this prose.

Reports: `playwright-results/` (JUnit XML), `playwright-report/` (HTML).
