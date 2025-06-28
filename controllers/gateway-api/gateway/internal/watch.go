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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func WatchGatewayClasses() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(requestsFromGatewayClass)
}

func WatchHTTPRoutes() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(requestsFromHTTPRoute)
}

func WatchKafkaRoutes() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(requestsFromKafkaRoute)
}

func WatchServices() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(requestsFromService)
}

func WatchSecrets() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(requestsFromSecret)
}

func requestsFromGatewayClass(ctx context.Context, obj client.Object) []reconcile.Request {
	gwc, ok := obj.(*gwAPIv1.GatewayClass)
	if !ok {
		return nil
	}
	if !gwc.DeletionTimestamp.IsZero() {
		return nil
	}
	listOpts := &client.ListOptions{
		Namespace: gwc.Namespace,
	}
	list := &gwAPIv1.GatewayList{}
	if err := k8s.GetClient().List(ctx, list, listOpts); err != nil {
		return nil
	}
	reqs := make([]reconcile.Request, len(list.Items))
	for i := range list.Items {
		gw := list.Items[i]
		if gw.Spec.GatewayClassName == gwAPIv1.ObjectName(gwc.Name) {
			reqs[i] = buildRequest(gw)
		}
	}
	return reqs
}

func requestsFromService(ctx context.Context, obj client.Object) []reconcile.Request {
	svc, ok := obj.(*coreV1.Service)
	if !ok {
		return nil
	}
	if !svc.DeletionTimestamp.IsZero() {
		return nil
	}
	if !k8s.IsGatewayComponent(svc) {
		return nil
	}
	listOpts := &client.ListOptions{
		Namespace: svc.Namespace,
	}
	list := &gwAPIv1.GatewayList{}
	if err := k8s.GetClient().List(ctx, list, listOpts); err != nil {
		return nil
	}
	reqs := make([]reconcile.Request, len(list.Items))
	for i := range list.Items {
		gw := list.Items[i]
		if k8s.IsGatewayDependent(gateway.WrapGateway(&gw), svc) {
			reqs[i] = buildRequest(gw)
		}
	}
	return reqs
}

func requestsFromHTTPRoute(ctx context.Context, obj client.Object) []reconcile.Request {
	httpRoute, ok := obj.(*gwAPIv1.HTTPRoute)
	if !ok {
		return nil
	}
	listOpts := &client.ListOptions{}
	list := &gwAPIv1.GatewayList{}
	if err := k8s.GetClient().List(ctx, list, listOpts); err != nil {
		return nil
	}
	var reqs []reconcile.Request
	for _, gw := range list.Items {
		for _, ref := range httpRoute.Spec.ParentRefs {
			if isParent(gw, ref) {
				reqs = append(reqs, buildRequest(gw))
			}
		}
	}
	return reqs
}

func requestsFromKafkaRoute(ctx context.Context, obj client.Object) []reconcile.Request {
	kafkaRoute, ok := obj.(*v1alpha1.KafkaRoute)
	if !ok {
		return nil
	}
	listOpts := &client.ListOptions{}
	list := &gwAPIv1.GatewayList{}
	if err := k8s.GetClient().List(ctx, list, listOpts); err != nil {
		return nil
	}
	var reqs []reconcile.Request
	for _, gw := range list.Items {
		for _, ref := range kafkaRoute.Spec.ParentRefs {
			if isParent(gw, ref) {
				reqs = append(reqs, buildRequest(gw))
			}
		}
	}
	return reqs
}

func requestsFromSecret(ctx context.Context, obj client.Object) []reconcile.Request {
	secret, ok := obj.(*coreV1.Secret)
	if !ok {
		return nil
	}
	if !secret.DeletionTimestamp.IsZero() {
		return nil
	}
	listOpts := &client.ListOptions{
		Namespace: secret.Namespace,
	}
	list := &gwAPIv1.GatewayList{}
	if err := k8s.GetClient().List(ctx, list, listOpts); err != nil {
		return nil
	}
	var reqs []reconcile.Request
	for _, gw := range list.Items {
		for _, l := range gw.Spec.Listeners {
			if hasSecretRef(l, secret) {
				reqs = append(reqs, buildRequest(gw))
			}
		}
	}
	return reqs
}

func hasSecretRef(
	listener gwAPIv1.Listener,
	secret *coreV1.Secret,
) bool {
	if listener.TLS == nil {
		return false
	}
	for _, ref := range listener.TLS.CertificateRefs {
		if k8s.IsSecretRef(secret, ref) {
			return true
		}
	}
	return false
}

func buildRequest(gateway gwAPIv1.Gateway) reconcile.Request {
	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: gateway.Namespace,
			Name:      gateway.Name,
		},
	}
}

func isParent(gw gwAPIv1.Gateway, ref gwAPIv1.ParentReference) bool {
	return k8s.IsGatewayKind(ref) && gw.Name == string(ref.Name)
}
