# AGENTS.md: platform-test (E2E suite)

@.agent/rules/e2e-test-authoring.md
@.agent/rules/e2e-framework-development.md

Guidance for AI coding agents working under `test/platform-test/`. This is the
Playwright (TypeScript) end-to-end suite that drives a **real** local Kubernetes
cluster running APIM + Gateway + the GKO operator, plus the Terraform APIM
provider. There are **no mocks**: every test mutates live cluster + APIM state.

## Start here

| You are… | Read | Then invoke |
|---|---|---|
| adding or changing a **test** | [`.agent/rules/e2e-test-authoring.md`](.agent/rules/e2e-test-authoring.md) | `write-e2e-test` |
| chasing a **failing test** | same | `investigate-e2e-failure` |
| adding a **provisioner** or engine | [`.agent/rules/e2e-framework-development.md`](.agent/rules/e2e-framework-development.md) | `add-provisioner` |
| extending the **assertion library** | same | `extend-platform-assertions` |

The two rules files are the invariants and are imported above, so they load with
this file. The four skills live in [`.claude/skills/`](.claude/skills/) and carry
the step-by-step procedures.

For environment bootstrap (`gck`, Helm, pre-flight checks) read
[`e2e/README.md`](e2e/README.md); for the assertion library API read
[`README.md`](README.md); for coverage status and what is left to migrate read
[`PARITY.md`](PARITY.md). Do not duplicate those here.

---

## One authoring model

Everything is provisioned through the **provisioner layer**. A test declares
*intent* once; a provisioner applies it through one path (GKO, Terraform, later
UI/HTTP); the same `mapi` / `gateway` assertions run against the resulting
platform state. A GKO-only concern is simply a scenario that declares one
provisioner.

```
test/platform-test/
  .agent/rules/                 # the invariants, per audience
  .claude/skills/               # the procedures, per audience
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

Full procedure: the [`write-e2e-test`](.claude/skills/write-e2e-test/SKILL.md)
skill. Two rules decide placement, and they are the ones most often got wrong:

1. **Persona picks the folder** — who performs this journey? `api-producer`,
   `api-consumer`, `platform-admin`. See the
   [catalog](./e2e/tests/user-journeys/README.md).
2. **A new folder only when the _story_ differs, not when the config does.**
   Variants of one story (message-API entrypoints, plan security types, api-key
   modes) are a **named variant table + loop** inside the journey's own
   `*.scenario.ts`, sharing one assertion body.

Then: confirm the Terraform path against the provider schema on `origin/main`
(a missing attribute is an Automation-API backlog item to file, not a reason to
drop the arm), author **one shared intent** whose body never branches on
`provisionerId`, delete any `tests/gko/` test the journey supersedes and carry
its Xray id onto the GKO arm, tag both arms, and update the catalog and
PARITY.md.

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

Reports: `playwright-results/` (JUnit XML), `playwright-report/` (HTML).
