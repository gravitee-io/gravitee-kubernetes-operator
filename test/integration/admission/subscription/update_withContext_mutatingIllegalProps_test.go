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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
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
	fixtures.Subscription.Namespace = constants.Namespace
	fixtures.Subscription.Spec.API.Namespace = constants.Namespace
	fixtures.Subscription.Spec.App.Namespace = constants.Namespace
	fixtures.Application.Spec.Settings.App.ClientID = &clientID

	fixtures.Apply()

	DescribeTable("mutating illegal property",
		func(modify func(v1alpha1.Subscription) *v1alpha1.Subscription, expected error) {
			Eventually(func() error {
				Expect(admissionCtrl.Default(ctx, fixtures.Subscription)).ToNot(HaveOccurred())
				newSub := modify(*fixtures.Subscription)
				_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Subscription, newSub)
				return assert.Equals("error", expected, err)
			}, constants.EventualTimeout, constants.Interval).Should(Succeed())

		},
		Entry(
			"API ref",
			func(sub v1alpha1.Subscription) *v1alpha1.Subscription {
				cp := sub.DeepCopy()
				cp.Spec.API.Name = "illegal"
				return cp
			},
			errors.NewSeveref(
				"API reference is immutable. Detected change from [%s] to [default/illegal]",
				fixtures.APIv4.GetRef(),
			),
		),
		Entry(
			"application ref",
			func(sub v1alpha1.Subscription) *v1alpha1.Subscription {
				cp := sub.DeepCopy()
				cp.Spec.App.Name = "illegal"
				return cp
			},
			errors.NewSeveref(
				"Application reference is immutable. Detected change from [%s] to [default/illegal]",
				fixtures.Application.GetRef(),
			),
		),
		Entry(
			"plan ref",
			func(sub v1alpha1.Subscription) *v1alpha1.Subscription {
				cp := sub.DeepCopy()
				cp.Spec.Plan = "illegal"
				return cp
			},
			errors.NewSevere(
				"Plan is immutable. Detected change from [JWT] to [illegal]",
			),
		),
	)
})
