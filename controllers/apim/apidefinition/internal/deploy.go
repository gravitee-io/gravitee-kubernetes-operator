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
	"errors"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

func (d *Delegate) deploy(api *v1alpha1.ApiDefinition) error {
	if api.Spec.IsLocal {
		return d.updateConfigMap(api)
	}

	// Is a not-local and need to deploy it directly on APIM Console, no ConfigMap needed
	if !d.HasContext() {
		return errors.New("a non-local API definition must have a reference to a ManagementContext")
	}

	if err := d.deleteConfigMap(api); err != nil {
		return err
	}

	return d.apim.APIs.Deploy(api.Spec.ID)
}

func (d *Delegate) updateState(api *v1alpha1.ApiDefinition) error {
	if api.Spec.IsLocal {
		return nil
	}

	return d.apim.APIs.UpdateState(api.Spec.ID, model.ApiStateToAction(api.Spec.State))
}
