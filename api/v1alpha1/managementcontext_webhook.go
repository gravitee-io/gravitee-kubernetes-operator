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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	wk "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/webhook"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ webhook.Defaulter = &ManagementContext{}
var _ webhook.Validator = &ManagementContext{}

func (ctx *ManagementContext) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(ctx).
		Complete()
}

func (ctx *ManagementContext) Default() {}

func (ctx *ManagementContext) ValidateCreate() (admission.Warnings, error) {
	return validateManagementContext(ctx)
}

func (ctx *ManagementContext) ValidateUpdate(_ runtime.Object) (admission.Warnings, error) {
	return validateManagementContext(ctx)
}

func (*ManagementContext) ValidateDelete() (admission.Warnings, error) {
	return admission.Warnings{}, nil
}

func validateManagementContext(ctx *ManagementContext) (admission.Warnings, error) {
	// Make sure the secret exist
	if ctx.Spec.HasSecretRef() {
		secret := new(corev1.Secret)
		err := k8s.GetClient().Get(context.Background(), ctx.Spec.SecretRef().NamespacedName(), secret)
		if err != nil {
			return admission.Warnings{}, fmt.Errorf("can't create management context [%s] because it is using "+
				"sercret [%v] that doesn't exist in the cluster", ctx.Name, ctx.Spec.SecretRef())
		}
	}

	ctxRef := refs.NewNamespacedName(ctx.Namespace, ctx.Name)
	if err := wk.CheckAPIMAvailability(&ctxRef); err != nil {
		return admission.Warnings{err.Error()}, nil //nolint:nilerr // changed to warning
	}

	return admission.Warnings{}, nil
}
