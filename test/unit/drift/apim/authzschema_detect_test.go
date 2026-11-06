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

var _ = Describe("Authz schema Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			model.AuthzSchemaDTO{},
			model.AuthzSchemaDTO{},
		),
		Entry("equal struct",
			completeAuthzSchemaDTO(),
			completeAuthzSchemaDTO(),
		),
		Entry("minimal CRD, no target, against what the platform stores for it",
			model.ToAuthzSchemaDTO(minimalAuthzSchemaCRD(), ""),
			model.ToAuthzSchemaDTO(minimalAuthzSchemaCRD(), ""),
		),
		Entry("hrid is identity, not payload",
			withAuthzSchemaHRID(completeAuthzSchemaDTO(), "default-drift-full"),
			withAuthzSchemaHRID(completeAuthzSchemaDTO(), "something-else"),
		),
	)

	DescribeTable("drifted values",
		func(crd, remote any) {
			result := drift.DetectWithNamespace(crd, remote, "")
			Expect(result.DriftDetected()).To(BeTrue())
		},
		Entry("name changed on the platform",
			completeAuthzSchemaDTO(),
			withAuthzSchemaName(completeAuthzSchemaDTO(), "Renamed in the Console"),
		),
		Entry("schema text edited on the platform",
			completeAuthzSchemaDTO(),
			withAuthzSchemaText(completeAuthzSchemaDTO(), "entity User;\n"),
		),
		Entry("a trailing newline added to the schema text: the platform stores it as sent",
			completeAuthzSchemaDTO(),
			withAuthzSchemaText(completeAuthzSchemaDTO(), completeAuthzSchemaDTO().Schema+"\n"),
		),
		Entry("moved to another gateway on the platform",
			completeAuthzSchemaDTO(),
			withAuthzSchemaTarget(completeAuthzSchemaDTO(), "drift-full@us"),
		),
		Entry("target removed on the platform",
			completeAuthzSchemaDTO(),
			withAuthzSchemaTarget(completeAuthzSchemaDTO(), ""),
		),
		Entry("target added on the platform to a schema declared without one",
			model.ToAuthzSchemaDTO(minimalAuthzSchemaCRD(), ""),
			withAuthzSchemaTarget(model.ToAuthzSchemaDTO(minimalAuthzSchemaCRD(), ""), "default"),
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeAuthzSchemaDTO(), completeAuthzSchemaDTO())
		})
	})
})

func completeAuthzSchemaDTO() model.AuthzSchemaDTO {
	GinkgoHelper()
	return loadFixture[model.AuthzSchemaDTO]("authzschema_full_dto.json")
}

func minimalAuthzSchemaCRD() *v1alpha1.AuthzSchema {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.AuthzSchema]("authzschema_minimal_crd.json")
	return &fixture
}

func withAuthzSchemaHRID(dto model.AuthzSchemaDTO, hrid string) model.AuthzSchemaDTO {
	dto.HRID = hrid
	return dto
}

func withAuthzSchemaName(dto model.AuthzSchemaDTO, name string) model.AuthzSchemaDTO {
	dto.Name = name
	return dto
}

func withAuthzSchemaText(dto model.AuthzSchemaDTO, text string) model.AuthzSchemaDTO {
	dto.Schema = text
	return dto
}

func withAuthzSchemaTarget(dto model.AuthzSchemaDTO, target string) model.AuthzSchemaDTO {
	dto.Target = target
	return dto
}
