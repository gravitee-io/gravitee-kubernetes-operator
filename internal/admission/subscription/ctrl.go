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

package subscription

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimsubscription "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ admission.Validator[*v1alpha1.Subscription] = AdmissionCtrl{}
var _ admission.Defaulter[*v1alpha1.Subscription] = AdmissionCtrl{}

type Lifecycle = lifecycle.AdmissionLifecycle[*v1alpha1.Subscription, apimsubscription.DTO, *apim.SDKClient]

// AdmissionCtrl validates Subscriptions. Its lifecycle runs the update checks; a create runs the same lifecycle
// with the create checks, which also cover the state of the API and the application.
type AdmissionCtrl struct {
	Lifecycle Lifecycle
}

func NewAdmissionCtrl() AdmissionCtrl {
	return AdmissionCtrl{
		Lifecycle: lifecycle.NewAdmissionLifecycle(Lifecycle{
			ResolveRefs:     apimsubscription.ResolveRefs,
			ClientFactory:   apimsubscription.CreateAPIMClient,
			PreCheck:        validateUpdate,
			ImmutableFields: validateImmutableProperties,
			ToDTO:           apimsubscription.ToResolvedDTO,
			DryRun:          apimsubscription.DryRun,
			GetRemote:       apimsubscription.GetRemote,
		}),
	}
}

func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	a = NewAdmissionCtrl()
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.Subscription{}).
		WithValidator(a).
		WithDefaulter(a).
		Complete()
}

func (a AdmissionCtrl) Default(_ context.Context, _ *v1alpha1.Subscription) error {
	return nil
}

func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context, obj *v1alpha1.Subscription,
) (admission.Warnings, error) {
	create := a.Lifecycle
	create.PreCheck = validateCreate
	return create.ValidateCreate(ctx, obj).Map()
}

func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context, oldObj *v1alpha1.Subscription, newObj *v1alpha1.Subscription,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateUpdate(ctx, oldObj, newObj).Map()
}

func (a AdmissionCtrl) ValidateDelete(
	ctx context.Context, obj *v1alpha1.Subscription,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateDelete(ctx, obj).Map()
}
