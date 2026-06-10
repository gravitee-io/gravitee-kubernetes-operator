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

// Type defines the specification of a Portal resource (next-gen developer portal).
type Type struct {
	// Display name of the portal.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// The portal's navigation hierarchy as an ordered, flat list of paths.
	// The order of entries in the list is preserved. Intermediate folders are
	// implicitly created by APIM if not listed explicitly.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=50
	// +listType=map
	// +listMapKey=path
	Navigation []*nav.NavigationPath `json:"navigation,omitempty"`
}
