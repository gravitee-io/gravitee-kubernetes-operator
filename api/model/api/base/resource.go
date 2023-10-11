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
)

type Resource struct {
	// +kubebuilder:validation:Optional
	Enabled       bool                    `json:"enabled"`
	Name          string                  `json:"name,omitempty"`
	ResourceType  string                  `json:"type,omitempty"`
	Configuration *utils.GenericStringMap `json:"configuration,omitempty"`
}

type ResourceOrRef struct {
	*Resource `json:",omitempty,inline"`
	Ref       *refs.NamespacedName `json:"ref,omitempty"`
}

func (r *ResourceOrRef) IsRef() bool {
	return r.Ref != nil
}
