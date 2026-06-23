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

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
)

var _ = Describe("Subscription update with v4 API", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should subscribe with metadata", func() {
		fixtures := fixture.Builder().
			WithApplication(constants.ApplicationWithClientIDFile).
			WithAPIv4(constants.ApiV4WithJWTPlanFile).
			WithContext(constants.ContextWithCredentialsFile).
			WithSubscription(constants.SubscriptionFile).
			Build()

		clientID := random.GetName()
		fixtures.Application.Spec.Settings.App.ClientID = &clientID
		fixtures.Subscription.Spec.API.Name = fixtures.APIv4.Name
		fixtures.Subscription.Spec.API.Kind = core.CRDApiV4DefinitionResource
		fixtures.Subscription.Spec.App.Name = fixtures.Application.Name

		fixtures.Apply()

		By("expecting subscription status to be completed")

		Eventually(func() error {
			return assert.SubscriptionCompleted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		Eventually(func() error {
			return assert.SubscriptionAccepted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		By("Expecting subscription to be updated")
		endDate := "2027-09-01T00:00:00Z"
		fixtures.Subscription.Spec.EndingAt = &endDate

		Expect(manager.UpdateSafely(ctx, fixtures.Subscription)).To(Succeed())

		By("Expecting API to keep the same subscription count")
		Consistently(func() error {
			err := manager.GetLatest(ctx, fixtures.APIv4)
			if err != nil {
				return err
			}
			return assert.Equals(
				"subscriptionCount",
				1,
				fixtures.APIv4.Status.SubscriptionCount,
			)
		}, constants.ConsistentTimeout, interval).ShouldNot(Succeed())

		By("Expecting App to keep the same subscription count")
		Consistently(func() error {
			err := manager.GetLatest(ctx, fixtures.Application)
			if err != nil {
				return err
			}
			return assert.Equals(
				"subscriptionCount",
				1,
				fixtures.Application.Status.SubscriptionCount,
			)
		}, constants.ConsistentTimeout, interval).ShouldNot(Succeed())

	})
})
