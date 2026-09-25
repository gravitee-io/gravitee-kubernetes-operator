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
package mcpproxy

import "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"

// The flow types below are the subset of the v4 flow model an MCP proxy accepts: MCP and
// CONDITION selectors, request and response phases. The platform refuses any other field, so
// reusing the ApiV4Definition flow type would admit manifests that can never be applied.

// FlowSelectorType discriminates a flow selector.
// +kubebuilder:validation:Enum=MCP;CONDITION
type FlowSelectorType string

const (
	FlowSelectorMCP       FlowSelectorType = "MCP"
	FlowSelectorCondition FlowSelectorType = "CONDITION"
)

// McpSelector scopes a flow to MCP JSON-RPC methods, e.g. tools/call.
type McpSelector struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Methods []string `json:"methods"`
}

// ConditionSelector scopes a flow to requests matching an EL condition.
type ConditionSelector struct {
	// +kubebuilder:validation:Required
	Condition string `json:"condition"`
}

// FlowSelector is discriminated by type with exactly one nested block named after it.
// +kubebuilder:validation:XValidation:rule="self.type != 'MCP' || has(self.mcp)",message="mcp must be set when type is MCP"
// +kubebuilder:validation:XValidation:rule="self.type == 'MCP' || !has(self.mcp)",message="mcp must not be set when type is not MCP"
// +kubebuilder:validation:XValidation:rule="self.type != 'CONDITION' || has(self.condition)",message="condition must be set when type is CONDITION"
// +kubebuilder:validation:XValidation:rule="self.type == 'CONDITION' || !has(self.condition)",message="condition must not be set when type is not CONDITION"
type FlowSelector struct {
	// +kubebuilder:validation:Required
	Type FlowSelectorType `json:"type"`
	// Required when type is MCP.
	// +kubebuilder:validation:Optional
	MCP *McpSelector `json:"mcp,omitempty"`
	// Required when type is CONDITION.
	// +kubebuilder:validation:Optional
	Condition *ConditionSelector `json:"condition,omitempty"`
}

// FlowStep is one policy execution in a flow phase.
type FlowStep struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Policy plugin id, e.g. rate-limit.
	// +kubebuilder:validation:Required
	Policy string `json:"policy"`
	// +kubebuilder:validation:Optional
	Description *string `json:"description,omitempty"`
	// +kubebuilder:default:=true
	Enabled bool `json:"enabled"`
	// +kubebuilder:validation:Optional
	Condition *string `json:"condition,omitempty"`
	// Policy configuration, free-form per plugin.
	// +kubebuilder:validation:Optional
	Configuration *utils.GenericStringMap `json:"configuration,omitempty"`
}

// Flow is a set of policies applied to the requests and responses its selectors match.
type Flow struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:default:=true
	Enabled bool `json:"enabled"`
	// +kubebuilder:validation:Optional
	Selectors []FlowSelector `json:"selectors,omitempty"`
	// +kubebuilder:validation:Optional
	Request []FlowStep `json:"request,omitempty"`
	// +kubebuilder:validation:Optional
	Response []FlowStep `json:"response,omitempty"`
	// +kubebuilder:validation:Optional
	Tags []string `json:"tags,omitempty"`
}
