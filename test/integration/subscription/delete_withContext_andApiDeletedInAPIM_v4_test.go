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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
)

var _ = Describe("Subscription deletion with v4 API already deleted in APIM", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should release the finalizer and decrement subscription counts", func() {
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

		By("deleting the API in APIM, behind the operator's back")

		Expect(manager.GetLatest(ctx, fixtures.APIv4)).To(Succeed())
		Expect(apim.NewClient(ctx).APIs.DeleteV4(fixtures.APIv4)).To(Succeed())

		By("deleting the subscription, expecting it to be gone")

		Expect(manager.Delete(ctx, fixtures.Subscription)).To(Succeed())

		Eventually(func() error {
			return assert.Deleted(ctx, "subscription", fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		By("expecting API and application subscription counts to be decremented")

		Eventually(func() error {
			if err := manager.GetLatest(ctx, fixtures.APIv4); err != nil {
				return err
			}
			return assert.Equals("API subscriptionCount", uint(0), fixtures.APIv4.Status.SubscriptionCount)
		}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)

		Eventually(func() error {
			if err := manager.GetLatest(ctx, fixtures.Application); err != nil {
				return err
			}
			return assert.Equals("application subscriptionCount", uint(0), fixtures.Application.Status.SubscriptionCount)
		}, timeout, interval).Should(Succeed(), fixtures.Application.Name)
	})
})
