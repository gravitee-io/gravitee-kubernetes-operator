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
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func Init(gw *gateway.Gateway) {
	k8s.SetCondition(
		gw,
		k8s.NewGatewayProgrammedConditionBuilder(gw.Object.Generation).
			Build(),
	)
	initListeners(gw)
}

func Program(
	ctx context.Context,
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
) error {
	if err := k8s.DeployGateway(ctx, gw, params); err != nil {
		return err
	}

	k8s.SetCondition(
		gw,
		k8s.NewGatewayProgrammedConditionBuilder(gw.Object.Generation).
			Program("all listeners have been programmed").
			Build(),
	)

	programListeners(gw)

	return nil
}

func programListeners(gw *gateway.Gateway) {
	listeners := gw.Object.Status.Listeners
	for i := range listeners {
		status := gateway.NewListenerStatus(&listeners[i])
		k8s.SetCondition(
			status,
			k8s.NewGatewayProgrammedConditionBuilder(gw.Object.Generation).
				Program("listener is programmed").
				Build(),
		)
	}
}

func initListeners(gw *gateway.Gateway) {
	statuses := make([]gwAPIv1.ListenerStatus, len(gw.Object.Spec.Listeners))
	for i, l := range gw.Object.Spec.Listeners {
		status := gateway.NewListenerStatus(
			&gwAPIv1.ListenerStatus{
				Name:           l.Name,
				SupportedKinds: k8s.GetSupportedRouteKinds(l),
			},
		)
		k8s.SetCondition(status,
			k8s.NewListenerProgrammedConditionBuilder(gw.Object.Generation).
				Build(),
		)
		statuses[i] = *status.Object
	}
	gw.Object.Status.Listeners = statuses
}
