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
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

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

	DescribeTable("with a management context",
		func(builder *fixture.FSBuilder, status int) {
			fixtures := builder.Build().Apply()

			By("expecting API V4 status to be completed")

			Expect(assert.ApiV4Completed(fixtures.APIv4)).To(Succeed())

			By("expecting to find config map")

			cm := &v1.ConfigMap{}
			Eventually(func() error {
				return manager.Client().Get(ctx, types.NamespacedName{
					Name:      fixtures.APIv4.Name,
					Namespace: fixtures.APIv4.Namespace,
				}, cm)
			}, timeout, interval).Should(Succeed())

			By(fmt.Sprintf("calling gateway endpoint, expecting status %d", status))

			endpoint := constants.BuildAPIV4Endpoint(fixtures.APIv4.Spec.Listeners[0])
			Eventually(func() error {
				res, callErr := httpClient.Get(endpoint)
				return assert.NoErrorAndHTTPStatus(callErr, res, status)
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
		// Entry("should import",
		//	fixture.Builder().
		//		WithAPIv4(constants.ApiV4WithContextFile).
		//		WithContext(constants.ContextWithSecretFile),
		//	200,
		// ),
		Entry("should import with health check",
			fixture.Builder().
				WithAPIv4(constants.ApiV4WithHCFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with disabled health check",
			fixture.Builder().
				WithAPIv4(constants.ApiV4WithDisabledHCFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with logging",
			fixture.Builder().
				WithAPIv4(constants.ApiV4WithLoggingFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		// Not supported yet
		// Entry("should import with metadata",
		//	fixture.Builder().
		//		WithAPIv4(constants.ApiV4WithMetadataFile).
		//		WithContext(constants.ContextWithSecretFile),
		//	200,
		// ),
		Entry("should import with cache redis resource",
			fixture.Builder().
				WithAPIv4(constants.ApiV4WithCacheRedisResourceFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with oauth2 generic resource",
			fixture.Builder().
				WithAPIv4(constants.ApiV4WithOAuth2GenericResourceFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with oauth2 am resource",
			fixture.Builder().
				WithAPIv4(constants.ApiV4WithOauth2AmResourceFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with keycloak adapter resource",
			fixture.Builder().
				WithAPIv4(constants.ApiV4WithKeycloakAdapterFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with LDAP auth provider",
			fixture.Builder().
				WithAPIv4(constants.ApiV4WithLDAPAuthProviderFile).
				WithContext(constants.ContextWithSecretFile),
			401,
		),
		Entry("should import with inline auth provider",
			fixture.Builder().
				WithAPIv4(constants.ApiV4WithInlineAuthProviderFile).
				WithContext(constants.ContextWithSecretFile),
			401,
		),
		Entry("should import with HTTP auth provider",
			fixture.Builder().
				WithAPIv4(constants.ApiV4WithHTTPAuthProviderFile).
				WithContext(constants.ContextWithSecretFile),
			401,
		),
	)
})
