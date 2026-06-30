# Use case: publish an API to the portal and serve traffic

Provision a V4 proxy API, published to the portal (PUBLIC + PUBLISHED) and
STARTED, then take it on and off the gateway. The same journey runs against both
provisioners.

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/api-started.yaml`](./gko/api-started.yaml) + [`gko/api-stopped.yaml`](./gko/api-stopped.yaml) | `ApiV4Definition`, keyless plan; re-apply the stopped variant to stop. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_apiv4`; `state` re-applied between STARTED and STOPPED. |

**What it proves:** a started, published, public API is reachable at the gateway
(200) and reports `state STARTED`, `visibility PUBLIC`, `lifecycleState PUBLISHED`
in APIM; stopping it returns 404; re-starting serves traffic again.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-1464
npm --prefix test/platform-test run e2e -- --grep @GKO-1464 --provision-with gko
npm --prefix test/platform-test run e2e -- --grep @GKO-1464 --provision-with terraform
```
