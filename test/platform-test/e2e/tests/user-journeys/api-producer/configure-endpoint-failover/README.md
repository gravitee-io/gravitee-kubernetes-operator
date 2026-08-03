# Journey: configure endpoint failover

**As an API producer, I put a failover policy on my API so a slow or failing
backend is retried instead of surfacing as an error.**

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/`](./gko/) | `ApiV4Definition.spec.failover`, two retry budgets under one CR name. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_apiv4.failover`. |

**What it proves:** every field of the failover block round-trips to APIM (a
nested block is exactly where defaults get silently substituted), the API with
failover still serves traffic, and tightening the retry budget updates the policy
in place.

Exercising a *real* endpoint failure needs a backend the suite can take down on
demand, which the test cluster does not provide; the traffic assertion here only
proves the policy does not break the happy path.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-859
npm --prefix test/platform-test run e2e -- --grep @GKO-859 --provision-with terraform
```
