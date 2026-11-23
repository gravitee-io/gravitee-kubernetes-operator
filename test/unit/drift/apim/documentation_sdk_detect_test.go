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
	apimdocs "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/docs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func sdkDocWith(page sdk.Documentation) apimdocs.DTO {
	return apimdocs.DTO{Documentation: page}
}

var _ = Describe("Documentation Drift detection on the APIM SDK types", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			apimdocs.DTO{},
			apimdocs.DTO{},
		),
		Entry("equal struct",
			sdkCompleteDocumentationDTO(),
			sdkCompleteDocumentationDTO(),
		),
		Entry("equal struct from CRD mapping",
			sdkCompleteDocumentationDTO(),
			sdkDocumentationToDTO(sdkCompleteDocumentationCRD()),
		),
		Entry("empty area equivalent to remote TOP_NAVBAR",
			apimdocs.DTO{},
			sdkDocWith(sdk.Documentation{Area: new(sdk.TOPNAVBAR)}),
		),
		Entry("unset visibility equivalent to any remote visibility resolved by APIM",
			apimdocs.DTO{},
			sdkDocWith(sdk.Documentation{Visibility: new(sdk.PortalVisibilityPRIVATE)}),
		),
		Entry("parent, which is not part of the payload",
			apimdocs.DTO{Parent: navigation.Parent{PortalHrid: "default-a"}},
			apimdocs.DTO{Parent: navigation.Parent{ApiHrid: "default-b"}},
		),
	)

	It("compares an area the CRD sets against the remote default", func() {
		crd := sdkDocWith(sdk.Documentation{Area: new(sdk.HOMEPAGE)})
		remote := sdkDocWith(sdk.Documentation{Area: new(sdk.TOPNAVBAR)})
		result := drift.DetectWithNamespace(crd, remote, "")
		Expect(result.DriftDetected()).To(BeTrue())
		Expect(result.String()).To(ContainSubstring(`area: HOMEPAGE != TOP_NAVBAR`))
	})

	It("detects drift when the CRD declares a visibility and the remote differs", func() {
		crd := sdkDocWith(sdk.Documentation{Visibility: new(sdk.PortalVisibilityPUBLIC)})
		remote := sdkDocWith(sdk.Documentation{Visibility: new(sdk.PortalVisibilityPRIVATE)})
		result := drift.DetectWithNamespace(crd, remote, "")
		Expect(result.DriftDetected()).To(BeTrue())
		Expect(result.String()).To(ContainSubstring(`visibility: PUBLIC != PRIVATE`))
	})

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(sdkCompleteDocumentationDTO(), sdkCompleteDocumentationDTO())
		})
	})
})

func sdkDocumentationToDTO(doc *v1alpha1.Documentation) apimdocs.DTO {
	GinkgoHelper()
	dto, err := apimdocs.ToDTO(doc)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func sdkCompleteDocumentationDTO() apimdocs.DTO {
	GinkgoHelper()
	return apimdocs.DTO{
		Documentation: loadFixture[sdk.Documentation]("documentation_dto.json"),
		Parent:        navigation.Parent{PortalHrid: "default-my-portal"},
	}
}

func sdkCompleteDocumentationCRD() *v1alpha1.Documentation {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.Documentation]("documentation_crd.json")
	return &fixture
}
