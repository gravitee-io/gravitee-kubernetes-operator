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

package catalogmcpserver

import (
	"context"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimcatalog "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/catalogmcpserver"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ admission.Validator[*v1alpha1.CatalogMcpServer] = AdmissionCtrl{}
var _ admission.Defaulter[*v1alpha1.CatalogMcpServer] = AdmissionCtrl{}

type Lifecycle = lifecycle.AdmissionLifecycle[*v1alpha1.CatalogMcpServer, sdk.AimCatalogMcpServer, *apim.SDKClient]

type AdmissionCtrl struct {
	Lifecycle Lifecycle
}

// NewAdmissionCtrl wires the CatalogMcpServer admission steps. The dry run discovers the upstream, so an
// unreachable server or a taken entityId is refused. Schema and CEL own the spec-only rules (entityId grammar
// and immutability, auth variant pairing).
func NewAdmissionCtrl() AdmissionCtrl {
	return AdmissionCtrl{
		Lifecycle: lifecycle.NewAdmissionLifecycle(Lifecycle{
			ResolveRefs:   apim.ResolveContextRef[*v1alpha1.CatalogMcpServer],
			ClientFactory: apimcatalog.CreateAPIMClient,
			ToDTO:         apimcatalog.ToDTO,
			DryRun:        apimcatalog.DryRun,
			GetRemote:     apimcatalog.GetRemote,
			DeleteCheck:   validateDelete,
		}),
	}
}

func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	a = NewAdmissionCtrl()
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.CatalogMcpServer{}).
		WithValidator(a).
		WithDefaulter(a).
		Complete()
}

func (a AdmissionCtrl) Default(_ context.Context, _ *v1alpha1.CatalogMcpServer) error {
	return nil
}

func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context, obj *v1alpha1.CatalogMcpServer,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateCreate(ctx, obj).Map()
}

func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context, oldObj *v1alpha1.CatalogMcpServer, newObj *v1alpha1.CatalogMcpServer,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateUpdate(ctx, oldObj, newObj).Map()
}

func (a AdmissionCtrl) ValidateDelete(
	ctx context.Context, obj *v1alpha1.CatalogMcpServer,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateDelete(ctx, obj).Map()
}
