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

package mapper

import (
	"fmt"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var serviceURIPattern = "http://%s.%s.svc.cluster.local:%d"

func buildEndpointGroups(route *gwAPIv1.HTTPRoute) []*v4.EndpointGroup {
	groups := []*v4.EndpointGroup{}
	for i, rule := range route.Spec.Rules {
		groups = append(groups, buildEndpointGroup(rule, i, route.Namespace))
	}
	return groups
}

func buildEndpointGroup(
	rule gwAPIv1.HTTPRouteRule,
	index int,
	ns string,
) *v4.EndpointGroup {
	endpointGroup := v4.NewHttpEndpointGroup(
		buildEndpointGroupName(index),
	)
	endpointGroup.Endpoints = buildEndpoints(rule, ns)
	return endpointGroup
}

func buildEndpointGroupName(index int) string {
	return fmt.Sprintf("rule-%d", index)
}

func buildEndpoints(rule gwAPIv1.HTTPRouteRule, ns string) []*v4.Endpoint {
	endpoints := []*v4.Endpoint{}
	for i, ref := range rule.BackendRefs {
		endpoints = append(endpoints, buildEndpoint(ref, i, ns))
	}
	return endpoints
}

func buildEndpoint(
	ref gwAPIv1.HTTPBackendRef,
	index int,
	ns string,
) *v4.Endpoint {
	endpoint := v4.NewHttpEndpoint(
		fmt.Sprintf("backend-%d", index),
	)
	endpoint.Config.Object["target"] = buildEndpointTarget(ref, ns)
	return endpoint
}

func buildEndpointTarget(
	ref gwAPIv1.HTTPBackendRef,
	ns string,
) string {
	return fmt.Sprintf(serviceURIPattern, ref.Name, ns, *ref.Port)
}
