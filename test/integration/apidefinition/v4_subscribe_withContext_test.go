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

package apidefinition_test

import (
	"context"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
)

var _ = Describe("Subscribe", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	httpClient := http.Client{Timeout: 5 * time.Second}

	It("should subscribe to API V4 key plan", func() {
		fixtures := fixture.Builder().
			WithAPIv4(constants.ApiV4WithApiKeyPlanFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting API V4 status to be completed")

		Expect(assert.ApiV4Completed(fixtures.APIv4)).To(Succeed())

		By("calling API endpoint, expecting status 401")

		endpoint := constants.BuildAPIV4Endpoint(fixtures.APIv4.Spec.Listeners[0])
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

		By("calling rest API expecting to find API V4 with plan")

		Expect(fixtures.APIv4.Status.Plans).ToNot(BeEmpty())
		planID := fixtures.APIv4.Status.Plans["API_KEY"]

		By("calling rest API expecting to application to subscribe to API")

		subscription, err := apim.Subscriptions.Subscribe(fixtures.APIv4.Status.ID, app.Id, planID)
		Expect(err).ToNot(HaveOccurred())

		By("calling rest API expecting to find subscription API key")

		keys, err := apim.Subscriptions.GetApiKeys(fixtures.APIv4.Status.ID, subscription.Id)
		Expect(err).ToNot(HaveOccurred())
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
