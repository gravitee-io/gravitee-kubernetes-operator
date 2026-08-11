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

package portallink

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

// Type defines the specification of a PortalLink resource.
// It attaches an external navigation link to a portal at a chosen location
// in the portal's navigation hierarchy.
type Type struct {
	// Reference to the Portal this link is attached to.
	// +kubebuilder:validation:Required
	Portal refs.NamespacedName `json:"portalRef"`
	// Display name of the link.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// The URL this link points to.
	// +kubebuilder:validation:Required
	Href string `json:"href"`
	// The path in the portal's navigation hierarchy where this link should
	// appear. The link is only visible on the portal if this matches a path
	// defined in the Portal's navigation.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^/`
	Location *string `json:"location,omitempty"`
	// Optional display order of this link relative to its siblings at the
	// same location.
	// +kubebuilder:validation:Optional
	Order *int32 `json:"order,omitempty"`
}

func (t *Type) GetPortalRef() core.ObjectRef {
	return &t.Portal
}
