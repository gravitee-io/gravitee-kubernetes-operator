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
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

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

	DescribeTable("with a management context in another namespace and a secret ref",
		func(builder *fixture.FSBuilder, status int) {
			fixtures := builder.Build()

			fixtures.API.Namespace = constants.GraviteeNamespace

			fixtures = fixtures.Apply()

			By("expecting API status to be completed")

			Expect(assert.ApiCompleted(fixtures.API)).To(Succeed())

			By("expecting to find config map")

			cm := &v1.ConfigMap{}
			Eventually(func() error {
				return manager.Client().Get(ctx, types.NamespacedName{
					Name:      fixtures.API.Name,
					Namespace: fixtures.API.Namespace,
				}, cm)
			}, timeout, interval).Should(Succeed())

			By(fmt.Sprintf("calling gateway endpoint, expecting status %d", status))

			endpoint := constants.BuildAPIEndpoint(fixtures.API)
			Eventually(func() error {
				res, callErr := httpClient.Get(endpoint)
				return assert.NoErrorAndHTTPStatus(callErr, res, status)
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
		Entry("should import",
			fixture.Builder().
				WithAPI(constants.ApiWithContextFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with health check",
			fixture.Builder().
				WithAPI(constants.ApiWithHCFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with disabled health check",
			fixture.Builder().
				WithAPI(constants.ApiWithDisabledHCFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with logging",
			fixture.Builder().
				WithAPI(constants.ApiWithLoggingFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with endpoint groups",
			fixture.Builder().
				WithAPI(constants.ApiWithEndpointGroupsFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with service discovery",
			fixture.Builder().
				WithAPI(constants.ApiWithServiceDiscoveryFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with metadata",
			fixture.Builder().
				WithAPI(constants.ApiWithMetadataFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with cache redis resource",
			fixture.Builder().
				WithAPI(constants.ApiWithCacheRedisResourceFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with oauth2 generic resource",
			fixture.Builder().
				WithAPI(constants.ApiWithOAuth2GenericResourceFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with oauth2 am resource",
			fixture.Builder().
				WithAPI(constants.ApiWithOauth2AmResourceFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with keycloak adapter resource",
			fixture.Builder().
				WithAPI(constants.ApiWithKeycloakAdapterFile).
				WithContext(constants.ContextWithSecretFile),
			200,
		),
		Entry("should import with LDAP auth provider",
			fixture.Builder().
				WithAPI(constants.ApiWithLDAPAuthProviderFile).
				WithContext(constants.ContextWithSecretFile),
			401,
		),
		Entry("should import with inline auth provider",
			fixture.Builder().
				WithAPI(constants.ApiWithInlineAuthProviderFile).
				WithContext(constants.ContextWithSecretFile),
			401,
		),
		Entry("should import with HTTP auth provider",
			fixture.Builder().
				WithAPI(constants.ApiWithHTTPAuthProviderFile).
				WithContext(constants.ContextWithSecretFile),
			401,
		),
	)
})
