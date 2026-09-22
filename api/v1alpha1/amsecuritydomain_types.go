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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am/domain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ContextAwareObject = &AMSecurityDomain{}
var _ core.Spec = &AMSecurityDomainSpec{}
var _ core.Status = &AMSecurityDomainStatus{}
var _ core.ConditionAware = &AMSecurityDomain{}

// AMSecurityDomainSpec defines the desired state of an AM security domain.
// +kubebuilder:object:generate=true
type AMSecurityDomainSpec struct {
	domain.Domain `json:",inline"`
	// +kubebuilder:validation:Required
	Context *refs.NamespacedName `json:"contextRef" ref:"amcontext"`
}

func (spec *AMSecurityDomainSpec) Hash() string {
	return hash.Calculate(spec)
}

// AMSecurityDomainStatus defines the observed state of an AM security domain.
type AMSecurityDomainStatus struct {
	domain.Status `json:",inline"`
}

func (s *AMSecurityDomainStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *AMSecurityDomain:
		t.Status.DeepCopyInto(s)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (s *AMSecurityDomainStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *AMSecurityDomain:
		s.DeepCopyInto(&t.Status)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (s *AMSecurityDomainStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unimplemented
}

func (s *AMSecurityDomainStatus) IsFailed() bool {
	if s.Conditions != nil {
		for _, condition := range s.Conditions {
			if condition.Status == metav1.ConditionFalse {
				return true
			}
		}
	}
	return false
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="Path",type=string,JSONPath=`.spec.path`
// +kubebuilder:printcolumn:name="Enabled",type=string,JSONPath=`.spec.enabled`
// +kubebuilder:resource:shortName=amsecuritydomains
// +kubebuilder:storageversion
type AMSecurityDomain struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AMSecurityDomainSpec   `json:"spec,omitempty"`
	Status AMSecurityDomainStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// AMSecurityDomainList contains a list of AM security domains.
type AMSecurityDomainList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AMSecurityDomain `json:"items"`
}

func (d *AMSecurityDomain) IsBeingDeleted() bool {
	return !d.ObjectMeta.DeletionTimestamp.IsZero()
}

func (d *AMSecurityDomain) GetSpec() core.Spec {
	return &d.Spec
}

func (d *AMSecurityDomain) GetStatus() core.Status {
	return &d.Status
}

func (d *AMSecurityDomain) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      d.Name,
		Namespace: d.Namespace,
	}
}

func (d *AMSecurityDomain) ContextRef() core.ObjectRef {
	return d.Spec.Context
}

func (d *AMSecurityDomain) HasContext() bool {
	return d.Spec.Context != nil
}

func (d *AMSecurityDomain) GetID() string {
	return d.Status.ID
}

func (d *AMSecurityDomain) GetOrgID() string {
	return d.Status.OrgID
}

func (d *AMSecurityDomain) GetEnvID() string {
	return d.Status.EnvID
}

func (d *AMSecurityDomain) PopulateIDs(_ core.ContextModel, _ bool) {
	// AM is Automation API only, no HRID/UUID migration.
}

func (d *AMSecurityDomain) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(d.Status.Conditions)
}

func (d *AMSecurityDomain) SetConditions(conditions []metav1.Condition) {
	d.Status.Conditions = conditions
}
