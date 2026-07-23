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
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func WatchConfigMaps() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(requestsFromConfigMap)
}

func requestsFromConfigMap(ctx context.Context, obj client.Object) []reconcile.Request {
	cm, ok := obj.(*coreV1.ConfigMap)
	if !ok {
		return nil
	}

	policies := &gwAPIv1.BackendTLSPolicyList{}
	if err := k8s.GetClient().List(ctx, policies, client.InNamespace(cm.Namespace)); err != nil {
		return nil
	}

	var requests []reconcile.Request
	for i := range policies.Items {
		policy := &policies.Items[i]
		if referencesConfigMap(policy, cm.Name) {
			requests = append(requests, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: policy.Namespace,
					Name:      policy.Name,
				},
			})
		}
	}
	return requests
}

func referencesConfigMap(policy *gwAPIv1.BackendTLSPolicy, cmName string) bool {
	for _, ref := range policy.Spec.Validation.CACertificateRefs {
		if (ref.Kind == "" || ref.Kind == "ConfigMap") && string(ref.Name) == cmName {
			return true
		}
	}
	return false
}
