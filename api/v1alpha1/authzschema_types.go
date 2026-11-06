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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzschema"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ContextAwareObject = &AuthzSchema{}
var _ core.Spec = &AuthzSchemaSpec{}
var _ core.Status = &AuthzSchemaStatus{}
var _ core.ConditionAware = &AuthzSchema{}

// AuthzSchemaSpec defines the desired state of an AuthzSchema.
// +kubebuilder:object:generate=true
type AuthzSchemaSpec struct {
	authzschema.Type `json:",inline"`
	// Reference to a ManagementContext that determines which APIM instance this schema document
	// is created in. Immutable: the operator deletes a schema through its current context only, so
	// a schema moved to another one would stay on the first platform, published to its target.
	// The message is static: refs.NamespacedName strings carry no maxLength, so a message naming
	// both values exceeds the API server's CEL cost budget and the CRD is refused. Required, so the
	// rule also holds without the admission webhook: CEL compares the field only when both sides have it.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="contextRef is immutable; delete and recreate the AuthzSchema"
	Context *refs.NamespacedName `json:"contextRef"`
}

func (spec *AuthzSchemaSpec) Hash() string {
	return hash.Calculate(spec)
}

// AuthzSchemaStatus defines the observed state of an AuthzSchema.
type AuthzSchemaStatus struct {
	authzschema.Status `json:",inline"`
}

func (s *AuthzSchemaStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *AuthzSchema:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *AuthzSchemaStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *AuthzSchema:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *AuthzSchemaStatus) IsFailed() bool {
	if s.Conditions != nil {
		for _, condition := range s.Conditions {
			if condition.Status == metav1.ConditionFalse {
				return true
			}
		}
	}
	return false
}

func (s *AuthzSchemaStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unused
}

// AuthzSchema is a GAPL schema document of Gravitee Authorization Management as a Kubernetes
// resource, published to the PDP gateway its target names. A PDP gateway holds one schema
// document: a target another schema holds is refused.
//
// A target names an AuthzPdp through pdpRef, or a routing scope no AuthzPdp declares through
// scope. A pdpRef waits on ResolvedRefs=False until the AuthzPdp has a routing scope. Deleting
// that AuthzPdp removes the target on the platform: the document stays stored and published
// nowhere, and the AuthzSchema waits on ResolvedRefs=False until its pdpRef names another AuthzPdp.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Target",type=string,JSONPath=`.status.target`
// +kubebuilder:resource:shortName=graviteeauthzschemas
// +kubebuilder:storageversion
type AuthzSchema struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AuthzSchemaSpec   `json:"spec,omitempty"`
	Status            AuthzSchemaStatus `json:"status,omitempty"`
}

// AuthzSchemaList contains a list of AuthzSchema resources.
// +kubebuilder:object:root=true
type AuthzSchemaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AuthzSchema `json:"items"`
}

func (t *AuthzSchema) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      t.Name,
		Namespace: t.Namespace,
	}
}

func (t *AuthzSchema) GetSpec() core.Spec {
	return &t.Spec
}

func (t *AuthzSchema) GetStatus() core.Status {
	return &t.Status
}

func (t *AuthzSchema) IsBeingDeleted() bool {
	return !t.ObjectMeta.DeletionTimestamp.IsZero()
}

func (t *AuthzSchema) HasContext() bool {
	return t.Spec.Context != nil
}

func (t *AuthzSchema) ContextRef() core.ObjectRef {
	return t.Spec.Context
}

func (t *AuthzSchema) GetEnvID() string {
	return t.Status.EnvID
}

func (t *AuthzSchema) GetID() string {
	return t.Status.ID
}

func (t *AuthzSchema) GetOrgID() string {
	return t.Status.OrgID
}

func (t *AuthzSchema) PopulateIDs(_ core.ContextModel, _ bool) {
	// done when calling the API
}

func (t *AuthzSchema) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(t.Status.Conditions)
}

func (t *AuthzSchema) SetConditions(conditions []metav1.Condition) {
	t.Status.Conditions = conditions
}
