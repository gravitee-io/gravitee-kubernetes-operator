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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logr "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var log = logr.Log.WithName("managementcontext-resource")

type validator struct{}

func (v *validator) ValidateCreate(ctx context.Context, obj runtime.Object) error {
	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		log.Error(err, "failed to get admission request")
	}
	log.Info("custom validate create", "request", req)
	return nil
}

func (v *validator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) error {
	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		log.Error(err, "failed to get admission request")
	}

	log.Info("custom validate update", "request", req)
	return nil
}

func (v *validator) ValidateDelete(ctx context.Context, obj runtime.Object) error {
	return nil
}

func (r *ManagementContext) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(&validator{}).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-gravitee-io-v1alpha1-managementcontext,mutating=true,failurePolicy=fail,sideEffects=None,groups=gravitee.io,resources=managementcontexts,verbs=create;update,versions=v1alpha1,name=mmanagementcontext.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &ManagementContext{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *ManagementContext) Default() {
	log.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-gravitee-io-v1alpha1-managementcontext,mutating=false,failurePolicy=fail,sideEffects=None,groups=gravitee.io,resources=managementcontexts,verbs=create;update,versions=v1alpha1,name=vmanagementcontext.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &ManagementContext{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *ManagementContext) ValidateCreate() error {
	log.Info("validate create", "name", r.Name)

	err := errors.NewInvalid(GroupVersion.WithKind("ManagementContext").GroupKind(), r.Name, field.ErrorList{
		field.Invalid(field.NewPath("spec").Child("field"), "r.Spec.Field", "invalid value"),
	})

	err.ErrStatus.Code = 299

	return err
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *ManagementContext) ValidateUpdate(old runtime.Object) error {
	log.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *ManagementContext) ValidateDelete() error {
	log.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
