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
	apimportallisting "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portallisting"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

func CreateOrUpdate(ctx context.Context, listing *v1alpha1.PortalListing) error {
	ns := listing.Namespace

	prtl, err := dynamic.ResolvePortal(ctx, listing.GetPortalRef(), ns)
	if err != nil {
		return err
	}

	if !prtl.HasContext() {
		return gerrors.NewIllegalStateError(
			fmt.Errorf("portal [%s] has no management context", prtl.GetName()),
		)
	}

	client, err := apimportallisting.CreateAPIMClient(ctx, listing)
	if err != nil {
		return err
	}

	dto, err := apimportallisting.ToDTO(listing)
	if err != nil {
		return gerrors.NewIllegalStateError(err)
	}

	resp, err := apimportallisting.Upsert(ctx, client, dto)
	if err != nil {
		return gerrors.NewControlPlaneError(err)
	}

	return apimportallisting.UpdateStatus(ctx, listing, resp)
}
