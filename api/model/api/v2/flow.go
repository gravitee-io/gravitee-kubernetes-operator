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

package v2

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

// +kubebuilder:validation:Enum=DEFAULT;BEST_MATCH;
type FlowMode string

const (
	BestMatchFlowMode = FlowMode("BEST_MATCH")
	DefaultFlowMode   = FlowMode("DEFAULT")
)

type PathOperator struct {
	Path string `json:"path,omitempty"`
	// +kubebuilder:default:=STARTS_WITH
	Operator string `json:"operator,omitempty"`
}

func NewPathOperator(path, operator string) *PathOperator {
	return &PathOperator{
		Path:     path,
		Operator: operator,
	}
}

type Flow struct {
	ID           string          `json:"id,omitempty"`
	Name         string          `json:"name,omitempty"`
	PathOperator *PathOperator   `json:"path-operator,omitempty"`
	Pre          []base.FlowStep `json:"pre,omitempty"`
	Post         []base.FlowStep `json:"post,omitempty"`
	// +kubebuilder:default:=true
	Enabled   bool              `json:"enabled"`
	Methods   []base.HttpMethod `json:"methods,omitempty"`
	Condition string            `json:"condition,omitempty"`
	Consumers []Consumer        `json:"consumers,omitempty"`
}

func NewFlow(name string) Flow {
	return Flow{
		Name:      name,
		Enabled:   true,
		Pre:       []base.FlowStep{},
		Post:      []base.FlowStep{},
		Methods:   []base.HttpMethod{},
		Consumers: []Consumer{},
		Condition: "",
	}
}

type Policy struct {
	Name          string                  `json:"name,omitempty"`
	Configuration *utils.GenericStringMap `json:"configuration,omitempty"`
}

type Rule struct {
	Methods     []base.HttpMethod `json:"methods,omitempty"`
	Policy      *Policy           `json:"policy,omitempty"`
	Description string            `json:"description,omitempty"`
	Enabled     bool              `json:"enabled,omitempty"`
}
