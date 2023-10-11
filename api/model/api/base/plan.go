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

// +kubebuilder:validation:Enum=API;CATALOG;
type PlanType string

// +kubebuilder:validation:Enum=STAGING;PUBLISHED;CLOSED;DEPRECATED;
type PlanStatus string

const (
	StagingPlanStatus   = PlanStatus("STAGING")
	PublishedPlanStatus = PlanStatus("PUBLISHED")
	ClosedPlanStatus    = PlanStatus("CLOSED")
)

// +kubebuilder:validation:Enum=AUTO;MANUAL;
type PlanValidation string

type Plan struct {
	Id          string   `json:"id,omitempty"`
	CrossId     string   `json:"crossId,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags,omitempty"`
	// +kubebuilder:default:=PUBLISHED
	Status          PlanStatus `json:"status,omitempty"`
	Characteristics []string   `json:"characteristics,omitempty"`
	// +kubebuilder:default:=AUTO
	Validation      PlanValidation `json:"validation,omitempty"`
	CommentRequired bool           `json:"comment_required,omitempty"`
	Order           int            `json:"order,omitempty"`
	// +kubebuilder:default:=API
	Type           PlanType `json:"type,omitempty"`
	ExcludedGroups []string `json:"excluded_groups,omitempty"`
}

func NewPlan(name, description string) *Plan {
	return &Plan{
		Name:            name,
		Description:     description,
		Tags:            []string{},
		Characteristics: []string{},
		ExcludedGroups:  []string{},
	}
}

func (plan *Plan) WithStatus(status PlanStatus) *Plan {
	plan.Status = status
	return plan
}

func (plan *Plan) WithID(id string) *Plan {
	plan.Id = id
	return plan
}

func (plan *Plan) WithCrossID(id string) *Plan {
	plan.CrossId = id
	return plan
}
