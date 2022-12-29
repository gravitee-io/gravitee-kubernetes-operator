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

package internal

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimModel "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/utils"
)

// Add a default keyless plan to the api definition if no plan is defined.
func addDefaultPlan(api *gio.ApiDefinition) {
	plans := api.Spec.Plans

	if len(plans) == 0 {
		api.Spec.Plans = []*model.Plan{
			{
				Name:     defaultPlanName,
				Security: defaultPlanSecurity,
				Status:   defaultPlanStatus,
			},
		}
	}
}

// For each plan, generate a CrossId from Api Id & Plan Name if not defined.
func generateEmptyPlanCrossIds(spec *gio.ApiDefinitionSpec) {
	plans := spec.Plans

	for _, plan := range plans {
		if plan.CrossId == "" {
			plan.CrossId = utils.ToUUID(spec.ID + separator + plan.Name)
		}
	}
}

// Retrieve the plan ids from the management apiEntity.
func retrieveMgmtPlanIds(spec *gio.ApiDefinitionSpec, mgmtApi *apimModel.ApiEntity) {
	plans := spec.Plans

	for _, plan := range plans {
		for _, mgmtPlan := range mgmtApi.Plans {
			if plan.CrossId == mgmtPlan.CrossId {
				plan.Id = mgmtPlan.Id
				plan.Api = mgmtPlan.Api
			}
		}
	}
}
