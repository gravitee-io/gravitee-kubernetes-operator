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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Subscription Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, "default"))
		},
		Entry("empty struct",
			sdk.Subscription{},
			sdk.Subscription{},
		),
		Entry("equal struct",
			completeSubscriptionDTO(),
			completeSubscriptionDTO(),
		),
		Entry("start/end date",
			sdk.Subscription{
				Hrid:            "sub-id",
				ApplicationHrid: "app-id",
				EndingAt:        new("2023-08-25T23:43:16Z"),
				ApiKeys:         []sdk.ApiKey{{Key: "key1", ExpireAt: new("2024-08-25T23:43:16-00:00")}},
			},
			sdk.Subscription{
				Hrid:            "default-sub-id",
				ApplicationHrid: "default-app-id",
				EndingAt:        new("2023-08-25T22:43:16-01:00"),
				ApiKeys:         []sdk.ApiKey{{Key: "key1", ExpireAt: new("2024-08-25T23:43:16+00:00")}},
			},
		),
		Entry("crd-only api keys ignored",
			sdk.Subscription{
				ApiKeys: []sdk.ApiKey{
					{Key: "key1"},
					{Key: "crd-only"},
				},
			},
			sdk.Subscription{
				ApiKeys: []sdk.ApiKey{
					{Key: "key1"},
				},
			},
		),
		Entry("empty collections",
			sdk.Subscription{},
			sdk.Subscription{
				ApiKeys:  []sdk.ApiKey{},
				Metadata: map[string]string{},
			},
		),
	)

	It("admits the defaults APIM saves in an entrypoint configuration", func() {
		crd := subscriptionWithEntrypointConfiguration(`{"callbackUrl": "https://example.com/hook", "ssl": {"hostnameVerifier": true}}`)
		remote := subscriptionWithEntrypointConfiguration(
			`{"callbackUrl": "https://example.com/hook", "ssl": {"hostnameVerifier": true, "trustAll": false}}`,
		)

		expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
	})

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeSubscriptionDTO(), completeSubscriptionDTO())
		})
	})
})

func completeSubscriptionDTO() sdk.Subscription {
	GinkgoHelper()
	return loadFixture[sdk.Subscription]("subscription_dto.json")
}

func subscriptionWithEntrypointConfiguration(configuration string) sdk.Subscription {
	GinkgoHelper()
	return parseJSON[sdk.Subscription](
		`{"consumerConfiguration": {"entrypointId": "webhook", "entrypointConfiguration": ` + configuration + `}}`,
	)
}
