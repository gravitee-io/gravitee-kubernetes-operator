# E2E Tests

End-to-end tests for GKO and the Terraform APIM provider, built on Playwright Test.

## Prerequisites

- APIM, Gateway, and GKO operator running (e.g. via Kind)
- `kubectl` configured and pointing at the cluster
- `terraform` CLI installed (for Terraform tests)
- `config.yaml` in `test/platform-test/` with APIM/Gateway endpoints and credentials (see sample in repo)

## Running

```bash
cd test/platform-test

# Run all E2E tests
npm run e2e

# Run a single Xray test
npm run e2e -- --grep @GKO-110

# Run regression suite
npm run e2e:regression
```

## Folder structure

```
e2e/
  playwright.config.ts   # Test runner config (serial, 1 worker)
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
