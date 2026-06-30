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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
				StartingAt: "2023-07-25T02:43:16+03:00",
				EndingAt:   "2023-08-25T23:43:16Z",
				ApiKeys:    []model.ApiKeySpec{{Key: "key1", ExpireAt: ptr("2024-08-25T23:43:16-00:00")}},
			},
			model.SubscriptionDTO{
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
})

func completeSubscriptionDTO() model.SubscriptionDTO {
	GinkgoHelper()
	return model.SubscriptionDTO{
		ID:         "12346",
		ApiID:      "456798",
		AppID:      "789123",
		PlanID:     "keyless",
		StartingAt: "2023-08-25T23:43:16Z",
		EndingAt:   "2024-08-25T23:43:16Z",
		Metadata:   map[string]string{"foo": "bar", "baz": "puk"},
		ApiKeys:    []model.ApiKeySpec{{Key: "key1", ExpireAt: ptr("2024-08-25T23:43:16Z")}},
		ConsumerConfiguration: subscription.ConsumerConfiguration{
			EntrypointID: "entrypoint-id",
			Channel:      "channel",
			EntrypointConfiguration: &utils.GenericStringMap{
				Unstructured: unstructured.Unstructured{
					Object: map[string]interface{}{
						"callbackUrl": "https://webhook.site/bbd53b8c-e330-4881-b5ad-ddca91c52af1",
						"headers": []map[string]string{
							{"name": "X-Gravitee-Custom", "value": "Hello"},
						},
						"auth": map[string]interface{}{
							"type": "basic",
							"basic": map[string]string{
								"username": "admin",
								"password": "admin",
							},
						},
						"ssl": map[string]interface{}{
							"hostnameVerifier": true,
							"trustAll":         false,
						},
					},
				},
			},
		},
	}

}
