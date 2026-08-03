# Journey: create a group with a member

**As a platform admin, I create a group so I can organise API members.**

A group created through either provisioner lands in APIM via the Automation API.

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/group.yaml`](./gko/group.yaml) | `Group` CR with one member. |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | `apim_group`. |

**What it proves:** a group created through either driver is recorded in APIM with
`origin: KUBERNETES`. Member reconciliation / drift / import stay in the
per-provisioner suites (`tests/gko/groups`, `tests/terraform/groups.test.ts`).

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-983
npm --prefix test/platform-test run e2e -- --grep @GKO-983 --provision-with terraform
```
