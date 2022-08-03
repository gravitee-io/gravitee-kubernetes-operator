package apis

import (
	"encoding/json"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func (d *Delegate) Update(api *gio.ApiDefinition) error {
	d.addPlan(api)
	setIds(api)
	setDeployedAt(api)

	apiJson, err := json.Marshal(api.Spec)
	if err != nil {
		d.log.Error(err, "Unable to marshall API definition as JSON")
		return err
	}

	_, err = d.updateConfigMap(api, apiJson)
	if err != nil {
		d.log.Error(err, "Unable to create or update ConfigMap from API definition")
		return err
	}

	api.Status.Generation = api.ObjectMeta.Generation

	err = d.cli.Status().Update(d.ctx, api)
	if err != nil {
		d.log.Error(err, "Unexpected error while updating API definition status")
		return err
	}

	return nil
}
