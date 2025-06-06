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

	"slices"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	kErrors "k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func Accept(ctx context.Context, route *v1alpha1.KafkaRoute) error {
	return acceptParents(ctx, route)
}

func acceptParents(ctx context.Context, route *v1alpha1.KafkaRoute) error {
	for i := range route.Spec.ParentRefs {
		status := gateway.WrapRouteParentStatus(&route.Status.Parents[i])
		if accepted, err := acceptParent(ctx, route, status); err != nil {
			return err
		} else {
			k8s.SetCondition(status, accepted)
			route.Status.Parents[i] = *status.Object
		}
	}
	return nil
}

func acceptParent(
	ctx context.Context,
	route *v1alpha1.KafkaRoute,
	status *gateway.RouteParentStatus,
) (*metaV1.Condition, error) {
	accepted := k8s.NewAcceptedConditionBuilder(route.Generation).Accept("route is accepted")

	if !k8s.IsResolved(status) {
		accepted.RejectNoMatchingParent("parent ref could not be resolved")
		return accepted.Build(), nil
	}

	gw, err := k8s.ResolveGateway(ctx, route.ObjectMeta, status.Object.ParentRef)
	if client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	if kErrors.IsNotFound(err) {
		accepted.RejectNoMatchingParent("unable to resolve parent reference")
		return accepted.Build(), nil
	}

	if !k8s.IsAccepted(gateway.WrapGateway(gw)) {
		accepted.RejectNoMatchingParent("parent is not accepted")
		return accepted.Build(), nil
	}

	if !supportsKafka(gw) {
		accepted.RejectNoMatchingParent("parent ref does not support Kafka routes")
		return accepted.Build(), nil
	}

	return accepted.Accept("route has been accepted").Build(), nil
}

func IsAccepted(
	ctx context.Context,
	route *v1alpha1.KafkaRoute,
) bool {
	for i := range route.Status.Parents {
		status := gateway.WrapRouteParentStatus(&route.Status.Parents[i])
		if !k8s.IsAccepted(status) {
			return false
		}
	}
	return true
}

func supportsKafka(gw *gwAPIv1.Gateway) bool {
	return slices.ContainsFunc(gw.Status.Listeners, k8s.IsKafkaListenerStatus)
}
