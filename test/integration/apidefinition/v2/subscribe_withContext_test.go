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

package v2

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	tHTTP "github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
)

var _ = Describe("Subscribe", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	httpClient := tHTTP.NewClient()

	It("should subscribe to API key plan", func() {
		fixtures := fixture.Builder().
			WithAPI(constants.ApiWithApiKeyPlanFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting API status to be completed")

		Expect(assert.ApiCompleted(fixtures.API)).To(Succeed())
		Expect(assert.ApiAccepted(fixtures.API)).To(Succeed())

		By("calling API endpoint, expecting status 401")

		endpoint := constants.BuildAPIEndpoint(fixtures.API)
		Eventually(func() error {
			res, err := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(err, res, http.StatusUnauthorized)
		}, timeout, interval).Should(Succeed())

		By("calling rest API, expecting to find an application")

		apim := apim.NewClient(ctx)
		apps, err := apim.Applications.Search("", "ACTIVE")
		Expect(err).ToNot(HaveOccurred())
		Expect(apps).NotTo(BeEmpty())
		app := apps[0]

		By("calling rest API expecting to find API with plan")

		api, err := apim.APIs.GetByID(fixtures.API.Status.ID)
		Expect(err).ToNot(HaveOccurred())
		Expect(api.Plans).ToNot(BeEmpty())
		planID := api.Plans[0].ID

		By("calling rest API expecting to application to subscribe to API")

		subscription, err := apim.Subscriptions.Subscribe(fixtures.API.Status.ID, app.ID, planID)
		Expect(err).ToNot(HaveOccurred())

		By("calling rest API expecting to find subscription API key")

		keys, err := apim.Subscriptions.GetApiKeys(fixtures.API.Status.ID, subscription.ID)
		Expect(err).ToNot(HaveOccurred(), fixtures.API.Name)
		Expect(keys).ToNot(BeEmpty())
		key := keys[0].Key

		By("calling API endpoint with API key, expecting status 200")

		Eventually(func() error {
			req, httpErr := http.NewRequest(http.MethodGet, endpoint, nil)
			if httpErr != nil {
				return httpErr
			}
			req.Header.Set("X-Gravitee-Api-Key", key)
			res, httpErr := httpClient.Do(req)
			return assert.NoErrorAndHTTPStatus(httpErr, res, http.StatusOK)
		}, timeout, interval).Should(Succeed())
	})
})
