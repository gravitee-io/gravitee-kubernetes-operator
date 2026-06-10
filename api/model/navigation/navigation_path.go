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

package navigation

// NavigationPath is a single entry in a portal or API definition navigation tree.
type NavigationPath struct {
	// A slash-separated path defining the navigation hierarchy.
	// Intermediate folders are implicitly created by APIM if not listed explicitly.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^/`
	Path string `json:"path"`
	// Optional human-friendly label for this node. Listing a path explicitly
	// is the only way to attach a display name to it.
	DisplayName *string `json:"displayName,omitempty"`
	// Optional display order of this node relative to its siblings at the same level.
	// Listing a path explicitly is the only way to attach an order.
	Order *int32 `json:"order,omitempty"`
}
