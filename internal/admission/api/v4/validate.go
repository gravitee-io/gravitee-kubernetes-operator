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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func validateCreate(ctx context.Context, api *v1alpha1.ApiV4Definition) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	errs = validateFlowsAndEndpoints(api, errs)
	if errs.IsSevere() {
		return errs
	}

	errs = validateSharedPolicyGroups(ctx, api)
	if errs.IsSevere() {
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

func validateFlowsAndEndpoints(api *v1alpha1.ApiV4Definition,
	errs *errors.AdmissionErrors) *errors.AdmissionErrors {
	impl := &api.Spec.Api

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

	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateNativeKafkaPlanPorts(impl))

	return errs
}

func validateNativeKafkaPlanPorts(api *v4.Api) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if api == nil || api.Type != nativeAPI || api.Plans == nil {
		return errs
	}

	for hrid, plan := range *api.Plans {
		if plan == nil || plan.BootstrapPort == nil {
			continue
		}

		if plan.BrokerRangeStart == nil || plan.BrokerRangeEnd == nil {
			errs.AddSeveref(
				"plan [%s] has bootstrapPort set but brokerRangeStart/brokerRangeEnd are missing in Native API [%s]",
				hrid,
				api.Name,
			)
			continue
		}

		start := *plan.BrokerRangeStart
		end := *plan.BrokerRangeEnd
		bootstrap := *plan.BootstrapPort

		if start >= end {
			errs.AddSeveref(
				"plan [%s] has invalid broker port range (brokerRangeStart must be < brokerRangeEnd) in Native API [%s]",
				hrid,
				api.Name,
			)
			continue
		}

		if bootstrap >= start && bootstrap <= end {
			errs.AddSeveref(
				"plan [%s] has bootstrapPort within broker port range [%d,%d] in Native API [%s]",
				hrid,
				start,
				end,
				api.Name,
			)
		}
	}

	return errs
}

func validateDryRun(ctx context.Context, api *v1alpha1.ApiV4Definition) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	cp := api.DeepCopy()

	apimClient, err := apim.FromContextRef(ctx, cp.ContextRef(), cp.GetNamespace())
	if err != nil {
		errs.AddSevere(err.Error())
		return errs
	}

	cp.PopulateIDs(apimClient.Context, k8s.IsAutomationAPIManaged(api))
	cp.SetDefinitionContext(v4.NewDefaultKubernetesContext().MergeWith(cp.GetDefinitionContext()))
	status, err := apimClient.APIs.DryRunImportV4(cp)
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
	oldApi *v1alpha1.ApiV4Definition,
	newApi *v1alpha1.ApiV4Definition,
) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if newApi.IsBeingDeleted() {
		return errs
	}

	base.DeleteDefinitionConfigMapIfNeeded(ctx, oldApi, newApi)

	errs.Add(validateApiType(oldApi, newApi))
	if errs.IsSevere() {
		return errs
	}

	errs.Add(base.ValidateSubscribedPlans(ctx, oldApi, newApi, search.ApiV4SubsField))
	if errs.IsSevere() {
		return errs
	}

	errs.MergeWith(validateCreate(ctx, newApi))
	if errs.IsSevere() {
		return errs
	}

	mergeDriftValidation(ctx, oldApi, newApi, errs)
	return errs
}
