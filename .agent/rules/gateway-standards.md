# Gateway API Standards

These rules apply to controllers in `controllers/gateway-api/`.

## Reconciliation Pipeline

Gateway API controllers follow a phased reconciliation chain:

```
Init -> Resolve -> Accept -> Program
```

- **Init**: reset status structures (rebuild listener statuses, seed conditions to `Pending`)
- **Resolve**: validate references (`ResolvedRefs` condition) — TLS secrets, parent gateways, backend Services, ReferenceGrants
- **Accept**: attachment and policy checks (`Accepted` condition) — listener compatibility, hostname intersection, namespace policy
- **Program**: materialize desired state — deploy workloads or create downstream CRs/ConfigMaps

Each phase sets its own condition. Do not skip phases or merge them.

## Status Conditions

- Always set both `Accepted` and `Programmed` conditions on Gateways
- Always set `Accepted` and `ResolvedRefs` per route parent in `RouteParentStatus`
- Use upstream Gateway API reason codes (`RouteReasonNoMatchingParent`, `ListenerReasonProtocolConflict`, etc.)
- Set `ObservedGeneration` on every condition update
- Use the decorator wrappers in `api/model/gateway/decorators.go` for condition manipulation

## Route Acceptance

- Build one `RouteParentStatus` entry per `spec.parentRefs[]` element
- Tag each with `controllerName: apim.gravitee.io/gateway`
- Validate: parent must be a Gateway, `sectionName`/`port` must match a listener, protocol must match, hostname must intersect
- If the parent Gateway status is not ready yet, return `ErrGatewayNotReady` and requeue
- Track `attachedRoutes` count on each listener status (sum of HTTP + Kafka routes)

## Conflict Resolution

- `conflict.go` detects protocol, hostname, and path conflicts across routes on the same listener
- Set the `Conflicted` condition on affected listeners and route parents
- Conflict detection uses gateway tags (`namespace/gateway-name`) on generated API specs

## Cross-Namespace References

- All cross-namespace references (backend Services, TLS secrets) require a `ReferenceGrant`
- Use `SupportsRouteNamespace()` to check `allowedRoutes` namespace policy on listeners
- Watch `ReferenceGrant` resources and map changes back to affected Gateways/Routes

## Route Programming

- Default mode: routes project into `ApiV4Definition` CRs (owned by the route via owner references)
- Alternate mode (`GATEWAY_API_SKIP_API_DEFINITION=true`): routes write gateway ConfigMaps directly
- Mappers in `internal/mapper/` translate Gateway API rules to the v4 API model
- Set `DefinitionContext{Origin: Kubernetes, SyncFrom: Kubernetes}` on generated API specs

## Conformance

- Changes to Gateway API controllers must pass `make conformance`
- The project targets `GatewayHTTPConformanceProfile` (core + extended)
- Conformance reports are versioned in `test/conformance/kubernetes.io/gateway-api/report/`
- Skipped tests must be documented with a reason (feature flag dependency or known limitation)
- When adding or modifying Gateway API features, you MUST run `make conformance` 3 consecutive times. All 3 runs must pass before the work can be considered complete. If any run fails, investigate and fix the issue, then restart the 3-pass cycle from scratch.
- A build passing, a linter passing and unit tests passing are **not** evidence about Gateway API behaviour. None of them executes a reconcile. Work that has not been through the suite is unverified and must be reported as unverified.
- See `.agent/skills/run-conformance.md` for how to bootstrap and run the suite locally.

## Verifying a fix for a flaky test

The three-pass rule above is the baseline. It is **not sufficient on its own for a test already
known to be non-deterministic**, and the arithmetic is why: if a test fails with probability p on
any given run, three consecutive passes happen by chance with probability (1-p)³. At p = 0.1 that
is 73%; at p = 0.3 it is 34%. Passing three times in a row is the most likely way to declare a
flaky bug fixed while it is still there.

So when the target is a proven flaky spot, the full suite is the integration check, not the
instrument. Required, in order:

1. **Build a targeted reproducer** — the smallest set of resources and requests that exercises the
   race, driven in a loop. Races of this kind usually need repeated polling to surface: a single
   request can look fine.
2. **Prove it reproduces before changing anything.** Measure the failure rate over N iterations on
   the unmodified code. If it will not fail on demand, it is not a reproducer and the spot is not
   understood yet — go back to evidence gathering rather than guessing at a fix.
3. **After the fix, bound the rate rather than counting a streak.** Require zero failures over
   enough iterations to put a real ceiling on p. By the rule of three, zero failures in n runs
   gives a 95% upper bound of about 3/n: n ≥ 30 bounds it under 10%, n ≥ 60 under 5%. Report the
   bound achieved, for example "0/50, p < 6% at 95%". Targeted runs take seconds, so this is cheap
   where the equivalent number of full-suite runs would not be.
4. **Then** the 3 consecutive full-suite runs, to catch interactions the targeted reproducer
   cannot. Any failure at any stage restarts from step 2.

A green full-suite run never clears a known flaky spot on its own. Keep a running tally of
consecutive greens per setup and state it next to the measured bound. Never report a flaky spot as
fixed without that bound: "3 runs passed" is not a result, "0 failures in 50 targeted iterations
plus 3 clean full runs" is.

Flakiness work must be driven by observed failures — a CI junit message, a cluster dump, an
operator or gateway log line — not by reading the code and reasoning about what could race. Code
reasoning generates candidates; only an observed failure justifies a change.
