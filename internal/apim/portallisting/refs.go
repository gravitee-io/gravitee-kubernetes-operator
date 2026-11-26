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

package portallisting

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

// ResolveRefsForAdmission reads the listing's portal and checks it has a management context, with the
// messages the webhook always reported. The API entries are checked by the webhook.
func ResolveRefsForAdmission(ctx context.Context, obj *v1alpha1.PortalListing, _ string) error {
	_, err := resolvePortal(ctx, obj, false)
	return err
}

// ResolveRefs resolves the portal and its management context for the reconcile. A portal or a context
// deleted first is a NotFound; a portal without a context, while deleting, means nothing was synced. The
// API entries are mapped from their references, never resolved. An upsert waits for a portal GKO has not
// synced yet: APIM refuses a listing of a portal it does not hold while the environment holds another, its
// default portal included.
func ResolveRefs(ctx context.Context, obj *v1alpha1.PortalListing, _ string) error {
	prtl, err := resolvePortal(ctx, obj, obj.IsBeingDeleted())
	if err != nil || !prtl.HasContext() {
		return err
	}
	if err = navigation.ResolveContext(ctx, prtl); err != nil || obj.IsBeingDeleted() {
		return err
	}
	return dynamic.CheckPortalSynced(prtl)
}

// ResolvePortal reads the listing's portal, with the message the webhook always reported.
func ResolvePortal(ctx context.Context, obj *v1alpha1.PortalListing) (*v1alpha1.Portal, error) {
	return resolvePortal(ctx, obj, false)
}

func resolvePortal(ctx context.Context, obj *v1alpha1.PortalListing, deleting bool) (*v1alpha1.Portal, error) {
	prtl, err := dynamic.ResolvePortal(ctx, obj.GetPortalRef(), obj.GetNamespace())
	if err != nil {
		return nil, gerrors.NewRefNotFoundError(
			fmt.Sprintf(
				"portal listing [%s] references portal [%v] that can't be resolved",
				obj.GetName(), obj.GetPortalRef(),
			), err,
		)
	}
	if !prtl.HasContext() && !deleting {
		return nil, fmt.Errorf(
			"referenced portal [%v] has no management context (spec.contextRef)", obj.GetPortalRef(),
		)
	}
	return prtl, nil
}
