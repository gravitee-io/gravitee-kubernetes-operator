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

package navigation

import (
	"context"
	"errors"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

// ErrNotExactlyOneParent is returned for a link or documentation page that references both a portal and
// an API, or neither.
var ErrNotExactlyOneParent = errors.New("exactly one of portalRef / apiRef must be set")

// Parent is the HRID of the portal or the API a link or documentation page belongs to: exactly one is
// set. It goes in the request path, never in the payload.
type Parent struct {
	PortalHrid string
	ApiHrid    string
}

// IsAPI reports a resource that belongs to an API.
func (p Parent) IsAPI() bool {
	return p.ApiHrid != ""
}

// ParentFromRefs computes the parent HRID from the refs, a ref without namespace pointing to ns, as
// resolving it would.
func ParentFromRefs(portalRef, apiRef *refs.NamespacedName, ns string) (Parent, error) {
	if (portalRef == nil) == (apiRef == nil) {
		return Parent{}, ErrNotExactlyOneParent
	}
	if apiRef != nil {
		return Parent{ApiHrid: hridOf(*apiRef, ns)}, nil
	}
	return Parent{PortalHrid: hridOf(*portalRef, ns)}, nil
}

func hridOf(ref refs.NamespacedName, ns string) string {
	if ref.IsMissingNamespace() {
		ref.SetNamespace(ns)
	}
	return ref.HRID()
}

// ResolveParent reads the portal or the API a link or documentation page references from the cluster.
func ResolveParent(
	ctx context.Context, portalRef, apiRef *refs.NamespacedName, ns string,
) (core.ContextAwareObject, error) {
	if (portalRef == nil) == (apiRef == nil) {
		return nil, ErrNotExactlyOneParent
	}
	if apiRef != nil {
		api, err := dynamic.ResolveAPI(ctx, apiRef, ns)
		if err != nil {
			return nil, err
		}
		return api, nil
	}
	prtl, err := dynamic.ResolvePortal(ctx, portalRef, ns)
	if err != nil {
		return nil, err
	}
	return prtl, nil
}

// ParentKind names the kind of a resolved parent the way the messages always have: "portal" or "API".
func ParentKind(parent core.ContextAwareObject) string {
	if _, ok := parent.(*v1alpha1.Portal); ok {
		return "portal"
	}
	return "API"
}

// CheckPortalSynced returns dynamic.ErrPortalNotSynced while parent is a portal GKO has not written to APIM
// yet. An API parent always passes: APIM keeps a page or link of an API it does not hold until the API lands.
func CheckPortalSynced(parent core.ContextAwareObject) error {
	if prtl, ok := parent.(*v1alpha1.Portal); ok {
		return dynamic.CheckPortalSynced(prtl)
	}
	return nil
}
