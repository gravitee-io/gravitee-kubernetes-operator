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

// +kubebuilder:object:generate=true
package subscription

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

var _ core.SubscriptionModel = &Type{}

// +kubebuilder:validation:Enum=ACCEPTED;PAUSED;
type Status string

type Type struct {
	// +kubebuilder:validation:Required
	API refs.NamespacedName `json:"api"`
	// +kubebuilder:validation:Required
	App refs.NamespacedName `json:"application"`
	// +kubebuilder:validation:Required
	Plan string `json:"plan"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Format:=date-time
	EndingAt *string `json:"endingAt,omitempty"`
}

type ApiRef struct {
	refs.NamespacedName `json:",inline"`
	// +kubebuilder:default:=ApiV4Definition
	Kind string `json:"kind,omitempty"`
}

func (t *Type) GetApiRef() core.ObjectRef {
	return &t.API
}

func (t *Type) GetAppRef() core.ObjectRef {
	return &t.App
}

func (t *Type) GetPlan() string {
	return t.Plan
}

func (t *Type) SetApiKind(kind string) {
	t.API.Kind = kind
}
