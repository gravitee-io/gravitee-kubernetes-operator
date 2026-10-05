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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/certificate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Certificate drift detection", func() {
	It("does not drift on the passwords and keystore file AM returns masked", func() {
		expectNoDrift(drift.DetectWithNamespace(certificateForDrift(), remoteCertificate(), ""))
	})

	It("detects drift on a configuration value AM does not mask", func() {
		remote := remoteCertificate()
		remote.Configuration = unstructured.StringifiedFrom(map[string]any{
			"content":   "********",
			"alias":     "signing",
			"storepass": "********",
			"keypass":   "********",
			"algorithm": "RS512",
			"use":       []any{"sig"},
		})

		expectDrift(drift.DetectWithNamespace(certificateForDrift(), remote, ""), `configuration:
  algorithm: "RS256" != "RS512"`)
	})

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			// system is false on a regular certificate, a zero value this check rejects.
			crd := certificateForDrift()
			crd.System = new(true)
			remote := remoteCertificate()
			remote.System = new(true)
			expectedEquivalentNotHavingAnyZeroValue(crd, remote)
		})
	})
})

func certificateForDrift() sdk.Certificate {
	GinkgoHelper()
	crd := loadFixture[v1alpha1.AMCertificate]("certificate_crd.json")
	dto, err := certificate.ToCertificateDTOForDrift(&crd)
	Expect(err).ToNot(HaveOccurred())
	return dto.Certificate
}

func remoteCertificate() sdk.Certificate {
	GinkgoHelper()
	return loadFixture[sdk.Certificate]("certificate_dto.json")
}
