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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	GwAPIv1HTTPRouteKind = "HTTPRoute"
	GwAPIv1GatewayKind   = "Gateway"
	CoreV1ServiceKind    = "Service"
	GwAPIv1APIVersion    = "gateway.networking.k8s.io/v1"

	GraviteeAPIVersion     = "gravitee.io/v1alpha1"
	GraviteeKafkaRouteKind = "KafkaRoute"
)

var (
	GwAPIv1Group  = gwAPIv1.Group(gwAPIv1.GroupVersion.Group)
	GraviteeGroup = gwAPIv1.Group(v1alpha1.GroupVersion.Group)
)

var SupportedGwAPIProtocols = sets.New(
	gwAPIv1.HTTPProtocolType,
	gwAPIv1.HTTPSProtocolType,
	gwAPIv1.TLSProtocolType,
)

var ProtocolToRouteKinds = map[gwAPIv1.ProtocolType][]gwAPIv1.RouteGroupKind{
	gwAPIv1.HTTPProtocolType: {
		{
			Group: &GwAPIv1Group,
			Kind:  GwAPIv1HTTPRouteKind,
		},
	},
	gwAPIv1.HTTPSProtocolType: {
		{
			Group: &GwAPIv1Group,
			Kind:  GwAPIv1HTTPRouteKind,
		},
	},
	gwAPIv1.TLSProtocolType: {
		{
			Group: &GraviteeGroup,
			Kind:  GraviteeKafkaRouteKind,
		},
	},
}

var ProtocolToServerType = map[gwAPIv1.ProtocolType]string{
	gwAPIv1.HTTPProtocolType:  "http",
	gwAPIv1.HTTPSProtocolType: "http",
}

func GetSupportedRouteKinds(listener gwAPIv1.Listener) []gwAPIv1.RouteGroupKind {
	if kinds, ok := ProtocolToRouteKinds[listener.Protocol]; ok {
		return kinds
	}
	return []gwAPIv1.RouteGroupKind{}
}

func IsListenerRef(
	gw *gwAPIv1.Gateway,
	listener gwAPIv1.Listener,
	ref gwAPIv1.ParentReference,
) bool {
	if !IsGatewayRef(gw, ref) {
		return false
	}
	if ref.SectionName == nil {
		return true
	}
	if *ref.SectionName != listener.Name {
		return false
	}
	if ref.Port == nil {
		return true
	}
	return *ref.Port == listener.Port
}

func HasHTTPRouteOwner(ownerRefs []metaV1.OwnerReference) bool {
	for _, ref := range ownerRefs {
		if ref.APIVersion == GwAPIv1APIVersion && ref.Kind == GwAPIv1HTTPRouteKind {
			return true
		}
	}
	return false
}

func IsGatewayKind(ref gwAPIv1.ParentReference) bool {
	switch {
	case ref.Group == nil:
		return false
	case *ref.Group != gwAPIv1.GroupName:
		return false
	case ref.Kind == nil:
		return false
	case *ref.Kind != gwAPIv1.Kind(GwAPIv1GatewayKind):
		return false
	default:
		return true
	}
}

func IsServiceKind(ref gwAPIv1.HTTPBackendRef) bool {
	switch {
	case ref.Group == nil:
		return false
	case *ref.Group != coreV1.GroupName:
		return false
	case ref.Kind == nil:
		return false
	case *ref.Kind != CoreV1ServiceKind:
		return false
	default:
		return true
	}
}

func IsGatewayRef(gw *gwAPIv1.Gateway, ref gwAPIv1.ParentReference) bool {
	if ref.Group != nil && *ref.Group != GwAPIv1Group {
		return false
	}
	if ref.Kind != nil && string(*ref.Kind) != gw.Kind {
		return false
	}
	if string(ref.Name) != gw.Name {
		return false
	}
	return true
}

func HasHTTPSupport(listener gwAPIv1.Listener) bool {
	if kinds, ok := ProtocolToRouteKinds[listener.Protocol]; !ok {
		return false
	} else {
		for _, k := range kinds {
			if k.Kind == GwAPIv1HTTPRouteKind {
				return true
			}
		}
	}
	return false
}

func IsGatewayComponent(obj client.Object) bool {
	return obj.GetLabels()[ComponentLabelKey] == GatewayComponentLabelValue
}

func GetMatchingGatewayLabels(obj client.Object) map[string]string {
	gwName, ok := obj.GetLabels()[InstanceLabelKey]
	if !ok {
		return map[string]string{}
	}
	return GwAPIv1GatewayLabels(gwName)
}

func IsGatewayDependent(gw *gateway.Gateway, obj client.Object) bool {
	for _, ref := range obj.GetOwnerReferences() {
		if ref.UID == gw.Object.UID {
			return true
		}
	}
	return false
}

func GetKafkaListener(gw *gateway.Gateway) *gwAPIv1.Listener {
	for i, l := range gw.Object.Spec.Listeners {
		if IsKafkaListener(l) {
			return &gw.Object.Spec.Listeners[i]
		}
	}
	return nil
}

func GetKafkaListenerStatus(gw *gateway.Gateway) *gwAPIv1.ListenerStatus {
	for i, l := range gw.Object.Spec.Listeners {
		if IsKafkaListener(l) {
			return &gw.Object.Status.Listeners[i]
		}
	}
	return nil
}

func IsKafkaListener(listener gwAPIv1.Listener) bool {
	switch {
	case listener.Protocol != gwAPIv1.TLSProtocolType:
		return false
	case listener.AllowedRoutes == nil:
		return false
	case len(listener.AllowedRoutes.Kinds) != 1:
		return false
	default:
		return IsKafkaRouteKind(listener.AllowedRoutes.Kinds[0])
	}
}

func IsKafkaRouteKind(routeKind gwAPIv1.RouteGroupKind) bool {
	switch {
	case routeKind.Group == nil:
		return false
	case *routeKind.Group != GraviteeGroup:
		return false
	case routeKind.Kind != GraviteeKafkaRouteKind:
		return false
	default:
		return true
	}
}
