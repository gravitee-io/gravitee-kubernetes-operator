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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApiResourceSpec defines the desired state of ApiResource.
// +kubebuilder:object:generate=true
type ApiResourceSpec struct {
	*base.Resource `json:",inline"`
}

type ApiResourceStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

type ApiResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiResourceSpec   `json:"spec,omitempty"`
	Status ApiResourceStatus `json:"status,omitempty"`
}

func (res *ApiResource) IsBeingDeleted() bool {
	return !res.ObjectMeta.DeletionTimestamp.IsZero()
}

//+kubebuilder:object:root=true

type ApiResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApiResource{}, &ApiResourceList{})
}
