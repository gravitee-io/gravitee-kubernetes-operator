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

func (d *GatewayClass) GetConditions() map[string]metav1.Condition {
	conditions := make(map[string]metav1.Condition)
	for _, condition := range d.Object.Status.Conditions {
		conditions[condition.Type] = condition
	}
	return conditions
}

func (d *GatewayClass) SetConditions(conditions []metav1.Condition) {
	d.Object.Status.Conditions = conditions
}
