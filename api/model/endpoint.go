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

package model

type EndpointStatus int

const (
	Down EndpointStatus = iota
	TransitionallyDown
	TransitionallyUp
	Up
)

type EndpointType string

const (
	HttpEndpointType EndpointType = "http"
	GrpcEndpointType EndpointType = "grpc"
)

type Endpoint struct {
	Name         string         `json:"name,omitempty"`
	Target       string         `json:"target,omitempty"`
	Weight       int            `json:"weight,omitempty"`
	Backup       bool           `json:"backup,omitempty"`
	Status       EndpointStatus `json:"-,omitempty"`
	Tenants      []string       `json:"tenants,omitempty"`
	EndpointType EndpointType   `json:"type,omitempty"`
	Inherit      bool           `json:"inherit,omitempty"`
}

type EndpointHealthCheckService struct {
	Inherit bool `json:"inherit,omitempty"`

	// HealthCheckService
	Steps    []Step `json:"steps,omitempty"`
	Schedule string `json:"schedule,omitempty"`
}

func NewEndpointHealthCheckService() *EndpointHealthCheckService {
	return &EndpointHealthCheckService{Schedule: "health-check"}
}

type HttpEndpoint struct {
	// From Endpoint
	Name         string         `json:"name,omitempty"`
	Target       string         `json:"target,omitempty"`
	Weight       int            `json:"weight,omitempty"`
	Backup       bool           `json:"backup,omitempty"`
	Status       EndpointStatus `json:"-,omitempty"`
	Tenants      []string       `json:"tenants,omitempty"`
	EndpointType EndpointType   `json:"type,omitempty"`
	Inherit      bool           `json:"inherit,omitempty"`

	HttpProxy            *HttpProxy                  `json:"proxy,omitempty"`
	HttpClientOptions    *HttpClientOptions          `json:"http,omitempty"`
	HttpClientSslOptions *HttpClientSslOptions       `json:"ssl,omitempty"`
	Headers              []HttpHeader                `json:"headers,omitempty"`
	HealthCheck          *EndpointHealthCheckService `json:"healthCheck,omitempty"`
}

type EndpointDiscoveryService struct {
	Name          string            `json:"name,omitempty"`
	Enabled       bool              `json:"enabled,omitempty"`
	Provider      string            `json:"provider,omitempty"`
	Configuration *GenericStringMap `json:"configuration,omitempty"`
}

type DynamicPropertyProvider int

const (
	HttpPropertyProvider DynamicPropertyProvider = iota
)

type DynamicPropertyService struct {
	Schedule string                  `json:"schedule,omitempty"`
	Provider DynamicPropertyProvider `json:"provider,omitempty"`
	// Configuration DynamicPropertyProviderConfiguration `json:"configuration,omitempty"`  // needs to be fixed later
}

func NewDynamicPropertyService() *DynamicPropertyService {
	return &DynamicPropertyService{
		Schedule: "dynamic-property",
	}
}

type Service struct {
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

func NewService() *Service {
	return &Service{
		Enabled: true,
	}
}

type Services struct {
	EndpointDiscoveryService *EndpointDiscoveryService `json:"discovery,omitempty"`
	HealthCheckService       *HealthCheckService       `json:"health-check,omitempty"`
	DynamicPropertyService   *DynamicPropertyService   `json:"dynamic-property,omitempty"`
}

type LoadBalancerType string

const (
	RoundRobin         LoadBalancerType = "ROUND_ROBIN"
	Random             LoadBalancerType = "RANDOM"
	WeightedRoundRobin LoadBalancerType = "WEIGHTED_ROUND_ROBIN"
	WeightedRandom     LoadBalancerType = "WEIGHTED_RANDOM"
)

type LoadBalancer struct {
	LoadBalancerType LoadBalancerType `json:"type,omitempty"`
}

type EndpointGroup struct {
	Name                 string                `json:"name,omitempty"`
	Endpoints            []*HttpEndpoint       `json:"endpoints,omitempty"`
	LoadBalancer         LoadBalancer          `json:"load_balancing,omitempty"`
	Services             *Services             `json:"services,omitempty"`
	HttpProxy            *HttpProxy            `json:"proxy,omitempty"`
	HttpClientOptions    *HttpClientOptions    `json:"http,omitempty"`
	HttpClientSslOptions *HttpClientSslOptions `json:"ssl,omitempty"`
	Headers              map[string]string     `json:"headers,omitempty"`
}
