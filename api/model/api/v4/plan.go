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
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

// +kubebuilder:validation:Enum=STANDARD;PUSH;
type PlanMode string

type DefinitionVersion string

var _ core.PlanModel = &Plan{}

const (
	PlanDefinitionVersion DefinitionVersion = "V4"
)

type Plan struct {
	*base.Plan `json:",inline"`

	// Plan display name, this will be the name displayed in the UI
	// if a management context is used to sync the API with APIM
	Name string `json:"name"`

	// Plan Description
	// +kubebuilder:validation:Optional
	Description *string `json:"description,omitempty"`

	// Plan definition version
	// +kubebuilder:default:=`V4`
	DefinitionVersion DefinitionVersion `json:"definitionVersion,omitempty"`

	// Plan security
	Security *PlanSecurity `json:"security,omitempty"`

	// The plan mode
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=`STANDARD`
	// +kubebuilder:validation:Enum=STANDARD;PUSH;
	Mode PlanMode `json:"mode,omitempty"`

	// Plan selection rule
	// +kubebuilder:validation:Optional
	SelectionRule *string `json:"selectionRule,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	// List of plan flows
	Flows []*Flow `json:"flows"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	ExcludedGroups []string `json:"excludedGroups"`
	// The general conditions defined to use this plan
	// +kubebuilder:validation:Optional
	GeneralConditions *string `json:"generalConditions,omitempty"`
}

func (plan *Plan) MarshalJSON() ([]byte, error) {
	type Alias Plan
	aux := struct {
		*Alias
	}{
		Alias: (*Alias)(plan),
	}

	if plan.Plan == nil {
		plan.Plan = &base.Plan{}
	}

	// Set default if zero value
	if plan.Plan.Status == nil {
		defaultStatus := base.PublishedPlanStatus
		plan.Status = &defaultStatus
	}
	if plan.Plan.Validation == nil {
		defaultValidation := base.PlanValidation("AUTO")
		plan.Validation = &defaultValidation
	}
	if plan.Plan.Type == nil {
		defaultPlanType := base.PlanType("API")
		plan.Type = &defaultPlanType
	}

	return json.Marshal(aux)
}

type GatewayDefinitionPlan struct {
	*Plan `json:",inline"`
	Name  string `json:"name"`
}

func (plan *Plan) GetSecurityType() string {
	return plan.Security.Type
}

func NewPlan() *Plan {
	return &Plan{
		Plan: base.NewPlan(),
	}
}

func (plan *Plan) WithSecurity(security *PlanSecurity) *Plan {
	plan.Security = security
	return plan
}

func (plan *Plan) ToGatewayDefinition(name string) *GatewayDefinitionPlan {
	def := &GatewayDefinitionPlan{Plan: plan, Name: name}
	if plan.Security != nil {
		def.Security.Type = Enum(plan.Security.Type).ToGatewayDefinition()
	}
	def.Mode = PlanMode(Enum(plan.Mode).ToGatewayDefinition())
	planStatus := base.PlanStatus(Enum(*plan.Status).ToGatewayDefinition())
	def.Status = &planStatus
	flows := make([]*Flow, len(plan.Flows))
	for i := range def.Flows {
		flows[i] = def.Flows[i].ToGatewayDefinition()
	}
	def.Flows = flows
	return def
}

type PlanSecurity struct {
	// +kubebuilder:validation:Required
	// Plan Security type
	Type string `json:"type"`

	// Plan security configuration, a map of arbitrary key-values
	// +kubebuilder:validation:Optional
	Config *utils.GenericStringMap `json:"configuration,omitempty"`
}

func NewPlanSecurity(kind string) PlanSecurity {
	return PlanSecurity{
		Type:   kind,
		Config: utils.NewGenericStringMap(),
	}
}
