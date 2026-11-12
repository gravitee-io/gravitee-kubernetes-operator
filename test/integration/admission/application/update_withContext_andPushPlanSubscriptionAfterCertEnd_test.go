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

package application

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate update", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := application.AdmissionCtrl{}

	// The certificate end-date check reads the security type of every subscribed plan.
	// A PUSH plan has none, which used to panic the webhook (GKO-3237).
	It("should not fail when a PUSH plan subscription ends after all client certificates", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4MessageWebhookPushPlanFile).
			WithApplication(constants.ApplicationWithClientCertsRefs).
			WithSubscription(constants.SubscriptionPushPlanFile).
			WithContext(constants.ContextWithCredentialsFile).
			AddSecret(constants.ClientCertSecretFile).
			AddConfigMap(constants.ClientCertConfigMapFile).
			Build()

		certsEndAt := "2030-01-01T00:00:00Z"
		for i := range fixtures.Application.Spec.Settings.TLS.ClientCertificates {
			fixtures.Application.Spec.Settings.TLS.ClientCertificates[i].EndsAt = certsEndAt
		}

		subEndsAt := "2035-01-01T00:00:00Z"
		fixtures.Subscription.Spec.App.Name = fixtures.Application.Name
		fixtures.Subscription.Spec.EndingAt = &subEndsAt

		fixtures.Apply()

		By("waiting for the subscription to be reconciled, so the application search finds it")

		Eventually(func() error {
			return assert.SubscriptionCompleted(fixtures.Subscription)
		}, constants.EventualTimeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		By("checking that application update validation succeeds")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Application, fixtures.Application.DeepCopy())
			return assert.Equals("error", nil, err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
