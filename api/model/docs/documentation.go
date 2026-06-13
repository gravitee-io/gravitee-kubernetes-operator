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

package docs

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

// PageType is the type of a Documentation page.
// +kubebuilder:validation:Enum=GRAVITEE_MARKDOWN;OPENAPI;ASYNCAPI
type PageType string

const (
	GraviteeMarkdown PageType = "GRAVITEE_MARKDOWN"
	OpenAPI          PageType = "OPENAPI"
	AsyncAPI         PageType = "ASYNCAPI"
)

// Type defines the specification of a Documentation resource.
// A Documentation page is attached to exactly one of a Portal (portalRef) or
// an API (apiRef) and is placed at a chosen location in the owning resource's
// navigation hierarchy.
type Type struct {
	// Display name of the documentation page.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// The type of the documentation page.
	// +kubebuilder:validation:Required
	PageType PageType `json:"type"`
	// The content of the documentation page.
	// +kubebuilder:validation:Required
	Content string `json:"content"`
	// The path in the owning resource's navigation hierarchy where this page
	// should appear. The page is only visible if this matches a path defined
	// in the Portal's or API's navigation.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^/`
	Location string `json:"location,omitempty"`
	// Optional display order of this page relative to its siblings at the same
	// location.
	// +kubebuilder:validation:Optional
	Order *int32 `json:"order,omitempty"`
	// Reference to the Portal this documentation page is attached to.
	// Mutually exclusive with apiRef; exactly one of the two must be set.
	// +kubebuilder:validation:Optional
	Portal *refs.NamespacedName `json:"portalRef,omitempty"`
	// Reference to the API this documentation page is attached to. Only v4 APIs
	// (ApiV4Definition) are supported by the next-gen portal; the referenced
	// kind defaults to ApiV4Definition when left empty.
	// Mutually exclusive with portalRef; exactly one of the two must be set.
	// +kubebuilder:validation:Optional
	API *refs.NamespacedName `json:"apiRef,omitempty"`
}

func (t *Type) IsPortalDoc() bool {
	return t.Portal != nil
}

func (t *Type) IsApiDoc() bool {
	return t.API != nil
}

func (t *Type) GetPortalRef() core.ObjectRef {
	if t.Portal == nil {
		return nil
	}
	return t.Portal
}

func (t *Type) GetApiRef() core.ObjectRef {
	if t.API == nil {
		return nil
	}
	return t.API
}
