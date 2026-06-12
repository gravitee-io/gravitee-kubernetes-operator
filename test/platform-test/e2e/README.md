# E2E Tests

End-to-end tests for the **Gravitee Kubernetes Operator** and the **Terraform
APIM provider**, built on Playwright Test. The suite drives a real local
Kubernetes cluster running APIM, Gateway, and the GKO operator — there are no
mocks.

This doc takes you from a fresh checkout to a green `npm run e2e`.

> Looking to run the GKO operator locally for development (not E2E)? See
> [CONTRIBUTING.adoc](../../../CONTRIBUTING.adoc) — that path runs the operator
> on your host and does **not** satisfy the in-cluster operator probe used by
> this suite.

## Prerequisites

Host tools (install before bootstrapping the cluster):

| Tool | Version | Notes |
|------|---------|-------|
| Node.js | ≥ 18 | per `package.json#engines` |
| `kubectl` | recent | tests shell out to it (`helpers/kubectl.ts`) |
| `gck` | recent | provisions the Kind cluster + APIM stack ([docs](https://gravitee-io-labs.github.io/gck/)) |
| `kind` | recent | `kind load` of the locally built operator image |
| `helm` | recent | installs the GKO chart |
| `terraform` | 1.12.1 | matches CI; Terraform tests assume it on `PATH` |
| Docker | recent | required by Kind and by `gck` |

The suite also expects:

- An APIM + Gateway stack reachable from the host (default URLs
  `http://localhost:30083` / `http://localhost:30082`).
- The GKO operator running **as a Deployment** in the cluster (any namespace).
- The GKO CRDs installed (notably `apiv4definitions.gravitee.io`).

These three are validated by [`global-setup.ts`](global-setup.ts) before any
test runs — see [Pre-flight checks](#pre-flight-checks).

## Bring up the cluster + APIM + GKO

The cluster + APIM stack is provisioned with
[`gck`](https://gravitee-io-labs.github.io/gck/) (Gravitee Cluster Kit) from
the checked-in [`gck.yaml`](../gck.yaml): it composes the registry context
`gravitee-io/oss/apim/mongodb` with suite-specific overrides (APIM master
nightlies from azurecr, gateway sync across all namespaces, coverage mount,
Elasticsearch and portal disabled).

Install gck (same version CI pins in `.circleci/config.yml`):

```bash
go install github.com/gravitee-io-labs/gck@v1.0.2
```

Mirrors what CI runs in `.circleci/config.yml` (`job-e2e-tests`):

```bash
# From the repo root. The image preload pulls APIM master nightlies:
docker login graviteeio.azurecr.io

# Kind cluster `gravitee` + APIM + MongoDB into namespace `gravitee`.
# Errors if the cluster already exists — `make delete-cluster` first.
make start-e2e-cluster

# Build the operator image, load it into Kind, install via Helm
IMG=gko TAG=latest make docker-build \
  && kind load docker-image gko:latest --name gravitee \
  && helm upgrade --install gko helm/gko -n default \
    --set manager.image.repository=gko \
    --set manager.image.tag=latest

kubectl rollout status deployment/gko-controller-manager -n default --timeout=120s
```

Useful gck commands:

```bash
gck info --config test/platform-test/gck.yaml   # preview composition + flags
gck list                                        # show gck-managed clusters
make delete-cluster                             # tear the cluster down
```

APIM lives in the `gravitee` namespace (in-cluster management API URL:
`http://apim-api.gravitee.svc:83`, hardcoded in the management-context
fixtures); the operator lives in `default`. The suite probes all namespaces
(`kubectl get deploy -A`), so both work. The integration-test (Ginkgo) suite
and `make start-cluster` still use the legacy
`hack/scripts/run-kind.mjs` flow, which installs the `apim3` chart into
`default` — the two cluster flavors share the name `gravitee`, so delete one
before creating the other.

## Build platform-test

```bash
cd test/platform-test
npm install
npm run build
```

## Configuration

The suite reads [`config.yaml`](../config.yaml) from `test/platform-test/`. It
is committed and defaults match the local cluster brought up by Option A or B,
so no edits are needed for the common case.

Environment variables override fields in `config.yaml`:

| Variable | Overrides |
|----------|-----------|
| `GRAVITEE_BASE_URL` | `apim.baseUrl` |
| `GRAVITEE_ENV_ID` | `apim.envId` |
| `GRAVITEE_USERNAME` | `apim.auth.username` |
| `GRAVITEE_PASSWORD` | `apim.auth.password` |
| `GRAVITEE_GATEWAY_URL` | `gateway.baseUrl` |
| `GRAVITEE_GATEWAY_MTLS_URL` | `gateway.mtlsBaseUrl` |

## Terraform provider

Terraform tests under [`tests/terraform/`](tests/terraform/) use the
registry-published `gravitee-io/apim` provider by default - no extra setup.

To exercise unreleased provider code, build the provider from source into a
local mirror that `helpers/terraform.ts` auto-detects:

```bash
bash test/platform-test/scripts/build-tf-provider.sh main
```


## Run the suite

```bash
cd test/platform-test

# All E2E tests
npm run e2e

# Regression suite only
npm run e2e:regression

# Single test by Xray tag
npm run e2e -- --grep @GKO-110
```

Reports land in `test/platform-test/playwright-results/` (JUnit XML) and
`test/platform-test/playwright-report/` (HTML).

## Pre-flight checks

[`global-setup.ts`](global-setup.ts) runs five checks before any test. If any
fail, the suite aborts with a clear message. Mapping of symptom → fix:

| Failure | What it means | Fix |
|---------|---------------|-----|
| Management API unreachable | `GRAVITEE_BASE_URL` doesn't respond | Verify APIM is running; check `kubectl get pods -A` and port-forward / NodePort mapping |
| Gateway unreachable | `GRAVITEE_GATEWAY_URL` doesn't respond | Verify the gateway pod is `Ready`; check the gateway service NodePort |
| `kubectl cluster-info` fails | No reachable cluster context | `kubectl config use-context kind-gravitee` (or your cluster's context) |
| `apiv4definitions.gravitee.io` CRD missing | GKO CRDs not installed | Re-run `helm upgrade --install gko helm/gko ...` (Option A) or `sew create` (Option B) |
| No `app.kubernetes.io/name=gko` deployment | Operator not running in-cluster | Same as above — the operator must be a Deployment, not `go run` on the host |

## Folder structure

```
e2e/
  playwright.config.ts   # Test runner config (serial, 1 worker, 30s timeout)
  global-setup.ts        # Pre-flight infra checks (APIM, Gateway, K8s, GKO)
  setup.ts               # Playwright fixtures (mapi, gateway, kubectl)
  helpers/
    kubectl.ts           # kubectl CLI wrapper
    terraform.ts         # Terraform workspace lifecycle & commands
    tags.ts              # Xray test ID constants
  fixtures/              # one folder per domain; see "Fixture convention" below
    admission-webhook/
    api-definitions/
    api-lifecycle/
    api-v4-definitions/
    applications/
    categories/
    ...
  tests/
    gko/                 # GKO operator tests
    terraform/           # Terraform provider tests
```

## Fixture convention

Every fixture lives under `fixtures/<domain>/<scenario>/`. A scenario directory contains:
- `crd.yaml` — Kubernetes CRD manifest(s) for the GKO-driven test (multi-doc YAML when multiple resources are part of the starting state)
- `main.tf` — Terraform configuration for the TF-provider-driven test

A scenario with both files is **paired** (same APIM behaviour exercised through both drivers). A scenario with only one file is single-driver — and the gap is visible at `ls` time.

**Naming**: domain folders mirror test folders under `tests/gko/` (e.g. `admission-webhook/`, `api-lifecycle/`, `categories/`, `policies/`). Scenario folder names describe *what's being tested*, not *what kind of CR sits at the top of the manifest* — e.g. a V4 API with a JWT plan goes under `plans/v4-jwt/`, not `api-v4-definitions/`. Inside domains that hold both V2 and V4 variants (`plans/`, `categories/`, `policies/`, `api-lifecycle/`, `admission-webhook/`), prefix scenario names with `v2-` / `v4-` to disambiguate.

The slimmed `api-definitions/` and `api-v4-definitions/` domains hold only scenarios that test the bare API CR itself (minimal shape, default-field behaviour, etc.) — anything that tests plans, lifecycle, categories, message-API entrypoints, etc. lives in the corresponding domain.
