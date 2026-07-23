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

package mapper

import (
	"context"
	"strings"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	coreV1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type backendTLSConfig struct {
	hostname string
	caData   string
	valid    bool
}

func lookupBackendTLSPolicy(
	ctx context.Context,
	serviceName string,
	namespace string,
	portName string,
) *backendTLSConfig {
	policies := &gwAPIv1.BackendTLSPolicyList{}
	if err := k8s.GetClient().List(ctx, policies, client.InNamespace(namespace)); err != nil {
		return nil
	}

	var winningPolicy *gwAPIv1.BackendTLSPolicy
	winningHasSection := false
	for i := range policies.Items {
		policy := &policies.Items[i]
		hasSection := matchesServiceWithSection(policy, serviceName, portName)
		hasGlobal := !hasSection && matchesService(policy, serviceName, portName)
		if !hasSection && !hasGlobal {
			continue
		}
		if winningPolicy == nil {
			winningPolicy = policy
			winningHasSection = hasSection
			continue
		}
		if hasSection && !winningHasSection {
			winningPolicy = policy
			winningHasSection = true
		} else if hasSection == winningHasSection && hasTLSPrecedence(policy, winningPolicy) {
			winningPolicy = policy
		}
	}

	if winningPolicy == nil {
		return nil
	}

	if isConflicted(ctx, winningPolicy, namespace) {
		return nil
	}

	caData := resolveCAData(ctx, winningPolicy, namespace)
	return &backendTLSConfig{
		hostname: string(winningPolicy.Spec.Validation.Hostname),
		caData:   caData,
		valid:    caData != "",
	}
}

func matchesServiceWithSection(policy *gwAPIv1.BackendTLSPolicy, serviceName, portName string) bool {
	for _, ref := range policy.Spec.TargetRefs {
		if ref.Kind != "" && ref.Kind != "Service" {
			continue
		}
		if string(ref.Name) != serviceName {
			continue
		}
		if ref.SectionName != nil && string(*ref.SectionName) == portName {
			return true
		}
	}
	return false
}

func matchesService(policy *gwAPIv1.BackendTLSPolicy, serviceName, portName string) bool {
	for _, ref := range policy.Spec.TargetRefs {
		if ref.Kind != "" && ref.Kind != "Service" {
			continue
		}
		if string(ref.Name) != serviceName {
			continue
		}
		if ref.SectionName != nil && string(*ref.SectionName) != portName {
			continue
		}
		return true
	}
	return false
}

func hasTLSPrecedence(a, b *gwAPIv1.BackendTLSPolicy) bool {
	if !a.CreationTimestamp.Equal(&b.CreationTimestamp) {
		return a.CreationTimestamp.Before(&b.CreationTimestamp)
	}
	aKey := a.Namespace + "/" + a.Name
	bKey := b.Namespace + "/" + b.Name
	return aKey < bKey
}

func isConflicted(ctx context.Context, policy *gwAPIv1.BackendTLSPolicy, namespace string) bool {
	policies := &gwAPIv1.BackendTLSPolicyList{}
	if err := k8s.GetClient().List(ctx, policies, client.InNamespace(namespace)); err != nil {
		return false
	}

	for i := range policies.Items {
		other := &policies.Items[i]
		if other.Name == policy.Name {
			continue
		}
		for _, ref := range policy.Spec.TargetRefs {
			for _, otherRef := range other.Spec.TargetRefs {
				if otherRef.Group != ref.Group || otherRef.Kind != ref.Kind || otherRef.Name != ref.Name {
					continue
				}
				sameSection := (ref.SectionName == nil && otherRef.SectionName == nil) ||
					(ref.SectionName != nil && otherRef.SectionName != nil && *ref.SectionName == *otherRef.SectionName)
				if sameSection && hasTLSPrecedence(other, policy) {
					return true
				}
			}
		}
	}
	return false
}

func resolveCAData(ctx context.Context, policy *gwAPIv1.BackendTLSPolicy, namespace string) string {
	var allCerts []string

	for _, ref := range policy.Spec.Validation.CACertificateRefs {
		if ref.Kind != "" && ref.Kind != "ConfigMap" {
			return ""
		}

		cm := &coreV1.ConfigMap{}
		key := client.ObjectKey{Namespace: namespace, Name: string(ref.Name)}
		if err := k8s.GetClient().Get(ctx, key, cm); err != nil {
			return ""
		}

		caCrt, ok := cm.Data["ca.crt"]
		if !ok || caCrt == "" {
			return ""
		}
		allCerts = append(allCerts, caCrt)
	}

	if len(allCerts) == 0 {
		return ""
	}
	return strings.Join(allCerts, "\n")
}

func applyBackendTLS(ctx context.Context, ep *v4.Endpoint, tlsCfg *backendTLSConfig, route *gwAPIv1.HTTPRoute) {
	if tlsCfg == nil {
		return
	}

	switchToHTTPS(ep)

	sslConfig := utils.NewGenericStringMap()
	sslConfig.Put("hostnameVerifier", true)
	sslConfig.Put("trustAll", false)

	if tlsCfg.valid && tlsCfg.caData != "" {
		sslConfig.Put("trustStore", map[string]any{
			"type":    "PEM",
			"content": tlsCfg.caData,
		})
	}

	if certPEM, keyPEM := resolveGatewayClientCert(ctx, route); certPEM != "" && keyPEM != "" {
		sslConfig.Put("keyStore", map[string]any{
			"type":        "PEM",
			"certContent": certPEM,
			"keyContent":  keyPEM,
		})
	}

	ep.ConfigOverride.Put("ssl", sslConfig)

	httpConfig := getOrCreateHTTPConfig(ep)
	httpConfig.Put("propagateClientHost", false)

	var headers []map[string]string
	if existing, ok := ep.ConfigOverride.Get("headers").([]map[string]string); ok {
		headers = existing
	}
	headers = append(headers, map[string]string{"name": "Host", "value": tlsCfg.hostname})
	ep.ConfigOverride.Put("headers", headers)
}

func resolveGatewayClientCert(ctx context.Context, route *gwAPIv1.HTTPRoute) (string, string) {
	for _, parentRef := range route.Spec.ParentRefs {
		if !k8s.IsGatewayKind(parentRef) {
			continue
		}

		gw, err := k8s.ResolveGateway(ctx, route.ObjectMeta, parentRef)
		if err != nil || gw == nil {
			continue
		}

		if gw.Spec.TLS == nil || gw.Spec.TLS.Backend == nil || gw.Spec.TLS.Backend.ClientCertificateRef == nil {
			continue
		}

		ref := *gw.Spec.TLS.Backend.ClientCertificateRef
		ns := ref.Namespace
		if ns == nil {
			gwNs := gwAPIv1.Namespace(gw.Namespace)
			ns = &gwNs
		}

		secret := &coreV1.Secret{}
		key := client.ObjectKey{Namespace: string(*ns), Name: string(ref.Name)}
		if err := k8s.GetClient().Get(ctx, key, secret); err != nil {
			continue
		}

		certPEM := string(secret.Data["tls.crt"])
		keyPEM := string(secret.Data["tls.key"])
		if certPEM != "" && keyPEM != "" {
			return certPEM, keyPEM
		}
	}
	return "", ""
}

func switchToHTTPS(ep *v4.Endpoint) {
	if target, ok := ep.Config.Get("target").(string); ok {
		ep.Config.Object["target"] = strings.Replace(target, "http://", "https://", 1)
	}
}
