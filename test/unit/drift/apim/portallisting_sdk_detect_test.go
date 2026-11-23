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
	apimportallisting "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/portallisting"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func sdkListingWithEntry(visibility *sdk.PortalVisibility) apimportallisting.DTO {
	return apimportallisting.DTO{PortalListing: sdk.PortalListing{
		Apis: []sdk.PortalListingApiEntry{{ApiHrid: "default-api", Location: "/alpha", Visibility: visibility}},
	}}
}

var _ = Describe("PortalListing Drift detection on the APIM SDK types", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			apimportallisting.DTO{},
			apimportallisting.DTO{},
		),
		Entry("equal struct",
			sdkCompletePortalListingDTO(),
			sdkCompletePortalListingDTO(),
		),
		Entry("equal struct from CRD mapping",
			sdkCompletePortalListingDTO(),
			sdkPortalListingToDTO(sdkCompletePortalListingCRD()),
		),
		Entry("unset visibility equivalent to any remote visibility resolved by APIM",
			sdkListingWithEntry(nil),
			sdkListingWithEntry(new(sdk.PortalVisibilityPRIVATE)),
		),
		Entry("portal HRID, which is not part of the payload",
			apimportallisting.DTO{PortalHrid: "default-a"},
			apimportallisting.DTO{PortalHrid: "default-b"},
		),
	)

	It("detects drift when the CRD declares a visibility and the remote differs", func() {
		crd := sdkListingWithEntry(new(sdk.PortalVisibilityPUBLIC))
		remote := sdkListingWithEntry(new(sdk.PortalVisibilityPRIVATE))
		result := drift.DetectWithNamespace(crd, remote, "")
		Expect(result.DriftDetected()).To(BeTrue())
		Expect(result.String()).To(ContainSubstring(`visibility: PUBLIC != PRIVATE`))
	})

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(sdkCompletePortalListingDTO(), sdkCompletePortalListingDTO())
		})
	})
})

func sdkPortalListingToDTO(listing *v1alpha1.PortalListing) apimportallisting.DTO {
	GinkgoHelper()
	dto, err := apimportallisting.ToDTO(listing)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func sdkCompletePortalListingDTO() apimportallisting.DTO {
	GinkgoHelper()
	return apimportallisting.DTO{
		PortalListing: loadFixture[sdk.PortalListing]("portal_listing_dto.json"),
		PortalHrid:    "default-my-portal",
	}
}

func sdkCompletePortalListingCRD() *v1alpha1.PortalListing {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.PortalListing]("portal_listing_crd.json")
	return &fixture
}
