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
	"os"
	"testing"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/conformance/kubernetes.io/gateway-api/impl"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/gateway-api/conformance"
	"sigs.k8s.io/gateway-api/conformance/tests"
	"sigs.k8s.io/gateway-api/conformance/utils/config"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"
	"sigs.k8s.io/gateway-api/pkg/features"
	"sigs.k8s.io/yaml"
)

var lazyTimeoutConfig = config.TimeoutConfig{
	TestIsolation:                      0 * time.Second,
	GWCMustBeAccepted:                  300 * time.Second,
	GatewayStatusMustHaveListeners:     300 * time.Second,
	GatewayListenersMustHaveConditions: 300 * time.Second,
	HTTPRouteMustNotHaveParents:        180 * time.Second,
	HTTPRouteMustHaveCondition:         180 * time.Second,
	TLSRouteMustHaveCondition:          180 * time.Second,
	RouteMustHaveParents:               180 * time.Second,
	GetTimeout:                         180 * time.Second,
}

func TestGatewayAPIConformance(t *testing.T) {
	flag.Parse()

	opts := conformance.DefaultOptions(t)

	opts.Implementation = impl.Manifest
	opts.ReportOutputPath = impl.GetReportOutputPath()

	opts.ConformanceProfiles = sets.New(
		suite.GatewayHTTPConformanceProfileName,
	)

	opts.SupportedFeatures = sets.New(
		features.GatewayFeature.Name,
		features.HTTPRouteFeature.Name,
		features.ReferenceGrantFeature.Name,
	)

	opts.Mode = "default"

	opts.TimeoutConfig = lazyTimeoutConfig
	opts.RestConfig.QPS = -1

	// Here you can specify test name for debug purpose

	// Failing tests

	//   HTTPRouteMatching
	//   HTTPRouteHTTPSListener
	//   HTTPRouteListenerHostnameMatching
	//   HTTPRouteMatchingAcrossRoutes

	opts.RunTest = ""
	opts.CleanupBaseResources = false

	opts.SkipTests = []string{}

	// We skip this test because right now we cannot accept different
	// routes with the same host and path.
	// For that reason the second route will get Accepted but not Programmed
	// because of the conflict.
	opts.SkipTests = append(opts.SkipTests, "HTTPRouteMatchingAcrossRoutes")

	// We skip this test in circle ci because for some reason
	// threads get blocked on the gatway side when
	// running it. Needs to investigate (possibly how java Atomic are handled by the underlying system)
	if os.Getenv("CIRCLECI") == env.TrueString {
		opts.SkipTests = append(opts.SkipTests, "HTTPRouteWeight")
	}

	// That one might be handled first because it looks like it sits in our code base.
	opts.SkipTests = append(opts.SkipTests, "HTTPRouteServiceTypes")

	opts.CleanupBaseResources = false

	cSuite, err := suite.NewConformanceTestSuite(opts)
	if err != nil {
		t.Fatalf("Error creating conformance test suite: %v", err)
	}
	cSuite.Setup(t, tests.ConformanceTests)
	if err := cSuite.Run(t, tests.ConformanceTests); err != nil {
		t.Fatalf("Error running conformance tests: %v", err)
	}

	generateReport(t, cSuite, opts)
}

func generateReport(t *testing.T, cSuite *suite.ConformanceTestSuite, opts suite.ConformanceOptions) {
	report, err := cSuite.Report()
	if err != nil {
		t.Fatalf("error generating conformance profile report: %v", err)
	}

	rawReport, err := yaml.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}

	if err = os.WriteFile(opts.ReportOutputPath, rawReport, 0o600); err != nil {
		t.Fatal(err)
	}
}
