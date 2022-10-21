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
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	managementapierror "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/clienterror"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (d *Delegate) Delete(
	apiDefinition *gio.ApiDefinition,
) error {
	// Do nothing if finalizer is already removed
	if !util.ContainsFinalizer(apiDefinition, keys.ApiDefinitionDeletionFinalizer) {
		return nil
	}

	if d.HasManagementContext() && apiDefinition.Status.ID != "" {
		d.log.Info("Delete API definition into Management API")
		err := d.apimClient.DeleteApi(apiDefinition.Status.ID)
		if managementapierror.IsNotFound(err) {
			d.log.Info("The API has already been deleted", "id", apiDefinition.Status.ID)
		}
		if err != nil && !managementapierror.IsNotFound(err) {
			d.log.Error(err, "Unable to delete API definition into Management API")
			return err
		}
	}

	// Remove finalizer when API definition is fully deleted
	util.RemoveFinalizer(apiDefinition, keys.ApiDefinitionDeletionFinalizer)

	return d.k8sClient.Update(d.ctx, apiDefinition)
}
