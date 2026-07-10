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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Group Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			model.GroupDTO{},
			model.GroupDTO{},
		),
		Entry("equal struct",
			completeGroupDTO(),
			completeGroupDTO(),
		),
		Entry("notifyMembers ignored and empty members equivalent to nil",
			model.GroupDTO{
				ID:   "123456",
				HRID: "my-group",
				Name: "My Group",
			},
			model.GroupDTO{
				ID:            "123456",
				HRID:          "my-group",
				Name:          "My Group",
				NotifyMembers: true,
				Members:       make([]model.Member, 0),
			},
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeGroupDTO(), completeGroupDTO())
		})
	})
})

func completeGroupDTO() model.GroupDTO {
	GinkgoHelper()
	return loadFixture[model.GroupDTO]("group_dto.json")
}
