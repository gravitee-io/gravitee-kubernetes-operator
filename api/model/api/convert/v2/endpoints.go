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
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
)

// Used to convert back the typed endpoint group
// configuration of the v2 model to the generic
// implementation we have in the v4 model.
type sharedConfig interface {
	base.HttpClientOptions | base.HttpClientSslOptions | base.HttpProxy
}

const (
	httpKey  = "http"
	sslKey   = "ssl"
	proxyKey = "proxy"
)

func toEndpointGroups(v4EndpointGroups []*v4.EndpointGroup) []*v2.EndpointGroup {
	if v4EndpointGroups == nil {
		return []*v2.EndpointGroup{}
	}
	var endpointGroups []*v2.EndpointGroup
	for _, v4EndpointGroup := range v4EndpointGroups {
		endpointGroups = append(endpointGroups, toEndpointGroup(v4EndpointGroup))
	}
	return endpointGroups
}

func toEndpointGroup(v4EndpointGroup *v4.EndpointGroup) *v2.EndpointGroup {
	endpointGroup := v2.NewHttpEndpointGroup(v4EndpointGroup.Name)
	endpointGroup.Endpoints = toEndpoints(v4EndpointGroup.Endpoints)
	if v4EndpointGroup.LoadBalancer != nil {
		lb := toLoadBalancer(v4EndpointGroup.LoadBalancer)
		endpointGroup.LoadBalancer = *lb
	}
	endpointGroup.Services = toEndpointGroupServices(v4EndpointGroup.Services)
	endpointGroup.HttpClientOptions = toSharedConfigOpts[base.HttpClientOptions](v4EndpointGroup, httpKey)
	endpointGroup.HttpClientSslOptions = toSharedConfigOpts[base.HttpClientSslOptions](v4EndpointGroup, sslKey)
	endpointGroup.HttpProxy = toSharedConfigOpts[base.HttpProxy](v4EndpointGroup, proxyKey)
	return endpointGroup
}

func toSharedConfigOpts[OPTS sharedConfig](v4EndpointGroup *v4.EndpointGroup, key string) *OPTS {
	if v4EndpointGroup.SharedConfig == nil {
		return nil
	}
	http := v4EndpointGroup.SharedConfig.Get(key)
	if http == nil {
		return nil
	}
	httpClientOptions, ok := http.(map[string]interface{})
	if !ok {
		return nil
	}
	return fromSharedConfig[OPTS](httpClientOptions)
}

func fromSharedConfig[OUT sharedConfig](in map[string]interface{}) *OUT {
	if in == nil {
		return nil
	}
	config := new(OUT)
	b, err := json.Marshal(in)
	if err != nil {
		return nil
	}
	if err = json.Unmarshal(b, config); err != nil {
		return nil
	}
	return config
}

func toLoadBalancer(v4LB *v4.LoadBalancer) *v2.LoadBalancer {
	switch v4LB.Type {
	case v4.Random:
		return v2.NewLoadBalancer(v2.Random)
	case v4.WeightedRoundRobin:
		return v2.NewLoadBalancer(v2.WeightedRoundRobin)
	case v4.WeightedRandom:
		return v2.NewLoadBalancer(v2.WeightedRandom)
	case v4.RoundRobin:
		fallthrough
	default:
		return v2.NewLoadBalancer(v2.RoundRobin)
	}
}

func toEndpoints(v4Endpoints []*v4.Endpoint) []*v2.Endpoint {
	if v4Endpoints == nil {
		return []*v2.Endpoint{}
	}
	var endpoints []*v2.Endpoint
	for _, v4Endpoint := range v4Endpoints {
		endpoints = append(endpoints, toEndpoint(v4Endpoint))
	}
	return endpoints
}

func toEndpoint(v4Endpoint *v4.Endpoint) *v2.Endpoint {
	endpoint := v2.NewHttpEndpoint(v4Endpoint.Name)
	endpoint.Weight = v4Endpoint.Weight
	endpoint.Inherit = v4Endpoint.Inherit
	endpoint.Tenants = v4Endpoint.Tenants
	endpoint.HealthCheck = toEndpointHealthCheck(v4Endpoint.Services)
	target, ok := v4Endpoint.Config.Object["target"].(string)
	if ok {
		endpoint.Target = target
	}
	return endpoint
}
