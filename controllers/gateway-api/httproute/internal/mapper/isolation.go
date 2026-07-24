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
	"regexp"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/el"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	hostNotMatchesCondition = el.Expression("(#request.headers['Host'][0] matches '%s') eq false")
)

// computeIsolationExclusions returns EL expressions that exclude hosts belonging
// to more-specific listeners on the same gateway port. This implements Gateway API
// HTTPListenerIsolation semantics.
func computeIsolationExclusions(ctx context.Context, route *gwAPIv1.HTTPRoute) ([]el.Expression, error) {
	if len(route.Spec.ParentRefs) == 0 {
		return nil, nil
	}

	ref := route.Spec.ParentRefs[0]
	if !k8s.IsGatewayKind(ref) || ref.SectionName == nil {
		return nil, nil
	}

	gw, err := k8s.ResolveGateway(ctx, route.ObjectMeta, ref)
	if err != nil {
		return nil, err
	}

	var attachedListener *gwAPIv1.Listener
	for i := range gw.Spec.Listeners {
		if gw.Spec.Listeners[i].Name == *ref.SectionName {
			attachedListener = &gw.Spec.Listeners[i]
			break
		}
	}

	if attachedListener == nil {
		return nil, nil
	}

	otherHostnames := collectMoreSpecificHostnames(gw, attachedListener)
	if len(otherHostnames) == 0 {
		return nil, nil
	}

	exclusions := make([]el.Expression, 0, len(otherHostnames))
	for _, h := range otherHostnames {
		pattern := hostnameToRegex(h)
		exclusions = append(exclusions, hostNotMatchesCondition.Format(pattern))
	}
	return exclusions, nil
}

// collectMoreSpecificHostnames returns hostnames from other listeners on the same port
// that are more specific than the attached listener's hostname and overlap with it.
func collectMoreSpecificHostnames(gw *gwAPIv1.Gateway, attached *gwAPIv1.Listener) []string {
	var hostnames []string

	for i := range gw.Spec.Listeners {
		other := &gw.Spec.Listeners[i]
		if other.Name == attached.Name || other.Port != attached.Port {
			continue
		}
		if other.Hostname == nil {
			continue
		}

		otherHost := string(*other.Hostname)

		if attached.Hostname == nil {
			hostnames = append(hostnames, otherHost)
			continue
		}

		attachedHost := string(*attached.Hostname)
		if isMoreSpecificAndOverlaps(otherHost, attachedHost) {
			hostnames = append(hostnames, otherHost)
		}
	}

	return hostnames
}

// isMoreSpecificAndOverlaps returns true if candidate is more specific than reference
// and a request matching candidate would also match reference in Gravitee's virtual hosting
// (which uses multi-level wildcards where *.example.com matches bar.foo.example.com).
func isMoreSpecificAndOverlaps(candidate, reference string) bool {
	candidateIsWildcard := strings.HasPrefix(candidate, "*.")
	referenceIsWildcard := strings.HasPrefix(reference, "*.")

	if !referenceIsWildcard {
		return false
	}

	refSuffix := reference[1:]

	if !candidateIsWildcard {
		return strings.HasSuffix(candidate, refSuffix)
	}

	candidateSuffix := candidate[1:]
	return strings.HasSuffix(candidateSuffix, refSuffix)
}

// hostnameToRegex converts a Gateway API hostname to a Java regex pattern for Gravitee EL.
func hostnameToRegex(hostname string) string {
	if strings.HasPrefix(hostname, "*.") {
		suffix := regexp.QuoteMeta(hostname[2:])
		return `^[^.]+\.` + suffix + `$`
	}
	return `^` + regexp.QuoteMeta(hostname) + `$`
}
