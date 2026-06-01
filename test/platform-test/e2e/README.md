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
| `kind` | recent | local Kubernetes cluster |
| `helm` | recent | installs the GKO chart |
| `terraform` | 1.12.1 | matches CI; Terraform tests assume it on `PATH` |
| Docker | recent | required by Kind and by `sew` |

The suite also expects:

- An APIM + Gateway stack reachable from the host (default URLs
  `http://localhost:30083` / `http://localhost:30082`).
- The GKO operator running **as a Deployment** in the cluster (any namespace).
- The GKO CRDs installed (notably `apiv4definitions.gravitee.io`).

These three are validated by [`global-setup.ts`](global-setup.ts) before any
test runs — see [Pre-flight checks](#pre-flight-checks).

## Bring up the cluster + APIM + GKO

Two supported paths. Pick one.

### Option A — CI-blessed path (recommended)

Mirrors what CI runs in `.circleci/config.yml` (`job-e2e-tests`):

```bash
# From the repo root
make start-cluster                       # Kind + APIM CE via hack/scripts/run-kind.mjs

# Build the operator image, load it into Kind, install via Helm
IMG=gko TAG=latest make docker-build \
  && kind load docker-image gko:latest --name gravitee \
  && helm upgrade --install gko helm/gko -n default \
    --set manager.image.repository=gko \
    --set manager.image.tag=latest

kubectl rollout status deployment/gko-controller-manager -n default --timeout=120s
```

### Option B — `sew` (one-shot stack composer)

[`sew`](https://a-cordier.github.io/sew/) is a third-party Kubernetes stack
composer that brings up APIM + GKO in a single command from its registry.
It is **not** part of this repo — you author your own `sew.yaml` and run
`sew create`.

Install (macOS):

```bash
brew install sew
```

Pick a starter composition for E2E:

```yaml
# sew.yaml — minimal: APIM (DB-less) + GKO
from:
  - gravitee-io/oss/apim/dbless
```

```yaml
# sew.yaml — full stack: APIM + MongoDB + Elasticsearch + GKO (closer to CI)
from:
  - gravitee-io/oss/apim/mongodb
```

Then:

```bash
sew info --from gravitee-io/oss/apim/mongodb   # preview the composition
sew create                                     # bring up the stack
sew delete                                     # tear it down
```

`sew` installs GKO into the `gravitee` namespace, while Option A uses
`default`. The suite probes all namespaces (`kubectl get deploy -A`), so both
work. EE registry variants (`gravitee-io/ee/apim/*`) exist but require a
licence file and are out of scope here — see the
[sew registry](https://a-cordier.github.io/sew/registry/).

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
  fixtures/
    crds/                # Kubernetes CRD manifests
    terraform/           # Terraform .tf files
  tests/
    gko/                 # GKO operator tests
    terraform/           # Terraform provider tests
```
