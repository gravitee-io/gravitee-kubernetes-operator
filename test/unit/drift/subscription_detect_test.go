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

var _ = Describe("Application Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.Detect(crd, remote))
		},
		Entry("empty struct",
			model.SubscriptionDTO{},
			model.SubscriptionDTO{},
		),
		Entry("equal struct",
			completeSubscriptionDTO(),
			completeSubscriptionDTO(),
		),
		Entry("start/end date",
			model.SubscriptionDTO{
				ID:       "sub-id",
				ApiID:    "api-id",
				AppID:    "app-id",
				EndingAt: "2023-08-25T23:43:16Z",
				ApiKeys:  []model.ApiKeySpec{{Key: "key1", ExpireAt: ptr("2024-08-25T23:43:16-00:00")}},
			},
			model.SubscriptionDTO{
				ID:         "default-sub-id",
				ApiID:      "default-api-id",
				AppID:      "default-app-id",
				StartingAt: "2023-07-24T20:43:16-03:00",
				EndingAt:   "2023-08-25T22:43:16-01:00",
				ApiKeys:    []model.ApiKeySpec{{Key: "key1", ExpireAt: ptr("2024-08-25T23:43:16+00:00")}},
			},
		),
		Entry("empty collections",
			model.SubscriptionDTO{},
			model.SubscriptionDTO{
				ApiKeys:  []model.ApiKeySpec{},
				Metadata: map[string]string{},
			},
		),
	)
	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeSubscriptionDTO(), completeSubscriptionDTO())
		})
	})
})

func completeSubscriptionDTO() model.SubscriptionDTO {
	GinkgoHelper()
	return loadFixture[model.SubscriptionDTO]("subscription_dto.json")
}
