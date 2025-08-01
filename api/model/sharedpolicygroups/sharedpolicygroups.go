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

package sharedpolicygroups

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

// +kubebuilder:validation:Enum=MESSAGE;PROXY;NATIVE;
type ApiType string

// +kubebuilder:validation:Enum=REQUEST;RESPONSE;INTERACT;CONNECT;PUBLISH;SUBSCRIBE;
type FlowPhase string

type SharedPolicyGroup struct {
	// CrossID to export SharedPolicyGroup into different environments
	CrossID *string `json:"crossId,omitempty"`
	// SharedPolicyGroup name
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// SharedPolicyGroup description
	Description *string `json:"description,omitempty"`
	// SharedPolicyGroup prerequisite Message
	PrerequisiteMessage *string `json:"prerequisiteMessage,omitempty"`
	// Specify the SharedPolicyGroup ApiType
	// +kubebuilder:validation:Required
	ApiType ApiType `json:"apiType"`
	// SharedPolicyGroup phase (REQUEST;RESPONSE;INTERACT;CONNECT;PUBLISH;SUBSCRIBE)
	// +kubebuilder:validation:Required
	Phase *FlowPhase `json:"phase"`
	// SharedPolicyGroup Steps
	Steps []*Step `json:"steps,omitempty"`
	// SharedPolicyGroup LifecycleState (UNDEPLOYED;DEPLOYED;PENDING)
}

type Step struct {
	// +kubebuilder:default:=true
	// Indicate if this FlowStep is enabled or not
	Enabled bool `json:"enabled"`
	// FlowStep policy
	// +kubebuilder:validation:Optional
	Policy *string `json:"policy,omitempty"`
	// FlowStep name
	// +kubebuilder:validation:Optional
	Name *string `json:"name,omitempty"`
	// FlowStep description
	// +kubebuilder:validation:Optional
	Description *string `json:"description,omitempty"`
	// FlowStep configuration is a map of arbitrary key-values
	// +kubebuilder:validation:Optional
	Configuration *utils.GenericStringMap `json:"configuration,omitempty"`
	// FlowStep condition
	// +kubebuilder:validation:Optional
	Condition *string `json:"condition,omitempty"`
}
