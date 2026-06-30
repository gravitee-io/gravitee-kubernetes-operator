package subscription

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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

var _ = Describe("Validate drift", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := admission.AdmissionCtrl{}
	clientID := random.GetName()
	fixtures := fixture.
		Builder().
		WithAPIv4(constants.ApiV4WithJWTPlanFile).
		WithApplication(constants.ApplicationWithClientIDFile).
		WithSubscription(constants.SubscriptionFile).
		WithContext(constants.ContextWithCredentialsFile).
		Build()
	fixtures.Application.Spec.Settings.App.ClientID = &clientID
	fixtures.Subscription.Namespace = constants.Namespace
	fixtures.Subscription.Spec.API.Namespace = constants.Namespace
	fixtures.Subscription.Spec.App.Namespace = constants.Namespace
	fixtures.Apply()

	It("should not drift on a simple update", func() {

		By("changing the subscription endingAt")
		newObj := fixtures.Subscription.DeepCopy()
		updatedEnding := "2070-01-01T00:00:00Z"
		newObj.Spec.EndingAt = &updatedEnding

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Subscription, newObj)
			return err
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should detect drift", func() {
		By("changing the remote subscription endingAt")

		newObj := fixtures.Subscription.DeepCopy()

		apim := apim.NewClient(ctx)

		_, err := apim.Subscriptions.Import(model.SubscriptionDTO{
			ID:       refs.NewNamespacedNameFromObject(fixtures.Subscription).HRID(),
			ApiID:    refs.NewNamespacedNameFromObject(fixtures.APIv4).HRID(),
			AppID:    refs.NewNamespacedNameFromObject(fixtures.Application).HRID(),
			PlanID:   fixtures.Subscription.Spec.GetPlan(),
			EndingAt: "2070-11-28T12:20:00Z",
		}, newObj, false, false, false)
		Expect(err).ToNot(HaveOccurred())

		By("changing the CRD subscription endingAt")

		localEnding := "2060-12-25T09:12:28Z"
		newObj.Spec.EndingAt = &localEnding

		By("checking that subscription validation returns error")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Subscription, newObj)
			return assert.DriftDetected("endingAt: \"2060-12-25T09:12:28Z\" != \"2070-11-28T12:20:00Z\"", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
