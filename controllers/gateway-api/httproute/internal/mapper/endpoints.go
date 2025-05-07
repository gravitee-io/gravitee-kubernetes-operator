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

var (
	serviceURIPattern           = "http://%s.%s.svc.cluster.local:%d"
	defaultEndpointWeight int32 = 1
)

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
	backendRefs := getActiveBackendRefs(rule.BackendRefs)

	if len(backendRefs) > 1 {
		endpointGroup.LoadBalancer = v4.NewLoadBalancer(v4.WeightedRoundRobin)
	}
	endpointGroup.Endpoints = buildEndpoints(backendRefs, ns)
	return endpointGroup
}

func buildEndpointGroupName(index int) string {
	return fmt.Sprintf("endpoints-%d", index)
}

func buildEndpoints(
	backendRefs []gwAPIv1.HTTPBackendRef,
	namespace string,
) []*v4.Endpoint {
	endpoints := []*v4.Endpoint{}
	for i, ref := range backendRefs {
		endpoints = append(endpoints, buildEndpoint(ref, i, namespace))
	}
	return endpoints
}

func buildEndpoint(
	ref gwAPIv1.HTTPBackendRef,
	index int,
	namespace string,
) *v4.Endpoint {
	endpoint := v4.NewHttpEndpoint(
		fmt.Sprintf("backend-%d", index),
	)
	endpoint.Weight = ref.Weight
	endpoint.Config.Object["target"] = buildEndpointTarget(ref, namespace)
	return endpoint
}

func buildEndpointTarget(
	ref gwAPIv1.HTTPBackendRef,
	namespace string,
) string {
	return fmt.Sprintf(serviceURIPattern, ref.Name, namespace, *ref.Port)
}

// If several backends are provided, skip backends with a weight defined to 0.
// See https://gateway-api.sigs.k8s.io/guides/traffic-splitting
func getActiveBackendRefs(refs []gwAPIv1.HTTPBackendRef) []gwAPIv1.HTTPBackendRef {
	if len(refs) == 1 {
		return refs
	}
	activeRefs := []gwAPIv1.HTTPBackendRef{}
	for _, ref := range refs {
		if ref.Weight == nil {
			ref.Weight = &defaultEndpointWeight
		}
		if *ref.Weight > 0 {
			activeRefs = append(activeRefs, ref)
		}
	}
	return activeRefs
}
