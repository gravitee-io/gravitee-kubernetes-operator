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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portallink"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.Object = &PortalLink{}
var _ core.Spec = &PortalLinkSpec{}
var _ core.Status = &PortalLinkStatus{}
var _ core.ConditionAware = &PortalLink{}

// PortalLinkSpec defines the desired state of a PortalLink.
// +kubebuilder:object:generate=true
// +kubebuilder:validation:XValidation:rule="has(self.portalRef) != has(self.apiRef)",message="exactly one of portalRef or apiRef must be set"
type PortalLinkSpec struct {
	portallink.Type `json:",inline"`
}

func (spec *PortalLinkSpec) Hash() string {
	return hash.Calculate(spec)
}

// PortalLinkStatus defines the observed state of a PortalLink.
type PortalLinkStatus struct {
	portallink.Status `json:",inline"`
}

func (s *PortalLinkStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *PortalLink:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *PortalLinkStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *PortalLink:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *PortalLinkStatus) IsFailed() bool {
	if s.Conditions != nil {
		for _, condition := range s.Conditions {
			if condition.Status == metav1.ConditionFalse {
				return true
			}
		}
	}
	return false
}

func (s *PortalLinkStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unused
}

// PortalLink attaches an external navigation link to exactly one of a Portal
// (portalRef) or an API (apiRef) at a chosen location in the owning
// resource's navigation hierarchy. The APIM management context is derived
// from the referenced Portal or API.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Portal",type=string,JSONPath=`.spec.portalRef.name`
// +kubebuilder:resource:shortName=graviteeportallinks
// +kubebuilder:storageversion
type PortalLink struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PortalLinkSpec   `json:"spec,omitempty"`
	Status            PortalLinkStatus `json:"status,omitempty"`
}

// PortalLinkList contains a list of PortalLink resources.
// +kubebuilder:object:root=true
type PortalLinkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PortalLink `json:"items"`
}

func (p *PortalLink) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      p.Name,
		Namespace: p.Namespace,
	}
}

func (p *PortalLink) GetSpec() core.Spec {
	return &p.Spec
}

func (p *PortalLink) GetStatus() core.Status {
	return &p.Status
}

func (p *PortalLink) IsBeingDeleted() bool {
	return !p.ObjectMeta.DeletionTimestamp.IsZero()
}

func (p *PortalLink) GetPortalRef() core.ObjectRef {
	return p.Spec.GetPortalRef()
}

func (p *PortalLink) GetApiRef() core.ObjectRef {
	return p.Spec.GetApiRef()
}

func (p *PortalLink) IsPortalLink() bool {
	return p.Spec.IsPortalLink()
}

func (p *PortalLink) IsApiLink() bool {
	return p.Spec.IsApiLink()
}

func (p *PortalLink) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(p.Status.Conditions)
}

func (p *PortalLink) SetConditions(conditions []metav1.Condition) {
	p.Status.Conditions = conditions
}
