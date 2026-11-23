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

package apim

import (
	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimportaltheme "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portaltheme"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Portal theme Drift detection on the APIM SDK types", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			sdk.PortalTheme{},
			sdk.PortalTheme{},
		),
		Entry("equal struct",
			sdkCompletePortalThemeDTO(),
			sdkCompletePortalThemeDTO(),
		),
		Entry("equal struct from CRD mapping",
			sdkCompletePortalThemeDTO(),
			sdkPortalThemeToDTO(sdkCompletePortalThemeCRD()),
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(sdkCompletePortalThemeDTO(), sdkCompletePortalThemeDTO())
		})
	})
})

func sdkCompletePortalThemeDTO() sdk.PortalTheme {
	GinkgoHelper()
	return loadFixture[sdk.PortalTheme]("portal_theme_dto.json")
}

func sdkCompletePortalThemeCRD() *v1alpha1.PortalTheme {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.PortalTheme]("portal_theme_crd.json")
	return &fixture
}

func sdkPortalThemeToDTO(thm *v1alpha1.PortalTheme) sdk.PortalTheme {
	GinkgoHelper()
	dto, err := apimportaltheme.ToDTO(thm)
	Expect(err).ToNot(HaveOccurred())
	return dto
}
