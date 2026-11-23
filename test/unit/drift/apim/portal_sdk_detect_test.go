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
	apimportal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Portal Drift detection on the APIM SDK types", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			sdk.Portal{},
			sdk.Portal{},
		),
		Entry("equal struct",
			sdkCompletePortalDTO(),
			sdkCompletePortalDTO(),
		),
		Entry("equal struct from CRD mapping",
			sdkCompletePortalDTO(),
			sdkPortalToDTO(sdkCompletePortalCRD()),
		),
		// Tests for the new structure field
		Entry("equal struct with structure",
			sdkCompletePortalDTOWithStructure(),
			sdkCompletePortalDTOWithStructure(),
		),
		Entry("equal struct with structure from CRD mapping",
			sdkCompletePortalDTOWithStructure(),
			sdkPortalToDTO(sdkCompletePortalCRDWithStructure()),
		),
		Entry("unset visibility resolved to PUBLIC matches the remote",
			apimportal.DriftDTO(sdkPortalWithTopNavbar(sdkTopNavbarEntry("/alpha", ""))),
			apimportal.DriftDTO(sdkPortalWithTopNavbar(sdkTopNavbarEntry("/alpha", sdk.PortalVisibilityPUBLIC))),
		),
		Entry("unset visibility under a PRIVATE parent matches the inherited remote",
			apimportal.DriftDTO(sdkPortalWithTopNavbar(
				sdkTopNavbarEntry("/alpha", sdk.PortalVisibilityPRIVATE),
				sdkTopNavbarEntry("/alpha/docs", ""),
			)),
			apimportal.DriftDTO(sdkPortalWithTopNavbar(
				sdkTopNavbarEntry("/alpha", sdk.PortalVisibilityPRIVATE),
				sdkTopNavbarEntry("/alpha/docs", sdk.PortalVisibilityPRIVATE),
			)),
		),
	)

	It("detects a console-side flip of an entry the CRD leaves unset", func() {
		crd := apimportal.DriftDTO(sdkPortalWithTopNavbar(sdkTopNavbarEntry("/alpha", "")))
		remote := sdkPortalWithTopNavbar(sdkTopNavbarEntry("/alpha", sdk.PortalVisibilityPRIVATE))
		result := drift.DetectWithNamespace(crd, remote, "")
		Expect(result.DriftDetected()).To(BeTrue())
		Expect(result.String()).To(ContainSubstring(`visibility: PUBLIC != PRIVATE`))
	})

	It("leaves the payload DriftDTO resolves untouched", func() {
		sent := sdkPortalWithTopNavbar(sdkTopNavbarEntry("/alpha", ""))
		apimportal.DriftDTO(sent)
		Expect(sent.Structure.TopNavbar[0].Visibility).To(BeNil())
	})

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(sdkCompletePortalDTO(), sdkCompletePortalDTO())
		})
		It("ensure no new property isn't tested for structure field", func() {
			expectedEquivalentNotHavingAnyZeroValue(sdkCompletePortalDTOWithStructure(), sdkCompletePortalDTOWithStructure())
		})
	})
})

func sdkTopNavbarEntry(path string, visibility sdk.PortalVisibility) sdk.PortalNavigationPath {
	entry := sdk.PortalNavigationPath{Path: path}
	if visibility != "" {
		entry.Visibility = &visibility
	}
	return entry
}

func sdkPortalWithTopNavbar(entries ...sdk.PortalNavigationPath) sdk.Portal {
	return sdk.Portal{Structure: &sdk.PortalNavigationStructure{TopNavbar: entries}}
}

func sdkPortalToDTO(prtl *v1alpha1.Portal) sdk.Portal {
	GinkgoHelper()
	dto, err := apimportal.ToDTO(prtl)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func sdkCompletePortalDTO() sdk.Portal {
	GinkgoHelper()
	return loadFixture[sdk.Portal]("portal_sdk_dto.json")
}

func sdkCompletePortalCRD() *v1alpha1.Portal {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.Portal]("portal_crd.json")
	return &fixture
}

func sdkCompletePortalDTOWithStructure() sdk.Portal {
	GinkgoHelper()
	return loadFixture[sdk.Portal]("portal_dto_with_structure.json")
}

func sdkCompletePortalCRDWithStructure() *v1alpha1.Portal {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.Portal]("portal_crd_with_structure.json")
	return &fixture
}
