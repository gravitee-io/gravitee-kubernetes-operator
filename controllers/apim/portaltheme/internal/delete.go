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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimportaltheme "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portaltheme"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Delete(ctx context.Context, thm *v1alpha1.PortalTheme) error {
	if !util.ContainsFinalizer(thm, core.PortalThemeFinalizer) {
		return nil
	}

	if err := search.AssertNoPortalThemeRef(ctx, thm); err != nil {
		return err
	}

	if !thm.HasContext() {
		// Nothing was ever synced to APIM without a context; let the finalizer be removed.
		return nil
	}

	client, err := apimportaltheme.CreateAPIMClient(ctx, thm)
	if err != nil {
		// ManagementContext already gone: APIM is unreachable, let the finalizer be removed.
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	dto, err := apimportaltheme.ToDTO(thm)
	if err != nil {
		return gerrors.NewIllegalStateError(err)
	}

	if err := apimportaltheme.Delete(ctx, client, dto); err != nil {
		return gerrors.NewControlPlaneError(err)
	}

	return nil
}
