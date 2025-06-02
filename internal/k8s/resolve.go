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

package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func ResolveGateway(
	ctx context.Context,
	routeMeta metav1.ObjectMeta,
	ref gwAPIv1.ParentReference,
) (*gwAPIv1.Gateway, error) {
	ns := ref.Namespace
	if ns == nil {
		routeNS := gwAPIv1.Namespace(routeMeta.Namespace)
		ns = &routeNS
	}

	gw := &gwAPIv1.Gateway{}
	key := client.ObjectKey{Namespace: string(*ns), Name: string(ref.Name)}
	if err := GetClient().Get(ctx, key, gw); err != nil {
		return nil, err
	}
	return gw, nil
}

func ResolveRouteHostnames(ctx context.Context, route *gwAPIv1.HTTPRoute) []string {
	hostnames := []string{}
	for i := range route.Spec.ParentRefs {
		refStatus := route.Status.Parents[i]
		ref := refStatus.ParentRef
		gw, err := ResolveGateway(ctx, route.ObjectMeta, ref)
		if err != nil {
			continue
		}
		parentHosts := GetHTTPHosts(route, gw, ref)
		hostnames = append(hostnames, parentHosts...)
	}
	return hostnames
}
