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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ core.ResourceObject = &ApiResource{}

// ApiResourceSpec defines the desired state of ApiResource.
// +kubebuilder:object:generate=true
type ApiResourceSpec struct {
	*base.Resource `json:",inline"`
}

// Hash implements core.Spec.
func (spec ApiResourceSpec) Hash() string {
	return hash.Calculate(spec)
}

type ApiResourceStatus struct {
}

// DeepCopyFrom implements core.Status.
func (s *ApiResourceStatus) DeepCopyFrom(obj client.Object) error {
	if res, ok := obj.(*ApiResource); ok {
		res.Status.DeepCopyInto(s)
		return nil
	}
	return fmt.Errorf("unknown type %T", obj)
}

// DeepCopyTo implements core.Status.
func (s *ApiResourceStatus) DeepCopyTo(obj client.Object) error {
	if res, ok := obj.(*ApiResource); ok {
		s.DeepCopyInto(&res.Status)
		return nil
	}
	return fmt.Errorf("unknown type %T", obj)
}

// SetProcessingStatus implements core.Status.
func (s *ApiResourceStatus) SetProcessingStatus(status core.ProcessingStatus) {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

type ApiResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiResourceSpec   `json:"spec,omitempty"`
	Status ApiResourceStatus `json:"status,omitempty"`
}

// GetResourceName implements core.ResourceModel.
func (res *ApiResource) GetResourceName() string {
	return res.Spec.GetResourceName()
}

// GetConfig implements core.ResourceModel.
func (res *ApiResource) GetConfig() *utils.GenericStringMap {
	return res.Spec.GetConfig()
}

// GetType implements core.ResourceModel.
func (res *ApiResource) GetType() string {
	return res.Spec.GetType()
}

func (res *ApiResource) IsBeingDeleted() bool {
	return !res.ObjectMeta.DeletionTimestamp.IsZero()
}

func (res *ApiResource) DeepCopyResource() core.Object {
	return res.DeepCopy()
}

func (res *ApiResource) GetSpec() core.Spec {
	return res.Spec
}

func (res *ApiResource) GetStatus() core.Status {
	return &res.Status
}

func (res *ApiResource) GetRef() core.ObjectRef {
	return &refs.NamespacedName{
		Namespace: res.Namespace,
		Name:      res.Name,
	}
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
