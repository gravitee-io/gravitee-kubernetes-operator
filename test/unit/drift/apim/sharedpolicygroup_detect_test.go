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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/sharedpolicygroups"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	spg "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/sharedpolicygroup"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SharedPolicyGroup Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			spg.DTO{},
			spg.DTO{},
		),
		Entry("equal struct",
			completeSharedPolicyGroupDTO(),
			completeSharedPolicyGroupDTO(),
		),
		Entry("equal struct from CRD mapping",
			completeSharedPolicyGroupDTO(),
			toSharedPolicyGroupDTO(completeSharedPolicyGroupSpec()),
		),
		Entry("empty steps equivalent to nil",
			spg.DTO{SharedPolicyGroup: sdk.SharedPolicyGroup{
				Hrid:    "my-spg",
				Name:    "My SPG",
				ApiType: "PROXY",
			}},
			spg.DTO{SharedPolicyGroup: sdk.SharedPolicyGroup{
				Hrid:    "my-spg",
				Name:    "My SPG",
				ApiType: "PROXY",
				Steps:   make([]sdk.StepV4, 0),
			}},
		),
	)

	It("admits the defaults APIM saves in a step configuration", func() {
		crd := sharedPolicyGroupWithStep(`{"addHeaders": true, "rate": {"limit": 10}}`)
		remote := sharedPolicyGroupWithStep(`{"addHeaders": true, "rate": {"limit": 10}, "async": false}`)

		expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
	})

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(toSharedPolicyGroupDTO(completeSharedPolicyGroupSpec()), completeSharedPolicyGroupDTO())
		})
	})
})

func completeSharedPolicyGroupSpec() sharedpolicygroups.SharedPolicyGroup {
	GinkgoHelper()
	return loadFixture[sharedpolicygroups.SharedPolicyGroup]("shared_policy_group_spec.json")
}

// toSharedPolicyGroupDTO maps a spec the way the operator does, through a CR that holds it.
func toSharedPolicyGroupDTO(spec sharedpolicygroups.SharedPolicyGroup) spg.DTO {
	GinkgoHelper()
	dto, err := spg.ToDTO(&v1alpha1.SharedPolicyGroup{Spec: v1alpha1.SharedPolicyGroupSpec{SharedPolicyGroup: &spec}})
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func completeSharedPolicyGroupDTO() spg.DTO {
	GinkgoHelper()
	return loadFixture[spg.DTO]("shared_policy_group_dto.json")
}

func sharedPolicyGroupWithStep(configuration string) spg.DTO {
	GinkgoHelper()
	return parseJSON[spg.DTO](`{"name": "spg", "apiType": "PROXY", "phase": "REQUEST", "steps": [` +
		`{"enabled": true, "policy": "rate-limit", "name": "Rate Limit", "configuration": ` + configuration + `}]}`)
}
