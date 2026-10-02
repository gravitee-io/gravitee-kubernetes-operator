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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am/reporter"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.AMDomainSubResource = &AMReporter{}
var _ core.ContextAwareObject = &AMReporter{}
var _ core.Spec = &AMReporterSpec{}
var _ core.Status = &AMReporterStatus{}
var _ core.ConditionAware = &AMReporter{}

// AMReporterSpec defines the desired state of an AM reporter.
// +kubebuilder:object:generate=true
type AMReporterSpec struct {
	reporter.Reporter `json:",inline"`
}

func (spec *AMReporterSpec) Hash() string {
	return hash.Calculate(spec)
}

// AMReporterStatus defines the observed state of an AM reporter.
type AMReporterStatus struct {
	am.ReporterStatus `json:",inline"`
}

func (s *AMReporterStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *AMReporter:
		t.Status.DeepCopyInto(s)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (s *AMReporterStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *AMReporter:
		s.DeepCopyInto(&t.Status)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (s *AMReporterStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unimplemented
}

func (s *AMReporterStatus) IsFailed() bool {
	for _, condition := range s.Conditions {
		if condition.Status == metav1.ConditionFalse {
			return true
		}
	}
	return false
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.status.name`
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.status.type`
// +kubebuilder:printcolumn:name="Domain",type=string,JSONPath=`.spec.domainRef.name`
// +kubebuilder:printcolumn:name="Enabled",type=boolean,JSONPath=`.spec.enabled`
// +kubebuilder:resource:shortName=graviteeamreporters
// +kubebuilder:storageversion
type AMReporter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AMReporterSpec   `json:"spec,omitempty"`
	Status AMReporterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// AMReporterList contains a list of AM reporters.
type AMReporterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AMReporter `json:"items"`
}

func (c *AMReporter) IsBeingDeleted() bool {
	return !c.ObjectMeta.DeletionTimestamp.IsZero()
}

func (c *AMReporter) GetSpec() core.Spec {
	return &c.Spec
}

func (c *AMReporter) GetStatus() core.Status {
	return &c.Status
}

func (c *AMReporter) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      c.Name,
		Namespace: c.Namespace,
	}
}

// ContextRef is nil: the AM context comes from the domain.
func (c *AMReporter) ContextRef() core.ObjectRef {
	return nil
}

func (c *AMReporter) HasContext() bool {
	return false
}

func (c *AMReporter) GetID() string {
	return c.Status.Key
}

func (c *AMReporter) GetOrgID() string {
	return c.Status.OrgID
}

func (c *AMReporter) GetEnvID() string {
	return c.Status.EnvID
}

func (c *AMReporter) PopulateIDs(_ core.ContextModel, _ bool) {
	// AM is Automation API only, no HRID/UUID migration.
}

func (c *AMReporter) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(c.Status.Conditions)
}

func (c *AMReporter) SetConditions(conditions []metav1.Condition) {
	c.Status.Conditions = conditions
}

func (c *AMReporter) GetDomainRef() core.ObjectRef {
	return &c.Spec.DomainRef
}
