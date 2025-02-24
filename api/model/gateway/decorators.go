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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type GatewayClass struct {
	Object *gAPIv1.GatewayClass
}

func NewGatewayClass(gwc *gAPIv1.GatewayClass) *GatewayClass {
	return &GatewayClass{Object: gwc}
}

func (gwc *GatewayClass) GetConditions() map[string]metav1.Condition {
	return mapConditions(gwc.Object.Status.Conditions)
}

func (gwc *GatewayClass) SetConditions(conditions []metav1.Condition) {
	gwc.Object.Status.Conditions = conditions
}

type Gateway struct {
	Object *gAPIv1.Gateway
}

func NewGateway(gateway *gAPIv1.Gateway) *Gateway {
	return &Gateway{Object: gateway}
}

func (gw *Gateway) GetConditions() map[string]metav1.Condition {
	return mapConditions(gw.Object.Status.Conditions)
}

func (gw *Gateway) SetConditions(conditions []metav1.Condition) {
	gw.Object.Status.Conditions = conditions
}

type ListenerStatus struct {
	Object *gAPIv1.ListenerStatus
}

func NewListenerStatus(lst *gAPIv1.ListenerStatus) *ListenerStatus {
	return &ListenerStatus{Object: lst}
}

func (lst *ListenerStatus) GetConditions() map[string]metav1.Condition {
	return mapConditions(lst.Object.Conditions)
}

func (lst *ListenerStatus) SetConditions(conditions []metav1.Condition) {
	lst.Object.Conditions = conditions
}

func mapConditions(conditionsSlice []metav1.Condition) map[string]metav1.Condition {
	conditions := make(map[string]metav1.Condition)
	for _, condition := range conditionsSlice {
		conditions[condition.Type] = condition
	}
	return conditions
}
