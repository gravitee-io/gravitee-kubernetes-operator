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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ContextAwareObject = &Portal{}
var _ core.Spec = &PortalSpec{}
var _ core.Status = &PortalStatus{}
var _ core.ConditionAware = &Portal{}

// PortalSpec defines the desired state of a Portal.
// +kubebuilder:object:generate=true
type PortalSpec struct {
	portal.Type `json:",inline"`
	// Reference to a ManagementContext that determines which APIM instance this portal is synced to.
	Context *refs.NamespacedName `json:"contextRef,omitempty"`
}

func (spec *PortalSpec) Hash() string {
	return hash.Calculate(spec)
}

// PortalStatus defines the observed state of a Portal.
type PortalStatus struct {
	portal.Status `json:",inline"`
}

func (s *PortalStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *Portal:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *PortalStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *Portal:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *PortalStatus) IsFailed() bool {
	if s.Conditions != nil {
		for _, condition := range s.Conditions {
			if condition.Status == metav1.ConditionFalse {
				return true
			}
		}
	}
	return false
}

func (s *PortalStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unused
}

// Portal is a Gravitee next-gen developer portal managed as a Kubernetes resource.
// A Portal is an environment-level object carrying its own navigation hierarchy.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Display Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:resource:shortName=graviteeportals
// +kubebuilder:storageversion
type Portal struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PortalSpec   `json:"spec,omitempty"`
	Status            PortalStatus `json:"status,omitempty"`
}

// PortalList contains a list of Portal resources.
// +kubebuilder:object:root=true
type PortalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Portal `json:"items"`
}

func (p *Portal) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      p.Name,
		Namespace: p.Namespace,
	}
}

func (p *Portal) GetSpec() core.Spec {
	return &p.Spec
}

func (p *Portal) GetStatus() core.Status {
	return &p.Status
}

func (p *Portal) IsBeingDeleted() bool {
	return !p.ObjectMeta.DeletionTimestamp.IsZero()
}

func (p *Portal) HasContext() bool {
	return p.Spec.Context != nil
}

func (p *Portal) ContextRef() core.ObjectRef {
	return p.Spec.Context
}

func (p *Portal) GetEnvID() string {
	return p.Status.EnvID
}

func (p *Portal) GetID() string {
	return p.Status.ID
}

func (p *Portal) GetOrgID() string {
	return p.Status.OrgID
}

func (p *Portal) PopulateIDs(_ core.ContextModel, _ bool) {
	// done when calling the API
}

func (p *Portal) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(p.Status.Conditions)
}

func (p *Portal) SetConditions(conditions []metav1.Condition) {
	p.Status.Conditions = conditions
}

