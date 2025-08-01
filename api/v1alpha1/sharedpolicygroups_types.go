/*
 * Copyright (C) 2015 The Gravitee team (http://gravitee.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/sharedpolicygroups"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ core.ContextAwareObject = &SharedPolicyGroup{}
var _ core.ConditionAware = &SharedPolicyGroup{}

// SharedPolicyGroupSpec
// +kubebuilder:object:generate=true
type SharedPolicyGroupSpec struct {
	*sharedpolicygroups.SharedPolicyGroup `json:",inline"`
	// +kubebuilder:validation:Required
	Context *refs.NamespacedName `json:"contextRef"`
}

// Hash implements custom.Spec.
func (spec *SharedPolicyGroupSpec) Hash() string {
	return hash.Calculate(spec)
}

// SharedPolicyGroupSpecStatus defines the observed state of an API Context.
type SharedPolicyGroupSpecStatus struct {
	sharedpolicygroups.Status `json:",inline"`
}

// SharedPolicyGroup
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="description",type=string,JSONPath=`.spec.description`
// +kubebuilder:printcolumn:name="apiType",type=string,JSONPath=`.spec.apiType`
// +kubebuilder:resource:shortName=sharedpolicygroups
// +kubebuilder:storageversion
type SharedPolicyGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SharedPolicyGroupSpec       `json:"spec,omitempty"`
	Status SharedPolicyGroupSpecStatus `json:"status,omitempty"`
}

// SharedPolicyGroupList contains a list of shared policy groups.
// +kubebuilder:object:root=true
type SharedPolicyGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SharedPolicyGroup `json:"items"`
}

func (s *SharedPolicyGroup) IsBeingDeleted() bool {
	return !s.ObjectMeta.DeletionTimestamp.IsZero()
}

func (s *SharedPolicyGroup) GetSpec() core.Spec {
	return &s.Spec
}

func (s *SharedPolicyGroup) GetStatus() core.Status {
	return &s.Status
}

func (s *SharedPolicyGroup) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      s.Name,
		Namespace: s.Namespace,
	}
}

func (s *SharedPolicyGroup) ContextRef() core.ObjectRef {
	return s.Spec.Context
}

func (s *SharedPolicyGroup) HasContext() bool {
	return s.Spec.Context != nil
}

func (s *SharedPolicyGroup) GetID() string {
	return s.Status.CrossID
}

func (s *SharedPolicyGroup) PopulateIDs(_ core.ContextModel) {
	if s.Spec.CrossID != nil {
		return
	}

	if s.Status.CrossID != "" {
		s.Spec.CrossID = &s.Status.CrossID
	} else {
		s.Spec.CrossID = utils.ToReference(string(s.UID))
	}
}

func (s *SharedPolicyGroup) GetOrgID() string {
	return s.Status.OrgID
}

func (s *SharedPolicyGroup) GetEnvID() string {
	return s.Status.EnvID
}

func (s *SharedPolicyGroup) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(s.Status.Conditions)
}

func (s *SharedPolicyGroup) SetConditions(conditions []metav1.Condition) {
	s.Status.Conditions = conditions
}

func (s *SharedPolicyGroupSpecStatus) SetProcessingStatus(status core.ProcessingStatus) {
	s.Status.ProcessingStatus = status
}

func (s *SharedPolicyGroupSpecStatus) IsFailed() bool {
	return s.ProcessingStatus == core.ProcessingStatusFailed
}

func (s *SharedPolicyGroupSpecStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *SharedPolicyGroup:
		t.Status.DeepCopyInto(s)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (s *SharedPolicyGroupSpecStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *SharedPolicyGroup:
		s.DeepCopyInto(&t.Status)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func init() {
	SchemeBuilder.Register(&SharedPolicyGroup{}, &SharedPolicyGroupList{})
}
