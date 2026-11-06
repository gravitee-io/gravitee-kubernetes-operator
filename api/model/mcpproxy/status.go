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

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StudioToolStatus is a tool a Studio exposes, with its catalog identity and the reference a
// policy names it by.
type StudioToolStatus struct {
	// HRID of the catalog server the tool comes from.
	// +kubebuilder:validation:Optional
	Server string `json:"server,omitempty"`
	// Runtime name of the tool on its server.
	// +kubebuilder:validation:Optional
	Tool string `json:"tool,omitempty"`
	// Name the Studio exposes the tool under, when one is declared.
	// +kubebuilder:validation:Optional
	Alias *string `json:"alias,omitempty"`
	// Catalog identity of the tool, e.g. mcp-tool.github.create_issue.
	// +kubebuilder:validation:Optional
	EntityID string `json:"entityId,omitempty"`
	// What an AuthzPolicy names the tool by, e.g. MCPTool::"github.create_issue": the resource the
	// gateway checks on tools/call when enableFGA is set. An aliased tool keeps the reference of the
	// tool it exposes.
	// +kubebuilder:validation:Optional
	AuthPolicyRef string `json:"authPolicyRef,omitempty"`
}

type Status struct {
	// The ID of the MCP proxy in the Gravitee API Management instance, also the id of the API
	// behind it
	// +kubebuilder:validation:Optional
	ID string `json:"id,omitempty"`
	// The organization ID defined in the management context
	// +kubebuilder:validation:Optional
	OrgID string `json:"organizationId,omitempty"`
	// The environment ID defined in the management context
	// +kubebuilder:validation:Optional
	EnvID string `json:"environmentId,omitempty"`
	// The human-readable ID the platform addresses this proxy by, derived from the
	// resource namespace and name.
	// +kubebuilder:validation:Optional
	HRID string `json:"hrid,omitempty"`
	// The lifecycle state the platform observes.
	// +kubebuilder:validation:Optional
	State string `json:"state,omitempty"`
	// The tools a Studio exposes, in the order of the platform's state, with their catalog
	// identities and the references policies name them by.
	// +kubebuilder:validation:Optional
	Tools []StudioToolStatus `json:"tools,omitempty"`
	// Conditions describe the current conditions of the McpProxy.
	//
	// Known condition types are:
	// * "Accepted"
	// * "ResolvedRefs"
	// * "AutomationAPIManaged"
	//
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:default={}
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// When the proxy has been created regardless of errors, this field is
	// used to persist the error message encountered during admission
	Errors status.Errors `json:"errors,omitempty"`
}
