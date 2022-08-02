package apis

import (
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

const (
	requestTimeoutSeconds = 5
	defaultPlanSecurity   = "KEY_LESS"
	defaultPlanStatus     = "PUBLISHED"
	defaultPlanName       = "G.K.O. Default"
	origin                = "kubernetes"
	mode                  = "fully_managed"
)

func (d *Delegate) Create(
	api *gio.ApiDefinition,
) error {
	// Plan is not required from the CRD, but is expected by the Gateway, so we must create at least one
	d.addPlan(api)

	// Ensure that IDs have been generated
	generateIds(api)
	setDeployedAt(api)

	api.Spec.DefinitionContext = &model.DefinitionContext{
		Origin: origin,
		Mode:   mode,
	}

	apiJson, err := json.Marshal(api.Spec)
	if err != nil {
		d.log.Error(err, "Unable to marshall API definition as JSON")
		return err
	}

	updated, err := d.updateConfigMap(api, apiJson)
	if err != nil {
		d.log.Error(err, "Unable to create or update ConfigMap from API definition")
		return err
	}

	if updated {
		err = d.importToManagementApi(api, apiJson)
		if err != nil {
			d.log.Error(err, "Unable to import to the Management API")
			return err
		}
	}

	api.Status.ApiID = api.Spec.CrossId

	err = d.cli.Status().Update(d.ctx, api)
	if err != nil {
		d.log.Error(err, "Unexpected error while updating status")
		return err
	}

	return nil
}

// Add a default keyless plan to the api definition if no plan is defined.
func (d *Delegate) addPlan(api *gio.ApiDefinition) {
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
