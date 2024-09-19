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

package mctx

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	ctrl "sigs.k8s.io/controller-runtime"
)

var _ admission.CustomValidator = AdmissionCtrl{}
var _ admission.CustomDefaulter = AdmissionCtrl{}

func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&v1alpha1.ManagementContext{}).
		WithValidator(a).
		WithDefaulter(a).
		Complete()
}

type AdmissionCtrl struct{}

// Default implements admission.CustomDefaulter.
func (a AdmissionCtrl) Default(ctx context.Context, obj runtime.Object) error {
	return SetDefaults(ctx, obj)
}

// ValidateCreate implements admission.CustomValidator.
func (a AdmissionCtrl) ValidateCreate(
	ctx context.Context,
	obj runtime.Object,
) (admission.Warnings, error) {
	return validateCreate(ctx, obj).Map()
}

// ValidateDelete implements admission.CustomValidator.
func (a AdmissionCtrl) ValidateDelete(
	ctx context.Context, obj runtime.Object,
) (admission.Warnings, error) {
	return admission.Warnings{}, nil
}

// ValidateUpdate implements admission.CustomValidator.
func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context,
	oldObj runtime.Object,
	newObj runtime.Object,
) (admission.Warnings, error) {
	return validateCreate(ctx, newObj).Map()
}
