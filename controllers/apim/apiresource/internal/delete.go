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
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/indexer"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Delete(
	ctx context.Context,
	client client.Client,
	resource *v1alpha1.ApiResource,
) error {
	if !util.ContainsFinalizer(resource, keys.ApiResourceFinalizer) {
		return nil
	}

	search := search.New(ctx, client)

	apis := &v1alpha1.ApiDefinitionList{}
	if err := search.FindByFieldReferencing(
		indexer.ResourceField,
		refs.NewNamespacedName(resource.Namespace, resource.Name),
		apis,
	); err != nil {
		err = fmt.Errorf("an error occurred while checking if the api resource is linked to an api definition: %w", err)
		return err
	}

	if len(apis.Items) > 0 {
		return fmt.Errorf("resource is referenced and will remain")
	}

	util.RemoveFinalizer(resource, keys.ApiResourceFinalizer)

	return client.Update(ctx, resource)
}
