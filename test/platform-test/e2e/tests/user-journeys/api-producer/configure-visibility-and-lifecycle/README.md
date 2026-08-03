# Journey: configure portal visibility and lifecycle state

**As an API producer, I decide who can discover my API in the developer portal,
then retire it.**

`visibility` controls *who* sees the API in the portal; `lifecycleState` controls
*whether it is listed there at all*. Neither takes the API off the gateway —
only `state` does, which is
[`publish-api-and-serve-traffic`](../publish-api-and-serve-traffic/).

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/`](./gko/) | Three `ApiV4Definition` variants under one CR name. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_apiv4.visibility` / `.lifecycle_state`, count-gated for `remove()`. |

**What it proves:** the three combinations the portal distinguishes
(PRIVATE+PUBLISHED, PUBLIC+PUBLISHED, PUBLIC+UNPUBLISHED) each land on the APIM
record, changing **in place** on re-apply; an unlisted API still serves traffic;
and deleting the API removes it from both APIM and the gateway.

The portal's own rendering is a downstream consumer of these two fields, so the
journey asserts them at the source rather than through a browser.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-172
npm --prefix test/platform-test run e2e -- --grep @GKO-172 --provision-with terraform
```
