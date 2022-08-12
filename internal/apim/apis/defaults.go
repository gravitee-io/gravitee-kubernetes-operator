package apis

import (
	"k8s.io/apimachinery/pkg/types"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/utils"
)

const separator = "/"

// Return Spec CrossId or generate a new one from api Name & Namespace.
func RetrieveCrossId(api *gio.ApiDefinition) string {
	// If a CrossID is defined at the API level, reuse it.
	// If not, just generate a new CrossID
	if api.Spec.CrossId == "" {
		// The ID of the API will be based on the API Name and Namespace to ensure consistency
		return utils.ToUUID(types.NamespacedName{Namespace: api.Namespace, Name: api.Name}.String())
	}

	return api.Spec.CrossId
}

// Add a default keyless plan to the api definition if no plan is defined.
func (d *Delegate) addDefaultPlan(api *gio.ApiDefinition) {
	plans := api.Spec.Plans

	if len(plans) == 0 {
		d.log.Info("Define default plan for API")
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
func (d *Delegate) retrievePlansCrossId(api *gio.ApiDefinition) {
	plans := api.Spec.Plans

	for _, plan := range plans {
		if plan.CrossId == "" {
			plan.CrossId = utils.ToUUID(api.Spec.Id + separator + plan.Name)
		}
	}
}
