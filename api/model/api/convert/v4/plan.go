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
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
)

func toV4Plans(v2Plans []*v2.Plan) []*v4.Plan {
	var plans []*v4.Plan
	for _, v3Plan := range v2Plans {
		plans = append(plans, toPlan(v3Plan))
	}
	return plans
}

func toPlan(v2Plan *v2.Plan) *v4.Plan {
	plan := v4.NewPlan(v2Plan.Plan)
	plan.Security = toV4Security(v2Plan)
	plan.Flows = toFlows(v2Plan.Flows)
	plan.Mode = v4.StandardPlanMode
	plan.DefinitionVersion = v4.PlanDefinitionVersion
	return plan
}

func toV4Security(v2Plan *v2.Plan) v4.PlanSecurity {
	security := v4.NewPlanSecurity(v2Plan.Security)
	security.Config = toV4SecurityConfig(v2Plan.SecurityDefinition)
	return security
}

func toV4SecurityConfig(securityDefinition string) *utils.GenericStringMap {
	securityConfig := utils.NewGenericStringMap()
	if err := securityConfig.UnmarshalJSON([]byte(securityDefinition)); err != nil {
		return nil
	}
	return securityConfig
}
