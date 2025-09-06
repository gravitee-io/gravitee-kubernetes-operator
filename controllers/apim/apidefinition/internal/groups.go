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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ResolveGroupRefs(ctx context.Context, referer core.ConditionAwareObject, refs []core.ObjectRef) ([]string, error) {
	groups := []string{}

	if refs == nil || reflect.ValueOf(refs).IsNil() {
		return groups, nil
	}

	for _, ref := range refs {
		group := new(v1alpha1.Group)
		nsn := getNamespacedName(ref, referer.GetNamespace())
		err := k8s.GetClient().Get(ctx, nsn, group)
		if client.IgnoreNotFound(err) != nil {
			return groups, err
		} else if err != nil {
			log.Debug(ctx, "Skipping group reference "+nsn.String()+" as it does not exist")
			k8s.SetCondition(
				referer,
				k8s.
					NewResolvedRefsConditionBuilder(referer.GetGeneration()).
					RejectGroupNotFound("Group "+nsn.String()+" could not be found").
					Build(),
			)
			continue
		}
		if !slices.Contains(groups, group.Spec.Name) {
			groups = append(groups, group.Spec.Name)
		}
	}
	return groups, nil
}
