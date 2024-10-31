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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.SubscriptionObject = &Subscription{}
var _ core.Spec = &SubscriptionSpec{}
var _ core.Status = &SubscriptionStatus{}

// +kubebuilder:object:generate=true
type SubscriptionSpec struct {
	subscription.Type `json:",inline"`
}

func (spec *SubscriptionSpec) Hash() string {
	return hash.Calculate(spec)
}

type SubscriptionStatus struct {
	// Subscription ID
	ID string `json:"id,omitempty"`
	// When the subscription was started and made available
	StartedAt string `json:"startedAt,omitempty"`
	// The expiry date for the subscription (no date means no expiry)
	EndingAt string `json:"endingAt,omitempty"`
	// This value is `Completed` if the sync with APIM succeeded, Failed otherwise.
	ProcessingStatus core.ProcessingStatus `json:"processingStatus,omitempty"`
}

func (s *SubscriptionStatus) DeepCopyFrom(obj client.Object) error {
	switch t := obj.(type) {
	case *Subscription:
		t.Status.DeepCopyInto(s)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *SubscriptionStatus) DeepCopyTo(obj client.Object) error {
	switch t := obj.(type) {
	case *Subscription:
		s.DeepCopyInto(&t.Status)
	default:
		return fmt.Errorf("unknown type %T", t)
	}

	return nil
}

func (s *SubscriptionStatus) IsFailed() bool {
	return s.ProcessingStatus == core.ProcessingStatusFailed
}

func (s *SubscriptionStatus) SetProcessingStatus(status core.ProcessingStatus) {
	s.ProcessingStatus = status
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Started at",type=string,JSONPath=`.status.startedAt`,description="The date from when the subscription starts"
// +kubebuilder:printcolumn:name="Ending at",type=string,JSONPath=`.status.endingAt`,description="The date when the subscription expires"
// +kubebuilder:storageversion
type Subscription struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubscriptionSpec   `json:"spec,omitempty"`
	Status SubscriptionStatus `json:"status,omitempty"`
}

func (s *Subscription) GetApiRef() core.ObjectRef {
	return s.Spec.GetApiRef()
}

func (s *Subscription) GetAppRef() core.ObjectRef {
	return s.Spec.GetAppRef()
}

func (s *Subscription) GetPlan() string {
	return s.Spec.GetPlan()
}

func (s *Subscription) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Name:      s.Name,
		Namespace: s.Namespace,
	}
}

func (s *Subscription) GetSpec() core.Spec {
	return &s.Spec
}

func (s *Subscription) GetStatus() core.Status {
	return &s.Status
}

func (s *Subscription) SetApiKind(kind string) {
	s.Spec.SetApiKind(kind)
}

// +kubebuilder:object:root=true
type SubscriptionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Subscription `json:"items"`
}

func (s *Subscription) IsBeingDeleted() bool {
	return !s.ObjectMeta.DeletionTimestamp.IsZero()
}

func init() {
	SchemeBuilder.Register(&Subscription{}, &SubscriptionList{})
}
