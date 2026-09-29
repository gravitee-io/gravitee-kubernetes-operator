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

package amsecuritydomain

import (
	"context"

	amsdk "github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/securitydomain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ admission.Validator[*v1alpha1.AMSecurityDomain] = AdmissionCtrl{}

type Lifecycle = lifecycle.AdmissionLifecycle[*v1alpha1.AMSecurityDomain, amsdk.Domain, *am.Client]

type AdmissionCtrl struct {
	Lifecycle Lifecycle
}

func NewAdmissionCtrl() AdmissionCtrl {
	return AdmissionCtrl{
		Lifecycle: lifecycle.NewAdmissionLifecycle(Lifecycle{
			ClientFactory: internal.CreateAMClient,
			PreCheck:      internal.ValidateKey,
			DryRun:        internal.DryRun,
			GetRemote:     internal.GetRemote,
			ToDTO:         internal.ToDomainDTO,
		}),
	}
}

func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	a = NewAdmissionCtrl()
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.AMSecurityDomain{}).
		WithValidator(a).
		Complete()
}

func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context,
	newObj *v1alpha1.AMSecurityDomain,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateCreate(ctx, newObj).Map()
}

func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context,
	oldObj *v1alpha1.AMSecurityDomain,
	newObj *v1alpha1.AMSecurityDomain,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateUpdate(ctx, oldObj, newObj).Map()
}

func (a AdmissionCtrl) ValidateDelete(
	ctx context.Context, obj *v1alpha1.AMSecurityDomain,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateDelete(ctx, obj).Map()
}
