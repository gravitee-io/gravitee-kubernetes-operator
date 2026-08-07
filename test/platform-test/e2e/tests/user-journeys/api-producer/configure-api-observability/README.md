# Journey: configure API observability

**As an API producer, I turn on analytics for my API so its traffic shows up in
the console, ship its logs and traces to OpenTelemetry, and later turn tracing
verbose to debug a problem.**

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/`](./gko/) | `ApiV4Definition.spec.analytics`, two variants per scenario under one CR name. |
| Terraform | [`terraform/proxy/`](./terraform/proxy/) · [`terraform/native/`](./terraform/native/) | `apim_apiv4.analytics`. |

**What it proves:** the whole analytics block round-trips to APIM (a nested block
is exactly where defaults get silently substituted), and changing one field
updates it in place rather than replacing the API.

Two scenarios, not one variant table, because the product splits the surface and
an API cannot cross the split with an update:

| Scenario | API type | Fields |
|---|---|---|
| Report traffic with analytics, OTel logs and tracing | `PROXY`, STARTED | `enabled`, `otelLogs.enabled`, `tracing.{enabled,verbose}` |
| Disable the connection-metrics reporter | `NATIVE` Kafka, STOPPED | `enabled`, `reporterMetricsEnabled` |

`tracing` and `logging` are documented "not for native APIs";
`reporterMetricsEnabled` applies only to native ones. The reporter is asserted
**off** first because the provider defaults it to `true`, so `false` is the value
that proves the declaration was transmitted rather than defaulted.

The native fixtures carry a listener host and broker port range unique to each
driver: APIM admission enforces global uniqueness on both, so sharing them across
arms would make the suite order-dependent.

Nothing is asserted at the gateway. Analytics is control-plane configuration, and
observing it on the data plane needs an OTel collector and a log sink the test
cluster does not run. `analytics.sampling` and `analytics.logging` are not
covered yet.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-3117
npm --prefix test/platform-test run e2e -- --grep @GKO-3118   # terraform arm
npm --prefix test/platform-test run e2e -- --grep @GKO-3119   # native reporter, gko
npm --prefix test/platform-test run e2e -- --grep @GKO-3120   # native reporter, terraform
```
