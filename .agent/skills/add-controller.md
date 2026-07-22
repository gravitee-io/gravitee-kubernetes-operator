# Skill: Add a New Controller

## When to Use

Use this skill when adding a new reconciler for an existing or new CRD.

## Steps

### 1. Create the Controller File

Create `controllers/<family>/<resource>/controller.go`:

```go
package resource

import (
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler struct {
    client.Client
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // Fetch the resource
    // Check DeletionTimestamp -> delete path
    // Otherwise -> update path
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&v1alpha1.MyResource{}).
        WithEventFilter(predicate.LastSpecHashPredicate{}).
        Complete(r)
}
```

### 2. Create the Internal Package

Create `controllers/<family>/<resource>/internal/` with:

- `update.go` — main reconciliation logic (create/update downstream resources)
- `delete.go` — finalizer-based cleanup logic
- `status.go` — status subresource updates

### 3. Register in main.go

Add the controller setup in `main.go` inside the appropriate feature-flag block:

```go
if err := (&resource.Reconciler{Client: mgr.GetClient()}).SetupWithManager(mgr); err != nil {
    setupLog.Error(err, "Unable to create controller", "controller", "MyResource")
    os.Exit(1)
}
```

### 4. Add Predicates and Watches

- Use `predicate.LastSpecHashPredicate` to skip no-op reconciliations
- Add `Watches()` calls for related resources (e.g. referenced Secrets, ManagementContexts)
- Use mapper functions from `internal/watch/` to map related resource events back to your CR

### 5. Add RBAC Markers

Add kubebuilder RBAC markers above the `Reconcile` method:

```go
// +kubebuilder:rbac:groups=gravitee.io,resources=myresources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gravitee.io,resources=myresources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gravitee.io,resources=myresources/finalizers,verbs=update
```

Then run `make manifests` to regenerate RBAC.

### 6. Create Integration Tests

Create test files in `test/integration/<resource>/`:
- Use the fixture builder pattern
- Test create, update, and delete paths
- Test error cases (missing references, invalid specs)

### 7. Final Checks

```bash
make generate && make manifests
make build
make -j4 lint-sources
make unit
```
