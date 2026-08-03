# Journey: secure an API with a plan

**As an API producer, I secure my API with a JWT plan and an OAuth2 plan.**

Provision a V4 proxy API secured with a JWT plan and an OAuth2 plan. The same
journey runs against both provisioners and asserts both plan security types land
in APIM.

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/api-with-plans.yaml`](./gko/api-with-plans.yaml) | `ApiV4Definition` with `JWTPlan` + `OAuth2Plan`. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_apiv4` with inline `plans[]` (JWT + OAuth2). |

**What it proves:** an API secured through either driver exposes a published plan
with `security.type = JWT` and one with `security.type = OAUTH2` in APIM. Gateway
enforcement with real tokens stays in the GKO subscription suites.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-162
npm --prefix test/platform-test run e2e -- --grep @GKO-162 --provision-with terraform
```
