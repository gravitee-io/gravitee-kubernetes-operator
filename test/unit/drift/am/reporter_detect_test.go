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

package am

import (
	"github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2/pkg/sdk"
	"github.com/gravitee-io/gravitee-automation-sdk/common/pkg/unstructured"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/reporter"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reporter drift detection", func() {
	It("does not drift on the password AM returns masked", func() {
		expectNoDrift(drift.DetectWithNamespace(reporterForDrift(), remoteReporter(), ""))
	})

	It("detects drift on a configuration value AM does not mask", func() {
		remote := remoteReporter()
		remote.Configuration = unstructured.StringifiedFrom(map[string]any{
			"bootstrapServers": "kafka:9092",
			"topic":            "renamed-in-am",
			"acks":             "1",
			"username":         "gravitee",
			"password":         "********",
		})

		expectDrift(drift.DetectWithNamespace(reporterForDrift(), remote, ""), `configuration:
  topic: "audit" != "renamed-in-am"`)
	})

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			// system is false on a regular reporter, a zero value this check rejects.
			crd := reporterForDrift()
			crd.System = new(true)
			remote := remoteReporter()
			remote.System = new(true)
			expectedEquivalentNotHavingAnyZeroValue(crd, remote)
		})
	})
})

func reporterForDrift() sdk.Reporter {
	GinkgoHelper()
	crd := loadFixture[v1alpha1.AMReporter]("reporter_crd.json")
	dto, err := reporter.ToReporterDTOForDrift(&crd)
	Expect(err).ToNot(HaveOccurred())
	return dto.Reporter
}

func remoteReporter() sdk.Reporter {
	GinkgoHelper()
	return loadFixture[sdk.Reporter]("reporter_dto.json")
}
