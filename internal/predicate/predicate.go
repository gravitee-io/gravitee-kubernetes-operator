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

//nolint:cyclop // as the number of types increases the package bigger and bigger switches
package predicate

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	corev1 "k8s.io/api/core/v1"
	netV1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type LastSpecHashPredicate struct {
	predicate.Funcs
}

// Create returns true if the Create event should be processed.
func (LastSpecHashPredicate) Create(e event.CreateEvent) bool {
	if e.Object.GetDeletionTimestamp() != nil {
		return true
	}

	if shouldAlwaysReconcile() {
		return true
	}

	switch t := e.Object.(type) {
	case *v1alpha1.ApiDefinition:
		return e.Object.GetAnnotations()[core.LastSpecHashAnnotation] != hash.Calculate(&t.Spec) ||
			t.Status.ProcessingStatus != core.ProcessingStatusCompleted
	case *v1alpha1.ApiV4Definition:
		return e.Object.GetAnnotations()[core.LastSpecHashAnnotation] != hash.Calculate(&t.Spec) ||
			t.Status.ProcessingStatus != core.ProcessingStatusCompleted
	case *v1alpha1.ManagementContext:
		return e.Object.GetAnnotations()[core.LastSpecHashAnnotation] != hash.Calculate(&t.Spec)
	case *v1alpha1.ApiResource:
		return e.Object.GetAnnotations()[core.LastSpecHashAnnotation] != hash.Calculate(&t.Spec)
	case *v1alpha1.Application:
		return e.Object.GetAnnotations()[core.LastSpecHashAnnotation] != hash.Calculate(&t.Spec) ||
			t.Status.ProcessingStatus != core.ProcessingStatusCompleted
	case *v1alpha1.Subscription:
		return e.Object.GetAnnotations()[core.LastSpecHashAnnotation] != hash.Calculate(&t.Spec)
	case *v1alpha1.SharedPolicyGroup:
		return e.Object.GetAnnotations()[core.LastSpecHashAnnotation] != hash.Calculate(&t.Spec)
	case *v1alpha1.Group:
		return e.Object.GetAnnotations()[core.LastSpecHashAnnotation] != hash.Calculate(&t.Spec)
	case *v1alpha1.Notification:
		return e.Object.GetAnnotations()[core.LastSpecHashAnnotation] != hash.Calculate(&t.Spec)
	case *netV1.Ingress:
		return e.Object.GetAnnotations()[core.LastSpecHashAnnotation] != hash.Calculate(&t.Spec)
	case *corev1.Secret:
		return e.Object.GetAnnotations()[core.LastSpecHashAnnotation] != hash.Calculate(&t.Data)
	case *v1alpha1.GatewayClassParameters:
		return e.Object.GetAnnotations()[core.LastSpecHashAnnotation] != hash.Calculate(&t.Spec)
	case *gwAPIv1.GatewayClass:
		return e.Object.GetAnnotations()[core.LastSpecHashAnnotation] != hash.Calculate(&t.Spec)
	case *gwAPIv1.HTTPRoute:
		return e.Object.GetAnnotations()[core.LastSpecHashAnnotation] != hash.Calculate(&t.Spec)
	case *v1alpha1.KafkaRoute:
		return e.Object.GetAnnotations()[core.LastSpecHashAnnotation] != hash.Calculate(&t.Spec)
	default:
		return false
	}
}

// Update implements default UpdateEvent filter for validating spec hash change.
//
//nolint:funlen // keep
func (LastSpecHashPredicate) Update(e event.UpdateEvent) bool {
	if e.ObjectOld == nil || e.ObjectNew == nil {
		return false
	}

	if e.ObjectNew.GetDeletionTimestamp() != nil {
		switch no := e.ObjectNew.(type) {
		case *corev1.Secret, *corev1.ConfigMap:
			return false
		default:
			return no.GetDeletionTimestamp() != nil
		}
	}

	switch no := e.ObjectNew.(type) {
	case *v1alpha1.ApiDefinition:
		oo, _ := e.ObjectOld.(*v1alpha1.ApiDefinition)
		return hash.Calculate(&no.Spec) != hash.Calculate(&oo.Spec)
	case *v1alpha1.ApiV4Definition:
		oo, _ := e.ObjectOld.(*v1alpha1.ApiV4Definition)
		return hash.Calculate(&no.Spec) != hash.Calculate(&oo.Spec)
	case *v1alpha1.ManagementContext:
		oo, _ := e.ObjectOld.(*v1alpha1.ManagementContext)
		return hash.Calculate(&no.Spec) != hash.Calculate(&oo.Spec)
	case *v1alpha1.ApiResource:
		oo, _ := e.ObjectOld.(*v1alpha1.ApiResource)
		return hash.Calculate(&no.Spec) != hash.Calculate(&oo.Spec)
	case *v1alpha1.Application:
		oo, _ := e.ObjectOld.(*v1alpha1.Application)
		return hash.Calculate(&no.Spec) != hash.Calculate(&oo.Spec)
	case *v1alpha1.Subscription:
		oo, _ := e.ObjectOld.(*v1alpha1.Subscription)
		return hash.Calculate(&no.Spec) != hash.Calculate(&oo.Spec)
	case *v1alpha1.SharedPolicyGroup:
		oo, _ := e.ObjectOld.(*v1alpha1.SharedPolicyGroup)
		return hash.Calculate(&no.Spec) != hash.Calculate(&oo.Spec)
	case *v1alpha1.Group:
		oo, _ := e.ObjectOld.(*v1alpha1.Group)
		return hash.Calculate(&no.Spec) != hash.Calculate(&oo.Spec)
	case *v1alpha1.Notification:
		oo, _ := e.ObjectOld.(*v1alpha1.Notification)
		return hash.Calculate(&no.Spec) != hash.Calculate(&oo.Spec)
	case *netV1.Ingress:
		oo, _ := e.ObjectOld.(*netV1.Ingress)
		return hash.Calculate(&no.Spec) != hash.Calculate(&oo.Spec)
	case *corev1.Secret:
		oo, _ := e.ObjectOld.(*corev1.Secret)
		return hash.Calculate(&no.Data) != hash.Calculate(&oo.Data)
	case *corev1.ConfigMap:
		oo, _ := e.ObjectOld.(*corev1.ConfigMap)
		return hash.Calculate(&no.Data) != hash.Calculate(&oo.Data) ||
			hash.Calculate(&no.BinaryData) != hash.Calculate(&oo.BinaryData)
	case *v1alpha1.GatewayClassParameters:
		oo, _ := e.ObjectOld.(*v1alpha1.GatewayClassParameters)
		return hash.Calculate(&no.Spec) != hash.Calculate(&oo.Spec)
	case *gwAPIv1.GatewayClass:
		oo, _ := e.ObjectOld.(*gwAPIv1.GatewayClass)
		return hash.Calculate(&no.Spec) != hash.Calculate(&oo.Spec)
	case *gwAPIv1.HTTPRoute:
		oo, _ := e.ObjectOld.(*gwAPIv1.HTTPRoute)
		return hash.Calculate(&no.Spec) != hash.Calculate(&oo.Spec)
	case *v1alpha1.KafkaRoute:
		oo, _ := e.ObjectOld.(*v1alpha1.KafkaRoute)
		return hash.Calculate(&no.Spec) != hash.Calculate(&oo.Spec)
	default:
		return false
	}
}

// Delete returns true if the Delete event should be processed.
func (LastSpecHashPredicate) Delete(e event.DeleteEvent) bool {
	switch t := e.Object.(type) {
	case *corev1.Secret, *corev1.ConfigMap:
		return false
	default:
		return t.GetDeletionTimestamp() != nil
	}
}

func shouldAlwaysReconcile() bool {
	return env.Config.ReconcileStrategy == "always"
}
