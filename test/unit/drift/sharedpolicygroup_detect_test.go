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

package drift

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/sharedpolicygroups"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("SharedPolicyGroup Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			model.SharedPolicyGroupDTO{},
			model.SharedPolicyGroupDTO{},
		),
		Entry("equal struct",
			completeSharedPolicyGroupDTO(),
			completeSharedPolicyGroupDTO(),
		),
		Entry("equal struct from CRD mapping",
			completeSharedPolicyGroupDTO(),
			model.ToSharePolicyGroupDTO(completeSharedPolicyGroupSpec()),
		),
		Entry("empty steps equivalent to nil",
			model.SharedPolicyGroupDTO{
				HRID:    "my-spg",
				Name:    "My SPG",
				ApiType: "PROXY",
			},
			model.SharedPolicyGroupDTO{
				HRID:    "my-spg",
				Name:    "My SPG",
				ApiType: "PROXY",
				Steps:   make([]model.StepDTO, 0),
			},
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(model.ToSharePolicyGroupDTO(completeSharedPolicyGroupSpec()), completeSharedPolicyGroupDTO())
		})
	})
})

func completeSharedPolicyGroupSpec() sharedpolicygroups.SharedPolicyGroup {
	GinkgoHelper()
	return loadFixture[sharedpolicygroups.SharedPolicyGroup]("shared_policy_group_spec.json")
}

func completeSharedPolicyGroupDTO() model.SharedPolicyGroupDTO {
	GinkgoHelper()
	return loadFixture[model.SharedPolicyGroupDTO]("shared_policy_group_dto.json")
}
