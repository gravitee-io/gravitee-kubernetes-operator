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
	"context"

	commonMutate "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/common/mutate"
	wk "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/webhook"
	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ webhook.Defaulter = &ApiV4Definition{}
var _ webhook.Validator = &ApiV4Definition{}

func (api *ApiV4Definition) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(api).
		Complete()
}

func (api *ApiV4Definition) Default() {
	commonMutate.SetDefaults(api)
}

func (api *ApiV4Definition) ValidateCreate() (admission.Warnings, error) {
	return wk.ValidateApiV4(context.Background(), &api.Spec.Api, api.Name, api.Namespace, api.Spec.Context)
}

func (api *ApiV4Definition) ValidateUpdate(_ runtime.Object) (admission.Warnings, error) {
	return wk.ValidateApiV4(context.Background(), &api.Spec.Api, api.Name, api.Namespace, api.Spec.Context)
}

func (*ApiV4Definition) ValidateDelete() (admission.Warnings, error) {
	return admission.Warnings{}, nil
}
