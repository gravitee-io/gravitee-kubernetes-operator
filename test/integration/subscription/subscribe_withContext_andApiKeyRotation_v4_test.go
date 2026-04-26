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

package subscription

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/subscription"
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

var _ = Describe("Subscription API key rotation with v4 API", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should rotate API keys by adding a new key and revoking the old one", func() {
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
		fixtures.Subscription.Spec.ApiKeys = []subscription.ApiKeySpec{
			{Key: "rotation-key-v1-at-least-32-chars"},
		}

		By("creating API, application and subscription with initial API key")
		fixtures.Apply()

		By("expecting subscription status to be completed")
		Eventually(func() error {
			return assert.SubscriptionCompleted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		Eventually(func() error {
			return assert.SubscriptionAccepted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		client := apim.NewClient(ctx)

		By("verifying initial API key is created in APIM")
		Eventually(func() error {
			apiKeys, err := client.Subscriptions.GetApiKeys(fixtures.APIv4.Status.ID, fixtures.Subscription.Status.ID)
			if err != nil {
				return err
			}
			for _, k := range apiKeys {
				if k.Key == "rotation-key-v1-at-least-32-chars" && !k.Revoked {
					return nil
				}
			}
			return fmt.Errorf("expected active rotation-key-v1-at-least-32-chars, got %v", apiKeys)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		By("rotating to a new key by updating the subscription spec")
		fixtures.Subscription.Spec.ApiKeys = []subscription.ApiKeySpec{
			{Key: "rotation-key-v2-at-least-32-chars"},
		}
		Expect(manager.UpdateSafely(ctx, fixtures.Subscription)).To(Succeed())

		By("verifying old key is revoked and new key is active")
		Eventually(func() error {
			apiKeys, err := client.Subscriptions.GetApiKeys(fixtures.APIv4.Status.ID, fixtures.Subscription.Status.ID)
			if err != nil {
				return err
			}

			var foundNewActive, foundOldRevoked bool
			for _, k := range apiKeys {
				if k.Key == "rotation-key-v2-at-least-32-chars" && !k.Revoked {
					foundNewActive = true
				}
				if k.Key == "rotation-key-v1-at-least-32-chars" && k.Revoked {
					foundOldRevoked = true
				}
			}

			if !foundNewActive {
				return fmt.Errorf("expected active rotation-key-v2-at-least-32-chars but not found in %v", apiKeys)
			}
			if !foundOldRevoked {
				return fmt.Errorf("expected revoked rotation-key-v1-at-least-32-chars but not found in %v", apiKeys)
			}
			return nil
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)
	})

	It("should support multiple active keys for gradual rotation", func() {
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
		fixtures.Subscription.Spec.ApiKeys = []subscription.ApiKeySpec{
			{Key: "gradual-rotation-key-v1-32-chars!"},
		}

		By("creating subscription with initial key")
		fixtures.Apply()

		Eventually(func() error {
			return assert.SubscriptionCompleted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		Eventually(func() error {
			return assert.SubscriptionAccepted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		By("adding a second key alongside the first for gradual rotation")
		fixtures.Subscription.Spec.ApiKeys = []subscription.ApiKeySpec{
			{Key: "gradual-rotation-key-v1-32-chars!"},
			{Key: "gradual-rotation-key-v2-32-chars!"},
		}
		Expect(manager.UpdateSafely(ctx, fixtures.Subscription)).To(Succeed())

		client := apim.NewClient(ctx)

		By("verifying both keys are active in APIM")
		Eventually(func() error {
			apiKeys, err := client.Subscriptions.GetApiKeys(fixtures.APIv4.Status.ID, fixtures.Subscription.Status.ID)
			if err != nil {
				return err
			}

			activeKeys := make(map[string]bool)
			for _, k := range apiKeys {
				if !k.Revoked {
					activeKeys[k.Key] = true
				}
			}

			if !activeKeys["gradual-rotation-key-v1-32-chars!"] {
				return fmt.Errorf("expected active gradual-rotation-key-v1-32-chars! in %v", apiKeys)
			}
			if !activeKeys["gradual-rotation-key-v2-32-chars!"] {
				return fmt.Errorf("expected active gradual-rotation-key-v2-32-chars! in %v", apiKeys)
			}
			return nil
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		By("completing rotation by removing the old key")
		fixtures.Subscription.Spec.ApiKeys = []subscription.ApiKeySpec{
			{Key: "gradual-rotation-key-v2-32-chars!"},
		}
		Expect(manager.UpdateSafely(ctx, fixtures.Subscription)).To(Succeed())

		By("verifying only the new key is active and old key is revoked")
		Eventually(func() error {
			apiKeys, err := client.Subscriptions.GetApiKeys(fixtures.APIv4.Status.ID, fixtures.Subscription.Status.ID)
			if err != nil {
				return err
			}

			var v1Revoked, v2Active bool
			for _, k := range apiKeys {
				if k.Key == "gradual-rotation-key-v1-32-chars!" && k.Revoked {
					v1Revoked = true
				}
				if k.Key == "gradual-rotation-key-v2-32-chars!" && !k.Revoked {
					v2Active = true
				}
			}
			if !v1Revoked {
				return fmt.Errorf("expected revoked gradual-rotation-key-v1-32-chars! in %v", apiKeys)
			}
			if !v2Active {
				return fmt.Errorf("expected active gradual-rotation-key-v2-32-chars! in %v", apiKeys)
			}
			return nil
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)
	})
})
