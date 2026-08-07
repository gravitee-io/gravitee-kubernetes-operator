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
	"sigs.k8s.io/gateway-api/conformance"
	"sigs.k8s.io/gateway-api/conformance/tests"
	"sigs.k8s.io/gateway-api/conformance/utils/config"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"
	"sigs.k8s.io/gateway-api/pkg/features"
	"sigs.k8s.io/yaml"
)

// Set to "true" to skip the tests that need more CPU than a small runner has.
// See the comment next to its only use below.
const skipStarvedTests = "CONFORMANCE_SKIP_STARVED_TESTS"

var lazyTimeoutConfig = config.TimeoutConfig{
	TestIsolation:                      100 * time.Millisecond,
	DefaultTestTimeout:                 180 * time.Second,
	MaxTimeToConsistency:               300 * time.Second,
	GWCMustBeAccepted:                  300 * time.Second,
	GatewayStatusMustHaveListeners:     300 * time.Second,
	GatewayListenersMustHaveConditions: 300 * time.Second,
	HTTPRouteMustNotHaveParents:        180 * time.Second,
	HTTPRouteMustHaveCondition:         180 * time.Second,
	TLSRouteMustHaveCondition:          180 * time.Second,
	RouteMustHaveParents:               180 * time.Second,
	GetTimeout:                         180 * time.Second,
	LatestObservedGenerationSet:        180 * time.Second,
	NamespacesMustBeReady:              600 * time.Second,
}

func TestGatewayAPIConformance(t *testing.T) {
	flag.Parse()

	opts := conformance.DefaultOptions(t)

	opts.Implementation = impl.Manifest
	opts.ReportOutputPath = impl.GetReportOutputPath()

	opts.ConformanceProfiles = []suite.ConformanceProfileName{
		suite.GatewayHTTPConformanceProfileName,
	}

	opts.SupportedFeatures = []features.FeatureName{
		features.GatewayFeature.Name,
		features.HTTPRouteFeature.Name,
		features.ReferenceGrantFeature.Name,
		features.SupportGatewayPort8080,
		features.SupportHTTPRouteSchemeRedirect,
		features.SupportHTTPRoutePortRedirect,
		features.SupportHTTPRoutePathRedirect,
		features.SupportHTTPRouteResponseHeaderModification,
		features.SupportHTTPRoutePathRewrite,
		features.SupportHTTPRouteMethodMatching,
		features.SupportHTTPRouteQueryParamMatching,
		features.SupportHTTPRouteHostRewrite,
		features.SupportHTTPRouteBackendRequestHeaderModification,
		features.SupportHTTPRouteNamedRouteRule,
		features.SupportGatewayHTTPListenerIsolation,
		features.SupportGatewayInfrastructurePropagation,
		features.SupportHTTPRouteBackendProtocolH2C,
		features.SupportGatewayFrontendClientCertificateValidation,
	}

	opts.Mode = "default"

	opts.TimeoutConfig = lazyTimeoutConfig
	opts.RestConfig.QPS = -1

	// Here you can specify test name for debug purpose
	opts.RunTest = ""
	opts.CleanupBaseResources = false

	opts.SkipTests = []string{}

	if os.Getenv(env.GatewayAPIMatchAcrossRoutes) != env.TrueString {
		opts.SkipTests = append(opts.SkipTests, "HTTPRouteMatchingAcrossRoutes")
	}

	// These two starve on an under-provisioned runner: HTTPRouteWeight blocks
	// threads on the gateway side, HTTPRouteRedirectPortAndScheme just fails.
	// They were gated on CIRCLECI, which made every CI report partial and so
	// unsubmittable. The gate is now explicit, so a job that runs on a big
	// enough box simply does not set it and gets the full suite.
	//
	// A run whose report is meant for submission must never set this.
	if os.Getenv(skipStarvedTests) == env.TrueString {
		opts.SkipTests = append(opts.SkipTests, "HTTPRouteWeight", "HTTPRouteRedirectPortAndScheme")
	}

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
