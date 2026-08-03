# Journey: subscribe to a token-secured plan

**As an API consumer, I subscribe my application to an API's JWT plan and to its
OAuth2 plan.**

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/`](./gko/) | `ApiV4Definition` + `Application` + one `Subscription` per plan. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_apiv4` + `apim_application` + two `apim_subscription`. |

**What it proves:** both subscriptions reach `ACCEPTED` in APIM **even though
both plans declare `MANUAL` validation** — a subscription written through the
Automation API is auto-validated regardless, and a regression there would leave
every declarative subscription stuck `PENDING`.

Two fixture details are load-bearing:

- The API sets `allowMultiJwtOauth2Subscriptions`. Without it APIM rejects the
  second subscription with *"An other OAuth2 or JWT plan is already subscribed by
  the same application"*.
- The subscriptions are applied only **after** the API and application have
  reconciled. APIM's subscription endpoint rejects an application it still
  considers archived, which a single multi-document apply runs straight into.

## No gateway assertion here — and a product bug

Minting a token the gateway would accept needs a JWT/OAuth2 resource wired to a
real issuer, which the test cluster does not provide. The *no-token* path cannot
be asserted on this API either:

> **GKO-3086 — an OAuth2 plan with no `oauthResource` configured makes the
> gateway answer HTTP 500, not 401.** The gateway throws
> `NullPointerException: Cannot invoke "OAuth2PolicyConfiguration.getOauthResource()"
> because "this.oAuth2PolicyConfiguration" is null`
> (`io.gravitee.policy.oauth2.Oauth2Policy.getOauth2Resource`). An unconfigured
> OAuth2 plan should be rejected at design time or fail closed with a 401.

The 401 on a JWT-only API is `@GKO-817` in `tests/gko/subscriptions`.

**Not covered here** (GKO-only): immutability of an accepted `Subscription` CR,
admission on a mismatched plan, a `syncFrom=KUBERNETES` API, an unstarted API,
deleting a subscribed plan, and everything V2.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-800
npm --prefix test/platform-test run e2e -- --grep @GKO-800 --provision-with terraform
```
