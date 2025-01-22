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
	"reflect"
	"slices"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ResolveGroupRefs(ctx context.Context, api core.ApiDefinitionObject) error {
	groupRefs := api.GetGroupRefs()

	if groupRefs == nil || reflect.ValueOf(groupRefs).IsNil() {
		return nil
	}

	groups := api.GetGroups()
	for _, ref := range groupRefs {
		group := new(v1alpha1.Group)
		nsn := getNamespacedName(ref, api.GetNamespace())
		err := k8s.GetClient().Get(ctx, nsn, group)
		if client.IgnoreNotFound(err) != nil {
			return err
		} else if err != nil {
			log.Debug(ctx, "Skipping group reference "+ref.String()+" as it does not exist")
			continue
		}
		if !slices.Contains(groups, group.Spec.Name) {
			groups = append(groups, group.Spec.Name)
		}
	}
	api.SetGroups(groups)
	return nil
}

func getNamespacedName(ref core.ObjectRef, apiNs string) types.NamespacedName {
	if ref.GetNamespace() == "" {
		return types.NamespacedName{
			Name:      ref.GetName(),
			Namespace: apiNs,
		}
	}
	return ref.NamespacedName()
}
