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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
)

func (d *Delegate) CreateOrUpdate(apiDefinition *v1beta1.ApiDefinition) error {
	cp := apiDefinition.DeepCopy()

	spec := &cp.Spec
	spec.ID = cp.PickID()
	spec.CrossID = cp.PickCrossID()
	spec.Plans = cp.PickPlanIDs()
	spec.DefinitionContext = v4.NewDefaultKubernetesContext().MergeWith(spec.DefinitionContext)

	if err := d.resolveResources(cp); err != nil {
		log.Error(d.ctx, err, "Unable to resolve API resources from references")
		return err
	}

	if d.HasContext() {
		log.Info(d.ctx, "Syncing API with APIM")
		status, err := d.apim.APIs.Import(&spec.Api)
		if err != nil {
			return err
		}
		apiDefinition.Status = *status
		log.Debug(d.ctx, "API ID: "+spec.ID)
	}

	if spec.DefinitionContext.SyncFrom == v4.OriginManagement || spec.State == base.StateStopped {
		log.Debug(
			d.ctx,
			"Deleting config map as API is not managed by operator or is stopped",
			"syncFrom", spec.DefinitionContext.SyncFrom,
			"state", spec.State,
		)
		if err := d.deleteConfigMap(cp); err != nil {
			return err
		}
	} else {
		log.Debug(d.ctx, "Saving config map")
		if err := d.saveConfigMap(cp); err != nil {
			return err
		}
	}
	return nil
}
