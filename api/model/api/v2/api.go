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
	// The definition context is used to inform a management API instance that this API definition
	// is managed using a kubernetes operator
	DefinitionContext *DefinitionContext `json:"definition_context,omitempty"`
	// +kubebuilder:default:=DEFAULT
	// The flow mode of the API. The value is either `DEFAULT` or `BEST_MATCH`.
	FlowMode FlowMode `json:"flow_mode,omitempty"`
	// The proxy of the API that specifies its VirtualHosts and Groups.
	Proxy *Proxy `json:"proxy,omitempty"`
	// Contains different services for the API (EndpointDiscovery, HealthCheck ...)
	Services *Services `json:"services,omitempty"`
	// +kubebuilder:default:={}
	// The flow of the API
	Flows []Flow `json:"flows,omitempty"`
	// API Path mapping
	PathMappings []string `json:"path_mappings,omitempty"`
	// +kubebuilder:default:={}
	// API plans
	Plans []*Plan `json:"plans,omitempty"`
	// local defines if the api is local or not.
	//
	// If true, the Operator will create the ConfigMaps for the Gateway and pushes the API to the Management API
	// but without setting the update flag in the datastore.
	//
	// If false, the Operator will not create the ConfigMaps for the Gateway.
	// Instead, it pushes the API to the Management API and forces it to update the event in the datastore.
	// This will cause Gateways to fetch the APIs from the datastore
	//
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	IsLocal bool `json:"local"`
}

const (
	ModeFullyManaged = "fully_managed"
	OriginKubernetes = "kubernetes"
)

type DefinitionContext struct {
	// +kubebuilder:default:=kubernetes
	Origin string `json:"origin,omitempty"`
	// +kubebuilder:default:=fully_managed
	Mode string `json:"mode,omitempty"`
}