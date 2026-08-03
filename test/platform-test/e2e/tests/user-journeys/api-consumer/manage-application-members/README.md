# Journey: manage application members

**As an application developer, I share my application with a teammate, promote
them when they take over its subscriptions, and remove them when they leave.**

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/`](./gko/) | `Application.spec.members`, one variant per stage under a single CR name. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_application.members` + `notify_members`. |

**What it proves:** the application's member list in APIM after each change — the
whole list, so a role change that *added* a second entry instead of updating the
existing one is caught.

The member must be a real user or it is silently dropped, so the journey creates
a `gravitee`-source service account as a precondition, the same as
[`manage-api-members`](../../platform-admin/manage-api-members/).

**Not covered here** (GKO-only, in `tests/gko/applications`): a member naming a
non-existent user, group or role, a member missing `source`, and a member
declared *without* a role — the Automation API marks it non-nullable, so
Terraform rejects omitting it.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-538
npm --prefix test/platform-test run e2e -- --grep @GKO-538 --provision-with terraform
```
