# The provisioner contract

Annotated reference for the interfaces in `src/provisioners/`, plus the two
implementations to read as worked examples. Read the source alongside this; the
files are short and heavily commented.

| File | Holds |
|---|---|
| [`registry.ts`](../../../../src/provisioners/registry.ts) | `PROVISIONER_ORDER`, `PROVISIONER_LANES`, the `ProvisionerId` type |
| [`types.ts`](../../../../src/provisioners/types.ts) | `Provisioner`, `Provisioned`, `Role`, `ProvisionerChecks` |
| [`view.ts`](../../../../src/provisioners/view.ts) | `ProvisionerView`, `ProvisionerState`, `assertProvisioner` |
| [`base.ts`](../../../../src/provisioners/base.ts) | `BaseProvisioned`, the shared id getters |
| [`gko/`](../../../../src/provisioners/gko/) | `GkoProvisioner`, `GkoView`, `subscriptionYaml` |
| [`terraform/`](../../../../src/provisioners/terraform/) | `TerraformProvisioner`, `TerraformView`, `TerraformChecks` |
| [`engines/`](../../../../src/provisioners/engines/) | `kubectl.ts`, `terraform-core.ts` |

---

## `Role`

```ts
type Role = "api" | "application" | "subscription" | "plan" | (string & {});
```

A logical name inside a scenario, mapped per provisioner to a concrete resource.
The union gives autocomplete for the common roles while still allowing
scenario-specific keys. `BaseProvisioned` builds the internal role string as
`kind` or `kind:label`, so `apiId("two-plans")` resolves the role `api:two-plans`.

## `Provisioner<P>`

```ts
interface Provisioner<P = unknown> {
  readonly provisionerId: ProvisionerId;
  provision(params: P): Promise<Provisioned<P>>;
  cleanup?(): Promise<void>;
  attach?(refs: Partial<Record<Role, { hrid: string }>>): Promise<Provisioned<P>>;
}
```

- `provision` stands the scenario up with initial params and returns a live
  handle.
- `cleanup` is best-effort teardown **without** a handle, for the case where
  `provision()` itself failed partway (a CR applied but reconcile stuck). Omit it
  when the driver self-cleans on a failed provision, as Terraform does.
- `attach` is the upgrade-testing seam: rebuild a handle from stable HRIDs after
  the original in-memory state is gone, which is what lets
  `tests/upgrade/survival.before.spec.ts` and `survival.after.spec.ts` run as two
  separate processes either side of a `gck patch`. GKO implements it; Terraform
  does not.

## `Provisioned<P>`

```ts
interface Provisioned<P = unknown> {
  readonly provisionerId: ProvisionerId;

  apiId(label?: string): Promise<string>;
  subscriptionId(label?: string): Promise<string>;
  applicationId(label?: string): Promise<string>;
  groupId(label?: string): Promise<string>;
  contextPath(): Promise<string>;

  update(params: P): Promise<void>;
  remove(role: Role): Promise<void>;
  destroy(): Promise<void>;

  readonly view: ProvisionerView;
  readonly checks?: ProvisionerChecks;
}
```

Extend `BaseProvisioned` and implement one `resolveId(role)`; the four getters
come from it. Ids are resolved once, then cached.

Adding a new kind is **one line in `base.ts`**, not a method per provisioner:

```ts
dictionaryId(label?: string): Promise<string> {
  return this.resolveId(roleFor("dictionary", label));
}
```

## `ProvisionerView`

```ts
type ProvisionerState = "applied" | "failed" | "gone";

interface ProvisionerView<Detail = unknown> {
  read(role: Role): Promise<{ state: ProvisionerState; detail: Detail }>;
}
```

Three outcomes, no fourth. An "unknown" state would let a half-implemented
`read()` pass silently, so a provisioner that cannot yet tell **must throw**.

`assertProvisioner(provisioned, role, expected, options?)` polls `read()`,
retrying both the transient throw and a not-yet-`expected` state. Callable from
shared scenario bodies without narrowing.

Where it is used: after `remove()` (expect `"gone"`) and on failure paths
(expect `"failed"`, whose `detail` carries the driver's reason). Not as a
convergence wait after `update()`, because GKO's condition `observedGeneration`
can lag after a re-apply (GKO-2940).

## `ProvisionerChecks`

```ts
interface ProvisionerChecks {
  readonly provisionerId: ProvisionerId;
}
```

The base carries only the discriminant; a concrete shape (`TerraformChecks`)
extends it, and a type guard (`isTerraform`) narrows it. Deliberately rare:
anything expressible as "did MY layer land this role?" belongs in `view`.

A provisioner with no unique assertions **omits the property entirely**, so
shared bodies cannot reach for it by habit. GKO does exactly this: its
control-plane readouts are the `view` question, and the remaining Kubernetes
primitives (admission rejection, events) are reached through the kubectl engine
in tests where nothing is provisioned.

Terraform's checks today: drift detection, sensitive-value redaction, plan exit
codes, taint.

---

## Worked example: GKO

```ts
interface GkoScenarioSpec<P = unknown> {
  manifests: string[];              // absolute paths, applied in order at provision()
  roles: GkoRoles;                  // role -> CR name (string) or { kind, name }
  dynamicRoles?: Role[];            // roles created by applyParams; awaited after it
  contextPath?: string;             // GKO has no output to read it from
  applyParams?: (k: KubectlEngine, params: P) => Promise<void>;
  namespace?: string;
}
```

- **Kind by convention** from the role name: `api` → `apiv4definition`,
  `application` → `application`, `subscription` → `subscription`,
  `group` → `group`. Any other role must use the full `{ kind, name }` form, and
  the provisioner throws a named error if it cannot infer one.
- `resolveId(role)` reads `.status.id` of the role's CR.
- `view.read(role)` reads the CR's conditions.
- `remove(role)` deletes that CR.
- `attach()` is implemented: HRID is `namespace + "-" + name`, which is stable
  across an upgrade.

## Worked example: Terraform

```ts
interface TfScenarioSpec<P = unknown> {
  fixtureDir: string;                              // absolute, contains main.tf
  env: Record<string, string>;                     // APIM_* auth, built by the e2e adapter
  outputs?: Partial<Record<Role, string>>;         // role -> output name
  addresses?: Partial<Record<Role, string>>;       // role -> "apim_apiv4.api", for view.read()
  contextPathOutput?: string;                      // default "api_context_path"
  toVars?: (params: P) => Record<string, unknown>; // params -> tfvars, per apply
  removeVars?: Partial<Record<Role, Record<string, unknown>>>;  // how remove(role) drops it
}
```

- `resolveId(role)` reads `terraform output`. Defaults: `api` → `api_id`,
  `subscription` → `sub_id`, `application` → `app_id`, `group` → `group_id`.
- `addresses` has **no default on purpose**: resource local names are not
  standardized across fixtures (`test`, `app`, `dict`, `group`, `spg` are all in
  use), so guessing one would point `view.read()` at the wrong resource instead
  of failing loudly. It throws when a role needs `view` and has no address.
- `remove(role)` merges `removeVars[role]` over the current vars and re-applies,
  which is how a `count`-gated resource drops out of desired state. That is a
  realistic user action, not a state surgery.
- `cleanup()` is omitted: a failed apply self-cleans.

---

## The e2e binding

`src/` must not import from `e2e/`. The two things it therefore cannot own are
supplied by [`e2e/helpers/provisioner-env.ts`](../../../../e2e/helpers/provisioner-env.ts):

```ts
export function xScenario<P = unknown>(input: XScenarioInput<P>): () => Provisioner<P> {
  return () => new XProvisioner<P>({ ...input, manifests: input.manifests.map(resolveFixturePath) });
}
```

- **APIM environment** from `config.yaml` plus env-var overrides, memoized per
  run (see `tfEnv()`).
- **Fixture-path resolution**: absolute paths (co-located journey fixtures built
  from `import.meta.url`) pass through unchanged; relative paths are rooted at
  `e2e/fixtures/`.

## How a scenario becomes tests

[`e2e/helpers/for-each-provisioner.ts`](../../../../e2e/helpers/for-each-provisioner.ts)
expands one `ScenarioDef` into one Playwright test per provisioner:

- The title gets the scenario title, that arm's Xray id(s), `@<provisionerId>`,
  the extra tags, and `@since-<v>` when declared.
- A provisioner with a factory runs the shared body; one in `pending` renders as
  `test.fixme`; one absent from both emits nothing.
- The active provision is tracked at **module scope** so a single `afterEach` can
  tear it down even when the test **times out**, which an inline `finally` cannot
  do. Safe because the suite runs serially with one worker.
- Timeout defaults to Playwright's 30s, or `TF_WORKSPACE_TIMEOUT_MS` for the
  Terraform arm, overridable per provisioner via `timeoutMs`.
