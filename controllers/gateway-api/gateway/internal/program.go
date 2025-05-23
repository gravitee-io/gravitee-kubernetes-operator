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
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var (
	IPAddressType       = gwAPIv1.IPAddressType
	HostnameAddressType = gwAPIv1.HostnameAddressType
)

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

	return setGatewayAddresses(ctx, gw)
}

func programListeners(gw *gateway.Gateway) {
	listeners := gw.Object.Status.Listeners
	for i := range listeners {
		status := gateway.WrapListenerStatus(&listeners[i])
		k8s.SetCondition(
			status,
			k8s.NewGatewayProgrammedConditionBuilder(gw.Object.Generation).
				Program("listener is programmed").
				Build(),
		)
	}
}

func setGatewayAddresses(ctx context.Context, gw *gateway.Gateway) error {
	svcList := &coreV1.ServiceList{}
	if err := k8s.GetClient().List(
		ctx,
		svcList,
		&client.ListOptions{
			Namespace:     gw.Object.Namespace,
			LabelSelector: labels.SelectorFromSet(k8s.GwAPIv1GatewayLabels(gw.Object.Name)),
		},
	); err != nil {
		return err
	}

	gw.Object.Status.Addresses = getGatewayServiceAddresses(gw, svcList.Items)

	return nil
}

func getGatewayServiceAddresses(gw *gateway.Gateway, svcList []coreV1.Service) []gwAPIv1.GatewayStatusAddress {
	for i := range svcList {
		if k8s.IsGatewayDependent(gw, &svcList[i]) {
			return getIngressAddresses(svcList[i])
		}
	}
	return []gwAPIv1.GatewayStatusAddress{}
}

func getIngressAddresses(svc coreV1.Service) []gwAPIv1.GatewayStatusAddress {
	if svc.Spec.Type == coreV1.ServiceTypeLoadBalancer {
		return getLBAddresses(svc.Status.LoadBalancer)
	}
	return []gwAPIv1.GatewayStatusAddress{newIPAddress(svc.Spec.ClusterIP)}
}

func getLBAddresses(lb coreV1.LoadBalancerStatus) []gwAPIv1.GatewayStatusAddress {
	addrs := make([]gwAPIv1.GatewayStatusAddress, 0)
	for _, addr := range lb.Ingress {
		if addr.IP != "" {
			addrs = append(addrs, newIPAddress(addr.IP))
		}
		if addr.Hostname != "" {
			addrs = append(addrs, newHostnameAddress(addr.Hostname))
		}
	}
	return addrs
}

func newIPAddress(ip string) gwAPIv1.GatewayStatusAddress {
	return gwAPIv1.GatewayStatusAddress{
		Value: ip,
		Type:  &IPAddressType,
	}
}

func newHostnameAddress(hostname string) gwAPIv1.GatewayStatusAddress {
	return gwAPIv1.GatewayStatusAddress{
		Value: hostname,
		Type:  &HostnameAddressType,
	}
}
