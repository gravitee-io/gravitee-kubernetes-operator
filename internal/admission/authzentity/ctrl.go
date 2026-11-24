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

package authzentity

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimauthzentity "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ admission.Validator[*v1alpha1.AuthzEntity] = AdmissionCtrl{}
var _ admission.Defaulter[*v1alpha1.AuthzEntity] = AdmissionCtrl{}

type Lifecycle = lifecycle.AdmissionLifecycle[*v1alpha1.AuthzEntity, sdk.AuthzEntity, *apim.SDKClient]

type AdmissionCtrl struct {
	Lifecycle Lifecycle
}

// NewAdmissionCtrl wires the AuthzEntity admission steps: the pdpRef kinds and the context, the derived HRID,
// then the platform dry run with the resolved targets, which refuses an entityId another source holds (AUTHZ-5),
// reserved attribute keys, duplicate or cyclic parents, an Action with targets and an entity off a strict schema.
// Schema and CEL own the spec-only rules (the pdpRef/scope union, the scope and type grammars, immutability).
// Deletion is always admitted, with a warning naming the entities the platform deletes with it: refusing would
// make pruning a directory depend on the order resources are deleted in.
func NewAdmissionCtrl() AdmissionCtrl {
	return AdmissionCtrl{
		Lifecycle: lifecycle.NewAdmissionLifecycle(Lifecycle{
			ResolveRefs:   resolveRefs,
			ClientFactory: apimauthzentity.CreateAPIMClient,
			PreCheck:      preCheck,
			ToDTO:         apimauthzentity.ToDTO,
			DryRun:        apimauthzentity.DryRun,
			DeferDryRun:   deferDryRun,
			GetRemote:     apimauthzentity.GetRemote,
			DeferDrift:    deferDrift,
			DeleteCheck:   validateDelete,
		}),
	}
}

func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	a = NewAdmissionCtrl()
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.AuthzEntity{}).
		WithValidator(a).
		WithDefaulter(a).
		Complete()
}

func (a AdmissionCtrl) Default(_ context.Context, _ *v1alpha1.AuthzEntity) error {
	return nil
}

func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context, obj *v1alpha1.AuthzEntity,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateCreate(ctx, obj).Map()
}

func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context, oldObj *v1alpha1.AuthzEntity, newObj *v1alpha1.AuthzEntity,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateUpdate(ctx, oldObj, newObj).Map()
}

func (a AdmissionCtrl) ValidateDelete(
	ctx context.Context, obj *v1alpha1.AuthzEntity,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateDelete(ctx, obj).Map()
}
