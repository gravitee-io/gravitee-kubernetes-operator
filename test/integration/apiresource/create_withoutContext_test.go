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

package apiresource

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
)

var _ = Describe("Create", labels.WithoutContext, func() {
	httpClient := http.Client{Timeout: 5 * time.Second}
	timeout := constants.EventualTimeout
	interval := constants.Interval

	DescribeTable("without a management context",
		func(builder *fixture.FSBuilder, status int) {
			fixtures := builder.Build().Apply()

			By("expecting API status to be completed")

			Expect(assert.ApiCompleted(fixtures.API)).To(Succeed())

			By(fmt.Sprintf("calling gateway endpoint, expecting status %d", status))

			endpoint := constants.BuildAPIEndpoint(fixtures.API)
			Eventually(func() error {
				res, callErr := httpClient.Get(endpoint)
				return assert.NoErrorAndHTTPStatus(callErr, res, status)
			}, timeout, interval).Should(Succeed())

			By("expecting API event to have been emitted")

			assert.EventsEmitted(fixtures.API, "UpdateStarted", "UpdateSucceeded")
		},
		Entry("should import with redis cache resource ref",
			fixture.Builder().
				WithAPI(constants.ApiWithCacheRedisResourceRefFile).
				WithResource(constants.ApiResourceCacheRedisFile),
			200,
		),
		Entry("should import with oauth2 generic resource ref",
			fixture.Builder().
				WithAPI(constants.ApiWithOAuth2GenericResourceRefFile).
				WithResource(constants.ApiResourceOauth2GenericFile),
			200,
		),
		Entry("should import with ldap auth provider resource ref",
			fixture.Builder().
				WithAPI(constants.ApiWithLDAPAuthProviderRefFile).
				WithResource(constants.ApiResourceLDAPAuthProviderFile),
			401,
		),
		Entry("should import with inline auth provider resource ref",
			fixture.Builder().
				WithAPI(constants.ApiWithInlineAuthProviderRefFile).
				WithResource(constants.ApiResourceInlineAuthProviderFile),
			401,
		),
		Entry("should import with http auth provider resource ref",
			fixture.Builder().
				WithAPI(constants.ApiWithHTTPAuthProviderRefFile).
				WithResource(constants.ApiResourceHTTPAuthProviderFile),
			401,
		),
	)
})
