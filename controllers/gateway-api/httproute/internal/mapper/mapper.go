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
	"sort"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	rootPath = "/"
)

var (
	keyLessSecurity = v4.NewPlanSecurity("KEY_LESS")
)

func Map(ctx context.Context, route *gwAPIv1.HTTPRoute) (*v1alpha1.ApiV4Definition, error) {
	api := newAPI(route)
	spec, err := MapSpec(ctx, route)
	if err != nil {
		return nil, err
	}
	api.Spec = spec

	return api, nil
}

func MapSpec(ctx context.Context, route *gwAPIv1.HTTPRoute) (v1alpha1.ApiV4DefinitionSpec, error) {
	return MapSpecWithPrefix(ctx, route, "")
}

func MapSpecWithPrefix(ctx context.Context, route *gwAPIv1.HTTPRoute, prefix string) (v1alpha1.ApiV4DefinitionSpec, error) {
	spec := newAPISpec(route)
	listeners, err := buildListeners(ctx, route)
	if err != nil {
		return v1alpha1.ApiV4DefinitionSpec{}, err
	}
	spec.Listeners = listeners
	endpointGroups, err := buildEndpointGroupsWithPrefix(ctx, route, prefix)
	if err != nil {
		return v1alpha1.ApiV4DefinitionSpec{}, err
	}
	spec.EndpointGroups = endpointGroups
	flows, err := buildFlowsWithPrefix(ctx, route, prefix)
	if err != nil {
		return v1alpha1.ApiV4DefinitionSpec{}, err
	}
	spec.Flows = flows
	spec.Tags = buildTags(route)
	return spec, nil
}

func MergeSpecs(specs []v1alpha1.ApiV4DefinitionSpec) v1alpha1.ApiV4DefinitionSpec {
	if len(specs) == 0 {
		return v1alpha1.ApiV4DefinitionSpec{}
	}

	merged := specs[0]

	for i := 1; i < len(specs); i++ {
		spec := specs[i]
		merged.Flows = append(merged.Flows, spec.Flows...)
		merged.EndpointGroups = append(merged.EndpointGroups, spec.EndpointGroups...)
		merged.Tags = mergeTags(merged.Tags, spec.Tags)
		if len(spec.Listeners) > 0 && len(merged.Listeners) > 0 {
			mergeListenerPaths(merged.Listeners[0], spec.Listeners[0])
		}
	}

	sortFlowsByWeight(merged.Flows)

	return merged
}

func sortFlowsByWeight(flows []*v4.Flow) {
	sort.SliceStable(flows, func(i, j int) bool {
		return flowWeight(flows[i]) > flowWeight(flows[j])
	})
}

func flowWeight(flow *v4.Flow) int {
	if flow == nil || len(flow.Selectors) < 2 {
		return 0
	}

	weight := 1 // base path condition

	if methods, ok := flow.Selectors[0].Get("methods").([]base.HttpMethod); ok && len(methods) > 0 {
		weight += methodPrecedenceBonus
	}

	cond := flow.Selectors[1].GetString("condition")
	weight += strings.Count(cond, "#request.headers") * headerWeight
	weight += strings.Count(cond, "#request.params") * queryParamWeight

	return weight
}

func mergeTags(a, b []string) []string {
	seen := make(map[string]struct{}, len(a))
	for _, t := range a {
		seen[t] = struct{}{}
	}
	for _, t := range b {
		if _, ok := seen[t]; !ok {
			a = append(a, t)
			seen[t] = struct{}{}
		}
	}
	return a
}

func mergeListenerPaths(dst, src *v4.GenericListener) {
	dstListener := dst.ToListener()
	srcListener := src.ToListener()
	httpDst, ok1 := dstListener.(*v4.HttpListener)
	httpSrc, ok2 := srcListener.(*v4.HttpListener)
	if !ok1 || !ok2 {
		return
	}
	existing := make(map[string]struct{})
	for _, p := range httpDst.Paths {
		existing[p.Host+"|"+p.Path] = struct{}{}
	}
	for _, p := range httpSrc.Paths {
		key := p.Host + "|" + p.Path
		if _, ok := existing[key]; !ok {
			httpDst.Paths = append(httpDst.Paths, p)
			existing[key] = struct{}{}
		}
	}
	*dst = *v4.ToGenericListener(httpDst)
}

func newAPI(route *gwAPIv1.HTTPRoute) *v1alpha1.ApiV4Definition {
	return &v1alpha1.ApiV4Definition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      route.Name,
			Namespace: route.Namespace,
		},
		Spec: newAPISpec(route),
	}
}

func newAPISpec(route *gwAPIv1.HTTPRoute) v1alpha1.ApiV4DefinitionSpec {
	return v1alpha1.ApiV4DefinitionSpec{
		Api: v4.Api{
			V4BaseApi: &v4.V4BaseApi{
				Type: "PROXY",
				FlowExecution: &v4.FlowExecution{
					Mode:          v4.FlowModeBestMatch,
					MatchRequired: true,
				},
				ApiBase: &base.ApiBase{
					Name:    buildAPIName(route),
					Version: "v1alpha1",
				},
				DefinitionContext: &v4.DefinitionContext{
					Origin:   v4.OriginKubernetes,
					SyncFrom: v4.OriginKubernetes,
				},
			},
			Plans: &map[string]*v4.Plan{
				"default": newKeyLessPlan(),
			},
		},
	}
}

func buildAPIName(route *gwAPIv1.HTTPRoute) string {
	return route.Name + "-" + route.Namespace
}

func newKeyLessPlan() *v4.Plan {
	return v4.NewPlan().WithSecurity(&keyLessSecurity)
}

func buildTags(route *gwAPIv1.HTTPRoute) []string {
	tags := []string{}
	for i, ref := range route.Spec.ParentRefs {
		if i >= len(route.Status.Parents) {
			tags = append(tags, buildTag(route, ref))
			continue
		}
		routeParentStatus := gateway.WrapRouteParentStatus(&route.Status.Parents[i])
		if k8s.IsAccepted(routeParentStatus) {
			tags = append(tags, buildTag(route, ref))
		}
	}
	return tags
}

func buildTag(route *gwAPIv1.HTTPRoute, ref gwAPIv1.ParentReference) string {
	ns := ref.Namespace
	if ns == nil {
		routeNS := gwAPIv1.Namespace(route.Namespace)
		ns = &routeNS
	}
	return k8s.BuildTag(string(*ns), string(ref.Name))
}
