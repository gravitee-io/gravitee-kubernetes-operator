// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package search

import (
	"context"
	"fmt"
	"sort"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AuthzEntityParentKey is how AuthzEntityParentField spells an entity: the ManagementContext it
// syncs through (see dynamic.ContextOf), then its type and id.
func AuthzEntityParentKey(contextOf, entityType, entityID string) string {
	return contextOf + "/" + entityType + "/" + entityID
}

// FindAuthzEntityCascadeDeletes returns, as namespace/name sorted, the RESOURCE entities among
// FindAuthzEntityCascade's: those the platform deletes with entity. The PRINCIPAL ones are only
// detached, and applied again as declared.
func FindAuthzEntityCascadeDeletes(ctx context.Context, entity *v1alpha1.AuthzEntity) ([]string, error) {
	cascade, err := findAuthzEntityCascade(ctx, entity)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(cascade))
	for _, child := range cascade {
		if child.Spec.Kind == authzentity.KindResource {
			names = append(names, client.ObjectKeyFromObject(child).String())
		}
	}
	sort.Strings(names)
	return names, nil
}

// FindAuthzEntityCascade returns the keys of the AuthzEntity resources of entity's ManagementContext
// that its delete changes on the platform, entity excluded. The platform deletes every RESOURCE child
// of a deleted entity, transitively, and detaches every PRINCIPAL child: the search follows the
// declared parents the same way, through RESOURCE children only. An entity another source created
// is not found, nor what the platform cascades through it.
func FindAuthzEntityCascade(ctx context.Context, entity *v1alpha1.AuthzEntity) ([]types.NamespacedName, error) {
	cascade, err := findAuthzEntityCascade(ctx, entity)
	if err != nil {
		return nil, err
	}
	keys := make([]types.NamespacedName, 0, len(cascade))
	for _, child := range cascade {
		keys = append(keys, client.ObjectKeyFromObject(child))
	}
	return keys, nil
}

// findAuthzEntityCascade returns the entities FindAuthzEntityCascade returns the keys of, in the
// order the search finds them.
func findAuthzEntityCascade(ctx context.Context, entity *v1alpha1.AuthzEntity) ([]*v1alpha1.AuthzEntity, error) {
	contextOf := dynamic.ContextOf(entity)
	if contextOf == "" {
		return nil, nil
	}

	seen := map[types.NamespacedName]bool{client.ObjectKeyFromObject(entity): true}
	found := make([]*v1alpha1.AuthzEntity, 0)
	deleted := []*v1alpha1.AuthzEntity{entity}
	for len(deleted) > 0 {
		parent := deleted[0]
		deleted = deleted[1:]
		children, err := listAuthzEntityChildren(ctx, contextOf, parent)
		if err != nil {
			return nil, err
		}
		for i := range children {
			child := &children[i]
			key := client.ObjectKeyFromObject(child)
			if seen[key] {
				continue
			}
			seen[key] = true
			found = append(found, child)
			if child.Spec.Kind == authzentity.KindResource {
				deleted = append(deleted, child)
			}
		}
	}
	return found, nil
}

// listAuthzEntityChildren lists the AuthzEntity resources of a ManagementContext that declare parent.
func listAuthzEntityChildren(
	ctx context.Context,
	contextOf string,
	parent *v1alpha1.AuthzEntity,
) ([]v1alpha1.AuthzEntity, error) {
	key := AuthzEntityParentKey(contextOf, parent.Spec.EntityType, parent.Spec.EntityID)
	children := &v1alpha1.AuthzEntityList{}
	filter := &client.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{AuthzEntityParentField.String(): key}),
	}
	if err := k8s.GetClient().List(ctx, children, filter); err != nil {
		return nil, fmt.Errorf("unable to list AuthzEntity: %w", err)
	}
	return children.Items, nil
}
