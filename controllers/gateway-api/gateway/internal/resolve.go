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
	"encoding/pem"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	kErrors "k8s.io/apimachinery/pkg/api/errors"

	coreV1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func Resolve(
	ctx context.Context,
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
) error {
	for i, listener := range gw.Object.Spec.Listeners {
		conditionBuilder := k8s.NewResolvedRefsConditionBuilder(gw.Object.Generation)

		if err := resolveTLS(ctx, gw, conditionBuilder, listener); err != nil {
			return err
		}

		status := gateway.WrapListenerStatus(&gw.Object.Status.Listeners[i])

		resolveRouteKinds(listener, status.Object, conditionBuilder)

		if k8s.IsKafkaListener(listener) {
			resolveKafkaListener(params, conditionBuilder)
		}

		k8s.SetCondition(status, conditionBuilder.Build())

		if httpRoutesCount, err := countAttachedHTTPRoutes(ctx, gw.Object, listener); err != nil {
			return err
		} else {
			status.Object.AttachedRoutes = httpRoutesCount
		}

		if kafkaRoutesCount, err := countAttachedKafkaRoutes(ctx, gw.Object, listener); err != nil {
			return err
		} else {
			status.Object.AttachedRoutes += kafkaRoutesCount
		}
	}
	return nil
}

func resolveKafkaListener(
	params *v1alpha1.GatewayClassParameters,
	builder *k8s.ConditionBuilder,
) {
	if !k8s.HasKafkaEnabled(params) {
		builder.
			RejectInvalidRouteKinds(
				"Kafka is not enabled on the gateway class",
			)
	}
}

func resolveTLS(
	ctx context.Context,
	gw *gateway.Gateway,
	builder *k8s.ConditionBuilder,
	listener gwAPIv1.Listener,
) error {
	if listener.TLS == nil {
		return nil
	}

	if len(listener.TLS.CertificateRefs) != 1 {
		builder.RejectTooManyCertificateRefs("listener should have exacty one certificate reference")
		return nil
	}

	for _, ref := range listener.TLS.CertificateRefs {
		if err := resolveTLSRef(ctx, gw, builder, ref); err != nil {
			return err
		}
	}

	return nil
}

func resolveTLSRef(
	ctx context.Context,
	gw *gateway.Gateway,
	builder *k8s.ConditionBuilder,
	ref gwAPIv1.SecretObjectReference,
) error {
	if hasInvalidSecretGroup(ref) {
		builder.RejectInvalidCertificateRef(
			fmt.Sprintf("TLS certificate group [%s] is invalid", *ref.Group),
		)
		return nil
	}

	if hasInvalidSecretKind(ref) {
		builder.RejectInvalidCertificateRef(
			fmt.Sprintf("TLS certificate kind [%s] is invalid", *ref.Kind),
		)
		return nil
	}

	ns := ref.Namespace
	if ns == nil {
		gwNs := gwAPIv1.Namespace(gw.Object.Namespace)
		ns = &gwNs
	}

	objectRef := gwAPIv1.ObjectReference{
		Name:      ref.Name,
		Group:     *ref.Group,
		Kind:      *ref.Kind,
		Namespace: ref.Namespace,
	}

	if granted, err := k8s.IsGrantedReference(ctx, gw.Object, objectRef); err != nil {
		return err
	} else if !granted {
		builder.RejectListenerRefNotPermitted(
			"Illegal TLS certificate reference",
		)
	}

	key := client.ObjectKey{Namespace: string(*ns), Name: string(ref.Name)}
	secret := &coreV1.Secret{}
	if err := k8s.GetClient().Get(ctx, key, secret); client.IgnoreNotFound(err) != nil {
		return err
	} else if kErrors.IsNotFound(err) {
		builder.RejectInvalidCertificateRef(
			fmt.Sprintf("TLS certificate secret [%s] could not be resolved", key.String()),
		)
		return nil
	}

	if isMalformedSecret(secret) {
		builder.RejectInvalidCertificateRef(
			fmt.Sprintf("TLS certificate secret [%s] is malformed", key.String()),
		)
	}

	return nil
}

func resolveRouteKinds(
	listener gwAPIv1.Listener,
	status *gwAPIv1.ListenerStatus,
	builder *k8s.ConditionBuilder,
) {
	if len(status.SupportedKinds) == 0 {
		builder.
			RejectInvalidRouteKinds("at least one supported kind is expected")
	}

	if len(listener.AllowedRoutes.Kinds) == 0 {
		return
	}

	validRoutesCount := len(listener.AllowedRoutes.Kinds)
	for _, k := range listener.AllowedRoutes.Kinds {
		if !isValidRouteKind(status.SupportedKinds, k) {
			validRoutesCount -= 1

			builder.
				RejectInvalidRouteKinds(
					fmt.Sprintf("route kind [%s] is not supported", k.Kind),
				)
		}
	}

	if validRoutesCount == 0 {
		status.SupportedKinds = []gwAPIv1.RouteGroupKind{}
	}
}

func isValidRouteKind(
	supportedKinds []gwAPIv1.RouteGroupKind,
	kind gwAPIv1.RouteGroupKind,
) bool {
	if kind.Group == nil {
		return false
	}
	if *kind.Group != k8s.GwAPIv1Group && *kind.Group != k8s.GraviteeGroup {
		return false
	}
	for _, k := range supportedKinds {
		if k.Kind == kind.Kind {
			return true
		}
	}
	return false
}

func countAttachedHTTPRoutes(
	ctx context.Context,
	gw *gwAPIv1.Gateway,
	listener gwAPIv1.Listener,
) (int32, error) {
	var count int32 = 0

	if !k8s.HasHTTPSupport(listener) {
		return 0, nil
	}

	opts := &client.ListOptions{}
	routesList := &gwAPIv1.HTTPRouteList{}
	if err := k8s.GetClient().List(ctx, routesList, opts); err != nil {
		return 0, err
	}
	for _, route := range routesList.Items {
		if k8s.IsAttachedHTTPRoute(gw, listener, route) {
			count += 1
		}
	}
	return count, nil
}

func countAttachedKafkaRoutes(
	ctx context.Context,
	gw *gwAPIv1.Gateway,
	listener gwAPIv1.Listener,
) (int32, error) {
	var count int32 = 0

	if !k8s.IsKafkaListener(listener) {
		return 0, nil
	}

	opts := &client.ListOptions{}
	routesList := &v1alpha1.KafkaRouteList{}
	if err := k8s.GetClient().List(ctx, routesList, opts); err != nil {
		return 0, err
	}
	for _, route := range routesList.Items {
		if k8s.IsAttachedKafkaRoute(gw, listener, route) {
			count += 1
		}
	}
	return count, nil
}

func hasInvalidSecretGroup(ref gwAPIv1.SecretObjectReference) bool {
	if ref.Group == nil {
		return false
	}
	if *ref.Group == "" {
		return false
	}
	return *ref.Group != gwAPIv1.Group(coreV1.SchemeGroupVersion.Group)
}

func hasInvalidSecretKind(ref gwAPIv1.SecretObjectReference) bool {
	if ref.Kind == nil {
		return false
	}
	if *ref.Kind == "" {
		return false
	}
	return *ref.Kind != gwAPIv1.Kind("Secret")
}

func isMalformedSecret(secret *coreV1.Secret) bool {
	var ok bool
	var crt, key []byte
	if crt, ok = secret.Data["tls.crt"]; !ok {
		return true
	}
	if key, ok = secret.Data["tls.key"]; !ok {
		return true
	}
	if p, _ := pem.Decode(crt); p == nil {
		return true
	}
	if p, _ := pem.Decode(key); p == nil {
		return true
	}
	return false
}
