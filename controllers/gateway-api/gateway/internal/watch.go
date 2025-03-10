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
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func WatchAttachedRoutes() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(listAttachedRoutes)
}

func listAttachedRoutes(ctx context.Context, obj client.Object) []reconcile.Request {
	httpRoute, ok := obj.(*gwAPIv1.HTTPRoute)
	if !ok {
		return nil
	}
	gateways := &gwAPIv1.GatewayList{}
	if err := k8s.GetClient().List(ctx, gateways); err != nil {
		return nil
	}
	var recs []reconcile.Request
	for _, gateway := range gateways.Items {
		for _, ref := range httpRoute.Spec.ParentRefs {
			if isParent(gateway, ref) {
				recs = append(recs, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Namespace: gateway.Namespace,
						Name:      gateway.Name,
					},
				})
			}
		}
	}
	return recs
}

func isParent(gw gwAPIv1.Gateway, ref gwAPIv1.ParentReference) bool {
	return k8s.IsGatewayKind(ref) && gw.Name == string(ref.Name)
}
