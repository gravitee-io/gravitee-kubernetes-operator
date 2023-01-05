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
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"k8s.io/apimachinery/pkg/util/errors"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (d *Delegate) Delete(
	apiDefinition *gio.ApiDefinition,
) error {
	if !util.ContainsFinalizer(apiDefinition, keys.ApiDefinitionDeletionFinalizer) {
		return nil
	}

	if d.HasContext() {
		return d.deleteWithContexts(apiDefinition)
	}

	util.RemoveFinalizer(apiDefinition, keys.ApiDefinitionDeletionFinalizer)

	return d.k8s.Update(d.ctx, apiDefinition)
}

func (d *Delegate) deleteWithContexts(api *gio.ApiDefinition) error {
	errs := make([]error, 0)

	for _, context := range d.contexts {
		if err := d.deleteWithContext(api, context); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		util.RemoveFinalizer(api, keys.ApiDefinitionDeletionFinalizer)
	}

	errs = append(errs, d.k8s.Update(d.ctx, api))

	return errors.NewAggregate(errs)
}

func (d *Delegate) deleteWithContext(api *gio.ApiDefinition, context DelegateContext) error {
	cp, err := context.compile(api)
	if err != nil {
		return err
	}

	return context.delete(cp)
}
