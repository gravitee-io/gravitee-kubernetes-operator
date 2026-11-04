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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ContextAwareObject = &AuthzPdp{}
var _ core.Spec = &AuthzPdpSpec{}
var _ core.Status = &AuthzPdpStatus{}
var _ core.ConditionAware = &AuthzPdp{}

// AuthzPdpSpec defines the desired state of an AuthzPdp.
// +kubebuilder:object:generate=true
// +kubebuilder:validation:XValidation:rule="self.targetPdpId == oldSelf.targetPdpId",messageExpression="'targetPdpId is immutable: [' + oldSelf.targetPdpId + '] cannot become [' + self.targetPdpId + ']; delete and recreate the AuthzPdp'"
// +kubebuilder:validation:XValidation:rule="(has(self.tag) ? self.tag : '') == (has(oldSelf.tag) ? oldSelf.tag : '')",messageExpression="'tag is immutable: [' + (has(oldSelf.tag) ? oldSelf.tag : '') + '] cannot become [' + (has(self.tag) ? self.tag : '') + ']; delete and recreate the AuthzPdp'"
type AuthzPdpSpec struct {
	authzpdp.Type `json:",inline"`
	// Reference to a ManagementContext that determines which APIM instance this PDP gateway is
	// created in. Immutable: the operator deletes a PDP through its current context only, so a PDP
	// moved to another one would stay on the first platform with its routing scope and endpoint.
	// The message is static: refs.NamespacedName strings carry no maxLength, so a message naming
	// both values exceeds the API server's CEL cost budget and the CRD is refused.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="contextRef is immutable; delete and recreate the AuthzPdp"
	Context *refs.NamespacedName `json:"contextRef,omitempty"`
}

func (spec *AuthzPdpSpec) Hash() string {
	return hash.Calculate(spec)
}

// AuthzPdpStatus defines the observed state of an AuthzPdp.
type AuthzPdpStatus struct {
	authzpdp.Status `json:",inline"`
}

func (s *AuthzPdpStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *AuthzPdp:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *AuthzPdpStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *AuthzPdp:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *AuthzPdpStatus) IsFailed() bool {
	if s.Conditions != nil {
		for _, condition := range s.Conditions {
			if condition.Status == metav1.ConditionFalse {
				return true
			}
		}
	}
	return false
}

func (s *AuthzPdpStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unused
}

// AuthzPdp is a PDP gateway of Gravitee Authorization Management as a Kubernetes resource: the
// engine it runs (targetPdpId), the gateways it runs on (tag) and its AuthZEN endpoint. Its
// status.routingScope is what the pdpRef of a schema, an entity or a policy resolves to.
//
// Deleting an AuthzPdp removes its routing scope from the targets of every schema, entity and
// policy that named it, and a resource whose pdpRef still names it waits on ResolvedRefs=False.
// The platform refuses a delete that would leave a policy or an entity with no target: the
// AuthzPdp then stays in Terminating, its Accepted condition naming the reason, until those are
// retargeted.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Routing Scope",type=string,JSONPath=`.status.routingScope`
// +kubebuilder:printcolumn:name="Endpoint",type=string,JSONPath=`.spec.endpoint.path`
// +kubebuilder:printcolumn:name="Strict",type=boolean,JSONPath=`.spec.strict`
// +kubebuilder:resource:shortName=graviteeauthzpdps
// +kubebuilder:storageversion
type AuthzPdp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AuthzPdpSpec   `json:"spec,omitempty"`
	Status            AuthzPdpStatus `json:"status,omitempty"`
}

// AuthzPdpList contains a list of AuthzPdp resources.
// +kubebuilder:object:root=true
type AuthzPdpList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AuthzPdp `json:"items"`
}

func (t *AuthzPdp) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      t.Name,
		Namespace: t.Namespace,
	}
}

func (t *AuthzPdp) GetSpec() core.Spec {
	return &t.Spec
}

func (t *AuthzPdp) GetStatus() core.Status {
	return &t.Status
}

func (t *AuthzPdp) IsBeingDeleted() bool {
	return !t.ObjectMeta.DeletionTimestamp.IsZero()
}

func (t *AuthzPdp) HasContext() bool {
	return t.Spec.Context != nil
}

func (t *AuthzPdp) ContextRef() core.ObjectRef {
	return t.Spec.Context
}

func (t *AuthzPdp) GetEnvID() string {
	return t.Status.EnvID
}

func (t *AuthzPdp) GetID() string {
	return t.Status.ID
}

func (t *AuthzPdp) GetOrgID() string {
	return t.Status.OrgID
}

func (t *AuthzPdp) PopulateIDs(_ core.ContextModel, _ bool) {
	// done when calling the API
}

func (t *AuthzPdp) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(t.Status.Conditions)
}

func (t *AuthzPdp) SetConditions(conditions []metav1.Condition) {
	t.Status.Conditions = conditions
}
