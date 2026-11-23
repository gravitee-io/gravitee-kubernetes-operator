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

package authzpdp

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimauthzpdp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hrid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ admission.Validator[*v1alpha1.AuthzPdp] = AdmissionCtrl{}
var _ admission.Defaulter[*v1alpha1.AuthzPdp] = AdmissionCtrl{}

type Lifecycle = lifecycle.AdmissionLifecycle[*v1alpha1.AuthzPdp, sdk.AuthzPdp, *apim.SDKClient]

type AdmissionCtrl struct {
	Lifecycle Lifecycle
}

// NewAdmissionCtrl wires the AuthzPdp admission steps: the context, the derived HRID, then the platform dry
// run, which refuses a routing scope another PDP holds, a strict value that disagrees with the engine and an
// endpoint path another API uses. Schema and CEL own the spec-only rules (the targetPdpId and tag grammar,
// their immutability).
func NewAdmissionCtrl() AdmissionCtrl {
	return AdmissionCtrl{
		Lifecycle: lifecycle.NewAdmissionLifecycle(Lifecycle{
			ResolveRefs:   apim.ResolveContextRef[*v1alpha1.AuthzPdp],
			ClientFactory: apimauthzpdp.CreateAPIMClient,
			PreCheck:      preCheck,
			ToDTO:         apimauthzpdp.ToDTO,
			DryRun:        apimauthzpdp.DryRun,
			GetRemote:     apimauthzpdp.GetRemote,
			DeleteCheck:   validateDelete,
		}),
	}
}

func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	a = NewAdmissionCtrl()
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.AuthzPdp{}).
		WithValidator(a).
		WithDefaulter(a).
		Complete()
}

func (a AdmissionCtrl) Default(_ context.Context, _ *v1alpha1.AuthzPdp) error {
	return nil
}

func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context, obj *v1alpha1.AuthzPdp,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateCreate(ctx, obj).Map()
}

func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context, oldObj *v1alpha1.AuthzPdp, newObj *v1alpha1.AuthzPdp,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateUpdate(ctx, oldObj, newObj).Map()
}

// ValidateDelete admits every deletion, with a warning naming the resources that still target the PDP:
// refusing would make pruning a directory depend on the order resources are deleted in.
func (a AdmissionCtrl) ValidateDelete(
	ctx context.Context, obj *v1alpha1.AuthzPdp,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateDelete(ctx, obj).Map()
}

// preCheck validates the HRID derived from the namespace and name.
func preCheck(_ context.Context, pdp *v1alpha1.AuthzPdp) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if err := hrid.Validate(refs.NewNamespacedNameFromObject(pdp).HRID()); err != nil {
		errs.AddSevere(err.Error())
	}
	return errs
}
