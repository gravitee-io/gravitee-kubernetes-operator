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
	Name string `json:"name,omitempty"`
	// +kubebuilder:default:=true
	Enabled   bool            `json:"enabled"`
	Selectors []*FlowSelector `json:"selectors,omitempty"`
	Request   []*FlowStep     `json:"request,omitempty"`
	Response  []*FlowStep     `json:"response,omitempty"`
	Subscribe []*FlowStep     `json:"subscribe,omitempty"`
	Publish   []*FlowStep     `json:"publish,omitempty"`
	Tags      []string        `json:"tags,omitempty"`
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

type FlowStep struct {
	base.FlowStep    `json:",inline"`
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
	FlowModeDefault   = FlowMode("default")
	FlowModeBestMatch = FlowMode("best-match")
)

type FlowExecution struct {
	Mode          FlowMode `json:"mode,omitempty"`
	MatchRequired bool     `json:"matchRequired"`
}

type SelectorType string

const (
	HTTPSelectorType      = SelectorType("http")
	ChannelSelectorType   = SelectorType("channel")
	ConditionSelectorType = SelectorType("condition")
)

func NewHTTPSelector(path string, operator base.Operator, methods []base.HttpMethod) *FlowSelector {
	selector := utils.NewGenericStringMap()
	selector.Put("path", path)
	selector.Put("pathOperator", operator)
	selector.Put("type", HTTPSelectorType)
	if methods != nil {
		selector.Put("methods", methods)
	}
	return &FlowSelector{
		GenericStringMap: selector,
	}
}

type ChannelOperation string

const (
	SubscribeChannelOperation = ChannelOperation("SUBSCRIBE")
	PublishChannelOperation   = ChannelOperation("PUBLISH")
)

type FlowSelector struct {
	*utils.GenericStringMap `json:"inline,omitempty"`
}

func (l *FlowSelector) UnmarshalJSON(data []byte) error {
	if l.GenericStringMap == nil {
		l.GenericStringMap = utils.NewGenericStringMap()
	}
	return l.GenericStringMap.UnmarshalJSON(data)
}
