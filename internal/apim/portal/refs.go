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

package portal

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gerrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

// ResolveRefs resolves the portal's management context and active theme for the reconcile. While
// deleting, a portal listings, links or documentation pages still reference is refused first, whether
// its context is there or not; a portal without a context has nothing in APIM and needs no context. The
// theme is not resolved while deleting: a theme deleted first must not release the portal's finalizer
// before the portal is deleted from APIM.
func ResolveRefs(ctx context.Context, obj *v1alpha1.Portal, _ string) error {
	if obj.IsBeingDeleted() {
		if err := DeleteGuard(ctx, obj); err != nil {
			return err
		}
		if !obj.HasContext() {
			return nil
		}
		return navigation.ResolveContext(ctx, obj)
	}
	if !obj.HasContext() {
		return fmt.Errorf("portal [%s] has no management context", obj.GetName())
	}
	if obj.HasThemeRef() {
		if kind := obj.Spec.Theme.Kind; kind != "" && dynamic.ResourceFromKind(kind) != core.CRDPortalThemeResource {
			return fmt.Errorf("themeRef [%s] is of kind %s, but only PortalTheme is supported", obj.Spec.Theme.Name, kind)
		}
	}
	if err := navigation.ResolveContext(ctx, obj); err != nil {
		return err
	}
	return resolveTheme(ctx, obj)
}

func resolveTheme(ctx context.Context, obj *v1alpha1.Portal) error {
	if !obj.HasThemeRef() {
		return nil
	}
	if _, err := dynamic.ResolvePortalTheme(ctx, obj.GetThemeRef(), obj.GetNamespace()); err != nil {
		return gerrors.NewRefNotFoundError(
			fmt.Sprintf("portal [%s] references theme [%v] that can't be resolved", obj.GetName(), obj.GetThemeRef()), err,
		)
	}
	return nil
}

// DeleteGuard refuses to delete a portal listings, links or documentation pages still reference.
func DeleteGuard(ctx context.Context, obj *v1alpha1.Portal) error {
	if err := search.AssertNoPortalListingRef(ctx, obj); err != nil {
		return err
	}
	if err := search.AssertNoPortalLinkRef(ctx, obj); err != nil {
		return err
	}
	return search.AssertNoPortalDocumentationRef(ctx, obj)
}
