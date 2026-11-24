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

package portal

import (
	"context"
	"errors"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/ctxref"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	apimportal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portal"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
)

var _ admission.Validator[*v1alpha1.Portal] = AdmissionCtrl{}
var _ admission.Defaulter[*v1alpha1.Portal] = AdmissionCtrl{}

type Lifecycle = lifecycle.AdmissionLifecycle[*v1alpha1.Portal, sdk.Portal, *apim.SDKClient]

// AdmissionCtrl needs no setup: AdmissionCtrl{} is ready to use.
type AdmissionCtrl struct{}

var admissionLifecycle = lifecycle.NewAdmissionLifecycle(Lifecycle{
	ClientFactory:  admissionClient,
	PreCheck:       resolveThemeRef,
	ToDTO:          apimportal.ToDTO,
	DryRun:         navigation.FindingsOnError(apimportal.DryRun),
	PayloadChanged: everyUpdate,
	GetRemote:      apimportal.GetRemote,
	DriftDTO:       apimportal.DriftDTO,
	// The messages are the ones the webhook reported before the lifecycle.
	UnprefixedResolveErrors: true,
})

// SetupWithManager keeps the no-op defaulter: Helm registers a mutating webhook for portals.
func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.Portal{}).
		WithValidator(a).
		WithDefaulter(a).
		Complete()
}

func (a AdmissionCtrl) Default(_ context.Context, _ *v1alpha1.Portal) error {
	return nil
}

func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context, obj *v1alpha1.Portal,
) (admission.Warnings, error) {
	errs := preLifecycle(obj)
	if errs.IsSevere() {
		return errs.Map()
	}
	errs.MergeWith(admissionLifecycle.ValidateCreate(ctx, obj))
	return errs.Map()
}

func (a AdmissionCtrl) ValidateDelete(
	_ context.Context, _ *v1alpha1.Portal,
) (admission.Warnings, error) {
	return admission.Warnings{}, nil
}

func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context, oldObj *v1alpha1.Portal, newObj *v1alpha1.Portal,
) (admission.Warnings, error) {
	if newObj.IsBeingDeleted() {
		return admission.Warnings{}, nil
	}
	errs := preLifecycle(newObj)
	if errs.IsSevere() {
		return errs.Map()
	}
	errs.MergeWith(admissionLifecycle.ValidateUpdate(ctx, oldObj, newObj))
	return errs.Map()
}

// preLifecycle runs the checks the webhook ran before the context and the theme, as before the lifecycle:
// the deprecated navigation, then the theme kind.
func preLifecycle(prtl *v1alpha1.Portal) *gerrors.AdmissionErrors {
	errs := validateNavigation(prtl)
	if errs.IsSevere() {
		return errs
	}
	errs.MergeWith(validateThemeKind(prtl))
	return errs
}

var errNoContext = errors.New("a management context reference (spec.contextRef) is required")

// admissionClient refuses a portal without context, or with a context the cluster does not hold, with the
// messages the webhook reported before the lifecycle, then builds the client.
func admissionClient(ctx context.Context, prtl *v1alpha1.Portal) (*apim.SDKClient, error) {
	if !prtl.HasContext() {
		return nil, errNoContext
	}
	if err := ctxref.Validate(ctx, prtl); err != nil {
		return nil, err
	}
	return apimportal.CreateAPIMClient(ctx, prtl)
}

// everyUpdate dry-runs every update, as before the lifecycle.
func everyUpdate(_, _ *v1alpha1.Portal) bool {
	return true
}
