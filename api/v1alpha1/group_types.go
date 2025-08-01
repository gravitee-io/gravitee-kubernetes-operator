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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ContextAwareObject = &Group{}
var _ core.Spec = &GroupSpec{}
var _ core.Status = &GroupStatus{}
var _ core.ConditionAware = &Group{}

// +kubebuilder:object:generate=true
type GroupSpec struct {
	*group.Type `json:",inline"`
	Context     *refs.NamespacedName `json:"contextRef,omitempty"`
}

func (spec *GroupSpec) Hash() string {
	return hash.Calculate(spec)
}

type GroupStatus struct {
	group.Status `json:",inline"`
}

func (s *GroupStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *Group:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *GroupStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *Group:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *GroupStatus) IsFailed() bool {
	return s.ProcessingStatus == core.ProcessingStatusFailed
}

func (s *GroupStatus) SetProcessingStatus(status core.ProcessingStatus) {
	s.ProcessingStatus = status
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Members at",type=string,JSONPath=`.status.members`,description="The number of members added to the group"
// +kubebuilder:storageversion
type Group struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              GroupSpec   `json:"spec,omitempty"`
	Status            GroupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type GroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Group `json:"items"`
}

func (g *Group) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      g.Name,
		Namespace: g.Namespace,
	}
}

func (g *Group) GetSpec() core.Spec {
	return &g.Spec
}

func (g *Group) GetStatus() core.Status {
	return &g.Status
}

func (g *Group) IsBeingDeleted() bool {
	return !g.ObjectMeta.DeletionTimestamp.IsZero()
}

func (g *Group) HasContext() bool {
	return g.Spec.Context != nil
}

func (g *Group) ContextRef() core.ObjectRef {
	return g.Spec.Context
}

func (g *Group) GetEnvID() string {
	return g.Status.EnvID
}

func (g *Group) GetID() string {
	return g.Status.ID
}

func (g *Group) GetOrgID() string {
	return g.Status.OrgID
}

func (g *Group) PopulateIDs(mCtx core.ContextModel) {
	g.Spec.ID = g.pickID(mCtx)
}

func (g *Group) pickID(mCtx core.ContextModel) string {
	if g.Status.ID != "" {
		return g.Status.ID
	}

	if g.Spec.ID != "" {
		return g.Spec.ID
	}

	if mCtx != nil {
		return uuid.FromStrings(g.Spec.Name, mCtx.GetOrgID(), mCtx.GetEnvID())
	}

	return string(g.UID)
}

func (g *Group) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(g.Status.Conditions)
}

func (g *Group) SetConditions(conditions []metav1.Condition) {
	g.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&Group{}, &GroupList{})
}
