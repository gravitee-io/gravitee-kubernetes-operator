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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

type ConsumerType int

var _ core.PlanModel = &Plan{}

const (
	TAG ConsumerType = iota
)

type Consumer struct {
	// Consumer type (possible values TAG)
	ConsumerType ConsumerType `json:"consumerType,omitempty"`
	// Consumer ID
	// +kubebuilder:validation:Optional
	ConsumerID *string `json:"consumerId,omitempty"`
}

type Plan struct {
	*base.Plan `json:",inline"`

	// Plan name
	Name string `json:"name"`
	// Plan Description
	Description string `json:"description"`
	// Plan Security
	Security string `json:"security"`
	// Plan Security definition
	// +kubebuilder:validation:Optional
	SecurityDefinition *string `json:"securityDefinition,omitempty"`
	// A map of different paths (alongside their Rules) for this Plan
	// +kubebuilder:validation:Optional
	Paths *map[string][]Rule `json:"paths,omitempty"`
	// Specify the API associated with this plan
	// +kubebuilder:validation:Optional
	Api *string `json:"api,omitempty"`
	// Plan selection rule
	// +kubebuilder:validation:Optional
	SelectionRule *string `json:"selection_rule,omitempty"`
	// List of different flows for this Plan
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Flows []Flow `json:"flows"`
	// List of excluded groups for this plan
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	ExcludedGroups []string `json:"excluded_groups"`
}

func NewPlan(base *base.Plan) *Plan {
	paths := make(map[string][]Rule)
	return &Plan{
		Plan:  base,
		Flows: []Flow{},
		Paths: &paths,
	}
}

func (plan *Plan) GetSecurityType() string {
	return plan.Security
}

func (plan *Plan) WithName(name string) *Plan {
	plan.Name = name
	return plan
}

func (plan *Plan) WithDescription(description string) *Plan {
	plan.Description = description
	return plan
}

func (plan *Plan) WithSecurity(security string) *Plan {
	plan.Security = security
	return plan
}

type Path struct {
	// Path
	// +kubebuilder:validation:Optional
	Path *string `json:"path,omitempty"`
	// Path Rules
	// +kubebuilder:validation:Optional
	Rules []*Rule `json:"rules"`
}
