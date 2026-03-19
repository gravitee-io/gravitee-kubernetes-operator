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
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/idpgroupmapping"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ContextAwareObject = &IDPGroupMapping{}
var _ core.Spec = &IDPGroupMappingSpec{}
var _ core.Status = &IDPGroupMappingStatus{}
var _ core.ConditionAware = &IDPGroupMapping{}

// +kubebuilder:object:generate=true
type IDPGroupMappingSpec struct {
	*idpgroupmapping.Type `json:",inline"`
	Context               *refs.NamespacedName `json:"contextRef,omitempty"`
}

func (spec *IDPGroupMappingSpec) Hash() string {
	return hash.Calculate(spec)
}

type IDPGroupMappingStatus struct {
	idpgroupmapping.Status `json:",inline"`
}

func (s *IDPGroupMappingStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *IDPGroupMapping:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *IDPGroupMappingStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *IDPGroupMapping:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *IDPGroupMappingStatus) IsFailed() bool {
	return s.ProcessingStatus == core.ProcessingStatusFailed
}

func (s *IDPGroupMappingStatus) SetProcessingStatus(status core.ProcessingStatus) {
	s.ProcessingStatus = status
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Members at",type=string,JSONPath=`.status.members`,description="The number of members added to the group"
// +kubebuilder:storageversion
type IDPGroupMapping struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              IDPGroupMappingSpec   `json:"spec,omitempty"`
	Status            IDPGroupMappingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type IDPGroupMappingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IDPGroupMapping `json:"items"`
}

func (g *IDPGroupMapping) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      g.Name,
		Namespace: g.Namespace,
	}
}

func (g *IDPGroupMapping) GetSpec() core.Spec {
	return &g.Spec
}

func (g *IDPGroupMapping) GetStatus() core.Status {
	return &g.Status
}

func (g *IDPGroupMapping) IsBeingDeleted() bool {
	return !g.ObjectMeta.DeletionTimestamp.IsZero()
}

func (g *IDPGroupMapping) HasContext() bool {
	return g.Spec.Context != nil
}

func (g *IDPGroupMapping) ContextRef() core.ObjectRef {
	return g.Spec.Context
}

func (g *IDPGroupMapping) GetEnvID() string {
	return g.Status.EnvID
}

func (g *IDPGroupMapping) GetID() string {
	return g.Status.ID
}

func (g *IDPGroupMapping) GetOrgID() string {
	return g.Status.OrgID
}

func (g *IDPGroupMapping) PopulateIDs(mCtx core.ContextModel) {
	g.Spec.ID = g.pickID(mCtx)
}

func (g *IDPGroupMapping) pickID(mCtx core.ContextModel) string {
	if g.Status.ID != "" {
		return g.Status.ID
	}

	if g.Spec.ID != "" {
		return g.Spec.ID
	}

	if mCtx != nil {
		return uuid.FromStrings(strings.Join(g.Spec.Groups, ","), g.Spec.Condition, mCtx.GetOrgID(), mCtx.GetEnvID())
	}

	return string(g.UID)
}

func (g *IDPGroupMapping) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(g.Status.Conditions)
}

func (g *IDPGroupMapping) SetConditions(conditions []metav1.Condition) {
	g.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&IDPGroupMapping{}, &IDPGroupMappingList{})
}
