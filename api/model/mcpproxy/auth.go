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

// UpstreamAuthType discriminates how the gateway authenticates against an upstream MCP server.
// +kubebuilder:validation:Enum=NONE;API_KEY;BEARER;BASIC;OAUTH2
type UpstreamAuthType string

const (
	UpstreamAuthTypeNone   UpstreamAuthType = "NONE"
	UpstreamAuthTypeAPIKey UpstreamAuthType = "API_KEY"
	UpstreamAuthTypeBearer UpstreamAuthType = "BEARER"
	UpstreamAuthTypeBasic  UpstreamAuthType = "BASIC"
	UpstreamAuthTypeOAuth2 UpstreamAuthType = "OAUTH2"
)

// APIKeyAuth is a static API key sent in a header.
type APIKeyAuth struct {
	// Name of the header carrying the key.
	// +kubebuilder:validation:Required
	Header string `json:"header"`
	// The key: a literal, a secret:// URI resolved by the gateway, or a templated Secret value
	// ([[ secret `my-secret/key` ]]). Never returned by the platform.
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}

// BearerAuth is a bearer token sent in the Authorization header.
type BearerAuth struct {
	// The token: a literal, a secret:// URI resolved by the gateway, or a templated Secret value.
	// Never returned by the platform.
	// +kubebuilder:validation:Required
	Token string `json:"token"`
}

// BasicAuth are HTTP basic credentials.
type BasicAuth struct {
	// +kubebuilder:validation:Required
	Username string `json:"username"`
	// A literal, a secret:// URI resolved by the gateway, or a templated Secret value.
	// Never returned by the platform.
	// +kubebuilder:validation:Required
	Password string `json:"password"`
}

// OAuth2Auth is the MCP authorization flow the gateway runs against the upstream server.
type OAuth2Auth struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format=uri
	AuthorizeURL string `json:"authorizeUrl"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format=uri
	TokenURL string `json:"tokenUrl"`
	// +kubebuilder:validation:Required
	ClientID string `json:"clientId"`
	// A literal, a secret:// URI resolved by the gateway, or a templated Secret value.
	// Never returned by the platform.
	// +kubebuilder:validation:Required
	ClientSecret string `json:"clientSecret"`
	// +kubebuilder:validation:Optional
	Scopes []string `json:"scopes,omitempty"`
}

// UpstreamAuth is the credential the gateway presents to an upstream MCP server, discriminated
// by type with exactly one nested block named after it. NONE passes the caller's credentials
// through. A reference to a vaulted credential is added later as one more type and one more
// block, so existing manifests are never reshaped.
// +kubebuilder:validation:XValidation:rule="self.type != 'API_KEY' || has(self.apiKey)",message="apiKey must be set when type is API_KEY"
// +kubebuilder:validation:XValidation:rule="self.type == 'API_KEY' || !has(self.apiKey)",message="apiKey must not be set when type is not API_KEY"
// +kubebuilder:validation:XValidation:rule="self.type != 'BEARER' || has(self.bearer)",message="bearer must be set when type is BEARER"
// +kubebuilder:validation:XValidation:rule="self.type == 'BEARER' || !has(self.bearer)",message="bearer must not be set when type is not BEARER"
// +kubebuilder:validation:XValidation:rule="self.type != 'BASIC' || has(self.basic)",message="basic must be set when type is BASIC"
// +kubebuilder:validation:XValidation:rule="self.type == 'BASIC' || !has(self.basic)",message="basic must not be set when type is not BASIC"
// +kubebuilder:validation:XValidation:rule="self.type != 'OAUTH2' || has(self.oauth2)",message="oauth2 must be set when type is OAUTH2"
// +kubebuilder:validation:XValidation:rule="self.type == 'OAUTH2' || !has(self.oauth2)",message="oauth2 must not be set when type is not OAUTH2"
type UpstreamAuth struct {
	// +kubebuilder:validation:Required
	Type UpstreamAuthType `json:"type"`
	// Required when type is API_KEY.
	// +kubebuilder:validation:Optional
	APIKey *APIKeyAuth `json:"apiKey,omitempty"`
	// Required when type is BEARER.
	// +kubebuilder:validation:Optional
	Bearer *BearerAuth `json:"bearer,omitempty"`
	// Required when type is BASIC.
	// +kubebuilder:validation:Optional
	Basic *BasicAuth `json:"basic,omitempty"`
	// Required when type is OAUTH2.
	// +kubebuilder:validation:Optional
	OAuth2 *OAuth2Auth `json:"oauth2,omitempty"`
}
