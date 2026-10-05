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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/securitydomain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Security domain drift detection", func() {
	It("does not drift between the CRD mapping and the AM response", func() {
		expectNoDrift(drift.DetectWithNamespace(securityDomainForDrift(), remoteSecurityDomain(), ""))
	})

	It("does not drift on the data plane AM resolves when the CRD leaves it unset", func() {
		crd := securityDomainForDrift()
		crd.DataPlaneId = nil

		expectNoDrift(drift.DetectWithNamespace(crd, remoteSecurityDomain(), ""))
	})

	It("detects drift on a remote change", func() {
		remote := remoteSecurityDomain()
		remote.CorsSettings.Enabled = new(false)

		expectDrift(drift.DetectWithNamespace(securityDomainForDrift(), remote, ""), `corsSettings:
  enabled: true != false`)
	})

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(securityDomainForDrift(), remoteSecurityDomain())
		})
	})
})

func securityDomainForDrift() sdk.Domain {
	GinkgoHelper()
	crd := loadFixture[v1alpha1.AMSecurityDomain]("securitydomain_crd.json")
	dto, err := securitydomain.ToDomainDTO(&crd)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func remoteSecurityDomain() sdk.Domain {
	GinkgoHelper()
	return loadFixture[sdk.Domain]("securitydomain_dto.json")
}
