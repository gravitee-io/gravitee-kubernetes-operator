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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portallink"
)

type PortalLinkDTO struct {
	HRID     string `json:"hrid,omitempty"`
	Name     string `json:"name,omitempty"`
	Href     string `json:"href,omitempty"`
	Location string `json:"location,omitempty"`
	Order    *int32 `json:"order,omitempty"`
	// Unset is omitted so that APIM resolves the visibility from the parent folder.
	Visibility nav.Visibility `json:"visibility,omitempty" drift:"ignore-unset"`
}

type PortalLinkState struct {
	PortalLinkDTO     `json:",omitempty"`
	portallink.Status `json:",omitempty"`
}

// ToPortalLinkDTO converts a PortalLink CRD to a PortalLinkDTO.
func ToPortalLinkDTO(crd portallink.Type, hrid string) PortalLinkDTO {
	dto := mapViaJSON[PortalLinkDTO](crd)
	dto.HRID = hrid
	return dto
}
