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

// FindAuthzEntityChildren returns, as namespace/name and sorted, the AuthzEntity resources that
// declare entity as a parent and sync through its ManagementContext. Entities created on the
// platform by another source are not found.
func FindAuthzEntityChildren(ctx context.Context, entity *v1alpha1.AuthzEntity) ([]string, error) {
	keys, err := FindAuthzEntityChildKeys(ctx, entity)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(keys))
	for _, key := range keys {
		names = append(names, key.String())
	}
	sort.Strings(names)
	return names, nil
}

// FindAuthzEntityChildKeys returns the keys of the AuthzEntity resources that declare entity as a
// parent and sync through its ManagementContext, entity itself excluded.
func FindAuthzEntityChildKeys(ctx context.Context, entity *v1alpha1.AuthzEntity) ([]types.NamespacedName, error) {
	contextOf := dynamic.ContextOf(entity)
	if contextOf == "" {
		return nil, nil
	}
	key := AuthzEntityParentKey(contextOf, entity.Spec.EntityType, entity.Spec.EntityID)
	children := &v1alpha1.AuthzEntityList{}
	filter := &client.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{AuthzEntityParentField.String(): key}),
	}
	if err := k8s.GetClient().List(ctx, children, filter); err != nil {
		return nil, fmt.Errorf("unable to list AuthzEntity: %w", err)
	}

	self := client.ObjectKeyFromObject(entity)
	keys := make([]types.NamespacedName, 0, len(children.Items))
	for i := range children.Items {
		if child := client.ObjectKeyFromObject(&children.Items[i]); child != self {
			keys = append(keys, child)
		}
	}
	return keys, nil
}
