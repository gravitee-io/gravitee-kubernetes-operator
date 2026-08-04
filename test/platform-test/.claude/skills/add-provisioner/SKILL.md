---
name: add-provisioner
description: Add a new provisioner (a creation path such as the Automation API over HTTP, the APIM console UI, or a new IaC driver) or a new engine to the Gravitee platform-test framework, so existing user journeys gain another arm. Use when asked to support a new way of creating APIM resources in the e2e suite, add a lane, extend the provisioner layer, or make journeys run through an additional driver.
allowed-tools: Bash, Read, Edit, Write, Glob, Grep
---

# add-provisioner

The provisioner layer is what lets one journey assert the same platform outcome
regardless of how the resources were created. Adding a provisioner should make
every existing journey capable of a new arm without touching a single assertion.

Invariants: [`.agent/rules/e2e-framework-development.md`](../../../.agent/rules/e2e-framework-development.md).
Annotated contracts and the two worked implementations:
[`references/provisioner-contract.md`](references/provisioner-contract.md).

Work from `test/platform-test/`.

---

## Step 1: register it

[`src/provisioners/registry.ts`](../../../src/provisioners/registry.ts) is the
single source of truth. Add the id and its lane tag:

```ts
export const PROVISIONER_ORDER = ["gko", "terraform", "<id>"] as const;

export const PROVISIONER_LANES: readonly ProvisionerLane[] = [
  { id: "gko", tag: "@gko" },
  { id: "terraform", tag: "@terraform" },
  { id: "<id>", tag: "@<id>" },
];
```

Everything else follows for free: the `ProvisionerId` union,
`forEachProvisioner`'s generation order, `--provision-with <id>` in
`scripts/e2e.mjs`, and the lane `grepInvert` in `e2e/playwright.config.ts`. **Do
not** add a hardcoded list anywhere else. If you find yourself wanting to, the
thing you are editing should derive from the registry instead.

Choose the tag lowercase and distinct from the uppercase `@GKO-NNNN` Xray
prefix; also add it to `PROVISIONER` in
[`e2e/helpers/tags.ts`](../../../e2e/helpers/tags.ts) so single-provisioner
tests can carry it.

## Step 2: build the engine (if there is one)

`src/provisioners/engines/` holds the mechanical drivers (`kubectl.ts`,
`terraform-core.ts`): process invocation, workspace lifecycle, error surfacing.
Keep this layer free of scenario semantics. If the new provisioner is a plain
HTTP client, reuse `src/utils/http/` rather than adding a dependency.

## Step 3: implement `Provisioner`

```ts
export class XProvisioner<P = unknown> implements Provisioner<P> {
  readonly provisionerId = "<id>";

  async provision(params: P): Promise<Provisioned<P>> { … }

  /** Best-effort teardown WITHOUT a live handle, for a provision() that failed
   *  partway. Idempotent and tolerant. Omit if the driver self-cleans. */
  async cleanup(): Promise<void> { … }

  /** Optional seam for upgrade testing: rebuild a handle from stable HRIDs
   *  after the in-memory state is gone. GKO implements it; Terraform does not. */
  async attach(refs): Promise<Provisioned<P>> { … }
}
```

The scenario spec it takes should mirror `GkoScenarioSpec` / `TfScenarioSpec`:
a role map, the parameterized apply hook, and whatever the driver needs to find
its own resources.

## Step 4: implement `Provisioned` on `BaseProvisioned`

Extend [`BaseProvisioned`](../../../src/provisioners/base.ts) so
`apiId()` / `applicationId()` / `subscriptionId()` / `groupId()` all come from
one `resolveId(role)`. Do not reimplement the getters.

Then implement `contextPath()`, `update(params)`, `remove(role)` and
`destroy()`.

- **`remove(role)` removes that role the way a user would**, leaving the rest of
  the scenario standing (GKO deletes that CR; Terraform drops the resource from
  desired state and re-applies). It is what lets a journey assert partial
  teardown.
- **`destroy()` must be idempotent and must never throw.** It runs in `finally`
  and in the `afterEach` safety net; a teardown failure must not mask a result.

## Step 5: implement `ProvisionerView`

```ts
async read(role: Role): Promise<ProvisionerViewResult<XViewDetail>>
```

Exactly three outcomes: `applied`, `failed`, `gone`. **There is no `unknown`.**
A provisioner that cannot yet tell what happened must **throw**, so
`assertProvisioner` retries it instead of a half-implemented `read()` passing
silently.

`detail` carries the driver's own evidence: a condition, a plan diff, an HTTP
error body.

## Step 6: `checks`, only if unavoidable

`ProvisionerChecks` is the narrowed escape hatch for assertions only this
provisioner can make. Before adding one, ask whether it is really "did MY layer
land this role?", which belongs in `view` and needs no narrowing. GKO declares
none, deliberately. If you do add one, ship the type guard alongside it
(`isTerraform` is the model) and export both from
`src/provisioners/index.ts`.

## Step 7: bind it for the e2e suite

`src/` cannot import from `e2e/`, so the environment and fixture-path resolution
are supplied by [`e2e/helpers/provisioner-env.ts`](../../../e2e/helpers/provisioner-env.ts).
Add an `xScenario()` factory there, taking fixture-relative paths and returning
`() => Provisioner<P>`, matching `gkoScenario` / `tfScenario`.

Scenario authors must never construct a provisioner directly.

## Step 8: verify

```bash
npm run build
npm run typecheck && npm run typecheck:e2e && npm run typecheck:examples
npm test
npm run check:lanes          # needs the build; the lane count must grow and still partition
npm run e2e -- --provision-with <id>
```

`check:lanes` is the guard that matters here: it asserts every test runs in
exactly one lane. A new lane whose tag is wrong shows up as tests in two lanes
or none, and both still pass an untagged full run.

## Step 9: wire CI and docs

- Add the lane to the matrix in `.circleci/config.yml` next to `gko` and
  `terraform`.
- [`PARITY.md`](../../../PARITY.md): what the new driver can and cannot express,
  and the resulting backlog.
- [The catalog](../../../e2e/tests/user-journeys/README.md): the new arm's column.
- [`AGENTS.md`](../../../AGENTS.md) and [`e2e/README.md`](../../../e2e/README.md):
  the lane in the tree and in the run commands.

## Step 10: adopt it journey by journey

Do not try to light up every journey at once.

- Add the arm to one journey, prove it end to end, then fan out.
- Every journey without the arm gets `pending: { <id>: "<reason>" }`, which
  renders as a visible `test.fixme`. A silent absence is a gap nobody sees; a
  provisioner absent from both `provisioners` and `pending` reads as "N/A by
  design", so only leave it out when that is true.
- Each arm carries **its own Xray id**: add `@GKO-TBD-*` placeholders and run
  `/xray-sync-tests`.

## Checklist

- [ ] Registered in `registry.ts` only; no hardcoded provisioner list added anywhere
- [ ] Lane tag added to `PROVISIONER` in `tags.ts`
- [ ] `Provisioned` extends `BaseProvisioned`; getters not reimplemented
- [ ] `destroy()` idempotent and non-throwing; `cleanup()` tolerant
- [ ] `view.read()` returns three states and throws when it cannot tell
- [ ] `checks` added only for something `view` genuinely cannot express
- [ ] `xScenario()` factory in `provisioner-env.ts`; `src/` still imports nothing from `e2e/`
- [ ] Apache 2.0 header on every new file; `exports` map updated if a new import path
- [ ] Full verification gate green, including a real `--provision-with <id>` run
- [ ] CI matrix, PARITY.md, catalog, AGENTS.md, e2e/README.md updated
- [ ] Journeys either have the arm or an explicit `pending` reason
