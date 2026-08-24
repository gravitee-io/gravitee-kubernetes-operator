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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Portal Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			model.PortalDTO{},
			model.PortalDTO{},
		),
		Entry("equal struct",
			completePortalDTO(),
			completePortalDTO(),
		),
		Entry("equal struct from CRD mapping",
			completePortalDTO(),
			model.ToPortalDTO(
				completePortalCRD().Spec.Type,
				refs.NewNamespacedNameFromObject(completePortalCRD()).HRID(),
				completePortalCRD().ActiveThemeHRID(),
			),
		),
		// Tests for the new structure field
		Entry("equal struct with structure",
			completePortalDTOWithStructure(),
			completePortalDTOWithStructure(),
		),
		Entry("equal struct with structure from CRD mapping",
			completePortalDTOWithStructure(),
			model.ToPortalDTO(
				completePortalCRDWithStructure().Spec.Type,
				refs.NewNamespacedNameFromObject(completePortalCRDWithStructure()).HRID(),
				completePortalCRDWithStructure().ActiveThemeHRID(),
			),
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completePortalDTO(), completePortalDTO())
		})
		It("ensure no new property isn't tested for structure field", func() {
			expectedEquivalentNotHavingAnyZeroValue(completePortalDTOWithStructure(), completePortalDTOWithStructure())
		})
	})
})

func completePortalDTO() model.PortalDTO {
	GinkgoHelper()
	return loadFixture[model.PortalDTO]("portal_dto.json")
}

func completePortalCRD() *v1alpha1.Portal {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.Portal]("portal_crd.json")
	return &fixture
}

func completePortalDTOWithStructure() model.PortalDTO {
	GinkgoHelper()
	return loadFixture[model.PortalDTO]("portal_dto_with_structure.json")
}

func completePortalCRDWithStructure() *v1alpha1.Portal {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.Portal]("portal_crd_with_structure.json")
	return &fixture
}
