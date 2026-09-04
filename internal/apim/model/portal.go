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
	nav "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portal"
)

type PortalDTO struct {
	HRID       string                  `json:"hrid,omitempty"`
	Name       string                  `json:"name"`
	Navigation []NavigationPathDTO     `json:"navigation,omitempty" drift:"ignore"`
	Structure  *NavigationStructureDTO `json:"structure,omitempty"`
	// An absent value is read by APIM as a request to deactivate the portal's theme.
	ActiveThemeHRID string `json:"activeThemeHrid,omitempty"`
}

type NavigationPathDTO struct {
	Path        string         `json:"path"`
	DisplayName *string        `json:"displayName,omitempty"`
	Order       *int32         `json:"order,omitempty"`
	Visibility  nav.Visibility `json:"visibility,omitempty" drift:"ignore-unset"`
}

type NavigationStructureDTO struct {
	TopNavbar []*NavigationEntryDTO `json:"topNavbar,omitempty"`
}

type NavigationEntryDTO struct {
	Path        string         `json:"path"`
	DisplayName *string        `json:"displayName,omitempty"`
	Visibility  nav.Visibility `json:"visibility,omitempty" drift:"ignore-unset"`
}

type PortalState struct {
	PortalDTO     `json:",omitempty"`
	portal.Status `json:",omitempty"`
}

func ToPortalDTO(crd portal.Type, hrid string, activeThemeHrid string) PortalDTO {
	dto := mapViaJSON[PortalDTO](crd)
	dto.HRID = hrid
	dto.ActiveThemeHRID = activeThemeHrid
	return dto
}
