# Runner examples

These are **illustrations, not part of the test suite**. They exist to
demonstrate one property of `@gravitee/platform-test`: it is **runner-agnostic**.
The library depends only on `node:*`, `yaml` and native `fetch`, so the same
assertions and the same provisioner layer work under any test runner.

The real suite lives in [`../e2e/`](../e2e/) and runs on Playwright. If you are
writing GKO tests, read [`../AGENTS.md`](../AGENTS.md) and work there — not here.

```
runner-examples/
  playwright/   the same tests under @playwright/test
  jest/         the same tests under jest + ts-jest
```

Each folder holds the *same* set of tests, so the diff between them is only the
runner's hook and assertion names:

| Test | Shows |
|---|---|
| `provision-and-assert` | the current authoring model: drive a `Provisioner`, read its own record with `assertProvisioner`, then assert the platform with `mapi` / `gateway` |
| `start-stop-api` | gateway + mAPI assertions against an API you already have |
| `one-category` | partial-match assertions on an API record |
| `subscribe-to-jwt-plan` | gateway assertions with a bearer token |

They import from `../../dist/…` because they live inside this repo; a real
consumer would import from `@gravitee/platform-test`. Build the library first
(`npm run build` in the package root) or the imports will not resolve.

## Keeping them honest

They are **not** run by CI (they need a live cluster), so the only thing standing
between them and silent rot is a type check. The Playwright example resolves
`@playwright/test` from the package root's `node_modules`, so it checks with no
extra install:

```bash
npm run typecheck:examples          # from test/platform-test/
```

The Jest example needs its own dependencies (`ts-jest`, `@types/jest`), so check
it after installing them:

```bash
npm --prefix runner-examples/jest install
npx tsc --noEmit -p runner-examples/jest
```

## Running one

Both need APIM + Gateway reachable and `config.yaml` (or the `GRAVITEE_*` env
vars) configured — see [`../README.md`](../README.md#configuration). The
provisioning example additionally needs a cluster with the GKO operator:

```bash
MANIFEST=/abs/path/to/api.yaml API_NAME=my-api API_PATH=/my-api \
  npx playwright test provision-and-assert    # from runner-examples/playwright/
```
