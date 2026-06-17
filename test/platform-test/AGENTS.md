# AGENTS.md: platform-test (E2E suite)

Guidance for AI coding agents working under `test/platform-test/`. This is the
Playwright (TypeScript) end-to-end suite that drives a **real** local Kubernetes
cluster running APIM + Gateway + the GKO operator, plus the Terraform APIM
provider. There are **no mocks**: every test mutates live cluster + APIM state.

> Read this before editing or adding tests. For environment bootstrap (`gck`,
> Helm, pre-flight checks) read [`e2e/README.md`](e2e/README.md) instead,
> do not duplicate that here. For the assertion library API read
> [`README.md`](README.md).

---

## Golden rules

> **Critical: these are the mistakes that have actually broken CI.**

1. **Every test cleans up everything it creates, with a safety net.** Inline
   cleanup (`kubectl.del()` for GKO, `terraform.destroyWorkspace()` for the
   Terraform provider) does **not** run if the test times out. A leaked APIM
   resource then poisons every later test that reuses the same name (see
   [Resource isolation](#resource-isolation)). Add a `test.afterEach`/`afterAll`
   safety net to any describe block that creates resources. Patterns below.
2. **Never fix a failure by raising a timeout or skipping the test.**
   Investigate the root cause first. A 30s timeout that needs 31s is hiding a
   real reconcile, apply, or consistency problem.
3. **Run the test you changed before reporting done.** `npm run e2e -- --grep
   @GKO-xxxx` (both GKO and Terraform tests live in the GKO Jira project, so
   both carry `@GKO-` tags). Don't claim green without a run.
4. **Use `npm run e2e`, never bare `npx playwright test`.** The bare command
   skips `globalSetup` (the infra pre-flight checks) and gives misleading
   failures. Scripts are in [`package.json`](package.json).
5. **Clean up in reverse dependency order:** subscriptions → applications →
   APIs. For GKO the admission webhook blocks deleting an application that still
   has subscriptions; for Terraform, `terraform destroy` walks the dependency
   graph for you. Never delete shared preconditions other tests rely on: the GKO
   `dev-ctx` ManagementContext, or the APIM org/environment the Terraform
   provider authenticates against.

---

## Layout

```
test/platform-test/
  src/                            # @gravitee/platform-test library (runner-agnostic, typechecked by `npm run typecheck`)
    assertions/apim/              # mapi, gateway (shared, driver-agnostic assertions)
    utils/match/                  # poll, deepPartialMatch
    provisioners/                 # provisioner layer: one pluggable Provisioner per creation path
      types.ts                    #   Provisioner / Provisioned / ResourceRef / DriverId / ProvisionerChecks
      engines/kubectl.ts          #   kubectl CLI wrapper (moved here from e2e/helpers)
      engines/terraform-core.ts   #   pure Terraform workspace mechanics (no config/fixture coupling)
      gko/  terraform/            #   one folder per provisioner: GkoProvisioner + GkoChecks, etc. (add more here)
  e2e/
    playwright.config.ts          # serial: workers=1, retries=0, 30s; testMatch *.test.ts + *.scenario.ts
    global-setup.ts               # pre-flight: APIM, Gateway, K8s, GKO reachable
    setup.ts                      # Playwright fixtures: { mapi, gateway, kubectl, mtlsGatewayBaseUrl } + fixture()
    helpers/
      kubectl.ts terraform.ts     # thin shims/adapters over src/provisioners/engines (existing imports still work)
      for-provisioners.ts         # forProvisioners(): expand a scenario into 1 tagged test per provisioner
      provisioner-env.ts          # gkoScenario()/tfScenario(): bind engines + config + fixture() to provisioners
      tags.ts                     # XRAY.* test-ID constants + TAGS.REGRESSION
    fixtures/<domain>/<scenario>/ # one folder per scenario, referenced via fixture("<domain>/<scenario>/...")
      crd.yaml                    # GKO CRD manifest(s), often multi-doc (---) for paired resources
      main.tf                     # Terraform provider config; outputs api_id / sub_id / api_context_path (+ <role>_id)
    tests/gko/<area>/             # GKO-only operator tests
    tests/terraform/              # Terraform-only provider tests
    tests/scenarios/<domain>/     # <name>.scenario.ts: one shared intent, run across every supported provisioner
```

> Fixtures are organised **by domain then scenario** (e.g.
> `subscriptions/apikey-auto/`), not by file type. A scenario folder holds the
> `crd.yaml` (and/or `main.tf`) for that one case; paired resources that must be
> applied together live as multi-doc sections in a single `crd.yaml`.

## Test authoring conventions

- **Imports:** `import { test, fixture, expect } from "../../../setup.js";`.
  The `test` is the extended one with the `mapi` / `gateway` / `kubectl`
  fixtures. Use the `.js` extension in import paths (ESM/NodeNext).
- **Fixtures (manifests):** resolve with
  `fixture("<domain>/<scenario>/crd.yaml")` (or `main.tf` for Terraform). Do not
  hardcode absolute paths. Add a new scenario as its own folder rather than
  dropping a loose file into an existing one.
- **Xray tagging:** every test title ends with its Xray ID:
  `` test(`Description ${XRAY.AREA.TEST_ID}`, ...) ``. IDs
  live in [`e2e/helpers/tags.ts`](e2e/helpers/tags.ts). New tests get a real
  Jira Test ID.
- **Steps:** group phases with `test.step("...", async () => {...})` for
  readable reports.
- **One test = one concern.** Prefer many small tests over a mega-scenario; it
  localises failures and shrinks the blast radius of a leaked resource.

## Provisioner layer: one intent, every provisioner

For behaviour that should hold no matter how a resource was created, write it ONCE as a
`*.scenario.ts` under `tests/scenarios/<domain>/` and let `forProvisioners` run it against every
provisioner the scenario supports. The current provisioners are GKO and Terraform, and the layer is
built to grow: adding another (e.g. a UI path) means implementing the `Provisioner` interface under
`src/provisioners/` and listing it in the scenario, with no change to the scenario bodies or the
shared assertions. The shared body uses only the provisioner-agnostic handle (`provision`) plus
`mapi`/`gateway`; each generated test's title carries its provisioner tag (e.g. `@gko`, `@terraform`)
and the per-provisioner Xray id.

**Selecting a provisioner lane:** `npm run e2e -- --provision-with gko` (or the `npm run e2e:gko` /
`e2e:terraform` shortcuts) runs only that provisioner's tests — the matrix arms plus the
`*-gko-only` / `*-tf-only` files, which carry the tag on their `describe`. Do **not** use
`--grep @gko`: Playwright's CLI `--grep` is case-insensitive, so `@gko` also matches every `@GKO-NNNN`
Xray tag and selects the whole suite. The lane filter is a case-sensitive grep applied in
`playwright.config.ts` from the `E2E_PROVISIONER` env var (which `--provision-with` sets); `--grep
@GKO-NNNN` still works for selecting a single test. The `scripts/e2e.mjs` wrapper also accepts
`--run-up-to-version <semver>` as a reserved seam for future version-gating (accepted but not yet
enforced).

```ts
import { forProvisioners } from "../../../../helpers/for-provisioners.js";
import { gkoScenario, tfScenario } from "../../../../helpers/provisioner-env.js";

forProvisioners<MyParams>(
  {
    title: "API is started and reachable",
    provisioners: {
      gko: gkoScenario<MyParams>({
        manifests: ["plans/v4-keyless/crd.yaml"],   // fixture-relative
        roles: { api: "e2e-v4-keyless" },           // role -> CR name (kind by convention)
        contextPath: "/e2e-v4-keyless",
      }),
      terraform: tfScenario<MyParams>({ fixture: "plans/v4-keyless" }), // folder with main.tf
    },
    xray: { gko: XRAY.X.GKO_ID, terraform: XRAY.X.TF_ID }, // a list is allowed per provisioner
    tags: [TAGS.REGRESSION],
    timeoutMs: { gko: 60_000 },                     // TF defaults to TF_WORKSPACE_TIMEOUT_MS
  },
  async ({ provision, mapi, gateway }) => {
    await mapi.waitForApiStarted(await provision.id("api"));
    await gateway.assertResponds(await provision.contextPath(), { status: 200 });
  },
  {} as MyParams,                                   // initial params
);
```

Rules of thumb:
- **Handle surface:** `provision.id(role?)`, `provision.ref(role?)`, `provision.contextPath()`,
  `provision.update(params)` (rotation-style re-provision), `provision.destroy()`. `id`/`contextPath`
  are resolved once then cached. The generator destroys the handle for you, with an `afterEach` safety
  net that survives a test timeout.
- **Roles -> ids:** GKO reads `.status.id` of the role's CR (kind derived by convention:
  `api`->apiv4definition, `application`->application, `subscription`->subscription; use the full
  `{ kind, name }` form otherwise). Terraform reads `terraform output` (`api`->`api_id`,
  `subscription`->`sub_id`, `application`->`app_id`; override via `outputs`). Paired TF fixtures MUST
  expose `api_id`, `sub_id`, `api_context_path` (and `<role>_id` for extra roles).
- **Parameterization** that differs structurally per provisioner (e.g. "set the api-keys") lives in a
  small co-located `params.ts` exposing one shared param type plus the per-provisioner apply closures
  (the GKO `applyParams` closure, the TF `toVars` closure). See
  `tests/scenarios/subscriptions/apikey/` for the reference pilot.
- **Provisioner-specific assertions** (no shared-layer home) go in `provision.checks`, narrowed by a
  per-provisioner type guard (`isGko(...)` / `isTerraform(...)`): GKO conditions/events/`.status`, TF
  drift/idempotency/redaction. Behaviour whose *assertion* (not just provisioning) is
  provisioner-specific stays in a plain `*-gko-only.test.ts` / `*-tf-only.test.ts` rather than the matrix.
- **Gaps without noise:** a planned-but-unimplemented provisioner goes in
  `pending: { terraform: "<reason or tracking ref>" }` and renders as a visible `test.fixme`, never a
  silent skip. A provisioner simply absent from `provisioners`/`pending` emits nothing (N/A by design).

## Polling & eventual consistency

Both `kubectl apply` (GKO) and `terraform apply` (provider) return before
APIM/Gateway have converged. Never assert immediately after an apply.

- **GKO:** use `kubectl.waitForCondition("apiv4definition", name, "Accepted")`
  to wait for the operator, then assert.
- Use `mapi.waitForApiMatches()` / `expect.poll()` / the `poll()` util for APIM
  and gateway state, not a single-shot assertion. This is the convergence check
  that matters for **both** drivers, since both ultimately write to APIM via the
  Automation API.
- **Combine polled checks atomically:** `expect.poll(() => fetch...).toMatchObject({...})`
  rather than polling one field then re-fetching for the rest, which avoids
  races where state changes between calls.
- **To trigger a reconcile, re-`kubectl apply -f` a modified CR file.** Do not
  use `kubectl patch`/`annotate`; the operator reconciles on spec changes via
  apply. (`helpers/kubectl.ts` exposes `apply`, `applyString`,
  `applyExpectFailure`, `delExpectFailure`, `getStatus`, `exists`,
  `waitForDeletion`, `rolloutRestart`, …) For Terraform, edit the `.tf`/vars and
  re-`apply` through the `terraform` helper (`initWorkspace`, `apply`, `plan`,
  `output`, `destroyWorkspace`).

## Resource isolation

The suite runs **serially with a single worker** and tests **share one APIM
backend**. Two consequences agents must respect:

- **API/App names are a shared global namespace.** The same name (e.g.
  `e2e-v4-sync-mgmt`) is reused across multiple test files. If one test leaks it,
  the next file's apply collides with stuck state and times out, so one root
  failure cascades into many. When adding a test, prefer a **unique, test-scoped
  name** over reusing an existing fixture's name.
- **APIM/MongoDB state persists across cluster restarts.** Only `kind delete
  cluster` or a full Helm uninstall + PV delete wipes it. A half-cleaned test
  leaves rows behind that survive `make start-e2e-cluster`.

### Safety-net cleanup pattern

Add to any describe block that creates resources, so cleanup runs even when the
test body times out before its inline cleanup. Pick the variant for your driver.

**GKO (`kubectl`):** module-level import, reverse dependency order, `del()`
ignores errors (resource may already be gone), never delete shared `dev-ctx`:

```ts
import * as kubectl from "../../../helpers/kubectl.js";

test.describe("…", () => {
  test.afterEach(async () => {
    for (const f of [
      "<domain>/<sub-scenario>/crd.yaml",
      "<domain>/<app-scenario>/crd.yaml",
      "<domain>/<api-scenario>/crd.yaml",
    ]) {
      await kubectl.del(fixture(f)).catch(() => {});
    }
  });

  // tests…
});
```

**Terraform (`terraform`):** track the workspace and tear it down in
`afterAll`. `destroyWorkspace` re-runs `destroy` as a no-op if a test already
destroyed inline, so it is always safe to call:

```ts
import * as terraform from "../../helpers/terraform.js";

test.describe("…", () => {
  let ws: terraform.TfWorkspace | undefined;

  test.afterAll(async () => {
    if (ws) await terraform.destroyWorkspace(ws).catch(() => {});
  });

  test("…", async () => {
    ws = await terraform.initWorkspace("<domain>/<scenario>"); // folder with main.tf
    await terraform.apply(ws);
    // assert against APIM via mapi/poll, then let afterAll destroy
  });
});
```

## APIM behaviours worth knowing (save yourself a debugging session)

These are quirks of the **APIM backend / Automation API**, not the operator.
They surface in e2e because the suite asserts against live APIM state.

- **Origin labels:** APIM `origin: MANAGEMENT` = written via mAPI;
  `origin: KUBERNETES` = written via the Automation API, which is the write path
  for **both** GKO and the Terraform provider (so origin alone does not tell you
  which driver created a resource).
- **API-key listing returns revoked/expired keys:** the endpoint has no
  server-side filter. Filter client-side on `revoked`/`expired`.
- **API-key values are unique per API:** including already-revoked entries.
  Custom-key tests must generate a per-run unique value.
- **`syncFrom: MANAGEMENT`** is the default for almost all V4 fixtures; it lets
  management-plane edits flow back. Not a discriminator when triaging failures.

The one GKO-specific correlation rule:

- **HRID → ID:** GKO's human-readable ID is `namespace + "-" + name`; APIM then
  derives a deterministic UUIDv3 from it. Use it to correlate CR ↔ APIM API.



## Committing

> **Critical: no AI attribution on commits or PRs.** Whatever agent you are
> (Claude, Cursor, Copilot, …), do **not** add an AI co-author or attribution
> trailer: no `Co-Authored-By: …`, no "Generated with …" / "Made with …"
> footer. Match the repo's existing commit style: a `test:` / `docs:` / `fix:`
> prefixed subject and a plain body. PR descriptions are Summary-only.

This is **enforced by committed config**, so you should not have to think about
it, but verify your trailers if your tool ignores the config (some CLIs do):

- **Claude Code:** [`.claude/settings.json`](.claude/settings.json) sets
  `attribution.commit` and `attribution.pr` to empty strings.
- **Cursor:** [`.cursor/cli-config.json`](.cursor/cli-config.json) sets
  `attribution.attributeCommitsToAgent` and `attributePRsToAgent` to `false`.

Adding a new agent? Drop its equivalent config under `test/platform-test/`
(this suite is self-contained and may move to its own repo) rather than relying
on this prose.

## Commands

```bash
npm run e2e                          # all E2E tests (runs globalSetup)
npm run e2e:regression               # @regression suite only
npm run e2e -- --provision-with gko  # only the GKO provisioner lane (matrix + *-gko-only)
npm run e2e:terraform                # only the Terraform lane (shortcut for --provision-with terraform)
npm run e2e -- --grep @GKO-176       # single test by Xray tag
npm run typecheck                    # tsc --noEmit, run before committing
```

Reports: `playwright-results/` (JUnit XML), `playwright-report/` (HTML).

## When a test fails

Triage in this order before touching the test:

1. **Is it a leaked resource / cascade?** A wave of generic 30s timeouts across
   unrelated suites usually means one earlier test leaked a shared-named
   resource. Look for the *first* failure and the missing safety-net cleanup.
2. **Is it eventual consistency?** A flaky single-shot assertion → convert to
   `poll()` / `expect.poll()`, don't bump the timeout.
3. **Is the APIM image too old?** Some tests need a fix not yet in the pinned
   APIM (the version comes from the `gravitee-io/gravitee` CircleCI orb, not this
   repo). Confirm the image contains the fix commit before re-enabling.
