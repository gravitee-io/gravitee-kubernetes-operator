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

const timeout = 180 * time.Second

func TestGatewayAPIConformance(t *testing.T) {
	flag.Parse()

	opts := conformance.DefaultOptions(t)
	opts.SupportedFeatures = sets.New(
		features.GatewayFeature.Name,
		features.HTTPRouteFeature.Name,
		// features.GRPCRouteFeature.Name,
		// features.ReferenceGrantFeature.Name,
	)

	opts.TimeoutConfig = config.DefaultTimeoutConfig()
	opts.TimeoutConfig.GatewayStatusMustHaveListeners = timeout
	opts.TimeoutConfig.GatewayListenersMustHaveConditions = timeout
	opts.TimeoutConfig.HTTPRouteMustHaveCondition = timeout
	opts.RestConfig.QPS = -1

	// Here you can specify test name for debug purpose

	// Failing tests

	//   HTTPRouteMatching
	//   HTTPRouteHTTPSListener
	//   HTTPRouteListenerHostnameMatching
	//   HTTPRouteMatchingAcrossRoutes

	opts.RunTest = ""
	opts.CleanupBaseResources = false

	cSuite, err := suite.NewConformanceTestSuite(opts)
	if err != nil {
		t.Fatalf("Error creating conformance test suite: %v", err)
	}
	cSuite.Setup(t, tests.ConformanceTests)
	if err := cSuite.Run(t, tests.ConformanceTests); err != nil {
		t.Fatalf("Error running conformance tests: %v", err)
	}
}
