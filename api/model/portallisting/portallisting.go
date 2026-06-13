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

package portallisting

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

// Type defines the specification of a PortalListing resource.
// It publishes one or more APIs to a portal at chosen locations in the
// portal's navigation hierarchy.
type Type struct {
	// Reference to the Portal this listing publishes APIs to.
	// +kubebuilder:validation:Required
	Portal refs.NamespacedName `json:"portalRef"`
	// The APIs to publish to the portal, each at a chosen location in the
	// portal's navigation. The order of entries in the list is preserved and
	// also determines display order relative to siblings sharing the same
	// location when no explicit order is set.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	APIs []ApiEntry `json:"apis"`
}

// ApiEntry places one API at a location in the portal's navigation.
type ApiEntry struct {
	// Reference to the API to publish. Only v4 APIs (ApiV4Definition) are
	// supported by the next-gen portal; the referenced kind defaults to
	// ApiV4Definition when left empty.
	// +kubebuilder:validation:Required
	Ref refs.NamespacedName `json:"ref"`
	// The path in the portal's navigation where this API should appear.
	// The API is only visible on the portal if this matches a path defined
	// in the Portal's navigation.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^/`
	Location string `json:"location"`
	// Optional display order of this API relative to its siblings at the same
	// location. The position in the list is also preserved.
	// +kubebuilder:validation:Optional
	Order *int32 `json:"order,omitempty"`
}

func (t *Type) GetPortalRef() core.ObjectRef {
	return &t.Portal
}

func (t *Type) GetApiRefs() []core.ObjectRef {
	out := make([]core.ObjectRef, 0, len(t.APIs))
	for i := range t.APIs {
		out = append(out, &t.APIs[i].Ref)
	}
	return out
}
