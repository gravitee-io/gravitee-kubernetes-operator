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

package base

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

var _ core.ResourceModel = &Resource{}
var _ core.ObjectOrRef[core.ResourceModel] = &ResourceOrRef{}

type Resource struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	// Is resource enabled or not?
	Enabled bool `json:"enabled"`
	// Resource Name
	// +kubebuilder:validation:Optional
	Name *string `json:"name,omitempty"`
	// Resource Type
	// +kubebuilder:validation:Optional
	Type *string `json:"type,omitempty"`
	// Resource Configuration, arbitrary map of key-values
	// +kubebuilder:validation:Optional
	Configuration *utils.GenericStringMap `json:"configuration,omitempty"`
}

func (r *Resource) GetType() string {
	if r.Type != nil {
		return *r.Type
	}

	return ""
}

func (r *Resource) GetResourceName() string {
	if r.Name != nil {
		return *r.Name
	}

	return ""
}

func (r *Resource) GetConfig() *utils.GenericStringMap {
	return r.Configuration
}

type ResourceOrRef struct {
	// Resource
	*Resource `json:",omitempty,inline"`
	// Reference to a resource
	Ref *refs.NamespacedName `json:"ref,omitempty"`
}

func (r *ResourceOrRef) IsRef() bool {
	return r.Ref != nil
}

func (r *ResourceOrRef) GetRef() core.ObjectRef {
	return r.Ref
}

func (r *ResourceOrRef) GetObject() core.ResourceModel {
	return r.Resource
}

func (r *ResourceOrRef) SetObject(obj core.ResourceModel) {
	if v, ok := obj.(*Resource); ok {
		r.Resource = v
	}
}
