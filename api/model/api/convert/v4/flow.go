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
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
)

func toFlowExecution(v2FlowMode v2.FlowMode) *v4.FlowExecution {
	if v2FlowMode == v2.BestMatchFlowMode {
		return &v4.FlowExecution{
			Mode: v4.FlowModeBestMatch,
		}
	}
	return &v4.FlowExecution{
		Mode: v4.FlowModeDefault,
	}
}

func toFlows(v2Flows []v2.Flow) []*v4.Flow {
	if len(v2Flows) == 0 {
		return []*v4.Flow{}
	}
	var flows []*v4.Flow
	for i := range v2Flows {
		flows = append(flows, toFlow(v2Flows[i]))
	}
	return flows
}

func toFlow(v2Flow v2.Flow) *v4.Flow {
	flow := v4.NewFlow(v2Flow.Name)
	flow.Selectors = append(flow.Selectors, toHttpSelector(v2Flow))
	if v2Flow.Condition != "" {
		flow.Selectors = append(flow.Selectors, v4.NewConditionSelector(v2Flow.Condition))
	}
	flow.Request = toFlowSteps(v2Flow.Pre)
	flow.Response = toFlowSteps(v2Flow.Post)
	flow.Enabled = v2Flow.Enabled
	return flow
}

func toFlowSteps(v2Step []base.FlowStep) []*v4.FlowStep {
	if v2Step == nil {
		return []*v4.FlowStep{}
	}
	var flowSteps []*v4.FlowStep
	for i := range v2Step {
		flowSteps = append(flowSteps, v4.NewFlowStep(v2Step[i]))
	}
	return flowSteps
}

func toHttpSelector(v2Flow v2.Flow) *v4.FlowSelector {
	path, operator := getPathAndOperator(v2Flow.PathOperator)
	return v4.NewHTTPSelector(path, operator, v2Flow.Methods)
}

func getPathAndOperator(pathOperator *v2.PathOperator) (string, string) {
	if pathOperator == nil {
		return "", ""
	}
	return pathOperator.Path, pathOperator.Operator
}
