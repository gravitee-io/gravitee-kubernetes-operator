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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.AMDomainSubResource = &AMApplication{}
var _ core.ContextAwareObject = &AMApplication{}
var _ core.Spec = &AMApplicationSpec{}
var _ core.Status = &AMApplicationStatus{}
var _ core.ConditionAware = &AMApplication{}

// AMApplicationSpec defines the desired state of an AM application.
// +kubebuilder:object:generate=true
type AMApplicationSpec struct {
	application.Application `json:",inline"`
}

func (spec *AMApplicationSpec) Hash() string {
	return hash.Calculate(spec)
}

// AMApplicationStatus defines the observed state of an AM application.
type AMApplicationStatus struct {
	am.ApplicationStatus `json:",inline"`
}

func (s *AMApplicationStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *AMApplication:
		t.Status.DeepCopyInto(s)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (s *AMApplicationStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *AMApplication:
		s.DeepCopyInto(&t.Status)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (s *AMApplicationStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unimplemented
}

func (s *AMApplicationStatus) IsFailed() bool {
	for _, condition := range s.Conditions {
		if condition.Status == metav1.ConditionFalse {
			return true
		}
	}
	return false
}

// AMApplication is an AM application (OAuth 2.0 / OIDC client) of a security domain.
// It has no client secret until an AMClientSecret references it; its effective clientId is in status.clientId.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.type`
// +kubebuilder:printcolumn:name="Client ID",type=string,JSONPath=`.status.clientId`
// +kubebuilder:printcolumn:name="Domain",type=string,JSONPath=`.spec.domainRef.name`
// +kubebuilder:printcolumn:name="Enabled",type=boolean,JSONPath=`.spec.enabled`
// +kubebuilder:resource:shortName=graviteeamapplications
// +kubebuilder:storageversion
type AMApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AMApplicationSpec   `json:"spec,omitempty"`
	Status AMApplicationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// AMApplicationList contains a list of AM applications.
type AMApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AMApplication `json:"items"`
}

func (c *AMApplication) IsBeingDeleted() bool {
	return !c.ObjectMeta.DeletionTimestamp.IsZero()
}

func (c *AMApplication) GetSpec() core.Spec {
	return &c.Spec
}

func (c *AMApplication) GetStatus() core.Status {
	return &c.Status
}

func (c *AMApplication) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      c.Name,
		Namespace: c.Namespace,
	}
}

// ContextRef is nil: the AM context comes from the domain.
func (c *AMApplication) ContextRef() core.ObjectRef {
	return nil
}

func (c *AMApplication) HasContext() bool {
	return false
}

func (c *AMApplication) GetID() string {
	return c.Status.Key
}

func (c *AMApplication) GetOrgID() string {
	return c.Status.OrgID
}

func (c *AMApplication) GetEnvID() string {
	return c.Status.EnvID
}

func (c *AMApplication) PopulateIDs(_ core.ContextModel, _ bool) {
	// AM is Automation API only, no HRID/UUID migration.
}

func (c *AMApplication) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(c.Status.Conditions)
}

func (c *AMApplication) SetConditions(conditions []metav1.Condition) {
	c.Status.Conditions = conditions
}

func (c *AMApplication) GetDomainRef() core.ObjectRef {
	return &c.Spec.DomainRef
}
