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

package v4

import (
	"context"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"

	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

var _ = Describe("Validate delete", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := adm.AdmissionCtrl{}

	fixtures := fixture.Builder().
		AddSecret(constants.SubscribeJWTUseCasePublicKeySecretFile).
		WithApplication(constants.SubscribeJWTUseCaseApplicationFile).
		WithAPIv4(constants.SubscribeJWTUseCaseAPIFile).
		WithContext(constants.SubscribeJWTUseCaseContextFile).
		WithSubscription(constants.SubscribeJWTUseCaseSubscriptionFile).
		Build()

	clientID := random.GetName()
	fixtures.Application.Spec.Settings.App.ClientID = &clientID
	fixtures.Subscription.Namespace = constants.Namespace

	fixtures.Apply()

	It("should fail with subscription", func() {

		Eventually(func() error {
			Expect(admissionCtrl.Default(ctx, fixtures.Subscription)).ToNot(HaveOccurred())

			Expect(manager.GetLatest(ctx, fixtures.APIv4)).ToNot(HaveOccurred())

			newApi := fixtures.APIv4.DeepCopy()
			newPlans := *newApi.Spec.Plans
			delete(newPlans, fixtures.Subscription.Spec.Plan)
			newPlans["KEY_LESS"] = &v4.Plan{
				Name: "KEY_LESS",
				Security: &v4.PlanSecurity{
					Type: "KEY_LESS",
				},
			}

			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.APIv4, newApi)
			return assert.Equals(
				"error",
				errors.NewSeveref(
					"Plan [%s] could not be found in API [%s] "+
						"but there is a subscription referencing it. "+
						"You can review the depending subscriptions using the following command: "+
						"kubectl get subscriptions.gravitee.io -A "+
						"-o jsonpath='{.items[?(@.spec.api.name==\"%s\")].metadata.name}'",
					fixtures.Subscription.Spec.Plan, newApi.GetRef(), newApi.GetName(),
				),
				err,
			)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})
})
