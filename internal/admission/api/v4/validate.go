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

package v4

import (
	"context"

	baseModel "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func validateCreate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	if api, ok := obj.(core.ApiDefinitionResource); ok {
		errs = errs.MergeWith(base.ValidateCreate(ctx, obj))
		if errs.IsSevere() {
			return errs
		}

		errs.Add(validateApiPlans(ctx, api))
		if errs.IsSevere() {
			return errs
		}

		if errs.IsSevere() {
			return errs
		}
		if api.HasContext() {
			errs = errs.MergeWith(validateDryRun(ctx, api))
		}
	}
	return errs
}

func validateApiPlans(_ context.Context, api core.ApiDefinitionResource) *errors.AdmissionError {
	cp, _ := api.DeepCopyResource().(core.ApiDefinitionResource)

	apiDef, ok := cp.GetDefinition().(*v4.Api)
	if !ok {
		return errors.NewSevere("unable to validate the CRD because it is not a v4 API")
	}

	if apiDef.State == baseModel.StateStarted &&
		len(apiDef.Plans) == 0 {
		return errors.NewSevere("cannot apply API [%s]. Its state is set to STARTED,"+
			" but the API has no plans. APIs must have at least one plan in order to be deployed.", apiDef.Name)
	}

	return nil
}

func validateDryRun(ctx context.Context, api core.ApiDefinitionResource) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp, _ := api.DeepCopyResource().(core.ApiDefinitionResource)

	apim, err := apim.FromContextRef(ctx, cp.ContextRef(), cp.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
	}

	cp.PopulateIDs(apim.Context)
	cp.SetDefinitionContext(v4.NewDefaultKubernetesContext().MergeWith(cp.GetDefinitionContext()))

	impl, ok := cp.GetDefinition().(*v4.Api)
	if !ok {
		errs.AddSevere("unable to call dry run import because api is not a v4 API")
	}

	status, err := apim.APIs.DryRunImportV4(impl)
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
