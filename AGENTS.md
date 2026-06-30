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
make e2e                       # Run e2e tests (Chainsaw) — test/e2e/chainsaw/
make e2e-focus                 # Run only e2e tests with label focus=true

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

### Internal Packages (internal/)

Key packages: `apim/` (APIM client logic), `core/` (shared interfaces), `env/` (config via env vars), `search/` (cache field indexers), `template/` (Go templating for CRD values), `watch/` (dynamic resource watching), `webhook/` (webhook server setup).

### Entry Point (main.go)

Initializes controller-runtime manager, registers all controllers and webhooks based on feature flags (`ENABLE_GATEWAY_API`, `ENABLE_INGRESS`, `ENABLE_WEBHOOK`, `ENABLE_TEMPLATING`), optionally applies CRDs from embedded Helm chart (`APPLY_CRDS`).

## Testing Patterns

- **Unit tests** (`test/unit/`): Ginkgo v2 suites. Dot-imports for `ginkgo/v2` and `gomega` are allowed.
- **Integration tests** (`test/integration/`): Ginkgo v2 suites requiring a running cluster. Use `test/internal/integration/fixture/` for building test fixtures from YAML files in `test/internal/integration/`. Use `test/internal/integration/constants/` for shared file paths and timeouts.
- **E2E tests** (`test/e2e/`): Declarative YAML-based tests using Chainsaw (Kyverno). Test specs are `chainsaw-test.yaml` files.
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