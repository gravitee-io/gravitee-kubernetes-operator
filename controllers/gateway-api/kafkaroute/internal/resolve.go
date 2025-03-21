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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	coreV1 "k8s.io/api/core/v1"
	kErrors "k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func Resolve(ctx context.Context, route *v1alpha1.KafkaRoute) error {
	route.Status.Parents = make([]gwAPIv1.RouteParentStatus, len(route.Spec.ParentRefs))
	if err := resolveParents(ctx, route); err != nil {
		return err
	}
	if err := resolveBackendRefs(ctx, route); err != nil {
		return err
	}
	return nil
}

func resolveParents(ctx context.Context, route *v1alpha1.KafkaRoute) error {
	for i, ref := range route.Spec.ParentRefs {
		if resolved, err := resolveParent(ctx, route, ref); err != nil {
			return err
		} else {
			status := gateway.InitRouteParentStatus(ref)
			k8s.SetCondition(status, resolved)
			route.Status.Parents[i] = *status.Object
		}
	}
	return nil
}

func resolveParent(
	ctx context.Context,
	route *v1alpha1.KafkaRoute,
	ref gwAPIv1.ParentReference,
) (*metaV1.Condition, error) {
	conditionBuilder := k8s.NewResolvedRefsConditionBuilder(route.Generation)

	if !k8s.IsGatewayKind(ref) {
		conditionBuilder.RejectInvalidGatewayKind("parent reference must be of Gateway kind")
		return conditionBuilder.Build(), nil
	}

	_, err := k8s.ResolveGateway(ctx, route.ObjectMeta, ref)
	if client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	if kErrors.IsNotFound(err) {
		conditionBuilder.RejectNoMatchingParent("unable to resolve parent reference")
		return conditionBuilder.Build(), nil
	}

	return conditionBuilder.ResolveRefs("parent has been resolved").Build(), nil
}

func resolveBackendRefs(ctx context.Context, route *v1alpha1.KafkaRoute) error {
	resolvedBuilder := k8s.NewResolvedRefsConditionBuilder(route.Generation)

	for i, ref := range route.Spec.BackendRefs {
		if !k8s.IsServiceKind(ref.BackendObjectReference) {
			resolvedBuilder.RejectInvalidBackendKind(
				fmt.Sprintf("backend %d is not of Service kind", i),
			)
		}
		if resolved, err := isResolvedBackend(
			ctx, route, ref.BackendObjectReference,
		); err != nil {
			return err
		} else if !resolved {
			resolvedBuilder.RejectBackendNotFound(
				fmt.Sprintf("backend %d", i),
			)
		}
	}

	resolved := resolvedBuilder.Build()

	if resolved.Status == metaV1.ConditionFalse {
		for i := range route.Status.Parents {
			ref := route.Status.Parents[i]
			status := gateway.WrapRouteParentStatus(&ref)
			k8s.SetCondition(status, resolved)
			route.Status.Parents[i] = *status.Object
		}
	}

	return nil
}

func isResolvedBackend(
	ctx context.Context,
	route *v1alpha1.KafkaRoute,
	ref gwAPIv1.BackendObjectReference,
) (bool, error) {
	ns := ref.Namespace
	if ns == nil {
		gwNs := gwAPIv1.Namespace(route.Namespace)
		ns = &gwNs
	}

	key := client.ObjectKey{Namespace: string(*ns), Name: string(ref.Name)}
	svc := &coreV1.Service{}

	if err := k8s.GetClient().Get(ctx, key, svc); client.IgnoreNotFound(err) != nil {
		return false, err
	} else if kErrors.IsNotFound(err) {
		return false, nil
	}

	return true, nil
}
