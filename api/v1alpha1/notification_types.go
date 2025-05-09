/*
 * Copyright (C) 2015 The Gravitee team (http://gravitee.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/notification"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ConsoleNotificationSettingsObject = &base.ConsoleNotificationConfiguration{}

// Notification defines notification settings in Gravitee
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Target",type=string,JSONPath=`.spec.target`,description="Target"
// +kubebuilder:printcolumn:name="Event Type",type=string,JSONPath=`.spec.eventType`,description="Event Type"
// +kubebuilder:resource:shortName=graviteenotif
type Notification struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NotificationSpec   `json:"spec,omitempty"`
	Status NotificationStatus `json:"status,omitempty"`
}

func (res *Notification) IsBeingDeleted() bool {
	return !res.ObjectMeta.DeletionTimestamp.IsZero()
}

func (res *Notification) GetSpec() core.Spec {
	return &res.Spec
}

func (res *Notification) GetStatus() core.Status {
	return &res.Status
}

func (res *Notification) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Namespace: res.Namespace,
		Name:      res.Name,
	}
}

//+kubebuilder:object:root=true

type NotificationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Notification `json:"items"`
}

// NotificationSpec defines the desired state of a Notification.
// It is to be referenced in an API.
// +kubebuilder:object:generate=true
type NotificationSpec struct {
	*notification.Type `json:",inline"`
}

// Hash implements core.Spec.
func (spec *NotificationSpec) Hash() string {
	return hash.Calculate(spec)
}

// NotificationStatus defines the observed state of the Notification.
type NotificationStatus struct {
	// Conditions are the condition that must be met by the Notification
	// "Accepted" condition is used to indicate if the `Notification` can be used by another resource.
	// "ResolveRef" condition is used to indicate if an error occurred while resolving console groups.
	Conditions *[]metav1.Condition `json:"conditions"`
}

// DeepCopyFrom implements core.Status.
func (s *NotificationStatus) DeepCopyFrom(obj client.Object) error {
	if res, ok := obj.(*Notification); ok {
		res.Status.DeepCopyInto(s)
		return nil
	}
	return fmt.Errorf("unknown type %T", obj)
}

// DeepCopyTo implements core.Status.
func (s *NotificationStatus) DeepCopyTo(obj client.Object) error {
	if res, ok := obj.(*Notification); ok {
		s.DeepCopyInto(&res.Status)
		return nil
	}
	return fmt.Errorf("unknown type %T", obj)
}

// SetProcessingStatus implements core.Status.
func (s *NotificationStatus) SetProcessingStatus(core.ProcessingStatus) {
	// unused
}

func (s *NotificationStatus) IsFailed() bool {
	if s.Conditions != nil {
		for _, condition := range *s.Conditions {
			if condition.Status == metav1.ConditionFalse {
				return true
			}
		}
	}
	return false
}

func init() {
	SchemeBuilder.Register(&Notification{}, &NotificationList{})
}
