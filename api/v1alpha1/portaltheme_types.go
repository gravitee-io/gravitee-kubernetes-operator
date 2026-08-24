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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portaltheme"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ContextAwareObject = &PortalTheme{}
var _ core.Spec = &PortalThemeSpec{}
var _ core.Status = &PortalThemeStatus{}
var _ core.ConditionAware = &PortalTheme{}

// PortalThemeSpec defines the desired state of a PortalTheme.
// +kubebuilder:object:generate=true
type PortalThemeSpec struct {
	portaltheme.Type `json:",inline"`
	// Reference to a ManagementContext that determines which APIM instance this theme is synced to.
	Context *refs.NamespacedName `json:"contextRef,omitempty"`
}

func (spec *PortalThemeSpec) Hash() string {
	return hash.Calculate(spec)
}

// PortalThemeStatus defines the observed state of a PortalTheme.
type PortalThemeStatus struct {
	portaltheme.Status `json:",inline"`
}

func (s *PortalThemeStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *PortalTheme:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *PortalThemeStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *PortalTheme:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *PortalThemeStatus) IsFailed() bool {
	if s.Conditions != nil {
		for _, condition := range s.Conditions {
			if condition.Status == metav1.ConditionFalse {
				return true
			}
		}
	}
	return false
}

func (s *PortalThemeStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unused
}

// PortalTheme is the look and feel of a Gravitee next-gen developer portal managed as a
// Kubernetes resource. Applying one makes it the active theme of its environment.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Display Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:resource:shortName=graviteeportalthemes
// +kubebuilder:storageversion
type PortalTheme struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PortalThemeSpec   `json:"spec,omitempty"`
	Status            PortalThemeStatus `json:"status,omitempty"`
}

// PortalThemeList contains a list of PortalTheme resources.
// +kubebuilder:object:root=true
type PortalThemeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PortalTheme `json:"items"`
}

func (t *PortalTheme) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      t.Name,
		Namespace: t.Namespace,
	}
}

func (t *PortalTheme) GetSpec() core.Spec {
	return &t.Spec
}

func (t *PortalTheme) GetStatus() core.Status {
	return &t.Status
}

func (t *PortalTheme) IsBeingDeleted() bool {
	return !t.ObjectMeta.DeletionTimestamp.IsZero()
}

func (t *PortalTheme) HasContext() bool {
	return t.Spec.Context != nil
}

func (t *PortalTheme) ContextRef() core.ObjectRef {
	return t.Spec.Context
}

func (t *PortalTheme) GetEnvID() string {
	return t.Status.EnvID
}

func (t *PortalTheme) GetID() string {
	return t.Status.ID
}

func (t *PortalTheme) GetOrgID() string {
	return t.Status.OrgID
}

func (t *PortalTheme) PopulateIDs(_ core.ContextModel, _ bool) {
	// done when calling the API
}

func (t *PortalTheme) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(t.Status.Conditions)
}

func (t *PortalTheme) SetConditions(conditions []metav1.Condition) {
	t.Status.Conditions = conditions
}
