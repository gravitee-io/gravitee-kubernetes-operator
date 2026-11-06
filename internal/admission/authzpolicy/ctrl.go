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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	ctrl "sigs.k8s.io/controller-runtime"
)

var _ admission.Validator[*v1alpha1.AuthzPolicy] = AdmissionCtrl{}
var _ admission.Defaulter[*v1alpha1.AuthzPolicy] = AdmissionCtrl{}

func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.AuthzPolicy{}).
		WithValidator(a).
		WithDefaulter(a).
		Complete()
}

type AdmissionCtrl struct{}

func (a AdmissionCtrl) Default(_ context.Context, _ *v1alpha1.AuthzPolicy) error {
	return nil
}

func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context, obj *v1alpha1.AuthzPolicy,
) (admission.Warnings, error) {
	return validateCreate(ctx, obj).Map()
}

// ValidateDelete admits every deletion with no warning and no platform call: a policy delete
// takes nothing else with it, and APIM 4.13 ignores dryRun on a DELETE and deletes.
func (a AdmissionCtrl) ValidateDelete(
	_ context.Context, _ *v1alpha1.AuthzPolicy,
) (admission.Warnings, error) {
	return admission.Warnings{}, nil
}

func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context, oldObj *v1alpha1.AuthzPolicy, newObj *v1alpha1.AuthzPolicy,
) (admission.Warnings, error) {
	if newObj.IsBeingDeleted() {
		return admission.Warnings{}, nil
	}

	return validateUpdate(ctx, oldObj, newObj).Map()
}
