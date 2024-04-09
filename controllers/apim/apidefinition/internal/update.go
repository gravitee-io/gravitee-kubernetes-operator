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
	"net/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
)

func (d *Delegate) CreateOrUpdate(apiDefinition client.Object) error {
	switch t := apiDefinition.(type) {
	case *v1alpha1.ApiDefinition:
		return d.createOrUpdateV2(t)
	case *v1alpha1.ApiDefinitionV4:
		return d.createOrUpdateV4(t)
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}

func (d *Delegate) createOrUpdateV2(apiDefinition *v1alpha1.ApiDefinition) error {
	cp := apiDefinition.DeepCopy()

	spec := &cp.Spec
	spec.ID = cp.PickID()
	spec.SetDefinitionContext()

	apiDefinition.Status.ID = cp.Spec.ID

	if err := d.resolveResources(spec.Resources); err != nil {
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

	return nil
}

func (d *Delegate) createOrUpdateV4(apiDefinition *v1alpha1.ApiDefinitionV4) error {
	cp := apiDefinition.DeepCopy()

	spec := &cp.Spec
	spec.ID = cp.PickID()
	spec.CrossID = cp.PickCrossID()
	spec.Plans = cp.PickPlanIDs()
	spec.DefinitionContext = v4.NewDefaultKubernetesContext().MergeWith(spec.DefinitionContext)

	if err := d.resolveResources(spec.Resources); err != nil {
		d.log.Error(err, "Unable to resolve API resources from references")
		return err
	}

	if d.HasContext() {
		d.log.Info("Syncing API with APIM")
		status, err := d.apim.APIs.ImportV4(&spec.Api)
		if err != nil {
			return err
		}
		apiDefinition.Status = *status
		d.log.Info("API successfully synced with APIM. ID: " + spec.ID)
	}

	if spec.DefinitionContext.SyncFrom == v4.OriginManagement || spec.State == base.StateStopped {
		d.log.Info(
			"Deleting config map as API is not managed by operator or is stopped",
			"syncFrom", spec.DefinitionContext.SyncFrom,
			"state", spec.State,
		)
		if err := d.deleteConfigMap(cp); err != nil {
			return err
		}
	} else {
		d.log.Info("Saving config map")
		if err := d.saveConfigMap(cp); err != nil {
			return err
		}
	}
	return nil
}

func (d *Delegate) updateWithContext(api *v1alpha1.ApiDefinition) error {
	spec := &api.Spec

	_, findErr := d.apim.APIs.GetByCrossID(spec.CrossID)
	if errors.IgnoreNotFound(findErr) != nil {
		return apim.NewContextError(findErr)
	}

	importMethod := http.MethodPost
	if findErr == nil {
		importMethod = http.MethodPut
	}

	mgmtApi, mgmtErr := d.apim.APIs.ImportV2(importMethod, &spec.Api)
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
