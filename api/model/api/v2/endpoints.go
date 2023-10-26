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
	Name    string         `json:"name,omitempty"`
	Target  string         `json:"target,omitempty"`
	Weight  int            `json:"weight,omitempty"`
	Backup  bool           `json:"backup,omitempty"`
	Status  EndpointStatus `json:"-,omitempty"`
	Tenants []string       `json:"tenants,omitempty"`
	Type    EndpointType   `json:"type,omitempty"`
	Inherit bool           `json:"inherit,omitempty"`

	HttpProxy            *base.HttpProxy             `json:"proxy,omitempty"`
	HttpClientOptions    *base.HttpClientOptions     `json:"http,omitempty"`
	HttpClientSslOptions *base.HttpClientSslOptions  `json:"ssl,omitempty"`
	Headers              []base.HttpHeader           `json:"headers,omitempty"`
	HealthCheck          *EndpointHealthCheckService `json:"healthcheck,omitempty"`
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
	Type LoadBalancerType `json:"type,omitempty"`
}

func NewLoadBalancer(algo LoadBalancerType) *LoadBalancer {
	return &LoadBalancer{
		Type: algo,
	}
}

type EndpointGroup struct {
	Name                 string                     `json:"name,omitempty"`
	Endpoints            []*Endpoint                `json:"endpoints,omitempty"`
	LoadBalancer         LoadBalancer               `json:"load_balancing,omitempty"`
	Services             *Services                  `json:"services,omitempty"`
	HttpProxy            *base.HttpProxy            `json:"proxy,omitempty"`
	HttpClientOptions    *base.HttpClientOptions    `json:"http,omitempty"`
	HttpClientSslOptions *base.HttpClientSslOptions `json:"ssl,omitempty"`
	Headers              map[string]string          `json:"headers,omitempty"`
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
	MaxAttempts  int            `json:"maxAttempts,omitempty"`
	RetryTimeout int64          `json:"retryTimeout,omitempty"`
	Cases        []FailoverCase `json:"cases,omitempty"`
}
