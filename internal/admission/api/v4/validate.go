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

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"k8s.io/apimachinery/pkg/runtime"
)

func validateCreate(ctx context.Context, obj runtime.Object) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	api, ok := obj.(core.ApiDefinitionObject)
	if !ok {
		return errs
	}

	errs = validateFlowsAndEndpoints(ctx, api, errs)
	if errs.IsSevere() {
		return errs
	}

	errs = validateSharedPolicyGroups(ctx, api)
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(base.ValidateCreate(ctx, obj))
	if errs.IsSevere() {
		return errs
	}

	if api.HasContext() {
		errs.MergeWith(validateDryRun(ctx, api))
	}

	return errs
}

func validateFlowsAndEndpoints(_ context.Context, api core.ApiDefinitionObject,
	errs *errors.AdmissionErrors) *errors.AdmissionErrors {
	cp, _ := api.DeepCopyObject().(core.ApiDefinitionObject)
	impl, ok := cp.GetDefinition().(*v4.Api)
	if !ok {
		errs.AddSevere("unable to API type because api is not a v4 API")
	}

	errs.MergeWith(validateApiFlows(impl))

	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateApiEndpointGroups(impl))

	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateApiFlowExecution(impl))

	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateApiResponseTemplates(impl))

	return errs
}

func validateDryRun(ctx context.Context, api core.ApiDefinitionObject) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp, _ := api.DeepCopyObject().(core.ApiDefinitionObject)

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

func validateUpdate(
	ctx context.Context,
	oldObj runtime.Object,
	newObj runtime.Object,
) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()
	oldApi, ook := oldObj.(core.ApiDefinitionObject)
	newApi, nok := newObj.(core.ApiDefinitionObject)
	if !ook || !nok {
		return errs
	}

	base.DeleteDefinitionConfigMapIfNeeded(ctx, oldApi, newApi)

	errs.Add(validateApiType(oldApi, newApi))
	if errs.IsSevere() {
		return errs
	}

	errs.Add(base.ValidateSubscribedPlans(ctx, oldApi, newApi, indexer.ApiV4SubsField))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateCreate(ctx, newApi))
	return errs
}
