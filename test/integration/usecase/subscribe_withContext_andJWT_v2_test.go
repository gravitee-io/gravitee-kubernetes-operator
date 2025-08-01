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

package usecase

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	tHTTP "github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/jwt"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
)

var _ = Describe("Usecase", labels.WithContext, func() {
	timeout := constants.EventualTimeout

	interval := constants.Interval

	ctx := context.Background()

	It("should subscribe to v2 API with JWT plan", func() {
		fixtures := fixture.Builder().
			WithApplication(constants.ApplicationWithClientIDFile).
			WithAPI(constants.ApiWithJWTPlan).
			WithContext(constants.ContextWithCredentialsFile).
			WithSubscription(constants.SubscriptionFile).
			Build()

		clientID := random.GetName()
		fixtures.Application.Spec.Settings.App.ClientID = &clientID
		fixtures.Subscription.Spec.API.Name = fixtures.API.Name
		fixtures.Subscription.Spec.API.Kind = core.CRDApiDefinitionResource
		fixtures.Subscription.Spec.App.Name = fixtures.Application.Name

		By("expecting subscription status to be completed")

		fixtures.Apply()

		Eventually(func() error {
			return assert.SubscriptionCompleted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		Eventually(func() error {
			return assert.SubscriptionAccepted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		By("calling API endpoint without token, expecting status 401")

		httpClient := tHTTP.NewClient()

		endpoint := constants.BuildAPIEndpoint(fixtures.API)
		Eventually(func() error {
			res, err := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(err, res, http.StatusUnauthorized)
		}, timeout, interval).Should(Succeed(), fixtures.API.Name)

		By("calling API endpoint with a token, expecting status 200")

		token, err := jwt.GetToken(clientID, constants.SubscribeJWTUseCasePrivateKeyFile)
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() error {
			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
			if err != nil {
				return err
			}
			bearer := "Bearer " + token
			req.Header.Set("Authorization", bearer)
			res, err := httpClient.Do(req)
			return assert.NoErrorAndHTTPStatus(err, res, http.StatusOK)
		}, timeout, interval).Should(Succeed(), fixtures.API.Name)

		By("deleting subscription, expecting status 401")

		Expect(manager.Delete(ctx, fixtures.Subscription)).To(Succeed())

		Eventually(func() error {
			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
			if err != nil {
				return err
			}
			bearer := "Bearer " + token
			req.Header.Set("Authorization", bearer)
			res, err := httpClient.Do(req)
			return assert.NoErrorAndHTTPStatus(err, res, http.StatusUnauthorized)
		}, timeout, interval).Should(Succeed(), fixtures.API.Name)
	})
})
