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

package standard

import (
	"flag"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/gateway-api/conformance"
	"sigs.k8s.io/gateway-api/conformance/tests"
	"sigs.k8s.io/gateway-api/conformance/utils/config"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"
	"sigs.k8s.io/gateway-api/pkg/features"
)

var lazyTimeoutConfig = config.TimeoutConfig{
	TestIsolation:                      3 * time.Second,
	GWCMustBeAccepted:                  300 * time.Second,
	GatewayStatusMustHaveListeners:     180 * time.Second,
	GatewayListenersMustHaveConditions: 180 * time.Second,
	HTTPRouteMustNotHaveParents:        180 * time.Second,
	HTTPRouteMustHaveCondition:         180 * time.Second,
	TLSRouteMustHaveCondition:          180 * time.Second,
	RouteMustHaveParents:               180 * time.Second,
}

func TestGatewayAPIConformance(t *testing.T) {
	flag.Parse()

	opts := conformance.DefaultOptions(t)

	opts.SupportedFeatures = sets.New(
		features.GatewayFeature.Name,
		features.HTTPRouteFeature.Name,
		features.ReferenceGrantFeature.Name,
		// features.GRPCRouteFeature.Name,
	)

	opts.TimeoutConfig = lazyTimeoutConfig
	opts.RestConfig.QPS = -1

	// Here you can specify test name for debug purpose
	opts.RunTest = ""

	opts.SkipTests = []string{}

	// We skip this test because right now we cannot accept different
	// routes with the same host and path.
	// For that reason the second route will get Accepted but not Programmed
	// because of the conflict.
	opts.SkipTests = append(opts.SkipTests, "HTTPRouteMatchingAcrossRoutes")

	// We skip this test because for some reason
	// threads get blocked on the gatway side when
	// running it. Needs to investigate
	opts.SkipTests = append(opts.SkipTests, "HTTPRouteWeight")

	// That one might be handled first because it looks like its half
	// baked implementation from our side.
	opts.SkipTests = append(opts.SkipTests, "HTTPRouteServiceTypes")

	cSuite, err := suite.NewConformanceTestSuite(opts)
	if err != nil {
		t.Fatalf("Error creating conformance test suite: %v", err)
	}

	cSuite.Setup(t, tests.ConformanceTests)
	if err := cSuite.Run(t, tests.ConformanceTests); err != nil {
		t.Fatalf("Error running conformance tests: %v", err)
	}
}
