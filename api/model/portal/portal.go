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

package portal

import (
	nav "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/navigation"
)

// NavigationEntry is a portal navigation entry, ordered by its position in the list.
type NavigationEntry struct {
	// A slash-separated path defining the navigation hierarchy.
	// Intermediate folders are implicitly created by APIM if not listed explicitly.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^/`
	Path string `json:"path"`
	// Optional human-friendly label for this node. Listing a path explicitly
	// is the only way to attach a display name to it.
	// +kubebuilder:validation:Optional
	DisplayName *string `json:"displayName,omitempty"`
	// Whether this entry is shown to anonymous visitors of the portal.
	// Left unset, APIM resolves it: the entry inherits the visibility of its
	// parent folder, and a root entry defaults to PUBLIC. A PUBLIC entry under
	// a PRIVATE parent is rejected.
	// +kubebuilder:validation:Optional
	Visibility *nav.Visibility `json:"visibility,omitempty"`
}

// NavigationStructure groups a portal's navigation by portal area. Each area
// carries its own tree, and leaving an area unset leaves it untouched.
type NavigationStructure struct {
	// Top navbar entries as an ordered, flat list of paths.
	// The order of entries in the list is preserved. Intermediate folders are
	// implicitly created by APIM if not listed explicitly.
	// +kubebuilder:validation:Optional
	// +listType=map
	// +listMapKey=path
	TopNavbar []*NavigationEntry `json:"topNavbar,omitempty"`
}

// Type defines the specification of a Portal resource (next-gen developer portal).
type Type struct {
	// Display name of the portal.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// The portal's navigation, grouped by portal area.
	// +kubebuilder:validation:Optional
	Structure *NavigationStructure `json:"structure,omitempty"`
	// Deprecated: use structure.topNavbar instead, the two cannot be set at the same time.
	// Still synced for portals created before navigation was grouped by area.
	// +kubebuilder:validation:Optional
	// +listType=map
	// +listMapKey=path
	Navigation []*nav.NavigationPath `json:"navigation,omitempty"`
}

func (t *Type) UsesDeprecatedNavigation() bool {
	return len(t.Navigation) > 0
}
