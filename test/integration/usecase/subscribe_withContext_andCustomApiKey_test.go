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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Usecase", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should subscribe with custom API key loaded from a secret", func() {
		const customAPIKey = "my-custom-api-key"

		fixtures := fixture.Builder().
			AddSecret("usecase/subscribe-to-api-key-plan/resources/custom-api-key-secret.yml").
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
		fixtures.Subscription.Spec.CustomApiKey = customAPIKey

		By("creating secret, API, application and subscription resources")
		fixtures.Apply()

		By("expecting subscription status to be completed")
		Eventually(func() error {
			return assert.SubscriptionCompleted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		Eventually(func() error {
			return assert.SubscriptionAccepted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		By("calling management API, expecting subscription API key to match customApiKey")
		client := apim.NewClient(ctx)
		Eventually(func() error {
			apiKeys, err := client.Subscriptions.GetApiKeys(fixtures.APIv4.Status.ID, fixtures.Subscription.Status.ID)
			if err != nil {
				return err
			}
			if len(apiKeys) == 0 {
				return fmt.Errorf("no API key returned for subscription")
			}
			return assert.Equals("api key", customAPIKey, apiKeys[0].Key)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)
	})
})
