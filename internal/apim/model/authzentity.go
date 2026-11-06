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

package model

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

// AuthzEntityDTO is the automation API wire representation of an AuthzEntity (PUT /authz/entities).
// The platform sorts targets, parents and attribute keys and answers {} and [] for absent ones, so
// none of them compares by position.
type AuthzEntityDTO struct {
	HRID       string                  `json:"hrid,omitempty" drift:"ignore"`
	EntityID   string                  `json:"entityId"`
	Kind       string                  `json:"kind"`
	EntityType string                  `json:"entityType"`
	Attributes *utils.GenericStringMap `json:"attributes,omitempty" drift:"unstructured"`
	Parents    []AuthzEntityParentDTO  `json:"parents,omitempty" drift:"set"`
	// Active is always sent, so that false is never dropped on the wire or in drift.
	Active  bool     `json:"active"`
	Targets []string `json:"targets,omitempty" drift:"set"`
}

// AuthzEntityParentDTO names a parent the way the engine spells it.
type AuthzEntityParentDTO struct {
	EntityType string `json:"entityType"`
	EntityID   string `json:"entityId"`
}

// AuthzEntityState is what a PUT or a GET answers: the payload and the envelope.
type AuthzEntityState struct {
	AuthzEntityDTO `json:",inline"`
	ID             string        `json:"id,omitempty"`
	OrgID          string        `json:"organizationId,omitempty"`
	EnvID          string        `json:"environmentId,omitempty"`
	Errors         status.Errors `json:"errors,omitempty"`
}

// ToAuthzEntityParentDTO maps a declared parent onto its wire shape.
func ToAuthzEntityParentDTO(parent authzentity.Parent) AuthzEntityParentDTO {
	return AuthzEntityParentDTO{EntityType: parent.EntityType, EntityID: parent.EntityID}
}

// ToAuthzEntityDTO maps the CRD onto the wire payload with the targets the caller resolved, in
// declaration order. It is the single mapping used by the sync path, the dry run and drift
// detection. The targets are a parameter rather than read from the spec: a pdpRef has no wire
// value until the cluster is read.
func ToAuthzEntityDTO(crd *v1alpha1.AuthzEntity, targets []string) AuthzEntityDTO {
	var parents []AuthzEntityParentDTO
	for _, parent := range crd.Spec.Parents {
		parents = append(parents, ToAuthzEntityParentDTO(parent))
	}
	return AuthzEntityDTO{
		HRID:       refs.NewNamespacedNameFromObject(crd).HRID(),
		EntityID:   crd.Spec.EntityID,
		Kind:       crd.Spec.Kind,
		EntityType: crd.Spec.EntityType,
		Attributes: crd.Spec.Attributes.DeepCopy(),
		Parents:    parents,
		Active:     crd.Spec.IsActive(),
		Targets:    targets,
	}
}
