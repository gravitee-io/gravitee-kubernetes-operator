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

import "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"

// +kubebuilder:validation:Enum=SHARED;EXCLUSIVE;UNSPECIFIED;
type KeyMode string

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

type MetaData struct {
	// +kubebuilder:validation:Required
	// Meta data Name
	Name string `json:"name"`
	// Meta data ApplicationId
	ApplicationId string `json:"applicationId,omitempty"`
	// Meta data Value
	Value string `json:"value,omitempty"`
	// Meta data DefaultValue
	DefaultValue string `json:"defaultValue,omitempty"`
	// Meta data Format
	Format *MetaDataFormat `json:"format,omitempty"`
	// Meta data is hidden or not?
	Hidden bool `json:"hidden,omitempty"`
	// Meta data Key
	Key string `json:"key,omitempty"`
}

type Application struct {
	// +kubebuilder:validation:Required
	// Application name
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	// Application Description
	Description string `json:"description"`
	// Application Type
	ApplicationType string `json:"type,omitempty"`
	// The ClientId identifying the application. This field is required when subscribing to an OAUTH2 / JWT plan.
	ClientId string `json:"clientId,omitempty"`
	// List of application Redirect Uris
	RedirectUris []string `json:"redirectUris,omitempty"`
	// Application Metadata, a map of arbitrary key-values
	Metadata *utils.GenericStringMap `json:"metadata,omitempty"`

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
	Settings *Setting `json:"settings,omitempty"`
	// The API key mode to use. If shared, the application will reuse the same API key across various subscriptions.
	AppKeyMode *KeyMode `json:"app_key_mode,omitempty"`
	// +kubebuilder:validation:Optional
	// Should members get notified when they are added to the application ?
	DisableMembershipNotifications bool `json:"disable_membership_notifications"`
	// Application meta data
	ApplicationMetaData *[]MetaData `json:"applicationMetaData,omitempty"`
}
