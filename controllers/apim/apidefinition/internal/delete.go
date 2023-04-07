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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (d *Delegate) Delete(
	apiDefinition *gio.ApiDefinition,
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

func (d *Delegate) deleteWithContext(api *gio.ApiDefinition) error {
	if err := errors.IgnoreNotFound(d.apim.APIs.Delete(api.Status.ID)); err != nil {
		return err
	}

	context := new(gio.ManagementContext)
	contextRef := api.Spec.Context
	ns := contextRef.ToK8sType()
	d.log.Info("Resolving API context", "namespace", ns.Namespace, "name", ns.Name)
	if err := d.k8s.Get(d.ctx, ns, context); err != nil {
		return err
	}

	if !util.ContainsFinalizer(context, keys.ManagementContextFinalizer) {
		return nil
	}

	apis := &gio.ApiDefinitionList{}
	if err := search.New(d.ctx, d.k8s).FindByFieldReferencing(
		indexer.ContextField,
		*contextRef,
		apis,
	); err != nil {
		return err
	}

	if len(apis.Items) == 1 {
		util.RemoveFinalizer(context, keys.ManagementContextFinalizer)
		if err := d.k8s.Update(d.ctx, context); err != nil {
			return err
		}
	}

	return nil
}
