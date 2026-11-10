// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authzpdp

import "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"

// Target names the PDP gateway a resource is published to: an AuthzPdp through pdpRef, or a
// routing scope no AuthzPdp declares (default, a gateway created in the Console).
// +kubebuilder:validation:XValidation:rule="has(self.pdpRef) != has(self.scope)",message="set pdpRef or scope, not both"
type Target struct {
	// The AuthzPdp whose status.routingScope is the target, its namespace defaulted to the
	// resource's. It must sync through the same ManagementContext. The resource waits on
	// ResolvedRefs=False until the AuthzPdp has a routing scope.
	// +kubebuilder:validation:Optional
	PdpRef *refs.NamespacedName `json:"pdpRef,omitempty"`
	// A routing scope written as the gateway reports it: targetPdpId@tag, targetPdpId when
	// untagged, default, or * where the kind allows it. Not resolved: a bare targetPdpId naming a
	// tagged gateway is refused by the platform.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^(\*|[A-Za-z0-9_.-]+(@-?[a-z0-9]+(-[a-z0-9]+)*)?)$`
	// +kubebuilder:validation:MaxLength=320
	Scope string `json:"scope,omitempty"`
}
