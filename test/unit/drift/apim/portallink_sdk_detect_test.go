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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	apimportallink "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portallink"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PortalLink Drift detection on the APIM SDK types", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			apimportallink.DTO{},
			apimportallink.DTO{},
		),
		Entry("equal struct",
			sdkCompletePortalLinkDTO(),
			sdkCompletePortalLinkDTO(),
		),
		Entry("equal struct from CRD mapping",
			sdkCompletePortalLinkDTO(),
			sdkPortalLinkToDTO(sdkCompletePortalLinkCRD()),
		),
		Entry("unset visibility equivalent to any remote visibility resolved by APIM",
			apimportallink.DTO{},
			apimportallink.DTO{PortalLink: sdk.PortalLink{Visibility: new(sdk.PortalVisibilityPRIVATE)}},
		),
		Entry("parent, which is not part of the payload",
			apimportallink.DTO{Parent: navigation.Parent{PortalHrid: "default-a"}},
			apimportallink.DTO{Parent: navigation.Parent{ApiHrid: "default-b"}},
		),
	)

	It("detects drift when the CRD declares a visibility and the remote differs", func() {
		crd := apimportallink.DTO{PortalLink: sdk.PortalLink{Visibility: new(sdk.PortalVisibilityPUBLIC)}}
		remote := apimportallink.DTO{PortalLink: sdk.PortalLink{Visibility: new(sdk.PortalVisibilityPRIVATE)}}
		result := drift.DetectWithNamespace(crd, remote, "")
		Expect(result.DriftDetected()).To(BeTrue())
		Expect(result.String()).To(ContainSubstring(`visibility: PUBLIC != PRIVATE`))
	})

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(sdkCompletePortalLinkDTO(), sdkCompletePortalLinkDTO())
		})
	})
})

func sdkPortalLinkToDTO(link *v1alpha1.PortalLink) apimportallink.DTO {
	GinkgoHelper()
	dto, err := apimportallink.ToDTO(link)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func sdkCompletePortalLinkDTO() apimportallink.DTO {
	GinkgoHelper()
	return apimportallink.DTO{
		PortalLink: loadFixture[sdk.PortalLink]("portallink_dto.json"),
		Parent:     navigation.Parent{PortalHrid: "default-my-portal"},
	}
}

func sdkCompletePortalLinkCRD() *v1alpha1.PortalLink {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.PortalLink]("portallink_crd.json")
	return &fixture
}
