# Use case: reuse a shared policy group across an API

Define a Shared Policy Group once and reuse it from a V4 API's request flow. The
SPG injects `X-SPG-Test`, which the echo endpoint reflects, proving the SPG runs
at the gateway. The same journey runs against both provisioners.

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/shared-policy-group.yaml`](./gko/shared-policy-group.yaml) + [`gko/api-with-spg.yaml`](./gko/api-with-spg.yaml) / [`gko/api-without-spg.yaml`](./gko/api-without-spg.yaml) | `SharedPolicyGroup` CR + API flow referencing it by HRID; detach = re-apply without the flow. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_shared_policy_group` + `apim_apiv4` flow; `attach_spg = false` drops the flow. |

**What it proves:** an API that reuses an SPG runs the SPG's policy at the gateway
(the injected header is reflected by the echo backend); detaching the SPG stops it.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-976
npm --prefix test/platform-test run e2e -- --grep @GKO-976 --provision-with gko
npm --prefix test/platform-test run e2e -- --grep @GKO-976 --provision-with terraform
```
