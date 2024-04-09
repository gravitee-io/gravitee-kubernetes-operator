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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (d *Delegate) Delete(apiDefinition client.Object) error {
	if !util.ContainsFinalizer(apiDefinition, keys.ApiDefinitionFinalizer) {
		return nil
	}

	if d.HasContext() {
		if err := d.deleteWithContext(apiDefinition); err != nil {
			return err
		}
	}

	util.RemoveFinalizer(apiDefinition, keys.ApiDefinitionFinalizer)

	return nil
}

func (d *Delegate) deleteWithContext(api client.Object) error {
	switch t := api.(type) {
	case *v1alpha1.ApiDefinition:
		return errors.IgnoreNotFound(d.apim.APIs.DeleteV2(t.Status.ID))
	case *v1alpha1.ApiDefinitionV4:
		return errors.IgnoreNotFound(d.apim.APIs.DeleteV4(t.Status.ID))
	default:
		return fmt.Errorf("unknown type %T", t)
	}
}
