# Skill: Add a New CRD

## When to Use

Use this skill when adding a new Custom Resource Definition to the operator.

## Steps

### 1. Define the Type

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

Add kubebuilder validation markers on spec fields as needed.

### 2. Register the Type

Add to the `SchemeBuilder` in `api/v1alpha1/package_markers.go`:

```go
func init() {
    SchemeBuilder.Register(&MyResource{}, &MyResourceList{})
}
```

### 3. Implement Core Interfaces

Implement the interfaces in `internal/core/interface.go` that apply:
- `Object` (all CRDs)
- `Spec` / `Status`
- `ContextAwareObject` (if the resource references a ManagementContext)

### 4. Generate Code

```bash
make generate    # DeepCopy methods
make manifests   # CRD YAML in crds/gravitee.io/
```

Verify the generated CRD appears in `crds/gravitee.io/`.

### 5. Create the Controller

Follow the `add-controller` skill.

### 6. Create the Webhook (if needed)

Follow the `add-webhook` skill.

### 7. Add Integration Test Fixtures

Create YAML fixtures in `test/internal/integration/`:
- A minimal valid CR
- Variants for different test scenarios

### 8. Add Search Indexers (if needed)

If other controllers need to look up this resource by field, add indexers in `internal/search/`.

### 9. Final Checks

```bash
make build
make -j4 lint-sources
make unit
```

Commit the generated files alongside your source changes.
