# Skill: Add a New CRD

## When to Use

Use this skill when adding a new Custom Resource Definition to the operator.

## 1. Define the Type

Create or edit a file in `api/v1alpha1/`:

```go
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type MyResource struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    Spec   MyResourceSpec   `json:"spec,omitempty"`
    Status MyResourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type MyResourceList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []MyResource `json:"items"`
}
```

Put the domain model in `api/model/<resource>/` and embed it in the spec, the way the other
resources do — that model is what the APIM DTO is mapped from, and keeping it separate from the
Kubernetes wrapper is what makes `mapViaJSON` viable.

Add kubebuilder validation markers on spec fields as needed.

## 2. Register the Type

Add to the `SchemeBuilder` in `api/v1alpha1/package_markers.go`:

```go
func init() {
    SchemeBuilder.Register(&MyResource{}, &MyResourceList{})
}
```

## 3. Implement Core Interfaces

Implement the interfaces in `internal/core/interface.go` that apply:
- `Object` (all CRDs)
- `Spec` / `Status`
- `ContextAwareObject` (if the resource references a ManagementContext)

Resources synced through the Automation API also need `PopulateIDs(mCtx, automationAPIManaged)`
and `GetID()`; the controller and the webhook both call them before every APIM request.

## 4. Generate Code

```bash
make generate manifests reference
```

Verify the generated CRD appears in `helm/gko/crds/gravitee.io/` and the API reference was updated.

## 5. Create the APIM Client Layer and the Controller

Follow the `add-controller` skill. It covers the DTO and service in `internal/apim/{model,service}`
before the reconciler itself — that ordering matters, the controller is the thin part.

## 6. Create the Webhook (if needed)

Follow the `add-webhook` skill.

## 7. Add Search Indexers (if needed)

If other controllers need to look up this resource by field — or if this resource is referenced by
others and those references must trigger a re-reconcile — add indexers in `internal/search/` and
register them in `search.InitCache`. A `Watches(...)` without a matching indexer silently enqueues
nothing.

## 8. Tests

Unit tests in `test/unit/<area>/` for the pure logic. Everything requiring a cluster or a live APIM
goes to [`gravitee-io/gravitee-platform-e2e`](https://github.com/gravitee-io/gravitee-platform-e2e),
with fixtures under `apim/fixtures/<area>/`. Do not add anything to `test/integration/`.

## 9. Final Checks

```bash
make build
make -j4 lint-sources
make unit
```

Commit the generated files alongside your source changes.
