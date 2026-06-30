# Use case: reuse a shared policy group across an API

Define a Shared Policy Group once and reuse it from a V4 API's request flow. The
SPG injects `X-SPG-Test`, which the echo endpoint reflects, proving the SPG runs
at the gateway.

> ⚠️ **Currently blocked — both arms are `pending`.** This documents the correct
> form and the two blockers; it does not run green until they are fixed.

| Driver | Fixture | Status |
|---|---|---|
| GKO | [`gko/shared-policy-group.yaml`](./gko/shared-policy-group.yaml) + [`gko/api-with-spg.yaml`](./gko/api-with-spg.yaml) (uses the correct `sharedPolicyGroupRef`) | ⛔ admission rejects the documented form — **GKO-3001** (the operator resolves the ref to the SPG crossId in the reconciler, but not before the admission dry-run). |
| Terraform | [`terraform/main.tf`](./terraform/main.tf) | ⛔ `apim_shared_policy_group` exposes only `id`, not the crossId; only the crossId executes the SPG at the gateway. |

**What it will prove (once unblocked):** an API that reuses an SPG runs the SPG's
policy at the gateway (the injected header is reflected by the echo backend);
detaching the SPG stops it.

The crossId works at the gateway today (verified manually); the blockers are
purely about wiring the reference. The GKO arm uses the correct `sharedPolicyGroupRef`
form already, so it will run as soon as GKO-3001 is fixed.

Run it (currently shows as skipped/pending):

```sh
npm --prefix test/platform-test run e2e -- --grep @GKO-976
```
