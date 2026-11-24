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

package portallink

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	apimportallink "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portallink"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
)

var _ admission.Validator[*v1alpha1.PortalLink] = AdmissionCtrl{}
var _ admission.Defaulter[*v1alpha1.PortalLink] = AdmissionCtrl{}

type Lifecycle = lifecycle.AdmissionLifecycle[*v1alpha1.PortalLink, apimportallink.DTO, *apim.SDKClient]

// AdmissionCtrl needs no setup: AdmissionCtrl{} is ready to use.
type AdmissionCtrl struct{}

var admissionLifecycle = lifecycle.NewAdmissionLifecycle(Lifecycle{
	ClientFactory:  admissionClient,
	ToDTO:          apimportallink.ToDTO,
	DryRun:         navigation.FindingsOnError(apimportallink.DryRun),
	PayloadChanged: everyUpdate,
	GetRemote:      apimportallink.GetRemote,
	// The messages are the ones the webhook reported before the lifecycle.
	UnprefixedResolveErrors: true,
})

// SetupWithManager keeps the no-op defaulter: Helm registers a mutating webhook for portal links.
func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.PortalLink{}).
		WithValidator(a).
		WithDefaulter(a).
		Complete()
}

func (a AdmissionCtrl) Default(_ context.Context, _ *v1alpha1.PortalLink) error {
	return nil
}

func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context, obj *v1alpha1.PortalLink,
) (admission.Warnings, error) {
	return admissionLifecycle.ValidateCreate(ctx, obj).Map()
}

func (a AdmissionCtrl) ValidateDelete(
	_ context.Context, _ *v1alpha1.PortalLink,
) (admission.Warnings, error) {
	return admission.Warnings{}, nil
}

func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context, oldObj *v1alpha1.PortalLink, newObj *v1alpha1.PortalLink,
) (admission.Warnings, error) {
	if newObj.IsBeingDeleted() {
		return admission.Warnings{}, nil
	}
	errs := preLifecycleUpdate(ctx, oldObj, newObj)
	if errs.IsSevere() {
		return errs.Map()
	}
	errs.MergeWith(admissionLifecycle.ValidateUpdate(ctx, oldObj, newObj))
	return errs.Map()
}

// admissionClient checks the parent references and reads the parent, with the messages the webhook
// reported before the lifecycle, then builds the client of the parent's context. Done here rather than in
// the lifecycle's ResolveRefs, which drift would also run on the previous version.
func admissionClient(ctx context.Context, obj *v1alpha1.PortalLink) (*apim.SDKClient, error) {
	if err := apimportallink.ResolveRefsForAdmission(ctx, obj, obj.GetNamespace()); err != nil {
		return nil, err
	}
	return apimportallink.CreateAPIMClient(ctx, obj)
}

// everyUpdate dry-runs every update, as before the lifecycle.
func everyUpdate(_, _ *v1alpha1.PortalLink) bool {
	return true
}

// preLifecycleUpdate refuses a change of parent first, as before the lifecycle: the move is the error to
// report, not the reference errors it leads to.
func preLifecycleUpdate(_ context.Context, oldObj, newObj *v1alpha1.PortalLink) *gerrors.AdmissionErrors {
	return verifyImmutableFields(oldObj, newObj)
}
