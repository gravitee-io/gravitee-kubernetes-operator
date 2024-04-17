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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimModel "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
)

// For each plan, generate a CrossId from Api Id & Plan Name if not defined.
func generateEmptyPlanCrossIds(spec *v1alpha1.ApiDefinitionV2Spec) {
	plans := spec.Plans

	for _, plan := range plans {
		if plan.CrossId == "" {
			plan.CrossId = uuid.FromStrings(spec.ID, separator, plan.Name)
		}
	}
}

func generatePageIDs(api *v1alpha1.ApiDefinition) {
	spec := &api.Spec
	pages := spec.Pages
	for name, page := range pages {
		page.API = spec.ID
		apiName := api.GetNamespacedName().String()
		if page.CrossID == "" {
			page.CrossID = uuid.FromStrings(apiName, separator, name)
		}
		if page.ID == "" {
			page.ID = uuid.FromStrings(spec.ID, separator, name)
		}
		if page.Parent != "" {
			page.ParentID = uuid.FromStrings(spec.ID, separator, page.Parent)
		}
	}
}

// Retrieve the plan ids from the management apiEntity.
func retrieveMgmtPlanIds(spec *v1alpha1.ApiDefinitionV2Spec, mgmtApi *apimModel.ApiEntity) {
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
