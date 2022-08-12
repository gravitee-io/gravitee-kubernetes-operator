package apis

import (
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func (d *Delegate) update(api *gio.ApiDefinition) error {
	d.addPlan(api)

	setSpecIdsFromStatus(api)

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
		err = d.saveConfigMap(api, apiJson)
		if err != nil {
			d.log.Error(err, "Unable to save ConfigMap from API definition")
			return err
		}
	}

	if d.IsConnectedToManagementApi() {
		err = d.apimClient.UpdateApi(apiJson)
		if err != nil {
			d.log.Error(err, "Unable to update API to the Management API")
			return err
		}

		d.log.Info("Api has been update to the Management API")
	}

	err = d.updateApiState(api)
	if err != nil {
		d.log.Error(err, "Unable to update api state to the Management API")
		return err
	}

	// Updated succeed, update Generation & Status
	api.Status.Generation = api.ObjectMeta.Generation
	err = d.k8sClient.Status().Update(d.ctx, api)
	if err != nil {
		d.log.Error(err, "Unexpected error while updating API definition status")
		return err
	}

	return nil
}
