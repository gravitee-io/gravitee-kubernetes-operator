// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v4

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apiv4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

var _ admission.Validator[*v1alpha1.ApiV4Definition] = AdmissionCtrl{}
var _ admission.Defaulter[*v1alpha1.ApiV4Definition] = AdmissionCtrl{}

type Lifecycle = lifecycle.AdmissionLifecycle[*v1alpha1.ApiV4Definition, apiv4.DTO, *apiv4.Client]

type AdmissionCtrl struct {
	Lifecycle Lifecycle
}

func NewAdmissionCtrl() AdmissionCtrl {
	return AdmissionCtrl{
		Lifecycle: lifecycle.NewAdmissionLifecycle(Lifecycle{
			ResolveRefs:    apiv4.ResolveRefsForAdmission,
			ClientFactory:  admissionClient,
			PreCheck:       preCheck,
			ToDTO:          apiv4.ToDTO,
			DryRun:         apiv4.DryRun,
			PayloadChanged: everyUpdate,
			GetRemote:      apiv4.GetRemote,
			DriftDTO:       apiv4.DriftDTO,
			DeleteCheck:    deleteCheck,
			// The v4 admission messages are the ones the webhook reported before the lifecycle.
			UnprefixedResolveErrors: true,
		}),
	}
}

// SetupWithManager keeps the no-op defaulter: Helm registers a mutating webhook for v4 APIs.
func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	a = NewAdmissionCtrl()
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.ApiV4Definition{}).
		WithValidator(a).
		WithDefaulter(a).
		Complete()
}

func (a AdmissionCtrl) Default(_ context.Context, _ *v1alpha1.ApiV4Definition) error {
	return nil
}

func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context,
	obj *v1alpha1.ApiV4Definition,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateCreate(ctx, obj).Map()
}

func (a AdmissionCtrl) ValidateDelete(
	ctx context.Context,
	obj *v1alpha1.ApiV4Definition,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateDelete(ctx, obj).Map()
}

func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context,
	oldObj *v1alpha1.ApiV4Definition,
	newObj *v1alpha1.ApiV4Definition,
) (admission.Warnings, error) {
	if newObj.IsBeingDeleted() {
		return nil, nil
	}
	// Best effort, as before the lifecycle: the reconcile never touches the ConfigMap of an API that syncs from
	// management, and some users deny ConfigMap access through RBAC.
	base.DeleteDefinitionConfigMapIfNeeded(ctx, oldObj, newObj)
	// The API type and the subscribed plans are checked first, as before the lifecycle: a forbidden change is the
	// error to report, not the reference or spec errors it leads to. The lifecycle's ImmutableFields would run after
	// them.
	errs := immutableFields(ctx, oldObj, newObj)
	if errs.IsSevere() {
		return errs.Map()
	}
	errs.MergeWith(a.Lifecycle.ValidateUpdate(ctx, oldObj, newObj))
	return errs.Map()
}

// admissionClient refuses a context reference that does not exist with the message the webhook reported before the
// lifecycle, then builds the API client.
func admissionClient(ctx context.Context, api *v1alpha1.ApiV4Definition) (*apiv4.Client, error) {
	if err := ctxref.Validate(ctx, api); err != nil {
		return nil, err
	}
	return apiv4.CreateAPIMClient(ctx, api)
}

// everyUpdate dry-runs every update, as before the lifecycle: a v4 update APIM would refuse is refused, whatever
// changed.
func everyUpdate(_, _ *v1alpha1.ApiV4Definition) bool {
	return true
}

// deleteCheck refuses to delete an API Subscriptions still reference.
func deleteCheck(ctx context.Context, obj *v1alpha1.ApiV4Definition) *errors.AdmissionErrors {
	return base.ValidateDelete(ctx, obj, search.ApiV4SubsField)
}
