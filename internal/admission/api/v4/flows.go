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
	apiV4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

const nativeAPI = "NATIVE"

func validateApiFlows(api *apiV4.Api) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	// Validate API Flows
	errs.MergeWith(validateFlow(api.Type, api.Name, api.Flows))

	if errs.IsSevere() {
		return errs
	}

	// Validate API Plan Flows
	if api.Plans != nil {
		errs.MergeWith(validateFlow(api.Type, api.Name, api.Flows))
	}

	return errs
}

func validateFlow(apiType apiV4.ApiType, apiName string, flows []*apiV4.Flow) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if apiType == nativeAPI { //nolint:nestif // normal complexity
		for _, flow := range flows {
			if len(flow.Request) > 0 {
				errs.AddSeveref("Request Flow is not supported in Native API [%s]", apiName)
			}
			if len(flow.Response) > 0 {
				errs.AddSeveref("Response Flow is not supported in Native API [%s]", apiName)
			}
		}
	} else {
		for _, flow := range flows {
			if len(flow.Connect) > 0 {
				errs.AddSeveref("Connect Flow is not supported in V4 API [%s]", apiName)
			}
			if len(flow.Interact) > 0 {
				errs.AddSeveref("Interact Flow is not supported in V4 API [%s]", apiName)
			}
		}
	}

	return errs
}

func validateApiEndpointGroups(api *apiV4.Api) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if api.Type == nativeAPI {
		for _, eg := range api.EndpointGroups {
			if eg.Services != nil {
				errs.AddSeveref("EndpointGroup services is not supported in Native API [%s]", api.Name)
			}
		}
	}

	return errs
}

func validateApiFlowExecution(api *apiV4.Api) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if api.Type == nativeAPI && api.FlowExecution != nil {
		errs.AddSeveref("FlowExecution is not supported in Native API [%s]", api.Name)
	}

	return errs
}

func validateApiResponseTemplates(api *apiV4.Api) *errors.AdmissionErrors {
	errs := errors.NewAdmissionErrors()

	if api.Type == nativeAPI && api.ResponseTemplates != nil {
		errs.AddSeveref("ResponseTemplates is not supported in Native API [%s]", api.Name)
	}

	return errs
}
