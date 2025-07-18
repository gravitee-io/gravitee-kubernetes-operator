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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:metadata:annotations={"gravitee.io/extends=gateway.networking.k8s.io"}
type GatewayClassParameters struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GatewayClassParametersSpec   `json:"spec,omitempty"`
	Status GatewayClassParametersStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type GatewayClassParametersList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GatewayClassParameters `json:"items"`
}

// GatewayClassParametersSpec defines the desired state of GatewayClassParameters
// +kubebuilder:object:generate=true
type GatewayClassParametersSpec struct {
	gateway.GatewayClassParameters `json:",inline"`
}

type GatewayClassParametersStatus struct {
	Conditions []metav1.Condition `json:"conditions"`
}

func (params *GatewayClassParameters) GetConditions() map[string]metav1.Condition {
	conditions := make(map[string]metav1.Condition)
	for _, condition := range params.Status.Conditions {
		conditions[condition.Type] = condition
	}
	return conditions
}

func (params *GatewayClassParameters) SetConditions(conditions []metav1.Condition) {
	params.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&GatewayClassParameters{}, &GatewayClassParametersList{})
}
