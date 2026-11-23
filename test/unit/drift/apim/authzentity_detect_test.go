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
	"slices"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimauthzentity "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// minimalAuthzEntityRemote is what the platform answers for the minimal CRD: it reads absent
// attributes, parents and targets as empty, and an absent active as true.
const minimalAuthzEntityRemote = `{
	"hrid": "default-drift-minimal", "entityId": "drift-minimal", "kind": "PRINCIPAL", "entityType": "User",
	"attributes": {}, "parents": [], "active": true, "targets": []
}`

var _ = Describe("Authz entity Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			sdk.AuthzEntity{},
			sdk.AuthzEntity{},
		),
		Entry("equal struct",
			completeAuthzEntityDTO(),
			completeAuthzEntityDTO(),
		),
		Entry("minimal CRD against what the platform answers for it",
			minimalAuthzEntityDTO(),
			parseJSON[sdk.AuthzEntity](minimalAuthzEntityRemote),
		),
		Entry("hrid is identity, not payload",
			withAuthzEntityHRID(completeAuthzEntityDTO(), "default-drift-full"),
			withAuthzEntityHRID(completeAuthzEntityDTO(), "something-else"),
		),
		Entry("attribute keys in another order: the platform sorts them",
			completeAuthzEntityDTO(),
			withAuthzEntityAttributes(completeAuthzEntityDTO(), `{
				"badges": ["gold", "silver"], "department": "travel", "address": {"country": "France", "city": "Paris"}
			}`),
		),
		Entry("parents in another order: the platform sorts them",
			completeAuthzEntityDTO(),
			withAuthzEntityParents(completeAuthzEntityDTO(), reversed(completeAuthzEntityDTO().Parents)),
		),
		Entry("targets in another order: the platform sorts them",
			completeAuthzEntityDTO(),
			withAuthzEntityTargets(completeAuthzEntityDTO(), []string{"default", "drift-full@eu"}),
		),
		// The platform always answers active: GKO's DTO held it as a bool, absent and false compare equal.
		Entry("active false against active left out",
			withAuthzEntityActive(completeAuthzEntityDTO(), new(false)),
			withAuthzEntityActive(completeAuthzEntityDTO(), nil),
		),
		Entry("absent parents against the platform's []",
			withAuthzEntityParents(completeAuthzEntityDTO(), nil),
			withAuthzEntityParents(completeAuthzEntityDTO(), []sdk.AuthzEntityParent{}),
		),
	)

	DescribeTable("drifted values",
		func(crd, remote any) {
			result := drift.DetectWithNamespace(crd, remote, "")
			Expect(result.DriftDetected()).To(BeTrue())
		},
		Entry("an attribute value changed on the platform",
			completeAuthzEntityDTO(),
			withAuthzEntityAttributes(completeAuthzEntityDTO(), `{
				"department": "sales", "address": {"city": "Paris", "country": "France"}, "badges": ["gold", "silver"]
			}`),
		),
		Entry("a nested attribute value changed on the platform",
			completeAuthzEntityDTO(),
			withAuthzEntityAttributes(completeAuthzEntityDTO(), `{
				"department": "travel", "address": {"city": "Lyon", "country": "France"}, "badges": ["gold", "silver"]
			}`),
		),
		Entry("a list inside an attribute reordered: the platform stores it as sent",
			completeAuthzEntityDTO(),
			withAuthzEntityAttributes(completeAuthzEntityDTO(), `{
				"department": "travel", "address": {"city": "Paris", "country": "France"}, "badges": ["silver", "gold"]
			}`),
		),
		Entry("an attribute added on the platform",
			minimalAuthzEntityDTO(),
			withAuthzEntityAttributes(parseJSON[sdk.AuthzEntity](minimalAuthzEntityRemote), `{"level": 3}`),
		),
		Entry("a parent added on the platform",
			completeAuthzEntityDTO(),
			withAuthzEntityParents(completeAuthzEntityDTO(), append(completeAuthzEntityDTO().Parents,
				sdk.AuthzEntityParent{EntityType: "Group", EntityId: "finance"})),
		),
		Entry("a parent removed on the platform",
			completeAuthzEntityDTO(),
			withAuthzEntityParents(completeAuthzEntityDTO(), completeAuthzEntityDTO().Parents[:1]),
		),
		Entry("a target removed on the platform",
			completeAuthzEntityDTO(),
			withAuthzEntityTargets(completeAuthzEntityDTO(), []string{"default"}),
		),
		Entry("deactivated on the platform",
			completeAuthzEntityDTO(),
			withAuthzEntityActive(completeAuthzEntityDTO(), new(false)),
		),
	)

	It("names the parent the platform added", func() {
		remote := withAuthzEntityParents(completeAuthzEntityDTO(), append(completeAuthzEntityDTO().Parents,
			sdk.AuthzEntityParent{EntityType: "Group", EntityId: "finance"}))

		result := drift.DetectWithNamespace(completeAuthzEntityDTO(), remote, "")

		Expect(result.String()).To(ContainSubstring(`+ {"entityId":"finance","entityType":"Group"}`))
	})

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeAuthzEntityDTO(), completeAuthzEntityDTO())
		})
	})

	It("refuses a parent APIM added when neither version of the entity sets parents", func() {
		remote := sdk.AuthzEntity{
			Kind:    "before",
			Parents: []sdk.AuthzEntityParent{{EntityType: "Group", EntityId: "admins"}},
		}
		merged := drift.Merge(
			drift.DetectWithNamespace(sdk.AuthzEntity{Kind: "before"}, remote, ""),
			drift.DetectWithNamespace(sdk.AuthzEntity{Kind: "after"}, remote, ""),
		)

		Expect(merged.DriftDetected()).To(BeTrue(), merged.String())
	})
})

func completeAuthzEntityDTO() sdk.AuthzEntity {
	GinkgoHelper()
	return loadFixture[sdk.AuthzEntity]("authzentity_full_dto.json")
}

func minimalAuthzEntityCRD() *v1alpha1.AuthzEntity {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.AuthzEntity]("authzentity_minimal_crd.json")
	return &fixture
}

func minimalAuthzEntityDTO() sdk.AuthzEntity {
	GinkgoHelper()
	dto, err := apimauthzentity.ToDTO(minimalAuthzEntityCRD())
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func withAuthzEntityHRID(dto sdk.AuthzEntity, hrid string) sdk.AuthzEntity {
	dto.Hrid = hrid
	return dto
}

func withAuthzEntityAttributes(dto sdk.AuthzEntity, attributes string) sdk.AuthzEntity {
	GinkgoHelper()
	dto.Attributes = parseJSON[sdk.AuthzEntity](`{"attributes": ` + attributes + `}`).Attributes
	return dto
}

func withAuthzEntityParents(dto sdk.AuthzEntity, parents []sdk.AuthzEntityParent) sdk.AuthzEntity {
	dto.Parents = parents
	return dto
}

func withAuthzEntityTargets(dto sdk.AuthzEntity, targets []string) sdk.AuthzEntity {
	dto.Targets = targets
	return dto
}

// withAuthzEntityActive sets active, or leaves it out when active is nil.
func withAuthzEntityActive(dto sdk.AuthzEntity, active *bool) sdk.AuthzEntity {
	dto.Active = active
	return dto
}

func reversed[T any](items []T) []T {
	out := slices.Clone(items)
	slices.Reverse(out)
	return out
}
