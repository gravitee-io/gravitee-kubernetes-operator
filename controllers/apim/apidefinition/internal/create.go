// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	managementapierror "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/clienterror"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/utils"
)

func (d *Delegate) Create(
	apiDefinition *gio.ApiDefinition,
) error {
	apiDefinition.Status.CrossID = retrieveCrossId(apiDefinition)
	apiDefinition.Status.State = model.StateStarted // API is considered started by default and updated later if needed

	// Generate new Id or use existing one if is found in Management API
	apiDefinition.Status.ID = utils.NewUUID()
	if d.IsConnectedToManagementApi() {
		api, findApiErr := d.apimClient.GetByCrossId(apiDefinition.Status.CrossID)
		var crossIdNotFoundError *managementapierror.CrossIdNotFoundError

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
	addDefaultPlan(apiDefinition)
	retrievePlansCrossId(apiDefinition)
	apiDefinition.Spec.DefinitionContext = &model.DefinitionContext{
		Origin: origin,
		Mode:   mode,
	}

	if d.IsConnectedToManagementApi() {
		apiJson, marshalErr := json.Marshal(apiDefinition.Spec)
		if marshalErr != nil {
			d.log.Error(marshalErr, "Unable to marshall API definition as JSON")
			return marshalErr
		}
		mgmtApi, mgmtErr := d.apimClient.CreateApi(apiJson)

		if mgmtErr != nil {
			d.log.Error(mgmtErr, "Unable to create API to the Management API")
			return mgmtErr
		}

		d.log.Info("Api has been created to the Management API")

		// Get Plan Id from the Management API to send it to the Gateway. (Used by the Gateway to find subscription)
		retrieveMgmtPlanIds(apiDefinition, mgmtApi)
	}

	err := d.saveConfigMap(apiDefinition)
	if err != nil {
		d.log.Error(err, "Unable to create or update ConfigMap from API definition")
		return err
	}

	err = d.updateApiState(apiDefinition)
	if err != nil {
		d.log.Error(err, "Unable to update api state to the Management API")
		return err
	}

	// Creation succeed, update Status
	apiDefinition.Status.Generation = apiDefinition.ObjectMeta.Generation
	apiDefinition.Status.ProcessingStatus = gio.ProcessingStatusCompleted
	err = d.k8sClient.Status().Update(d.ctx, apiDefinition.DeepCopy())

	if err != nil {
		d.log.Error(err, "Unable to update API definition status")
		return err
	}

	return nil
}
