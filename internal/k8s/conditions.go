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
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	ConditionAccepted     = "Accepted"
	ConditionProgrammed   = "Programmed"
	ConditionConflicted   = "Conflicted"
	ConditionPending      = "Pending"
	ConditionResolvedRefs = "ResolvedRefs"

	ConditionStatusTrue  = "True"
	ConditionStatusFalse = "False"
)

type ConditionBuilder struct {
	condition *metav1.Condition
}

func NewResolvedRefsConditionBuilder(generation int64) *ConditionBuilder {
	return NewConditionBuilder(ConditionResolvedRefs).
		ObservedGeneration(generation).
		Status(ConditionStatusTrue).
		Reason(ConditionResolvedRefs)
}

func NewAcceptedConditionBuilder(generation int64) *ConditionBuilder {
	return NewConditionBuilder(ConditionAccepted).
		ObservedGeneration(generation).
		Status(ConditionStatusFalse)
}

func NewGatewayProgrammedConditionBuilder(generation int64) *ConditionBuilder {
	return NewConditionBuilder(ConditionProgrammed).
		ObservedGeneration(generation).
		Status(ConditionStatusFalse).
		Reason(string(gwAPIv1.GatewayReasonPending))
}

func NewListenerProgrammedConditionBuilder(generation int64) *ConditionBuilder {
	return NewConditionBuilder(ConditionProgrammed).
		ObservedGeneration(generation).
		Status(ConditionStatusFalse).
		Reason(string(gwAPIv1.ListenerReasonPending))
}

func NewListenerConflictedConditionBuilder(generation int64) *ConditionBuilder {
	return NewConditionBuilder(ConditionConflicted).
		ObservedGeneration(generation).
		Status(ConditionStatusFalse).
		Reason(string(gwAPIv1.ListenerReasonNoConflicts))
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
		Reason(ConditionAccepted).
		Status(metav1.ConditionTrue).
		Message(msg)
}

func (b *ConditionBuilder) Program(msg string) *ConditionBuilder {
	return b.
		Reason(ConditionProgrammed).
		Status(metav1.ConditionTrue).
		Message(msg)
}

func (b *ConditionBuilder) ResolveRefs(msg string) *ConditionBuilder {
	return b.
		Reason(ConditionResolvedRefs).
		Status(metav1.ConditionTrue).
		Message(msg)
}

func (b *ConditionBuilder) RejectInvalidRouteKinds(msg string) *ConditionBuilder {
	return b.
		Reason(string(gwAPIv1.ListenerReasonInvalidRouteKinds)).
		Status(metav1.ConditionFalse).
		Message(msg)
}

func (b *ConditionBuilder) RejectInvalidCertificateRef(msg string) *ConditionBuilder {
	return b.
		Reason(string(gwAPIv1.ListenerReasonInvalidCertificateRef)).
		Status(metav1.ConditionFalse).
		Message(msg)
}

func (b *ConditionBuilder) RejectInvalidParameters(msg string) *ConditionBuilder {
	return b.
		Reason(string(gwAPIv1.GatewayClassReasonInvalidParameters)).
		Status(metav1.ConditionFalse).
		Message(msg)
}

func (b *ConditionBuilder) RejectUnsupportedProtocol(msg string) *ConditionBuilder {
	return b.
		Reason(string(gwAPIv1.ListenerReasonUnsupportedProtocol)).
		Status(metav1.ConditionFalse).
		Message(msg)
}

func (b *ConditionBuilder) RejectListenersNotValid(msg string) *ConditionBuilder {
	return b.
		Reason(string(gwAPIv1.GatewayReasonListenersNotValid)).
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

func GetCondition(obj core.ConditionAware, conditionType string) *metav1.Condition {
	conditions := obj.GetConditions()
	condition, ok := conditions[conditionType]
	if !ok {
		return nil
	}
	return &condition
}

func IsConflicted(obj core.ConditionAware) bool {
	conflicted := GetCondition(obj, ConditionConflicted)
	return conflicted != nil && conflicted.Status == ConditionStatusTrue
}

func IsAccepted(obj core.ConditionAware) bool {
	accepted := GetCondition(obj, ConditionAccepted)
	return accepted != nil && accepted.Status == ConditionStatusTrue
}

func HasUnresolvedRefs(obj core.ConditionAware) bool {
	resolvedRefs := GetCondition(obj, ConditionResolvedRefs)
	return resolvedRefs != nil && resolvedRefs.Status == ConditionStatusFalse
}
