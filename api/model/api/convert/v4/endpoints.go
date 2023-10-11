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
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

func toEndpointGroups(proxy *v2.Proxy) []*v4.EndpointGroup {
	if proxy == nil || proxy.Groups == nil {
		return nil
	}

	var endpointGroups []*v4.EndpointGroup
	for _, group := range proxy.Groups {
		endpointGroups = append(endpointGroups, toEndpointGroup(group))
	}
	return endpointGroups
}

func toEndpointGroup(v2Group *v2.EndpointGroup) *v4.EndpointGroup {
	endpointGroup := v4.NewHttpEndpointGroup(v2Group.Name)
	endpointGroup.LoadBalancer = toLoadBalancer(v2Group.LoadBalancer)
	endpointGroup.Services = toEndpointGroupServices(v2Group.Services)
	for _, v3Endpoint := range v2Group.Endpoints {
		endpointGroup.Endpoints = append(endpointGroup.Endpoints, toEndpoint(v3Endpoint))
	}
	return endpointGroup
}

func toLoadBalancer(v2LB v2.LoadBalancer) *v4.LoadBalancer {
	switch v2LB.Type {
	case v2.RoundRobin:
		return v4.NewLoadBalancer(v4.RoundRobin)
	case v2.Random:
		return v4.NewLoadBalancer(v4.Random)
	case v2.WeightedRoundRobin:
		return v4.NewLoadBalancer(v4.WeightedRoundRobin)
	case v2.WeightedRandom:
		return v4.NewLoadBalancer(v4.WeightedRandom)
	default:
		return nil
	}
}

func toEndpoint(v2Endpoint *v2.Endpoint) *v4.Endpoint {
	endpoint := v4.NewHttpEndpoint(v2Endpoint.Name)
	endpoint.Weight = v2Endpoint.Weight
	endpoint.Inherit = v2Endpoint.Inherit
	endpoint.Tenants = v2Endpoint.Tenants
	endpoint.Config = utils.NewGenericStringMap()
	endpoint.Config.Object["target"] = v2Endpoint.Target
	endpoint.Services = toEndpointServices(v2Endpoint.HealthCheck)
	return endpoint
}
