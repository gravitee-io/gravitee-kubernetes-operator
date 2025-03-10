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
)

func Accept(gw *gateway.Gateway) {
	acceptListeners(gw)

	accepted := k8s.NewAcceptedConditionBuilder(gw.Object.Generation).Accept("gateway is accepted")
	for i := range gw.Object.Status.Listeners {
		status := gateway.WrapListenerStatus(&gw.Object.Status.Listeners[i])
		if !k8s.IsAccepted(status) {
			accepted.RejectListenersNotValid(fmt.Sprintf("listener [%d] is not valid", i))
		}
		if k8s.IsConflicted(status) {
			accepted.RejectListenersNotValid(fmt.Sprintf("listener [%d] conflicts", i))
		}
		if k8s.HasUnresolvedRefs(status) {
			accepted.RejectListenersNotValid(fmt.Sprintf("listener [%d] has unresolved refs", i))
		}
	}

	k8s.SetCondition(gw, accepted.Build())
}

func acceptListeners(gw *gateway.Gateway) {
	for i, l := range gw.Object.Spec.Listeners {
		listenerStatus := gateway.WrapListenerStatus(&gw.Object.Status.Listeners[i])
		condition := k8s.NewAcceptedConditionBuilder(gw.Object.Generation)
		switch {
		case !k8s.SupportedGwAPIProtocols.Has(l.Protocol):
			condition.RejectUnsupportedProtocol(fmt.Sprintf("protocol [%s] is not supported", l.Protocol))
		case k8s.IsConflicted(listenerStatus):
			condition.RejectListenersNotValid("listener conflicts")
		case k8s.HasUnresolvedRefs(listenerStatus):
			condition.RejectListenersNotValid("listener has unresolved refs")
		default:
			condition.Accept("listener has been accepted")
		}
		k8s.SetCondition(listenerStatus, condition.Build())
	}
}
