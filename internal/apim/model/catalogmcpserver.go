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
package model

import (
	"encoding/json"
	goerrors "errors"
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/catalogmcpserver"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

// CatalogMcpServerDTO is the automation API wire representation of a CatalogMcpServer
// (PUT /aim/catalog/mcp-servers). The auth variants are flat on the wire, discriminated by type.
type CatalogMcpServerDTO struct {
	HRID        string                        `json:"hrid,omitempty" drift:"ignore"`
	EntityID    string                        `json:"entityId"`
	Description *string                       `json:"description,omitempty" drift:"empty-is-nil"`
	Connection  CatalogMcpServerConnectionDTO `json:"connection"`
}

// CatalogMcpServerConnectionDTO is where the upstream answers and how the platform authenticates.
type CatalogMcpServerConnectionDTO struct {
	Endpoint  string `json:"endpoint"`
	Transport string `json:"transport"`
	// Auth is always emitted by the mapper ({type: NONE} when the spec has none), because the
	// platform reads an absent auth back as {type: NONE} and drift would otherwise flag it.
	Auth *CatalogMcpServerAuthDTO `json:"auth,omitempty"`
}

// CatalogMcpServerAuthDTO is the wire shape of the upstream authentication. Credential values
// are never returned by the platform, so they are ignored for drift.
type CatalogMcpServerAuthDTO struct {
	Type         string  `json:"type"`
	Name         *string `json:"name,omitempty"`
	Value        *string `json:"value,omitempty" drift:"ignore"`
	ClientID     *string `json:"clientId,omitempty"`
	ClientSecret *string `json:"clientSecret,omitempty" drift:"ignore"`
	TokenURL     *string `json:"tokenUrl,omitempty"`
	Scope        *string `json:"scope,omitempty" drift:"empty-is-nil"`
}

// CatalogMcpServerState is what a PUT or a GET answers: the payload plus the platform's state.
// It deliberately does not embed catalogmcpserver.Status: both halves would promote a hrid
// field and encoding/json silently drops a name promoted twice at the same depth.
type CatalogMcpServerState struct {
	CatalogMcpServerDTO         `json:",inline"`
	ID                          string        `json:"id,omitempty"`
	OrgID                       string        `json:"organizationId,omitempty"`
	EnvID                       string        `json:"environmentId,omitempty"`
	Errors                      status.Errors `json:"errors,omitempty"`
	catalogmcpserver.Discovered `json:",inline"`
}

// CatalogMcpServerRefusal reads the findings of a refused apply. The platform refuses with a 400
// whose body is the resource state, credentials removed, with the findings in errors.severe; a
// 400 the host answers before the module sees the request has another shape and is not a refusal.
func CatalogMcpServerRefusal(err error) (status.Errors, bool) {
	serverError := &errors.ServerError{}
	if !goerrors.As(err, serverError) || serverError.StatusCode != http.StatusBadRequest {
		return status.Errors{}, false
	}

	state := new(CatalogMcpServerState)
	if json.Unmarshal([]byte(serverError.Body), state) != nil || len(state.Errors.Severe) == 0 {
		return status.Errors{}, false
	}

	return state.Errors, true
}

// ToCatalogMcpServerDTO maps the CRD onto the wire payload. It is the single mapping used by
// the sync path and by drift detection.
func ToCatalogMcpServerDTO(crd *v1alpha1.CatalogMcpServer) CatalogMcpServerDTO {
	spec := crd.Spec.Type

	transport := string(spec.Connection.Transport)
	if transport == "" {
		transport = string(catalogmcpserver.TransportHTTP)
	}

	return CatalogMcpServerDTO{
		HRID:        refs.NewNamespacedNameFromObject(crd).HRID(),
		EntityID:    spec.EntityID,
		Description: spec.Description,
		Connection: CatalogMcpServerConnectionDTO{
			Endpoint:  spec.Connection.Endpoint,
			Transport: transport,
			Auth:      toCatalogMcpServerAuthDTO(spec.Connection.Auth),
		},
	}
}

func toCatalogMcpServerAuthDTO(auth *catalogmcpserver.Auth) *CatalogMcpServerAuthDTO {
	if auth == nil || auth.Type == "" {
		return &CatalogMcpServerAuthDTO{Type: string(catalogmcpserver.AuthTypeNone)}
	}

	dto := &CatalogMcpServerAuthDTO{Type: string(auth.Type)}

	switch auth.Type {
	case catalogmcpserver.AuthTypeHeader:
		if auth.Header != nil {
			dto.Name = new(auth.Header.Name)
			dto.Value = new(auth.Header.Value)
		}
	case catalogmcpserver.AuthTypeOAuth2:
		if auth.OAuth2 != nil {
			dto.ClientID = new(auth.OAuth2.ClientID)
			dto.ClientSecret = new(auth.OAuth2.ClientSecret)
			dto.TokenURL = new(auth.OAuth2.TokenURL)
			if auth.OAuth2.Scope != nil {
				dto.Scope = new(*auth.OAuth2.Scope)
			}
		}
	case catalogmcpserver.AuthTypeNone:
	}

	return dto
}
