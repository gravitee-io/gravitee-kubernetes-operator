package apis

import (
	"encoding/json"

	apimclientmodel "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client/model"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func (d *Delegate) update(api *gio.ApiDefinition) error {
	d.addPlan(api)
	setIds(api)

	apiJson, err := json.Marshal(api.Spec)
	if err != nil {
		d.log.Error(err, "Unable to marshall API definition as JSON")
		return err
	}

	// Handle Gateway with ConfigMap
	switch {
	case api.Spec.State == model.StateStopped:
		err = d.deleteConfigMap(api.Namespace, api.Name)
		if err != nil {
			d.log.Error(err, "Unable to delete ConfigMap from API definition")
			return err
		}
	default:
		_, err = d.saveConfigMap(api, apiJson)
		if err != nil {
			d.log.Error(err, "Unable to save ConfigMap from API definition")
			return err
		}
	}

	// Handle Management with rest-api Import
	err = d.importToManagementApi(api, apiJson)
	if err != nil {
		d.log.Error(err, "Unable to import to the Management API")
		return err
	}

	if api.Spec.State != "" && d.apimClient != nil {
		stateToAction := map[model.State]apimclientmodel.Action{
			model.StateStarted: apimclientmodel.ActionStart,
			model.StateStopped: apimclientmodel.ActionStop,
		}

		err = d.apimClient.UpdateApiState(api.Status.ID, stateToAction[api.Spec.State])
		if err != nil {
			d.log.Error(err, "Unable to update api state to the Management API")
			return err
		}
	}

	api.Status.Generation = api.ObjectMeta.Generation

	err = d.k8sClient.Status().Update(d.ctx, api)
	if err != nil {
		d.log.Error(err, "Unexpected error while updating API definition status")
		return err
	}

	return nil
}
