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
	// AccountBlockedDuration is the duration, in seconds, for which the account remains blocked after
	// too many failed login attempts.
	// +kubebuilder:validation:Optional
	AccountBlockedDuration *int32 `json:"accountBlockedDuration,omitempty"`

	// AutoLoginAfterRegistration controls whether the user is automatically logged in after completing
	// registration. Defaults to false.
	// +kubebuilder:validation:Optional
	AutoLoginAfterRegistration *bool `json:"autoLoginAfterRegistration,omitempty"`

	// AutoLoginAfterResetPassword controls whether the user is automatically logged in after a
	// password reset. Defaults to false.
	// +kubebuilder:validation:Optional
	AutoLoginAfterResetPassword *bool `json:"autoLoginAfterResetPassword,omitempty"`

	// CompleteRegistrationWhenResetPassword controls whether resetting a password also completes a
	// pending registration. Defaults to false.
	// +kubebuilder:validation:Optional
	CompleteRegistrationWhenResetPassword *bool `json:"completeRegistrationWhenResetPassword,omitempty"`

	// DefaultIdentityProviderForRegistration is the key of an identity provider that exists under this
	// domain, used as the default for user registration. Resolved against the domain's identity
	// providers when applied; a value that does not match an existing identity provider is rejected
	// with a 400 response.
	// +kubebuilder:validation:Optional
	DefaultIdentityProviderForRegistration *string `json:"defaultIdentityProviderForRegistration,omitempty"`

	// DeletePasswordlessDevicesAfterResetPassword controls whether passwordless (WebAuthn) devices are
	// deleted when the password is reset. Defaults to false.
	// +kubebuilder:validation:Optional
	DeletePasswordlessDevicesAfterResetPassword *bool `json:"deletePasswordlessDevicesAfterResetPassword,omitempty"`

	// DynamicUserRegistration controls whether dynamic (self-service) user registration is enabled.
	// Defaults to false.
	// +kubebuilder:validation:Optional
	DynamicUserRegistration *bool `json:"dynamicUserRegistration,omitempty"`

	// Inherited controls whether account settings are inherited from the parent (domain). When true,
	// the other fields are ignored. Has no effect when applied to domains. Defaults to true.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`

	// LoginAttemptsDetectionEnabled controls whether brute-force authentication attempts are detected
	// and blocked. Defaults to false.
	// +kubebuilder:validation:Optional
	LoginAttemptsDetectionEnabled *bool `json:"loginAttemptsDetectionEnabled,omitempty"`

	// LoginAttemptsResetTime is the time, in seconds, after which the login attempt counter is reset
	// when the maximum has not been reached.
	// +kubebuilder:validation:Optional
	LoginAttemptsResetTime *int32 `json:"loginAttemptsResetTime,omitempty"`

	// MaxLoginAttempts is the maximum number of failed login attempts before the account is blocked.
	// +kubebuilder:validation:Optional
	MaxLoginAttempts *int32 `json:"maxLoginAttempts,omitempty"`

	// MfaChallengeAttemptsDetectionEnabled controls whether failed MFA challenge attempts are detected
	// and blocked. Defaults to false.
	// +kubebuilder:validation:Optional
	MfaChallengeAttemptsDetectionEnabled *bool `json:"mfaChallengeAttemptsDetectionEnabled,omitempty"`

	// MfaChallengeAttemptsResetTime is the time, in seconds, after which the MFA challenge attempt
	// counter is reset.
	// +kubebuilder:validation:Optional
	MfaChallengeAttemptsResetTime *int32 `json:"mfaChallengeAttemptsResetTime,omitempty"`

	// MfaChallengeMaxAttempts is the maximum number of failed MFA challenge attempts before the user
	// is blocked.
	// +kubebuilder:validation:Optional
	MfaChallengeMaxAttempts *int32 `json:"mfaChallengeMaxAttempts,omitempty"`

	// MfaChallengeSendVerifyAlertEmail controls whether to send an alert email after too many failed
	// MFA challenge attempts. Defaults to false.
	// +kubebuilder:validation:Optional
	MfaChallengeSendVerifyAlertEmail *bool `json:"mfaChallengeSendVerifyAlertEmail,omitempty"`

	// RedirectUriAfterRegistration is the URL the user is redirected to after registration.
	// +kubebuilder:validation:Optional
	RedirectUriAfterRegistration *string `json:"redirectUriAfterRegistration,omitempty"`

	// RedirectUriAfterResetPassword is the URL the user is redirected to after a password reset.
	// +kubebuilder:validation:Optional
	RedirectUriAfterResetPassword *string `json:"redirectUriAfterResetPassword,omitempty"`

	// RememberMe controls whether users can remain logged in for a fixed duration (remember-me).
	// Defaults to false.
	// +kubebuilder:validation:Optional
	RememberMe *bool `json:"rememberMe,omitempty"`

	// RememberMeDuration is the duration, in seconds, for which a remembered session stays valid.
	// +kubebuilder:validation:Optional
	RememberMeDuration *int32 `json:"rememberMeDuration,omitempty"`

	// ResetPasswordConfirmIdentity controls whether the user must confirm their identity before
	// resetting a password. Defaults to false.
	// +kubebuilder:validation:Optional
	ResetPasswordConfirmIdentity *bool `json:"resetPasswordConfirmIdentity,omitempty"`

	// ResetPasswordCustomForm controls whether a custom form is used for the password-reset step.
	// Defaults to false.
	// +kubebuilder:validation:Optional
	ResetPasswordCustomForm *bool `json:"resetPasswordCustomForm,omitempty"`

	// ResetPasswordCustomFormFields lists the custom fields rendered on the password-reset form.
	// +kubebuilder:validation:Optional
	ResetPasswordCustomFormFields []FormField `json:"resetPasswordCustomFormFields,omitempty"`

	// ResetPasswordInvalidateTokens controls whether existing tokens are invalidated when the password
	// is reset. Defaults to false.
	// +kubebuilder:validation:Optional
	ResetPasswordInvalidateTokens *bool `json:"resetPasswordInvalidateTokens,omitempty"`

	// SendRecoverAccountEmail controls whether to send an account-recovery email. Defaults to false.
	// +kubebuilder:validation:Optional
	SendRecoverAccountEmail *bool `json:"sendRecoverAccountEmail,omitempty"`

	// SendVerifyRegistrationAccountEmail controls whether to send a registration-verification email.
	// Defaults to false.
	// +kubebuilder:validation:Optional
	SendVerifyRegistrationAccountEmail *bool `json:"sendVerifyRegistrationAccountEmail,omitempty"`
}

// CertificateSettings are domain-level certificate settings.
type CertificateSettings struct {
	// FallbackCertificate is the key of a certificate managed under this domain, used as the fallback
	// certificate when a client does not specify one. Must reference a certificate created via the
	// domain's certificate endpoints.
	// +kubebuilder:validation:Optional
	FallbackCertificate *string `json:"fallbackCertificate,omitempty"`
}

// CorsSettings is the Cross-Origin Resource Sharing configuration for the domain.
type CorsSettings struct {
	// AllowCredentials controls whether the browser may send credentials (cookies, authorization
	// headers) with cross-origin requests. Defaults to false.
	// +kubebuilder:validation:Optional
	AllowCredentials *bool `json:"allowCredentials,omitempty"`

	// AllowedHeaders lists the request headers permitted on cross-origin requests.
	// +kubebuilder:validation:Optional
	// +listType=set
	AllowedHeaders []string `json:"allowedHeaders,omitempty"`

	// AllowedMethods lists the HTTP methods permitted on cross-origin requests.
	// +kubebuilder:validation:Optional
	// +listType=set
	AllowedMethods []string `json:"allowedMethods,omitempty"`

	// AllowedOrigins lists the origins permitted to make cross-origin requests. Use "*" to allow any
	// origin.
	// +kubebuilder:validation:Optional
	// +listType=set
	AllowedOrigins []string `json:"allowedOrigins,omitempty"`

	// Enabled controls whether CORS handling is enabled for the domain when not inherited. Defaults to
	// false.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// Inherited controls whether CORS settings are inherited from the gateway defaults (gravitee.yml).
	// When null, legacy behaviour applies: enabled=true overrides and enabled=false inherits. Defaults
	// to true.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`

	// MaxAge is how long, in seconds, a browser may cache the result of a preflight request. Defaults
	// to 86400.
	// +kubebuilder:validation:Optional
	MaxAge *int32 `json:"maxAge,omitempty"`
}

// FormField is a custom field rendered on a domain form.
type FormField struct {
	// Key is the identifier of the field, mapped to a user attribute.
	// +kubebuilder:validation:Optional
	Key *string `json:"key,omitempty"`

	// Label is the label displayed for the field.
	// +kubebuilder:validation:Optional
	Label *string `json:"label,omitempty"`

	// Type is the input type of the field.
	// +kubebuilder:validation:Optional
	Type *string `json:"type,omitempty"`
}

// KeyRetrievalSettings are the fetch, SSRF and cache limits
// applied to every trusted domain in the security domain.
type KeyRetrievalSettings struct {
	// AllowPrivateIpAddress controls whether key material can be fetched from private IP addresses.
	// Defaults to false.
	// +kubebuilder:validation:Optional
	AllowPrivateIpAddress *bool `json:"allowPrivateIpAddress,omitempty"`

	// AllowUnsecuredHttpUri controls whether key material can be fetched over unsecured HTTP URIs.
	// Defaults to false.
	// +kubebuilder:validation:Optional
	AllowUnsecuredHttpUri *bool `json:"allowUnsecuredHttpUri,omitempty"`

	// CacheMaxEntries is the maximum number of key material entries retained in the cache. Defaults to
	// 50.
	// +kubebuilder:validation:Optional
	CacheMaxEntries *int32 `json:"cacheMaxEntries,omitempty"`

	// CacheTtlSeconds is the time-to-live, in seconds, for cached key material. Defaults to 300.
	// +kubebuilder:validation:Optional
	CacheTtlSeconds *int32 `json:"cacheTtlSeconds,omitempty"`

	// FetchTimeoutMs is the timeout, in milliseconds, for fetching key material. Defaults to 5000.
	// +kubebuilder:validation:Optional
	FetchTimeoutMs *int32 `json:"fetchTimeoutMs,omitempty"`

	// MaxResponseSizeKb is the maximum key material response size, in kilobytes. Defaults to 32.
	// +kubebuilder:validation:Optional
	MaxResponseSizeKb *int32 `json:"maxResponseSizeKb,omitempty"`
}

// LoginSettings is the configuration of the domain's login flow.
type LoginSettings struct {
	// CertificateBasedAuthEnabled controls whether certificate-based authentication is offered.
	// Defaults to false.
	// +kubebuilder:validation:Optional
	CertificateBasedAuthEnabled *bool `json:"certificateBasedAuthEnabled,omitempty"`

	// CertificateBasedAuthUrl is the URL used for certificate-based authentication.
	// +kubebuilder:validation:Optional
	CertificateBasedAuthUrl *string `json:"certificateBasedAuthUrl,omitempty"`

	// ForgotPasswordEnabled controls whether users can initiate a forgot-password flow from the login
	// page. Defaults to false.
	// +kubebuilder:validation:Optional
	ForgotPasswordEnabled *bool `json:"forgotPasswordEnabled,omitempty"`

	// HideForm controls whether the login form is hidden (for example when only social or
	// identifier-first login is offered). Defaults to false.
	// +kubebuilder:validation:Optional
	HideForm *bool `json:"hideForm,omitempty"`

	// IdentifierFirstEnabled controls whether identifier-first login is enabled, prompting for the
	// username before the password. Defaults to false.
	// +kubebuilder:validation:Optional
	IdentifierFirstEnabled *bool `json:"identifierFirstEnabled,omitempty"`

	// Inherited controls whether these login settings are inherited from a parent scope rather than
	// defined here. When true, the other fields are ignored. Defaults to true.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`

	// MagicLinkAuthEnabled controls whether magic-link authentication is offered. Defaults to false.
	// +kubebuilder:validation:Optional
	MagicLinkAuthEnabled *bool `json:"magicLinkAuthEnabled,omitempty"`

	// PasswordlessDeviceNamingEnabled controls whether users can name their passwordless devices.
	// Defaults to false.
	// +kubebuilder:validation:Optional
	PasswordlessDeviceNamingEnabled *bool `json:"passwordlessDeviceNamingEnabled,omitempty"`

	// PasswordlessEnabled controls whether passwordless (WebAuthn) authentication is offered. Defaults
	// to false.
	// +kubebuilder:validation:Optional
	PasswordlessEnabled *bool `json:"passwordlessEnabled,omitempty"`

	// PasswordlessEnforcePasswordEnabled controls whether a password is still required alongside
	// passwordless authentication. Defaults to false.
	// +kubebuilder:validation:Optional
	PasswordlessEnforcePasswordEnabled *bool `json:"passwordlessEnforcePasswordEnabled,omitempty"`

	// PasswordlessEnforcePasswordMaxAge is the period, in seconds, after which the user's credentials
	// must be re-entered to keep using passwordless authentication.
	// +kubebuilder:validation:Optional
	PasswordlessEnforcePasswordMaxAge *int32 `json:"passwordlessEnforcePasswordMaxAge,omitempty"`

	// PasswordlessRememberDeviceEnabled controls whether a passwordless device can be remembered to
	// skip future challenges. Defaults to false.
	// +kubebuilder:validation:Optional
	PasswordlessRememberDeviceEnabled *bool `json:"passwordlessRememberDeviceEnabled,omitempty"`

	// RegisterEnabled controls whether users can self-register from the login page. Defaults to false.
	// +kubebuilder:validation:Optional
	RegisterEnabled *bool `json:"registerEnabled,omitempty"`

	// RememberMeEnabled controls whether the login page offers a remember-me option. Defaults to
	// false.
	// +kubebuilder:validation:Optional
	RememberMeEnabled *bool `json:"rememberMeEnabled,omitempty"`

	// ResetPasswordOnExpiration controls whether the user is forced to reset their password once it
	// expires.
	// +kubebuilder:validation:Optional
	ResetPasswordOnExpiration *bool `json:"resetPasswordOnExpiration,omitempty"`
}

// OidcSettings holds OpenID Connect settings for the domain.
type OidcSettings struct {
	// CibaSettings holds the Client-Initiated Backchannel Authentication (CIBA) settings for the
	// domain. CIBA lets a relying party initiate end-user authentication from a separate consumption
	// device, without redirecting the user through the browser. Authentication device notifiers are
	// not managed by the Automation API and are not exposed here.
	// +kubebuilder:validation:Optional
	CibaSettings *CIBASettings `json:"cibaSettings,omitempty"`

	// ClientRegistrationSettings holds the OpenID Connect Dynamic Client Registration configuration
	// for the domain.
	// +kubebuilder:validation:Optional
	ClientRegistrationSettings *ClientRegistrationSettings `json:"clientRegistrationSettings,omitempty"`

	// PostLogoutRedirectUris lists the URLs the user may be redirected to after sign-out
	// (post_logout_redirect_uri).
	// +kubebuilder:validation:Optional
	PostLogoutRedirectUris []string `json:"postLogoutRedirectUris,omitempty"`

	// RedirectUriStrictMatching controls whether redirect_uri and post_logout_redirect_uri values are
	// matched strictly during OpenID Connect flows. Defaults to false.
	// +kubebuilder:validation:Optional
	RedirectUriStrictMatching *bool `json:"redirectUriStrictMatching,omitempty"`

	// RequestUris lists the allowed request_uri values for passing OpenID Connect request objects by
	// reference.
	// +kubebuilder:validation:Optional
	RequestUris []string `json:"requestUris,omitempty"`

	// SecurityProfileSettings holds the Financial-grade API (FAPI) security profile configuration for
	// the domain.
	// +kubebuilder:validation:Optional
	SecurityProfileSettings *SecurityProfileSettings `json:"securityProfileSettings,omitempty"`

	// WorkloadIdentitySettings holds the workload identity (SPIFFE) settings for the domain.
	// +kubebuilder:validation:Optional
	WorkloadIdentitySettings *SpiffeDomainSettings `json:"workloadIdentitySettings,omitempty"`
}

// SpiffeDomainSettings are the workload identity (SPIFFE) settings for the domain.
// Key retrieval limits are configured in keyRetrievalSettings.
type SpiffeDomainSettings struct {
	// ClockSkewSeconds is the allowed clock skew, in seconds, when validating JWT temporal claims.
	// Defaults to 30.
	// +kubebuilder:validation:Optional
	ClockSkewSeconds *int32 `json:"clockSkewSeconds,omitempty"`

	// DefaultAllowedAlgorithms is the default allowlist of signature algorithms accepted for SPIFFE
	// JWT validation.
	// +kubebuilder:validation:Optional
	DefaultAllowedAlgorithms []string `json:"defaultAllowedAlgorithms,omitempty"`

	// Enabled controls whether SPIFFE workload identity support is enabled for the domain. Defaults to
	// false.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// MaxJwtLifetimeSeconds is the maximum accepted JWT lifetime, in seconds, computed as exp minus
	// iat. Defaults to 300.
	// +kubebuilder:validation:Optional
	MaxJwtLifetimeSeconds *int32 `json:"maxJwtLifetimeSeconds,omitempty"`
}

// CIBASettings is the Client-Initiated Backchannel Authentication configuration.
type CIBASettings struct {
	// AuthReqExpiry is the default validity period, in seconds, of the issued auth_req_id.
	// +kubebuilder:validation:Optional
	AuthReqExpiry *int32 `json:"authReqExpiry,omitempty"`

	// BindingMessageLength is the maximum number of characters accepted for the binding_message
	// parameter.
	// +kubebuilder:validation:Optional
	BindingMessageLength *int32 `json:"bindingMessageLength,omitempty"`

	// Enabled controls whether Client-Initiated Backchannel Authentication is enabled for the domain.
	// Defaults to false.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// TokenReqInterval is the minimum delay, in seconds, that a client must wait between two polls of
	// the token endpoint for the same auth_req_id (POLL or PING delivery mode).
	// +kubebuilder:validation:Optional
	TokenReqInterval *int32 `json:"tokenReqInterval,omitempty"`
}

// ClientRegistrationSettings is the OpenID Connect Dynamic Client Registration configuration.
type ClientRegistrationSettings struct {
	// AllowHttpSchemeRedirectUri controls whether the unsecured http scheme is permitted in redirect
	// URIs. Defaults to false.
	// +kubebuilder:validation:Optional
	AllowHttpSchemeRedirectUri *bool `json:"allowHttpSchemeRedirectUri,omitempty"`

	// AllowLocalhostRedirectUri controls whether localhost is permitted as a redirect URI host.
	// Defaults to false.
	// +kubebuilder:validation:Optional
	AllowLocalhostRedirectUri *bool `json:"allowLocalhostRedirectUri,omitempty"`

	// AllowRedirectUriParamsExpressionLanguage controls whether expression language is permitted in
	// redirect URI parameters. Defaults to false.
	// +kubebuilder:validation:Optional
	AllowRedirectUriParamsExpressionLanguage *bool `json:"allowRedirectUriParamsExpressionLanguage,omitempty"`

	// AllowWildCardRedirectUri controls whether wildcards are permitted in redirect URIs. Defaults to
	// false.
	// +kubebuilder:validation:Optional
	AllowWildCardRedirectUri *bool `json:"allowWildCardRedirectUri,omitempty"`

	// AllowedScopes lists the scopes permitted on client registration requests when the allowed list
	// is enabled.
	// +kubebuilder:validation:Optional
	AllowedScopes []string `json:"allowedScopes,omitempty"`

	// AllowedScopesEnabled controls whether registered client scopes are restricted to an allowed
	// list. Defaults to false.
	// +kubebuilder:validation:Optional
	AllowedScopesEnabled *bool `json:"allowedScopesEnabled,omitempty"`

	// ClientTemplateEnabled controls whether a client may be used as a template for dynamic client
	// registration. Defaults to false.
	// +kubebuilder:validation:Optional
	ClientTemplateEnabled *bool `json:"clientTemplateEnabled,omitempty"`

	// DefaultScopes lists the default scopes added to every client registration request.
	// +kubebuilder:validation:Optional
	DefaultScopes []string `json:"defaultScopes,omitempty"`

	// DynamicClientRegistrationEnabled controls whether Dynamic Client Registration is enabled for the
	// domain. Defaults to false.
	// +kubebuilder:validation:Optional
	DynamicClientRegistrationEnabled *bool `json:"dynamicClientRegistrationEnabled,omitempty"`

	// OpenDynamicClientRegistrationEnabled controls whether open (unauthenticated) Dynamic Client
	// Registration is enabled for the domain. Defaults to false.
	// +kubebuilder:validation:Optional
	OpenDynamicClientRegistrationEnabled *bool `json:"openDynamicClientRegistrationEnabled,omitempty"`
}

// SecurityProfileSettings holds FAPI security profile settings.
type SecurityProfileSettings struct {
	// EnableFapiBrazil controls whether the Open Banking Brasil Financial-grade API security profile
	// (version 1.0) is applied. Defaults to false.
	// +kubebuilder:validation:Optional
	EnableFapiBrazil *bool `json:"enableFapiBrazil,omitempty"`

	// EnablePlainFapi controls whether the standard Financial-grade API security profile (version 1.0)
	// is applied. Defaults to false.
	// +kubebuilder:validation:Optional
	EnablePlainFapi *bool `json:"enablePlainFapi,omitempty"`
}

// PasswordSettings is the password policy applied to domain users.
type PasswordSettings struct {
	// ExcludePasswordsInDictionary controls whether passwords found in a common-password dictionary
	// are rejected.
	// +kubebuilder:validation:Optional
	ExcludePasswordsInDictionary *bool `json:"excludePasswordsInDictionary,omitempty"`

	// ExcludeUserProfileInfoInPassword controls whether passwords containing the user's profile
	// information are rejected.
	// +kubebuilder:validation:Optional
	ExcludeUserProfileInfoInPassword *bool `json:"excludeUserProfileInfoInPassword,omitempty"`

	// ExpiryDuration is the number of days after which a password expires and must be changed.
	// +kubebuilder:validation:Optional
	ExpiryDuration *int32 `json:"expiryDuration,omitempty"`

	// IncludeNumbers controls whether a password must contain at least one number.
	// +kubebuilder:validation:Optional
	IncludeNumbers *bool `json:"includeNumbers,omitempty"`

	// IncludeSpecialCharacters controls whether a password must contain at least one special
	// character.
	// +kubebuilder:validation:Optional
	IncludeSpecialCharacters *bool `json:"includeSpecialCharacters,omitempty"`

	// Inherited controls whether these password settings are inherited from a parent scope rather than
	// defined here. When true, the other fields are ignored. Defaults to true.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`

	// LettersInMixedCase controls whether a password must contain both uppercase and lowercase
	// letters.
	// +kubebuilder:validation:Optional
	LettersInMixedCase *bool `json:"lettersInMixedCase,omitempty"`

	// MaxConsecutiveLetters is the maximum number of identical consecutive characters allowed in a
	// password.
	// +kubebuilder:validation:Optional
	MaxConsecutiveLetters *int32 `json:"maxConsecutiveLetters,omitempty"`

	// MaxLength is the maximum number of characters a password may contain. Defaults to 128.
	// +kubebuilder:validation:Optional
	MaxLength *int32 `json:"maxLength,omitempty"`

	// MinLength is the minimum number of characters a password must contain. Defaults to 8.
	// +kubebuilder:validation:Optional
	MinLength *int32 `json:"minLength,omitempty"`

	// OldPasswords is the number of previous passwords retained in history and barred from reuse.
	// +kubebuilder:validation:Optional
	OldPasswords *int32 `json:"oldPasswords,omitempty"`

	// PasswordHistoryEnabled controls whether password history is enforced to prevent reuse of recent
	// passwords. Defaults to false.
	// +kubebuilder:validation:Optional
	PasswordHistoryEnabled *bool `json:"passwordHistoryEnabled,omitempty"`
}

// SamlSettings holds SAML 2.0 identity provider settings for the domain.
type SamlSettings struct {
	// Certificate is the key of a certificate managed under this domain, used to sign SAML responses.
	// Must reference a certificate created via the domain's certificate endpoints.
	// +kubebuilder:validation:Optional
	Certificate *string `json:"certificate,omitempty"`

	// Enabled controls whether the domain exposes the SAML 2.0 IdP protocol. Defaults to false.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// EntityId is the URL or URN that uniquely identifies this IdP (the SAML entity ID).
	// +kubebuilder:validation:Optional
	EntityId *string `json:"entityId,omitempty"`
}

// SCIMSettings is the SCIM 2.0 provisioning configuration.
type SCIMSettings struct {
	// Enabled controls whether the SCIM provisioning API is enabled for the domain. Defaults to false.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// IdpSelectionEnabled controls whether an identity provider is selected for SCIM-provisioned users
	// using a selection rule. Defaults to false.
	// +kubebuilder:validation:Optional
	IdpSelectionEnabled *bool `json:"idpSelectionEnabled,omitempty"`

	// IdpSelectionRule is the expression that selects the identity provider for a SCIM-provisioned
	// user.
	// +kubebuilder:validation:Optional
	IdpSelectionRule *string `json:"idpSelectionRule,omitempty"`
}

// SecretExpirationSettings controls client secret expiration.
type SecretExpirationSettings struct {
	// Enabled controls whether client-secret expiration is enabled.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// ExpiryTimeSeconds is the lifetime, in seconds, of a client secret before it expires.
	// +kubebuilder:validation:Optional
	ExpiryTimeSeconds *int64 `json:"expiryTimeSeconds,omitempty"`
}

// SelfServiceAccountManagementSettings controls end-user self-service account management.
type SelfServiceAccountManagementSettings struct {
	// Enabled controls whether self-service account management is enabled for end users. Defaults to
	// false.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// ResetPassword holds the rules applied to a self-service password reset.
	// +kubebuilder:validation:Optional
	ResetPassword *ResetPasswordSettings `json:"resetPassword,omitempty"`
}

// ResetPasswordSettings controls user-initiated password resets.
type ResetPasswordSettings struct {
	// OldPasswordRequired controls whether the user must supply their current password to set a new
	// one. Defaults to false.
	// +kubebuilder:validation:Optional
	OldPasswordRequired *bool `json:"oldPasswordRequired,omitempty"`

	// TokenAge is the lifetime, in seconds, of the password-reset token.
	// +kubebuilder:validation:Optional
	TokenAge *int32 `json:"tokenAge,omitempty"`
}

// TokenExchangeSettings is the OAuth 2.0 Token Exchange (RFC 8693) configuration.
type TokenExchangeSettings struct {
	// AllowDelegation controls whether delegation is allowed, where an actor acts on behalf of the
	// subject and an "act" claim is added to the issued token. At least one of allowImpersonation or
	// allowDelegation must be enabled. Defaults to false.
	// +kubebuilder:validation:Optional
	AllowDelegation *bool `json:"allowDelegation,omitempty"`

	// AllowImpersonation controls whether impersonation is allowed, where the issued token represents
	// the subject directly. At least one of allowImpersonation or allowDelegation must be enabled.
	// Defaults to true.
	// +kubebuilder:validation:Optional
	AllowImpersonation *bool `json:"allowImpersonation,omitempty"`

	// AllowedActorTokenTypes lists the token types accepted as the actor token when delegating.
	// +kubebuilder:validation:Optional
	AllowedActorTokenTypes []string `json:"allowedActorTokenTypes,omitempty"`

	// AllowedRequestedTokenTypes lists the token types that may be requested as the result of an
	// exchange.
	// +kubebuilder:validation:Optional
	AllowedRequestedTokenTypes []string `json:"allowedRequestedTokenTypes,omitempty"`

	// AllowedSubjectTokenTypes lists the token types accepted as the subject token in an exchange.
	// +kubebuilder:validation:Optional
	AllowedSubjectTokenTypes []string `json:"allowedSubjectTokenTypes,omitempty"`

	// Enabled controls whether token exchange is enabled for the domain. Defaults to false.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// IdJagSettings holds the ID-JAG issuance behavior of token exchange.
	// +kubebuilder:validation:Optional
	IdJagSettings *IdJagSettings `json:"idJagSettings,omitempty"`

	// MaxDelegationDepth is the maximum depth of the delegation chain (nested "act" claims). Clamped
	// to the range 1–100. Defaults to 25.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=100
	MaxDelegationDepth *int32 `json:"maxDelegationDepth,omitempty"`

	// TokenExchangeOAuthSettings holds the OAuth-specific token-exchange behavior, such as how scopes
	// are handled, with optional inheritance from domain defaults.
	// +kubebuilder:validation:Optional
	TokenExchangeOAuthSettings *TokenExchangeOAuthSettings `json:"tokenExchangeOAuthSettings,omitempty"`
}

// TokenExchangeOAuthSettings is the OAuth-specific token-exchange behavior.
type TokenExchangeOAuthSettings struct {
	// Inherited controls whether these settings are inherited from the domain defaults rather than
	// defined here. Defaults to true.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`

	// ScopeHandling is how scopes are handled when issuing the exchanged token. downscoping restricts
	// the issued token to a subset of the original scopes. Defaults to downscoping.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=downscoping;permissive
	ScopeHandling *string `json:"scopeHandling,omitempty"`
}

// IdJagSettings is the ID-JAG issuance behavior of token exchange.
type IdJagSettings struct {
	// LaxValidation also accepts an access token issued to the requesting client as the subject token.
	// By default only an ID token is accepted.
	// +kubebuilder:validation:Optional
	LaxValidation *bool `json:"laxValidation,omitempty"`
}

// UMASettings is the User-Managed Access (UMA 2.0) configuration.
type UMASettings struct {
	// Enabled controls whether User-Managed Access is enabled for the domain. Defaults to false.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`
}

// VirtualHost is a virtual host the domain is exposed on.
type VirtualHost struct {
	// Host is the hostname the domain is served on.
	// +kubebuilder:validation:Optional
	Host *string `json:"host,omitempty"`

	// OverrideEntrypoint controls whether this virtual host overrides the organization entry point.
	// Defaults to false.
	// +kubebuilder:validation:Optional
	OverrideEntrypoint *bool `json:"overrideEntrypoint,omitempty"`

	// Path is the context path the domain is served under on this host.
	// +kubebuilder:validation:Optional
	Path *string `json:"path,omitempty"`
}

// WebAuthnSettings is the WebAuthn (FIDO2) relying-party configuration.
type WebAuthnSettings struct {
	// AttestationConveyancePreference is the relying-party preference for attestation conveyance
	// during credential creation. none requests no attestation, indirect allows anonymized
	// attestation, and direct requests the authenticator's attestation statement. Defaults to none.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=direct;indirect;none
	AttestationConveyancePreference *string `json:"attestationConveyancePreference,omitempty"`

	// AuthenticatorAttachment is the preferred authenticator attachment. platform selects
	// authenticators bound to the device (such as a fingerprint reader); cross_platform selects
	// roaming authenticators (such as a security key).
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=cross_platform;platform
	AuthenticatorAttachment *string `json:"authenticatorAttachment,omitempty"`

	// Certificates holds the trusted device-attestation X.509 certificates, keyed by name.
	// +kubebuilder:validation:Optional
	Certificates map[string]string `json:"certificates,omitempty"`

	// EnforceAuthenticatorIntegrity controls whether to periodically re-verify that registered
	// authenticators remain valid against the FIDO2 Metadata Service. Defaults to false.
	// +kubebuilder:validation:Optional
	EnforceAuthenticatorIntegrity *bool `json:"enforceAuthenticatorIntegrity,omitempty"`

	// EnforceAuthenticatorIntegrityMaxAge is the maximum elapsed time, in seconds, since an
	// authenticator was last verified before it is re-checked on the next passwordless login.
	// +kubebuilder:validation:Optional
	EnforceAuthenticatorIntegrityMaxAge *int32 `json:"enforceAuthenticatorIntegrityMaxAge,omitempty"`

	// ForceRegistration controls whether to reject registration of a credential already registered to
	// a different user. Defaults to false.
	// +kubebuilder:validation:Optional
	ForceRegistration *bool `json:"forceRegistration,omitempty"`

	// Origin is the relying-party origin; must match the browser's window.location.origin during
	// registration and authentication ceremonies.
	// +kubebuilder:validation:Optional
	Origin *string `json:"origin,omitempty"`

	// RelyingPartyId is the relying-party identifier: a domain string that scopes credentials to this
	// entity. A credential can only be used with the relying party it was registered against.
	// +kubebuilder:validation:Optional
	RelyingPartyId *string `json:"relyingPartyId,omitempty"`

	// RelyingPartyName is the human-readable relying-party name shown to users during ceremonies.
	// +kubebuilder:validation:Optional
	RelyingPartyName *string `json:"relyingPartyName,omitempty"`

	// RequireResidentKey controls whether the authenticator must create a client-side resident
	// (discoverable) credential. Defaults to false.
	// +kubebuilder:validation:Optional
	RequireResidentKey *bool `json:"requireResidentKey,omitempty"`

	// UserVerification is the relying-party requirement regarding user verification during a ceremony.
	// required enforces verification, preferred requests it when available, and discouraged avoids it.
	// Defaults to preferred.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=discouraged;preferred;required
	UserVerification *string `json:"userVerification,omitempty"`
}

// WebProtectionSettings are HTTP security headers for login and consent pages.
type WebProtectionSettings struct {
	// Csp holds the Content Security Policy configuration for the domain's login and consent pages.
	// +kubebuilder:validation:Optional
	Csp *CspSettings `json:"csp,omitempty"`

	// Xframe controls whether the domain's pages may be embedded in frames on other origins.
	// +kubebuilder:validation:Optional
	Xframe *XFrameSettings `json:"xframe,omitempty"`

	// Xss controls the legacy X-XSS-Protection response header.
	// +kubebuilder:validation:Optional
	Xss *XssProtectionSettings `json:"xss,omitempty"`
}

// CspSettings is the Content Security Policy configuration.
type CspSettings struct {
	// Directives lists the CSP directives, one per entry, in the form "directive-name value". A
	// trailing semicolon is optional. Directive names must be valid CSP tokens and must not repeat;
	// values are not interpreted. Directives that take no value, such as "upgrade-insecure-requests",
	// may be supplied on their own. When reportOnly is enabled, a "report-uri" or "report-to"
	// directive is required.
	// +kubebuilder:validation:Optional
	Directives []string `json:"directives,omitempty"`

	// Enabled controls whether CSP is enabled for the domain when not inherited. Defaults to false.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// Inherited controls whether CSP settings are inherited from the gateway defaults (gravitee.yml).
	// When null, legacy behaviour applies: enabled=true overrides and enabled=false inherits. Defaults
	// to true.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`

	// ReportOnly delivers the policy as Content-Security-Policy-Report-Only when true.
	// +kubebuilder:validation:Optional
	ReportOnly *bool `json:"reportOnly,omitempty"`

	// ScriptInlineNonce controls whether inline scripts are allowed via a per-request nonce. Defaults
	// to true.
	// +kubebuilder:validation:Optional
	ScriptInlineNonce *bool `json:"scriptInlineNonce,omitempty"`
}

// XFrameSettings is the X-Frame-Options configuration.
type XFrameSettings struct {
	// Action is the X-Frame-Options action. Supported values: DENY, SAMEORIGIN. Leave unset to omit
	// the header.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=DENY;SAMEORIGIN
	Action *string `json:"action,omitempty"`

	// Enabled controls whether X-Frame-Options is enabled for the domain when not inherited. Defaults
	// to false.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// Inherited controls whether X-Frame-Options settings are inherited from the gateway defaults
	// (gravitee.yml). When null, legacy behaviour applies: enabled=true overrides and enabled=false
	// inherits. Defaults to true.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`
}

// XssProtectionSettings is the X-XSS-Protection configuration.
type XssProtectionSettings struct {
	// Action is the value of the X-XSS-Protection header.
	// +kubebuilder:validation:Optional
	Action *string `json:"action,omitempty"`

	// Enabled controls whether X-XSS-Protection is enabled for the domain when not inherited. Defaults
	// to false.
	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty"`

	// Inherited controls whether X-XSS-Protection settings are inherited from the gateway defaults
	// (gravitee.yml). When null, legacy behaviour applies: enabled=true overrides and enabled=false
	// inherits. Defaults to true.
	// +kubebuilder:validation:Optional
	Inherited *bool `json:"inherited,omitempty"`
}
