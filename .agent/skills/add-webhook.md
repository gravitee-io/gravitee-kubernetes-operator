# Skill: Add a New Admission Webhook

## When to Use

Use this skill when adding validation or defaulting webhooks for a CRD.

## Steps

### 1. Create the Webhook Controller

Create `internal/admission/<resource>/ctrl.go`:

```go
package resource

import (
    "context"

    "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type AdmissionCtrl struct{}

func (a *AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.MyResource{}).
        WithValidator(a).
        WithDefaulter(a).
        Complete()
}

// Implement admission.Validator[*v1alpha1.MyResource]
func (a *AdmissionCtrl) ValidateCreate(ctx context.Context, obj *v1alpha1.MyResource) (admission.Warnings, error) {
    return validate(ctx, obj)
}

func (a *AdmissionCtrl) ValidateUpdate(ctx context.Context, oldObj, newObj *v1alpha1.MyResource) (admission.Warnings, error) {
    return validate(ctx, newObj)
}

func (a *AdmissionCtrl) ValidateDelete(ctx context.Context, obj *v1alpha1.MyResource) (admission.Warnings, error) {
    return nil, nil
}

// Implement admission.Defaulter[*v1alpha1.MyResource]
func (a *AdmissionCtrl) Default(ctx context.Context, obj *v1alpha1.MyResource) error {
    return nil
}
```

### 2. Create the Validation Logic

Create `internal/admission/<resource>/validate.go`:

```go
package resource

import (
    "context"

    "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
    "sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func validate(ctx context.Context, obj *v1alpha1.MyResource) (admission.Warnings, error) {
    var errs []error
    // Add validation checks
    return nil, errors.Join(errs...)
}
```

### 3. Register in main.go

Add inside the `ENABLE_WEBHOOK` block:

```go
if err := (&resource.AdmissionCtrl{}).SetupWithManager(mgr); err != nil {
    setupLog.Error(err, "Unable to create webhook", "webhook", "MyResource")
    os.Exit(1)
}
```

### 4. Choose Validator, Defaulter, or Both

- Use `WithValidator()` only if you need validation without defaults
- Use `WithDefaulter()` only if you need mutation without validation
- Use both when the resource needs both paths
- Only implement the interfaces you register (compile will enforce this)

### 5. Create Admission Tests

Create tests in `test/integration/admission/<resource>/`:
- Test that valid objects pass
- Test that invalid objects are rejected with clear error messages
- Test defaulting behavior if applicable

### 6. Final Checks

```bash
make build
make -j4 lint-sources
make unit
```
