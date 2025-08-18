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
	"context"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwAPIv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

const (
	GwAPIv1HTTPRouteKind = "HTTPRoute"
	GwAPIv1GatewayKind   = "Gateway"
	CoreV1ServiceKind    = "Service"
	CoreV1SecretKind     = "Secret"
	GwAPIv1APIVersion    = "gateway.networking.k8s.io/v1"

	GraviteeAPIVersion     = "gravitee.io/v1alpha1"
	GraviteeKafkaRouteKind = "KafkaRoute"

	ServiceURIPattern = "http://%s.%s.svc.cluster.local:%d"
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

func IsAttachedHTTPRoute(
	gw *gwAPIv1.Gateway,
	listener gwAPIv1.Listener,
	route gwAPIv1.HTTPRoute,
) bool {
	for i := range route.Status.Parents {
		ref := route.Spec.ParentRefs[i]
		if IsListenerRef(gw, listener, ref) {
			status := gateway.WrapRouteParentStatus(&route.Status.Parents[i])
			return IsAccepted(status)
		}
	}
	return false
}

func IsAttachedKafkaRoute(
	gw *gwAPIv1.Gateway,
	listener gwAPIv1.Listener,
	route v1alpha1.KafkaRoute,
) bool {
	for i := range route.Status.Parents {
		ref := route.Spec.ParentRefs[i]
		if IsListenerRef(gw, listener, ref) {
			status := gateway.WrapRouteParentStatus(&route.Status.Parents[i])
			return IsAccepted(status)
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

func IsServiceKind(ref gwAPIv1.BackendObjectReference) bool {
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

func IsSecretRef(secret *coreV1.Secret, ref gwAPIv1.SecretObjectReference) bool {
	if ref.Group != nil && *ref.Group != "" {
		return false
	}
	if ref.Kind != nil && string(*ref.Kind) != CoreV1SecretKind {
		return false
	}
	if string(ref.Name) != secret.Name {
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

func GetKafkaListener(gw *gwAPIv1.Gateway) *gwAPIv1.Listener {
	for i, l := range gw.Spec.Listeners {
		if IsKafkaListener(l) {
			return &gw.Spec.Listeners[i]
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

func HasHafkaListener(gw *gateway.Gateway) bool {
	for _, l := range gw.Object.Spec.Listeners {
		if IsKafkaListener(l) {
			return true
		}
	}
	return false
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

func IsKafkaListenerStatus(listener gwAPIv1.ListenerStatus) bool {
	switch {
	case listener.SupportedKinds == nil:
		return false
	case len(listener.SupportedKinds) != 1:
		return false
	default:
		return IsKafkaRouteKind(listener.SupportedKinds[0])
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

func HasGatewayClassParameters(gw *gwAPIv1.GatewayClass, params *v1alpha1.GatewayClassParameters) bool {
	paramsRef := gw.Spec.ParametersRef
	switch {
	case paramsRef.Group != "gravitee.io":
		return false
	case paramsRef.Kind != "GatewayClassParameters":
		return false
	case paramsRef.Name != params.Name:
		return false
	case paramsRef.Namespace == nil:
		return gw.Namespace == params.Namespace
	default:
		return string(*paramsRef.Namespace) == params.Namespace
	}
}

func HasKafkaEnabled(params *v1alpha1.GatewayClassParameters) bool {
	return params.Spec.Gravitee != nil &&
		params.Spec.Gravitee.Kafka != nil &&
		params.Spec.Gravitee.Kafka.Enabled
}

func HasGraviteeYAML(params *v1alpha1.GatewayClassParameters) bool {
	return params.Spec.Gravitee != nil && params.Spec.Gravitee.YAML != nil
}

func HasGraviteeLicense(params *v1alpha1.GatewayClassParameters) bool {
	return params.Spec.Gravitee != nil && params.Spec.Gravitee.LicenseRef != nil
}

func FindListenerIndexBySectionName(gw *gwAPIv1.Gateway, sectionName gwAPIv1.SectionName) int {
	for i := range gw.Status.Listeners {
		lst := gw.Status.Listeners[i]
		if lst.Name == sectionName {
			return i
		}
	}
	return -1
}

func HasHTTPListenerAtIndex(gw *gwAPIv1.Gateway, index int) bool {
	if index < 0 || index >= len(gw.Status.Listeners) {
		return false
	}
	lst := gw.Status.Listeners[index]
	for j := range lst.SupportedKinds {
		if lst.SupportedKinds[j].Kind == gwAPIv1.Kind(GwAPIv1HTTPRouteKind) {
			return true
		}
	}
	return false
}

func GetHTTPHosts(route *gwAPIv1.HTTPRoute, gw *gwAPIv1.Gateway, ref gwAPIv1.ParentReference) []string {
	if ref.SectionName != nil {
		listenerIndex := FindListenerIndexBySectionName(gw, *ref.SectionName)
		if !HasHTTPListenerAtIndex(gw, listenerIndex) {
			return nil
		}
		listener := gw.Spec.Listeners[listenerIndex]
		return GetHTTPHostsForListenerAndRoute(route, listener)
	}
	hostnames := []string{}
	for i := range gw.Spec.Listeners {
		if !HasHTTPListenerAtIndex(gw, i) {
			continue
		}
		listener := gw.Spec.Listeners[i]
		hostnames = append(hostnames, GetHTTPHostsForListenerAndRoute(route, listener)...)
	}
	return hostnames
}

func GetHTTPHostsForListenerAndRoute(httpRoute *gwAPIv1.HTTPRoute, listener gwAPIv1.Listener) []string {
	hostnames := sets.New[string]()
	if listener.Hostname == nil {
		for i := range httpRoute.Spec.Hostnames {
			routeHostname := string(httpRoute.Spec.Hostnames[i])
			hostnames.Insert(routeHostname)
		}
		return hostnames.UnsortedList()
	}

	listenerHostname := string(*listener.Hostname)
	for i := range httpRoute.Spec.Hostnames {
		routeHostname := string(httpRoute.Spec.Hostnames[i])
		if intersect(listenerHostname, routeHostname) {
			if instersecFromRouteWildcard(listenerHostname, routeHostname) {
				hostnames.Insert(listenerHostname)
			} else {
				hostnames.Insert(routeHostname)
			}
		}
	}

	return hostnames.UnsortedList()
}

func HasIntersectingHostName(route *gwAPIv1.HTTPRoute, gw *gwAPIv1.Gateway, ref gwAPIv1.ParentReference) bool {
	if ref.SectionName != nil {
		i := FindListenerIndexBySectionName(gw, *ref.SectionName)
		if i < 0 {
			return false
		}
		return HasIntersectingHostNameAtIndex(route, gw, i)
	}
	for i := range gw.Spec.Listeners {
		if HasIntersectingHostNameAtIndex(route, gw, i) {
			return true
		}
	}
	return false
}

func HasIntersectingHostNameAtIndex(route *gwAPIv1.HTTPRoute, gw *gwAPIv1.Gateway, index int) bool {
	listener := gw.Spec.Listeners[index]
	listenerHostname := listener.Hostname

	if listenerHostname == nil {
		return true
	}

	if len(route.Spec.Hostnames) == 0 {
		return true
	}

	for _, routeHostname := range route.Spec.Hostnames {
		if intersect(string(*listenerHostname), string(routeHostname)) {
			return true
		}
	}
	return false
}

func intersect(listenerHostname, routeHostname string) bool {
	if strings.EqualFold(listenerHostname, routeHostname) {
		return true
	}

	if instersecFromListenerWildcard(listenerHostname, routeHostname) {
		return true
	}

	if instersecFromRouteWildcard(listenerHostname, routeHostname) {
		return true
	}

	return false
}

func instersecFromListenerWildcard(listenerHostname, routeHostname string) bool {
	if strings.HasPrefix(listenerHostname, "*") {
		return intersectWithWildcard(routeHostname, listenerHostname)
	}
	return false
}

func instersecFromRouteWildcard(listenerHostname, routeHostname string) bool {
	if strings.HasPrefix(routeHostname, "*") {
		return intersectWithWildcard(listenerHostname, routeHostname)
	}
	return false
}

func intersectWithWildcard(hostname, wildcardHostname string) bool {
	if !strings.HasSuffix(hostname, strings.TrimPrefix(wildcardHostname, "*")) {
		return false
	}

	wildcardMatch := strings.TrimSuffix(hostname, strings.TrimPrefix(wildcardHostname, "*"))
	return len(wildcardMatch) > 0
}

func IsGrantedReference(ctx context.Context, from client.Object, to gwAPIv1.ObjectReference) (bool, error) {
	grantNs := GetRefNs(from, to.Namespace)

	if from.GetNamespace() == grantNs {
		return true, nil
	}

	grantList := &gwAPIv1beta1.ReferenceGrantList{}
	opts := &client.ListOptions{Namespace: grantNs}
	if err := GetClient().List(ctx, grantList, opts); err != nil {
		return false, err
	}
	for _, grant := range grantList.Items {
		if hasGrantedFrom(from, grant) && hasGrantedTo(to, grant) {
			return true, nil
		}
	}
	return false, nil
}

func hasGrantedTo(to gwAPIv1.ObjectReference, grant gwAPIv1beta1.ReferenceGrant) bool {
	for _, grantTo := range grant.Spec.To {
		if isGrantedTo(to, grantTo) {
			return true
		}
	}
	return false
}

func isGrantedTo(to gwAPIv1.ObjectReference, grantTo gwAPIv1beta1.ReferenceGrantTo) bool {
	switch {
	case grantTo.Name != nil && to.Name != *grantTo.Name:
		return false
	case to.Kind != grantTo.Kind:
		return false
	case to.Group != grantTo.Group:
		return false
	default:
		return true
	}
}

func hasGrantedFrom(from client.Object, grant gwAPIv1beta1.ReferenceGrant) bool {
	for _, grantFrom := range grant.Spec.From {
		if isGrantedFrom(from, grantFrom) {
			return true
		}
	}
	return false
}

func isGrantedFrom(from client.Object, grantFrom gwAPIv1beta1.ReferenceGrantFrom) bool {
	switch {
	case from.GetNamespace() != string(grantFrom.Namespace):
		return false
	case from.GetObjectKind().GroupVersionKind().Group != string(grantFrom.Group):
		return false
	case from.GetObjectKind().GroupVersionKind().Kind != string(grantFrom.Kind):
		return false
	default:
		return true
	}
}

func GetRefNs(referencer client.Object, refNs *gwAPIv1.Namespace) string {
	if refNs != nil {
		return string(*refNs)
	}
	return referencer.GetNamespace()
}
