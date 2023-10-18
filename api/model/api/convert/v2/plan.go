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
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
)

func toPlans(v4Plans []*v4.Plan) []*v2.Plan {
	if len(v4Plans) == 0 {
		return []*v2.Plan{}
	}
	var plans []*v2.Plan
	for _, v4Plan := range v4Plans {
		plans = append(plans, toPlan(v4Plan))
	}
	return plans
}

func toPlan(v4Plan *v4.Plan) *v2.Plan {
	plan := v2.NewPlan(v4Plan.Plan)
	security := v4Plan.Security
	plan.Security = security.Type
	if security.Config != nil {
		if json, err := security.Config.MarshalJSON(); err == nil {
			plan.SecurityDefinition = string(json)
		}
	}
	plan.Tags = v4Plan.Tags
	plan.Flows = toFlows(v4Plan.Flows)
	return plan
}
