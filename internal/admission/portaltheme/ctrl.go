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

package portaltheme

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
	apimportaltheme "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portaltheme"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
)

var _ admission.Validator[*v1alpha1.PortalTheme] = AdmissionCtrl{}
var _ admission.Defaulter[*v1alpha1.PortalTheme] = AdmissionCtrl{}

type Lifecycle = lifecycle.AdmissionLifecycle[*v1alpha1.PortalTheme, sdk.PortalTheme, *apim.SDKClient]

// AdmissionCtrl needs no setup: AdmissionCtrl{} is ready to use.
type AdmissionCtrl struct{}

var admissionLifecycle = lifecycle.NewAdmissionLifecycle(Lifecycle{
	ClientFactory:  admissionClient,
	ToDTO:          apimportaltheme.ToDTO,
	DryRun:         navigation.FindingsOnError(apimportaltheme.DryRun),
	PayloadChanged: everyUpdate,
	GetRemote:      apimportaltheme.GetRemote,
	DeleteGuard:    apimportaltheme.DeleteGuard,
	// The messages are the ones the webhook reported before the lifecycle.
	UnprefixedResolveErrors: true,
})

// SetupWithManager keeps the no-op defaulter: Helm registers a mutating webhook for themes.
func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.PortalTheme{}).
		WithValidator(a).
		WithDefaulter(a).
		Complete()
}

func (a AdmissionCtrl) Default(_ context.Context, _ *v1alpha1.PortalTheme) error {
	return nil
}

func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context, obj *v1alpha1.PortalTheme,
) (admission.Warnings, error) {
	return admissionLifecycle.ValidateCreate(ctx, obj).Map()
}

// ValidateDelete keeps a theme that portals still activate from being deleted, so that the rejection is
// reported to whoever runs the delete rather than to the reconcile loop.
func (a AdmissionCtrl) ValidateDelete(
	ctx context.Context, obj *v1alpha1.PortalTheme,
) (admission.Warnings, error) {
	return admissionLifecycle.ValidateDelete(ctx, obj).Map()
}

func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context, oldObj *v1alpha1.PortalTheme, newObj *v1alpha1.PortalTheme,
) (admission.Warnings, error) {
	return admissionLifecycle.ValidateUpdate(ctx, oldObj, newObj).Map()
}

var errNoContext = errors.New("a management context reference (spec.contextRef) is required")

// admissionClient refuses a theme without context, or with a context the cluster does not hold, with the
// messages the webhook reported before the lifecycle, then builds the client.
func admissionClient(ctx context.Context, thm *v1alpha1.PortalTheme) (*apim.SDKClient, error) {
	if !thm.HasContext() {
		return nil, errNoContext
	}
	if err := ctxref.Validate(ctx, thm); err != nil {
		return nil, err
	}
	return apimportaltheme.CreateAPIMClient(ctx, thm)
}

// everyUpdate dry-runs every update, as before the lifecycle.
func everyUpdate(_, _ *v1alpha1.PortalTheme) bool {
	return true
}
