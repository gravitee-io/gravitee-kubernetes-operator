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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
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
	fixtures.Subscription.Namespace = constants.Namespace
	fixtures.APIv4.Spec.State = base.StateStopped
	fixtures.Apply()

	It("should fail if API is not started", func() {
		Eventually(func() error {
			Expect(admissionCtrl.Default(ctx, fixtures.Subscription)).ToNot(HaveOccurred())
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.Subscription)
			return assert.Equals(
				"error",
				errors.NewSeveref(
					"unable to subscribe to API [%s] because it is not started",
					fixtures.APIv4.GetRef(),
				),
				err,
			)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})
})
