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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/log"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwAPIv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

func WatchReferenceGrants() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(requestFromReferenceGrant)
}

func requestFromReferenceGrant(ctx context.Context, obj client.Object) []reconcile.Request {
	_, ok := obj.(*gwAPIv1beta1.ReferenceGrant)
	if !ok {
		return nil
	}

	listOpts := &client.ListOptions{}
	list := &gwAPIv1.HTTPRouteList{}
	if err := k8s.GetClient().List(ctx, list, listOpts); err != nil {
		return nil
	}
	reqs := make([]reconcile.Request, len(list.Items))
	for i, route := range list.Items {
		reqs[i] = buildRequest(route)
	}
	return reqs
}

func buildRequest(route gwAPIv1.HTTPRoute) reconcile.Request {
	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: route.Namespace,
			Name:      route.Name,
		},
	}
}

func WatchGateways() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(requestsFromGateway)
}

func requestsFromGateway(ctx context.Context, obj client.Object) []reconcile.Request {
	gw, ok := obj.(*gwAPIv1.Gateway)
	if !ok {
		return nil
	}

	listOpts := &client.ListOptions{}
	list := &gwAPIv1.HTTPRouteList{}
	if err := k8s.GetClient().List(ctx, list, listOpts); err != nil {
		log.Error(ctx, err, "failed to list HTTPRoutes when watching gateway")
		return []reconcile.Request{}
	}

	var reqs []reconcile.Request
	for _, route := range list.Items {
		if referencesGateway(&route, gw) {
			reqs = append(reqs, buildRequest(route))
		}
	}

	return reqs
}

func referencesGateway(route *gwAPIv1.HTTPRoute, gw *gwAPIv1.Gateway) bool {
	for _, ref := range route.Spec.ParentRefs {
		if !k8s.IsGatewayKind(ref) {
			continue
		}
		if string(ref.Name) != gw.Name {
			continue
		}
		refNS := ref.Namespace
		if refNS == nil {
			if route.Namespace == gw.Namespace {
				return true
			}
		} else if string(*refNS) == gw.Namespace {
			return true
		}
	}
	return false
}
