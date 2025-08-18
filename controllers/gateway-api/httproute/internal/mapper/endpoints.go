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
	"context"
	"fmt"
	"strings"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var (
	serviceURIPattern           = "http://%s.%s.svc.cluster.local:%d%s"
	discardURI                  = "http://not.a.service.cluster.local"
	defaultEndpointWeight int32 = 1
)

func buildEndpointGroups(ctx context.Context, route *gwAPIv1.HTTPRoute) ([]*v4.EndpointGroup, error) {
	groups := []*v4.EndpointGroup{}
	for ruleIndex, rule := range route.Spec.Rules {
		for matchIndex, match := range rule.Matches {
			epg, err := buildEndpointGroup(
				ctx,
				route,
				rule,
				ruleIndex,
				match,
				matchIndex,
			)
			if err != nil {
				return nil, err
			}
			groups = append(groups, epg)
		}
	}
	return groups, nil
}

func buildEndpointGroup(
	ctx context.Context,
	route *gwAPIv1.HTTPRoute,
	rule gwAPIv1.HTTPRouteRule,
	ruleIndex int,
	match gwAPIv1.HTTPRouteMatch,
	matchIndex int,
) (*v4.EndpointGroup, error) {
	endpointGroup := v4.NewHttpEndpointGroup(
		buildEndpointGroupName(ruleIndex, matchIndex),
	)

	backendRefs := getActiveBackendRefs(rule.BackendRefs)

	if len(backendRefs) > 1 {
		endpointGroup.LoadBalancer = v4.NewLoadBalancer(v4.WeightedRoundRobin)
	}
	eps, err := buildEndpoints(ctx, route, match, matchIndex, backendRefs)
	if err != nil {
		return nil, err
	}
	endpointGroup.Endpoints = eps
	return endpointGroup, nil
}

func buildEndpointGroupName(ruleIndex, matchIndex int) string {
	return fmt.Sprintf("endpoints-%d-%d", ruleIndex, matchIndex)
}

func buildEndpoints(
	ctx context.Context,
	route *gwAPIv1.HTTPRoute,
	match gwAPIv1.HTTPRouteMatch,
	matchIndex int,
	backendRefs []gwAPIv1.HTTPBackendRef,
) ([]*v4.Endpoint, error) {
	endpoints := []*v4.Endpoint{}
	if len(backendRefs) == 0 {
		// in case of HTTP redirect, there is no backend ref
		return append(endpoints, buildDummyEndpoint()), nil
	}
	for backendIndex, backendRef := range backendRefs {
		objectRef := gwAPIv1.ObjectReference{
			Name:      backendRef.Name,
			Group:     *backendRef.Group,
			Kind:      *backendRef.Kind,
			Namespace: backendRef.Namespace,
		}
		if granted, err := k8s.IsGrantedReference(ctx, route, objectRef); err != nil {
			return nil, err
		} else if !granted {
			endpoints = append(endpoints, buildDummyEndpoint())
		} else {
			endpoints = append(
				endpoints,
				buildEndpoint(
					backendRef,
					backendIndex,
					match,
					matchIndex,
					k8s.GetRefNs(route, backendRef.Namespace),
				),
			)
		}
	}
	return endpoints, nil
}

func buildDummyEndpoint() *v4.Endpoint {
	ep := buildEndpoint(gwAPIv1.HTTPBackendRef{}, 0, gwAPIv1.HTTPRouteMatch{}, 0, "")
	ep.Secondary = true
	return ep
}

func buildEndpoint(
	backendRef gwAPIv1.HTTPBackendRef,
	backendIndex int,
	match gwAPIv1.HTTPRouteMatch,
	matchIndex int,
	namespace string,
) *v4.Endpoint {
	endpoint := v4.NewHttpEndpoint(
		fmt.Sprintf("backend-%d-match-%d", backendIndex, matchIndex),
	)

	endpoint.Weight = backendRef.Weight
	endpoint.Config.Object["target"] = buildEndpointTarget(match, backendRef, namespace)
	endpoint.Inherit = false

	httpConfig := utils.NewGenericStringMap()
	httpConfig.Put("propagateClientHost", true)
	endpoint.ConfigOverride.Put("http", httpConfig)

	return endpoint
}

func buildEndpointTarget(
	match gwAPIv1.HTTPRouteMatch,
	backendRef gwAPIv1.HTTPBackendRef,
	namespace string,
) string {
	if !k8s.IsServiceKind(backendRef.BackendObjectReference) {
		return discardURI
	}
	return fmt.Sprintf(
		serviceURIPattern,
		backendRef.Name,
		namespace,
		*backendRef.Port,
		getEndpointPath(match),
	)
}

func getEndpointPath(match gwAPIv1.HTTPRouteMatch) string {
	if match.Path == nil {
		return ""
	}
	return strings.TrimSuffix(*match.Path.Value, "/")
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
