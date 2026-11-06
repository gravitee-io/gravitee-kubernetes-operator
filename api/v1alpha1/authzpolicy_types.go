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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpolicy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ContextAwareObject = &AuthzPolicy{}
var _ core.Spec = &AuthzPolicySpec{}
var _ core.Status = &AuthzPolicyStatus{}
var _ core.ConditionAware = &AuthzPolicy{}

// AuthzPolicySpec defines the desired state of an AuthzPolicy.
// +kubebuilder:object:generate=true
// +kubebuilder:validation:XValidation:rule="self.entityId == oldSelf.entityId",messageExpression="'entityId is immutable: [' + oldSelf.entityId + '] cannot become [' + self.entityId + ']; delete and recreate the AuthzPolicy'"
// +kubebuilder:validation:XValidation:rule="self.policyType == oldSelf.policyType",messageExpression="'policyType is immutable: [' + oldSelf.policyType + '] cannot become [' + self.policyType + ']; delete and recreate the AuthzPolicy'"
type AuthzPolicySpec struct {
	authzpolicy.Type `json:",inline"`
	// Reference to a ManagementContext that determines which APIM instance this policy is created
	// in. Immutable: the operator deletes a policy through its current context only, so a policy
	// moved to another one would stay on the first platform, deployed to its targets.
	// The message is static: refs.NamespacedName strings carry no maxLength, so a message naming
	// both values exceeds the API server's CEL cost budget and the CRD is refused. Required, so the
	// rule also holds without the admission webhook: CEL compares the field only when both sides have it.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="contextRef is immutable; delete and recreate the AuthzPolicy"
	Context *refs.NamespacedName `json:"contextRef"`
}

func (spec *AuthzPolicySpec) Hash() string {
	return hash.Calculate(spec)
}

// AuthzPolicyStatus defines the observed state of an AuthzPolicy.
type AuthzPolicyStatus struct {
	authzpolicy.Status `json:",inline"`
}

func (s *AuthzPolicyStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *AuthzPolicy:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *AuthzPolicyStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *AuthzPolicy:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *AuthzPolicyStatus) IsFailed() bool {
	if s.Conditions != nil {
		for _, condition := range s.Conditions {
			if condition.Status == metav1.ConditionFalse {
				return true
			}
		}
	}
	return false
}

func (s *AuthzPolicyStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unused
}

// AuthzPolicy is a GAPL policy of Gravitee Authorization Management as a Kubernetes resource,
// deployed to the PDP gateways its targets name. One policy of the environment holds an entityId:
// an id another policy holds is refused, naming that policy.
//
// A target names an AuthzPdp through pdpRef, or a routing scope no AuthzPdp declares through
// scope. A pdpRef waits on ResolvedRefs=False until the AuthzPdp has a routing scope.
//
// Every apply deploys the policy: a policy disabled in the Console is deployed again by the next
// spec change, with a warning that says so. The ManagementContext's token needs the
// ENVIRONMENT_AUTHZ_PDP UPDATE permission to apply a policy and to delete a deployed one: a delete
// refused for it leaves the resource Terminating, its Accepted condition naming the reason.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Entity ID",type=string,JSONPath=`.spec.entityId`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.policyStatus`
// +kubebuilder:resource:shortName=graviteeauthzpolicies
// +kubebuilder:storageversion
type AuthzPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AuthzPolicySpec   `json:"spec,omitempty"`
	Status            AuthzPolicyStatus `json:"status,omitempty"`
}

// AuthzPolicyList contains a list of AuthzPolicy resources.
// +kubebuilder:object:root=true
type AuthzPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AuthzPolicy `json:"items"`
}

func (t *AuthzPolicy) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      t.Name,
		Namespace: t.Namespace,
	}
}

func (t *AuthzPolicy) GetSpec() core.Spec {
	return &t.Spec
}

func (t *AuthzPolicy) GetStatus() core.Status {
	return &t.Status
}

func (t *AuthzPolicy) IsBeingDeleted() bool {
	return !t.ObjectMeta.DeletionTimestamp.IsZero()
}

func (t *AuthzPolicy) HasContext() bool {
	return t.Spec.Context != nil
}

func (t *AuthzPolicy) ContextRef() core.ObjectRef {
	return t.Spec.Context
}

func (t *AuthzPolicy) GetEnvID() string {
	return t.Status.EnvID
}

func (t *AuthzPolicy) GetID() string {
	return t.Status.ID
}

func (t *AuthzPolicy) GetOrgID() string {
	return t.Status.OrgID
}

func (t *AuthzPolicy) PopulateIDs(_ core.ContextModel, _ bool) {
	// done when calling the API
}

func (t *AuthzPolicy) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(t.Status.Conditions)
}

func (t *AuthzPolicy) SetConditions(conditions []metav1.Condition) {
	t.Status.Conditions = conditions
}
