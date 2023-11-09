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
)

type ConsumerType int

const (
	TAG ConsumerType = iota
)

type Consumer struct {
	ConsumerType ConsumerType `json:"consumerType,omitempty"`
	ConsumerId   string       `json:"consumerId,omitempty"`
}

type Plan struct {
	*base.Plan         `json:",inline"`
	Security           string            `json:"security"`
	SecurityDefinition string            `json:"securityDefinition,omitempty"`
	Paths              map[string][]Rule `json:"paths,omitempty"`
	Api                string            `json:"api,omitempty"`
	SelectionRule      string            `json:"selectionRule,omitempty"`
	Flows              []Flow            `json:"flows,omitempty"`
}

func NewPlan(base *base.Plan) *Plan {
	return &Plan{
		Plan:  base,
		Flows: []Flow{},
		Paths: make(map[string][]Rule),
	}
}

func (plan *Plan) WithSecurity(security string) *Plan {
	plan.Security = security
	return plan
}

type Path struct {
	Path  string  `json:"path,omitempty"`
	Rules []*Rule `json:"rules,omitempty"`
}
