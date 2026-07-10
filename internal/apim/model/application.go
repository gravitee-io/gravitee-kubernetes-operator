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

package model

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

type ApplicationDTO struct {
	ID            string                   `json:"id,omitempty" drift:"ignore"`
	HRID          string                   `json:"hrid,omitempty" drift:"ignore"`
	Name          string                   `json:"name,omitempty"`
	Status        string                   `json:"status,omitempty" drift:"ignore"`
	Description   string                   `json:"description,omitempty" drift:"trimmed"`
	Settings      *ApplicationSettingsDTO  `json:"settings,omitempty" drift:"empty-is-nil"`
	Background    string                   `json:"background,omitempty"`
	Domain        string                   `json:"domain,omitempty"`
	Groups        []string                 `json:"groups,omitempty" drift:"empty-is-nil"`
	Picture       string                   `json:"picture,omitempty"`
	PictureURL    string                   `json:"pictureUrl,omitempty"`
	NotifyMembers *bool                    `json:"notifyMembers" drift:"empty-is-nil"`
	Metadata      []ApplicationMetadataDTO `json:"metadata" drift:"empty-is-nil"`
	Members       []ApplicationMemberDTO   `json:"members" drift:"empty-is-nil"`
}

type ApplicationSettingsDTO struct {
	App   *ApplicationSimpleSettingsDTO      `json:"app,omitempty"`
	Oauth *ApplicationOAuthClientSettingsDTO `json:"oauth,omitempty"`
	TLS   *ApplicationTLSSettingsDTO         `json:"tls,omitempty"`
}

type ApplicationSimpleSettingsDTO struct {
	AppType  string  `json:"type"`
	ClientID *string `json:"clientId,omitempty"`
}

type ApplicationOAuthClientSettingsDTO struct {
	ApplicationType application.OauthType   `json:"applicationType"`
	GrantTypes      []application.GrantType `json:"grantTypes"`
	RedirectUris    []string                `json:"redirectUris"`
}

type ApplicationTLSSettingsDTO struct {
	ClientCertificate  string                            `json:"clientCertificate,omitempty" drift:"trimmed"`
	ClientCertificates []ApplicationClientCertificateDTO `json:"clientCertificates,omitempty" drift:"empty-is-nil"`
}

type ApplicationClientCertificateDTO struct {
	Name     string                        `json:"name,omitempty"`
	Content  string                        `json:"content,omitempty" drift:"trimmed"`
	Ref      *ApplicationCertificateRefDTO `json:"ref,omitempty" drift:"ignore"`
	StartsAt string                        `json:"startsAt,omitempty" drift:"rfc3339"`
	EndsAt   string                        `json:"endsAt,omitempty" drift:"rfc3339"`
	Encoded  bool                          `json:"encoded,omitempty"`
}

type ApplicationCertificateRefDTO struct {
	Kind      string `json:"kind,omitempty"`
	Name      string `json:"name"`
	Key       string `json:"key,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type ApplicationMetadataDTO struct {
	Metadata `json:",inline"`
	Format   *application.MetaDataFormat `json:"format,omitempty"`
}

type ApplicationMemberDTO struct {
	Source   string `json:"source"`
	SourceID string `json:"sourceId"`
	Role     string `json:"role,omitempty"`
}

type ApplicationMetaData struct {
	Name          string                      `json:"name"`
	ApplicationID string                      `json:"applicationId,omitempty"`
	Value         string                      `json:"value,omitempty"`
	DefaultValue  string                      `json:"defaultValue,omitempty"`
	Format        *application.MetaDataFormat `json:"format,omitempty"`
	Hidden        bool                        `json:"hidden,omitempty"`
	Key           string                      `json:"key,omitempty"`
}

func ToApplicationDTO(spec v1alpha1.ApplicationSpec) ApplicationDTO {
	return mapViaJSON[ApplicationDTO](spec.Application)
}
