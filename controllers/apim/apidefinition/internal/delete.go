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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apiv4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Delete(ctx context.Context, api core.ApiDefinitionObject) error {
	if !util.ContainsFinalizer(api, core.ApiDefinitionFinalizer) {
		return nil
	}

	if err := search.AssertNoApiDocumentationRef(ctx, api); err != nil {
		return err
	}

	if err := search.AssertNoApiPortalListingRef(ctx, api); err != nil {
		return err
	}

	if err := search.AssertNoApiPortalLinkRef(ctx, api); err != nil {
		return err
	}

	if v4API, ok := api.(*v1alpha1.ApiV4Definition); ok {
		return deleteV4(ctx, v4API)
	}

	if api.HasContext() {
		if err := deleteWithContext(ctx, api); err != nil {
			return err
		}
	}

	return nil
}

func deleteV4(ctx context.Context, api *v1alpha1.ApiV4Definition) error {
	client, err := apiv4.CreateAPIMClient(ctx, api)
	if err != nil {
		return err
	}
	dto, err := apiv4.ToDTO(api)
	if err != nil {
		return err
	}
	return apiv4.Delete(ctx, client, dto)
}

func deleteWithContext(ctx context.Context, api core.ApiDefinitionObject) error {
	apimClient, err := apim.FromContextRef(ctx, api.ContextRef(), api.GetNamespace())
	if err != nil {
		return err
	}
	switch {
	case api.GetDefinitionVersion() == core.ApiV2:
		return errors.IgnoreNotFound(apimClient.APIs.DeleteV2(api.GetID()))
	default:
		return fmt.Errorf("unknown version %s", api.GetDefinitionVersion())
	}
}
