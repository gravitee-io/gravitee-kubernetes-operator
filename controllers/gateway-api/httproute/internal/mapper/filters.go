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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	trafficShadowingPolicyName = "traffic-shadowing"
)

func buildRequestFilters(ctx context.Context, route *gwAPIv1.HTTPRoute, rule gwAPIv1.HTTPRouteRule) []*v4.FlowStep {
	steps := []*v4.FlowStep{}
	for _, f := range rule.Filters {
		if f.RequestRedirect != nil {
			steps = append(steps, buildHTTPRedirect(ctx, route, *f.RequestRedirect))
		}
		if f.RequestHeaderModifier != nil {
			steps = append(steps, buildHeaderTransformer(*f.RequestHeaderModifier))
		}
		if f.RequestMirror != nil {
			steps = append(steps, buildRequestMirror(route, *f.RequestMirror))
		}
	}
	return steps
}

func buildBackendRequestFilters(rule gwAPIv1.HTTPRouteRule) []*v4.FlowStep {
	if len(rule.BackendRefs) != 1 {
		return nil
	}
	steps := []*v4.FlowStep{}
	for _, f := range rule.BackendRefs[0].Filters {
		if f.RequestHeaderModifier != nil {
			steps = append(steps, buildHeaderTransformer(*f.RequestHeaderModifier))
		}
	}
	return steps
}

func buildResponseFilters(rule gwAPIv1.HTTPRouteRule) []*v4.FlowStep {
	steps := []*v4.FlowStep{}
	for _, f := range rule.Filters {
		if f.ResponseHeaderModifier != nil {
			steps = append(steps, buildHeaderTransformer(*f.ResponseHeaderModifier))
		}
	}
	return steps
}

func buildRequestMirror(route *gwAPIv1.HTTPRoute, mirror gwAPIv1.HTTPRequestMirrorFilter) *v4.FlowStep {
	ns := route.Namespace
	if mirror.BackendRef.Namespace != nil {
		ns = string(*mirror.BackendRef.Namespace)
	}

	port := int32(80)
	if mirror.BackendRef.Port != nil {
		port = int32(*mirror.BackendRef.Port)
	}

	target := fmt.Sprintf(serviceURIPattern, mirror.BackendRef.Name, ns, port, "")

	policyName := trafficShadowingPolicyName
	configuration := utils.NewGenericStringMap().
		Put("target", target)

	return v4.NewFlowStep(base.FlowStep{
		Policy:        &policyName,
		Enabled:       true,
		Configuration: configuration,
	})
}
