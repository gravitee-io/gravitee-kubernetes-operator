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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApiContext represents the configuration for a specific environment
// +kubebuilder:object:generate=true
type ApiContextSpec struct {
	Management *model.Management `json:"management,omitempty"`
	Values     map[string]string `json:"values,omitempty"`
}

// ApiContextStatus defines the observed state of an API Context.
type ApiContextStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="BaseUrl",type=string,JSONPath=`.spec.baseUrl`
// +kubebuilder:resource:shortName=graviteecontexts
type ApiContext struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiContextSpec   `json:"spec,omitempty"`
	Status ApiContextStatus `json:"status,omitempty"`
}

func (context *ApiContext) GetNamespacedName() model.NamespacedName {
	return model.NamespacedName{Namespace: context.Namespace, Name: context.Name}
}

// +kubebuilder:object:root=true
// ApiContextList contains a list of API Contexts.
type ApiContextList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiContext `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApiContext{}, &ApiContextList{})
}
