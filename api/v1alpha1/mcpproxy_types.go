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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ContextAwareObject = &McpProxy{}
var _ core.Spec = &McpProxySpec{}
var _ core.Status = &McpProxyStatus{}
var _ core.ConditionAware = &McpProxy{}

// McpProxySpec defines the desired state of an McpProxy.
// +kubebuilder:object:generate=true
// +kubebuilder:validation:XValidation:rule="self.mode != 'PROXY' || has(self.proxy)",message="proxy must be set when mode is PROXY"
// +kubebuilder:validation:XValidation:rule="self.mode == 'PROXY' || !has(self.proxy)",message="proxy must not be set when mode is not PROXY"
// +kubebuilder:validation:XValidation:rule="self.mode != 'STUDIO' || has(self.studio)",message="studio must be set when mode is STUDIO"
// +kubebuilder:validation:XValidation:rule="self.mode == 'STUDIO' || !has(self.studio)",message="studio must not be set when mode is not STUDIO"
// +kubebuilder:validation:XValidation:rule="self.mode == oldSelf.mode",message="mode is immutable"
// +kubebuilder:validation:XValidation:rule="self.entityId == oldSelf.entityId",message="entityId is immutable"
// +kubebuilder:validation:XValidation:rule="self.name == oldSelf.name",message="name is immutable"
// +kubebuilder:validation:XValidation:rule="self.protocolVersion == oldSelf.protocolVersion",message="protocolVersion is immutable"
// +kubebuilder:validation:XValidation:rule="has(self.description) == has(oldSelf.description) && (!has(self.description) || self.description == oldSelf.description)",message="description is immutable"
type McpProxySpec struct {
	mcpproxy.Type `json:",inline"`
	// Reference to a ManagementContext that determines which APIM instance this proxy is created in.
	Context *refs.NamespacedName `json:"contextRef,omitempty"`
}

func (spec *McpProxySpec) Hash() string {
	return hash.Calculate(spec)
}

// McpProxyStatus defines the observed state of a McpProxy.
type McpProxyStatus struct {
	mcpproxy.Status `json:",inline"`
}

func (s *McpProxyStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *McpProxy:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *McpProxyStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *McpProxy:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *McpProxyStatus) IsFailed() bool {
	if s.Conditions != nil {
		for _, condition := range s.Conditions {
			if condition.Status == metav1.ConditionFalse {
				return true
			}
		}
	}
	return false
}

func (s *McpProxyStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unused
}

// McpProxy is a Gravitee MCP proxy as a Kubernetes resource: either a single upstream MCP
// server fronted whole (mode PROXY) or a Studio exposing tools selected across
// CatalogMcpServer resources (mode STUDIO).
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Entity ID",type=string,JSONPath=`.spec.entityId`
// +kubebuilder:printcolumn:name="Mode",type=string,JSONPath=`.spec.mode`
// +kubebuilder:printcolumn:name="Context Path",type=string,JSONPath=`.spec.contextPath`
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`
// +kubebuilder:resource:shortName=graviteemcpproxies
// +kubebuilder:storageversion
type McpProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              McpProxySpec   `json:"spec,omitempty"`
	Status            McpProxyStatus `json:"status,omitempty"`
}

// McpProxyList contains a list of McpProxy resources.
// +kubebuilder:object:root=true
type McpProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []McpProxy `json:"items"`
}

func (t *McpProxy) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      t.Name,
		Namespace: t.Namespace,
	}
}

func (t *McpProxy) GetSpec() core.Spec {
	return &t.Spec
}

func (t *McpProxy) GetStatus() core.Status {
	return &t.Status
}

func (t *McpProxy) IsBeingDeleted() bool {
	return !t.ObjectMeta.DeletionTimestamp.IsZero()
}

func (t *McpProxy) HasContext() bool {
	return t.Spec.Context != nil
}

func (t *McpProxy) ContextRef() core.ObjectRef {
	return t.Spec.Context
}

func (t *McpProxy) GetEnvID() string {
	return t.Status.EnvID
}

func (t *McpProxy) GetID() string {
	return t.Status.ID
}

func (t *McpProxy) GetOrgID() string {
	return t.Status.OrgID
}

func (t *McpProxy) PopulateIDs(_ core.ContextModel, _ bool) {
	// done when calling the API
}

func (t *McpProxy) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(t.Status.Conditions)
}

func (t *McpProxy) SetConditions(conditions []metav1.Condition) {
	t.Status.Conditions = conditions
}
