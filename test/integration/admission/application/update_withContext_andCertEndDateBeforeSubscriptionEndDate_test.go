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

	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"

	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

var _ = Describe("Validate update", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := adm.AdmissionCtrl{}

	It("should reject update when all cert end dates are before subscription ending date", func() {
		endingAt := "2028-01-01T00:00:00Z"
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.SubscribeMTLSUseCaseAPIFile).
			WithApplication(constants.ApplicationWithClientCertsAndDates).
			WithSubscription(constants.SubscriptionMTLSWithCertDatesFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		fixtures.Subscription.Spec.EndingAt = &endingAt
		fixtures.Subscription.Namespace = constants.Namespace

		fixtures.Apply()

		Eventually(func() error {
			Expect(manager.GetLatest(ctx, fixtures.Application)).ToNot(HaveOccurred())
			oldApp := fixtures.Application.DeepCopy()
			_, err := admissionCtrl.ValidateUpdate(ctx, oldApp, fixtures.Application)
			return assert.Equals(
				"error",
				errors.NewSeveref(
					"subscription [%s/%s] ending date [%s] is after all client certificate "+
						"end dates in application [%s]",
					fixtures.Subscription.Namespace,
					fixtures.Subscription.Name,
					endingAt,
					fixtures.Application.GetRef(),
				),
				err,
			)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})

	It("should allow update when at least one cert has no end date", func() {
		endingAt := "2028-01-01T00:00:00Z"
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.SubscribeMTLSUseCaseAPIFile).
			WithApplication(constants.ApplicationWithClientCertsAndDates).
			WithSubscription(constants.SubscriptionMTLSWithCertDatesFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		fixtures.Subscription.Spec.EndingAt = &endingAt
		fixtures.Subscription.Namespace = constants.Namespace
		// Clear endsAt on the certificate so it has no expiry constraint
		fixtures.Application.Spec.Settings.TLS.ClientCertificates[0].EndsAt = ""

		fixtures.Apply()

		Eventually(func() error {
			Expect(manager.GetLatest(ctx, fixtures.Application)).ToNot(HaveOccurred())
			oldApp := fixtures.Application.DeepCopy()
			_, err := admissionCtrl.ValidateUpdate(ctx, oldApp, fixtures.Application)
			return assert.Equals("error", nil, err)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})
})
