package apis

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client/clienterror"
)

const (
	requestTimeoutSeconds = 5
	defaultPlanSecurity   = "KEY_LESS"
	defaultPlanStatus     = "PUBLISHED"
	defaultPlanName       = "G.K.O. Default"
	origin                = "kubernetes"
	mode                  = "fully_managed"
)

func (d *Delegate) create(
	apiDefinition *gio.ApiDefinition,
) error {
	apiDefinition.Status.CrossID = RetrieveCrossId(apiDefinition)
	apiDefinition.Status.State = model.StateStarted // API is considered started by default and updated later if needed

	// Generate new Id or use existing one if is found in Management API
	apiDefinition.Status.ID = generateId()
	if d.IsConnectedToManagementApi() {
		api, findApiErr := d.apimClient.GetByCrossId(apiDefinition.Status.CrossID)
		var crossIdNotFoundError *clienterror.CrossIdNotFoundError

		switch {
		case findApiErr != nil && errors.As(findApiErr, &crossIdNotFoundError):
			// Do nothing. API is just not existing in the Management API
		case findApiErr != nil:
			d.log.Error(findApiErr, "Error while trying to find API in the Management API")
			return findApiErr
		default:
			// Api found in the Management API
			// Update status with the found ID to trigger new reconcile to update the existing API
			apiDefinition.Status.ID = api.Id
			d.log.Info(fmt.Sprintf("API \"%s\" found in the Management API. Continue with update process", api.Name),
				"id", api.Id, "crossId", apiDefinition.Status.CrossID, "name", api.Name)
			return d.k8sClient.Status().Update(d.ctx, apiDefinition.DeepCopy())
		}
	}

	// Add required fields to the API definition spec
	// ⚠️ This filed should not be added in ApiDefinition resource
	apiDefinition.Spec.Id = apiDefinition.Status.ID
	apiDefinition.Spec.CrossId = apiDefinition.Status.CrossID
	// Plan is not required from the CRD, but is expected by the Gateway, so we must create at least one
	d.addDefaultPlan(apiDefinition)
	d.retrievePlansCrossId(apiDefinition)
	apiDefinition.Spec.DefinitionContext = &model.DefinitionContext{
		Origin: origin,
		Mode:   mode,
	}

	apiJson, err := json.Marshal(apiDefinition.Spec)
	if err != nil {
		d.log.Error(err, "Unable to marshall API definition as JSON")
		return err
	}

	err = d.saveConfigMap(apiDefinition, apiJson)
	if err != nil {
		d.log.Error(err, "Unable to create or update ConfigMap from API definition")
		return err
	}

	if d.IsConnectedToManagementApi() {
		err = d.apimClient.CreateApi(apiJson)
		if err != nil {
			d.log.Error(err, "Unable to create API to the Management API")
			return err
		}

		d.log.Info("Api has been created to the Management API")
	}

	err = d.updateApiState(apiDefinition)
	if err != nil {
		d.log.Error(err, "Unable to update api state to the Management API")
		return err
	}

	// Creation succeed, update Generation & Status
	apiDefinition.Status.Generation = apiDefinition.ObjectMeta.Generation
	err = d.k8sClient.Status().Update(d.ctx, apiDefinition.DeepCopy())

	if err != nil {
		d.log.Error(err, "Unexpected error while updating API definition status")
		return err
	}

	return nil
}
