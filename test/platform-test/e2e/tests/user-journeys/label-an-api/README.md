# Journey: label an API

**As an API producer, I organise my API with portal labels.**

Set portal labels on a V4 API, then strip them. Labels are an **inline attribute**
of the API (`apim_apiv4.labels`) — there is no standalone resource, but the
journey is still fully expressible through both provisioners.

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/api-with-labels.yaml`](./gko/api-with-labels.yaml) + [`gko/api-without-labels.yaml`](./gko/api-without-labels.yaml) | `ApiV4Definition.spec.labels`; strip = re-apply without them. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_apiv4.labels`; `with_labels = false` empties the list. |

**What it proves:** labels set through either driver land in APIM; stripping them
removes them. This is the pattern for the other inline `apim_apiv4` attributes
(`categories`, `groups`, `metadata`, inline `pages[]`) that have no standalone
Terraform resource but are still parity-testable at the API level.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-1473
npm --prefix test/platform-test run e2e -- --grep @GKO-1473 --provision-with terraform
```
