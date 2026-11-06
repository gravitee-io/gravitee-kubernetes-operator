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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
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
			model.AuthzEntityDTO{},
			model.AuthzEntityDTO{},
		),
		Entry("equal struct",
			completeAuthzEntityDTO(),
			completeAuthzEntityDTO(),
		),
		Entry("minimal CRD against what the platform answers for it",
			model.ToAuthzEntityDTO(minimalAuthzEntityCRD(), nil),
			parseJSON[model.AuthzEntityDTO](minimalAuthzEntityRemote),
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
		Entry("absent parents against the platform's []",
			withAuthzEntityParents(completeAuthzEntityDTO(), nil),
			withAuthzEntityParents(completeAuthzEntityDTO(), []model.AuthzEntityParentDTO{}),
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
			model.ToAuthzEntityDTO(minimalAuthzEntityCRD(), nil),
			withAuthzEntityAttributes(parseJSON[model.AuthzEntityDTO](minimalAuthzEntityRemote), `{"level": 3}`),
		),
		Entry("a parent added on the platform",
			completeAuthzEntityDTO(),
			withAuthzEntityParents(completeAuthzEntityDTO(), append(completeAuthzEntityDTO().Parents,
				model.AuthzEntityParentDTO{EntityType: "Group", EntityID: "finance"})),
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
			withAuthzEntityActive(completeAuthzEntityDTO(), false),
		),
	)

	It("names the parent the platform added", func() {
		remote := withAuthzEntityParents(completeAuthzEntityDTO(), append(completeAuthzEntityDTO().Parents,
			model.AuthzEntityParentDTO{EntityType: "Group", EntityID: "finance"}))

		result := drift.DetectWithNamespace(completeAuthzEntityDTO(), remote, "")

		Expect(result.String()).To(ContainSubstring(`+ {"entityType":"Group","entityId":"finance"}`))
	})

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeAuthzEntityDTO(), completeAuthzEntityDTO())
		})
	})
})

func completeAuthzEntityDTO() model.AuthzEntityDTO {
	GinkgoHelper()
	return loadFixture[model.AuthzEntityDTO]("authzentity_full_dto.json")
}

func minimalAuthzEntityCRD() *v1alpha1.AuthzEntity {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.AuthzEntity]("authzentity_minimal_crd.json")
	return &fixture
}

func withAuthzEntityHRID(dto model.AuthzEntityDTO, hrid string) model.AuthzEntityDTO {
	dto.HRID = hrid
	return dto
}

func withAuthzEntityAttributes(dto model.AuthzEntityDTO, attributes string) model.AuthzEntityDTO {
	GinkgoHelper()
	dto.Attributes = parseJSON[model.AuthzEntityDTO](`{"attributes": ` + attributes + `}`).Attributes
	return dto
}

func withAuthzEntityParents(dto model.AuthzEntityDTO, parents []model.AuthzEntityParentDTO) model.AuthzEntityDTO {
	dto.Parents = parents
	return dto
}

func withAuthzEntityTargets(dto model.AuthzEntityDTO, targets []string) model.AuthzEntityDTO {
	dto.Targets = targets
	return dto
}

func withAuthzEntityActive(dto model.AuthzEntityDTO, active bool) model.AuthzEntityDTO {
	dto.Active = active
	return dto
}

func reversed[T any](items []T) []T {
	out := slices.Clone(items)
	slices.Reverse(out)
	return out
}
