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

package usecase

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	submodel "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Usecase", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should create and update subscription consumer configuration for a webhook push plan", func() {
		fixtures := fixture.Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithAPIv4(constants.ApiV4WithApiKeyPlanFile).
			WithApplication(constants.ApplicationWithClientIDFile).
			WithSubscription(constants.SubscriptionFile).
			Build()

		clientID := random.GetName()
		fixtures.Application.Spec.Settings.App.ClientID = &clientID
		fixtures.Subscription.Spec.API.Name = fixtures.APIv4.Name
		fixtures.Subscription.Spec.App.Name = fixtures.Application.Name
		fixtures.Subscription.Spec.Plan = "API_KEY"
		fixtures.Subscription.Spec.ConsumerConfiguration = &submodel.ConsumerConfiguration{
			EntrypointID: "webhook",
			Channel:      "orders",
			EntrypointConfiguration: utils.ToGenericStringMap(map[string]interface{}{
				"callbackUrl": "https://example.com/webhook",
				"headers": []interface{}{
					map[string]interface{}{
						"name":  "demoHeader",
						"value": "my-value",
					},
					map[string]interface{}{
						"name":  "anotherHeader",
						"value": "my-value2",
					},
				},
			}),
		}

		By("creating API, application and subscription resources")
		fixtures.Apply()

		By("expecting subscription status to be completed and accepted")
		Eventually(func() error {
			return assert.SubscriptionCompleted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)
		Eventually(func() error {
			return assert.SubscriptionAccepted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)
		Eventually(func() error {
			return assert.ManagedByAutomationAPI(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		By("calling automation API and expecting consumer configuration to be saved")
		client := apim.NewClient(ctx)
		Eventually(func() error {
			apiHrid := refs.NewNamespacedNameFromObject(fixtures.APIv4).HRID()
			subHrid := refs.NewNamespacedNameFromObject(fixtures.Subscription).HRID()
			sub, err := client.Subscriptions.GetByHRID(apiHrid, subHrid)
			if err != nil {
				return err
			}

			err = assert.Equals("consumerConfiguration.entrypointId", "webhook", sub.ConsumerConfiguration.EntrypointID)
			if err != nil {
				return err
			}

			if err := assert.Equals("consumerConfiguration.channel", "orders", sub.ConsumerConfiguration.Channel); err != nil {
				return err
			}
			return assert.Equals(
				"consumerConfiguration.entrypointConfiguration",
				map[string]interface{}{
					"callbackUrl": "https://example.com/webhook",
					"headers": []interface{}{
						map[string]interface{}{
							"name":  "demoHeader",
							"value": "my-value",
						},
						map[string]interface{}{
							"name":  "anotherHeader",
							"value": "my-value2",
						},
					},
				},
				sub.ConsumerConfiguration.EntrypointConfiguration.Object,
			)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		By("updating consumer configuration and expecting the change in automation API")
		fixtures.Subscription.Spec.ConsumerConfiguration = &submodel.ConsumerConfiguration{
			EntrypointID: "webhook",
			Channel:      "payments",
			EntrypointConfiguration: utils.ToGenericStringMap(map[string]interface{}{
				"callbackUrl": "https://example.com/webhook/v2",
				"headers": []interface{}{
					map[string]interface{}{
						"name":  "demoHeader",
						"value": "my-updated-value",
					},
					map[string]interface{}{
						"name":  "anotherHeader",
						"value": "my-updated-value2",
					},
				},
			}),
		}

		Expect(manager.UpdateSafely(ctx, fixtures.Subscription)).To(Succeed())

		Eventually(func() error {
			apiHrid := refs.NewNamespacedNameFromObject(fixtures.APIv4).HRID()
			subHrid := refs.NewNamespacedNameFromObject(fixtures.Subscription).HRID()
			sub, err := client.Subscriptions.GetByHRID(apiHrid, subHrid)
			if err != nil {
				return err
			}

			if err := assert.Equals("consumerConfiguration.channel", "payments", sub.ConsumerConfiguration.Channel); err != nil {
				return err
			}
			return assert.Equals(
				"consumerConfiguration.entrypointConfiguration",
				map[string]interface{}{
					"callbackUrl": "https://example.com/webhook/v2",
					"headers": []interface{}{
						map[string]interface{}{
							"name":  "demoHeader",
							"value": "my-updated-value",
						},
						map[string]interface{}{
							"name":  "anotherHeader",
							"value": "my-updated-value2",
						},
					},
				},
				sub.ConsumerConfiguration.EntrypointConfiguration.Object,
			)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)
	})
})
