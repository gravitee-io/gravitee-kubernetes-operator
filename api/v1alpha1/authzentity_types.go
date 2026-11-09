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

package v1alpha1

import (
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ContextAwareObject = &AuthzEntity{}
var _ core.Spec = &AuthzEntitySpec{}
var _ core.Status = &AuthzEntityStatus{}
var _ core.ConditionAware = &AuthzEntity{}

// AuthzEntitySpec defines the desired state of an AuthzEntity.
// +kubebuilder:object:generate=true
// +kubebuilder:validation:XValidation:rule="self.entityId == oldSelf.entityId",messageExpression="'entityId is immutable: [' + oldSelf.entityId + '] cannot become [' + self.entityId + ']; delete and recreate the AuthzEntity'"
// +kubebuilder:validation:XValidation:rule="self.kind == oldSelf.kind",messageExpression="'kind is immutable: [' + oldSelf.kind + '] cannot become [' + self.kind + ']; delete and recreate the AuthzEntity'"
// +kubebuilder:validation:XValidation:rule="self.entityType == oldSelf.entityType",messageExpression="'entityType is immutable: [' + oldSelf.entityType + '] cannot become [' + self.entityType + ']; delete and recreate the AuthzEntity'"
type AuthzEntitySpec struct {
	authzentity.Type `json:",inline"`
	// Reference to a ManagementContext that determines which APIM instance this entity is created
	// in. Immutable: the operator deletes an entity through its current context only, so an entity
	// moved to another one would stay on the first platform, published to its targets.
	// The message is static: refs.NamespacedName strings carry no maxLength, so a message naming
	// both values exceeds the API server's CEL cost budget and the CRD is refused. Required, so the
	// rule also holds without the admission webhook: CEL compares the field only when both sides have it.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="contextRef is immutable; delete and recreate the AuthzEntity"
	Context *refs.NamespacedName `json:"contextRef"`
}

func (spec *AuthzEntitySpec) Hash() string {
	return hash.Calculate(spec)
}

// AuthzEntityStatus defines the observed state of an AuthzEntity.
type AuthzEntityStatus struct {
	authzentity.Status `json:",inline"`
}

func (s *AuthzEntityStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *AuthzEntity:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *AuthzEntityStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *AuthzEntity:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *AuthzEntityStatus) IsFailed() bool {
	if s.Conditions != nil {
		for _, condition := range s.Conditions {
			if condition.Status == metav1.ConditionFalse {
				return true
			}
		}
	}
	return false
}

func (s *AuthzEntityStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unused
}

// AuthzEntity is a principal or a resource of Gravitee Authorization Management as a Kubernetes
// resource, published to the PDP gateways its targets name. One entity of the environment holds
// an entityId: an id another entity holds is refused, naming the source that declared it.
//
// A target names an AuthzPdp through pdpRef, or a routing scope no AuthzPdp declares through
// scope. A pdpRef waits on ResolvedRefs=False until the AuthzPdp has a routing scope.
//
// Deleting an AuthzEntity deletes its RESOURCE children and detaches its PRINCIPAL children on the
// platform, and each AuthzEntity declaring one of them is applied again: a PRINCIPAL as declared, a
// RESOURCE refused, its Accepted condition naming the parent that no longer exists. A delete that
// would remove more entities than the platform's cascade limit is refused (413): the resource stays
// Terminating, its Accepted condition naming the reason.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Entity ID",type=string,JSONPath=`.spec.entityId`
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.entityType`
// +kubebuilder:resource:shortName=graviteeauthzentities
// +kubebuilder:storageversion
type AuthzEntity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AuthzEntitySpec   `json:"spec,omitempty"`
	Status            AuthzEntityStatus `json:"status,omitempty"`
}

// AuthzEntityList contains a list of AuthzEntity resources.
// +kubebuilder:object:root=true
type AuthzEntityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AuthzEntity `json:"items"`
}

func (t *AuthzEntity) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      t.Name,
		Namespace: t.Namespace,
	}
}

func (t *AuthzEntity) GetSpec() core.Spec {
	return &t.Spec
}

func (t *AuthzEntity) GetStatus() core.Status {
	return &t.Status
}

func (t *AuthzEntity) IsBeingDeleted() bool {
	return !t.ObjectMeta.DeletionTimestamp.IsZero()
}

func (t *AuthzEntity) HasContext() bool {
	return t.Spec.Context != nil
}

func (t *AuthzEntity) ContextRef() core.ObjectRef {
	return t.Spec.Context
}

func (t *AuthzEntity) GetEnvID() string {
	return t.Status.EnvID
}

func (t *AuthzEntity) GetID() string {
	return t.Status.ID
}

func (t *AuthzEntity) GetOrgID() string {
	return t.Status.OrgID
}

func (t *AuthzEntity) PopulateIDs(_ core.ContextModel, _ bool) {
	// done when calling the API
}

func (t *AuthzEntity) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(t.Status.Conditions)
}

func (t *AuthzEntity) SetConditions(conditions []metav1.Condition) {
	t.Status.Conditions = conditions
}
