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
	if !d.HasManagementContext() {
		return nil
	}

	// Do noting if state did not change
	if apiDefinition.Spec.State == apiDefinition.Status.State {
		return nil
	}

	err := d.apimClient.UpdateApiState(apiDefinition.Status.ID, stateToAction[apiDefinition.Spec.State])
	if err != nil {
		d.log.Error(err, "Unable to update api state to the Management API")
		return err
	}

	d.log.Info(fmt.Sprintf("API state updated to \"%s\" to the Management API ", apiDefinition.Spec.State))

	return nil
}
