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

package k8s

import (
	"maps"
	"slices"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type ConditionBuilder struct {
	condition *metav1.Condition
}

func NewAcceptedConditionBuilder(generation int64) *ConditionBuilder {
	return NewConditionBuilder(string(gAPIv1.GatewayClassConditionStatusAccepted)).
		ObservedGeneration(generation)
}

func NewConditionBuilder(cType string) *ConditionBuilder {
	return &ConditionBuilder{
		condition: &metav1.Condition{
			Type: cType,
		},
	}
}

func (b *ConditionBuilder) Accept(msg string) *ConditionBuilder {
	return b.
		Reason("Accepted").
		Status(metav1.ConditionTrue).
		Message(msg)
}

func (b *ConditionBuilder) RejectInvalidParameters(msg string) *ConditionBuilder {
	return b.
		Reason(string(gAPIv1.GatewayClassReasonInvalidParameters)).
		Status(metav1.ConditionFalse).
		Message(msg)
}

func (b *ConditionBuilder) Reason(reason string) *ConditionBuilder {
	b.condition.Reason = reason
	return b
}

func (b *ConditionBuilder) Status(status metav1.ConditionStatus) *ConditionBuilder {
	b.condition.Status = status
	return b
}

func (b *ConditionBuilder) Message(msg string) *ConditionBuilder {
	b.condition.Message = msg
	return b
}

func (b *ConditionBuilder) ObservedGeneration(gen int64) *ConditionBuilder {
	b.condition.ObservedGeneration = gen
	return b
}

func (b *ConditionBuilder) Build() *metav1.Condition {
	b.condition.LastTransitionTime = metav1.Now()
	return b.condition
}

func SetCondition(obj core.ConditionAware, condition *metav1.Condition) {
	if condition != nil {
		conditions := obj.GetConditions()
		conditions[condition.Type] = *condition
		obj.SetConditions(slices.Collect(maps.Values(conditions)))
	}
}
