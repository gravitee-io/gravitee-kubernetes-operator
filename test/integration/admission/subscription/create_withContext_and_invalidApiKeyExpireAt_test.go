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
	"fmt"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/subscription"
	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/subscription"
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

	invalidExpireAt := "not-a-date"

	fixtures := fixture.
		Builder().
		WithAPIv4(constants.ApiV4WithApiKeyPlanFile).
		WithApplication(constants.ApplicationWithClientIDFile).
		WithSubscription(constants.SubscriptionFile).
		WithContext(constants.ContextWithCredentialsFile).
		Build()

	clientId := random.GetName()
	fixtures.Application.Spec.Settings.App.ClientID = &clientId
	fixtures.Subscription.Spec.API.Name = fixtures.APIv4.Name
	fixtures.Subscription.Spec.App.Name = fixtures.Application.Name
	fixtures.Subscription.Namespace = constants.Namespace

	fixtures.Apply()

	It("should fail if apiKeys has an invalid expireAt format", func() {
		fixtures.Subscription.Spec.Plan = "API_KEY"
		fixtures.Subscription.Spec.ApiKeys = []subscription.ApiKeySpec{
			{Key: "my-key-value-at-least-32-chars!!!!", ExpireAt: &invalidExpireAt},
		}
		Eventually(func() error {
			Expect(admissionCtrl.Default(ctx, fixtures.Subscription)).ToNot(HaveOccurred())
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.Subscription)
			if err == nil {
				return fmt.Errorf("expected validation error but got none")
			}
			if !strings.Contains(err.Error(), "invalid expireAt for key [my-key-value-at-least-32-chars!!!!]") {
				return fmt.Errorf("expected error about invalid expireAt, got: %s", err.Error())
			}
			return nil
		}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	})
})
