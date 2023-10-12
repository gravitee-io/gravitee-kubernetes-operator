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

package v1beta1

import (
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApiDefinitionSpec defines the desired state of ApiDefinition.
type ApiDefinitionSpec struct {
	v4.Api `json:",inline"`
	// We don't add the context here because APIM is not ready for that.
}

// ApiDefinitionStatus defines the observed state of ApiDefinition.
type ApiDefinitionStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ApiDefinition is the Schema for the apidefinitions API.
// The v1beta1 API version is compatible with APIM 4.x features.
type ApiDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiDefinitionSpec   `json:"spec,omitempty"`
	Status ApiDefinitionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApiDefinitionList contains a list of ApiDefinition.
type ApiDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApiDefinition{}, &ApiDefinitionList{})
}
