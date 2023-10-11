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
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

func toApiServices(v4services *v4.ApiServices) *v2.Services {
	if v4services == nil {
		return nil
	}
	return &v2.Services{
		DynamicPropertyService: toDynamicPropertyService(v4services.DynamicProperty),
	}
}

func toEndpointGroupServices(v4services *v4.EndpointGroupServices) *v2.Services {
	if v4services == nil {
		return nil
	}
	return &v2.Services{
		EndpointDiscoveryService: toEndpointDiscoveryService(v4services.Discovery),
		HealthCheckService:       toHealthCheckService(v4services.HealthCheck),
	}
}

func toDynamicPropertyService(v4Service *v4.Service) *v2.DynamicPropertyService {
	if v4Service == nil {
		return nil
	}

	return &v2.DynamicPropertyService{
		ScheduledService: toScheduleService(v4Service),
		Config:           v4Service.Config,
		Provider:         v2.HttpPropertyProvider,
	}
}

func toEndpointHealthCheck(v4Services *v4.EndpointServices) *v2.EndpointHealthCheckService {
	if v4Services == nil || v4Services.HealthCheck == nil {
		return nil
	}

	healthCheck := v4Services.HealthCheck

	return &v2.EndpointHealthCheckService{
		HealthCheckService: toHealthCheckService(healthCheck),
		Inherit:            !healthCheck.OverrideConfig,
	}
}

func toHealthCheckService(v4Service *v4.Service) *v2.HealthCheckService {
	if v4Service == nil {
		return nil
	}

	return &v2.HealthCheckService{
		ScheduledService: toScheduleService(v4Service),
		Steps:            toHealthCheckSteps(v4Service.Config),
	}
}

func toHealthCheckSteps(v4Config *utils.GenericStringMap) []*v2.HealthCheckStep {
	return []*v2.HealthCheckStep{toHealthCheckStep(v4Config)}
}

func toHealthCheckStep(v4Config *utils.GenericStringMap) *v2.HealthCheckStep {
	return &v2.HealthCheckStep{
		Request:  getRequest(v4Config),
		Response: getResponse(v4Config),
	}
}

func getRequest(v4Config *utils.GenericStringMap) v2.HealthCheckRequest {
	return v2.HealthCheckRequest{
		Method:   base.HttpMethod(v4Config.GetString("method")),
		Path:     v4Config.GetString("target"),
		Headers:  getHttpHeaders(v4Config),
		Body:     v4Config.GetString("body"),
		FromRoot: v4Config.GetBool("overrideEndpointPath"),
	}
}

func getHttpHeaders(config *utils.GenericStringMap) []base.HttpHeader {
	sl := config.GetSlice("headers")
	if sl == nil {
		return nil
	}

	headers := make([]base.HttpHeader, 0)

	for _, v := range sl {
		if m, ok := v.(map[string]interface{}); ok {
			headers = append(headers, base.HttpHeader{
				Name:  m["name"].(string),
				Value: m["value"].(string),
			})
		}
	}

	return headers
}

func getResponse(v4Config *utils.GenericStringMap) v2.HealthCheckResponse {
	return v2.HealthCheckResponse{
		Assertions: []string{getAssertion(v4Config)},
	}
}

func getAssertion(v4Config *utils.GenericStringMap) string {
	assertion := v4Config.GetString("assertion")
	return strings.TrimSuffix(strings.TrimPrefix(assertion, "{"), "}")
}

func toScheduleService(v4Service *v4.Service) *v2.ScheduledService {
	if v4Service.Config == nil {
		return nil
	}
	schedule := v4Service.Config.GetString("schedule")
	if schedule == "" {
		return nil
	}

	v4Service.Config.Remove("schedule")

	return &v2.ScheduledService{
		Service:  toService(v4Service),
		Schedule: schedule,
	}
}

func toService(v4Service *v4.Service) *v2.Service {
	return &v2.Service{
		Enabled: v4Service.Enabled,
	}
}

func toEndpointDiscoveryService(v4Service *v4.Service) *v2.EndpointDiscoveryService {
	if v4Service == nil || v4Service.Config == nil {
		return nil
	}

	return &v2.EndpointDiscoveryService{
		Provider: v4Service.Config.GetString("provider"),
		Service:  toService(v4Service),
		Config:   v4Service.Config,
	}
}
