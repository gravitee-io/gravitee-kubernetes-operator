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
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"

	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

var _ = Describe("Validate create", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := adm.AdmissionCtrl{}

	fixtures := fixture.
		Builder().
		WithAPIv4(constants.ApiV4WithJWTPlanFile).
		WithApplication(constants.ApplicationWithClientIDFile).
		WithSubscription(constants.SubscriptionFile).
		WithContext(constants.ContextWithCredentialsFile).
		Build()

	clientID := random.GetName()
	fixtures.Application.Spec.Settings.App.ClientID = &clientID

	// An API is not required to reference a management context, while an application is.
	fixtures.APIv4.Spec.Context = nil

	// The subscription is submitted to the admission controller directly and never applied:
	// the reconciler cannot build an APIM client for an API without a contextRef, so the
	// resource would never reach a terminal state.
	subscription := fixtures.Subscription
	fixtures.Subscription = nil

	fixtures.Apply()

	It("should reject the subscription if the API has no context", func() {
		Eventually(func() error {
			Expect(admissionCtrl.Default(ctx, subscription)).ToNot(HaveOccurred())
			_, err := admissionCtrl.ValidateCreate(ctx, subscription)
			return assert.Equals(
				"error",
				errors.NewSeveref(
					"unable to subscribe to API [%s] because it does not reference a management context",
					fixtures.APIv4.GetRef(),
				),
				err,
			)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})
})
