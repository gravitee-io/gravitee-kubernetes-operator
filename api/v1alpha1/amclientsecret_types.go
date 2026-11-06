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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ContextAwareObject = &AMClientSecret{}
var _ core.Spec = &AMClientSecretSpec{}
var _ core.Status = &AMClientSecretStatus{}
var _ core.ConditionAware = &AMClientSecret{}

// AMClientSecretSpec defines the desired state of an AM client secret.
// +kubebuilder:object:generate=true
type AMClientSecretSpec struct {
	// The AMApplication the secret belongs to.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="applicationRef is immutable"
	// +kubebuilder:validation:XValidation:rule="size(self.name) > 0",message="applicationRef.name must not be empty"
	ApplicationRef refs.NamespacedName `json:"applicationRef"`
	// The Secret entry holding the value of the client secret.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="secretRef is immutable"
	SecretRef am.SecretKeyRef `json:"secretRef"`
}

func (spec *AMClientSecretSpec) Hash() string {
	return hash.Calculate(spec)
}

// AMClientSecretStatus defines the observed state of an AM client secret.
type AMClientSecretStatus struct {
	am.ClientSecretStatus `json:",inline"`
}

func (s *AMClientSecretStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *AMClientSecret:
		t.Status.DeepCopyInto(s)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (s *AMClientSecretStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *AMClientSecret:
		s.DeepCopyInto(&t.Status)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (s *AMClientSecretStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unimplemented
}

func (s *AMClientSecretStatus) IsFailed() bool {
	for _, condition := range s.Conditions {
		if condition.Status == metav1.ConditionFalse {
			return true
		}
	}
	return false
}

// ClientSecretResolved holds what the reference resolvers read for the AM call.
// It is never serialized, so the value never reaches the API server.
type ClientSecretResolved struct {
	DomainKey string
	Value     string
}

// AMClientSecret is a client secret of an AMApplication. The value is read from a Secret entry and never returned.
// A secret is immutable: to rotate it, create a new Secret and AMClientSecret, move the clients,
// then delete the old AMClientSecret.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Application",type=string,JSONPath=`.spec.applicationRef.name`
// +kubebuilder:printcolumn:name="Secret",type=string,JSONPath=`.spec.secretRef.name`
// +kubebuilder:printcolumn:name="Expires At",type=string,JSONPath=`.status.expiresAt`
// +kubebuilder:resource:shortName=graviteeamclientsecrets
// +kubebuilder:storageversion
type AMClientSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AMClientSecretSpec   `json:"spec,omitempty"`
	Status AMClientSecretStatus `json:"status,omitempty"`

	Resolved ClientSecretResolved `json:"-"`
}

// +kubebuilder:object:root=true
// AMClientSecretList contains a list of AM client secrets.
type AMClientSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AMClientSecret `json:"items"`
}

func (c *AMClientSecret) IsBeingDeleted() bool {
	return !c.ObjectMeta.DeletionTimestamp.IsZero()
}

func (c *AMClientSecret) GetSpec() core.Spec {
	return &c.Spec
}

func (c *AMClientSecret) GetStatus() core.Status {
	return &c.Status
}

func (c *AMClientSecret) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      c.Name,
		Namespace: c.Namespace,
	}
}

// ContextRef is nil: the AM context comes from the application's domain.
func (c *AMClientSecret) ContextRef() core.ObjectRef {
	return nil
}

func (c *AMClientSecret) HasContext() bool {
	return false
}

func (c *AMClientSecret) GetID() string {
	return c.Status.Key
}

func (c *AMClientSecret) GetOrgID() string {
	return c.Status.OrgID
}

func (c *AMClientSecret) GetEnvID() string {
	return c.Status.EnvID
}

func (c *AMClientSecret) PopulateIDs(_ core.ContextModel, _ bool) {
	// AM is Automation API only, no HRID/UUID migration.
}

func (c *AMClientSecret) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(c.Status.Conditions)
}

func (c *AMClientSecret) SetConditions(conditions []metav1.Condition) {
	c.Status.Conditions = conditions
}
