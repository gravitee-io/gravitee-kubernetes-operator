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
	documentation "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/docs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

// DocumentationDTO is the automation API wire representation of a Documentation
// page. The owning Portal / API is carried by the request path, not the body,
// so the Kubernetes portalRef / apiRef fields are intentionally excluded.
type DocumentationDTO struct {
	HRID     string                 `json:"hrid,omitempty"`
	Name     string                 `json:"name,omitempty"`
	PageType documentation.PageType `json:"type,omitempty"`
	Content  string                 `json:"content,omitempty"`
	Location string                 `json:"location,omitempty"`
	Order    *int32                 `json:"order,omitempty"`
	// Unset is omitted so that APIM applies its own TOP_NAVBAR default.
	Area documentation.PageArea `json:"area,omitempty" drift:"ignore-remote:TOP_NAVBAR"`
}

type DocumentationState struct {
	DocumentationDTO     `json:",omitempty"`
	documentation.Status `json:",omitempty"`
}

func ToDocumentationDTO(crd *v1alpha1.Documentation) DocumentationDTO {
	dto := mapViaJSON[DocumentationDTO](crd.Spec.Type)
	dto.HRID = refs.NewNamespacedNameFromObject(crd).HRID()
	return dto
}
