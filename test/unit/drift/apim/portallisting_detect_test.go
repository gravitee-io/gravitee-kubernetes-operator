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
	nav "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("PortalListing Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			model.PortalListingDTO{},
			model.PortalListingDTO{},
		),
		Entry("equal struct",
			completePortalListingDTO(),
			completePortalListingDTO(),
		),
		Entry("equal struct from CRD mapping",
			completePortalListingDTO(),
			*service.ToPortalListingDTO(completePortalListingCRD()),
		),
		Entry("unset visibility equivalent to any remote visibility resolved by APIM",
			model.PortalListingDTO{
				APIs: []model.PortalListingApiEntryDTO{{ApiHrid: "default-api", Location: "/alpha"}},
			},
			model.PortalListingDTO{
				APIs: []model.PortalListingApiEntryDTO{{ApiHrid: "default-api", Location: "/alpha", Visibility: nav.Private}},
			},
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completePortalListingDTO(), completePortalListingDTO())
		})
	})
})

func completePortalListingDTO() model.PortalListingDTO {
	GinkgoHelper()
	return loadFixture[model.PortalListingDTO]("portal_listing_dto.json")
}

func completePortalListingCRD() *v1alpha1.PortalListing {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.PortalListing]("portal_listing_crd.json")
	return &fixture
}
