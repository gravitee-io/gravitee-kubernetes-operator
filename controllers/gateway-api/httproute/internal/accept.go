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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	coreV1 "k8s.io/api/core/v1"
	kErrors "k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func Accept(ctx context.Context, route *gwAPIv1.HTTPRoute) error {
	return acceptParents(ctx, route)
}

func acceptParents(ctx context.Context, route *gwAPIv1.HTTPRoute) error {
	for i, ref := range route.Spec.ParentRefs {
		status := gateway.WrapRouteParentStatus(&route.Status.Parents[i])
		if accepted, err := acceptParent(ctx, route, ref, status); err != nil {
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
	route *gwAPIv1.HTTPRoute,
	ref gwAPIv1.ParentReference,
	status *gateway.RouteParentStatus,
) (*metaV1.Condition, error) {
	accepted := k8s.NewAcceptedConditionBuilder(route.Generation).Accept("route is accepted")

	if !k8s.IsResolved(status) {
		unresolved := k8s.GetCondition(status, k8s.ConditionResolvedRefs)
		if unresolved.Reason == string(gwAPIv1.RouteReasonNoMatchingParent) {
			accepted.RejectNoMatchingParent("parent ref could not be resolved")
			return accepted.Build(), nil
		}
	}

	gw, err := k8s.ResolveGateway(ctx, route.ObjectMeta, status.Object.ParentRef)
	if client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	if kErrors.IsNotFound(err) {
		accepted.RejectNoMatchingParent("unable to resolve parent reference")
		return accepted.Build(), nil
	}

	if ref.SectionName != nil {
		if k8s.FindListenerIndexBySectionName(gw, *ref.SectionName) == -1 {
			accepted.RejectNoMatchingParent("section name does not exist")
			return accepted.Build(), nil
		}
	}

	if hasNsSupport, err := supportsRouteNamespace(ctx, gw, ref, route); err != nil {
		return nil, err
	} else if !hasNsSupport {
		accepted.RejectNotAllowedByListeners("parent namespace policy does not allow route namespace")
		return accepted.Build(), nil
	}

	if !supportsHTTP(gw, ref) {
		accepted.RejectNoMatchingParent("parent ref does not support HTTP routes")
		return accepted.Build(), nil
	}

	if !k8s.HasIntersectingHostName(route, gw, ref) {
		accepted.RejectNoMatchingListenerHostname("parent hostname and route hostnames do not intersect")
		return accepted.Build(), nil
	}

	return accepted.Build(), nil
}

func supportsHTTP(gw *gwAPIv1.Gateway, ref gwAPIv1.ParentReference) bool {
	if ref.SectionName != nil {
		lIdx := k8s.FindListenerIndexBySectionName(gw, *ref.SectionName)
		return k8s.HasHTTPListenerAtIndex(gw, lIdx)
	}
	for i := range gw.Status.Listeners {
		if k8s.HasHTTPListenerAtIndex(gw, i) {
			return true
		}
	}
	return false
}

func supportsRouteNamespace(
	ctx context.Context,
	gw *gwAPIv1.Gateway,
	ref gwAPIv1.ParentReference,
	route *gwAPIv1.HTTPRoute,
) (bool, error) {
	if ref.SectionName != nil {
		lIdx := k8s.FindListenerIndexBySectionName(gw, *ref.SectionName)
		return supportsRouteNamespaceAtListenerIndex(
			ctx, gw, ref, route, lIdx,
		)
	}
	for i := range gw.Spec.Listeners {
		if hasNsSupport, err := supportsRouteNamespaceAtListenerIndex(
			ctx, gw, ref, route, i,
		); err != nil {
			return false, err
		} else if hasNsSupport {
			return true, nil
		}
	}
	return false, nil
}

func supportsRouteNamespaceAtListenerIndex(
	ctx context.Context,
	gw *gwAPIv1.Gateway,
	ref gwAPIv1.ParentReference,
	route *gwAPIv1.HTTPRoute,
	lIdx int,
) (bool, error) {
	listener := gw.Spec.Listeners[lIdx]
	if *listener.AllowedRoutes.Namespaces.From == gwAPIv1.NamespacesFromAll {
		return true, nil
	}
	if *listener.AllowedRoutes.Namespaces.From == gwAPIv1.NamespacesFromSame {
		return ref.Namespace == nil || string(*ref.Namespace) == route.Namespace, nil
	}
	if *listener.AllowedRoutes.Namespaces.From == gwAPIv1.NamespacesFromSelector {
		ns, err := resolveNS(ctx, route.Namespace)
		if err != nil {
			return false, err
		}
		nsLabels := ns.Labels
		selectorLabels := listener.AllowedRoutes.Namespaces.Selector
		for k := range selectorLabels.MatchLabels {
			if nsLabels[k] != selectorLabels.MatchLabels[k] {
				return false, nil
			}
		}
		// For now we don't support label expressions
		return true, nil
	}
	return false, nil
}

func resolveNS(ctx context.Context, name string) (*coreV1.Namespace, error) {
	ns := &coreV1.Namespace{}
	err := k8s.GetClient().Get(ctx, client.ObjectKey{Name: name}, ns)
	if err != nil {
		return nil, err
	}
	return ns, nil
}
