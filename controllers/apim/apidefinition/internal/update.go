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
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func (d *Delegate) CreateOrUpdate(apiDefinition *gio.ApiDefinition) error {
	cp := apiDefinition.DeepCopy()

	spec := &cp.Spec
	spec.ID = cp.PickID()

	apiDefinition.Status.ID = cp.Spec.ID

	if err := d.resolveResources(spec); err != nil {
		d.log.Error(err, "unable to resolve resources")
		return err
	}

	generateEmptyPlanCrossIds(spec)
	stateUpdated := false
	if d.HasContext() {
		spec.CrossID = cp.PickCrossID()
		stateUpdated = apiDefinition.Status.State != spec.State
		apiDefinition.Status.EnvID = d.apim.EnvID()
		apiDefinition.Status.OrgID = d.apim.OrgID()
		apiDefinition.Status.CrossID = spec.CrossID
		if err := d.updateWithContext(cp); err != nil {
			return err
		}
		apiDefinition.Status.ID = spec.ID
		apiDefinition.Status.State = spec.State
	}

	if err := d.deploy(cp); err != nil {
		return err
	}

	if stateUpdated {
		if err := d.updateState(cp); err != nil {
			return err
		}
	}

	apiDefinition.Status.Status = gio.ProcessingStatusCompleted

	return nil
}

func (d *Delegate) updateWithContext(api *gio.ApiDefinition) error {
	spec := &api.Spec

	spec.SetDefinitionContext()

	_, findErr := d.apim.APIs.GetByCrossID(spec.CrossID)
	if errors.IgnoreNotFound(findErr) != nil {
		return apim.NewContextError(findErr)
	}

	importMethod := http.MethodPost
	if findErr == nil {
		importMethod = http.MethodPut
	}

	mgmtApi, mgmtErr := d.apim.APIs.Import(importMethod, &spec.Api)
	if mgmtErr != nil {
		return apim.NewContextError(mgmtErr)
	}

	retrieveMgmtPlanIds(spec, mgmtApi)
	spec.ID = mgmtApi.ID

	if mgmtApi.ShouldSetKubernetesContext() {
		if err := d.apim.APIs.SetKubernetesContext(mgmtApi.ID); err != nil {
			return apim.NewContextError(err)
		}
	}

	return nil
}
