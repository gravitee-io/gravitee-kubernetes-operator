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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func DetectConflicts(gw *gateway.Gateway) error {
	specLen := len(gw.Object.Spec.Listeners)
	statusLen := len(gw.Object.Status.Listeners)
	if statusLen != specLen {
		return fmt.Errorf("listener status array length (%d) does not match spec listeners length (%d)", statusLen, specLen)
	}

	conflicts := detectListenerConflicts(gw.Object.Spec.Listeners)
	setConflictConditions(gw, conflicts)

	return nil
}

func detectListenerConflicts(listeners []gwAPIv1.Listener) map[int]string {
	conflicts := make(map[int]string)

	for i, l1 := range listeners {
		for j, l2 := range listeners {
			if i == j {
				continue
			}

			if hasProtocolConflict(l1, l2) {
				setConflictIfNotExists(conflicts, i, string(gwAPIv1.ListenerReasonProtocolConflict))
				setConflictIfNotExists(conflicts, j, string(gwAPIv1.ListenerReasonProtocolConflict))
				continue
			}

			if hasHostnameConflict(l1, l2) {
				setConflictIfNotExists(conflicts, i, string(gwAPIv1.ListenerReasonHostnameConflict))
				setConflictIfNotExists(conflicts, j, string(gwAPIv1.ListenerReasonHostnameConflict))
				continue
			}

			if hasKafkaConflict(l1, l2) {
				setConflictIfNotExists(conflicts, i, k8s.ListenerReasonKafkaConflict)
				setConflictIfNotExists(conflicts, j, k8s.ListenerReasonKafkaConflict)
				continue
			}
		}
	}

	return conflicts
}

func hasProtocolConflict(l1, l2 gwAPIv1.Listener) bool {
	return l1.Port == l2.Port && l1.Protocol != l2.Protocol
}

func hasHostnameConflict(l1, l2 gwAPIv1.Listener) bool {
	return l1.Hostname != nil && l2.Hostname != nil && *l1.Hostname == *l2.Hostname
}

func hasKafkaConflict(l1, l2 gwAPIv1.Listener) bool {
	return k8s.IsKafkaListener(l1) && k8s.IsKafkaListener(l2)
}

func setConflictIfNotExists(conflicts map[int]string, index int, reason string) {
	if _, exists := conflicts[index]; !exists {
		conflicts[index] = reason
	}
}

func setConflictConditions(gw *gateway.Gateway, conflicts map[int]string) {
	for i := range gw.Object.Spec.Listeners {
		listenerStatus := gateway.WrapListenerStatus(&gw.Object.Status.Listeners[i])
		condition := k8s.NewListenerConflictedConditionBuilder(gw.Object.Generation)

		if reason, hasConflict := conflicts[i]; hasConflict {
			condition.Status(k8s.ConditionStatusTrue)
			condition.Reason(reason)
		} else {
			condition.Status(k8s.ConditionStatusFalse)
		}

		k8s.SetCondition(listenerStatus, condition.Build())
	}
}
