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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func Init(route *gwAPIv1.HTTPRoute) {
	parents := make([]gwAPIv1.RouteParentStatus, len(route.Spec.ParentRefs))
	for i := range route.Spec.ParentRefs {
		parents[i] = gwAPIv1.RouteParentStatus{
			ParentRef:      route.Spec.ParentRefs[i],
			ControllerName: core.GraviteeGatewayClassController,
		}
	}
}
