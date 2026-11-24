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
	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimauthzpdp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/authzpdp"
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
			sdk.AuthzPdp{},
			sdk.AuthzPdp{},
		),
		Entry("equal struct",
			completeAuthzPdpDTO(),
			completeAuthzPdpDTO(),
		),
		Entry("minimal CRD against what the platform stores for it",
			minimalAuthzPdpDTO(),
			minimalAuthzPdpDTO(),
		),
		Entry("no endpoint against an empty endpoint",
			withAuthzPdpEndpoint(completeAuthzPdpDTO(), nil),
			withAuthzPdpEndpoint(completeAuthzPdpDTO(), &sdk.AuthzPdpEndpoint{}),
		),
		Entry("an empty endpoint against no endpoint",
			withAuthzPdpEndpoint(completeAuthzPdpDTO(), &sdk.AuthzPdpEndpoint{}),
			withAuthzPdpEndpoint(completeAuthzPdpDTO(), nil),
		),
		Entry("strict false against strict left out",
			withAuthzPdpStrict(completeAuthzPdpDTO(), new(false)),
			withAuthzPdpStrict(completeAuthzPdpDTO(), nil),
		),
		// GKO's DTO held the tag as a string: left out and "" compare equal.
		Entry("tag left out against an empty tag",
			withAuthzPdpTag(completeAuthzPdpDTO(), nil),
			withAuthzPdpTag(completeAuthzPdpDTO(), new("")),
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
			withAuthzPdpStrict(completeAuthzPdpDTO(), new(false)),
		),
		Entry("endpoint path changed on the platform",
			completeAuthzPdpDTO(),
			withAuthzPdpEndpoint(completeAuthzPdpDTO(), &sdk.AuthzPdpEndpoint{Path: "/authz/elsewhere"}),
		),
		Entry("endpoint added on the platform",
			minimalAuthzPdpDTO(),
			withAuthzPdpEndpoint(minimalAuthzPdpDTO(),
				&sdk.AuthzPdpEndpoint{Path: "/authz/drift-minimal"}),
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeAuthzPdpDTO(), completeAuthzPdpDTO())
		})
	})
})

func completeAuthzPdpDTO() sdk.AuthzPdp {
	GinkgoHelper()
	return loadFixture[sdk.AuthzPdp]("authzpdp_full_dto.json")
}

func minimalAuthzPdpCRD() *v1alpha1.AuthzPdp {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.AuthzPdp]("authzpdp_minimal_crd.json")
	return &fixture
}

func minimalAuthzPdpDTO() sdk.AuthzPdp {
	GinkgoHelper()
	dto, err := apimauthzpdp.ToDTO(minimalAuthzPdpCRD())
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func withAuthzPdpHRID(dto sdk.AuthzPdp, hrid string) sdk.AuthzPdp {
	dto.Hrid = hrid
	return dto
}

func withAuthzPdpName(dto sdk.AuthzPdp, name string) sdk.AuthzPdp {
	dto.Name = name
	return dto
}

func withAuthzPdpStrict(dto sdk.AuthzPdp, strict *bool) sdk.AuthzPdp {
	dto.Strict = strict
	return dto
}

// withAuthzPdpTag sets the tag, or leaves it out when tag is nil.
func withAuthzPdpTag(dto sdk.AuthzPdp, tag *string) sdk.AuthzPdp {
	dto.Tag = tag
	return dto
}

func withAuthzPdpEndpoint(dto sdk.AuthzPdp, endpoint *sdk.AuthzPdpEndpoint) sdk.AuthzPdp {
	dto.Endpoint = endpoint
	return dto
}
