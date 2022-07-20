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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var managementcontextlog = logf.Log.WithName("managementcontext-resource")

func (r *ManagementContext) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-gravitee-io-v1alpha1-managementcontext,mutating=true,failurePolicy=fail,sideEffects=None,groups=gravitee.io,resources=managementcontexts,verbs=create;update,versions=v1alpha1,name=mmanagementcontext.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &ManagementContext{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *ManagementContext) Default() {
	managementcontextlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-gravitee-io-v1alpha1-managementcontext,mutating=false,failurePolicy=fail,sideEffects=None,groups=gravitee.io,resources=managementcontexts,verbs=create;update,versions=v1alpha1,name=vmanagementcontext.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &ManagementContext{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ManagementContext) ValidateCreate() error {
	managementcontextlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ManagementContext) ValidateUpdate(old runtime.Object) error {
	managementcontextlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ManagementContext) ValidateDelete() error {
	managementcontextlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
