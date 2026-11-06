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
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// minimalAuthzPolicyRemote is what the platform answers for the minimal CRD: it reads absent
// targets as empty.
const minimalAuthzPolicyRemote = `{
	"hrid": "default-drift-minimal", "entityId": "policy.drift-minimal", "name": "Drift minimal",
	"policyType": "CUSTOM", "policyText": "permit(principal, action, resource);", "targets": []
}`

var _ = Describe("Authz policy Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			model.AuthzPolicyDTO{},
			model.AuthzPolicyDTO{},
		),
		Entry("equal struct",
			completeAuthzPolicyDTO(),
			completeAuthzPolicyDTO(),
		),
		Entry("minimal CRD against what the platform answers for it",
			model.ToAuthzPolicyDTO(minimalAuthzPolicyCRD(), nil),
			parseJSON[model.AuthzPolicyDTO](minimalAuthzPolicyRemote),
		),
		Entry("hrid is identity, not payload",
			withAuthzPolicyHRID(completeAuthzPolicyDTO(), "default-drift-full"),
			withAuthzPolicyHRID(completeAuthzPolicyDTO(), "something-else"),
		),
		Entry("targets in another order: the platform sorts them",
			completeAuthzPolicyDTO(),
			withAuthzPolicyTargets(completeAuthzPolicyDTO(), []string{"default", "drift-full@eu"}),
		),
	)

	// The platform stores and answers the text as sent: any difference, whitespace included, is a
	// change someone made there.
	DescribeTable("drifted values",
		func(crd, remote any) {
			result := drift.DetectWithNamespace(crd, remote, "")
			Expect(result.DriftDetected()).To(BeTrue())
		},
		Entry("the text changed on the platform",
			completeAuthzPolicyDTO(),
			withAuthzPolicyText(completeAuthzPolicyDTO(), "permit(principal, action, resource);\n"),
		),
		Entry("the text lost its trailing newline on the platform",
			completeAuthzPolicyDTO(),
			withAuthzPolicyText(completeAuthzPolicyDTO(),
				strings.TrimSuffix(completeAuthzPolicyDTO().PolicyText, "\n")),
		),
		Entry("renamed on the platform",
			completeAuthzPolicyDTO(),
			withAuthzPolicyName(completeAuthzPolicyDTO(), "Drift renamed"),
		),
		Entry("a target removed on the platform",
			completeAuthzPolicyDTO(),
			withAuthzPolicyTargets(completeAuthzPolicyDTO(), []string{"default"}),
		),
		Entry("a target added on the platform",
			model.ToAuthzPolicyDTO(minimalAuthzPolicyCRD(), nil),
			withAuthzPolicyTargets(parseJSON[model.AuthzPolicyDTO](minimalAuthzPolicyRemote), []string{"default"}),
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeAuthzPolicyDTO(), completeAuthzPolicyDTO())
		})
	})
})

func completeAuthzPolicyDTO() model.AuthzPolicyDTO {
	GinkgoHelper()
	return loadFixture[model.AuthzPolicyDTO]("authzpolicy_full_dto.json")
}

func minimalAuthzPolicyCRD() *v1alpha1.AuthzPolicy {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.AuthzPolicy]("authzpolicy_minimal_crd.json")
	return &fixture
}

func withAuthzPolicyHRID(dto model.AuthzPolicyDTO, hrid string) model.AuthzPolicyDTO {
	dto.HRID = hrid
	return dto
}

func withAuthzPolicyText(dto model.AuthzPolicyDTO, text string) model.AuthzPolicyDTO {
	dto.PolicyText = text
	return dto
}

func withAuthzPolicyName(dto model.AuthzPolicyDTO, name string) model.AuthzPolicyDTO {
	dto.Name = name
	return dto
}

func withAuthzPolicyTargets(dto model.AuthzPolicyDTO, targets []string) model.AuthzPolicyDTO {
	dto.Targets = targets
	return dto
}
