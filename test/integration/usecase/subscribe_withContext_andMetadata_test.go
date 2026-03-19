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

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
)

var _ = Describe("Usecase", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should subscribe with metadata", func() {
		fixtures := fixture.Builder().
			WithApplication(constants.ApplicationWithClientIDFile).
			WithAPI(constants.ApiWithJWTPlan).
			WithContext(constants.ContextWithCredentialsFile).
			WithSubscription(constants.SubscriptionFile).
			Build()

		clientID := random.GetName()
		fixtures.Application.Spec.Settings.App.ClientID = &clientID
		fixtures.Subscription.Spec.API.Name = fixtures.API.Name
		fixtures.Subscription.Spec.API.Kind = core.CRDApiDefinitionResource
		fixtures.Subscription.Spec.App.Name = fixtures.Application.Name
		fixtures.Subscription.Spec.Metadata = map[string]string{
			"team":        "platform",
			"cost-center": "engineering",
		}

		fixtures.Apply()

		By("expecting subscription status to be completed")

		Eventually(func() error {
			return assert.SubscriptionCompleted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		Eventually(func() error {
			return assert.SubscriptionAccepted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		By("calling management API, expecting metadata to be saved")

		client := apim.NewClient(ctx)

		Eventually(func() error {
			sub, err := client.Subscriptions.GetByID(
				fixtures.API.Status.ID,
				fixtures.Subscription.Status.ID,
			)
			if err != nil {
				return err
			}
			expectedMetadata := map[string]string{
				"team":        "platform",
				"cost-center": "engineering",
			}
			return assert.Equals("metadata", expectedMetadata, sub.Metadata)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		By("updating metadata and expecting it to be updated in management API")

		fixtures.Subscription.Spec.Metadata = map[string]string{
			"team":        "platform-v2",
			"cost-center": "engineering",
			"env":         "staging",
		}

		Expect(manager.UpdateSafely(ctx, fixtures.Subscription)).To(Succeed())

		Eventually(func() error {
			sub, err := client.Subscriptions.GetByID(
				fixtures.API.Status.ID,
				fixtures.Subscription.Status.ID,
			)
			if err != nil {
				return err
			}
			expectedMetadata := map[string]string{
				"team":        "platform-v2",
				"cost-center": "engineering",
				"env":         "staging",
			}
			return assert.Equals("metadata", expectedMetadata, sub.Metadata)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)
	})
})
