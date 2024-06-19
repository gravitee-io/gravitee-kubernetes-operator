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

import (
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

type Page struct {
	ID      string      `json:"id,omitempty"`
	Content string      `json:"content,omitempty"`
	Type    string      `json:"type,omitempty"`
	Source  *PageSource `json:"source,omitempty"`
}

type PageSource struct {
	Type          string                  `json:"type,omitempty"`
	Configuration *utils.GenericStringMap `json:"configuration,omitempty"`
}

type PageImport struct {
	*v2.Page          `json:",inline"`
	DefinitionContext *v2.DefinitionContext `json:"definition_context"`
}

type PagesQuery struct {
	Type string
}

func NewPageQuery() *PagesQuery {
	return &PagesQuery{}
}

func (query *PagesQuery) WithType(t string) *PagesQuery {
	query.Type = t
	return query
}

func (query *PagesQuery) AsMap() map[string]string {
	return map[string]string{
		"type": query.Type,
	}
}
