// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package apidefinition

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
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Create", labels.WithContext, func() {
	httpClient := http.Client{Timeout: 5 * time.Second}

	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	DescribeTable("with wrong management context",
		func(builder *fixture.FSBuilder) {
			fixtures := builder.Build().Apply()

			By("expecting API status to be failed")

			Expect(assert.ApiFailed(fixtures.API)).To(Succeed(), fixtures.API.Name)

			By("expecting API events to have been emitted")

			assert.EventsEmitted(fixtures.API, "UpdateStarted", "UpdateFailed")

			endpoint := constants.BuildAPIEndpoint(fixtures.API)
			Eventually(func() error {
				res, callErr := httpClient.Get(endpoint)
				return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusNotFound)
			}, timeout, interval).Should(Succeed())

			By("fixing the management context")

			fixed := fixtures.Context.DeepCopy()
			fixed.Spec = fixture.Builder().
				WithContext(constants.ContextWithCredentialsFile).
				Build().Context.Spec

			Eventually(func() error {
				return manager.UpdateSafely(ctx, fixed)
			}, timeout, interval).Should(Succeed())

			By("expecting API status to be completed")

			Eventually(func() error {
				err := manager.GetLatest(ctx, fixtures.API)
				if err != nil {
					return err
				}
				return assert.ApiCompleted(fixtures.API)
			}, timeout, interval).Should(Succeed())

			By("calling gateway endpoint, expecting status 200")

			Eventually(func() error {
				res, callErr := httpClient.Get(endpoint)
				return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).Should(Succeed())

			By("calling rest API, expecting API to match status cross ID")

			apim := apim.NewClient(ctx)

			Eventually(func() error {
				api, apiErr := apim.APIs.GetByID(fixtures.API.Status.ID)
				if apiErr != nil {
					return apiErr
				}
				return assert.Equals("API entity crossId", fixtures.API.Status.CrossID, api.CrossID)
			}, timeout, interval).ShouldNot(HaveOccurred(), fixtures.API.Name)

			By("expecting API event to have been emitted")

			assert.EventsEmitted(fixtures.API, "UpdateStarted", "UpdateSucceeded")
		},
		Entry("should reconcile on bad credentials updates",
			fixture.Builder().
				WithAPI(constants.ApiWithContextFile).
				WithContext(constants.ContextWithBadCredentialsFile),
		),
		Entry("should reconcile on bad URL update",
			fixture.Builder().
				WithAPI(constants.ApiWithContextFile).
				WithContext(constants.ContextWithBadURLFile),
		))

})
