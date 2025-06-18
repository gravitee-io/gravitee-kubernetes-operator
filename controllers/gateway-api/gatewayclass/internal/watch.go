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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func WatchGatewayClassParameters() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(requestsFromGatewayClassParameters)
}

func requestsFromGatewayClassParameters(ctx context.Context, obj client.Object) []reconcile.Request {
	gwcp, ok := obj.(*v1alpha1.GatewayClassParameters)
	if !ok {
		return nil
	}
	if !gwcp.DeletionTimestamp.IsZero() {
		return nil
	}
	listOpts := &client.ListOptions{}
	list := &gwAPIv1.GatewayClassList{}
	if err := k8s.GetClient().List(ctx, list, listOpts); err != nil {
		return nil
	}
	reqs := make([]reconcile.Request, 0)
	for i := range list.Items {
		gwc := list.Items[i]
		if k8s.HasGatewayClassParameters(&gwc, gwcp) {
			reqs = append(reqs, buildRequest(gwc))
		}
	}
	return reqs
}

func buildRequest(gwc gwAPIv1.GatewayClass) reconcile.Request {
	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: gwc.Namespace,
			Name:      gwc.Name,
		},
	}
}
