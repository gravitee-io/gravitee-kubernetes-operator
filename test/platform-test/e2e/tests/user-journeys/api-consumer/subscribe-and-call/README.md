# Journey: subscribe an application to an API and call it

**As an application developer, I subscribe my application to an api-key plan and call the API through the gateway.**

The api-key subscription lifecycle (auto-generated, custom, expiry, rotation,
reactivation, staggered expiry, revoke-on-delete, keyless coexistence) is
exercised against both provisioners from one shared intent.

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/api.yaml`](./gko/api.yaml) + [`gko/application.yaml`](./gko/application.yaml) (+ [`gko/api-mixed.yaml`](./gko/api-mixed.yaml) for the keyless-coexistence case) | API-key-plan API + application; the Subscription is built at run time. |
| Terraform | [`terraform/apikey-auto`](./terraform/apikey-auto), [`terraform/apikey-custom`](./terraform/apikey-custom), [`terraform/apikey-mixed`](./terraform/apikey-mixed) | `apim_apiv4` + `apim_application` + `apim_subscription` with an `api_keys` list. |

Provisioner-specific behaviour lives in `apikey-gko-only.test.ts` (admission,
templating) and `apikey-tf-only.test.ts` (drift, Sensitive redaction, plan exit
codes).

**What it proves:** a subscription's api-key(s) authenticate calls at the gateway,
and key rotation/expiry/revocation behave identically regardless of the driver.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-2825
npm --prefix test/platform-test run e2e -- --grep @GKO-2879 --provision-with terraform
```
