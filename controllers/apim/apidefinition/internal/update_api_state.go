package internal

import (
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	managementapimodel "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/model"
)

var stateToAction = map[model.State]managementapimodel.Action{
	model.StateStarted: managementapimodel.ActionStart,
	model.StateStopped: managementapimodel.ActionStop,
}

func (d *Delegate) updateApiState(
	apiDefinition *gio.ApiDefinition,
) error {
	// Check if Management context is provided
	if !d.IsConnectedToManagementApi() {
		return nil
	}

	// Do noting if state not change
	if apiDefinition.Spec.State == apiDefinition.Status.State {
		return nil
	}

	err := d.apimClient.UpdateApiState(apiDefinition.Status.ID, stateToAction[apiDefinition.Spec.State])
	if err != nil {
		d.log.Error(err, "Unable to update api state to the Management API")
		return err
	}

	d.log.Info(fmt.Sprintf("API state updated to \"%s\" to the Management API ", apiDefinition.Spec.State))

	apiDefinition.Status.State = apiDefinition.Spec.State
	return nil
}
