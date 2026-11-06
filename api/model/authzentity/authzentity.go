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

package authzentity

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

// Parent is an entity this one is a member of, spelled the way the engine spells it.
type Parent struct {
	// Type of the parent, spelled the way the engine spells it: Group, not group.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[_a-zA-Z][_a-zA-Z0-9]*(?:::[_a-zA-Z][_a-zA-Z0-9]*)*$`
	// +kubebuilder:validation:MaxLength=255
	EntityType string `json:"entityType"`
	// Id of the parent, under the same rule as entityId.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	EntityID string `json:"entityId"`
}

// Type defines the specification of an AuthzEntity resource: a principal or a resource of
// Authorization Management, its parents, its attributes and the PDP gateways it is published to.
type Type struct {
	// Id of the entity in policies, without its type: agents, not Group::"agents". One entity of
	// the environment holds it: an id another entity holds is refused, naming that entity's source.
	// Immutable.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	EntityID string `json:"entityId"`
	// Whether the entity acts (PRINCIPAL) or is acted upon (RESOURCE). Immutable.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=PRINCIPAL;RESOURCE
	Kind string `json:"kind"`
	// Type of the entity, spelled the way the engine spells it: Group, not group. A built-in type
	// sets the kind. Immutable.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[_a-zA-Z][_a-zA-Z0-9]*(?:::[_a-zA-Z][_a-zA-Z0-9]*)*$`
	// +kubebuilder:validation:MaxLength=255
	EntityType string `json:"entityType"`
	// Attributes of the entity, replaced whole on every apply. A key starting with _ is reserved
	// to the platform and refused.
	// +kubebuilder:validation:Optional
	Attributes *utils.GenericStringMap `json:"attributes,omitempty"`
	// The entities this one is a member of, replaced whole on every apply. Each is declared once,
	// and never in a cycle.
	// +kubebuilder:validation:Optional
	Parents []Parent `json:"parents,omitempty"`
	// Whether the entity is published to the PDP gateways it targets.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	Active *bool `json:"active,omitempty"`
	// The PDP gateways the entity is published to. Empty, it is published to the default engine on
	// every gateway. An Action entity takes none. Deleting the AuthzPdp a pdpRef names removes its
	// scope here; the platform refuses that delete when the entity would be left with no target.
	// +kubebuilder:validation:Optional
	Targets []authzpdp.Target `json:"targets,omitempty"`
}

// IsActive tells if the entity is published to its targets: an unset active reads as true, as on
// the platform.
func (t *Type) IsActive() bool {
	return t.Active == nil || *t.Active
}
