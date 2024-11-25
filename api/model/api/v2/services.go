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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

type Service struct {
	// Service name
	// +kubebuilder:validation:Optional
	Name *string `json:"name,omitempty"`
	// +kubebuilder:default:=false
	// +kubebuilder:validation:Optional
	// Is service enabled or not?
	Enabled bool `json:"enabled"`
}

type ScheduledService struct {
	*Service `json:",inline"`
	// +kubebuilder:validation:Optional
	Schedule *string `json:"schedule,omitempty"`
}

type EndpointDiscoveryService struct {
	*Service `json:",inline"`
	// Provider name
	// +kubebuilder:validation:Optional
	Provider *string `json:"provider,omitempty"`
	// Configuration, arbitrary map of key-values
	// +kubebuilder:validation:Optional
	Config *utils.GenericStringMap `json:"configuration,omitempty"`
	// Is it secondary or not?
	// +kubebuilder:validation:Optional
	Secondary *bool `json:"secondary,omitempty"`
	// List of tenants
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Tenants []string `json:"tenants"`
}

// +kubebuilder:validation:Enum=HTTP;
type DynamicPropertyProvider string

const (
	HttpPropertyProvider DynamicPropertyProvider = "HTTP"
)

type DynamicPropertyService struct {
	*ScheduledService `json:",inline"`
	Provider          DynamicPropertyProvider `json:"provider,omitempty"`
	// Configuration, arbitrary map of key-values
	// +kubebuilder:validation:Optional
	Config *utils.GenericStringMap `json:"configuration,omitempty"`
}

type Services struct {
	// Endpoint Discovery Service
	EndpointDiscoveryService *EndpointDiscoveryService `json:"discovery,omitempty"`
	// Health Check Service
	HealthCheckService *HealthCheckService `json:"health-check,omitempty"`
	// Dynamic Property Service
	DynamicPropertyService *DynamicPropertyService `json:"dynamic-property,omitempty"`
}

type HealthCheckService struct {
	*ScheduledService `json:",inline"`
	// List of health check steps
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Steps []*HealthCheckStep `json:"steps"`
}

type EndpointHealthCheckService struct {
	*HealthCheckService `json:",inline"`
	// Is service inherited or not?
	Inherit bool `json:"inherit"`
}

type HealthCheckStep struct {
	// Health Check Step Name
	// +kubebuilder:validation:Optional
	Name *string `json:"name,omitempty"`
	// Health Check Step Request
	Request HealthCheckRequest `json:"request,omitempty"`
	// Health Check Step Response
	Response HealthCheckResponse `json:"response,omitempty"`
}

type HealthCheckRequest struct {
	// The path of the endpoint handling the health check request
	// +kubebuilder:validation:Optional
	Path *string `json:"path,omitempty"`
	// The HTTP method to use when issuing the health check request
	Method base.HttpMethod `json:"method,omitempty"`
	// List of HTTP headers to include in the health check request
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Headers []base.HttpHeader `json:"headers"`
	// Health Check Request Body
	// +kubebuilder:validation:Optional
	Body *string `json:"body,omitempty"`
	// If true, the health check request will be issued without prepending the context path of the API.
	FromRoot bool `json:"fromRoot"`
}

type HealthCheckResponse struct {
	// +kubebuilder:validation:Optional
	Assertions []string `json:"assertions"`
}
