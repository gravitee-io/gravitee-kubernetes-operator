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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (d *Delegate) Delete(
	apiDefinition *v1beta1.ApiDefinition,
) error {
	if !util.ContainsFinalizer(apiDefinition, keys.ApiDefinitionDeletionFinalizer) {
		return nil
	}

	if d.HasContext() {
		if err := d.deleteWithContext(apiDefinition); err != nil {
			return err
		}
	}

	util.RemoveFinalizer(apiDefinition, keys.ApiDefinitionDeletionFinalizer)

	return d.k8s.Update(d.ctx, apiDefinition)
}

func (d *Delegate) deleteWithContext(api *v1beta1.ApiDefinition) error {
	if api.Status.State == base.StateStarted {
		if err := errors.IgnoreNotFound(d.apim.APIs.Stop(api.Status.ID)); err != nil {
			return err
		}
	}
	return errors.IgnoreNotFound(d.apim.APIs.Delete(api.Status.ID))
}
