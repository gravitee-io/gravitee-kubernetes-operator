# GKO - End-to-End Tests

This document explains how to run the E2E test suite for the Gravitee Kubernetes Operator.  
E2E tests use [Kyverno Chainsaw](https://github.com/kyverno/chainsaw) to exercise real cluster behavior (applying CRDs, deploying APIs, asserting runtime responses).

## 1. Prerequisites

Install and have in PATH:
- Docker
- kind
- kubectl
- helm
- Node.js (for zx scripts invoked by tests)
- Azure CLI (to log into the private registry)
- GNU make

A local kind cluster is used (created by provided make target).

## 2. Install Chainsaw binary

Install the Chainsaw test runner locally into the repository `bin/` directory:

```sh
GOBIN=$(pwd)/bin go install github.com/kyverno/chainsaw@v0.2.13
```



## 3. Start a local cluster with APIM 

Log in to the private Gravitee Azure registry and start the APIM stack inside Kind with the latest images.

```sh
az acr login --name graviteeio
APIM_IMAGE_REGISTRY=graviteeio.azurecr.io APIM_IMAGE_TAG=master-latest make start-cluster
```

This:
- Creates a kind cluster named `gravitee`
- Deploys APIM components using Helm values in the repository

## 4. Build and deploy the Operator via Helm

Build a local operator image, load it into the kind cluster, and install the Helm release.

```sh
IMG=gko TAG=dev make docker-build \
  && kind load docker-image gko:dev --name gravitee \
  && helm upgrade --install gko helm/gko \
    --set manager.image.repository=gko \
    --set manager.image.tag=dev \
    --set manager.metrics.enabled=false
```



## 5. Directory layout (simplified)

```
test/
    ├── e2e/
        ├── chainsaw/
            ├── commands/
            ├── tests/
                ├── apis/
                    ├── startApi/
                        ├── v2/
                        |   ├── chainsaw-test.yaml
                        ├── v4/
                            ├── chainsaw-test.yaml
```

Each test folder contains a `chainsaw-test.yaml` manifest describing test steps (create resources, run scripts, capture pod logs, etc.).

## 6. Run the full E2E suite

Runs every Chainsaw test.

```sh
make e2e
```

This target:
- Invokes Chainsaw against `test/e2e/chainsaw/tests/**/chainsaw-test.yaml`

## 7. Focus a single test

To iterate quickly on one test:
1. Add the metadata label in that test file:

```yaml
metadata:
  labels:
    focus: "true"
```

Example (excerpt):
```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: v2-start-stopped-api
  labels:
    focus: "true"
```

2. Run the focused suite:

```sh
make e2e-focus
```

The `e2e-focus` target filters tests by the label `focus=true`.  


## 8. Test authoring guidelines

- Prefer small, single-purpose test directories (one test per folder)
- Use descriptive test `metadata.name`
- Leverage `try` / `catch` blocks to capture events or pod logs for troubleshooting
- Keep helper scripts in `test/e2e/commands` (invoked via `npx zx ...`)
- Use environment bindings (`bindings:`) in Chainsaw to avoid hardcoding values


## 9. Clean up

```sh
kind delete cluster --name gravitee
```

Recreate via section 3 when needed.

## 10. Summary of core commands

```sh
# Install Chainsaw
GOBIN=$(pwd)/bin go install github.com/kyverno/chainsaw@v0.2.13

# Start cluster + APIM
APIM_IMAGE_REGISTRY=graviteeio.azurecr.io APIM_IMAGE_TAG=master-latest make start-cluster

# Build + deploy operator
IMG=gko TAG=dev make docker-build
kind load docker-image gko:dev --name gravitee
helm upgrade --install gko helm/gko \
  --set manager.image.repository=gko \
  --set manager.image.tag=dev \
  --set manager.metrics.enabled=false

# Run all e2e tests
make e2e

# Focused run (requires focus: "true" label in test files)
make e2e-focus
```
