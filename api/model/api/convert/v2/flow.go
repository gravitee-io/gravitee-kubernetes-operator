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
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
)

func toFlows(v4Flows []*v4.Flow) []v2.Flow {
	if len(v4Flows) == 0 {
		return []v2.Flow{}
	}
	var flows []v2.Flow
	for _, v4Flow := range v4Flows {
		flows = append(flows, toFlow(v4Flow))
	}
	return flows
}

func toFlow(v4Flow *v4.Flow) v2.Flow {
	flow := v2.NewFlow(v4Flow.Name)
	flow.Enabled = v4Flow.Enabled
	flow.Pre = toFlowSteps(v4Flow.Request)
	flow.Post = toFlowSteps(v4Flow.Response)
	flow.PathOperator = toPathOperator(v4Flow.Selectors)
	return flow
}

func toPathOperator(selectors []*v4.FlowSelector) *v2.PathOperator {
	if len(selectors) == 0 {
		return nil
	}
	selector := selectors[0]
	if selector.GetString("type") != "http" {
		return nil
	}

	path := selector.GetString("path")
	operator := selector.GetString("pathOperator")
	return v2.NewPathOperator(path, base.Operator(operator))
}

func toFlowSteps(v4Steps []*v4.FlowStep) []base.FlowStep {
	if len(v4Steps) == 0 {
		return []base.FlowStep{}
	}
	var flowSteps []base.FlowStep
	for _, v4Step := range v4Steps {
		flowSteps = append(flowSteps, v4Step.FlowStep)
	}
	return flowSteps
}
