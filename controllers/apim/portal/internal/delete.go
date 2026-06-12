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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Delete(ctx context.Context, prtl *v1alpha1.Portal) error {
	if !util.ContainsFinalizer(prtl, core.PortalFinalizer) {
		return nil
	}

	if err := search.AssertNoPortalListingRef(ctx, prtl); err != nil {
		return err
	}

	if err := search.AssertNoPortalDocumentationRef(ctx, prtl); err != nil {
		return err
	}

	if !prtl.HasContext() {
		// Nothing was ever synced to APIM without a context; let the finalizer be removed.
		return nil
	}

	apimClient, err := apim.FromContextRef(ctx, prtl.ContextRef(), prtl.GetNamespace())
	if err != nil {
		return err
	}

	if err := gerrors.IgnoreNotFound(apimClient.Portals.Delete(prtl)); err != nil {
		return gerrors.NewControlPlaneError(err)
	}

	return nil
}
