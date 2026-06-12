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

	documentation "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/docs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.Object = &Documentation{}
var _ core.Spec = &DocumentationSpec{}
var _ core.Status = &DocumentationStatus{}
var _ core.ConditionAware = &Documentation{}

// DocumentationSpec defines the desired state of a Documentation.
// +kubebuilder:object:generate=true
// +kubebuilder:validation:XValidation:rule="has(self.portalRef) != has(self.apiRef)",message="exactly one of portalRef or apiRef must be set"
type DocumentationSpec struct {
	documentation.Type `json:",inline"`
}

func (spec *DocumentationSpec) Hash() string {
	return hash.Calculate(spec)
}

// DocumentationStatus defines the observed state of a Documentation.
type DocumentationStatus struct {
	documentation.Status `json:",inline"`
}

func (s *DocumentationStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *Documentation:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *DocumentationStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *Documentation:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *DocumentationStatus) IsFailed() bool {
	if s.Conditions != nil {
		for _, condition := range s.Conditions {
			if condition.Status == metav1.ConditionFalse {
				return true
			}
		}
	}
	return false
}

func (s *DocumentationStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unused
}

// Documentation is a single page of documentation (Gravitee Markdown, OpenAPI
// or AsyncAPI) for the next-gen developer portal. A Documentation page is
// attached to exactly one of a Portal or an API; the APIM management context is
// derived from the referenced resource.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.type`
// +kubebuilder:printcolumn:name="Location",type=string,JSONPath=`.spec.location`
// +kubebuilder:resource:shortName=graviteedocumentations
// +kubebuilder:storageversion
type Documentation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DocumentationSpec   `json:"spec,omitempty"`
	Status            DocumentationStatus `json:"status,omitempty"`
}

// DocumentationList contains a list of Documentation resources.
// +kubebuilder:object:root=true
type DocumentationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Documentation `json:"items"`
}

func (d *Documentation) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      d.Name,
		Namespace: d.Namespace,
	}
}

func (d *Documentation) GetSpec() core.Spec {
	return &d.Spec
}

func (d *Documentation) GetStatus() core.Status {
	return &d.Status
}

func (d *Documentation) IsBeingDeleted() bool {
	return !d.ObjectMeta.DeletionTimestamp.IsZero()
}

func (d *Documentation) IsPortalDoc() bool {
	return d.Spec.IsPortalDoc()
}

func (d *Documentation) IsApiDoc() bool {
	return d.Spec.IsApiDoc()
}

func (d *Documentation) GetPortalRef() core.ObjectRef {
	return d.Spec.GetPortalRef()
}

func (d *Documentation) GetApiRef() core.ObjectRef {
	return d.Spec.GetApiRef()
}

func (d *Documentation) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(d.Status.Conditions)
}

func (d *Documentation) SetConditions(conditions []metav1.Condition) {
	d.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&Documentation{}, &DocumentationList{})
}
