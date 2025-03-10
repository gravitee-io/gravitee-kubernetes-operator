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
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func WatchAttachedRoutes() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(requestsFromHTTPRoutes)
}

func WatchServices() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(requestFromServices)
}

func requestFromServices(ctx context.Context, obj client.Object) []reconcile.Request {
	svc, ok := obj.(*coreV1.Service)
	if !ok {
		return nil
	}
	if !k8s.IsGatewayComponent(svc) {
		return nil
	}
	if !svc.DeletionTimestamp.IsZero() {
		return nil
	}
	listOpts := &client.ListOptions{
		Namespace: svc.Namespace,
	}
	gateways := &gwAPIv1.GatewayList{}
	if err := k8s.GetClient().List(ctx, gateways, listOpts); err != nil {
		return nil
	}
	reqs := make([]reconcile.Request, len(gateways.Items))
	for i := range gateways.Items {
		gw := gateways.Items[i]
		if k8s.IsGatewayDependent(gateway.WrapGateway(&gw), svc) {
			reqs[i] = reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: gw.Namespace,
					Name:      gw.Name,
				},
			}
		}
	}
	return reqs
}

func requestsFromHTTPRoutes(ctx context.Context, obj client.Object) []reconcile.Request {
	httpRoute, ok := obj.(*gwAPIv1.HTTPRoute)
	if !ok {
		return nil
	}
	listOpts := &client.ListOptions{
		Namespace: httpRoute.Namespace,
	}
	gateways := &gwAPIv1.GatewayList{}
	if err := k8s.GetClient().List(ctx, gateways, listOpts); err != nil {
		return nil
	}
	var reqs []reconcile.Request
	for _, gateway := range gateways.Items {
		for _, ref := range httpRoute.Spec.ParentRefs {
			if isParent(gateway, ref) {
				reqs = append(reqs, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Namespace: gateway.Namespace,
						Name:      gateway.Name,
					},
				})
			}
		}
	}
	return reqs
}

func isParent(gw gwAPIv1.Gateway, ref gwAPIv1.ParentReference) bool {
	return k8s.IsGatewayKind(ref) && gw.Name == string(ref.Name)
}
