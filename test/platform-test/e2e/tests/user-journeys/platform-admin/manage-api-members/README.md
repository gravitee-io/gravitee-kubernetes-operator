# Journey: manage API members

**As a platform admin, I grant a teammate access to an API, adjust their role,
and revoke it when they move on.**

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/`](./gko/) | `ApiV4Definition.spec.members`, one variant per stage under a single CR name. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_apiv4.members` + `notify_members`. |

**What it proves:** the API's member list in APIM after each change — the whole
list, not just the member under test, so a role change that *added* a second
entry instead of updating the existing one is caught.

The member must be a real user: a member naming a user that does not exist is
silently dropped and APIM records only the primary owner, so the journey creates
a `gravitee`-source service account as a precondition (same as
[`create-group-with-member`](../create-group-with-member/)).

**Not covered here** (GKO-only, in `tests/gko/members`): a member naming a
non-existent user or group, `PRIMARY_OWNER` declared in `members`, primary owner
resolved from the ManagementContext user, a member missing `source`, and a member
declared *without* a role — the Automation API's member schema marks `role`
non-nullable, so the Terraform provider rejects omitting it.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-247
npm --prefix test/platform-test run e2e -- --grep @GKO-247 --provision-with terraform
```
