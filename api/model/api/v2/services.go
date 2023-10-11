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
	Name string `json:"name,omitempty"`
	// +kubebuilder:default:=false
	Enabled bool `json:"enabled"`
}

type ScheduledService struct {
	*Service `json:",inline"`
	Schedule string `json:"schedule,omitempty"`
}

type EndpointDiscoveryService struct {
	*Service  `json:",inline"`
	Provider  string                  `json:"provider,omitempty"`
	Config    *utils.GenericStringMap `json:"configuration,omitempty"`
	Secondary bool                    `json:"secondary,omitempty"`
	Tenants   []string                `json:"tenants,omitempty"`
}

// +kubebuilder:validation:Enum=HTTP;
type DynamicPropertyProvider string

const (
	HttpPropertyProvider DynamicPropertyProvider = "HTTP"
)

type DynamicPropertyService struct {
	*ScheduledService `json:",inline"`
	Provider          DynamicPropertyProvider `json:"provider,omitempty"`
	Config            *utils.GenericStringMap `json:"configuration,omitempty"`
}

type Services struct {
	EndpointDiscoveryService *EndpointDiscoveryService `json:"discovery,omitempty"`
	HealthCheckService       *HealthCheckService       `json:"health-check,omitempty"`
	DynamicPropertyService   *DynamicPropertyService   `json:"dynamic-property,omitempty"`
}

type HealthCheckService struct {
	*ScheduledService `json:",inline"`
	Steps             []*HealthCheckStep `json:"steps,omitempty"`
}

type EndpointHealthCheckService struct {
	*HealthCheckService `json:",inline"`
	Inherit             bool `json:"inherit"`
}

type HealthCheckStep struct {
	Name     string              `json:"name,omitempty"`
	Request  HealthCheckRequest  `json:"request,omitempty"`
	Response HealthCheckResponse `json:"response,omitempty"`
}

type HealthCheckRequest struct {
	Path     string            `json:"path,omitempty"`
	Method   base.HttpMethod   `json:"method,omitempty"`
	Headers  []base.HttpHeader `json:"headers,omitempty"`
	Body     string            `json:"body,omitempty"`
	FromRoot bool              `json:"fromRoot"`
}

type HealthCheckResponse struct {
	Assertions []string `json:"assertions,omitempty"`
}
