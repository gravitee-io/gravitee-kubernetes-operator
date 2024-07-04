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

var _ webhook.Defaulter = &ApiResource{}
var _ webhook.Validator = &ApiResource{}

func (res *ApiResource) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(res).
		Complete()
}

func (res *ApiResource) Default() {}

func (res *ApiResource) ValidateCreate() (admission.Warnings, error) {
	return admission.Warnings{}, nil
}

func (res *ApiResource) ValidateUpdate(_ runtime.Object) (admission.Warnings, error) {
	return admission.Warnings{}, nil
}

func (*ApiResource) ValidateDelete() (admission.Warnings, error) {
	return admission.Warnings{}, nil
}
