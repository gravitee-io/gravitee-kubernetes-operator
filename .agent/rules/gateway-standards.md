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
