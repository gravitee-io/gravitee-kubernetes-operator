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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
)

type EndpointType string

const (
	HttpEndpointType EndpointType = "http"
	GrpcEndpointType EndpointType = "grpc"
)

type Endpoint struct {
	// Name of the endpoint
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// The end target of this endpoint (backend)
	// +kubebuilder:validation:Optional
	Target string `json:"target,omitempty"`
	// Endpoint weight used for load-balancing
	// +kubebuilder:validation:Optional
	Weight int `json:"weight,omitempty"`
	// Indicate that this ia a back-end endpoint
	// +kubebuilder:validation:Optional
	Backup bool `json:"backup,omitempty"`
	// The endpoint tenants
	// +kubebuilder:validation:Optional
	Tenants []string `json:"tenants"`
	// The type of endpoint (HttpEndpointType or GrpcEndpointType)
	Type EndpointType `json:"type,omitempty"`
	// Is endpoint inherited or not
	// +kubebuilder:validation:Optional
	Inherit bool `json:"inherit,omitempty"`
	// Configure the HTTP Proxy settings to reach target if needed
	HttpProxy *base.HttpProxy `json:"proxy,omitempty"`
	// Custom HTTP client options used for this endpoint
	HttpClientOptions *base.HttpClientOptions `json:"http,omitempty"`
	// Custom HTTP SSL client options used for this endpoint
	HttpClientSslOptions *base.HttpClientSslOptions `json:"ssl,omitempty"`
	// List of headers for this endpoint
	// +kubebuilder:validation:Optional
	Headers []base.HttpHeader `json:"headers"`
	// Specify EndpointHealthCheck service settings
	HealthCheck *EndpointHealthCheckService `json:"healthcheck,omitempty"`
}

func NewHttpEndpoint(name string) *Endpoint {
	return &Endpoint{
		Type: HttpEndpointType,
		Name: name,
	}
}

type LoadBalancerType string

const (
	RoundRobin         LoadBalancerType = "ROUND_ROBIN"
	Random             LoadBalancerType = "RANDOM"
	WeightedRoundRobin LoadBalancerType = "WEIGHTED_ROUND_ROBIN"
	WeightedRandom     LoadBalancerType = "WEIGHTED_RANDOM"
)

type LoadBalancer struct {
	// Type of the LoadBalancer (RoundRobin, Random, WeightedRoundRobin, WeightedRandom)
	Type LoadBalancerType `json:"type,omitempty"`
}

func NewLoadBalancer(algo LoadBalancerType) *LoadBalancer {
	return &LoadBalancer{
		Type: algo,
	}
}

type EndpointGroup struct {
	// EndpointGroup name
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// List of Endpoints belonging to this group
	// +kubebuilder:validation:Optional
	Endpoints []*Endpoint `json:"endpoints"`
	// The LoadBalancer Type
	LoadBalancer LoadBalancer `json:"load_balancing,omitempty"`
	// Specify different Endpoint Services
	Services *Services `json:"services,omitempty"`
	// Configure the HTTP Proxy settings for this EndpointGroup if needed
	HttpProxy *base.HttpProxy `json:"proxy,omitempty"`
	// Custom HTTP SSL client options used for this EndpointGroup
	HttpClientOptions *base.HttpClientOptions `json:"http,omitempty"`
	// Custom HTTP SSL client options used for this EndpointGroup
	HttpClientSslOptions *base.HttpClientSslOptions `json:"ssl,omitempty"`
	// List of headers needed for this EndpointGroup
	// +kubebuilder:validation:Optional
	Headers map[string]string `json:"headers,omitempty"`
}

func NewHttpEndpointGroup(name string) *EndpointGroup {
	return &EndpointGroup{
		Name:      name,
		Endpoints: []*Endpoint{},
		Headers:   map[string]string{},
	}
}

type FailoverCase string

type Failover struct {
	// Maximum number of attempts
	// +kubebuilder:validation:Optional
	MaxAttempts int `json:"maxAttempts,omitempty"`
	// Retry timeout
	// +kubebuilder:validation:Optional
	RetryTimeout int64 `json:"retryTimeout,omitempty"`
	// List of Failover cases
	// +kubebuilder:validation:Optional
	Cases []FailoverCase `json:"cases"`
}
