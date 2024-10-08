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

import "github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"

var _ core.ApplicationModel = &Application{}
var _ core.ApplicationSettings = &Setting{}

type SimpleSettings struct {
	// Application Type
	AppType string `json:"type"`
	// ClientID is the client id of the application
	ClientID string `json:"clientId,omitempty"`
}

type OAuthClientSettings struct {
	// Oauth client application type
	ApplicationType string `json:"applicationType"`
	// List of Oauth client grant types
	GrantTypes []string `json:"grantTypes"`
	// List of Oauth client redirect uris
	RedirectUris []string `json:"redirectUris,omitempty"`
}

type Setting struct {
	App   *SimpleSettings      `json:"app,omitempty"`
	Oauth *OAuthClientSettings `json:"oauth,omitempty"`
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
	return in.Oauth.ApplicationType
}

// IsSimple implements core.ApplicationSettings.
func (in *Setting) IsSimple() bool {
	return in.App != nil
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

type Member struct {
	// Member source
	// +kubebuilder:validation:Required
	// +kubebuilder:example:=gravitee
	Source string `json:"source"`
	// Member source ID
	// +kubebuilder:validation:Required
	// +kubebuilder:example:=user@email.com
	SourceID string `json:"sourceId"`
	// Member display name
	DisplayName string `json:"displayName,omitempty"`
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
	Background string `json:"background,omitempty"`
	// Application domain
	Domain string `json:"domain,omitempty"`
	// Application groups
	Groups []string `json:"groups,omitempty"`
	// The base64 encoded picture to use for this application when displaying it on the portal (if not relying on an URL)
	Picture string `json:"picture,omitempty"`
	// A URL pointing to the picture to use when displaying the application on the portal
	PictureURL string `json:"pictureUrl,omitempty"`
	// Application settings
	// +kubebuilder:validation:Required
	Settings *Setting `json:"settings"`
	// +kubebuilder:validation:Optional
	// Notify members when they are added to the application
	NotifyMembers bool `json:"notifyMembers"`
	// Application metadata
	Metadata *[]Metadata `json:"metadata,omitempty"`
	// Application members
	Members *[]Member `json:"members,omitempty"`
}

// GetSettings implements core.ApplicationModel.
func (in *Application) GetSettings() core.ApplicationSettings {
	return in.Settings
}
