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
package v4

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
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Create", labels.WithContext, func() {
	httpClient := tHTTP.NewClient()

	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	DescribeTable("with wrong management context",
		func(builder *fixture.FSBuilder) {
			Skip(`
				This test was migrated and moved to e2e test suite
			`)

			fixtures := builder.Build().Apply()

			By("expecting API V4 status to be failed")

			Expect(assert.ApiV4Failed(fixtures.APIv4)).To(Succeed())

			By("expecting API V4 events to have been emitted")

			assert.EventsEmitted(fixtures.APIv4, "UpdateStarted", "UpdateFailed")

			endpoint := constants.BuildAPIV4Endpoint(fixtures.APIv4.Spec.Listeners[0])
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

			By("expecting API V4 status to be completed")

			Eventually(func() error {
				err := manager.GetLatest(ctx, fixtures.APIv4)
				if err != nil {
					return err
				}
				return assert.ApiV4Completed(fixtures.APIv4)
			}, timeout, interval).Should(Succeed())

			By("calling gateway endpoint, expecting status 200")

			Eventually(func() error {
				res, callErr := httpClient.Get(endpoint)
				return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).Should(Succeed())

			By("calling rest API, expecting API V4 to match status cross ID")

			apim := apim.NewClient(ctx)

			Eventually(func() error {
				api, apiErr := apim.APIs.GetV4ByID(fixtures.APIv4.Status.ID)
				if apiErr != nil {
					return apiErr
				}
				return assert.Equals("API V4 crossId", fixtures.APIv4.Status.CrossID, api.CrossID)
			}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)

			By("expecting API V4 event to have been emitted")

			assert.EventsEmitted(fixtures.APIv4, "UpdateStarted", "UpdateSucceeded")
		},
		Entry("should reconcile on bad credentials updates",
			fixture.Builder().
				WithAPIv4(constants.ApiV4WithContextFile).
				WithContext(constants.ContextWithBadCredentialsFile),
		),
		Entry("should reconcile on bad URL update",
			fixture.Builder().
				WithAPIv4(constants.ApiV4WithContextFile).
				WithContext(constants.ContextWithBadURLFile),
		))

})
