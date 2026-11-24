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

package portallisting

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	internaladmission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	apimportallisting "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portallisting"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
)

var _ admission.Validator[*v1alpha1.PortalListing] = AdmissionCtrl{}
var _ admission.Defaulter[*v1alpha1.PortalListing] = AdmissionCtrl{}

type Lifecycle = lifecycle.AdmissionLifecycle[*v1alpha1.PortalListing, apimportallisting.DTO, *apim.SDKClient]

// AdmissionCtrl needs no setup: AdmissionCtrl{} is ready to use.
type AdmissionCtrl struct{}

var admissionLifecycle = lifecycle.NewAdmissionLifecycle(Lifecycle{
	ClientFactory:  admissionClient,
	PreCheck:       validateApis,
	ToDTO:          apimportallisting.ToDTO,
	DryRun:         navigation.FindingsOnError(apimportallisting.DryRun),
	PayloadChanged: everyUpdate,
	GetRemote:      apimportallisting.GetRemote,
	// The messages are the ones the webhook reported before the lifecycle.
	UnprefixedResolveErrors: true,
})

// SetupWithManager keeps the no-op defaulter: Helm registers a mutating webhook for portal listings.
func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.PortalListing{}).
		WithValidator(a).
		WithDefaulter(a).
		Complete()
}

func (a AdmissionCtrl) Default(_ context.Context, _ *v1alpha1.PortalListing) error {
	return nil
}

func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context, obj *v1alpha1.PortalListing,
) (admission.Warnings, error) {
	errs := preLifecycle(ctx, obj)
	if errs.IsSevere() {
		return errs.Map()
	}
	errs.MergeWith(admissionLifecycle.ValidateCreate(ctx, obj))
	return errs.Map()
}

func (a AdmissionCtrl) ValidateDelete(
	_ context.Context, _ *v1alpha1.PortalListing,
) (admission.Warnings, error) {
	return admission.Warnings{}, nil
}

func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context, oldObj *v1alpha1.PortalListing, newObj *v1alpha1.PortalListing,
) (admission.Warnings, error) {
	if newObj.IsBeingDeleted() {
		return admission.Warnings{}, nil
	}
	// A change of portal is refused first, as before the lifecycle: the move is the error to report, not the
	// reference errors it leads to.
	errs := validateImmutableFields(oldObj, newObj)
	if errs.IsSevere() {
		return errs.Map()
	}
	errs.MergeWith(preLifecycle(ctx, newObj))
	if errs.IsSevere() {
		return errs.Map()
	}
	errs.MergeWith(admissionLifecycle.ValidateUpdate(ctx, oldObj, newObj))
	return errs.Map()
}

// preLifecycle compiles templates and checks the API kinds before the portal, as before the lifecycle.
func preLifecycle(ctx context.Context, listing *v1alpha1.PortalListing) *gerrors.AdmissionErrors {
	errs := gerrors.NewAdmissionErrors()
	errs.Add(internaladmission.CompileAndValidateTemplate(ctx, listing))
	if errs.IsSevere() {
		return errs
	}
	errs.MergeWith(validateApiKinds(listing))
	return errs
}

// admissionClient reads the portal and checks it has a management context, with the messages the webhook
// reported before the lifecycle, then builds the client of the portal's context. Done here rather than in
// the lifecycle's ResolveRefs, which drift would also run on the previous version.
func admissionClient(ctx context.Context, obj *v1alpha1.PortalListing) (*apim.SDKClient, error) {
	if err := apimportallisting.ResolveRefsForAdmission(ctx, obj, obj.GetNamespace()); err != nil {
		return nil, err
	}
	return apimportallisting.CreateAPIMClient(ctx, obj)
}

// everyUpdate dry-runs every update, as before the lifecycle.
func everyUpdate(_, _ *v1alpha1.PortalListing) bool {
	return true
}
