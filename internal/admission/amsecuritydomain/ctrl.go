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

	amsdk "github.com/gravitee-io-labs/gravitee-automation-tools/am-sdk/pkg/sdk/domain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/securitydomain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/lifecycle"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ admission.Validator[*v1alpha1.AMSecurityDomain] = AdmissionCtrl{}

type AdmissionCtrl struct {
	Lifecycle lifecycle.AdmissionLifecycle[*v1alpha1.AMSecurityDomain, amsdk.Domain, *am.Client]
}

func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	a.Lifecycle.ClientFactory = internal.CreateAMClient
	a.Lifecycle.DryRun = internal.DryRun
	a.Lifecycle.GetRemote = internal.GetRemote
	a.Lifecycle.ToDTO = internal.ToDomainDTO
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
	if newObj.IsBeingDeleted() {
		return admission.Warnings{}, nil
	}
	return a.Lifecycle.ValidateUpdate(ctx, oldObj, newObj).Map()
}

func (a AdmissionCtrl) ValidateDelete(
	ctx context.Context, obj *v1alpha1.AMSecurityDomain,
) (admission.Warnings, error) {
	return a.Lifecycle.ValidateDelete(ctx, obj).Map()
}
