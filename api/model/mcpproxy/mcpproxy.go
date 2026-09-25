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
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
)

// Mode is the shape of an MCP proxy: one upstream fronted whole, or a Studio composing tools
// selected across catalog servers.
// +kubebuilder:validation:Enum=PROXY;STUDIO
type Mode string

const (
	ModeProxy  Mode = "PROXY"
	ModeStudio Mode = "STUDIO"
)

// LifecycleState is the lifecycle an MCP proxy is declared or observed in.
// +kubebuilder:validation:Enum=STARTED;STOPPED
type LifecycleState string

const (
	StateStarted LifecycleState = "STARTED"
	StateStopped LifecycleState = "STOPPED"
)

// Proxy is the PROXY mode: a single upstream fronted whole.
type Proxy struct {
	// Raw URL of the upstream MCP server. The platform does not contact it on write; the gateway
	// reaches it at runtime.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format=uri
	ServerURL string `json:"serverUrl"`
	// Credential the gateway presents to the upstream. Omitted or NONE passes the caller's
	// credentials through.
	// +kubebuilder:validation:Optional
	UpstreamAuth *UpstreamAuth `json:"upstreamAuth,omitempty"`
}

// StudioTool selects one tool a catalog server exposes.
type StudioTool struct {
	// The CatalogMcpServer exposing the tool. The namespace defaults to the proxy's.
	// +kubebuilder:validation:Required
	ServerRef refs.NamespacedName `json:"serverRef"`
	// Name of the tool as the server advertises it (see the server's status.tools).
	// +kubebuilder:validation:Required
	Tool string `json:"tool"`
	// Name the Studio exposes the tool under.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^[a-zA-Z0-9_=-]+$`
	Alias *string `json:"alias,omitempty"`
}

// StudioUpstreamAuth is the credential the gateway presents to one catalog server.
type StudioUpstreamAuth struct {
	// +kubebuilder:validation:Required
	ServerRef refs.NamespacedName `json:"serverRef"`
	// Use type NONE when the server needs no credential.
	// +kubebuilder:validation:Required
	Auth UpstreamAuth `json:"auth"`
}

// Studio is the STUDIO mode: a curated tool surface composed from catalog servers.
type Studio struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Tools []StudioTool `json:"tools"`
	// One entry per server a selected tool comes from.
	// +kubebuilder:validation:Optional
	UpstreamAuth []StudioUpstreamAuth `json:"upstreamAuth,omitempty"`
	// Insert the tool-level authorization policy on tools/call. The platform owns that flow:
	// declared flows must not carry the authz-pep policy.
	// +kubebuilder:validation:Optional
	EnableFGA bool `json:"enableFGA,omitempty"`
}

// Type defines the specification of an McpProxy resource.
type Type struct {
	// Stable identity of the proxy, the name authorization policies reference: lowercase,
	// dot-separated segments, first segment `mcp-proxy`. Immutable.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^mcp-proxy\.[a-z0-9_-]+(\.[a-z0-9_-]+)*$`
	// +kubebuilder:validation:MaxLength=255
	EntityID string `json:"entityId"`
	// Display name. Immutable.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Immutable.
	// +kubebuilder:validation:Optional
	Description *string `json:"description,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^/`
	ContextPath string `json:"contextPath"`
	// MCP protocol version the proxy speaks. Immutable.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum="2024-11-05";"2025-03-26";"2025-11-25"
	ProtocolVersion string `json:"protocolVersion"`
	// PROXY requires the proxy block, STUDIO the studio block. Immutable.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=PROXY
	Mode Mode `json:"mode,omitempty"`
	// Required when mode is PROXY.
	// +kubebuilder:validation:Optional
	Proxy *Proxy `json:"proxy,omitempty"`
	// Required when mode is STUDIO.
	// +kubebuilder:validation:Optional
	Studio *Studio `json:"studio,omitempty"`
	// STARTED starts the proxy and redeploys it when it changed; STOPPED stops it and keeps later
	// changes stored until it is started again.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=STARTED
	State LifecycleState `json:"state,omitempty"`
	// +kubebuilder:validation:Optional
	FlowExecution *v4.FlowExecution `json:"flowExecution,omitempty"`
	// Proxy-level flows. The manifest owns them: omitted means none.
	// +kubebuilder:validation:Optional
	Flows []Flow `json:"flows,omitempty"`
	// +kubebuilder:validation:Optional
	// +listType=map
	// +listMapKey=name
	IdentityProviders []IdentityProvider `json:"identityProviders,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=name
	Plans []Plan `json:"plans"`
}

// ServerRefs returns every CatalogMcpServer the studio references, through a tool or an
// upstreamAuth entry, once each, in declaration order, with the namespace defaulted to ns.
func (t *Type) ServerRefs(ns string) []refs.NamespacedName {
	if t.Studio == nil {
		return nil
	}
	seen := make(map[string]bool)
	servers := make([]refs.NamespacedName, 0)
	add := func(ref refs.NamespacedName) {
		ref = WithNamespace(ref, ns)
		if !seen[ref.String()] {
			seen[ref.String()] = true
			servers = append(servers, ref)
		}
	}
	for i := range t.Studio.Tools {
		add(t.Studio.Tools[i].ServerRef)
	}
	for i := range t.Studio.UpstreamAuth {
		add(t.Studio.UpstreamAuth[i].ServerRef)
	}
	return servers
}

// ServersWithoutUpstreamAuth returns every server a selected tool comes from that has no
// studio.upstreamAuth entry, once each, with the namespace defaulted to ns.
func (t *Type) ServersWithoutUpstreamAuth(ns string) []refs.NamespacedName {
	if t.Studio == nil {
		return nil
	}
	declared := make(map[string]bool)
	for i := range t.Studio.UpstreamAuth {
		ref := WithNamespace(t.Studio.UpstreamAuth[i].ServerRef, ns)
		declared[ref.String()] = true
	}
	missing := make([]refs.NamespacedName, 0)
	for i := range t.Studio.Tools {
		ref := WithNamespace(t.Studio.Tools[i].ServerRef, ns)
		if !declared[ref.String()] {
			declared[ref.String()] = true
			missing = append(missing, ref)
		}
	}
	return missing
}

// WithNamespace returns ref with its namespace defaulted to ns and its kind dropped.
func WithNamespace(ref refs.NamespacedName, ns string) refs.NamespacedName {
	if ref.Namespace == "" {
		ref.Namespace = ns
	}
	ref.Kind = ""
	return ref
}
