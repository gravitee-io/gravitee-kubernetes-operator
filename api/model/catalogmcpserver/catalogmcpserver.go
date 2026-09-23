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

// Transport is the MCP transport of an upstream server.
// +kubebuilder:validation:Enum=HTTP
type Transport string

const (
	// TransportHTTP is the streamable HTTP transport, the only one discovery supports today.
	TransportHTTP Transport = "HTTP"
)

// AuthType discriminates how the platform authenticates against the upstream server.
// +kubebuilder:validation:Enum=NONE;HEADER;OAUTH2
type AuthType string

const (
	AuthTypeNone   AuthType = "NONE"
	AuthTypeHeader AuthType = "HEADER"
	AuthTypeOAuth2 AuthType = "OAUTH2"
)

// HeaderAuth is a static header sent on every request to the upstream server.
type HeaderAuth struct {
	// Name of the header, e.g. Authorization.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Full header value: `Bearer <token>` for a bearer token, `Basic <base64>` for basic
	// credentials, or the raw key for an API-key header. Use templating to read it from a
	// Secret: [[ secret `my-secret/token` ]]. Never returned by the platform.
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}

// OAuth2Auth are OAuth 2.0 client credentials; the platform fetches a token from tokenUrl
// before reaching the upstream server.
type OAuth2Auth struct {
	// +kubebuilder:validation:Required
	ClientID string `json:"clientId"`
	// Client secret, a literal or a templated Secret value. Never returned by the platform.
	// +kubebuilder:validation:Required
	ClientSecret string `json:"clientSecret"`
	// Token endpoint of the authorization server.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format=uri
	TokenURL string `json:"tokenUrl"`
	// +kubebuilder:validation:Optional
	Scope *string `json:"scope,omitempty"`
}

// Auth is the authentication the platform uses against the upstream MCP server, discriminated
// by type with exactly one nested block named after it. New ways of authenticating (a reference
// to a vaulted credential, for instance) are added as one more type and one more block, so
// existing manifests are never reshaped.
// +kubebuilder:validation:XValidation:rule="self.type != 'HEADER' || has(self.header)",message="header must be set when type is HEADER"
// +kubebuilder:validation:XValidation:rule="self.type == 'HEADER' || !has(self.header)",message="header must not be set when type is not HEADER"
// +kubebuilder:validation:XValidation:rule="self.type != 'OAUTH2' || has(self.oauth2)",message="oauth2 must be set when type is OAUTH2"
// +kubebuilder:validation:XValidation:rule="self.type == 'OAUTH2' || !has(self.oauth2)",message="oauth2 must not be set when type is not OAUTH2"
type Auth struct {
	// +kubebuilder:validation:Required
	Type AuthType `json:"type"`
	// Required when type is HEADER.
	// +kubebuilder:validation:Optional
	Header *HeaderAuth `json:"header,omitempty"`
	// Required when type is OAUTH2.
	// +kubebuilder:validation:Optional
	OAuth2 *OAuth2Auth `json:"oauth2,omitempty"`
}

// Connection is where the upstream MCP server answers and how the platform authenticates to it.
type Connection struct {
	// URL of the upstream MCP server.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format=uri
	Endpoint string `json:"endpoint"`
	// MCP transport of the upstream server.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=HTTP
	Transport Transport `json:"transport,omitempty"`
	// Authentication used to discover the server's capabilities. Omitted means none.
	// +kubebuilder:validation:Optional
	Auth *Auth `json:"auth,omitempty"`
}

// Type defines the specification of a CatalogMcpServer resource: an upstream MCP server
// registered in the AI Catalog. Tools, prompts and resources are discovered by the platform
// and reported in the status only.
type Type struct {
	// Stable catalog identity of the server, the name authorization policies reference:
	// lowercase, dot-separated segments, first segment `mcp-server`. Validated, never
	// repaired, and immutable once the server exists.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^mcp-server\.[a-z0-9_-]+(\.[a-z0-9_-]+)*$`
	// +kubebuilder:validation:MaxLength=255
	EntityID string `json:"entityId"`
	// +kubebuilder:validation:Optional
	Description *string `json:"description,omitempty"`
	// +kubebuilder:validation:Required
	Connection Connection `json:"connection"`
}
