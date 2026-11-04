// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apim

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Authz PDP Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			model.AuthzPdpDTO{},
			model.AuthzPdpDTO{},
		),
		Entry("equal struct",
			completeAuthzPdpDTO(),
			completeAuthzPdpDTO(),
		),
		Entry("minimal CRD against what the platform stores for it",
			model.ToAuthzPdpDTO(minimalAuthzPdpCRD()),
			model.ToAuthzPdpDTO(minimalAuthzPdpCRD()),
		),
		Entry("no endpoint against an empty endpoint",
			withAuthzPdpEndpoint(completeAuthzPdpDTO(), nil),
			withAuthzPdpEndpoint(completeAuthzPdpDTO(), &model.AuthzPdpEndpointDTO{}),
		),
		Entry("an empty endpoint against no endpoint",
			withAuthzPdpEndpoint(completeAuthzPdpDTO(), &model.AuthzPdpEndpointDTO{}),
			withAuthzPdpEndpoint(completeAuthzPdpDTO(), nil),
		),
		Entry("hrid is identity, not payload",
			withAuthzPdpHRID(completeAuthzPdpDTO(), "default-drift-full"),
			withAuthzPdpHRID(completeAuthzPdpDTO(), "something-else"),
		),
	)

	DescribeTable("drifted values",
		func(crd, remote any) {
			result := drift.DetectWithNamespace(crd, remote, "")
			Expect(result.DriftDetected()).To(BeTrue())
		},
		Entry("name changed on the platform",
			completeAuthzPdpDTO(),
			withAuthzPdpName(completeAuthzPdpDTO(), "Renamed in the Console"),
		),
		Entry("strict switched off by a sibling of the same targetPdpId",
			completeAuthzPdpDTO(),
			withAuthzPdpStrict(completeAuthzPdpDTO(), false),
		),
		Entry("endpoint path changed on the platform",
			completeAuthzPdpDTO(),
			withAuthzPdpEndpoint(completeAuthzPdpDTO(), &model.AuthzPdpEndpointDTO{Path: "/authz/elsewhere"}),
		),
		Entry("endpoint added on the platform",
			model.ToAuthzPdpDTO(minimalAuthzPdpCRD()),
			withAuthzPdpEndpoint(model.ToAuthzPdpDTO(minimalAuthzPdpCRD()),
				&model.AuthzPdpEndpointDTO{Path: "/authz/drift-minimal"}),
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeAuthzPdpDTO(), completeAuthzPdpDTO())
		})
	})
})

func completeAuthzPdpDTO() model.AuthzPdpDTO {
	GinkgoHelper()
	return loadFixture[model.AuthzPdpDTO]("authzpdp_full_dto.json")
}

func minimalAuthzPdpCRD() *v1alpha1.AuthzPdp {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.AuthzPdp]("authzpdp_minimal_crd.json")
	return &fixture
}

func withAuthzPdpHRID(dto model.AuthzPdpDTO, hrid string) model.AuthzPdpDTO {
	dto.HRID = hrid
	return dto
}

func withAuthzPdpName(dto model.AuthzPdpDTO, name string) model.AuthzPdpDTO {
	dto.Name = name
	return dto
}

func withAuthzPdpStrict(dto model.AuthzPdpDTO, strict bool) model.AuthzPdpDTO {
	dto.Strict = strict
	return dto
}

func withAuthzPdpEndpoint(dto model.AuthzPdpDTO, endpoint *model.AuthzPdpEndpointDTO) model.AuthzPdpDTO {
	dto.Endpoint = endpoint
	return dto
}
