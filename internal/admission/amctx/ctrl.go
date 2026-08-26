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

package amctx

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ admission.Validator[*v1alpha1.AMContext] = AdmissionCtrl{}

type AdmissionCtrl struct{}

func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.AMContext{}).
		WithValidator(a).
		Complete()
}

func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context,
	obj *v1alpha1.AMContext,
) (admission.Warnings, error) {
	return validateCreate(ctx, obj).Map()
}

func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context,
	_ *v1alpha1.AMContext,
	newObj *v1alpha1.AMContext,
) (admission.Warnings, error) {
	if newObj.IsBeingDeleted() {
		return admission.Warnings{}, nil
	}
	return validateCreate(ctx, newObj).Map()
}

func (a AdmissionCtrl) ValidateDelete(
	_ context.Context, _ *v1alpha1.AMContext,
) (admission.Warnings, error) {
	// Nothing references an AMContext until AM-7506.
	return admission.Warnings{}, nil
}
