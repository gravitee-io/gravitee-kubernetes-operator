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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func DetectConflicts(gw *gateway.Gateway) {
	for i, l1 := range gw.Object.Spec.Listeners {
		for j, l2 := range gw.Object.Spec.Listeners {
			if i == j {
				continue
			}
			condition := k8s.NewListenerConflictedConditionBuilder(gw.Object.Generation)
			if l1.Port == l2.Port && l1.Protocol != l2.Protocol {
				condition.Status(k8s.ConditionStatusTrue)
				condition.Reason(string(gwAPIv1.ListenerReasonProtocolConflict))
				break
			}
			if l1.Hostname != nil && l2.Hostname != nil && *l1.Hostname == *l2.Hostname {
				condition.Status(k8s.ConditionStatusTrue)
				condition.Reason(string(gwAPIv1.ListenerReasonHostnameConflict))
				break
			}
			listenerStatus := gw.Object.Status.Listeners[i]
			k8s.SetCondition(&gateway.ListenerStatus{Object: &listenerStatus}, condition.Build())
		}
	}
}
