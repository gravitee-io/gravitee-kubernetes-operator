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
	// Operator path
	Path string `json:"path,omitempty"`
	// +kubebuilder:default:=STARTS_WITH
	// Operator (possible values STARTS_WITH or EQUALS)
	Operator base.Operator `json:"operator,omitempty"`
}

func NewPathOperator(path string, operator base.Operator) *PathOperator {
	return &PathOperator{
		Path:     path,
		Operator: operator,
	}
}

type Flow struct {
	// Flow ID
	ID string `json:"id,omitempty"`
	// Flow name
	// +kubebuilder:validation:Optional
	Name string `json:"name"`
	// List of path operators
	PathOperator *PathOperator `json:"path-operator,omitempty"`
	// Flow pre step
	// +kubebuilder:validation:Optional
	Pre []base.FlowStep `json:"pre"`
	// Flow post step
	// +kubebuilder:validation:Optional
	Post []base.FlowStep `json:"post"`
	// +kubebuilder:default:=true
	// Indicate if this flow is enabled or disabled
	Enabled bool `json:"enabled"`
	// A list of methods  for this flow (GET;POST;PUT;PATCH;DELETE;OPTIONS;HEAD;CONNECT;TRACE;OTHER)
	// +kubebuilder:validation:Optional
	Methods []base.HttpMethod `json:"methods"`
	// Flow condition
	// +kubebuilder:validation:Optional
	Condition string `json:"condition,omitempty"`
	// List of the consumers of this Flow
	// +kubebuilder:validation:Optional
	Consumers []Consumer `json:"consumers"`
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
	// Policy name
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// Policy configuration is a map of arbitrary key-values
	// +kubebuilder:validation:Optional
	Configuration *utils.GenericStringMap `json:"configuration,omitempty"`
}

type Rule struct {
	// List of http methods for this Rule (GET;POST;PUT;PATCH;DELETE;OPTIONS;HEAD;CONNECT;TRACE;OTHER)
	// +kubebuilder:validation:Optional
	Methods []base.HttpMethod `json:"methods"`
	// Rule policy
	Policy *Policy `json:"policy,omitempty"`
	// Rule description
	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty"`
	// Indicate if the Rule is enabled or not
	// +kubebuilder:validation:Optional
	Enabled bool `json:"enabled,omitempty"`
}
