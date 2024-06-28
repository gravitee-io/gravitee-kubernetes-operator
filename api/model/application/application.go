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

// +kubebuilder:object:generate=true
package application

// +kubebuilder:validation:Enum=SHARED;EXCLUSIVE;UNSPECIFIED;
type AppKeyMode string

type SimpleSettings struct {
	// Application Type
	AppType string `json:"type,omitempty"`
	// ClientId is the client id of the application
	ClientId string `json:"client_id,omitempty"`
}

type OAuthClientSettings struct {
	// Oauth client application type
	ApplicationType string `json:"application_type,omitempty"`
	// Oauth client id
	ClientId string `json:"client_id,omitempty"`
	// Oauth client secret
	ClientSecret string `json:"client_secret,omitempty"`
	// Oauth client uri
	ClientURI string `json:"client_uri,omitempty"`
	// List of Oauth client grant types
	GrantTypes []string `json:"grant_types,omitempty"`
	// Oauth client logo uri
	LogoURI string `json:"logo_uri,omitempty"`
	// List of Oauth client redirect uris
	RedirectUris []string `json:"redirect_uris,omitempty"`
	// Whether client secret renewing is supported or not
	RenewClientSecretSupported bool `json:"renew_client_secret_supported,omitempty"`
	// List of Oauth client response types
	ResponseTypes []string `json:"response_types,omitempty"`
}

type Setting struct {
	App   *SimpleSettings      `json:"app,omitempty"`
	Oauth *OAuthClientSettings `json:"oauth,omitempty"`
}

// +kubebuilder:validation:Enum=STRING;NUMERIC;BOOLEAN;DATE;MAIL;URL;
type MetaDataFormat string

type Metadata struct {
	// +kubebuilder:validation:Required
	// Metadata Name
	Name string `json:"name"`
	// Metadata Value
	Value string `json:"value,omitempty"`
	// Metadata DefaultValue
	DefaultValue string `json:"defaultValue,omitempty"`
	// Metadata Format
	Format *MetaDataFormat `json:"format,omitempty"`
	// Metadata is hidden or not?
	Hidden bool `json:"hidden,omitempty"`
}

type Application struct {
	// +kubebuilder:validation:Required
	// Application name
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	// Application Description
	Description string `json:"description,omitempty"`
	// Application Type
	Type string `json:"type,omitempty"`
	// The ClientId identifying the application. This field is required when subscribing to an OAUTH2 / JWT plan.
	ClientId string `json:"clientId,omitempty"`
	// List of application Redirect Uris
	RedirectUris []string `json:"redirectUris,omitempty"`

	// The origin which is used to create this Application
	// +kubebuilder:validation:Enum=kubernetes;
	Origin string `json:"origin,omitempty"`
	// io.gravitee.definition.model.Application
	// Application ID
	ID string `json:"id,omitempty"`
	// The base64 encoded background to use for this application when displaying it on the portal
	Background string `json:"background,omitempty"`
	// Application domain
	Domain string `json:"domain,omitempty"`
	// Application groups
	Groups []string `json:"groups,omitempty"`
	// The base64 encoded picture to use for this application when displaying it on the portal (if not relying on an URL)
	Picture string `json:"picture,omitempty"`
	// An URL pointing to the picture to use when displaying the application on the portal
	PictureURL string `json:"picture_url,omitempty"`
	// Application settings
	// +kubebuilder:validation:Required
	Settings *Setting `json:"settings"`
	// The API key mode to use. If shared, the application will reuse the same API key across various subscriptions.
	AppKeyMode *AppKeyMode `json:"app_key_mode,omitempty"`
	// Should membership notifications be disabled or not ?
	DisableMembershipNotifications bool `json:"disable_membership_notifications,omitempty"`
	// Application metadata
	Metadata *[]Metadata `json:"metadata,omitempty"`
}
