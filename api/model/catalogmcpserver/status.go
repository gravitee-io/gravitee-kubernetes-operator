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
package catalogmcpserver

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServerInfo is what the upstream server says about itself.
type ServerInfo struct {
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Optional
	Version string `json:"version,omitempty"`
}

// Tool is a tool the upstream server exposes, with the catalog identity a policy names.
type Tool struct {
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Optional
	EntityID string `json:"entityId,omitempty"`
}

// Prompt is a prompt the upstream server exposes, with the catalog identity a policy names.
type Prompt struct {
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Optional
	EntityID string `json:"entityId,omitempty"`
}

// Resource is a resource the upstream server exposes, with the catalog identity a policy names.
type Resource struct {
	// +kubebuilder:validation:Optional
	URI string `json:"uri,omitempty"`
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Optional
	EntityID string `json:"entityId,omitempty"`
}

// Discovered is what the platform learned from the upstream server at the last sync.
// It is read-only and never compared for drift.
type Discovered struct {
	// When the platform last discovered the server, as an RFC 3339 date-time.
	// +kubebuilder:validation:Optional
	LastSyncedAt string `json:"lastSyncedAt,omitempty"`
	// MCP protocol version negotiated with the server.
	// +kubebuilder:validation:Optional
	ProtocolVersion string `json:"protocolVersion,omitempty"`
	// +kubebuilder:validation:Optional
	ServerInfo *ServerInfo `json:"serverInfo,omitempty"`
	// +kubebuilder:validation:Optional
	Tools []Tool `json:"tools,omitempty"`
	// +kubebuilder:validation:Optional
	Prompts []Prompt `json:"prompts,omitempty"`
	// +kubebuilder:validation:Optional
	Resources []Resource `json:"resources,omitempty"`
}

type Status struct {
	// The ID of the catalog server in the Gravitee API Management instance
	// +kubebuilder:validation:Optional
	ID string `json:"id,omitempty"`
	// The organization ID defined in the management context
	// +kubebuilder:validation:Optional
	OrgID string `json:"organizationId,omitempty"`
	// The environment ID defined in the management context
	// +kubebuilder:validation:Optional
	EnvID string `json:"environmentId,omitempty"`
	// The human-readable ID the platform addresses this server by, derived from the
	// resource namespace and name.
	// +kubebuilder:validation:Optional
	HRID string `json:"hrid,omitempty"`
	// Conditions describe the current conditions of the CatalogMcpServer.
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
	// When the server has been created regardless of errors, this field is
	// used to persist the error message encountered during admission
	Errors status.Errors `json:"errors,omitempty"`
	// What the platform discovered from the upstream server.
	Discovered `json:",inline"`
}
