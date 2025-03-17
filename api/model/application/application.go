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

package application

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

var _ core.ApplicationModel = &Application{}
var _ core.ApplicationSettings = &Setting{}

type SimpleSettings struct {
	// Application Type
	AppType string `json:"type"`
	// ClientID is the client id of the application
	ClientID *string `json:"clientId,omitempty"`
}

// +kubebuilder:validation:Enum=authorization_code;client_credentials;refresh_token;password;implicit
type GrantType string

// +kubebuilder:validation:Enum=BACKEND_TO_BACKEND;NATIVE;BROWSER;WEB
type OauthType string

const (
	GrantTypeClientCredentials GrantType = "client_credentials"
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeRefreshToken      GrantType = "refresh_token"
	GrantTypePassword          GrantType = "password"
	GrantTypeImplicit          GrantType = "implicit"
)

type OAuthClientSettings struct {
	// Oauth client application type
	// +kubebuilder:validation:Required
	ApplicationType OauthType `json:"applicationType"`
	// List of Oauth client grant types
	GrantTypes []GrantType `json:"grantTypes"`
	// List of Oauth client redirect uris
	// +kubebuilder:validation:Optional
	RedirectUris []string `json:"redirectUris"`
}

// TLS settings are used to configure client side TLS in order
// to be able to subscribe to a MTLS plan.
type TLSSettings struct {
	// This client certificate is mandatory to subscribe to a TLS plan.
	// +kubebuilder:validation:Required
	ClientCertificate string `json:"clientCertificate"`
}

type Setting struct {
	App   *SimpleSettings      `json:"app,omitempty"`
	Oauth *OAuthClientSettings `json:"oauth,omitempty"`
	TLS   *TLSSettings         `json:"tls,omitempty"`
}

// HasTLS implements core.ApplicationSettings.
func (in *Setting) HasTLS() bool {
	return in.TLS != nil
}

func (in *Setting) GetClientCertificate() string {
	return in.TLS.ClientCertificate
}

// IsOAuth implements core.ApplicationSettings.
func (in *Setting) IsOAuth() bool {
	return in.Oauth != nil
}

// GetOAuthType implements core.ApplicationSettings.
func (in *Setting) GetOAuthType() string {
	if !in.IsOAuth() {
		return ""
	}
	return string(in.Oauth.ApplicationType)
}

// IsSimple implements core.ApplicationSettings.
func (in *Setting) IsSimple() bool {
	return in.App != nil
}

func (in *Setting) GetClientID() string {
	if in.App != nil {
		clientID := in.App.ClientID
		return *clientID
	}
	return ""
}

// +kubebuilder:validation:Enum=STRING;NUMERIC;BOOLEAN;DATE;MAIL;URL;
type MetaDataFormat string

type Metadata struct {
	// +kubebuilder:validation:Required
	// Metadata Name
	Name string `json:"name"`
	// Metadata Value
	// +kubebuilder:validation:Optional
	Value *string `json:"value,omitempty"`
	// Metadata DefaultValue
	// +kubebuilder:validation:Optional
	DefaultValue *string `json:"defaultValue,omitempty"`
	// Metadata Format
	Format *MetaDataFormat `json:"format,omitempty"`
	// Metadata is hidden or not?
	// +kubebuilder:validation:Optional
	Hidden *bool `json:"hidden,omitempty"`
}

type Member struct {
	// Member source
	// +kubebuilder:validation:Required
	// +kubebuilder:example:=gravitee
	Source string `json:"source"`
	// Member source ID
	// +kubebuilder:validation:Required
	// +kubebuilder:example:=user@email.com
	SourceID string `json:"sourceId"`
	// The API role associated with this Member
	// +kubebuilder:default:=USER
	Role string `json:"role,omitempty"`
}

type Application struct {
	// +kubebuilder:validation:Required
	// Application name
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	// Application Description
	Description string `json:"description"`
	// io.gravitee.definition.model.Application
	// Application ID
	ID string `json:"id,omitempty"`
	// The base64 encoded background to use for this application when displaying it on the portal
	// +kubebuilder:validation:Optional
	Background *string `json:"background,omitempty"`
	// Application domain
	// +kubebuilder:validation:Optional
	Domain *string `json:"domain,omitempty"`
	// Application groups
	// +kubebuilder:validation:Optional
	Groups []string `json:"groups"`
	// The base64 encoded picture to use for this application when displaying it on the portal (if not relying on an URL)
	// +kubebuilder:validation:Optional
	Picture *string `json:"picture,omitempty"`
	// A URL pointing to the picture to use when displaying the application on the portal
	// +kubebuilder:validation:Optional
	PictureURL *string `json:"pictureUrl,omitempty"`
	// Application settings
	// +kubebuilder:validation:Required
	Settings *Setting `json:"settings"`
	// +kubebuilder:validation:Optional
	// Notify members when they are added to the application
	NotifyMembers *bool `json:"notifyMembers"`
	// Application metadata
	// +kubebuilder:validation:Optional
	Metadata *[]Metadata `json:"metadata"`
	// Application members
	// +kubebuilder:validation:Optional
	Members *[]Member `json:"members"`
}

// GetSettings implements core.ApplicationModel.
func (in *Application) GetSettings() core.ApplicationSettings {
	return in.Settings
}
