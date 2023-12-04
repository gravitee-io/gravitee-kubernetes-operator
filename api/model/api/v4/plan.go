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

package v4

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

type PlanMode string

const (
	StandardPlanMode = PlanMode("STANDARD")
	PushPlanMode     = PlanMode("PUSH")
)

type definitionVersion string

const (
	PlanDefinitionVersion = definitionVersion("V4")
)

type Plan struct {
	*base.Plan        `json:",inline"`
	DefinitionVersion definitionVersion `json:"definitionVersion,omitempty"`
	Security          PlanSecurity      `json:"security,omitempty"`
	Mode              PlanMode          `json:"mode,omitempty"`
	SelectionRule     string            `json:"selectionRule,omitempty"`
	Flows             []*Flow           `json:"flows,omitempty"`
}

type GatewayDefinitionPlan struct {
	*Plan `json:",inline"`
	Name  string `json:"name"`
}

func NewPlan(base *base.Plan) *Plan {
	return &Plan{
		Plan:              base,
		DefinitionVersion: PlanDefinitionVersion,
		Mode:              StandardPlanMode,
		Flows:             []*Flow{},
	}
}

func (plan *Plan) WithSecurity(security PlanSecurity) *Plan {
	plan.Security = security
	return plan
}

func (plan *Plan) ToGatewayDefinition(name string) *GatewayDefinitionPlan {
	def := &GatewayDefinitionPlan{Plan: plan, Name: name}
	def.Security.Type = Enum(plan.Security.Type).ToGatewayDefinition()
	def.Mode = PlanMode(Enum(plan.Mode).ToGatewayDefinition())
	def.Status = base.PlanStatus(Enum(plan.Status).ToGatewayDefinition())
	for i := range def.Flows {
		def.Flows[i] = def.Flows[i].ToGatewayDefinition()
	}
	return def
}

type PlanSecurity struct {
	// +kubebuilder:validation:Required
	Type   string                  `json:"type"`
	Config *utils.GenericStringMap `json:"configuration,omitempty"`
}

func NewPlanSecurity(kind string) PlanSecurity {
	return PlanSecurity{
		Type:   kind,
		Config: utils.NewGenericStringMap(),
	}
}