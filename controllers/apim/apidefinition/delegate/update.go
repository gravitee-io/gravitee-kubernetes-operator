package delegate

import (
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func (d *Delegate) Update(apiDefinition *gio.ApiDefinition) error {
	// Add required fields to the API definition spec
	// ⚠️ This filed should not be added in ApiDefinition resource
	apiDefinition.Spec.CrossId = apiDefinition.Status.CrossID
	apiDefinition.Spec.Id = apiDefinition.Status.ID
	addDefaultPlan(apiDefinition)
	retrievePlansCrossId(apiDefinition)

	if d.IsConnectedToManagementApi() {
		apiJson, err := json.Marshal(apiDefinition.Spec)
		if err != nil {
			d.log.Error(err, "Unable to marshall API definition as JSON")
			return err
		}

		mgmtApi, err := d.apimClient.UpdateApi(apiJson)
		if err != nil {
			d.log.Error(err, "Unable to update API to the Management API")
			return err
		}

		d.log.Info("Api has been update to the Management API")

		// Get Plan Id from the Management API to send it to the Gateway. (Used by the Gateway to find subscription)
		retrieveMgmtPlanIds(apiDefinition, mgmtApi)
	}

	// Handle Gateway with ConfigMap
	// Delete ConfigMap if api is stopped or save it
	switch {
	case apiDefinition.Spec.State == model.StateStopped:
		err := d.deleteConfigMap(apiDefinition.Namespace, apiDefinition.Name)
		if err != nil {
			d.log.Error(err, "Unable to delete ConfigMap from API definition")
			return err
		}
	default:
		err := d.saveConfigMap(apiDefinition)
		if err != nil {
			d.log.Error(err, "Unable to save ConfigMap from API definition")
			return err
		}
	}

	err := d.updateApiState(apiDefinition)
	if err != nil {
		d.log.Error(err, "Unable to update api state to the Management API")
		return err
	}

	// Updated succeed, update Generation & Status
	apiDefinition.Status.Generation = apiDefinition.ObjectMeta.Generation
	err = d.k8sClient.Status().Update(d.ctx, apiDefinition.DeepCopy())
	if err != nil {
		d.log.Error(err, "Unexpected error while updating API definition status")
		return err
	}

	return nil
}
