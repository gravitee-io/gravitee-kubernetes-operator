# AI Agent Context

@.agent/rules/go-standards.md
@.agent/rules/operator-standards.md
@.agent/rules/gateway-standards.md

## Project Overview

Gravitee Kubernetes Operator (GKO) is a Kubernetes operator built with **Kubebuilder/controller-runtime** that manages Gravitee API Management (APIM) resources through Custom Resource Definitions (CRDs). It allows users to define, deploy, and publish APIs to the Gravitee API Portal and Gateway declaratively.

**Language:** Go 1.26 | **API Group:** `gravitee.io/v1alpha1` | **Module:** `github.com/gravitee-io/gravitee-kubernetes-operator`

## Build & Development Commands

```bash
# Build
make build                     # Build manager binary (runs code generation first)
make generate                  # Generate DeepCopy methods via controller-gen
make manifests                 # Generate CRD manifests into helm/gko/crds/gravitee.io

# Lint
make lint-fix                  # Auto-fix lint issues + add license headers
make add-license               # Add Apache 2.0 license headers to all Go files
make lint-commits              # Lint HEAD commit message (commitlint)
# Lint every commit on the branch (same rules as CI `job-lint-commits`)
npx --yes -p @commitlint/cli -p @commitlint/config-conventional \
  commitlint --from origin/master

# Test
make unit                      # Run unit tests (Ginkgo) — test/unit/...
make it                        # Run integration tests (Ginkgo, requires cluster) — test/integration/...
make e2e                       # Run e2e tests (Playwright) — test/platform-test/

# Run a single unit test suite
go tool ginkgo test/unit/apim/...
# Run a single integration test file (use --focus to filter by description)
go tool ginkgo --focus "should ..." test/integration/apidefinition/v2/...

# Local development
make start-cluster             # Create local KinD cluster with APIM
make delete-cluster            # Delete local KinD cluster
make install                   # Install CRDs into current cluster
make run                       # Run operator locally (APPLY_CRDS=true ENABLE_GATEWAY_API=false)
```

## Architecture

### CRDs (api/v1alpha1/)

All CRDs belong to the `gravitee.io` API group, version `v1alpha1`:

- **APIM resources:** `ApiDefinition` (v2), `ApiV4Definition` (v4), `ManagementContext` (cluster-scoped), `ApiResource`, `Application`, `Subscription`, `Group`, `Notification`, `SharedPolicyGroup`
- **Gateway API resources:** `GatewayClassParameters`, plus standard `HTTPRoute`/`KafkaRoute`

Type definitions live in `api/v1alpha1/`, data models in `api/model/`. Core interfaces that all CRD types implement are in `internal/core/interface.go` (`Object`, `Spec`, `Status`, `ContextAwareObject`, etc.).

### Controllers (controllers/)

Two controller families, each under `controllers/`:
- **`apim/`** — 9 controllers (apidefinition v2/v4, apiresource, application, group, ingress, managementcontext, notification, sharedpolicygroups, subscription)
- **`gateway-api/`** — 5 controllers (gateway, gatewayclass, gatewayclassparameters, httproute, kafkaroute)

Each controller follows the standard Kubebuilder reconciler pattern:
- `*_controller.go` — `Reconciler` struct with `Reconcile()` and `SetupWithManager()`
- `internal/` subpackage — `update.go`, `delete.go`, `status.go` for reconciliation logic

Controllers use a **watch system** (`internal/watch/`) to react to changes in related resources (contexts, resources, groups, notifications). The `predicate.LastSpecHashPredicate` prevents reconciliation when the spec hasn't changed.

### Admission Webhooks (internal/admission/)

Validation and mutation webhooks organized by resource type (`api/v2/`, `api/v4/`, `application/`, `mctx/`, `subscription/`, `group/`, `policygroups/`). Each has a `ctrl.go` (webhook handler) and `validate.go`. Controllers implement generic `admission.Validator[T]` / `admission.Defaulter[T]`; private `validateCreate` / `validateUpdate` / `validateDelete` take the concrete CRD type (e.g. `*v1alpha1.Application`) so they do not type-assert from `runtime.Object`.

### APIM client (internal/apim/)

Everything that talks to APIM goes through `internal/apim`. This is the layer that makes a CRD real, and it is where a new resource needs the most work.

| Package | Role |
|---------|------|
| `internal/apim/apim.go` | The `APIM` facade: one field per service, built per reconcile from a `ManagementContext` |
| `internal/apim/client/` | `Client` (HTTP + `URLs`) and the target builders `OrgTarget`, `EnvV1Target`, `EnvV2Target`, `AutomationTarget` |
| `internal/apim/service/` | One file per resource: the HTTP calls (`CreateOrUpdate`, `DryRunCreateOrUpdate`, `Delete`, `GetByHRID`, `GetByID`) |
| `internal/apim/model/` | The **wire payloads** (`*DTO`) sent to and received from APIM, plus their `To*DTO` mappers |

Controllers and webhooks never build a URL or a payload themselves. They obtain a client with `apim.FromContextRef(ctx, obj.ContextRef(), obj.GetNamespace())` and call a service method.

#### Management API vs Automation API

`client.NewURLs` builds four roots from the `ManagementContext`:

| Target | Shape | Used by |
|--------|-------|---------|
| `AutomationTarget(path)` | `{base}/automation/organizations/{org}/environments/{env}{path}` | All HRID-managed resources (v4 APIs, applications, subscriptions, groups, shared policy groups, dictionaries, portals, documentations) |
| `EnvV1Target(path)` | `{base}/management/organizations/{org}/environments/{env}{path}` | Legacy v2 API import, some read-only lookups |
| `EnvV2Target(path)` | `{base}/management/v2/organizations/{org}/environments/{env}{path}` | Management API v2 reads |
| `OrgTarget(path)` | `{base}/management/organizations/{org}{path}` | Org-level calls |

Cloud contexts swap the base paths (`/apim/automation`, `/apim/rest`); a `spec.path` override on the context wins over both. **New resources are Automation API only** — do not add Management API endpoints for them.

The Automation API is HRID-addressed: `PUT` on the collection with the payload carrying `hrid`, `DELETE`/`GET` on `{collection}/{hrid}`, and `?dryRun=true|false` on writes. HRIDs come from `refs.NewNamespacedNameFromObject(obj).HRID()` (namespace + name).

Resources created before APIM 4.12 are still UUID-addressed. The service branches on the `AutomationAPIManaged` condition, sends the UUID in the HRID slot and sets `?hridContainsUUID=true`; on the first successful HRID-based write it calls `k8s.AddAutomationAPIManagedCondition(obj)` so subsequent reconciles use the HRID. Every service that supports the migration carries the same `get<Resource>ID` helper returning `(identifier, hridContainsUUID)`.

#### DTOs are the API payload

`internal/apim/model/` types are the **Automation API request and response bodies**, not comparison helpers. `SharedPolicyGroup.createOrUpdate` builds `model.ToSharePolicyGroupDTO(...)` and `PUT`s it; `GetByHRID` unmarshals the response into the same `SharedPolicyGroupDTO`. Drift detection is a *second consumer* of these types — that is the only reason `drift:` struct tags live on them.

Mappers are plain functions in `internal/apim/model/` (`ToApplicationDTO`, `ToAPIV4DTO`, `ToSharePolicyGroupDTO`, `ToGroupDTO`, `ToDictionaryDTO`, `ToPortalDTO`, …), not methods on `api/model/` types. Most delegate to the shared `mapViaJSON[T]` JSON round-trip; write explicit mapping only where the wire shape genuinely differs from the CRD (`ToAPIV4DTO` turns the spec's plan and page maps into sorted slices; `SubscriptionDTO.ToAutomation()` / `AutomationSubscriptionDTO.ToLegacy()` bridge two field namings).

Two deviations worth knowing before copying a pattern:

- **Application** `PUT`s `app.Spec` directly; `ApplicationDTO` exists for the typed `GET` response and for drift.
- **Subscription** `PUT`s `AutomationSubscriptionDTO` (`apiHrid`/`applicationHrid`/`planHrid`) but compares on `SubscriptionDTO` (`apiId`/`applicationId`/`planId`).

### Drift detection

Drift detection rejects admission **updates** when APIM was changed outside the operator while the CRD still reflects the old desired state. It is disabled globally by default (`DRIFT_DETECTION_ENABLED`, Helm `manager.driftDetection.enabled`) and can be overridden per resource with the `gravitee.io/drift-detection` annotation (`true` / `false`). Policy env vars: `DRIFT_DETECTION_POLICY`, `DRIFT_DETECTION_ON_REMOTE_MISSING`, `DRIFT_DETECTION_ON_FETCH_FAILURE` (`deny` | `warn` | `allow`).

Two packages are involved:

| Package | Role |
|---------|------|
| `internal/drift/` | Comparison engine: struct walk, equivalence registry, `Detect`, `Merge`, `Result.String` |
| `internal/admission/drift/` | Admission glue: template compile, ref resolution, remote fetch, DTO mapping |

`drift.Init()` must run at startup (`main.go`) and in any unit suite that exercises drift (`BeforeSuite`).

#### Adding drift detection to a resource

Hook into `validateUpdate` in `internal/admission/<resource>/validate.go`, after existing validations and before returning errors:

```go
errs.MergeWith(drift.ValidateDrift(ctx, oldObj, newObj, resolveRefs, getRemote, drift.MapDTO(toDTO)))
// or, when the APIM client comes from a related resource (e.g. subscription uses the application's context):
errs.MergeWith(drift.ValidateDriftWithContext(ctx, oldObj, newObj, resolveContext, resolveRefs, getRemote, dtoMapper))
```

Provide four callbacks:

1. **`RefResolver`** (`func(ctx context.Context, obj *v1alpha1.MyResource) error`) — resolve inlined references (Secrets, ConfigMaps, templates) on **both** old and new deep copies before comparison. Return a non-nil error to abort with a severe admission error (e.g. application client-certificate resolution in `resolveAppRefs`).
2. **`RemoteObjectGetter`** (`func(*apim.APIM, *v1alpha1.MyResource) (any, error)`) — fetch the live APIM object. Branch on `k8s.IsAutomationAPIManaged` (HRID + Automation API) vs legacy (UUID + Management API). Return the APIM client error as-is (including HTTP 404). `ValidateDriftWithContext` uses `errors.IsNotFound` to apply `OnRemoteMissing` vs `OnFetchFailure`.
3. **`DTOMapper`** — map each CRD copy to the **same struct type** returned by the remote getter. Use `drift.MapDTO(func(cr *v1alpha1.MyResource) model.MyDTO { ... })` for type safety. The DTO must represent what is (or would be) sent to APIM, not the raw CRD spec.
4. **`ContextResolver`** (only for `ValidateDriftWithContext`) — when the CRD has no `ManagementContext` ref of its own but depends on a related resource's context (subscription → application).

Reuse dependencies already resolved in `validateUpdate` (API, application, plan, etc.) inside closures passed to `getRemote` / `dtoMapper` — do not resolve them again (nil-deref risk).

#### Drift tags on the DTO

Do **not** define a DTO for drift. Reuse the resource's existing `internal/apim/model/` payload — the one the service already sends and receives (see [APIM client](#apim-client-internalapim)) — and add `drift:"<equivalence>"` struct tags to it. Drift tags are inert for JSON serialization, so tagging a live payload type is safe.

| Tag | Use for |
|-----|---------|
| `ignore` | Server-managed or identity fields not in the CRD payload (`id`, `hrid`, `crossId`, `status`) |
| `empty-is-nil` | Optional slices, maps, pointers, zero-value structs |
| `empty-is-true` | Booleans APIM defaults to `true` when absent |
| `trimmed` | Strings with insignificant whitespace |
| `rfc3339` | Date-time strings (timezone-tolerant) |
| `case-insensitive` | Enums APIM may echo back in a different case |
| `unstructured` | `GenericStringMap` / `unstructured.Unstructured` JSON blobs |
| `ignore-remote:A,B` | Strings where the listed remote values are server defaults |
| `ignore-namespace-prefix` | Strings APIM prefixes with the namespace |

Fields without a tag use `reflect.DeepEqual`. Only tag fields that are part of the **spec payload**; if APIM returns a field the mapper never sets and both sides end up empty, comparison is a no-op — explicit `ignore` is optional belt-and-suspenders.

Add drift tags on nested `api/model/` types when the same struct is embedded in the DTO (e.g. TLS certificate fields on `application.ClientCertificate`).

The equivalence registry is split in two: generic tags in `internal/drift/equivalences.go`, APIM-specific ones (`ignore-remote-only-metadata`, `ignore-unknown-crd-groups`) in `internal/apim/drift/equivalences.go`. Both are wired by `drift.Init()`.

#### Merge semantics (why both old and new are compared)

`ValidateDriftWithContext` compares **old CRD → remote** and **new CRD → remote**, then `drift.Merge`:

- Remote-only change → **drift** (reject)
- CRD update that realigns with remote → **ok** (allow)
- Unchanged CRD, remote changed → **drift** (reject)

See `internal/drift/types.go` (`Merge` comment) and `internal/drift/doc.go`.

#### Testing

Unit only, in `test/unit/drift/apim/`: table-driven `drift.Detect` tests over your DTO, covering each tag you added and `To*DTO` parity. Do not re-test the framework. Call `drift.Init()` in the suite `BeforeSuite`. Behaviour against a live APIM belongs in the e2e repo (see [Testing](#testing)).

Reference implementations: `internal/admission/application/drift.go`, `internal/admission/subscription/drift.go`, `test/unit/drift/apim/`.

### Internal Packages (internal/)

Key packages: `apim/` (APIM client — see [APIM client](#apim-client-internalapim)), `core/` (shared interfaces), `env/` (config via env vars), `k8s/dynamic/` (unstructured resolution of referenced CRs, Secrets and ConfigMaps), `search/` (cache field indexers), `template/` (Go templating for CRD values — delimiters are ``[[`` / ``]]``, not the Go default, with `secret` and `configmap` functions: ``[[ secret `my-secret/token` ]]``), `watch/` (dynamic resource watching), `webhook/` (webhook server setup).

### Entry Point (main.go)

Initializes controller-runtime manager, registers all controllers and webhooks based on feature flags (`ENABLE_GATEWAY_API`, `ENABLE_INGRESS`, `ENABLE_WEBHOOK`, `ENABLE_TEMPLATING`), optionally applies CRDs from embedded Helm chart (`APPLY_CRDS`).

## Testing

**New work is unit tests here and e2e tests in the platform repo. Do not add integration tests.**

| Layer | Where | What belongs there |
|-------|-------|--------------------|
| Unit | `test/unit/<area>/` in this repo | Pure logic: DTO mapping, drift tags, validation predicates, templating, helpers. Ginkgo v2; dot-imports for `ginkgo/v2` and `gomega` are allowed |
| E2E | [`gravitee-io/gravitee-platform-e2e`](https://github.com/gravitee-io/gravitee-platform-e2e) | Anything requiring a cluster or a live APIM: reconciliation, `.status`, admission rejection, drift, deletion |
| Helm | `helm/gko/tests/` | helm-unittest YAML tests |

In the e2e repo, operator-specific coverage goes in `apim/tests/gko/<area>/` with fixtures in `apim/fixtures/<area>/`; behaviour a customer could also reach through Terraform goes in `apim/tests/user-journeys/<persona>/<journey>/`. That repo carries its own `AGENTS.md` and a `write-e2e-test` skill — follow those, do not infer its conventions from this file.

`test/integration/` and `test/platform-test/` still exist and still run in CI. Keep them green, but do not extend them.

## Code Generation

After modifying CRD types in `api/`:
```bash
make generate manifests reference
```
- `generate` — DeepCopy methods (`zz_generated.deepcopy.go`)
- `manifests` — CRD YAML in `helm/gko/crds/gravitee.io/`
- `reference` — API docs in `docs/api/reference.md`

After modifying Helm values (`helm/gko/values.yaml`):
```bash
make helm-reference
```
Regenerates `helm/gko/README.md`.

All four targets are checked in CI.

## Workflow
I work plan-first. I write structured prompts in `prompts/` and expect
a reviewed plan in `plans/` before any implementation.

- Never implement directly from a prompt file. Always produce a plan first.
- Plans are markdown files that I review and edit before implementation.
- When implementing, follow the plan strictly. Stop and ask if something
  doesn't match reality.

## Conventions

- **Commit style:** Conventional Commits (enforced by commitlint via `commitlint.config.js`). Body lines must be ≤ 120 characters (`body-max-line-length`). After rewriting history, lint the whole branch before finishing: `npx --yes -p @commitlint/cli -p @commitlint/config-conventional commitlint --from origin/master`
- **License headers:** Apache 2.0 on all `.go` files (enforced by `addlicense`, template in `LICENSE_TEMPLATE.txt`)
- **Linting:** `go vet` + `revive` + `staticcheck` (run in parallel via `make -j4 lint-sources`). Config in `.revive.toml`. Max cyclomatic complexity 30. Coding standards in `.agent/rules/`
- **Naming:** Lint excludes `Api/Url/Http` vs `API/URL/HTTP` casing warnings
- **Makefile:** Modular structure in `hack/make/*.mk`; Go tools managed via `go.mod` `tool` directive (invoked as `go tool <name>`)