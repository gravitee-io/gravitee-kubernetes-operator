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

package model

import "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/dictionary"

type DictionaryDTO struct {
	HRID           string                    `json:"hrid,omitempty"`
	Name           string                    `json:"name"`
	Description    string                    `json:"description,omitempty"`
	Deployed       bool                      `json:"deployed"`
	DictionaryType dictionary.DictionaryType `json:"type"`
	Manual         *ManualSpec               `json:"manual,omitempty"`
	Dynamic        *DynamicSpec              `json:"dynamic,omitempty"`
}

type ManualSpec struct {
	Properties map[string]string `json:"properties"`
}

type DynamicSpec struct {
	Provider *Provider `json:"provider"`
	Trigger  *Trigger  `json:"trigger"`
}

type Provider struct {
	ProviderType   string           `json:"type"`
	URL            string           `json:"url"`
	Method         string           `json:"method"`
	Specification  string           `json:"specification"`
	Body           string           `json:"body,omitempty"`
	UseSystemProxy bool             `json:"useSystemProxy,omitempty"`
	Headers        []ProviderHeader `json:"headers,omitempty" drift:"empty-is-nil"`
}

type ProviderHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Trigger struct {
	Rate int64                  `json:"rate"`
	Unit dictionary.TriggerUnit `json:"unit"`
}

type DictionaryState struct {
	DictionaryDTO     `json:",omitempty"`
	dictionary.Status `json:",omitempty"`
}

func ToDictionaryDTO(crd dictionary.Type, hrid string) DictionaryDTO {
	dto := mapViaJSON[DictionaryDTO](crd)
	dto.HRID = hrid
	return dto
}
