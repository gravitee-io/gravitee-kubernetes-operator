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

import "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"

// +kubebuilder:validation:Enum=STARTS_WITH;EQUALS;
type Operator string

const (
	StartWithOperator = Operator("STARTS_WITH")
	EqualsOperator    = Operator("EQUALS")
)

type FlowStep struct {
	// +kubebuilder:default:=true
	// Indicate if this FlowStep is enabled or not
	Enabled bool `json:"enabled"`
	// FlowStep policy
	Policy string `json:"policy,omitempty"`
	// FlowStep name
	Name string `json:"name,omitempty"`
	// FlowStep description
	Description string `json:"description,omitempty"`
	// FlowStep configuration is a map of arbitrary key-values
	Configuration *utils.GenericStringMap `json:"configuration,omitempty"`
	// FlowStep condition
	Condition string `json:"condition,omitempty"`
}
