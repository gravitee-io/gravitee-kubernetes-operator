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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am/identityprovider"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.AMDomainSubResource = &AMIdentityProvider{}
var _ core.ContextAwareObject = &AMIdentityProvider{}
var _ core.Spec = &AMIdentityProviderSpec{}
var _ core.Status = &AMIdentityProviderStatus{}
var _ core.ConditionAware = &AMIdentityProvider{}
var _ core.AMDomainSubResource = &AMIdentityProvider{}

// AMIdentityProviderSpec defines the desired state of an AM identity provider.
// +kubebuilder:object:generate=true
type AMIdentityProviderSpec struct {
	identityprovider.IdentityProvider `json:",inline"`
}

func (spec *AMIdentityProviderSpec) Hash() string {
	return hash.Calculate(spec)
}

// AMIdentityProviderStatus defines the observed state of an AM identity provider.
type AMIdentityProviderStatus struct {
	am.IdentityProviderStatus `json:",inline"`
}

func (s *AMIdentityProviderStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *AMIdentityProvider:
		t.Status.DeepCopyInto(s)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (s *AMIdentityProviderStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *AMIdentityProvider:
		s.DeepCopyInto(&t.Status)
		return nil
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (s *AMIdentityProviderStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unimplemented
}

func (s *AMIdentityProviderStatus) IsFailed() bool {
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
// +kubebuilder:resource:shortName=graviteeamidentityproviders
// +kubebuilder:storageversion
type AMIdentityProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AMIdentityProviderSpec   `json:"spec,omitempty"`
	Status AMIdentityProviderStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// AMIdentityProviderList contains a list of AM identity providers.
type AMIdentityProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AMIdentityProvider `json:"items"`
}

func (p *AMIdentityProvider) IsBeingDeleted() bool {
	return !p.ObjectMeta.DeletionTimestamp.IsZero()
}

func (p *AMIdentityProvider) GetSpec() core.Spec {
	return &p.Spec
}

func (p *AMIdentityProvider) GetStatus() core.Status {
	return &p.Status
}

func (p *AMIdentityProvider) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      p.Name,
		Namespace: p.Namespace,
	}
}

// ContextRef is nil: the AM context comes from the domain.
func (p *AMIdentityProvider) ContextRef() core.ObjectRef {
	return nil
}

func (p *AMIdentityProvider) HasContext() bool {
	return false
}

func (p *AMIdentityProvider) GetID() string {
	return p.Status.Key
}

func (p *AMIdentityProvider) GetOrgID() string {
	return p.Status.OrgID
}

func (p *AMIdentityProvider) GetEnvID() string {
	return p.Status.EnvID
}

func (p *AMIdentityProvider) PopulateIDs(_ core.ContextModel, _ bool) {
	// AM is Automation API only, no HRID/UUID migration.
}

func (p *AMIdentityProvider) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(p.Status.Conditions)
}

func (p *AMIdentityProvider) SetConditions(conditions []metav1.Condition) {
	p.Status.Conditions = conditions
}

func (p *AMIdentityProvider) GetDomainRef() core.ObjectRef {
	return &p.Spec.DomainRef
}
