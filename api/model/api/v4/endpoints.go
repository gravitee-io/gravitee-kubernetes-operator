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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

type EndpointType string

const (
	EndpointTypeHTTP = EndpointType("http-proxy")
)

type Endpoint struct {
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Required
	Type           string                  `json:"type,omitempty"`
	Weight         int                     `json:"weight,omitempty"`
	Inherit        bool                    `json:"inheritConfiguration"`
	Config         *utils.GenericStringMap `json:"configuration,omitempty"`
	ConfigOverride *utils.GenericStringMap `json:"sharedConfigurationOverride,omitempty"`
	Services       *EndpointServices       `json:"services,omitempty"`
	Secondary      bool                    `json:"secondary"`
	Tenants        []string                `json:"tenants,omitempty"`
}

func NewHttpEndpoint(name string) *Endpoint {
	return &Endpoint{
		Name: name,
		Type: string(EndpointTypeHTTP),
	}
}

type LoadBalancerType string

const (
	RoundRobin         LoadBalancerType = "round-robin"
	Random             LoadBalancerType = "random"
	WeightedRoundRobin LoadBalancerType = "weighted-round-robin"
	WeightedRandom     LoadBalancerType = "weighted-random"
)

type LoadBalancer struct {
	Type LoadBalancerType `json:"type,omitempty"`
}

func NewLoadBalancer(algo LoadBalancerType) *LoadBalancer {
	return &LoadBalancer{
		Type: algo,
	}
}

type EndpointGroup struct {
	Name                 string                     `json:"name,omitempty"`
	Type                 string                     `json:"type,omitempty"`
	LoadBalancer         *LoadBalancer              `json:"loadBalancer,omitempty"`
	SharedConfig         *utils.GenericStringMap    `json:"sharedConfiguration,omitempty"`
	Endpoints            []*Endpoint                `json:"endpoints,omitempty"`
	Services             *EndpointGroupServices     `json:"services,omitempty"`
	HttpClientOptions    *base.HttpClientOptions    `json:"http,omitempty"`
	HttpClientSslOptions *base.HttpClientSslOptions `json:"ssl,omitempty"`
	Headers              map[string]string          `json:"headers,omitempty"`
}

func NewHttpEndpointGroup(name string) *EndpointGroup {
	return &EndpointGroup{
		Name: name,
		Type: string(EndpointTypeHTTP),
	}
}
