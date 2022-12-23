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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/errors"
)

func (d *Delegate) CreateOrUpdate(api *gio.ApiDefinition) error {
	addDefaultPlan(api)

	api.Spec.DefinitionContext = &model.DefinitionContext{
		Origin: origin,
		Mode:   mode,
	}

	if d.HasContext() {
		return d.UpdateWithContext(api)
	}

	return d.UpdateWithoutContext(api)
}

func (d *Delegate) UpdateWithContext(api *gio.ApiDefinition) error {
	errs := make([]error, 0)

	for i, context := range d.contexts {
		cp, err := context.compile(api)

		if err != nil {
			errs = append(errs, err)
		}

		if err = d.ResolveResources(cp); err != nil {
			errs = append(errs, err)
		}

		if err = context.update(cp); err != nil {
			errs = append(errs, err)
		}

		if err = d.updateConfigMap(cp, &d.contexts[i]); err != nil {
			errs = append(errs, err)
		}

		api.Status.Contexts[context.Location] = cp.Status.Contexts[context.Location]
	}

	return errors.NewAggregate(errs)
}

func (d *Delegate) UpdateWithoutContext(api *gio.ApiDefinition) error {
	spec := &api.Spec

	generateEmptyPlanCrossIds(spec)

	errs := make([]error, 0)

	if err := d.ResolveResources(api); err != nil {
		errs = append(errs, err)
	}

	if err := d.updateConfigMap(api, nil); err != nil {
		errs = append(errs, err)
	}

	return errors.NewAggregate(errs)
}
