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
	apimportallink "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portallink"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

// resolvedParent is what reconcile checks on the PortalLink's owning portal or API
// before syncing: whether it carries a management context.
type resolvedParent struct {
	name       string
	kind       string
	hasContext bool
}

// resolveParent resolves the link's owning Portal or API and derives the
// management context to sync against. Exactly one of portalRef / apiRef is
// expected to be set (enforced by admission). Whether the parent carries a
// management context is reported via hasContext so the delete path can
// tolerate a parent that was never synced.
func resolveParent(ctx context.Context, link *v1alpha1.PortalLink) (*resolvedParent, error) {
	ns := link.Namespace

	if link.IsPortalLink() == link.IsApiLink() {
		// both unset or both set — admission normally rejects this; guard here for
		// when webhooks are disabled/bypassed so we never nil-deref in ResolveAPI
		// nor silently ignore one ref.
		return nil, gerrors.NewIllegalStateError(
			fmt.Errorf(
				"portal link [%s/%s] must reference exactly one of portalRef or apiRef",
				link.Namespace, link.Name,
			),
		)
	}

	if link.IsPortalLink() {
		prtl, err := dynamic.ResolvePortal(ctx, link.GetPortalRef(), ns)
		if err != nil {
			return nil, err
		}
		return &resolvedParent{
			name:       prtl.GetName(),
			kind:       "portal",
			hasContext: prtl.HasContext(),
		}, nil
	}

	api, err := dynamic.ResolveAPI(ctx, link.GetApiRef(), ns)
	if err != nil {
		return nil, err
	}
	return &resolvedParent{
		name:       api.GetName(),
		kind:       "API",
		hasContext: api.HasContext(),
	}, nil
}

func CreateOrUpdate(ctx context.Context, link *v1alpha1.PortalLink) error {
	parent, err := resolveParent(ctx, link)
	if err != nil {
		return err
	}

	if !parent.hasContext {
		return gerrors.NewIllegalStateError(
			fmt.Errorf("portal link %s parent [%s] has no management context", parent.kind, parent.name),
		)
	}

	client, err := apimportallink.CreateAPIMClient(ctx, link)
	if err != nil {
		return err
	}

	dto, err := apimportallink.ToDTO(link)
	if err != nil {
		return gerrors.NewIllegalStateError(err)
	}

	resp, err := apimportallink.Upsert(ctx, client, dto)
	if err != nil {
		return gerrors.NewControlPlaneError(err)
	}

	return apimportallink.UpdateStatus(ctx, link, resp)
}
