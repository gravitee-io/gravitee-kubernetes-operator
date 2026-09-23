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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/catalogmcpserver"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ContextAwareObject = &CatalogMcpServer{}
var _ core.Spec = &CatalogMcpServerSpec{}
var _ core.Status = &CatalogMcpServerStatus{}
var _ core.ConditionAware = &CatalogMcpServer{}

// CatalogMcpServerSpec defines the desired state of a CatalogMcpServer.
// +kubebuilder:object:generate=true
// +kubebuilder:validation:XValidation:rule="self.entityId == oldSelf.entityId",message="entityId is immutable"
type CatalogMcpServerSpec struct {
	catalogmcpserver.Type `json:",inline"`
	// Reference to a ManagementContext that determines which APIM instance this server is registered in.
	Context *refs.NamespacedName `json:"contextRef,omitempty"`
}

func (spec *CatalogMcpServerSpec) Hash() string {
	return hash.Calculate(spec)
}

// CatalogMcpServerStatus defines the observed state of a CatalogMcpServer.
type CatalogMcpServerStatus struct {
	catalogmcpserver.Status `json:",inline"`
}

func (s *CatalogMcpServerStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *CatalogMcpServer:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *CatalogMcpServerStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *CatalogMcpServer:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *CatalogMcpServerStatus) IsFailed() bool {
	if s.Conditions != nil {
		for _, condition := range s.Conditions {
			if condition.Status == metav1.ConditionFalse {
				return true
			}
		}
	}
	return false
}

func (s *CatalogMcpServerStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unused
}

// CatalogMcpServer is an upstream MCP server registered in the Gravitee AI Catalog as a
// Kubernetes resource. The platform discovers the tools, prompts and resources the server
// exposes and reports them in the status.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Entity ID",type=string,JSONPath=`.spec.entityId`
// +kubebuilder:printcolumn:name="Endpoint",type=string,JSONPath=`.spec.connection.endpoint`
// +kubebuilder:printcolumn:name="Synced",type=string,JSONPath=`.status.lastSyncedAt`
// +kubebuilder:resource:shortName=graviteecatalogmcpservers
// +kubebuilder:storageversion
type CatalogMcpServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CatalogMcpServerSpec   `json:"spec,omitempty"`
	Status            CatalogMcpServerStatus `json:"status,omitempty"`
}

// CatalogMcpServerList contains a list of CatalogMcpServer resources.
// +kubebuilder:object:root=true
type CatalogMcpServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CatalogMcpServer `json:"items"`
}

func (t *CatalogMcpServer) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      t.Name,
		Namespace: t.Namespace,
	}
}

func (t *CatalogMcpServer) GetSpec() core.Spec {
	return &t.Spec
}

func (t *CatalogMcpServer) GetStatus() core.Status {
	return &t.Status
}

func (t *CatalogMcpServer) IsBeingDeleted() bool {
	return !t.ObjectMeta.DeletionTimestamp.IsZero()
}

func (t *CatalogMcpServer) HasContext() bool {
	return t.Spec.Context != nil
}

func (t *CatalogMcpServer) ContextRef() core.ObjectRef {
	return t.Spec.Context
}

func (t *CatalogMcpServer) GetEnvID() string {
	return t.Status.EnvID
}

func (t *CatalogMcpServer) GetID() string {
	return t.Status.ID
}

func (t *CatalogMcpServer) GetOrgID() string {
	return t.Status.OrgID
}

func (t *CatalogMcpServer) PopulateIDs(_ core.ContextModel, _ bool) {
	// done when calling the API
}

func (t *CatalogMcpServer) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(t.Status.Conditions)
}

func (t *CatalogMcpServer) SetConditions(conditions []metav1.Condition) {
	t.Status.Conditions = conditions
}
