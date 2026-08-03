# Journey: associate a group with an API

**As a platform admin, I hand an API to a team by associating it with that team's
group, and detach it when ownership moves.**

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/`](./gko/) | `Group` CR + `ApiV4Definition.spec.groups`. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_group` + `apim_apiv4.groups`. |

**What it proves:** the API record in APIM carries the group's id while
associated and nothing once detached.

The group is created by the same provisioner as the API and must reconcile
**first** — APIM silently drops a reference to a group that does not exist yet.
GKO orders this through `dynamicRoles`, Terraform through the
`apim_group.group.name` reference. That silent drop is also why the assertion
resolves the group's APIM id rather than trusting the apply.

A reference to a group that does *not* exist is operator behaviour (a `.status`
warning naming the missing group) and stays in `tests/gko/members`.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-257
npm --prefix test/platform-test run e2e -- --grep @GKO-257 --provision-with terraform
```
