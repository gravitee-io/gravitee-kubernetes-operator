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

package application

import (
	"context"
	"strconv"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	ctrl "sigs.k8s.io/controller-runtime"
)

var _ webhook.CustomValidator = AdmissionCtrl{}
var _ webhook.CustomDefaulter = AdmissionCtrl{}

func (a AdmissionCtrl) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&v1alpha1.Application{}).
		WithValidator(a).
		WithDefaulter(a).
		Complete()
}

type AdmissionCtrl struct{}

// Default implements admission.CustomDefaulter.
func (a AdmissionCtrl) Default(_ context.Context, obj runtime.Object) error {
	app, ok := obj.(*v1alpha1.Application)
	if !ok {
		return nil
	}
	defaultClientCertificates(app)
	return nil
}

func defaultClientCertificates(app *v1alpha1.Application) {
	if app.Spec.Settings == nil || app.Spec.Settings.TLS == nil {
		return
	}

	for i := range app.Spec.Settings.TLS.ClientCertificates {
		cert := &app.Spec.Settings.TLS.ClientCertificates[i]

		if cert.Name == "" {
			cert.Name = app.Spec.Name + "-" + strconv.Itoa(i)
		}

		if cert.Ref != nil {
			if cert.Ref.Kind == "" {
				cert.Ref.Kind = "secrets"
			}
			if cert.Ref.Key == "" {
				cert.Ref.Key = "tls.crt"
			}
		}
	}
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
	return validateDelete(ctx, obj).Map()
}

// ValidateUpdate implements admission.CustomValidator.
func (a AdmissionCtrl) ValidateUpdate(
	ctx context.Context,
	oldObj runtime.Object,
	newObj runtime.Object,
) (admission.Warnings, error) {
	return validateUpdate(ctx, oldObj, newObj).Map()
}
