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
	AppType  string `json:"type,omitempty"`
	ClientId string `json:"client_id,omitempty"`
}

type OAuthClientSettings struct {
	ApplicationType            string   `json:"application_type,omitempty"`
	ClientId                   string   `json:"client_id,omitempty"`
	ClientSecret               string   `json:"client_secret,omitempty"`
	ClientURI                  string   `json:"client_uri,omitempty"`
	GrantTypes                 []string `json:"grant_types,omitempty"`
	LogoURI                    string   `json:"logo_uri,omitempty"`
	RedirectUris               []string `json:"redirect_uris,omitempty"`
	RenewClientSecretSupported bool     `json:"renew_client_secret_supported,omitempty"`
	ResponseTypes              []string `json:"response_types,omitempty"`
}

type Setting struct {
	App   *SimpleSettings      `json:"app,omitempty"`
	Oauth *OAuthClientSettings `json:"oauth,omitempty"`
}

// +kubebuilder:validation:Enum=STRING;NUMERIC;BOOLEAN;DATE;MAIL;URL;
type MetaDataFormat string

type MetaData struct {
	// +kubebuilder:validation:Required
	Name          string          `json:"name"`
	ApplicationId string          `json:"applicationId,omitempty"`
	Value         string          `json:"value,omitempty"`
	DefaultValue  string          `json:"defaultValue,omitempty"`
	Format        *MetaDataFormat `json:"format,omitempty"`
	Hidden        bool            `json:"hidden,omitempty"`
	Key           string          `json:"key,omitempty"`
}

type Application struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	Description     string                  `json:"description,omitempty"`
	ApplicationType string                  `json:"type,omitempty"`
	ClientId        string                  `json:"clientId,omitempty"`
	RedirectUris    []string                `json:"redirectUris,omitempty"`
	Metadata        *utils.GenericStringMap `json:"metadata,omitempty"`

	// The origin which is used to create this Application
	// +kubebuilder:validation:Enum=kubernetes;
	Origin string `json:"origin,omitempty"`
	// io.gravitee.definition.model.Application
	ID                             string      `json:"id,omitempty"`
	Background                     string      `json:"background,omitempty"`
	Domain                         string      `json:"domain,omitempty"`
	Groups                         []string    `json:"groups,omitempty"`
	Picture                        string      `json:"picture,omitempty"`
	PictureURL                     string      `json:"picture_url,omitempty"`
	Settings                       *Setting    `json:"settings,omitempty"`
	AppKeyMode                     *KeyMode    `json:"app_key_mode,omitempty"`
	DisableMembershipNotifications bool        `json:"disable_membership_notifications,omitempty"`
	ApplicationMetaData            *[]MetaData `json:"applicationMetaData,omitempty"`
}
