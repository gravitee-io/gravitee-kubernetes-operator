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

// in APIM there are different state (STAGING;PUBLISHED;CLOSED;DEPRECATED;) but in GKO we only have a single state
// +kubebuilder:validation:Enum=PUBLISHED;
type PlanStatus string

const (
	PublishedPlanStatus = PlanStatus("PUBLISHED")
)

// +kubebuilder:validation:Enum=AUTO;MANUAL;
type PlanValidation string

type Plan struct {
	// Plan ID
	ID string `json:"id,omitempty"`
	// The plan Cross ID.
	// This field is used to identify plans defined for an API
	// that has been promoted between different environments.
	CrossID string `json:"crossId,omitempty"`
	// List of plan tags
	// +kubebuilder:validation:Optional
	Tags []string `json:"tags"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=PUBLISHED
	// The plan status
	Status PlanStatus `json:"status,omitempty"`
	// List of plan characteristics
	// +kubebuilder:validation:Optional
	Characteristics []string `json:"characteristics"`
	// +kubebuilder:default:=AUTO
	// Plan validation strategy
	Validation PlanValidation `json:"validation,omitempty"`
	// Indicate of comment is required for this plan or not
	// +kubebuilder:validation:Optional
	CommentRequired *bool `json:"comment_required,omitempty"`
	// Plan order
	// +kubebuilder:validation:Optional
	Order *int `json:"order,omitempty"`
	// +kubebuilder:default:=API
	// Plan type
	Type PlanType `json:"type,omitempty"`
}

func NewPlan() *Plan {
	return &Plan{
		Tags:            []string{},
		Characteristics: []string{},
	}
}

func (plan *Plan) WithStatus(status PlanStatus) *Plan {
	plan.Status = status
	return plan
}

func (plan *Plan) GetID() string {
	return plan.ID
}

func (plan *Plan) WithID(id string) *Plan {
	plan.ID = id
	return plan
}

func (plan *Plan) WithCrossID(id string) *Plan {
	plan.CrossID = id
	return plan
}
