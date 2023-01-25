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
	api.Spec.DefinitionContext = &model.DefinitionContext{
		Origin: origin,
		Mode:   mode,
	}

	if d.HasContext() {
		return d.updateWithContexts(api)
	}

	return d.updateWithoutContext(api)
}

func (d *Delegate) updateWithContexts(api *gio.ApiDefinition) error {
	errs := make([]error, 0)

	for _, context := range d.contexts {
		if err := d.updateWithContext(api, context); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.NewAggregate(errs)
}

func (d *Delegate) updateWithContext(api *gio.ApiDefinition, context DelegateContext) error {
	log := d.log.WithValues("context", context.Location).WithValues("api", api.GetNamespacedName())

	cp, err := context.compile(api)
	if err != nil {
		log.Error(err, "unable to compile api definition")
		return err
	}

	if err = d.ResolveResources(cp); err != nil {
		log.Error(err, "unable to resolve resources")
		return err
	}

	if err = context.update(cp); err != nil {
		log.Error(err, "unable to update api definition")
		return err
	}

	if err = d.updateConfigMap(cp, &context); err != nil {
		log.Error(err, "unable to update config map")
		return err
	}

	statusContext := cp.Status.Contexts[context.Location]

	statusContext.ID = cp.Spec.ID
	statusContext.CrossID = cp.Spec.CrossID
	statusContext.State = cp.Spec.State
	statusContext.Status = gio.ProcessingStatusCompleted

	api.Status.Contexts[context.Location] = statusContext

	return nil
}

func (d *Delegate) updateWithoutContext(api *gio.ApiDefinition) error {
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
