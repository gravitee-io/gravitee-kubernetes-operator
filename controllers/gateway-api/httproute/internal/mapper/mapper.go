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

func Map(ctx context.Context, route *gwAPIv1.HTTPRoute) *v1alpha1.ApiV4Definition {
	api := newAPI(route)
	api.Spec = MapSpec(ctx, route)
	return api
}

func MapSpec(ctx context.Context, route *gwAPIv1.HTTPRoute) v1alpha1.ApiV4DefinitionSpec {
	spec := newAPISpec(route)
	spec.Listeners = buildListeners(ctx, route)
	spec.EndpointGroups = buildEndpointGroups(route)
	spec.Flows = buildFlows(route)
	spec.Tags = buildTags(route)
	return spec
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
			Type: "PROXY",
			Plans: &map[string]*v4.Plan{
				"default": newKeyLessPlan(),
			},
			FlowExecution: &v4.FlowExecution{
				Mode:          v4.FlowModeDefault,
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
	}
}

func buildAPIName(route *gwAPIv1.HTTPRoute) string {
	return route.Name + "-" + route.Namespace
}

func newKeyLessPlan() *v4.Plan {
	plan := v4.NewPlan().WithSecurity(&keyLessSecurity)
	plan.Status = base.PublishedPlanStatus
	return plan
}

func buildTags(route *gwAPIv1.HTTPRoute) []string {
	tags := []string{}
	for i, ref := range route.Spec.ParentRefs {
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
