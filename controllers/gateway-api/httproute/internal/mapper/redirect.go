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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var (
	httpRedirectPolicyName = "http-redirect"

	rulesKey        = "rules"
	rulePathKey     = "path"
	ruleLocationKey = "location"
	ruleStatusKey   = "status"

	rulePath = "(.*)"

	locationPathWhenNoGivenPath   = "{#request.contextPath.replaceAll('/$', '')}{#request.pathInfo}"
	locationPathWhenReplacePrefix = "/%s{#request.pathInfo}"

	locationSchemeDefault = "{#request.scheme}"
	locationHostDefault   = "{#request.host}"

	statusCodeDefault = 302
)

func buildHTTPRedirect(
	ctx context.Context,
	route *gwAPIv1.HTTPRoute,
	filter gwAPIv1.HTTPRequestRedirectFilter,
) *v4.FlowStep {
	return v4.NewFlowStep(base.FlowStep{
		Policy:        &httpRedirectPolicyName,
		Enabled:       true,
		Configuration: buildHTTPRedirectConfig(ctx, route, filter),
	})
}

func buildHTTPRedirectConfig(
	ctx context.Context,
	route *gwAPIv1.HTTPRoute,
	filter gwAPIv1.HTTPRequestRedirectFilter,
) *utils.GenericStringMap {
	config := utils.NewGenericStringMap()
	rules := []any{buildRedirectRule(ctx, route, filter)}
	config.Put(rulesKey, rules)
	return config
}

func buildRedirectRule(
	ctx context.Context,
	route *gwAPIv1.HTTPRoute,
	filter gwAPIv1.HTTPRequestRedirectFilter,
) map[string]any {
	rule := make(map[string]any)
	rule[rulePathKey] = rulePath
	rule[ruleLocationKey] = buildRedirectLocation(ctx, route, filter)
	rule[ruleStatusKey] = getStatusCode(filter)
	return rule
}

func buildRedirectLocation(
	ctx context.Context,
	route *gwAPIv1.HTTPRoute,
	filter gwAPIv1.HTTPRequestRedirectFilter,
) string {
	scheme := getLocationScheme(filter)
	host := getLocationHost(ctx, route, filter)
	path := getLocationPath(filter)
	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}

func getLocationScheme(filter gwAPIv1.HTTPRequestRedirectFilter) string {
	if filter.Scheme != nil {
		return *filter.Scheme
	}
	return locationSchemeDefault
}

func getLocationHost(ctx context.Context, route *gwAPIv1.HTTPRoute, filter gwAPIv1.HTTPRequestRedirectFilter) string {
	portSuffix := getHostPortSuffix(ctx, route, filter)

	if filter.Hostname != nil {
		host := string(*filter.Hostname)
		return host + portSuffix
	}
	return locationHostDefault + portSuffix
}

func getHostPortSuffix(
	ctx context.Context,
	route *gwAPIv1.HTTPRoute,
	filter gwAPIv1.HTTPRequestRedirectFilter,
) string {
	scheme := getLocationScheme(filter)

	if filter.Port != nil {
		return getActualPortSuffix(ctx, route, scheme, filter.Port)
	}

	port := inferPort(ctx, route, filter)

	return getActualPortSuffix(ctx, route, scheme, port)
}

func inferPort(
	ctx context.Context,
	route *gwAPIv1.HTTPRoute,
	filter gwAPIv1.HTTPRequestRedirectFilter,
) *gwAPIv1.PortNumber {
	scheme := getLocationScheme(filter)
	if scheme == "http" {
		port := gwAPIv1.PortNumber(80)
		return &port
	}
	if scheme == "https" {
		port := gwAPIv1.PortNumber(443)
		return &port
	}
	return getListenerPort(ctx, route)
}

func getActualPortSuffix(
	ctx context.Context,
	route *gwAPIv1.HTTPRoute,
	scheme string,
	port *gwAPIv1.PortNumber,
) string {
	if port == nil || shouldOmitPort(ctx, route, scheme, port) {
		return ""
	}
	return fmt.Sprintf(":%d", *port)
}

func shouldOmitPort(ctx context.Context, route *gwAPIv1.HTTPRoute, scheme string, port *gwAPIv1.PortNumber) bool {
	if *port == 80 {
		if scheme == "http" {
			return true
		}
		if scheme == locationSchemeDefault && isHTTPListener(ctx, route) {
			return true
		}
	}

	if *port == 443 {
		if scheme == "https" {
			return true
		}
		if scheme == locationSchemeDefault && isHTTPSListener(ctx, route) {
			return true
		}
	}

	return false
}

func getListenerPort(ctx context.Context, route *gwAPIv1.HTTPRoute) *gwAPIv1.PortNumber {
	for i := range route.Status.Parents {
		parentStatus := route.Status.Parents[i]
		if !k8s.IsAccepted(gateway.WrapRouteParentStatus(&parentStatus)) {
			continue
		}

		parentRef := parentStatus.ParentRef
		gw, err := resolveGateway(ctx, route.ObjectMeta, parentRef, k8s.ResolveGateway)
		if err != nil {
			continue
		}

		port := getPortFromGateway(gw, parentRef)
		if port != nil {
			return port
		}
	}
	return nil
}

func getPortFromGateway(gw *gwAPIv1.Gateway, parentRef gwAPIv1.ParentReference) *gwAPIv1.PortNumber {
	if parentRef.SectionName != nil {
		listenerIndex := findListenerIndexBySectionName(gw, parentRef.SectionName)
		if listenerIndex >= 0 && listenerIndex < len(gw.Spec.Listeners) {
			port := gw.Spec.Listeners[listenerIndex].Port
			return &port
		}
		return nil
	}
	return findFirstHTTPListenerPort(gw, parentRef)
}

func findFirstHTTPListenerPort(gw *gwAPIv1.Gateway, parentRef gwAPIv1.ParentReference) *gwAPIv1.PortNumber {
	if parentRef.Port != nil {
		for _, listener := range gw.Spec.Listeners {
			if (listener.Protocol == gwAPIv1.HTTPProtocolType || listener.Protocol == gwAPIv1.HTTPSProtocolType) &&
				listener.Port == *parentRef.Port {
				port := listener.Port
				return &port
			}
		}
	}
	for _, listener := range gw.Spec.Listeners {
		if listener.Protocol == gwAPIv1.HTTPProtocolType || listener.Protocol == gwAPIv1.HTTPSProtocolType {
			port := listener.Port
			return &port
		}
	}
	return nil
}

func findListenerIndexBySectionName(gw *gwAPIv1.Gateway, sectionName *gwAPIv1.SectionName) int {
	listenerIndex := k8s.FindListenerIndexBySectionName(gw, *sectionName)
	if listenerIndex < 0 {
		listenerIndex = k8s.FindListenerIndexBySectionNameInSpec(gw, *sectionName)
	}
	return listenerIndex
}

func isHTTPListener(ctx context.Context, route *gwAPIv1.HTTPRoute) bool {
	return hasProtocolListener(ctx, route, gwAPIv1.HTTPProtocolType)
}

func isHTTPSListener(ctx context.Context, route *gwAPIv1.HTTPRoute) bool {
	return hasProtocolListener(ctx, route, gwAPIv1.HTTPSProtocolType)
}

func hasProtocolListener(ctx context.Context, route *gwAPIv1.HTTPRoute, protocol gwAPIv1.ProtocolType) bool {
	for i := range route.Status.Parents {
		if i >= len(route.Spec.ParentRefs) {
			continue
		}
		parentStatus := route.Status.Parents[i]
		if !k8s.IsAccepted(gateway.WrapRouteParentStatus(&parentStatus)) {
			continue
		}

		parentRef := parentStatus.ParentRef
		gw, err := resolveGateway(ctx, route.ObjectMeta, parentRef, k8s.ResolveGateway)
		if err != nil {
			continue
		}

		if hasParentWithProtocol(gw, parentRef.SectionName, protocol) {
			return true
		}
	}
	return false
}

func hasParentWithProtocol(gw *gwAPIv1.Gateway, sectionName *gwAPIv1.SectionName, protocol gwAPIv1.ProtocolType) bool {
	if sectionName != nil {
		listenerIndex := k8s.FindListenerIndexBySectionName(gw, *sectionName)
		if listenerIndex < 0 {
			listenerIndex = k8s.FindListenerIndexBySectionNameInSpec(gw, *sectionName)
		}
		if listenerIndex >= 0 && listenerIndex < len(gw.Spec.Listeners) {
			return gw.Spec.Listeners[listenerIndex].Protocol == protocol
		}
	} else {
		for _, listener := range gw.Spec.Listeners {
			if listener.Protocol == protocol {
				return true
			}
		}
	}
	return false
}

func getLocationPath(filter gwAPIv1.HTTPRequestRedirectFilter) string {
	if filter.Path == nil {
		return locationPathWhenNoGivenPath
	}
	if filter.Path.ReplaceFullPath != nil {
		return *filter.Path.ReplaceFullPath
	}
	return fmt.Sprintf(locationPathWhenReplacePrefix, *filter.Path.ReplacePrefixMatch)
}

func getStatusCode(filter gwAPIv1.HTTPRequestRedirectFilter) int {
	if filter.StatusCode != nil {
		return *filter.StatusCode
	}
	return statusCodeDefault
}
