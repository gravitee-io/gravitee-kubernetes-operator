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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am/certificate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.AMDomainSubResource = &AMCertificate{}
var _ core.ContextAwareObject = &AMCertificate{}
var _ core.Spec = &AMCertificateSpec{}
var _ core.Status = &AMCertificateStatus{}
var _ core.ConditionAware = &AMCertificate{}

// AMCertificateSpec defines the desired state of an AM certificate.
// +kubebuilder:object:generate=true
type AMCertificateSpec struct {
	certificate.Certificate `json:",inline"`
}

func (spec *AMCertificateSpec) Hash() string {
	return hash.Calculate(spec)
}

// AMCertificateStatus defines the observed state of an AM certificate.
type AMCertificateStatus struct {
	am.CertificateStatus `json:",inline"`
}

func (s *AMCertificateStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *AMCertificate:
		t.Status.DeepCopyInto(s)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (s *AMCertificateStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *AMCertificate:
		s.DeepCopyInto(&t.Status)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (s *AMCertificateStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unimplemented
}

func (s *AMCertificateStatus) IsFailed() bool {
	for _, condition := range s.Conditions {
		if condition.Status == metav1.ConditionFalse {
			return true
		}
	}
	return false
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.type`
// +kubebuilder:printcolumn:name="Domain",type=string,JSONPath=`.spec.domainRef.name`
// +kubebuilder:printcolumn:name="Expires",type=date,JSONPath=`.status.expiresAt`
// +kubebuilder:resource:shortName=graviteeamcertificates
// +kubebuilder:storageversion
type AMCertificate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AMCertificateSpec   `json:"spec,omitempty"`
	Status AMCertificateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// AMCertificateList contains a list of AM certificates.
type AMCertificateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AMCertificate `json:"items"`
}

func (c *AMCertificate) IsBeingDeleted() bool {
	return !c.ObjectMeta.DeletionTimestamp.IsZero()
}

func (c *AMCertificate) GetSpec() core.Spec {
	return &c.Spec
}

func (c *AMCertificate) GetStatus() core.Status {
	return &c.Status
}

func (c *AMCertificate) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      c.Name,
		Namespace: c.Namespace,
	}
}

// ContextRef is nil: the AM context comes from the domain.
func (c *AMCertificate) ContextRef() core.ObjectRef {
	return nil
}

func (c *AMCertificate) HasContext() bool {
	return false
}

func (c *AMCertificate) GetID() string {
	return c.Status.Key
}

func (c *AMCertificate) GetOrgID() string {
	return c.Status.OrgID
}

func (c *AMCertificate) GetEnvID() string {
	return c.Status.EnvID
}

func (c *AMCertificate) PopulateIDs(_ core.ContextModel, _ bool) {
	// AM is Automation API only, no HRID/UUID migration.
}

func (c *AMCertificate) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(c.Status.Conditions)
}

func (c *AMCertificate) SetConditions(conditions []metav1.Condition) {
	c.Status.Conditions = conditions
}

func (c *AMCertificate) GetDomainRef() core.ObjectRef {
	return &c.Spec.DomainRef
}
