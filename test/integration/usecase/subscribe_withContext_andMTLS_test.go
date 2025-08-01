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

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	tHTTP "github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Usecase", labels.WithContext, func() {
	timeout := constants.EventualTimeout

	interval := constants.Interval

	ctx := context.Background()

	It("should subscribe to v4 API with MTLS plan", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.SubscribeMTLSUseCaseTLSSecretFile).
			WithApplication(constants.SubscribeMTLSUseCaseApplicationFile).
			WithAPIv4(constants.SubscribeMTLSUseCaseAPIFile).
			WithContext(constants.SubscribeMTLSUseCaseContextFile).
			WithSubscription(constants.SubscribeMTLSUseCaseSubscriptionFile).
			Build().
			Apply()

		Eventually(func() error {
			return assert.SubscriptionCompleted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		Eventually(func() error {
			return assert.SubscriptionAccepted(fixtures.Subscription)
		}, timeout, interval).Should(Succeed(), fixtures.Subscription.Name)

		By("calling API endpoint without client auth, expecting status 401")

		endpoint := constants.BuildAPIV4EndpointForTLS(fixtures.APIv4.Spec.Listeners[0])

		httpClient := tHTTP.NewClient()

		Eventually(func() error {
			res, err := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(err, res, http.StatusUnauthorized)
		}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)

		By("calling API endpoint with client auth, expecting status 200")

		mtlsClient := tHTTP.NewMTLSClient(
			constants.SubscribeMTLSUseCaseRootCAFile,
			constants.SubscribeMTLSUseCaseClientCertFile,
			constants.SubscribeMTLSUseCaseClientKeyFile,
		)

		Eventually(func() error {
			res, err := mtlsClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(err, res, http.StatusOK)
		}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)

		By("deleting subscription, expecting error")

		Expect(manager.Delete(ctx, fixtures.Subscription)).To(Succeed())

		Eventually(func() error {
			_, err := mtlsClient.Get(endpoint)
			return assert.NotNil("error", err)
		}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)
	})
})
