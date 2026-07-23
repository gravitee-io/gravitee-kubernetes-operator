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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func Reconcile(ctx context.Context, policy *gwAPIv1.BackendTLSPolicy) ([]gwAPIv1.PolicyAncestorStatus, error) {
	gateways, err := findAncestorGateways(ctx, policy)
	if err != nil {
		return nil, err
	}

	refsValid, refsCondition := resolveRefs(ctx, policy)

	conflicting := detectConflicts(ctx, policy)

	ancestors := make([]gwAPIv1.PolicyAncestorStatus, 0, len(gateways))
	for _, gw := range gateways {
		ns := gwAPIv1.Namespace(gw.Namespace)
		ancestor := gwAPIv1.PolicyAncestorStatus{
			AncestorRef: gwAPIv1.ParentReference{
				Group:     ptrTo(gwAPIv1.Group(gwAPIv1.GroupName)),
				Kind:      ptrTo(gwAPIv1.Kind("Gateway")),
				Namespace: &ns,
				Name:      gwAPIv1.ObjectName(gw.Name),
			},
			ControllerName: gwAPIv1.GatewayController(core.GraviteeGatewayClassController),
			Conditions:     buildConditions(policy, refsValid, refsCondition, conflicting),
		}
		ancestors = append(ancestors, ancestor)
	}

	return ancestors, nil
}

func buildConditions(
	policy *gwAPIv1.BackendTLSPolicy,
	refsValid bool,
	refsCondition metav1.Condition,
	conflicting bool,
) []metav1.Condition {
	refsCondition.ObservedGeneration = policy.Generation

	acceptedCond := metav1.Condition{
		Type:               string(gwAPIv1.PolicyConditionAccepted),
		ObservedGeneration: policy.Generation,
		LastTransitionTime: metav1.Now(),
	}

	switch {
	case conflicting:
		acceptedCond.Status = metav1.ConditionFalse
		acceptedCond.Reason = string(gwAPIv1.PolicyReasonConflicted)
		acceptedCond.Message = "Another BackendTLSPolicy targeting the same service takes precedence"
	case !refsValid:
		acceptedCond.Status = metav1.ConditionFalse
		acceptedCond.Reason = string(gwAPIv1.BackendTLSPolicyReasonNoValidCACertificate)
		acceptedCond.Message = "No valid CA certificate reference"
	default:
		acceptedCond.Status = metav1.ConditionTrue
		acceptedCond.Reason = string(gwAPIv1.PolicyReasonAccepted)
		acceptedCond.Message = "Policy has been accepted"
	}

	return []metav1.Condition{acceptedCond, refsCondition}
}

func resolveRefs(ctx context.Context, policy *gwAPIv1.BackendTLSPolicy) (bool, metav1.Condition) {
	condition := metav1.Condition{
		Type:               string(gwAPIv1.BackendTLSPolicyConditionResolvedRefs),
		LastTransitionTime: metav1.Now(),
		ObservedGeneration: policy.Generation,
	}

	validation := policy.Spec.Validation

	if len(validation.CACertificateRefs) == 0 {
		if validation.WellKnownCACertificates != nil && *validation.WellKnownCACertificates != "" {
			condition.Status = metav1.ConditionFalse
			condition.Reason = string(gwAPIv1.PolicyReasonInvalid)
			condition.Message = "WellKnownCACertificates is not supported"
			return false, condition
		}
		condition.Status = metav1.ConditionFalse
		condition.Reason = string(gwAPIv1.BackendTLSPolicyReasonInvalidCACertificateRef)
		condition.Message = "No CA certificate references specified"
		return false, condition
	}

	for _, ref := range validation.CACertificateRefs {
		if ref.Kind != "" && ref.Kind != "ConfigMap" {
			condition.Status = metav1.ConditionFalse
			condition.Reason = string(gwAPIv1.BackendTLSPolicyReasonInvalidKind)
			condition.Message = fmt.Sprintf("Unsupported CACertificateRef kind: %s", ref.Kind)
			return false, condition
		}

		if ref.Group != "" && ref.Group != "core" {
			condition.Status = metav1.ConditionFalse
			condition.Reason = string(gwAPIv1.BackendTLSPolicyReasonInvalidCACertificateRef)
			condition.Message = fmt.Sprintf("Unsupported CACertificateRef group: %s", ref.Group)
			return false, condition
		}

		cm := &coreV1.ConfigMap{}
		key := client.ObjectKey{
			Namespace: policy.Namespace,
			Name:      string(ref.Name),
		}

		if err := k8s.GetClient().Get(ctx, key, cm); err != nil {
			condition.Status = metav1.ConditionFalse
			condition.Reason = string(gwAPIv1.BackendTLSPolicyReasonInvalidCACertificateRef)
			condition.Message = fmt.Sprintf("ConfigMap %q not found", key.String())
			return false, condition
		}

		if _, ok := cm.Data["ca.crt"]; !ok {
			condition.Status = metav1.ConditionFalse
			condition.Reason = string(gwAPIv1.BackendTLSPolicyReasonInvalidCACertificateRef)
			condition.Message = fmt.Sprintf("ConfigMap %q does not contain key ca.crt", key.String())
			return false, condition
		}

		if cm.Data["ca.crt"] == "" {
			condition.Status = metav1.ConditionFalse
			condition.Reason = string(gwAPIv1.BackendTLSPolicyReasonInvalidCACertificateRef)
			condition.Message = fmt.Sprintf("ConfigMap %q has empty ca.crt key", key.String())
			return false, condition
		}
	}

	condition.Status = metav1.ConditionTrue
	condition.Reason = string(gwAPIv1.BackendTLSPolicyReasonResolvedRefs)
	condition.Message = "All CA certificate references are valid"
	return true, condition
}

func detectConflicts(ctx context.Context, policy *gwAPIv1.BackendTLSPolicy) bool {
	allPolicies := &gwAPIv1.BackendTLSPolicyList{}
	if err := k8s.GetClient().List(ctx, allPolicies, client.InNamespace(policy.Namespace)); err != nil {
		return false
	}

	for _, targetRef := range policy.Spec.TargetRefs {
		for i := range allPolicies.Items {
			other := &allPolicies.Items[i]
			if other.Name == policy.Name {
				continue
			}
			if conflictsOn(policy, other, targetRef) {
				if hasPrecedence(other, policy) {
					return true
				}
			}
		}
	}
	return false
}

func conflictsOn(
	_ *gwAPIv1.BackendTLSPolicy,
	other *gwAPIv1.BackendTLSPolicy,
	targetRef gwAPIv1.LocalPolicyTargetReferenceWithSectionName,
) bool {
	for _, otherRef := range other.Spec.TargetRefs {
		if otherRef.Group != targetRef.Group || otherRef.Kind != targetRef.Kind || otherRef.Name != targetRef.Name {
			continue
		}
		if targetRef.SectionName == nil && otherRef.SectionName == nil {
			return true
		}
		if targetRef.SectionName != nil && otherRef.SectionName != nil &&
			*targetRef.SectionName == *otherRef.SectionName {
			return true
		}
	}
	return false
}

func hasPrecedence(a, b *gwAPIv1.BackendTLSPolicy) bool {
	if !a.CreationTimestamp.Equal(&b.CreationTimestamp) {
		return a.CreationTimestamp.Before(&b.CreationTimestamp)
	}
	aKey := a.Namespace + "/" + a.Name
	bKey := b.Namespace + "/" + b.Name
	return aKey < bKey
}

func findAncestorGateways(ctx context.Context, policy *gwAPIv1.BackendTLSPolicy) ([]client.ObjectKey, error) {
	targetServices := make(map[string]bool)
	for _, ref := range policy.Spec.TargetRefs {
		if ref.Kind == "" || ref.Kind == "Service" {
			targetServices[string(ref.Name)] = true
		}
	}

	routes := &gwAPIv1.HTTPRouteList{}
	if err := k8s.GetClient().List(ctx, routes, client.InNamespace(policy.Namespace)); err != nil {
		return nil, err
	}

	gwSet := make(map[client.ObjectKey]bool)
	for i := range routes.Items {
		route := &routes.Items[i]
		if !routeReferencesServices(route, targetServices) {
			continue
		}
		for _, parentRef := range route.Spec.ParentRefs {
			if parentRef.Kind != nil && *parentRef.Kind != "Gateway" {
				continue
			}
			ns := policy.Namespace
			if parentRef.Namespace != nil {
				ns = string(*parentRef.Namespace)
			}
			gwKey := client.ObjectKey{Namespace: ns, Name: string(parentRef.Name)}
			if isOurGateway(ctx, gwKey) {
				gwSet[gwKey] = true
			}
		}
	}

	gateways := make([]client.ObjectKey, 0, len(gwSet))
	for gw := range gwSet {
		gateways = append(gateways, gw)
	}
	return gateways, nil
}

func routeReferencesServices(route *gwAPIv1.HTTPRoute, services map[string]bool) bool {
	for _, rule := range route.Spec.Rules {
		for _, ref := range rule.BackendRefs {
			if ref.Kind != nil && *ref.Kind != "Service" {
				continue
			}
			if services[string(ref.Name)] {
				return true
			}
		}
	}
	return false
}

func isOurGateway(ctx context.Context, gwKey client.ObjectKey) bool {
	gw := &gwAPIv1.Gateway{}
	if err := k8s.GetClient().Get(ctx, gwKey, gw); err != nil {
		return false
	}

	gwc := &gwAPIv1.GatewayClass{}
	if err := k8s.GetClient().Get(ctx, client.ObjectKey{Name: string(gw.Spec.GatewayClassName)}, gwc); err != nil {
		return false
	}

	return string(gwc.Spec.ControllerName) == core.GraviteeGatewayClassController
}

func ptrTo[T any](v T) *T {
	return &v
}
