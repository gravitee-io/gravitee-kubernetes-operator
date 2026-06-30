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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/subscription"
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

var _ = Describe("Validate update", labels.WithContext, func() {
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
	fixtures.Subscription.Namespace = constants.Namespace

	fixtures.Apply()

	It("should fail if the ending date is changed on the remote", func() {

		remoteEnding := "2070-11-28T12:20:00Z"
		localEnding := *fixtures.Subscription.Spec.EndingAt

		newObj := fixtures.Subscription.DeepCopy()
		By("updating the subscription with a future endingAt date")
		apim := apim.NewClient(ctx)
		_, err := apim.Subscriptions.Import(model.SubscriptionDTO{
			ID:       refs.NewNamespacedNameFromObject(fixtures.Subscription).HRID(),
			ApiID:    refs.NewNamespacedNameFromObject(fixtures.APIv4).HRID(),
			AppID:    refs.NewNamespacedNameFromObject(fixtures.Application).HRID(),
			PlanID:   fixtures.Subscription.Spec.GetPlan(),
			EndingAt: remoteEnding,
		}, newObj, false, false, false)
		Expect(err).ToNot(HaveOccurred())

		By("validating the subscription update")
		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Subscription, newObj)
			return assert.DriftDetected("endingAt: \""+localEnding+"\" != \""+remoteEnding+"\"", err)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})
})
