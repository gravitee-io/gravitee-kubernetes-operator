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

func toFlowExecution(v3FlowMode v2.FlowMode) *v4.FlowExecution {
	if v3FlowMode == "" {
		return nil
	}
	if v3FlowMode == v2.DefaultFlowMode {
		return &v4.FlowExecution{
			Mode: v4.FlowModeDefault,
		}
	}
	return &v4.FlowExecution{
		Mode: v4.FlowModeBestMatch,
	}
}

func toFlows(v3Flows []v2.Flow) []*v4.Flow {
	if len(v3Flows) == 0 {
		return []*v4.Flow{}
	}
	var flows []*v4.Flow
	for _, v3Flow := range v3Flows {
		flows = append(flows, toFlow(v3Flow))
	}
	return flows
}

func toFlow(v3Flow v2.Flow) *v4.Flow {
	flow := v4.NewFlow(v3Flow.Name)
	flow.Selectors = append(flow.Selectors, toHttpSelector(v3Flow))
	flow.Request = toFlowSteps(v3Flow.Pre)
	flow.Response = toFlowSteps(v3Flow.Post)
	flow.Enabled = v3Flow.Enabled
	return flow
}

func toFlowSteps(v3Step []base.FlowStep) []*v4.FlowStep {
	if v3Step == nil {
		return []*v4.FlowStep{}
	}
	var flowSteps []*v4.FlowStep
	for i := range v3Step {
		flowSteps = append(flowSteps, v4.NewFlowStep(v3Step[i]))
	}
	return flowSteps
}

func toHttpSelector(v3Flow v2.Flow) *v4.FlowSelector {
	path, operator, methods := v3Flow.PathOperator.Path, v3Flow.PathOperator.Operator, v3Flow.Methods
	return v4.NewHTTPSelector(path, operator, methods)
}
