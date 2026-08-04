---
name: write-e2e-test
description: Author or migrate a Playwright e2e test in test/platform-test - pick the right shape (shared user journey vs GKO-only vs Terraform-only), write the scenario and both provisioners' fixtures, tag it for Xray, run it, and update the catalog and PARITY docs. Use when asked to add e2e coverage, automate an Xray test, cover a new APIM/GKO/Terraform behaviour, or migrate a tests/gko test into a user journey.
allowed-tools: Bash, Read, Edit, Write, Glob, Grep
---

# write-e2e-test

Procedure for adding coverage to the Gravitee platform e2e suite at
`test/platform-test/`. The invariants this rests on are in
[`.agent/rules/e2e-test-authoring.md`](../../../.agent/rules/e2e-test-authoring.md);
read them if they are not already in context.

Work from `test/platform-test/` unless stated otherwise.

## References

Load these on demand, not upfront:

| File | When |
|---|---|
| `references/journey-template.md` | writing a new journey: skeletons for the scenario, `gko/*.yaml`, `terraform/main.tf`, `README.md` |
| `references/assertions-cheatsheet.md` | choosing an assertion: the full `mapi` / `gateway` / `poll` surface |
| `references/gko-only-test.md` | writing a hand-rolled `tests/gko/` test: safety-net cleanup, kubectl engine, admission rejection |

---

## Step 1: decide the shape

| The system under test is… | Goes in | Shape |
|---|---|---|
| the platform (APIM/gateway behaviour a customer sees) | `e2e/tests/user-journeys/<persona>/<journey>/` | `*.scenario.ts` via `forEachProvisioner` |
| the operator itself (admission, `.status`, templating, V2) | `e2e/tests/gko/<area>/` | `*.test.ts` tagged `PROVISIONER.GKO` |
| the provider itself (drift, plan exit codes, redaction, import) | `e2e/tests/terraform/` | `*.test.ts` tagged `PROVISIONER.TERRAFORM` |
| one journey's provisioner-specific slice | the journey folder | `<journey>-gko-only.test.ts` / `-tf-only.test.ts` |

`tests/gko/` is not a lane, it is the operator under test. If a customer could do
the same thing through Terraform, it is a journey.

**Default to a journey.** Only fall back to a single-provisioner test when the
thing being asserted is a Kubernetes or HCL mechanic, not a platform behaviour.

For a `tests/gko/` or `tests/terraform/` test, jump to
`references/gko-only-test.md` and skip to Step 6.

## Step 2: place the journey

1. **Persona picks the folder.** `api-producer` (designing, publishing,
   securing, documenting an API), `api-consumer` (registering an application,
   subscribing, calling the gateway), `platform-admin` (groups, members,
   dictionaries, environment config). Check
   [the catalog](../../../e2e/tests/user-journeys/README.md) for an existing home.
2. **A new folder only when the _story_ differs, not when the config does.**
   Variants of one story (message-API entrypoints, plan security types, api-key
   modes) are a **named variant table plus a loop** inside the journey's own
   `*.scenario.ts`, sharing one assertion body. Adding a folder per config
   permutation is the single most common mistake here.

Folder shape:

```
<persona>/<journey>/
  <journey>.scenario.ts   # the shared intent, run against every provisioner
  params.ts               # only if the two drivers need structurally different params
  gko/        *.yaml      # the GKO custom resources
  terraform/  main.tf     # the Terraform equivalent
  README.md               # "As a <persona>, I ..." + how to run it
```

## Step 3: confirm the Terraform arm before writing it

The provider (`gravitee-io/apim`) is Speakeasy-generated from APIM's Automation
API OpenAPI spec, so its 10 resources map 1:1 onto the Automation API paths.
**The local checkout drifts: always read `origin/main`.**

```sh
cd ~/dev/src/terraform-provider-apim && git fetch origin
git show origin/main:internal/provider/provider.go | sed -n '/func.*Resources(ctx/,/^}/p'
git show origin/main:internal/provider/apiv4_resource.go | grep -oE '^\t\t\t"[a-z_0-9]+": ' | tr -d '\t":' | sort
```

See [PARITY.md → Regenerating](../../../PARITY.md#regenerating-this-document).

- **Attribute exists** → write the arm.
- **Attribute missing** → this is an **Automation API backlog item**, not an
  external constraint. File it, then declare the arm
  `pending: { terraform: "<reason + ticket>" }` so it renders as a visible
  `test.fixme`. Do not silently drop the journey to GKO-only.
- **V2 API** → permanently GKO-only; the Automation API is V4-only by design.

## Step 4: write one shared intent

```ts
forEachProvisioner<MyParams>(
  {
    title: "…",
    provisioners: { gko: gkoScenario<MyParams>({…}), terraform: tfScenario<MyParams>({…}) },
    xray: { gko: XRAY.AREA.GKO_ID, terraform: XRAY.TERRAFORM.TF_ID },
    tags: [TAGS.REGRESSION],
    since: { gko: "4.12", terraform: "4.12" },   // only if the feature has a minimum APIM version
    timeoutMs: { gko: 90_000 },                  // only if 30s is genuinely not enough for the work
  },
  async ({ provisioned, mapi, gateway }) => { … },
  initialParams,
);
```

Rules for the body:

- **Never branch on `provisionerId`.** If you need to, the difference belongs in
  the fixtures or in `params.ts`, or the assertion belongs in a `-gko-only` /
  `-tf-only` file.
- **Assert through `mapi` / `gateway` only.** They are the provisioner-agnostic
  platform assertions. Use `provisioned.view` only after `remove()` (expect
  `"gone"`) and on failure paths, never as a convergence wait after `update()`.
- **Assert the fields that carry regression value** (`visibility`,
  `lifecycleState`, `security.type`, the policy's effect at the gateway), not
  just that the resource exists. A use-case framing must not cost granularity.
- **Poll, never assert immediately after an apply.** `waitFor*` or
  `expect.poll(...).toMatchObject({...})`, combining checks atomically.
- **Wrap stages in `test.step`** so a failure reports which stage broke.

Fixture rules:

- Co-locate `gko/*.yaml` and `terraform/main.tf` in the journey folder; pass
  **absolute** paths built from `import.meta.url`.
- GKO CRs carry the `gravitee.io/e2e: "true"` label and reference the shared
  `dev-ctx` ManagementContext.
- Every `main.tf` exposes `output "api_id"`, `output "sub_id"`,
  `output "app_id"`, `output "group_id"` and `output "api_context_path"` as
  applicable, or passes an `outputs` map. A role used with
  `assertProvisioner` also needs `addresses: { role: "apim_x.name" }`.
- Give resources **unique, journey-scoped names**. API/App names are a shared
  global namespace across the whole serial suite; the GKO and TF arms must not
  collide with each other (the convention is a `-tf` suffix on the TF arm).
- **Structurally different parameterization goes in `params.ts`**: one shared
  param type plus the per-provisioner closures (GKO `applyParams`, TF `toVars`).
  Model: `api-consumer/subscribe-and-call/params.ts`.

## Step 5: tag both arms

Ids live only in [`e2e/helpers/tags.ts`](../../../e2e/helpers/tags.ts).

- Reusing an existing Xray test: use its constant, and delete the test it
  supersedes (Step 7).
- New coverage: add a `@GKO-TBD-<slug>` placeholder per arm under the right
  category, with the `// ── GKO-NNNN: <parent story>` comment convention.
- GKO and Terraform are **different tickets**, so each arm gets its own id. A
  list is allowed per arm when one provisioner splits into several tickets what
  the other does in one.
- Never put a Story or bug key in a test title, docstring, comment or fixture.

## Step 6: run it

```bash
npm run e2e -- --grep @GKO-xxxx                          # the GKO arm
npm run e2e -- --grep @GKO-yyyy --provision-with terraform # the TF arm
npm run build && npm run check:lanes                     # lane partition guard
npm run typecheck:e2e
```

Never report done without a run. If it fails, use the
`investigate-e2e-failure` skill. Do not raise the timeout and do not skip.

## Step 7: finish the job

- **`/xray-sync-tests`** to turn every `@GKO-TBD-*` into a real Jira Test ticket
  and rewrite `tags.ts` in place.
- **Delete superseded tests.** Once the journey covers what a standalone
  `tests/gko/` test did, remove that test and carry its Xray id onto the journey's
  GKO arm. Admission or rejection cases stay in `tests/gko/<area>/` rather than
  moving into the journey.
- **Update [the catalog](../../../e2e/tests/user-journeys/README.md)**: one row,
  with the journey link, what it demonstrates, and both arms' Xray ids.
- **Update [PARITY.md](../../../PARITY.md)**: move the area from "Migratable, not
  yet migrated" to "Migrated", or add a row to "Not expressible through
  Terraform" with the reason.
- **Write the journey `README.md`**: "As a `<persona>`, I …" plus how to run it.

## Checklist before reporting done

- [ ] Placed by persona; a new folder only because the story differs
- [ ] Terraform arm confirmed against `origin/main`, or `pending` with a filed reason
- [ ] Body does not branch on `provisionerId`; asserts through `mapi`/`gateway`
- [ ] Regression-carrying fields asserted, not just existence
- [ ] Every assertion after an apply is polled
- [ ] Both arms tagged; ids only in `tags.ts`; no story/bug keys anywhere
- [ ] Fixtures co-located, uniquely named, TF outputs contract honoured
- [ ] Both arms actually run and passed (paste the result)
- [ ] `npm run check:lanes` green
- [ ] `/xray-sync-tests` run; catalog and PARITY.md updated; superseded tests deleted
- [ ] Commit prefixed `test:`, body wrapped at ~72 cols, no AI attribution
