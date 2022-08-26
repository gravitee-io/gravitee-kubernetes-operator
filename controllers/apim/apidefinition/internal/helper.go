package internal

import (
	"k8s.io/apimachinery/pkg/types"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	managementapimodel "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/utils"
)

// Return Spec CrossId or generate a new one from api Name & Namespace.
func retrieveCrossId(api *gio.ApiDefinition) string {
	// If a CrossID is defined at the API level, reuse it.
	// If not, just generate a new CrossID
	if api.Spec.CrossId == "" {
		// The ID of the API will be based on the API Name and Namespace to ensure consistency
		return utils.ToUUID(types.NamespacedName{Namespace: api.Namespace, Name: api.Name}.String())
	}

	return api.Spec.CrossId
}

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
func retrievePlansCrossId(api *gio.ApiDefinition) {
	plans := api.Spec.Plans

	for _, plan := range plans {
		if plan.CrossId == "" {
			plan.CrossId = utils.ToUUID(api.Spec.Id + separator + plan.Name)
		}
	}
}

// Retrieve the plan ids from the management apiEntity.
func retrieveMgmtPlanIds(apiDefinition *gio.ApiDefinition, mgmtApi *managementapimodel.ApiEntity) {
	plans := apiDefinition.Spec.Plans

	for _, plan := range plans {
		for _, mgmtPlan := range mgmtApi.Plans {
			if plan.CrossId == mgmtPlan.CrossId {
				plan.Id = mgmtPlan.Id
			}
		}
	}
}
