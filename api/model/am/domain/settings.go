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

// AccountSettings are user account settings for the domain:
// brute-force protection, registration, password reset, remember-me, and MFA challenge behavior.
type AccountSettings struct {
	// AccountBlockedDuration is the duration, in seconds, for which the account
	// remains blocked after too many failed login attempts.
	// +kubebuilder:validation:Optional
	AccountBlockedDuration *int32 `json:"accountBlockedDuration,omitempty"`

	// AutoLoginAfterRegistration controls whether the user is automatically
	// logged in after completing registration.
	// +kubebuilder:validation:Optional
	AutoLoginAfterRegistration *bool `json:"autoLoginAfterRegistration,omitempty"`

	// AutoLoginAfterResetPassword controls whether the user is automatically
	// logged in after a password reset.
	// +kubebuilder:validation:Optional
	AutoLoginAfterResetPassword *bool `json:"autoLoginAfterResetPassword,omitempty"`

	// CompleteRegistrationWhenResetPassword controls whether resetting a
	// password also completes a pending registration.
	// +kubebuilder:validation:Optional
	CompleteRegistrationWhenResetPassword *bool `json:"completeRegistrationWhenResetPassword,omitempty"`

	// DefaultIdentityProviderForRegistration is the key of an identity provider
	// that exists under this domain, used as the default for user registration.
	// +kubebuilder:validation:Optional
	DefaultIdentityProviderForRegistration *string `json:"defaultIdentityProviderForRegistration,omitempty"`

	// DeletePasswordlessDevicesAfterResetPassword controls whether passwordless
	// (WebAuthn) devices are deleted when the password is reset.
	// +kubebuilder:validation:Optional
	DeletePasswordlessDevicesAfterResetPassword *bool `json:"deletePasswordlessDevicesAfterResetPassword,omitempty"`

	// DynamicUserRegistration controls whether dynamic (self-service) user
	// registration is enabled.
	// +kubebuilder:validation:Optional
	DynamicUserRegistration *bool `json:"dynamicUserRegistration,omitempty"`

	// Inherited controls whether account settings are inherited from the parent.
	// When true, the other fields are ignored.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`

	// LoginAttemptsDetectionEnabled controls whether brute-force authentication
	// attempts are detected and blocked.
	// +kubebuilder:validation:Optional
	LoginAttemptsDetectionEnabled *bool `json:"loginAttemptsDetectionEnabled,omitempty"`

	// LoginAttemptsResetTime is the time, in seconds, after which the login
	// attempt counter is reset when the maximum has not been reached.
	// +kubebuilder:validation:Optional
	LoginAttemptsResetTime *int32 `json:"loginAttemptsResetTime,omitempty"`

	// MaxLoginAttempts is the maximum number of failed login attempts
	// before the account is blocked.
	// +kubebuilder:validation:Optional
	MaxLoginAttempts *int32 `json:"maxLoginAttempts,omitempty"`

	// MfaChallengeAttemptsDetectionEnabled controls whether failed MFA
	// challenge attempts are detected and blocked.
	// +kubebuilder:validation:Optional
	MfaChallengeAttemptsDetectionEnabled *bool `json:"mfaChallengeAttemptsDetectionEnabled,omitempty"`

	// MfaChallengeAttemptsResetTime is the time, in seconds, after which
	// the MFA challenge attempt counter is reset.
	// +kubebuilder:validation:Optional
	MfaChallengeAttemptsResetTime *int32 `json:"mfaChallengeAttemptsResetTime,omitempty"`

	// MfaChallengeMaxAttempts is the maximum number of failed MFA challenge
	// attempts before the user is blocked.
	// +kubebuilder:validation:Optional
	MfaChallengeMaxAttempts *int32 `json:"mfaChallengeMaxAttempts,omitempty"`

	// MfaChallengeSendVerifyAlertEmail controls whether to send an alert
	// email after too many failed MFA challenge attempts.
	// +kubebuilder:validation:Optional
	MfaChallengeSendVerifyAlertEmail *bool `json:"mfaChallengeSendVerifyAlertEmail,omitempty"`

	// RedirectUriAfterRegistration is the URL the user is redirected to after registration.
	// +kubebuilder:validation:Optional
	RedirectUriAfterRegistration *string `json:"redirectUriAfterRegistration,omitempty"`

	// RedirectUriAfterResetPassword is the URL the user is redirected to after a password reset.
	// +kubebuilder:validation:Optional
	RedirectUriAfterResetPassword *string `json:"redirectUriAfterResetPassword,omitempty"`

	// RememberMe controls whether users can remain logged in for a fixed
	// duration (remember-me).
	// +kubebuilder:validation:Optional
	RememberMe *bool `json:"rememberMe,omitempty"`

	// RememberMeDuration is the duration, in seconds, for which a remembered
	// session stays valid.
	// +kubebuilder:validation:Optional
	RememberMeDuration *int32 `json:"rememberMeDuration,omitempty"`

	// ResetPasswordConfirmIdentity controls whether the user must confirm
	// their identity before resetting a password.
	// +kubebuilder:validation:Optional
	ResetPasswordConfirmIdentity *bool `json:"resetPasswordConfirmIdentity,omitempty"`

	// ResetPasswordCustomForm controls whether a custom form is used for
	// the password-reset step.
	// +kubebuilder:validation:Optional
	ResetPasswordCustomForm *bool `json:"resetPasswordCustomForm,omitempty"`

	// ResetPasswordCustomFormFields are the custom fields rendered
	// on the password-reset form.
	// +kubebuilder:validation:Optional
	ResetPasswordCustomFormFields []FormField `json:"resetPasswordCustomFormFields,omitempty"`

	// ResetPasswordInvalidateTokens controls whether existing tokens are
	// invalidated when the password is reset.
	// +kubebuilder:validation:Optional
	ResetPasswordInvalidateTokens *bool `json:"resetPasswordInvalidateTokens,omitempty"`

	// SendRecoverAccountEmail controls whether to send an account-recovery email.
	// +kubebuilder:validation:Optional
	SendRecoverAccountEmail *bool `json:"sendRecoverAccountEmail,omitempty"`

	// SendVerifyRegistrationAccountEmail controls whether to send a
	// registration-verification email.
	// +kubebuilder:validation:Optional
	SendVerifyRegistrationAccountEmail *bool `json:"sendVerifyRegistrationAccountEmail,omitempty"`
}

// CertificateSettings are domain-level certificate settings.
type CertificateSettings struct {
	// FallbackCertificate is the key of a certificate managed under this domain,
	// used as the fallback when a client does not specify one.
	// +kubebuilder:validation:Optional
	FallbackCertificate *string `json:"fallbackCertificate,omitempty"`
}

// CorsSettings is the Cross-Origin Resource Sharing configuration for the domain.
type CorsSettings struct {
	// AllowCredentials controls whether credentials (cookies, authorization headers,
	// TLS client certificates) are included in CORS responses.
	// +kubebuilder:validation:Optional
	AllowCredentials *bool `json:"allowCredentials,omitempty"`

	// AllowedHeaders lists the HTTP headers permitted on CORS requests.
	// +kubebuilder:validation:Optional
	AllowedHeaders []string `json:"allowedHeaders,omitempty"`

	// AllowedMethods lists the HTTP methods permitted on CORS requests.
	// +kubebuilder:validation:Optional
	AllowedMethods []string `json:"allowedMethods,omitempty"`

	// AllowedOrigins lists the origins permitted to make CORS requests.
	// +kubebuilder:validation:Optional
	AllowedOrigins []string `json:"allowedOrigins,omitempty"`

	// Enabled controls whether CORS is enabled for the domain.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// Inherited controls whether CORS settings are inherited from the parent.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`

	// MaxAge is the maximum time, in seconds, a preflight response may be cached.
	// +kubebuilder:validation:Optional
	MaxAge *int32 `json:"maxAge,omitempty"`
}

// FormField is a custom field rendered on a domain form.
type FormField struct {
	// Key is the field identifier.
	// +kubebuilder:validation:Optional
	Key *string `json:"key,omitempty"`

	// Label is the human-readable label for the field.
	// +kubebuilder:validation:Optional
	Label *string `json:"label,omitempty"`

	// Type is the field type (text, select, etc.).
	// +kubebuilder:validation:Optional
	Type *string `json:"type,omitempty"`
}

// KeyRetrievalSettings are the fetch, SSRF and cache limits
// applied to every trusted domain in the security domain.
type KeyRetrievalSettings struct {
	// AllowPrivateIpAddress controls whether private IP addresses are
	// allowed in key retrieval URIs.
	// +kubebuilder:validation:Optional
	AllowPrivateIpAddress *bool `json:"allowPrivateIpAddress,omitempty"`

	// AllowUnsecuredHttpUri controls whether unsecured HTTP URIs are
	// allowed for key retrieval.
	// +kubebuilder:validation:Optional
	AllowUnsecuredHttpUri *bool `json:"allowUnsecuredHttpUri,omitempty"`

	// CacheMaxEntries is the maximum number of entries in the key retrieval cache.
	// +kubebuilder:validation:Optional
	CacheMaxEntries *int32 `json:"cacheMaxEntries,omitempty"`

	// CacheTtlSeconds is the time-to-live, in seconds, for cached key retrieval entries.
	// +kubebuilder:validation:Optional
	CacheTtlSeconds *int32 `json:"cacheTtlSeconds,omitempty"`

	// FetchTimeoutMs is the timeout, in milliseconds, for key retrieval HTTP requests.
	// +kubebuilder:validation:Optional
	FetchTimeoutMs *int32 `json:"fetchTimeoutMs,omitempty"`

	// MaxResponseSizeKb is the maximum response size, in kilobytes, for key retrieval.
	// +kubebuilder:validation:Optional
	MaxResponseSizeKb *int32 `json:"maxResponseSizeKb,omitempty"`
}

// LoginSettings is the configuration of the domain's login flow.
type LoginSettings struct {
	// CertificateBasedAuthEnabled controls whether certificate-based
	// authentication is available on the login page.
	// +kubebuilder:validation:Optional
	CertificateBasedAuthEnabled *bool `json:"certificateBasedAuthEnabled,omitempty"`

	// CertificateBasedAuthUrl is the URL endpoint for certificate-based authentication.
	// +kubebuilder:validation:Optional
	CertificateBasedAuthUrl *string `json:"certificateBasedAuthUrl,omitempty"`

	// ForgotPasswordEnabled controls whether the forgot-password link is shown.
	// +kubebuilder:validation:Optional
	ForgotPasswordEnabled *bool `json:"forgotPasswordEnabled,omitempty"`

	// HideForm controls whether the login form is hidden
	// (useful when only external identity providers are used).
	// +kubebuilder:validation:Optional
	HideForm *bool `json:"hideForm,omitempty"`

	// IdentifierFirstEnabled controls whether identifier-first login flow is used.
	// +kubebuilder:validation:Optional
	IdentifierFirstEnabled *bool `json:"identifierFirstEnabled,omitempty"`

	// Inherited controls whether login settings are inherited from the parent.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`

	// MagicLinkAuthEnabled controls whether magic-link (passwordless email) login is available.
	// +kubebuilder:validation:Optional
	MagicLinkAuthEnabled *bool `json:"magicLinkAuthEnabled,omitempty"`

	// PasswordlessDeviceNamingEnabled controls whether passwordless devices can be named.
	// +kubebuilder:validation:Optional
	PasswordlessDeviceNamingEnabled *bool `json:"passwordlessDeviceNamingEnabled,omitempty"`

	// PasswordlessEnabled controls whether passwordless login is available.
	// +kubebuilder:validation:Optional
	PasswordlessEnabled *bool `json:"passwordlessEnabled,omitempty"`

	// PasswordlessEnforcePasswordEnabled controls whether a password is still
	// required alongside passwordless authentication.
	// +kubebuilder:validation:Optional
	PasswordlessEnforcePasswordEnabled *bool `json:"passwordlessEnforcePasswordEnabled,omitempty"`

	// PasswordlessEnforcePasswordMaxAge is the maximum age, in seconds,
	// before a password confirmation is re-required for passwordless users.
	// +kubebuilder:validation:Optional
	PasswordlessEnforcePasswordMaxAge *int32 `json:"passwordlessEnforcePasswordMaxAge,omitempty"`

	// PasswordlessRememberDeviceEnabled controls whether passwordless
	// devices are remembered across sessions.
	// +kubebuilder:validation:Optional
	PasswordlessRememberDeviceEnabled *bool `json:"passwordlessRememberDeviceEnabled,omitempty"`

	// RegisterEnabled controls whether the registration link is shown on the login page.
	// +kubebuilder:validation:Optional
	RegisterEnabled *bool `json:"registerEnabled,omitempty"`

	// RememberMeEnabled controls whether the remember-me option is shown on the login page.
	// +kubebuilder:validation:Optional
	RememberMeEnabled *bool `json:"rememberMeEnabled,omitempty"`

	// ResetPasswordOnExpiration controls whether a password reset is forced
	// when the password has expired.
	// +kubebuilder:validation:Optional
	ResetPasswordOnExpiration *bool `json:"resetPasswordOnExpiration,omitempty"`
}

// OidcSettings holds OpenID Connect settings for the domain.
type OidcSettings struct {
	// CibaSettings is the Client-Initiated Backchannel Authentication configuration.
	// +kubebuilder:validation:Optional
	CibaSettings *CIBASettings `json:"cibaSettings,omitempty"`

	// ClientRegistrationSettings is the Dynamic Client Registration configuration.
	// +kubebuilder:validation:Optional
	ClientRegistrationSettings *ClientRegistrationSettings `json:"clientRegistrationSettings,omitempty"`

	// PostLogoutRedirectUris is the list of URIs allowed for post-logout redirection.
	// +kubebuilder:validation:Optional
	PostLogoutRedirectUris []string `json:"postLogoutRedirectUris,omitempty"`

	// RedirectUriStrictMatching controls whether redirect URI matching is strict
	// (no wildcard or partial matching).
	// +kubebuilder:validation:Optional
	RedirectUriStrictMatching *bool `json:"redirectUriStrictMatching,omitempty"`

	// RequestUris is the list of pre-registered request URIs.
	// +kubebuilder:validation:Optional
	RequestUris []string `json:"requestUris,omitempty"`

	// SecurityProfileSettings holds FAPI security profile settings.
	// +kubebuilder:validation:Optional
	SecurityProfileSettings *SecurityProfileSettings `json:"securityProfileSettings,omitempty"`

	// WorkloadIdentitySettings are the workload identity (SPIFFE) settings.
	// +kubebuilder:validation:Optional
	WorkloadIdentitySettings *SpiffeDomainSettings `json:"workloadIdentitySettings,omitempty"`
}

// SpiffeDomainSettings are the workload identity (SPIFFE) settings for the domain.
// Key retrieval limits are configured in keyRetrievalSettings.
type SpiffeDomainSettings struct {
	// ClockSkewSeconds is the allowed clock skew, in seconds, when validating JWT temporal claims.
	// +kubebuilder:validation:Optional
	ClockSkewSeconds *int32 `json:"clockSkewSeconds,omitempty"`

	// DefaultAllowedAlgorithms is the default allowlist of signature algorithms
	// accepted for SPIFFE JWT validation.
	// +kubebuilder:validation:Optional
	DefaultAllowedAlgorithms []string `json:"defaultAllowedAlgorithms,omitempty"`

	// Enabled controls whether SPIFFE workload identity support is enabled.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// MaxJwtLifetimeSeconds is the maximum accepted JWT lifetime, in seconds, computed as exp minus iat.
	// +kubebuilder:validation:Optional
	MaxJwtLifetimeSeconds *int32 `json:"maxJwtLifetimeSeconds,omitempty"`
}

// CIBASettings is the Client-Initiated Backchannel Authentication configuration.
type CIBASettings struct {
	// AuthReqExpiry is the default validity period, in seconds, of the issued auth_req_id.
	// +kubebuilder:validation:Optional
	AuthReqExpiry *int32 `json:"authReqExpiry,omitempty"`

	// BindingMessageLength is the maximum number of characters accepted
	// for the binding_message parameter.
	// +kubebuilder:validation:Optional
	BindingMessageLength *int32 `json:"bindingMessageLength,omitempty"`

	// Enabled controls whether CIBA is enabled for the domain.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// TokenReqInterval is the minimum delay, in seconds, between two polls
	// of the token endpoint for the same auth_req_id.
	// +kubebuilder:validation:Optional
	TokenReqInterval *int32 `json:"tokenReqInterval,omitempty"`
}

// ClientRegistrationSettings is the OpenID Connect Dynamic Client Registration configuration.
type ClientRegistrationSettings struct {
	// AllowHttpSchemeRedirectUri controls whether the unsecured http scheme
	// is permitted in redirect URIs.
	// +kubebuilder:validation:Optional
	AllowHttpSchemeRedirectUri *bool `json:"allowHttpSchemeRedirectUri,omitempty"`

	// AllowLocalhostRedirectUri controls whether localhost is permitted
	// as a redirect URI host.
	// +kubebuilder:validation:Optional
	AllowLocalhostRedirectUri *bool `json:"allowLocalhostRedirectUri,omitempty"`

	// AllowRedirectUriParamsExpressionLanguage controls whether expression
	// language is permitted in redirect URI parameters.
	// +kubebuilder:validation:Optional
	AllowRedirectUriParamsExpressionLanguage *bool `json:"allowRedirectUriParamsExpressionLanguage,omitempty"`

	// AllowWildCardRedirectUri controls whether wildcards are permitted in redirect URIs.
	// +kubebuilder:validation:Optional
	AllowWildCardRedirectUri *bool `json:"allowWildCardRedirectUri,omitempty"`

	// AllowedScopes lists scopes permitted on client registration requests.
	// +kubebuilder:validation:Optional
	AllowedScopes []string `json:"allowedScopes,omitempty"`

	// AllowedScopesEnabled controls whether registered client scopes
	// are restricted to the allowed list.
	// +kubebuilder:validation:Optional
	AllowedScopesEnabled *bool `json:"allowedScopesEnabled,omitempty"`

	// ClientTemplateEnabled controls whether a client template is used
	// for dynamic registration.
	// +kubebuilder:validation:Optional
	ClientTemplateEnabled *bool `json:"clientTemplateEnabled,omitempty"`

	// DefaultScopes are added to every client registration request.
	// +kubebuilder:validation:Optional
	DefaultScopes []string `json:"defaultScopes,omitempty"`

	// DynamicClientRegistrationEnabled controls whether dynamic client
	// registration is enabled.
	// +kubebuilder:validation:Optional
	DynamicClientRegistrationEnabled *bool `json:"dynamicClientRegistrationEnabled,omitempty"`

	// OpenDynamicClientRegistrationEnabled controls whether open (unauthenticated)
	// dynamic client registration is enabled.
	// +kubebuilder:validation:Optional
	OpenDynamicClientRegistrationEnabled *bool `json:"openDynamicClientRegistrationEnabled,omitempty"`
}

// SecurityProfileSettings holds FAPI security profile settings.
type SecurityProfileSettings struct {
	// EnableFapiBrazil controls whether the FAPI Brazil profile is enabled.
	// +kubebuilder:validation:Optional
	EnableFapiBrazil *bool `json:"enableFapiBrazil,omitempty"`

	// EnablePlainFapi controls whether the plain FAPI profile is enabled.
	// +kubebuilder:validation:Optional
	EnablePlainFapi *bool `json:"enablePlainFapi,omitempty"`
}

// PasswordSettings is the password policy applied to domain users.
type PasswordSettings struct {
	// ExcludePasswordsInDictionary controls whether dictionary words are
	// rejected as passwords.
	// +kubebuilder:validation:Optional
	ExcludePasswordsInDictionary *bool `json:"excludePasswordsInDictionary,omitempty"`

	// ExcludeUserProfileInfoInPassword controls whether user profile
	// information (name, email) is rejected in passwords.
	// +kubebuilder:validation:Optional
	ExcludeUserProfileInfoInPassword *bool `json:"excludeUserProfileInfoInPassword,omitempty"`

	// ExpiryDuration is the maximum lifetime, in seconds, of a password
	// before the user must change it.
	// +kubebuilder:validation:Optional
	ExpiryDuration *int32 `json:"expiryDuration,omitempty"`

	// IncludeNumbers controls whether passwords must include at least one digit.
	// +kubebuilder:validation:Optional
	IncludeNumbers *bool `json:"includeNumbers,omitempty"`

	// IncludeSpecialCharacters controls whether passwords must include
	// at least one special character.
	// +kubebuilder:validation:Optional
	IncludeSpecialCharacters *bool `json:"includeSpecialCharacters,omitempty"`

	// Inherited controls whether password settings are inherited from the parent.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`

	// LettersInMixedCase controls whether passwords must include both
	// uppercase and lowercase letters.
	// +kubebuilder:validation:Optional
	LettersInMixedCase *bool `json:"lettersInMixedCase,omitempty"`

	// MaxConsecutiveLetters is the maximum number of consecutive identical
	// characters allowed in a password.
	// +kubebuilder:validation:Optional
	MaxConsecutiveLetters *int32 `json:"maxConsecutiveLetters,omitempty"`

	// MaxLength is the maximum password length.
	// +kubebuilder:validation:Optional
	MaxLength *int32 `json:"maxLength,omitempty"`

	// MinLength is the minimum password length.
	// +kubebuilder:validation:Optional
	MinLength *int32 `json:"minLength,omitempty"`

	// OldPasswords is the number of previous passwords that cannot be reused.
	// +kubebuilder:validation:Optional
	OldPasswords *int32 `json:"oldPasswords,omitempty"`

	// PasswordHistoryEnabled controls whether password history is enforced.
	// +kubebuilder:validation:Optional
	PasswordHistoryEnabled *bool `json:"passwordHistoryEnabled,omitempty"`
}

// SamlSettings holds SAML 2.0 identity provider settings for the domain.
type SamlSettings struct {
	// Certificate is the X.509 certificate for SAML assertion signing.
	// +kubebuilder:validation:Optional
	Certificate *string `json:"certificate,omitempty"`

	// Enabled controls whether the SAML 2.0 IdP is enabled.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// EntityId is the SAML entity identifier for this domain.
	// +kubebuilder:validation:Optional
	EntityId *string `json:"entityId,omitempty"`
}

// SCIMSettings is the SCIM 2.0 provisioning configuration.
type SCIMSettings struct {
	// Enabled controls whether the SCIM 2.0 endpoints are enabled.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// IdpSelectionEnabled controls whether an identity provider can be
	// selected for SCIM-provisioned users.
	// +kubebuilder:validation:Optional
	IdpSelectionEnabled *bool `json:"idpSelectionEnabled,omitempty"`

	// IdpSelectionRule is an expression that selects the identity provider
	// for SCIM-provisioned users.
	// +kubebuilder:validation:Optional
	IdpSelectionRule *string `json:"idpSelectionRule,omitempty"`
}

// SecretExpirationSettings controls client secret expiration.
type SecretExpirationSettings struct {
	// Enabled controls whether client secrets expire.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// ExpiryTimeSeconds is the lifetime, in seconds, of a client secret.
	// +kubebuilder:validation:Optional
	ExpiryTimeSeconds *int64 `json:"expiryTimeSeconds,omitempty"`
}

// SelfServiceAccountManagementSettings controls end-user self-service account management.
type SelfServiceAccountManagementSettings struct {
	// Enabled controls whether self-service account management is enabled.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// ResetPassword holds settings for user-initiated password resets.
	// +kubebuilder:validation:Optional
	ResetPassword *ResetPasswordSettings `json:"resetPassword,omitempty"`
}

// ResetPasswordSettings controls user-initiated password resets.
type ResetPasswordSettings struct {
	// OldPasswordRequired controls whether the old password must be
	// provided when resetting.
	// +kubebuilder:validation:Optional
	OldPasswordRequired *bool `json:"oldPasswordRequired,omitempty"`

	// TokenAge is the validity period, in seconds, of the password reset token.
	// +kubebuilder:validation:Optional
	TokenAge *int32 `json:"tokenAge,omitempty"`
}

// TokenExchangeSettings is the OAuth 2.0 Token Exchange (RFC 8693) configuration.
type TokenExchangeSettings struct {
	// AllowDelegation controls whether token delegation is allowed.
	// +kubebuilder:validation:Optional
	AllowDelegation *bool `json:"allowDelegation,omitempty"`

	// AllowImpersonation controls whether token impersonation is allowed.
	// +kubebuilder:validation:Optional
	AllowImpersonation *bool `json:"allowImpersonation,omitempty"`

	// AllowedActorTokenTypes lists the token types accepted as the actor_token.
	// +kubebuilder:validation:Optional
	AllowedActorTokenTypes []string `json:"allowedActorTokenTypes,omitempty"`

	// AllowedRequestedTokenTypes lists the token types that may be requested.
	// +kubebuilder:validation:Optional
	AllowedRequestedTokenTypes []string `json:"allowedRequestedTokenTypes,omitempty"`

	// AllowedSubjectTokenTypes lists the token types accepted as the subject_token.
	// +kubebuilder:validation:Optional
	AllowedSubjectTokenTypes []string `json:"allowedSubjectTokenTypes,omitempty"`

	// Enabled controls whether token exchange is enabled.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// IdJagSettings is the ID-JAG issuance behavior of token exchange.
	// +kubebuilder:validation:Optional
	IdJagSettings *IdJagSettings `json:"idJagSettings,omitempty"`

	// MaxDelegationDepth is the maximum depth of delegation chains.
	// +kubebuilder:validation:Optional
	MaxDelegationDepth *int32 `json:"maxDelegationDepth,omitempty"`

	// TokenExchangeOAuthSettings are the OAuth-specific token-exchange behavior,
	// with optional inheritance from the domain defaults.
	// +kubebuilder:validation:Optional
	TokenExchangeOAuthSettings *TokenExchangeOAuthSettings `json:"tokenExchangeOAuthSettings,omitempty"`
}

// TokenExchangeOAuthSettings is the OAuth-specific token-exchange behavior.
type TokenExchangeOAuthSettings struct {
	// Inherited controls whether these settings are inherited from the domain defaults.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`

	// ScopeHandling is how scopes are handled when issuing the exchanged token.
	// downscoping restricts the issued token to a subset of the original scopes.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=downscoping;permissive
	ScopeHandling *string `json:"scopeHandling,omitempty"`
}

// IdJagSettings is the ID-JAG issuance behavior of token exchange.
type IdJagSettings struct {
	// LaxValidation also accepts an access token issued to the requesting client
	// as the subject token. By default only an ID token is accepted.
	// +kubebuilder:validation:Optional
	LaxValidation *bool `json:"laxValidation,omitempty"`
}

// UMASettings is the User-Managed Access (UMA 2.0) configuration.
type UMASettings struct {
	// Enabled controls whether UMA 2.0 is enabled.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`
}

// VirtualHost is a virtual host the domain is exposed on.
type VirtualHost struct {
	// Host is the hostname.
	// +kubebuilder:validation:Optional
	Host *string `json:"host,omitempty"`

	// OverrideEntrypoint controls whether the virtual host overrides the
	// domain's default entrypoint.
	// +kubebuilder:validation:Optional
	OverrideEntrypoint *bool `json:"overrideEntrypoint,omitempty"`

	// Path is the context path for this virtual host.
	// +kubebuilder:validation:Optional
	Path *string `json:"path,omitempty"`
}

// WebAuthnSettings is the WebAuthn (FIDO2) relying-party configuration.
type WebAuthnSettings struct {
	// AttestationConveyancePreference controls how attestation data is conveyed.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=direct;indirect;none
	AttestationConveyancePreference *string `json:"attestationConveyancePreference,omitempty"`

	// AuthenticatorAttachment constrains the type of authenticator allowed.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=cross_platform;platform
	AuthenticatorAttachment *string `json:"authenticatorAttachment,omitempty"`

	// Certificates maps certificate aliases to their PEM-encoded values.
	// +kubebuilder:validation:Optional
	Certificates map[string]string `json:"certificates,omitempty"`

	// EnforceAuthenticatorIntegrity controls whether authenticator integrity
	// is enforced on each use.
	// +kubebuilder:validation:Optional
	EnforceAuthenticatorIntegrity *bool `json:"enforceAuthenticatorIntegrity,omitempty"`

	// EnforceAuthenticatorIntegrityMaxAge is the maximum age, in seconds,
	// before an authenticator integrity re-check is required.
	// +kubebuilder:validation:Optional
	EnforceAuthenticatorIntegrityMaxAge *int32 `json:"enforceAuthenticatorIntegrityMaxAge,omitempty"`

	// ForceRegistration controls whether WebAuthn registration is forced
	// on every login.
	// +kubebuilder:validation:Optional
	ForceRegistration *bool `json:"forceRegistration,omitempty"`

	// Origin is the expected origin for WebAuthn assertions.
	// +kubebuilder:validation:Optional
	Origin *string `json:"origin,omitempty"`

	// RelyingPartyId is the relying party identifier.
	// +kubebuilder:validation:Optional
	RelyingPartyId *string `json:"relyingPartyId,omitempty"`

	// RelyingPartyName is the human-readable relying party name.
	// +kubebuilder:validation:Optional
	RelyingPartyName *string `json:"relyingPartyName,omitempty"`

	// RequireResidentKey controls whether a resident key (discoverable credential) is required.
	// +kubebuilder:validation:Optional
	RequireResidentKey *bool `json:"requireResidentKey,omitempty"`

	// UserVerification controls the user verification requirement.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=discouraged;preferred;required
	UserVerification *string `json:"userVerification,omitempty"`
}

// WebProtectionSettings are HTTP security headers for login and consent pages.
type WebProtectionSettings struct {
	// Csp is the Content Security Policy configuration.
	// +kubebuilder:validation:Optional
	Csp *CspSettings `json:"csp,omitempty"`

	// Xframe is the X-Frame-Options configuration.
	// +kubebuilder:validation:Optional
	Xframe *XFrameSettings `json:"xframe,omitempty"`

	// Xss is the X-XSS-Protection configuration.
	// +kubebuilder:validation:Optional
	Xss *XssProtectionSettings `json:"xss,omitempty"`
}

// CspSettings is the Content Security Policy configuration.
type CspSettings struct {
	// Directives are the CSP directives, one per entry, in the form "directive-name value"
	// (e.g. "default-src 'self'"). Directives that take no value may be supplied on their own.
	// +kubebuilder:validation:Optional
	Directives []string `json:"directives,omitempty"`

	// Enabled controls whether CSP headers are sent.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// Inherited controls whether CSP settings are inherited from the parent.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`

	// ReportOnly controls whether the CSP is enforced or report-only.
	// +kubebuilder:validation:Optional
	ReportOnly *bool `json:"reportOnly,omitempty"`

	// ScriptInlineNonce controls whether inline scripts are allowed via a per-request nonce.
	// +kubebuilder:validation:Optional
	ScriptInlineNonce *bool `json:"scriptInlineNonce,omitempty"`
}

// XFrameSettings is the X-Frame-Options configuration.
type XFrameSettings struct {
	// Action is the X-Frame-Options value (DENY, SAMEORIGIN, etc.).
	// +kubebuilder:validation:Optional
	Action *string `json:"action,omitempty"`

	// Enabled controls whether the X-Frame-Options header is sent.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// Inherited controls whether X-Frame-Options settings are inherited from the parent.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`
}

// XssProtectionSettings is the X-XSS-Protection configuration.
type XssProtectionSettings struct {
	// Action is the X-XSS-Protection value.
	// +kubebuilder:validation:Optional
	Action *string `json:"action,omitempty"`

	// Enabled controls whether the X-XSS-Protection header is sent.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// Inherited controls whether X-XSS-Protection settings are inherited from the parent.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`
}
