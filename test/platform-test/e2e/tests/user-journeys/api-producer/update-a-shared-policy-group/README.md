# Journey: update a shared policy group

**As an API producer, I define a shared policy group once and change its policy
step as my needs change.**

| Driver | Fixture | Notes |
|---|---|---|
| GKO | [`gko/`](./gko/) | `SharedPolicyGroup` CR, two variants under one CR name so the second apply is an update. |
| Terraform | [`terraform/`](./terraform/) | `apim_shared_policy_group`. Only `hrid` forces replacement, so `steps` and `description` change in place. |

| Stage | Description | Step |
|---|---|---|
| as authored | `Shared policy group as first authored` | `transform-headers` injecting `X-SPG-Test: spg-header` |
| updated | `Shared policy group after its update` | step renamed, injecting `X-SPG-Test: spg-header-updated` |

**What it proves:** the whole declared group reaches APIM (api type, phase,
description and the step's own configuration), and an update rewrites *that same
group* rather than creating a second one or stalling as an undeployed draft.

Both stages assert the group's content, not just that it survived. The update
changes the step name and the injected header value, so an update that reaches
APIM and silently no-ops fails here. Id stability comes for free: the second
stage polls the **original** id for the **new** content, so a group replaced
instead of updated makes that id 404. `lifecycleState` is asserted per stage
because an update APIM accepts but never deploys is invisible to consumers.

Nothing is asserted at the gateway. Reusing a group inside an API flow, which is
what makes it observable on the data plane, belongs to
[`reuse-shared-policy-group`](../reuse-shared-policy-group/) and stays blocked on
GKO-3001 and the provider's missing `cross_id`.

Each arm uses its own group name (`update-spg` / `update-spg-tf`) so the two
never collide in the shared environment.

Run it:

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-981
npm --prefix test/platform-test run e2e -- --grep @GKO-3121   # terraform arm
```
