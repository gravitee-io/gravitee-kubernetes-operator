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
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func buildRequestFilters(rule gwAPIv1.HTTPRouteRule) []*v4.FlowStep {
	steps := []*v4.FlowStep{}
	for _, f := range rule.Filters {
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
