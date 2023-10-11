// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +kubebuilder:object:generate=true
package v4

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
)

// +kubebuilder:validation:Enum=PROXY;MESSAGE;
type ApiType string

const (
	ProxyType   = ApiType("proxy")
	MessageType = ApiType("message")
)

type Api struct {
	*base.ApiBase `json:",inline"`
	// +kubebuilder:default:=`4.0.0`
	// +kubebuilder:validation:Enum=`4.0.0`;
	DefinitionVersion base.DefinitionVersion `json:"definitionVersion,omitempty"`
	ApiVersion        string                 `json:"apiVersion,omitempty"`
	// +kubebuilder:validation:Required
	Type ApiType `json:"type,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems:=1
	Listeners []*Listener `json:"listeners,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems:=1
	EndpointGroups []*EndpointGroup `json:"endpointGroups,omitempty"`
	Plans          []*Plan          `json:"plans,omitempty"`
	FlowExecution  *FlowExecution   `json:"flowExecution,omitempty"`
	Flows          []*Flow          `json:"flows,omitempty"`
	Analytics      *Analytics       `json:"analytics,omitempty"`
	Services       *ApiServices     `json:"services,omitempty"`
}
