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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portallisting"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.Object = &PortalListing{}
var _ core.Spec = &PortalListingSpec{}
var _ core.Status = &PortalListingStatus{}
var _ core.ConditionAware = &PortalListing{}

// PortalListingSpec defines the desired state of a PortalListing.
// +kubebuilder:object:generate=true
type PortalListingSpec struct {
	portallisting.Type `json:",inline"`
}

func (spec *PortalListingSpec) Hash() string {
	return hash.Calculate(spec)
}

// PortalListingStatus defines the observed state of a PortalListing.
type PortalListingStatus struct {
	portallisting.Status `json:",inline"`
}

func (s *PortalListingStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *PortalListing:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *PortalListingStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *PortalListing:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *PortalListingStatus) IsFailed() bool {
	if s.Conditions != nil {
		for _, condition := range s.Conditions {
			if condition.Status == metav1.ConditionFalse {
				return true
			}
		}
	}
	return false
}

func (s *PortalListingStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unused
}

// PortalListing publishes one or more APIs to a Portal at chosen locations in
// the portal's navigation hierarchy. The APIM management context is derived
// from the referenced Portal.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Portal",type=string,JSONPath=`.spec.portalRef.name`
// +kubebuilder:resource:shortName=graviteeportallistings
// +kubebuilder:storageversion
type PortalListing struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PortalListingSpec   `json:"spec,omitempty"`
	Status            PortalListingStatus `json:"status,omitempty"`
}

// PortalListingList contains a list of PortalListing resources.
// +kubebuilder:object:root=true
type PortalListingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PortalListing `json:"items"`
}

func (p *PortalListing) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      p.Name,
		Namespace: p.Namespace,
	}
}

func (p *PortalListing) GetSpec() core.Spec {
	return &p.Spec
}

func (p *PortalListing) GetStatus() core.Status {
	return &p.Status
}

func (p *PortalListing) IsBeingDeleted() bool {
	return !p.ObjectMeta.DeletionTimestamp.IsZero()
}

func (p *PortalListing) GetPortalRef() core.ObjectRef {
	return p.Spec.GetPortalRef()
}

func (p *PortalListing) GetApiRefs() []core.ObjectRef {
	return p.Spec.GetApiRefs()
}

func (p *PortalListing) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(p.Status.Conditions)
}

func (p *PortalListing) SetConditions(conditions []metav1.Condition) {
	p.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&PortalListing{}, &PortalListingList{})
}
