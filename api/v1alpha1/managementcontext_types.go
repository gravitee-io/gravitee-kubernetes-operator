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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/management"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ManagementContext represents the configuration for a specific environment
// +kubebuilder:object:generate=true
type ManagementContextSpec struct {
	*management.Context `json:",inline"`
}

// ManagementContextStatus defines the observed state of an API Context.
type ManagementContextStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="BaseUrl",type=string,JSONPath=`.spec.baseUrl`
// +kubebuilder:resource:shortName=graviteecontexts
type ManagementContext struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagementContextSpec   `json:"spec,omitempty"`
	Status ManagementContextStatus `json:"status,omitempty"`
}

func (ctx *ManagementContext) GetNamespacedName() *refs.NamespacedName {
	return &refs.NamespacedName{Namespace: ctx.Namespace, Name: ctx.Name}
}

// +kubebuilder:object:root=true
// ManagementContextList contains a list of API Contexts.
type ManagementContextList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagementContext `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ManagementContext{}, &ManagementContextList{})
}

func (ctx *ManagementContext) IsBeingDeleted() bool {
	return !ctx.ObjectMeta.DeletionTimestamp.IsZero()
}
