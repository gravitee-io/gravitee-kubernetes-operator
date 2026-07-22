# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gravitee Kubernetes Operator (GKO) is a Kubernetes operator built with **Kubebuilder/controller-runtime** that manages Gravitee API Management (APIM) resources through Custom Resource Definitions (CRDs). It allows users to define, deploy, and publish APIs to the Gravitee API Portal and Gateway declaratively.

**Language:** Go 1.25.5 | **API Group:** `gravitee.io/v1alpha1` | **Module:** `github.com/gravitee-io/gravitee-kubernetes-operator`

## Build & Development Commands

```bash
# Build
make build                     # Build manager binary (runs code generation first)
make generate                  # Generate DeepCopy methods via controller-gen
make manifests                 # Generate CRD manifests into helm/gko/crds/gravitee.io

# Lint
make lint-fix                  # Auto-fix lint issues + add license headers
make add-license               # Add Apache 2.0 license headers to all Go files

# Test
make unit                      # Run unit tests (Ginkgo) — test/unit/...
make it                        # Run integration tests (Ginkgo, requires cluster) — test/integration/...
make e2e                       # Run e2e tests (Playwright) — test/platform-test/

# Run a single unit test suite
./bin/ginkgo test/unit/apim/...
# Run a single integration test file (use --focus to filter by description)
./bin/ginkgo --focus "should ..." test/integration/apidefinition/v2/...

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

Validation and mutation webhooks organized by resource type (`api/v2/`, `api/v4/`, `application/`, `mctx/`, `subscription/`, `group/`, `policygroups/`). Each has a `ctrl.go` (webhook handler) and `validate.go`.

### Drift detection

Drift detection rejects admission **updates** when APIM was changed outside the operator while the CRD still reflects the old desired state. It is enabled globally by default (`ENABLE_DRIFT_DETECTION`, Helm `manager.driftDetection.enabled`) and can be overridden per resource with the `gravitee.io/drift-detection` annotation (`true` / `false`).

Two packages are involved:

| Package | Role |
|---------|------|
| `internal/drift/` | Comparison engine: struct walk, equivalence registry, `Detect`, `Merge`, `Result.String` |
| `internal/admission/drift/` | Admission glue: template compile, ref resolution, remote fetch, DTO mapping |

`drift.Init()` must run at startup (`main.go`) and in unit/integration suites that exercise drift (`BeforeSuite`).

#### Adding drift detection to a resource

Hook into `validateUpdate` in `internal/admission/<resource>/validate.go`, after existing validations and before returning errors:

```go
errs.MergeWith(drift.ValidateDrift(ctx, oldObj, newObj, resolveRefs, getRemote, drift.MapDTO(toDTO)))
// or, when the APIM client comes from a related resource (e.g. subscription uses the application's context):
errs.MergeWith(drift.ValidateDriftWithContext(ctx, oldObj, newObj, resolveContext, resolveRefs, getRemote, dtoMapper))
```

Provide four callbacks:

1. **`RefResolver`** (`func(ctx, runtime.Object) error`) — resolve inlined references (Secrets, ConfigMaps, templates) on **both** old and new deep copies before comparison. Return a non-nil error to abort with a severe admission error (e.g. application client-certificate resolution in `resolveAppRefs`).
2. **`RemoteObjectGetter`** — fetch the live APIM object. Branch on `k8s.IsAutomationAPIManaged` (HRID + Automation API) vs legacy (UUID + Management API). Report fetch failures via `errs.AddSeveref`.
3. **`DTOMapper`** — map each CRD copy to the **same struct type** returned by the remote getter. Use `drift.MapDTO(func(cr *v1alpha1.MyResource) model.MyDTO { ... })` for type safety. The DTO must represent what is (or would be) sent to APIM, not the raw CRD spec.
4. **`ContextResolver`** (only for `ValidateDriftWithContext`) — when the CRD has no `ManagementContext` ref of its own but depends on a related resource's context (subscription → application).

Reuse dependencies already resolved in `validateUpdate` (API, application, plan, etc.) inside closures passed to `getRemote` / `dtoMapper` — do not resolve them again (nil-deref risk).

#### DTO design

Define comparison DTOs in `internal/apim/model/` (e.g. `ApplicationDTO`, `SubscriptionDTO`). Tag fields with `drift:"<equivalence>"` struct tags:

| Tag | Use for |
|-----|---------|
| `ignore` | Server-managed or identity fields not in the CRD payload (`id`, `hrid`, `status` on applications) |
| `empty-is-nil` | Optional slices, maps, pointers, zero-value structs |
| `trimmed` | Strings with insignificant whitespace |
| `rfc3339` | Date-time strings (timezone-tolerant) |
| `unstructured` | `GenericStringMap` / `unstructured.Unstructured` JSON blobs |

Fields without a tag use `reflect.DeepEqual`. Only tag fields that are part of the **spec payload**; if APIM returns a field the DTO mapper never sets and both sides end up empty, comparison is a no-op — explicit `ignore` is optional belt-and-suspenders.

Add drift tags on nested `api/model/` types when the same struct is embedded in the DTO (e.g. TLS certificate fields on `application.ClientCertificate`).

Optionally add `Spec.ToDTO()` on the v1alpha1 type when the mapping is stable and reused by the controller.

#### Merge semantics (why both old and new are compared)

`ValidateDriftWithContext` compares **old CRD → remote** and **new CRD → remote**, then `drift.Merge`:

- Remote-only change → **drift** (reject)
- CRD update that realigns with remote → **ok** (allow)
- Unchanged CRD, remote changed → **drift** (reject)

See `internal/drift/types.go` (`Merge` comment) and `internal/drift/doc.go`.

#### Testing

- **Unit** (`test/unit/drift/`): table-driven `drift.Detect` tests for DTO equivalence / tag behaviour. Do not re-test the framework; test your DTO tags and `ToDTO()` parity. Call `drift.Init()` in the suite `BeforeSuite`.
- **Integration** (`test/integration/admission/<resource>/`): apply fixtures, mutate APIM out-of-band, call `AdmissionCtrl.ValidateUpdate`, assert with `test/internal/integration/assert.DriftDetected`. Call `drift.Init()` in `SynchronizedBeforeSuite`. Use `labels.WithContext` when a `ManagementContext` is required.

Reference implementations: `internal/admission/application/validate.go`, `internal/admission/subscription/validate.go`, `test/unit/drift/application_detect_test.go`, `test/unit/drift/subscription_detect_test.go`.

### Internal Packages (internal/)

Key packages: `apim/` (APIM client logic), `core/` (shared interfaces), `env/` (config via env vars), `search/` (cache field indexers), `template/` (Go templating for CRD values), `watch/` (dynamic resource watching), `webhook/` (webhook server setup).

### Entry Point (main.go)

Initializes controller-runtime manager, registers all controllers and webhooks based on feature flags (`ENABLE_GATEWAY_API`, `ENABLE_INGRESS`, `ENABLE_WEBHOOK`, `ENABLE_TEMPLATING`), optionally applies CRDs from embedded Helm chart (`APPLY_CRDS`).

## Testing Patterns

- **Unit tests** (`test/unit/`): Ginkgo v2 suites. Dot-imports for `ginkgo/v2` and `gomega` are allowed.
- **Integration tests** (`test/integration/`): Ginkgo v2 suites requiring a running cluster. Use `test/internal/integration/fixture/` for building test fixtures from YAML files in `test/internal/integration/`. Use `test/internal/integration/constants/` for shared file paths and timeouts.
- **E2E tests** (`test/platform-test/e2e/`): Playwright (TypeScript) suites running against a real cluster with APIM and the operator. Fixtures live in `test/platform-test/e2e/fixtures/`. Run via `make e2e` or `npm --prefix test/platform-test run e2e`.
- **Helm tests** (`helm/gko/tests/`): helm-unittest YAML tests.

Integration test fixtures are YAML manifests loaded via the fixture builder pattern:
```go
fixture.Builder().
    AddSecret(constants.ContextSecretFile).
    Build().
    Apply()
```

## Code Generation

After modifying CRD types in `api/`:
1. `make generate` — regenerates `zz_generated.deepcopy.go` files
2. `make manifests` — regenerates CRD YAML in `helm/gko/crds/gravitee.io/`

Both steps are required when changing model or API type structs.

## Workflow
I work plan-first. I write structured prompts in `prompts/` and expect
a reviewed plan in `plans/` before any implementation.

- Never implement directly from a prompt file. Always produce a plan first.
- Plans are markdown files that I review and edit before implementation.
- When implementing, follow the plan strictly. Stop and ask if something
  doesn't match reality.

## Conventions

- **Commit style:** Conventional Commits (enforced by commitlint)
- **License headers:** Apache 2.0 on all `.go` files (enforced by `addlicense`, template in `LICENSE_TEMPLATE.txt`)
- **Linting:** Strict golangci-lint config (`.golangci.yml`) — max function length 100 lines, max cyclomatic complexity 30, strict error checking, no naked returns
- **Naming:** Lint excludes `Api/Url/Http` vs `API/URL/HTTP` casing warnings
- **Makefile:** Modular structure in `hack/make/*.mk`; tools installed to `./bin/`