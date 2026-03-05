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

	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"

	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

var _ = FDescribe("Validate create", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := adm.AdmissionCtrl{}

	It("should fail to create if subscription ending date is after all client certificate end dates", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.SubscribeMTLSUseCaseAPIFile).
			WithApplication(constants.ApplicationWithClientCertsAndDates).
			WithSubscription(constants.SubscriptionMTLSWithCertDatesFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		// The application fixture has clientCertificates with endsAt: 2027-06-01T00:00:00Z
		// Set subscription endingAt to after all cert endsAt dates
		endingAt := "2028-01-01T00:00:00Z"
		fixtures.Subscription.Spec.EndingAt = &endingAt
		fixtures.Subscription.Namespace = constants.Namespace

		fixtures.Apply()

		Eventually(func() error {
			Expect(admissionCtrl.Default(ctx, fixtures.Subscription)).ToNot(HaveOccurred())
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.Subscription)
			return assert.Equals(
				"error",
				errors.NewSeveref(
					"subscription ending date [%s] is after all client certificate end dates in application [%s]",
					endingAt,
					fixtures.Application.GetRef(),
				),
				err,
			)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})

	It("should fail to update if subscription ending date is after all client certificate end dates", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.SubscribeMTLSUseCaseAPIFile).
			WithApplication(constants.ApplicationWithClientCertsAndDates).
			WithSubscription(constants.SubscriptionMTLSWithCertDatesFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build()
		fixtures.Apply()

		Eventually(func() error {
			Expect(admissionCtrl.Default(ctx, fixtures.Subscription)).ToNot(HaveOccurred())
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.Subscription)
			return assert.Equals("error", nil, err)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())

		// The application fixture has clientCertificates with endsAt: 2027-06-01T00:00:00Z
		// Set subscription endingAt to after all cert endsAt dates
		endingAt := "2028-01-01T00:00:00Z"
		newSubscription := fixtures.Subscription.DeepCopy()
		newSubscription.Spec.EndingAt = &endingAt

		Eventually(func() error {
			Expect(admissionCtrl.Default(ctx, fixtures.Subscription)).ToNot(HaveOccurred())
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Subscription, newSubscription)
			return assert.Equals(
				"error",
				errors.NewSeveref(
					"subscription ending date [%s] is after all client certificate end dates in application [%s]",
					endingAt,
					fixtures.Application.GetRef(),
				),
				err,
			)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())

	})
})
