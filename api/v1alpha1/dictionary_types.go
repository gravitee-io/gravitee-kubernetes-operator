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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/dictionary"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ContextAwareObject = &Dictionary{}
var _ core.Spec = &DictionarySpec{}
var _ core.Status = &DictionaryStatus{}
var _ core.ConditionAware = &Dictionary{}

// DictionarySpec defines the desired state of a Dictionary.
// +kubebuilder:object:generate=true
type DictionarySpec struct {
	dictionary.Type `json:",inline"`
	// Reference to a ManagementContext that determines which APIM instance this dictionary is synced to.
	Context *refs.NamespacedName `json:"contextRef,omitempty"`
}

func (spec *DictionarySpec) Hash() string {
	return hash.Calculate(spec)
}

// DictionaryStatus defines the observed state of a Dictionary.
type DictionaryStatus struct {
	dictionary.Status `json:",inline"`
}

func (s *DictionaryStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *Dictionary:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *DictionaryStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *Dictionary:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *DictionaryStatus) IsFailed() bool {
	if s.Conditions != nil {
		for _, condition := range s.Conditions {
			if condition.Status == metav1.ConditionFalse {
				return true
			}
		}
	}
	return false
}

func (s *DictionaryStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unused
}

// Dictionary is a Gravitee APIM dictionary managed as a Kubernetes resource.
// Dictionaries provide key/value data that can be referenced in API policies
// using Gravitee EL expressions: `{#dictionaries['hrid']['key']}`.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.type`
// +kubebuilder:resource:shortName=graviteedictionaries
// +kubebuilder:storageversion
type Dictionary struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DictionarySpec   `json:"spec,omitempty"`
	Status            DictionaryStatus `json:"status,omitempty"`
}

// DictionaryList contains a list of Dictionary resources.
// +kubebuilder:object:root=true
type DictionaryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Dictionary `json:"items"`
}

func (d *Dictionary) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      d.Name,
		Namespace: d.Namespace,
	}
}

func (d *Dictionary) GetSpec() core.Spec {
	return &d.Spec
}

func (d *Dictionary) GetStatus() core.Status {
	return &d.Status
}

func (d *Dictionary) IsBeingDeleted() bool {
	return !d.ObjectMeta.DeletionTimestamp.IsZero()
}

func (d *Dictionary) HasContext() bool {
	return d.Spec.Context != nil
}

func (d *Dictionary) ContextRef() core.ObjectRef {
	return d.Spec.Context
}

func (d *Dictionary) GetEnvID() string {
	return d.Status.EnvID
}

func (d *Dictionary) GetID() string {
	return d.Status.ID
}

func (d *Dictionary) GetOrgID() string {
	return d.Status.OrgID
}

func (d *Dictionary) PopulateIDs(_ core.ContextModel, _ bool) {
	// done when calling the API
}

func (d *Dictionary) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(d.Status.Conditions)
}

func (d *Dictionary) SetConditions(conditions []metav1.Condition) {
	d.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&Dictionary{}, &DictionaryList{})
}
