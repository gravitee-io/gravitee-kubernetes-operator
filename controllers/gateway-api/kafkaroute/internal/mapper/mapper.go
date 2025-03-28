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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var (
	keyLessSecurity = v4.NewPlanSecurity("KEY_LESS")
)

func Map(route *v1alpha1.KafkaRoute, gw *gwAPIv1.Gateway) *v1alpha1.ApiV4Definition {
	api := newAPI(route.ObjectMeta)
	api.Spec = MapSpec(route, gw)
	return api
}

func MapSpec(route *v1alpha1.KafkaRoute, gw *gwAPIv1.Gateway) v1alpha1.ApiV4DefinitionSpec {
	spec := newAPISpec(route.ObjectMeta)
	spec.Listeners = buildListeners(route, gw)
	spec.EndpointGroups = buildEndpointGroups(route)
	spec.Flows = buildFlows(route)
	spec.Tags = buildTags(route)
	return spec
}

func newAPI(meta metav1.ObjectMeta) *v1alpha1.ApiV4Definition {
	return &v1alpha1.ApiV4Definition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ApiV4Definition",
			APIVersion: "v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{Name: meta.Name, Namespace: meta.Namespace},
		Spec: v1alpha1.ApiV4DefinitionSpec{
			Api: v4.Api{
				Type: "NATIVE",
				Plans: &map[string]*v4.Plan{
					"default": newKeyLessPlan(),
				},
				ApiBase: &base.ApiBase{
					Name:    buildAPIName(meta),
					Version: "v1alpha1",
				},
				DefinitionContext: &v4.DefinitionContext{
					Origin:   v4.OriginKubernetes,
					SyncFrom: v4.OriginKubernetes,
				},
			},
		},
	}
}

func newAPISpec(meta metav1.ObjectMeta) v1alpha1.ApiV4DefinitionSpec {
	return v1alpha1.ApiV4DefinitionSpec{
		Api: v4.Api{
			Type: "NATIVE",
			Plans: &map[string]*v4.Plan{
				"default": newKeyLessPlan(),
			},
			ApiBase: &base.ApiBase{
				Name:    buildAPIName(meta),
				Version: "v1alpha1",
			},
			DefinitionContext: &v4.DefinitionContext{
				Origin:   v4.OriginKubernetes,
				SyncFrom: v4.OriginKubernetes,
			},
		},
	}
}

func buildAPIName(meta metav1.ObjectMeta) string {
	return meta.Name + "-" + meta.Namespace
}

func buildTags(route *v1alpha1.KafkaRoute) []string {
	tags := []string{}
	for _, ref := range route.Spec.ParentRefs {
		tags = append(tags, buildTag(route, ref))
	}
	return tags
}

func buildTag(route *v1alpha1.KafkaRoute, ref gwAPIv1.ParentReference) string {
	ns := ref.Namespace
	if ns == nil {
		routeNS := gwAPIv1.Namespace(route.Namespace)
		ns = &routeNS
	}
	return k8s.BuildTag(string(*ns), string(ref.Name))
}

func newKeyLessPlan() *v4.Plan {
	plan := v4.NewPlan().WithSecurity(&keyLessSecurity)
	plan.Status = base.PublishedPlanStatus
	return plan
}
