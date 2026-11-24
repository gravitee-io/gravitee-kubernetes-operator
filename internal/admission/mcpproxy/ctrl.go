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

package mcpproxy

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimmcpproxy "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ admission.Validator[*v1alpha1.McpProxy] = AdmissionCtrl{}
var _ admission.Defaulter[*v1alpha1.McpProxy] = AdmissionCtrl{}

type Lifecycle = lifecycle.AdmissionLifecycle[*v1alpha1.McpProxy, sdk.AimMcpProxy, *apim.SDKClient]

type AdmissionCtrl struct {
	Lifecycle Lifecycle
}

// NewAdmissionCtrl wires the McpProxy admission steps. Schema and CEL own the spec-only rules (entityId
// grammar, union pairing, mode and the immutable fields). Deletion is always admitted: nothing references
// an McpProxy.
func NewAdmissionCtrl() AdmissionCtrl {
	return AdmissionCtrl{
		Lifecycle: lifecycle.NewAdmissionLifecycle(Lifecycle{
			ResolveRefs:   apim.ResolveContextRef[*v1alpha1.McpProxy],
			ClientFactory: apimmcpproxy.CreateAPIMClient,
			PreCheck:      preCheck,
			ToDTO:         apimmcpproxy.ToDTO,
			DryRun:        apimmcpproxy.DryRun,
			DeferDryRun:   deferDryRun,
			GetRemote:     apimmcpproxy.GetRemote,
		}),
	}
}

func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	a = NewAdmissionCtrl()
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.McpProxy{}).
		WithValidator(a).
		WithDefaulter(a).
		Complete()
}

func (a AdmissionCtrl) Default(_ context.Context, _ *v1alpha1.McpProxy) error {
	return nil
}

func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context, obj *v1alpha1.McpProxy,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateCreate(ctx, obj).Map()
}

func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context, oldObj *v1alpha1.McpProxy, newObj *v1alpha1.McpProxy,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateUpdate(ctx, oldObj, newObj).Map()
}

func (a AdmissionCtrl) ValidateDelete(
	ctx context.Context, obj *v1alpha1.McpProxy,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateDelete(ctx, obj).Map()
}
