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

type Flow struct {
	// Flow name
	Name string `json:"name,omitempty"`

	// +kubebuilder:default:=true
	// Is flow enabled or not?
	Enabled bool `json:"enabled"`

	// List of Flow selectors
	Selectors []*FlowSelector `json:"selectors,omitempty"`

	// List of Request flow steps
	Request []*FlowStep `json:"request,omitempty"`

	// List of Response flow steps
	Response []*FlowStep `json:"response,omitempty"`

	// List of Subscribe flow steps
	Subscribe []*FlowStep `json:"subscribe,omitempty"`

	// List of Publish flow steps
	Publish []*FlowStep `json:"publish,omitempty"`

	// List of tags
	Tags []string `json:"tags,omitempty"`
}

func NewFlow(name string) *Flow {
	return &Flow{
		Name:      name,
		Enabled:   true,
		Selectors: []*FlowSelector{},
		Request:   []*FlowStep{},
		Response:  []*FlowStep{},
		Subscribe: []*FlowStep{},
		Publish:   []*FlowStep{},
		Tags:      []string{},
	}
}

func (fl Flow) ToGatewayDefinition() *Flow {
	for i := range fl.Selectors {
		fl.Selectors[i] = fl.Selectors[i].ToGatewayDefinition()
	}
	return &fl
}

type FlowStep struct {
	base.FlowStep `json:",inline"`
	// The message condition (supports EL expressions)
	MessageCondition string `json:"messageCondition,omitempty"`
}

func NewFlowStep(base base.FlowStep) *FlowStep {
	return &FlowStep{
		FlowStep: base,
	}
}

func (step *FlowStep) WithMessageCondition(messageCondition string) *FlowStep {
	step.MessageCondition = messageCondition
	return step
}

type FlowMode string

const (
	FlowModeDefault   = FlowMode("DEFAULT")
	FlowModeBestMatch = FlowMode("BEST_MATCH")
)

type FlowExecution struct {
	// The flow mode to use
	Mode FlowMode `json:"mode,omitempty"`
	// Is match required or not ? If set to true, a 404 status response will be returned if no matching flow was found.
	MatchRequired bool `json:"matchRequired"`
}

func DefaultFlowExecution() *FlowExecution {
	return &FlowExecution{
		Mode: FlowModeDefault,
	}
}

const (
	HTTPSelectorType      = "HTTP"
	ChannelSelectorType   = "CHANNEL"
	ConditionSelectorType = "CONDITION"
)

func NewHTTPSelector(path, operator string, methods []base.HttpMethod) *FlowSelector {
	impl := utils.NewGenericStringMap()
	impl.Put("type", HTTPSelectorType)
	impl.Put("path", path)
	impl.Put("pathOperator", operator)
	if methods != nil {
		impl.Put("methods", methods)
	}
	return &FlowSelector{
		GenericStringMap: impl,
	}
}

func NewConditionSelector(condition string) *FlowSelector {
	impl := utils.NewGenericStringMap()
	impl.Put("type", ConditionSelectorType)
	impl.Put("condition", condition)
	return &FlowSelector{
		GenericStringMap: impl,
	}
}

type ChannelOperation string

const (
	SubscribeChannelOperation = ChannelOperation("SUBSCRIBE")
	PublishChannelOperation   = ChannelOperation("PUBLISH")
)

type FlowSelector struct {
	*utils.GenericStringMap `json:",inline"`
}

func (fls FlowSelector) ToGatewayDefinition() *FlowSelector {
	fls.Put("type", Enum(fls.GetString("type")).ToGatewayDefinition())
	return &fls
}

func (fls *FlowSelector) UnmarshalJSON(data []byte) error {
	if fls.GenericStringMap == nil {
		fls.GenericStringMap = utils.NewGenericStringMap()
	}
	return fls.GenericStringMap.UnmarshalJSON(data)
}
