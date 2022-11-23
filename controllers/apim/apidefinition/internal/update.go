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
	"fmt"
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/clienterror"
	"k8s.io/client-go/util/retry"
)

func (d *Delegate) CreateOrUpdate(
	apiDefinition *gio.ApiDefinition,
) error {
	api := apiDefinition.DeepCopy()
	api.Status.CrossID = getOrGenerateCrossId(api)
	api.Status.ID = getOrGenerateId(api)
	api.Spec.Id = api.Status.ID
	api.Spec.CrossId = api.Status.CrossID

	// TODO Check if Management context is provided and don't add default plan if it is the case ?
	addDefaultPlan(api)

	generateEmptyPlanCrossIds(api)

	if err := d.ResolveResources(api); err != nil {
		return err
	}

	api.Spec.DefinitionContext = &model.DefinitionContext{
		Origin: origin,
		Mode:   mode,
	}

	if d.HasManagementContext() {
		apiJson, marshalErr := json.Marshal(api.Spec)
		if marshalErr != nil {
			d.log.Error(marshalErr, "Unable to marshall API definition as JSON")
			return marshalErr
		}

		_, findErr := d.apimClient.GetByCrossId(api.Status.CrossID)

		if findErr != nil && !clienterror.IsNotFound(findErr) {
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

		d.log.Info(fmt.Sprintf("Api has been %s to the Management API", importMethod))

		// Get Plan Id from the Management API to send it to the Gateway. (Used by the Gateway to find subscription)
		retrieveMgmtPlanIds(api, mgmtApi)

		// Make sure status ID will match APIM ID (could be different if APIM generated it)
		api.Spec.Id = mgmtApi.Id
	}

	if api.Spec.State == model.StateStopped {
		if err := d.deleteConfigMap(api.Namespace, api.Name); err != nil {
			d.log.Error(err, "Unable to delete ConfigMap from API definition")
			return err
		}
	} else {
		if err := d.saveConfigMap(api); err != nil {
			d.log.Error(err, "Unable to create or update ConfigMap from API definition")
			return err
		}
	}

	// Creation succeeded, update Status
	status := api.Status.DeepCopy()
	status.ObservedGeneration = api.ObjectMeta.Generation
	status.ProcessingStatus = gio.ProcessingStatusCompleted
	status.State = api.Spec.State
	status.ID = api.Spec.Id

	api.Status = *status

	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		return d.k8sClient.Status().Update(d.ctx, api)
	}); err != nil {
		d.log.Error(err, "Unable to update API definition status")
		return err
	}

	return nil
}
