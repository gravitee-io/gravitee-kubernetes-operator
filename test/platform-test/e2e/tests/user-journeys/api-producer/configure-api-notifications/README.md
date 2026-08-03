# Journey: configure console notifications on an API

**As an API producer, I have the console notify a team's group when my API starts
or stops, then widen that to expiring api-keys.**

The two drivers **model this differently**, which is exactly why it is worth one
shared journey:

| Driver | Fixture | Model |
|---|---|---|
| GKO | [`gko/`](./gko/) | a standalone `Notification` CR the API references via `notificationsRefs`. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | inline `apim_apiv4.console_notification`. |

**What it proves:** both land the same PORTAL notification setting — the same
hooks and the same target group — on the API in APIM, and widening the event list
updates it in place.

APIM only lets a console notification target groups the API itself belongs to, so
the group association is part of the setup here, not the assertion (that is
[`associate-groups-with-an-api`](../../platform-admin/associate-groups-with-an-api/)).

**Turning the notification off is not in this journey** (GKO-3085). The
Automation API answers HTTP 500 (`organizationId must not be empty`) to a console notification
with an empty `events` list, and setting the block to `null` is a Terraform no-op
that leaves the previous settings in place — so the Terraform arm cannot express
it. The GKO side removes the reference and is covered by `@GKO-1238` in
`tests/gko/notifications`.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-1231
npm --prefix test/platform-test run e2e -- --grep @GKO-1231 --provision-with terraform
```
