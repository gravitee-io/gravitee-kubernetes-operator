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

package policygroups

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	spg "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/sharedpolicygroup"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
)

var _ admission.Validator[*v1alpha1.SharedPolicyGroup] = AdmissionCtrl{}
var _ admission.Defaulter[*v1alpha1.SharedPolicyGroup] = AdmissionCtrl{}

type Lifecycle = lifecycle.AdmissionLifecycle[*v1alpha1.SharedPolicyGroup, spg.DTO, *apim.SDKClient]

type AdmissionCtrl struct {
	Lifecycle Lifecycle
}

func NewAdmissionCtrl() AdmissionCtrl {
	return AdmissionCtrl{
		Lifecycle: lifecycle.NewAdmissionLifecycle(Lifecycle{
			ClientFactory:  admissionClient,
			ToDTO:          spg.ToDTO,
			DryRun:         spg.DryRun,
			PayloadChanged: spg.PayloadChanged,
			GetRemote:      spg.GetRemote,
			DeleteGuard:    spg.DeleteGuard,
			// The shared policy group admission messages are the ones the webhook reported before the lifecycle.
			UnprefixedResolveErrors: true,
		}),
	}
}

// SetupWithManager keeps the no-op defaulter: Helm registers a mutating webhook for shared policy groups.
func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	a = NewAdmissionCtrl()
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.SharedPolicyGroup{}).
		WithValidator(a).
		WithDefaulter(a).
		Complete()
}

func (a AdmissionCtrl) Default(_ context.Context, _ *v1alpha1.SharedPolicyGroup) error {
	return nil
}

func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context,
	obj *v1alpha1.SharedPolicyGroup,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateCreate(ctx, obj).Map()
}

// ValidateUpdate checks the fields APIM cannot change first: a forbidden change is the error to report, not the
// dry run errors it leads to.
func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context,
	oldObj *v1alpha1.SharedPolicyGroup,
	newObj *v1alpha1.SharedPolicyGroup,
) (admission.Warnings, error) {
	if newObj.IsBeingDeleted() {
		return nil, nil
	}
	errs := errors.NewAdmissionErrors()
	errs.Add(validateImmutableFields(oldObj, newObj))
	if errs.IsSevere() {
		return errs.Map()
	}
	errs.MergeWith(a.Lifecycle.ValidateUpdate(ctx, oldObj, newObj))
	return errs.Map()
}

func (a AdmissionCtrl) ValidateDelete(
	ctx context.Context, obj *v1alpha1.SharedPolicyGroup,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateDelete(ctx, obj).Map()
}

// admissionClient refuses a context reference that does not exist with the message the webhook reported before the
// lifecycle, then builds the SDK client.
func admissionClient(ctx context.Context, obj *v1alpha1.SharedPolicyGroup) (*apim.SDKClient, error) {
	if err := ctxref.Validate(ctx, obj); err != nil {
		return nil, err
	}
	return spg.CreateAPIMClient(ctx, obj)
}
