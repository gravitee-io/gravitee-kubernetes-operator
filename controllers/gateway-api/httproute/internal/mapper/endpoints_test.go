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

// In-package: test/unit cannot import this internal package.
package mapper

import (
	"testing"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	. "github.com/onsi/gomega"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func twoEndpoints() []*v4.Endpoint {
	return []*v4.Endpoint{
		newEndpoint(gwAPIv1.HTTPBackendRef{}, 0, 0, "http://a.default.svc.cluster.local:8080"),
		newEndpoint(gwAPIv1.HTTPBackendRef{}, 1, 0, "http://b.default.svc.cluster.local:8080"),
	}
}

func ruleWithBackendRequest(d gwAPIv1.Duration) gwAPIv1.HTTPRouteRule {
	return gwAPIv1.HTTPRouteRule{Timeouts: &gwAPIv1.HTTPRouteTimeouts{BackendRequest: &d}}
}

func endpointHTTPConfig(ep *v4.Endpoint) *utils.GenericStringMap {
	return ep.ConfigOverride.Get("http").(*utils.GenericStringMap)
}

// The endpoints do not inherit the group configuration: the timeout must be in their own.
func TestApplyBackendTimeoutSetsEveryEndpoint(t *testing.T) {
	g := NewWithT(t)
	eps := twoEndpoints()

	applyBackendTimeout(eps, ruleWithBackendRequest("2s"))

	for _, ep := range eps {
		g.Expect(ep.Inherit).To(BeFalse())
		g.Expect(endpointHTTPConfig(ep).Get("readTimeout")).To(BeEquivalentTo(2000))
		g.Expect(endpointHTTPConfig(ep).GetBool("propagateClientHost")).To(BeTrue())
	}
}

func TestApplyBackendTimeoutZeroDisablesTheTimeout(t *testing.T) {
	g := NewWithT(t)
	eps := twoEndpoints()

	applyBackendTimeout(eps, ruleWithBackendRequest("0s"))

	for _, ep := range eps {
		g.Expect(endpointHTTPConfig(ep).Get("readTimeout")).To(BeEquivalentTo(0))
	}
}

func TestApplyBackendTimeoutWithoutTimeoutsLeavesEndpointsAlone(t *testing.T) {
	g := NewWithT(t)
	eps := twoEndpoints()

	applyBackendTimeout(eps, gwAPIv1.HTTPRouteRule{})

	for _, ep := range eps {
		g.Expect(endpointHTTPConfig(ep).Get("readTimeout")).To(BeNil())
	}
}
