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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"golang.org/x/net/context"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func Delete(
	ctx context.Context,
	client client.Client,
	instance *gio.ManagementContext,
) error {
	if !util.ContainsFinalizer(instance, keys.ManagementContextFinalizer) {
		return nil
	}

	apis := &gio.ApiDefinitionList{}
	err := search.New(ctx, client).FindByFieldReferencing(
		indexer.ContextField,
		model.NewNamespacedName(instance.Namespace, instance.Name),
		apis,
	)

	if err != nil {
		err = fmt.Errorf("an error occurred while checking if the management context is linked to an api definition: %w", err)
		return err
	}

	if len(apis.Items) > 0 {
		log.FromContext(ctx).Info("context is referenced and will remain", "refCount", len(apis.Items))
		return nil
	}

	apps := &gio.ApplicationList{}
	err = search.New(ctx, client).FindByFieldReferencing(
		indexer.AppContextField,
		model.NewNamespacedName(instance.Namespace, instance.Name),
		apps,
	)

	if err != nil {
		err = fmt.Errorf("an error occurred while checking if the management context is linked to an application: %w", err)
		return err
	}

	if len(apps.Items) > 0 {
		return fmt.Errorf("can not delete %s because %d application(s) are relying on this context",
			instance.Name, len(apps.Items))
	}

	util.RemoveFinalizer(instance, keys.ManagementContextFinalizer)

	return client.Update(ctx, instance)
}
