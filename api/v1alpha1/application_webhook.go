/*
 * Copyright (C) 2015 The Gravitee team (http://gravitee.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ webhook.Defaulter = &Application{}
var _ webhook.Validator = &Application{}

func (app *Application) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(app).
		Complete()
}

func (app *Application) Default() {}

func (app *Application) ValidateCreate() (admission.Warnings, error) {
	return admission.Warnings{}, nil
}

func (app *Application) ValidateUpdate(_ runtime.Object) (admission.Warnings, error) {
	return admission.Warnings{}, nil
}

func (*Application) ValidateDelete() (admission.Warnings, error) {
	return admission.Warnings{}, nil
}
