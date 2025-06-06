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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/gateway-api/httproute/internal/mapper"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func DetectConflicts(ctx context.Context, route *gwAPIv1.HTTPRoute) error {
	api := mapper.Map(ctx, route)
	conflict, err := base.FindConflictingPath(ctx, api)
	if err != nil {
		return err
	}

	if conflict.IsZero() {
		return nil
	}

	for i, ref := range route.Spec.ParentRefs {
		if hasMatchingTag(conflict, route, ref) {
			condition := k8s.NewListenerConflictedConditionBuilder(route.Generation)
			status := gateway.WrapRouteParentStatus(&route.Status.Parents[i])
			k8s.SetCondition(status, condition.RejectConflictingPath(conflict.ID).Build())
		}
	}

	return nil
}

func hasMatchingTag(
	conflict base.Conflict,
	route *gwAPIv1.HTTPRoute,
	parentRef gwAPIv1.ParentReference,
) bool {
	parentTag := buildTag(route, parentRef)
	for _, tag := range conflict.Tags {
		if tag == parentTag {
			return true
		}
	}
	return false
}

func buildTag(route *gwAPIv1.HTTPRoute, ref gwAPIv1.ParentReference) string {
	ns := ref.Namespace
	if ns == nil {
		routeNS := gwAPIv1.Namespace(route.Namespace)
		ns = &routeNS
	}
	return k8s.BuildTag(string(*ns), string(ref.Name))
}
