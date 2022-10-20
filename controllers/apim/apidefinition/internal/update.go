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
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/clienterror"
)

//nolint:gocognit
func (d *Delegate) CreateOrUpdate(
	apiDefinition *gio.ApiDefinition,
) error {
	updating := apiDefinition.Status.CrossID != ""

	apiDefinition.Status.CrossID = getOrGenerateCrossId(apiDefinition)
	apiDefinition.Status.ID = getOrGenerateId(apiDefinition)
	apiDefinition.Spec.Id = apiDefinition.Status.ID
	apiDefinition.Spec.CrossId = apiDefinition.Status.CrossID

	// TODO Check if Management context is provided and don't add default plan if it is the case ?
	addDefaultPlan(apiDefinition)

	// TODO do we really need to do this for plans or should it be done by the management api ?
	generateEmptyPlanCrossIds(apiDefinition)

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

		_, findErr := d.apimClient.GetByCrossId(apiDefinition.Status.CrossID)

		var notFoundError *clienterror.CrossIdNotFoundError
		if findErr != nil && !errors.As(findErr, &notFoundError) {
			return findErr
		}

		importMethod := http.MethodPost
		if findErr == nil {
			importMethod = http.MethodPut
		}

		mgmtApi, mgmtErr := d.apimClient.ImportApi(importMethod, apiJson)
		if mgmtErr != nil {
			d.log.Error(mgmtErr, "Unable to create API to the Management API")
			return mgmtErr
		}

		d.log.Info("Api has been created to the Management API")

		// Get Plan Id from the Management API to send it to the Gateway. (Used by the Gateway to find subscription)
		retrieveMgmtPlanIds(apiDefinition, mgmtApi)
		apiDefinition.Spec.Id = mgmtApi.Id
	}

	if updating || apiDefinition.Spec.State == model.StateStopped {
		if err := d.updateApiState(apiDefinition); err != nil {
			d.log.Error(err, "Unable to update api state to the Management API")
			return err
		}
	}

	if apiDefinition.Spec.State == model.StateStopped {
		if err := d.deleteConfigMap(apiDefinition.Namespace, apiDefinition.Name); err != nil {
			d.log.Error(err, "Unable to delete ConfigMap from API definition")
			return err
		}
	} else {
		if err := d.saveConfigMap(apiDefinition); err != nil {
			d.log.Error(err, "Unable to create or update ConfigMap from API definition")
			return err
		}
	}

	// Creation succeeded, update Status
	apiDefinition.Status.ObservedGeneration = apiDefinition.ObjectMeta.Generation
	apiDefinition.Status.ProcessingStatus = gio.ProcessingStatusCompleted
	apiDefinition.Status.State = apiDefinition.Spec.State
	apiDefinition.Status.ID = apiDefinition.Spec.Id

	if err := d.k8sClient.Status().Update(d.ctx, apiDefinition.DeepCopy()); err != nil {
		d.log.Error(err, "Unable to update API definition status")
		return err
	}

	return nil
}
