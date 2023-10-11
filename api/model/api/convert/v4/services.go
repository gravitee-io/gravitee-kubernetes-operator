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
	"strings"

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

func toApiServices(v2Services *v2.Services) *v4.ApiServices {
	if v2Services == nil {
		return nil
	}
	return &v4.ApiServices{
		DynamicProperty: toDynamicPropertyService(v2Services.DynamicPropertyService),
	}
}

func toDynamicPropertyService(v2Service *v2.DynamicPropertyService) *v4.Service {
	if v2Service == nil {
		return nil
	}
	service := v4.NewService("http-dynamic-property", v2Service.Enabled)
	service.Config = v2Service.Config
	service.Config.Put("schedule", v2Service.Schedule)
	return service
}

func toEndpointGroupServices(v2Services *v2.Services) *v4.EndpointGroupServices {
	if v2Services == nil {
		return nil
	}
	return &v4.EndpointGroupServices{
		HealthCheck: toHealthCheckService(v2Services.HealthCheckService),
		Discovery:   toDiscovery(v2Services.EndpointDiscoveryService),
	}
}

func toHealthCheckService(v2Service *v2.HealthCheckService) *v4.Service {
	if v2Service == nil {
		return nil
	}
	service := v4.NewService("http-health-check", v2Service.Enabled)
	service.Config = utils.NewGenericStringMap()
	configureHealthCheck(v2Service, service.Config)
	return service
}

func toEndpointServices(v2Service *v2.EndpointHealthCheckService) *v4.EndpointServices {
	if v2Service == nil {
		return nil
	}
	service := toHealthCheckService(v2Service.HealthCheckService)
	service.OverrideConfig = !v2Service.Inherit
	return &v4.EndpointServices{
		HealthCheck: service,
	}
}

func configureHealthCheck(v2Service *v2.HealthCheckService, config *utils.GenericStringMap) {
	config.Put("schedule", v2Service.Schedule)
	config.Put("failureThreshold", 1)
	config.Put("successThreshold", 1)
	if len(v2Service.Steps) > 0 {
		step := v2Service.Steps[0]
		putRequest(step.Request, config)
		putResponse(step.Response, config)
	}
}

func putRequest(request v2.HealthCheckRequest, config *utils.GenericStringMap) {
	config.Put("method", request.Method)
	config.Put("target", request.Path)
	config.Put("headers", request.Headers)
	config.Put("body", request.Body)
	config.Put("overrideEndpointPath", request.FromRoot)
}

func putResponse(response v2.HealthCheckResponse, config *utils.GenericStringMap) {
	if len(response.Assertions) > 0 {
		config.Put("assertion", toAssertion(response.Assertions[0]))
	}
}

func toAssertion(v2Assertion string) string {
	return "{" + strings.TrimSuffix(strings.TrimPrefix(v2Assertion, "{"), "}") + "}"
}

func toDiscovery(v2Service *v2.EndpointDiscoveryService) *v4.Service {
	if v2Service == nil {
		return nil
	}
	service := v4.NewService(v2Service.Provider, v2Service.Enabled)
	service.Config = v2Service.Config
	return service
}
