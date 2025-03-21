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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/kafka"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:metadata:annotations={"gravitee.io/extends=gateway.networking.k8s.io"}
type KafkaRoute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KafkaRouteSpec   `json:"spec,omitempty"`
	Status KafkaRouteStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type KafkaRouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KafkaRoute `json:"items"`
}

type KafkaRouteSpec struct {
	kafka.KafKaRoute `json:",inline"`
}

type KafkaRouteStatus struct {
	gwAPIv1.RouteStatus `json:",inline"`
}

func init() {
	SchemeBuilder.Register(&KafkaRoute{}, &KafkaRouteList{})
}
