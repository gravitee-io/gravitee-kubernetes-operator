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

package authzpolicy

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimauthzpolicy "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/authzpolicy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ admission.Validator[*v1alpha1.AuthzPolicy] = AdmissionCtrl{}
var _ admission.Defaulter[*v1alpha1.AuthzPolicy] = AdmissionCtrl{}

type Lifecycle = lifecycle.AdmissionLifecycle[*v1alpha1.AuthzPolicy, sdk.AuthzPolicy, *apim.SDKClient]

type AdmissionCtrl struct {
	Lifecycle Lifecycle
}

// NewAdmissionCtrl wires the AuthzPolicy admission steps: the pdpRef kinds and the context, the derived HRID, then
// the platform dry run with the resolved targets, which refuses an entityId another policy holds, a text that
// does not parse, a template and a violation of a strict gateway's schema, and warns about an MCP policy the
// default engine never sees and a tool named by its catalog id. Schema and CEL own the spec-only rules (the
// pdpRef/scope union, the entityId and scope grammars, immutability). Deletion is always admitted with no warning
// and no platform call: a policy delete takes nothing else with it, and APIM 4.13 ignores dryRun on a DELETE.
func NewAdmissionCtrl() AdmissionCtrl {
	return AdmissionCtrl{
		Lifecycle: lifecycle.NewAdmissionLifecycle(Lifecycle{
			ResolveRefs:   resolveRefs,
			ClientFactory: apimauthzpolicy.CreateAPIMClient,
			PreCheck:      preCheck,
			ToDTO:         apimauthzpolicy.ToDTO,
			DryRun:        apimauthzpolicy.DryRun,
			DeferDryRun:   deferDryRun,
			GetRemote:     apimauthzpolicy.GetRemote,
			DeferDrift:    deferDrift,
		}),
	}
}

func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	a = NewAdmissionCtrl()
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.AuthzPolicy{}).
		WithValidator(a).
		WithDefaulter(a).
		Complete()
}

func (a AdmissionCtrl) Default(_ context.Context, _ *v1alpha1.AuthzPolicy) error {
	return nil
}

func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context, obj *v1alpha1.AuthzPolicy,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateCreate(ctx, obj).Map()
}

func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context, oldObj *v1alpha1.AuthzPolicy, newObj *v1alpha1.AuthzPolicy,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateUpdate(ctx, oldObj, newObj).Map()
}

func (a AdmissionCtrl) ValidateDelete(
	ctx context.Context, obj *v1alpha1.AuthzPolicy,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateDelete(ctx, obj).Map()
}
