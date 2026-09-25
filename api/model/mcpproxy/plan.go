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

// PlanSecurityType discriminates how consumers authenticate against a plan.
// +kubebuilder:validation:Enum=KEY_LESS;API_KEY;OAUTH2
type PlanSecurityType string

const (
	PlanSecurityKeyLess PlanSecurityType = "KEY_LESS"
	PlanSecurityAPIKey  PlanSecurityType = "API_KEY"
	PlanSecurityOAuth2  PlanSecurityType = "OAUTH2"
)

// APIKeySource is where the gateway reads a consumer's API key.
// +kubebuilder:validation:Enum=HEADER;BEARER;QUERY_PARAMETER
type APIKeySource string

// PlanAPIKey configures an API key plan.
type PlanAPIKey struct {
	// +kubebuilder:validation:Required
	Source APIKeySource `json:"source"`
	// Custom header carrying the key. Omitted means the gateway default.
	// +kubebuilder:validation:Optional
	Header *string `json:"header,omitempty"`
	// Forward the key to the upstream server.
	// +kubebuilder:validation:Optional
	Propagate bool `json:"propagate,omitempty"`
}

// PlanOAuth2 configures an OAuth 2.0 plan.
type PlanOAuth2 struct {
	// Name of the identityProviders entry that validates tokens.
	// +kubebuilder:validation:Required
	Provider string `json:"provider"`
	// Scopes required on the token. Only a GRAVITEE_AM provider accepts them.
	// +kubebuilder:validation:Optional
	Scopes []string `json:"scopes,omitempty"`
}

// PlanSecurity is the consumer authentication of a plan, discriminated by type with at most
// one nested block named after it.
// +kubebuilder:validation:XValidation:rule="self.type != 'API_KEY' || has(self.apiKey)",message="apiKey must be set when type is API_KEY"
// +kubebuilder:validation:XValidation:rule="self.type == 'API_KEY' || !has(self.apiKey)",message="apiKey must not be set when type is not API_KEY"
// +kubebuilder:validation:XValidation:rule="self.type != 'OAUTH2' || has(self.oauth2)",message="oauth2 must be set when type is OAUTH2"
// +kubebuilder:validation:XValidation:rule="self.type == 'OAUTH2' || !has(self.oauth2)",message="oauth2 must not be set when type is not OAUTH2"
type PlanSecurity struct {
	// +kubebuilder:validation:Required
	Type PlanSecurityType `json:"type"`
	// Required when type is API_KEY.
	// +kubebuilder:validation:Optional
	APIKey *PlanAPIKey `json:"apiKey,omitempty"`
	// Required when type is OAUTH2.
	// +kubebuilder:validation:Optional
	OAuth2 *PlanOAuth2 `json:"oauth2,omitempty"`
}

// Plan is a consumer plan of the proxy. Plans converge by name: a declared plan the proxy does
// not have is created and published, an open plan absent from the declaration is closed, which
// ends its subscriptions, and changing the security of an existing plan is refused.
type Plan struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	Security PlanSecurity `json:"security"`
	// +kubebuilder:validation:Optional
	Flows []Flow `json:"flows,omitempty"`
}

// IdentityProviderType discriminates the authorization server an OAUTH2 plan relies on.
// +kubebuilder:validation:Enum=GRAVITEE_AM;OAUTH2_GENERIC;OAUTH2_AUTH0
type IdentityProviderType string

const (
	IdentityProviderGraviteeAM    IdentityProviderType = "GRAVITEE_AM"
	IdentityProviderOAuth2Generic IdentityProviderType = "OAUTH2_GENERIC"
	IdentityProviderOAuth2Auth0   IdentityProviderType = "OAUTH2_AUTH0"
)

// IntrospectionMethod is the HTTP method of a token introspection call.
// +kubebuilder:validation:Enum=GET;POST
type IntrospectionMethod string

// OAuth2GenericProvider is any OAuth 2.0 server that introspects tokens.
type OAuth2GenericProvider struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format=uri
	IssuerURL string `json:"issuerUrl"`
	// +kubebuilder:validation:Required
	IntrospectionEndpoint string `json:"introspectionEndpoint"`
	// +kubebuilder:validation:Required
	IntrospectionEndpointMethod IntrospectionMethod `json:"introspectionEndpointMethod"`
	// +kubebuilder:validation:Optional
	ClientID *string `json:"clientId,omitempty"`
	// +kubebuilder:validation:Optional
	UserInfoEndpoint *string `json:"userInfoEndpoint,omitempty"`
	// A literal, a secret:// URI resolved by the gateway, or a templated Secret value.
	// Never returned by the platform.
	// +kubebuilder:validation:Optional
	ClientSecret *string `json:"clientSecret,omitempty"`
}

// Auth0Provider is an Auth0 tenant.
type Auth0Provider struct {
	// +kubebuilder:validation:Required
	Domain string `json:"domain"`
	// +kubebuilder:validation:Required
	Audience string `json:"audience"`
}

// IdentityProvider is an authorization server an OAUTH2 plan names, discriminated by type with
// at most one nested block. GRAVITEE_AM needs none and is accepted only when the proxy is created.
// +kubebuilder:validation:XValidation:rule="self.type != 'OAUTH2_GENERIC' || has(self.oauth2Generic)",message="oauth2Generic must be set when type is OAUTH2_GENERIC"
// +kubebuilder:validation:XValidation:rule="self.type == 'OAUTH2_GENERIC' || !has(self.oauth2Generic)",message="oauth2Generic must not be set when type is not OAUTH2_GENERIC"
// +kubebuilder:validation:XValidation:rule="self.type != 'OAUTH2_AUTH0' || has(self.auth0)",message="auth0 must be set when type is OAUTH2_AUTH0"
// +kubebuilder:validation:XValidation:rule="self.type == 'OAUTH2_AUTH0' || !has(self.auth0)",message="auth0 must not be set when type is not OAUTH2_AUTH0"
type IdentityProvider struct {
	// Name an OAUTH2 plan references in security.oauth2.provider.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	Type IdentityProviderType `json:"type"`
	// Required when type is OAUTH2_GENERIC.
	// +kubebuilder:validation:Optional
	OAuth2Generic *OAuth2GenericProvider `json:"oauth2Generic,omitempty"`
	// Required when type is OAUTH2_AUTH0.
	// +kubebuilder:validation:Optional
	Auth0 *Auth0Provider `json:"auth0,omitempty"`
}
