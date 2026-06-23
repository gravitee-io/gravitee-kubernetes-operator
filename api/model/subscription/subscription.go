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

package subscription

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

var _ core.SubscriptionModel = &Type{}

// +kubebuilder:validation:Enum=ACCEPTED;PAUSED;
type Status string

// +kubebuilder:object:generate=true
type ApiKeySpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=32
	// +kubebuilder:validation:MaxLength=256
	Key string `json:"key"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Format:=date-time
	ExpireAt *string `json:"expireAt,omitempty"`
}

type ConsumerConfiguration struct {
	// +kubebuilder:validation:Required
	EntrypointID string `json:"entrypointId"`
	// +kubebuilder:validation:Optional
	Channel string `json:"channel,omitempty"`
	// +kubebuilder:validation:Optional
	EntrypointConfiguration *utils.GenericStringMap `json:"entrypointConfiguration,omitempty"`
}

func (k *ApiKeySpec) GetKey() string {
	return k.Key
}

func (k *ApiKeySpec) GetExpireAt() *string {
	return k.ExpireAt
}

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
	// +kubebuilder:validation:Optional
	Metadata map[string]string `json:"metadata,omitempty"`
	// +kubebuilder:validation:Optional
	ApiKeys []ApiKeySpec `json:"apiKeys,omitempty"`
	// +kubebuilder:validation:Optional
	ConsumerConfiguration *ConsumerConfiguration `json:"consumerConfiguration,omitempty"`
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

func (t *Type) GetEndingAt() *string {
	return t.EndingAt
}

func (t *Type) GetMetadata() map[string]string {
	if t.Metadata == nil {
		return nil
	}
	metadataCopy := make(map[string]string, len(t.Metadata))
	for k, v := range t.Metadata {
		metadataCopy[k] = v
	}
	return metadataCopy
}

func (t *Type) GetApiKeys() []core.ApiKeyModel {
	keys := make([]core.ApiKeyModel, len(t.ApiKeys))
	for i := range t.ApiKeys {
		keys[i] = &t.ApiKeys[i]
	}
	return keys
}

type AutomationApiKeySpec struct {
	Key      string  `json:"key"`
	ExpireAt *string `json:"expireAt,omitempty"`
}

type AutomationSubscription struct {
	HRID                  string                 `json:"hrid"`
	ApplicationHrid       string                 `json:"applicationHrid"`
	PlanHrid              string                 `json:"planHrid"`
	ApiHrid               string                 `json:"apiHrid"`
	Status                string                 `json:"status"`
	StartingAt            string                 `json:"startingAt"`
	EndingAt              string                 `json:"endingAt"`
	Metadata              map[string]string      `json:"metadata,omitempty"`
	ApiKeys               []AutomationApiKeySpec `json:"apiKeys,omitempty"`
	ConsumerConfiguration *ConsumerConfiguration `json:"consumerConfiguration,omitempty"`
}
