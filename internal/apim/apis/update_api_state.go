package apis

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimclientmodel "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client/model"
)

var stateToAction = map[model.State]apimclientmodel.Action{
	model.StateStarted: apimclientmodel.ActionStart,
	model.StateStopped: apimclientmodel.ActionStop,
}

func (d *Delegate) updateApiState(
	apiDefinition *gio.ApiDefinition,
) error {
	// Check if Management context is provided
	if d.apimClient == nil {
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

	apiDefinition.Status.State = apiDefinition.Spec.State
	return nil
}
