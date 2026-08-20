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

package v2

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func validateCreate(ctx context.Context, api *v1alpha1.ApiDefinition) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if api.GetAnnotations()[core.IngressTemplateAnnotation] == env.TrueString {
		log.Global.Debugf("skipping validation for ingress template %s", api.GetName())
		return errs
	}
	errs.MergeWith(base.ValidateCreate(ctx, api))
	if errs.IsSevere() {
		return errs
	}
	if api.HasContext() {
		errs.MergeWith(validateDryRun(ctx, api))
	}
	return errs
}

func validateDryRun(ctx context.Context, api *v1alpha1.ApiDefinition) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp := api.DeepCopy()

	apimClient, err := apim.FromContextRef(ctx, cp.ContextRef(), cp.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	cp.PopulateIDs(apimClient.Context, k8s.IsAutomationAPIManaged(api))

	status, err := apimClient.APIs.DryRunImportV2(&cp.Spec.Api)
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}
	for _, severe := range status.Errors.Severe {
		errs.AddSevere(severe)
	}
	if errs.IsSevere() {
		return errs
	}
	for _, warning := range status.Errors.Warning {
		errs.AddWarning(warning)
	}
	return errs
}

func validateUpdate(
	ctx context.Context,
	oldApi *v1alpha1.ApiDefinition,
	newApi *v1alpha1.ApiDefinition,
) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if newApi.IsBeingDeleted() {
		return errs
	}

	base.DeleteDefinitionConfigMapIfNeeded(ctx, oldApi, newApi)

	errs.Add(base.ValidateSubscribedPlans(ctx, oldApi, newApi, search.ApiV2SubsField))
	if errs.IsSevere() {
		return errs
	}
	errs.MergeWith(validateCreate(ctx, newApi))
	return errs
}
