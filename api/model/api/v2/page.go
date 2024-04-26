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

package v2

import "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"

type PageSource struct {
	// +kubebuilder:validation:Required
	Type string `json:"type"`
	// +kubebuilder:validation:Required
	Configuration *utils.GenericStringMap `json:"configuration"`
}

type Page struct {
	// +kubebuilder:validation:Optional
	// The ID of the page. This field is mostly required when you are applying
	// an API exported from APIM to make the operator take control over it.
	// If not set, this ID will be generated in a predictable manner based on
	// the map key associated to this entry in the API.
	ID string `json:"id,omitempty"`
	// +kubebuilder:validation:Optional
	// CrossID is designed to identified a page across environments.
	// If not set, this ID will be generated in a predictable manner based on
	// the map key associated to this entry in the API.
	CrossID string `json:"crossId,omitempty"`
	// +kubebuilder:validation:Required
	// This is the display name of the page in APIM and on the portal.
	// This field can be edited safely if you want to rename a page.
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=MARKDOWN;SWAGGER;ASYNCAPI;ASCIIDOC;FOLDER;SYSTEM_FOLDER
	// The type of the documentation page or folder.
	Type string `json:"type,omitempty"`
	// +kubebuilder:validation:Optional
	// The content of the page, if any.
	Content string `json:"content,omitempty"`
	// +kubebuilder:validation:Optional
	// The order used to display the page in APIM and on the portal.
	Order uint64 `json:"order"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	// If true, the page will be accessible from the portal (default is false)
	Published bool `json:"published"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=PUBLIC
	// +kubebuilder:validation:Enum=PUBLIC;
	// The visibility of the page. Only public pages are supported at the moment.
	Visibility string `json:"visibility,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	// If true, this page will be displayed as the homepage of your API documentation.
	HomePage bool `json:"homepage"`
	// +kubebuilder:validation:Optional
	// If your page contains a folder, setting this field to the map key associated to the
	// folder entry will be reflected into APIM by making the page a child of this folder.
	Parent string `json:"parent,omitempty"`
	// +kubebuilder:validation:Optional
	// The parent ID of the page. This field is mostly required when you are applying
	// an API exported from APIM to make the operator take control over it. Use `Parent`
	// in any other case.
	ParentID string `json:"parentId,omitempty"`
	// +kubebuilder:validation:Optional
	// The API of the page. If empty, will be set automatically to the generated ID of the API.
	API string `json:"api,omitempty"`
	// +kubebuilder:validation:Optional
	// Source allow you to fetch pages from various external sources, overriding page content
	// each time the source is fetched.
	Source *PageSource `json:"source,omitempty"`
	// +kubebuilder:validation:Optional
	// Legacy page configuration support can be added using this field
	Configuration map[string]string `json:"configuration,omitempty"`
}
