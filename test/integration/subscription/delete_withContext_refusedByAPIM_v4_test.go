// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
)

var _ = Describe("Subscription deletion refused by APIM", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should keep the subscription with a control plane error, then delete it once APIM accepts", func() {
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

		Eventually(func() error {
			return assert.SubscriptionCompleted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		By("giving the context a password APIM refuses")

		Expect(manager.GetLatest(ctx, fixtures.Context)).To(Succeed())
		fixtures.Context.Spec.SetCredentials("admin", "refused")
		Expect(manager.UpdateSafely(ctx, fixtures.Context)).To(Succeed())

		By("deleting the subscription, expecting it kept with APIM's refusal on Accepted")

		Expect(manager.Delete(ctx, fixtures.Subscription)).To(Succeed())

		Eventually(func() error {
			if err := manager.GetLatest(ctx, fixtures.Subscription); err != nil {
				return err
			}
			accepted := k8s.GetCondition(fixtures.Subscription, k8s.ConditionAccepted)
			if accepted == nil || accepted.Status != metav1.ConditionFalse {
				return fmt.Errorf("expected Accepted=False, got %v", accepted)
			}
			if err := assert.Equals("Accepted reason", "ControlPlaneError", accepted.Reason); err != nil {
				return err
			}
			if !strings.Contains(accepted.Message, "status 401") {
				return fmt.Errorf("expected APIM's 401 in the message, got %q", accepted.Message)
			}
			return nil
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		Expect(assert.HasFinalizer(fixtures.Subscription, core.SubscriptionFinalizer)).To(Succeed())

		By("restoring the password, expecting the retried delete to remove the subscription")

		Expect(manager.GetLatest(ctx, fixtures.Context)).To(Succeed())
		fixtures.Context.Spec.SetCredentials("admin", "admin")
		Expect(manager.UpdateSafely(ctx, fixtures.Context)).To(Succeed())

		Eventually(func() error {
			return assert.Deleted(ctx, "subscription", fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)
	})
})
