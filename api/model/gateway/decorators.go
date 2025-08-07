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

package gateway

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type GatewayClass struct {
	Object *gwAPIv1.GatewayClass
}

func WrapGatewayClass(gwc *gwAPIv1.GatewayClass) *GatewayClass {
	return &GatewayClass{Object: gwc}
}

func (gwc *GatewayClass) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(gwc.Object.Status.Conditions)
}

func (gwc *GatewayClass) SetConditions(conditions []metav1.Condition) {
	gwc.Object.Status.Conditions = conditions
}

type Gateway struct {
	Object *gwAPIv1.Gateway
}

func WrapGateway(gateway *gwAPIv1.Gateway) *Gateway {
	return &Gateway{Object: gateway}
}

func (gw *Gateway) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(gw.Object.Status.Conditions)
}

func (gw *Gateway) SetConditions(conditions []metav1.Condition) {
	gw.Object.Status.Conditions = conditions
}

type ListenerStatus struct {
	Object *gwAPIv1.ListenerStatus
}

func WrapListenerStatus(lst *gwAPIv1.ListenerStatus) *ListenerStatus {
	return &ListenerStatus{Object: lst}
}

func (lst *ListenerStatus) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(lst.Object.Conditions)
}

func (lst *ListenerStatus) SetConditions(conditions []metav1.Condition) {
	lst.Object.Conditions = conditions
}

type RouteParentStatus struct {
	Object *gwAPIv1.RouteParentStatus
}

func InitRouteParentStatus(ref gwAPIv1.ParentReference) *RouteParentStatus {
	return &RouteParentStatus{
		Object: &gwAPIv1.RouteParentStatus{
			ParentRef:      ref,
			ControllerName: core.GraviteeGatewayClassController,
			Conditions:     []metav1.Condition{},
		},
	}
}

func WrapRouteParentStatus(status *gwAPIv1.RouteParentStatus) *RouteParentStatus {
	return &RouteParentStatus{
		Object: status,
	}
}

func (st *RouteParentStatus) GetConditions() map[string]metav1.Condition {
	return utils.MapConditions(st.Object.Conditions)
}

func (st *RouteParentStatus) SetConditions(conditions []metav1.Condition) {
	st.Object.Conditions = conditions
}
