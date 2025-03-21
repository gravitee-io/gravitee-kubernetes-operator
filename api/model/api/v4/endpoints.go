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
	nameGen "github.com/moby/moby/pkg/namesgenerator"
)

type EndpointType string

const (
	EndpointTypeHTTP  = EndpointType("http-proxy")
	EndpointTypeKafka = EndpointType("native-kafka")
)

type Endpoint struct {
	// The endpoint name (this value should be unique across endpoints)
	// +kubebuilder:validation:Optional
	Name *string `json:"name,omitempty"`

	// +kubebuilder:validation:Required
	// Endpoint Type
	Type string `json:"type,omitempty"`

	// Endpoint Weight
	// +kubebuilder:validation:Optional
	Weight *int `json:"weight,omitempty"`

	// Should endpoint group configuration be inherited or not ?
	Inherit bool `json:"inheritConfiguration"`

	// Endpoint Configuration, arbitrary map of key-values
	// +kubebuilder:validation:Optional
	Config *utils.GenericStringMap `json:"configuration,omitempty"`

	// Endpoint Configuration Override, arbitrary map of key-values
	// +kubebuilder:validation:Optional
	ConfigOverride *utils.GenericStringMap `json:"sharedConfigurationOverride,omitempty"`

	// Endpoint Services
	Services *EndpointServices `json:"services,omitempty"`

	// Endpoint is secondary or not?
	Secondary bool `json:"secondary"`

	// List of endpoint tenants
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Tenants []string `json:"tenants"`
}

func NewHttpEndpoint(name string) *Endpoint {
	return &Endpoint{
		Name:   &name,
		Type:   string(EndpointTypeHTTP),
		Config: utils.NewGenericStringMap(),
	}
}

func NewKafkaEndpoint(name string) *Endpoint {
	return &Endpoint{
		Name:   &name,
		Type:   string(EndpointTypeKafka),
		Config: utils.NewGenericStringMap(),
	}
}

// +kubebuilder:validation:Enum=ROUND_ROBIN;RANDOM;WEIGHTED_ROUND_ROBIN;WEIGHTED_RANDOM;
type LoadBalancerType string

func (lt LoadBalancerType) toGatewayDefinition() LoadBalancerType {
	return LoadBalancerType(Enum(lt).ToGatewayDefinition())
}

const (
	RoundRobin         LoadBalancerType = "ROUND_ROBIN"
	Random             LoadBalancerType = "RANDOM"
	WeightedRoundRobin LoadBalancerType = "WEIGHTED_ROUND_ROBIN"
	WeightedRandom     LoadBalancerType = "WEIGHTED_RANDOM"
)

type LoadBalancer struct {
	// +kubebuilder:default:=`ROUND_ROBIN`
	Type LoadBalancerType `json:"type"`
}

func NewLoadBalancer(algo LoadBalancerType) *LoadBalancer {
	return &LoadBalancer{
		Type: algo,
	}
}

type EndpointGroup struct {
	// +kubebuilder:validation:Required
	// Endpoint group name
	Name string `json:"name"`
	// Endpoint group type
	// +kubebuilder:validation:Optional
	Type *string `json:"type,omitempty"`
	// Endpoint group load balancer
	LoadBalancer *LoadBalancer `json:"loadBalancer,omitempty"`
	// Endpoint group shared configuration, arbitrary map of key-values
	// +kubebuilder:validation:Optional
	SharedConfig *utils.GenericStringMap `json:"sharedConfiguration,omitempty"`
	// List of endpoint for the group
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Endpoints []*Endpoint `json:"endpoints"`
	// Endpoint group services
	Services *EndpointGroupServices `json:"services,omitempty"`
	// Endpoint group http client options
	HttpClientOptions *base.HttpClientOptions `json:"http,omitempty"`
	// Endpoint group http client SSL options
	HttpClientSslOptions *base.HttpClientSslOptions `json:"ssl,omitempty"`
	// Endpoint group headers, arbitrary map of key-values
	// +kubebuilder:validation:Optional
	Headers *map[string]string `json:"headers,omitempty"`
}

func NewHttpEndpointGroup(name string) *EndpointGroup {
	t := string(EndpointTypeHTTP)
	return &EndpointGroup{
		Name:         name,
		Type:         &t,
		Endpoints:    []*Endpoint{},
		SharedConfig: utils.NewGenericStringMap(),
	}
}

func NewKafkaEndpointGroup(name string) *EndpointGroup {
	t := string(EndpointTypeKafka)
	return &EndpointGroup{
		Name:         name,
		Type:         &t,
		Endpoints:    []*Endpoint{},
		SharedConfig: utils.NewGenericStringMap(),
	}
}

// If the API has been converted from a v1alpha1 version, the name might be empty
// Because a name is required by the GW for v4 API deserialization, we generate a random name
// using the docker name generator.
func (e EndpointGroup) ToGatewayDefinition() *EndpointGroup {
	if e.Name == "" {
		e.Name = nameGen.GetRandomName(0)
	}
	if e.LoadBalancer != nil {
		e.LoadBalancer.Type = e.LoadBalancer.Type.toGatewayDefinition()
	}
	return &e
}
