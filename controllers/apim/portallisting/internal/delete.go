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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func Delete(ctx context.Context, listing *v1alpha1.PortalListing) error {
	if !util.ContainsFinalizer(listing, core.PortalListingFinalizer) {
		return nil
	}

	ns := listing.Namespace

	prtl, err := dynamic.ResolvePortal(ctx, listing.GetPortalRef(), ns)
	if err != nil {
		// Parent Portal already gone: nothing left to sync against, let the finalizer be removed.
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	if !prtl.HasContext() {
		// Nothing was ever synced to APIM without a context; let the finalizer be removed.
		return nil
	}

	apimClient, err := apim.FromContextRef(ctx, prtl.ContextRef(), prtl.GetNamespace())
	if err != nil {
		// ManagementContext already gone: APIM is unreachable, let the finalizer be removed.
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	if err := gerrors.IgnoreNotFound(apimClient.Listings.Delete(listing, prtl)); err != nil {
		return gerrors.NewControlPlaneError(err)
	}

	return nil
}
