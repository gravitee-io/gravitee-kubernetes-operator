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

package domain

// Domain is the spec of a security domain managed by the Automation API.
// Certificates, identity providers, and reporters are not embedded;
// they are managed via the domain's sub-resource endpoints.
type Domain struct {
	// AccountSettings holds the user account settings for the domain: brute-force protection,
	// registration, password reset, remember-me, and MFA challenge behavior.
	// +kubebuilder:validation:Optional
	AccountSettings *AccountSettings `json:"accountSettings,omitempty"`

	// AlertEnabled controls whether alerting is enabled for the domain.
	// +kubebuilder:validation:Optional
	AlertEnabled *bool `json:"alertEnabled,omitempty"`

	// CertificateSettings holds the domain-level certificate settings.
	// +kubebuilder:validation:Optional
	CertificateSettings *CertificateSettings `json:"certificateSettings,omitempty"`

	// CorsSettings holds the Cross-Origin Resource Sharing configuration controlling which web origins
	// may call the domain's endpoints from a browser.
	// +kubebuilder:validation:Optional
	CorsSettings *CorsSettings `json:"corsSettings,omitempty"`

	// DataPlaneId is the identifier of the data plane this domain is connected to. Optional at
	// creation and resolved from the environment's data planes when omitted. Immutable afterwards: an
	// apply that names a different one is rejected.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxLength=255
	DataPlaneId *string `json:"dataPlaneId,omitempty"`

	// Description is the human-readable description of the domain.
	// +kubebuilder:validation:Optional
	Description *string `json:"description,omitempty"`

	// Enabled controls whether the domain handles incoming authentication and authorization requests.
	// Defaults to true.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	Enabled *bool `json:"enabled,omitempty"`

	// KeyRetrievalSettings holds the fetch, SSRF and cache limits applied to every trusted domain in
	// the security domain.
	// +kubebuilder:validation:Optional
	KeyRetrievalSettings *KeyRetrievalSettings `json:"keyRetrievalSettings,omitempty"`

	// LoginSettings holds the configuration of the domain's login flow and the features offered on the
	// sign-in page.
	// +kubebuilder:validation:Optional
	LoginSettings *LoginSettings `json:"loginSettings,omitempty"`

	// Master controls whether this is the master domain of its environment. A master domain may
	// perform cross-domain token introspection. Defaults to false.
	// +kubebuilder:validation:Optional
	Master *bool `json:"master,omitempty"`

	// Name is the human-readable name of the domain.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	Name string `json:"name"`

	// Oidc holds the OpenID Connect settings for the domain. CIMD (client identity metadata document)
	// settings are not exposed by the Automation API and are reset on update.
	// +kubebuilder:validation:Optional
	Oidc *OidcSettings `json:"oidc,omitempty"`

	// PasswordSettings holds the password policy applied to users of the domain: complexity
	// requirements, expiry, and history.
	// +kubebuilder:validation:Optional
	PasswordSettings *PasswordSettings `json:"passwordSettings,omitempty"`

	// Path is the context path the domain is served under, relative to the gateway. Must start with a
	// slash.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^/.*`
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	Path string `json:"path"`

	// Saml holds the settings for the domain acting as a SAML 2.0 identity provider (IdP).
	// +kubebuilder:validation:Optional
	Saml *SamlSettings `json:"saml,omitempty"`

	// Scim holds the configuration of the domain's SCIM 2.0 provisioning endpoints.
	// +kubebuilder:validation:Optional
	Scim *SCIMSettings `json:"scim,omitempty"`

	// SecretExpirationSettings controls whether client secrets in the domain expire and after how
	// long.
	// +kubebuilder:validation:Optional
	SecretExpirationSettings *SecretExpirationSettings `json:"secretExpirationSettings,omitempty"`

	// SelfServiceAccountManagementSettings controls whether end users can manage their own account
	// (for example, reset their password) and the rules that apply.
	// +kubebuilder:validation:Optional
	SelfServiceAccountManagementSettings *SelfServiceAccountManagementSettings `json:"selfServiceAccountManagementSettings,omitempty"`

	// Tags lists the sharding tags that control which gateways deploy this domain.
	// +kubebuilder:validation:Optional
	// +listType=set
	Tags []string `json:"tags,omitempty"`

	// TokenExchangeSettings holds the OAuth 2.0 Token Exchange (RFC 8693) configuration for the
	// domain, covering impersonation and delegation.
	// +kubebuilder:validation:Optional
	TokenExchangeSettings *TokenExchangeSettings `json:"tokenExchangeSettings,omitempty"`

	// Uma holds the configuration of the domain's User-Managed Access (UMA 2.0) authorization
	// features.
	// +kubebuilder:validation:Optional
	Uma *UMASettings `json:"uma,omitempty"`

	// VhostMode controls whether the domain is exposed through its virtual hosts rather than the
	// default context path. When true, vhosts must be supplied. Defaults to false.
	// +kubebuilder:validation:Optional
	VhostMode *bool `json:"vhostMode,omitempty"`

	// Vhosts lists the virtual hosts the domain is exposed on, overriding the default context path.
	// +kubebuilder:validation:Optional
	Vhosts []VirtualHost `json:"vhosts,omitempty"`

	// WebAuthnSettings holds the WebAuthn (FIDO2) relying-party configuration governing passwordless
	// and multi-factor authentication for the domain.
	// +kubebuilder:validation:Optional
	WebAuthnSettings *WebAuthnSettings `json:"webAuthnSettings,omitempty"`

	// WebProtectionSettings holds the HTTP security headers applied to the domain's login and consent
	// pages.
	// +kubebuilder:validation:Optional
	WebProtectionSettings *WebProtectionSettings `json:"webProtectionSettings,omitempty"`
}
