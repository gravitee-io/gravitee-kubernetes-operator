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

// +kubebuilder:object:generate=true
package v2

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
)

type Api struct {
	*base.ApiBase `json:",inline"`
	// Shows the time that the API is deployed
	DeployedAt uint64 `json:"deployedAt,omitempty"`
	// +kubebuilder:default:=`2.0.0`
	// The definition version of the API. For v1alpha1 resources, this field should always set to `2.0.0`.
	DefinitionVersion base.DefinitionVersion `json:"gravitee,omitempty"`
	// +kubebuilder:validation:Required
	// API version
	Version string `json:"version,omitempty"`
	// +kubebuilder:default:=DEFAULT
	// The flow mode of the API. The value is either `DEFAULT` or `BEST_MATCH`.
	FlowMode FlowMode `json:"flow_mode,omitempty"`
	// The proxy of the API that specifies its VirtualHosts and Groups.
	Proxy *Proxy `json:"proxy,omitempty"`
	// Contains different services for the API (EndpointDiscovery, HealthCheck ...)
	Services *Services `json:"services,omitempty"`
	// +kubebuilder:validation:Optional
	// The flow of the API
	Flows []Flow `json:"flows"`
	// +kubebuilder:validation:Optional
	// API Path mapping
	PathMappings []string `json:"path_mappings"`
	// +kubebuilder:validation:Optional
	// API plans
	Plans []*Plan `json:"plans"`
}
