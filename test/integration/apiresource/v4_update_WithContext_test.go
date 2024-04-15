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

package apiresource_test

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

var _ = Describe("Update", labels.WithContext, func() {
	httpClient := http.Client{Timeout: 5 * time.Second}
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	DescribeTable("with a management context",
		func(builder *fixture.FSBuilder) {
			fixtures := builder.Build().Apply()

			By("expecting API status to be completed")

			Expect(assert.ApiV4Completed(fixtures.APIv4)).To(Succeed())

			By("calling gateway endpoint, expecting status 401")

			endpoint := constants.BuildAPIV4Endpoint(fixtures.APIv4.Spec.Listeners[0])
			Eventually(func() error {
				res, callErr := httpClient.Get(endpoint)
				return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusUnauthorized)
			}, timeout, interval).Should(Succeed())

			By("calling rest API, expecting to find resources in API")

			apim := apim.NewClient(ctx)

			Eventually(func() error {
				api, apiErr := apim.APIs.GetV4ByID(fixtures.APIv4.Status.ID)
				if apiErr != nil {
					return apiErr
				}
				return assert.NotEmptySlice("resource", api.Resources)
			}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)

			By("disabling the resource")

			updated := fixtures.Resource.DeepCopy()
			updated.Spec.Enabled = false

			Expect(manager.UpdateSafely(updated)).To(Succeed())

			By("calling rest API, expecting resource to be updated")

			Eventually(func() error {
				api, apiErr := apim.APIs.GetV4ByID(fixtures.APIv4.Status.ID)
				if apiErr != nil {
					return apiErr
				}
				if err := assert.NotEmptySlice("resource", api.Resources); err != nil {
					return err
				}
				return assert.Equals("enabled", false, api.Resources[0].Enabled)
			}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)

			By("expecting API event to have been emitted")

			assert.EventsEmitted(fixtures.APIv4, "UpdateStarted", "UpdateSucceeded")
		},
		Entry("should import with ldap auth provider resource ref",
			fixture.Builder().
				WithAPIv4(constants.ApiV4WithLDAPAuthProviderRefFile).
				WithResource(constants.ApiResourceLDAPAuthProviderFile).
				WithContext(constants.ContextWithSecretFile),
		),
		Entry("should import with inline auth provider resource ref",
			fixture.Builder().
				WithAPIv4(constants.ApiV4WithInlineAuthProviderRef).
				WithResource(constants.ApiResourceInlineAuthProviderFile).
				WithContext(constants.ContextWithSecretFile),
		),
		Entry("should import with http auth provider resource ref",
			fixture.Builder().
				WithAPIv4(constants.ApiV4WithHTTPAuthProviderRefFile).
				WithResource(constants.ApiResourceHTTPAuthProviderFile).
				WithContext(constants.ContextWithSecretFile),
		),
	)
})
