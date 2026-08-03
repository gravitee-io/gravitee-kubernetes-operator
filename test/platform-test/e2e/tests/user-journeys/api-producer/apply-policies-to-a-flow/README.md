# Journey: apply a policy to a flow

**As an API producer, I add a transform-headers policy to my API's response
phase, change which header it adds, and remove it again.**

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/`](./gko/) | `ApiV4Definition.spec.flows`, one variant per stage under a single CR name. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | inline `apim_apiv4.flows`. |

**What it proves:** the flow reaches APIM **and** the gateway actually applies
it — the response carries the header while the policy is in place, carries the
new one (and not the old one) after a rewrite, and carries neither once the flow
is removed.

The gateway assertion is the point. APIM will happily record a flow the gateway
never runs, so a control-plane-only check proves nothing about the policy
executing; the definition is checked too, so a failure tells you whether the flow
never reached APIM or reached it and was not executed.

Plan-level policy validation (general conditions referencing a page that does not
exist) stays GKO-only under `tests/gko/policies`.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-94
npm --prefix test/platform-test run e2e -- --grep @GKO-94 --provision-with terraform
```
