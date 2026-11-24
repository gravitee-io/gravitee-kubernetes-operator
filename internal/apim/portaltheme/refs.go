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

package portaltheme

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

// ResolveRefs resolves the theme's management context for the reconcile. While deleting, a theme portals
// still activate is refused first, whether its context is there or not; a theme without a context has
// nothing in APIM and needs no context.
func ResolveRefs(ctx context.Context, obj *v1alpha1.PortalTheme, _ string) error {
	if obj.IsBeingDeleted() {
		if err := DeleteGuard(ctx, obj); err != nil {
			return err
		}
		if !obj.HasContext() {
			return nil
		}
	}
	if !obj.HasContext() {
		return fmt.Errorf("theme [%s] has no management context", obj.GetName())
	}
	return navigation.ResolveContext(ctx, obj)
}

// DeleteGuard refuses to delete a theme portals still activate.
func DeleteGuard(ctx context.Context, obj *v1alpha1.PortalTheme) error {
	return search.AssertNoPortalThemeRef(ctx, obj)
}
