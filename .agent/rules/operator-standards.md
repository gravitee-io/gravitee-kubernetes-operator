# Operator Standards

## Reconciliation

- Reconcile must be idempotent: safe to run any number of times for the same input
- Level-triggered: react to current state, not to the event that triggered reconciliation
- No side effects on the read path (fetching current state must not mutate anything)
- Return a non-nil error for exponential backoff requeue
- Use `ctrl.Result{RequeueAfter: d}` for known transient waits (e.g. waiting for a LoadBalancer IP)
- Never block inside Reconcile; return and requeue instead

## Status

- Always update via the status subresource (`Status().Update()` or `Status().Patch()`), never a full object update
- Set `ObservedGeneration` in status conditions to track which spec generation was last reconciled
- Use `metav1.Condition` with `Type`, `Status`, `Reason`, `Message`, and `ObservedGeneration`
- Do not update status if nothing changed (avoids unnecessary watch events)

## Finalizers

- Use finalizers only when external resources (APIM, cloud) need cleanup on CR deletion
- Always check `DeletionTimestamp` before adding a finalizer to avoid racing with deletion
- Remove the finalizer only after all cleanup is confirmed complete
- Keep finalizer names namespaced: `core.gravitee.io/<purpose>`

## Logging

- Use structured key-value logging via `log.FromContext(ctx)`
- Start messages with a capital letter: `"Reconciling ApiDefinition"`
- Use past tense for completed actions: `"Created Deployment"`
- Always include the resource kind and namespaced name in log context
- Never use `fmt.Sprintf` for log messages; pass key-value pairs instead

## Watches

- Use `Watches()` or `WatchesRawSource()` in `SetupWithManager()` for cross-resource triggers
- Keep mapper functions (event → reconcile request) focused and side-effect-free
- Use `predicate.LastSpecHashPredicate` to skip reconciliation when spec has not changed
- Avoid watching resources that change frequently but rarely affect your controller

## Generated Files

- Never hand-edit `zz_generated*.go` files or CRD YAML in `crds/`
- Always commit regenerated output alongside the source change that caused it
- Run `make generate && make manifests` after any change to types in `api/`
