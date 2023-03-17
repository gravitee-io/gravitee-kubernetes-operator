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

package test

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	model "github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"
)

var _ = Describe("Create", func() {
	const (
		origin = "kubernetes"
		mode   = "fully_managed"
	)

	httpClient := http.Client{Timeout: 5 * time.Second}

	Context("a basic spec without a management context", func() {
		var apiDefinitionFixture *gio.ApiDefinition

		var apiLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Initializing the API definition fixture")
			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api: internal.BasicApiFile,
			})

			Expect(err).ToNot(HaveOccurred())

			apiDefinitionFixture = fixtures.Api
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
		})

		It("should create an API Definition", func() {
			By("Creating an API definition resource without a management context")

			Expect(k8sClient.Create(ctx, apiDefinitionFixture)).Should(Succeed())

			By("Getting created resource and expect to find it")

			apiDefinition := new(gio.ApiDefinition)
			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, apiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			var endpoint = internal.GatewayUrl + apiDefinition.Spec.Proxy.VirtualHosts[0].Path

			expectedApiName := apiDefinitionFixture.Spec.Name
			Expect(apiDefinition.Spec.Name).Should(Equal(expectedApiName))

			By("Calling gateway endpoint and expect the API to be available")

			Eventually(func() error {
				res, callErr := httpClient.Get(endpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})
	})

	Context("a basic spec with a management context", func() {
		var apiDefinitionFixture *gio.ApiDefinition
		var managementContextFixture *gio.ManagementContext
		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName

		BeforeEach(func() {
			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:     internal.ApiWithContextFile,
				Context: internal.ContextWithSecretFile,
			})

			Expect(err).ToNot(HaveOccurred())

			apiDefinitionFixture = fixtures.Api
			managementContextFixture = fixtures.Context

			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
			contextLookupKey = types.NamespacedName{Name: managementContextFixture.Name, Namespace: namespace}
		})

		It("should create an API Definition", func() {
			By("Create a management context to synchronize with the REST API")
			Expect(k8sClient.Create(ctx, managementContextFixture)).Should(Succeed())

			managementContext := new(gio.ManagementContext)
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContext)
			}, timeout, interval).Should(Succeed())

			By("Create an API definition resource referencing the management context")
			Expect(k8sClient.Create(ctx, apiDefinitionFixture)).Should(Succeed())

			apiDefinition := new(gio.ApiDefinition)
			Eventually(func() error {
				if err := k8sClient.Get(ctx, apiLookupKey, apiDefinition); err != nil {
					return err
				}
				return internal.AssertStatusIsSet(apiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			expectedApiName := apiDefinitionFixture.Spec.Name
			Expect(apiDefinition.Spec.Name).Should(Equal(expectedApiName))

			By("Call gateway endpoint and expect the API to be available")

			var endpoint = internal.GatewayUrl + apiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() error {
				res, callErr := httpClient.Get(endpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Call rest API and expect one API matching status cross ID")

			apim, err := internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() error {
				api, apiErr := apim.APIs.GetByID(apiDefinition.Status.ID)
				if apiErr != nil {
					return apiErr
				}

				return internal.AssertApiEntityMatchesStatus(api, apiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Check events")
			Expect(
				getEventsReason(apiDefinition.GetNamespace(), apiDefinition.GetName()),
			).Should(
				ContainElements([]string{"UpdateStarted", "UpdateSucceeded"}),
			)
		})

		It("should create a STOPPED API Definition", func() {
			apiDefinitionFixture.Spec.State = model.StateStopped

			By("Create a management context to synchronize with the REST API")
			Expect(k8sClient.Create(ctx, managementContextFixture)).Should(Succeed())

			managementContext := new(gio.ManagementContext)
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContext)
			}, timeout, interval).Should(Succeed())

			By("Create an API definition resource referencing the management context")
			Expect(k8sClient.Create(ctx, apiDefinitionFixture)).Should(Succeed())

			apiDefinition := new(gio.ApiDefinition)
			Eventually(func() error {
				if err := k8sClient.Get(ctx, apiLookupKey, apiDefinition); err != nil {
					return err
				}
				return internal.AssertStatusIsSet(apiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			expectedApiName := apiDefinitionFixture.Spec.Name
			Expect(apiDefinition.Spec.Name).Should(Equal(expectedApiName))

			By("Call gateway endpoint and expect the API not to be available")

			var endpoint = internal.GatewayUrl + apiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() bool {
				res, callErr := httpClient.Get(endpoint)
				return callErr == nil && res.StatusCode == 404
			}, timeout, interval).Should(BeTrue())

			By("Call rest API and expect one API matching status cross ID and state STOPPED")

			apim, err := internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() error {
				api, apiErr := apim.APIs.GetByID(apiDefinition.Status.ID)
				if apiErr != nil {
					return apiErr
				}

				if err = internal.AssertApiEntityMatchesStatus(api, apiDefinition); err != nil {
					return err
				}

				if api.State != "STOPPED" {
					return fmt.Errorf("expected state STOPPED, got %s", api.State)
				}

				return nil
			}, timeout, interval).ShouldNot(HaveOccurred())
		})

		It("should create an API Definition with existing api in Management Api", func() {
			apim, err := internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())

			By("Init existing api in management api")
			existingApiSpec := apiDefinitionFixture.Spec.DeepCopy()
			existingApiSpec.ID = uuid.NewV4String()
			existingApiSpec.CrossID = uuid.FromStrings(apiDefinitionFixture.GetNamespacedName().String())
			existingApiSpec.DefinitionContext = &model.DefinitionContext{
				Origin: origin,
				Mode:   mode,
			}
			existingApiSpec.Plans = []*model.Plan{
				{
					Id:       uuid.NewV4String(),
					Name:     "G.K.O. Default",
					Security: "KEY_LESS",
					Status:   "PUBLISHED",
				},
			}

			_, err = apim.APIs.Import(http.MethodPost, &existingApiSpec.Api)
			Expect(err).ToNot(HaveOccurred())

			By("Create a management context to synchronize with the REST API")
			Expect(k8sClient.Create(ctx, managementContextFixture)).Should(Succeed())

			managementContext := new(gio.ManagementContext)
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContext)
			}, timeout, interval).Should(Succeed())

			By("Create an API definition resource referencing the management context")
			Expect(k8sClient.Create(ctx, apiDefinitionFixture)).Should(Succeed())

			apiDefinition := new(gio.ApiDefinition)
			Eventually(func() error {
				if getErr := k8sClient.Get(ctx, apiLookupKey, apiDefinition); getErr != nil {
					return getErr
				}
				return internal.AssertStatusIsSet(apiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			expectedApiName := apiDefinitionFixture.Spec.Name
			Expect(apiDefinition.Spec.Name).Should(Equal(expectedApiName))

			By("Call gateway endpoint and expect the API to be available")

			var endpoint = internal.GatewayUrl + apiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() error {
				res, callErr := httpClient.Get(endpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Call rest API and expect one API matching status cross ID")

			Eventually(func() error {
				api, apiErr := apim.APIs.GetByID(apiDefinition.Status.ID)
				if apiErr != nil {
					return apiErr
				}
				return internal.AssertApiEntityMatchesStatus(api, apiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})

	})

	DescribeTable("a featured API spec with a management context",
		func(specFile string, expectedGatewayStatusCode int) {
			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:     specFile,
				Context: internal.ContextWithSecretFile,
			})

			Expect(err).ToNot(HaveOccurred())

			apiDefinitionFixture := fixtures.Api
			managementContextFixture := fixtures.Context

			apiLookupKey := types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
			contextLookupKey := types.NamespacedName{Name: managementContextFixture.Name, Namespace: namespace}

			By("Creating a management context to synchronize with the REST API")
			Expect(k8sClient.Create(ctx, managementContextFixture)).Should(Succeed())

			managementContext := new(gio.ManagementContext)
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContext)
			}, timeout, interval).Should(Succeed())

			By("Creating an API definition resource referencing the management context")
			Expect(k8sClient.Create(ctx, apiDefinitionFixture)).Should(Succeed())

			apiDefinition := new(gio.ApiDefinition)
			Eventually(func() error {
				if err = k8sClient.Get(ctx, apiLookupKey, apiDefinition); err != nil {
					return err
				}
				return internal.AssertStatusIsSet(apiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			expectedApiName := apiDefinitionFixture.Spec.Name
			Expect(apiDefinition.Spec.Name).Should(Equal(expectedApiName))

			By("Calling gateway endpoint, expecting the API to be available")

			var endpoint = internal.GatewayUrl + apiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() error {
				res, callErr := httpClient.Get(endpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, expectedGatewayStatusCode)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Calling rest API, expecting one API to match status cross ID")

			apim, err := internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() error {
				api, apiErr := apim.APIs.GetByID(apiDefinition.Status.ID)
				if apiErr != nil {
					return apiErr
				}
				return internal.AssertApiEntityMatchesStatus(api, apiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())
		},
		Entry("should import with health check", internal.ApiWithHCFile, 200),
		Entry("should import with disabled health check", internal.ApiWithDisabledHCFile, 200),
		Entry("should import with logging", internal.ApiWithLoggingFile, 200),
		Entry("should import with endpoint groups", internal.ApiWithEndpointGroupsFile, 200),
		Entry("should import with service discovery", internal.ApiWithServiceDiscoveryFile, 200),
		Entry("should import with metadata", internal.ApiWithMetadataFile, 200),
		Entry("should import with cache resource", internal.ApiWithCacheResourceFile, 200),
		Entry("should import with cache redis resource", internal.ApiWithCacheRedisResourceFile, 200),
		Entry("should import with oauth2 generic resource", internal.ApiWithOAuth2GenericResourceFile, 200),
		Entry("should import with oauth2 am resource", internal.ApiWithOauth2AmResourceFile, 200),
		Entry("should import with keycloak adapter resource", internal.ApiWithKeycloakAdapterFile, 200),
		Entry("should import with LDAP auth provider", internal.ApiWithLDAPAuthProviderFile, 401),
		Entry("should import with inline auth provider", internal.ApiWithInlineAuthProviderFile, 401),
		Entry("should import with HTTP auth provider", internal.ApiWithHTTPAuthProviderFile, 401),
	)

	DescribeTable("a featured API spec with a management context and a resource ref",
		func(resourceFile, specFile string, expectedGatewayStatusCode int) {
			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:      specFile,
				Context:  internal.ContextWithSecretFile,
				Resource: resourceFile,
			})

			Expect(err).ToNot(HaveOccurred())

			By("Creating a reusable resource to reference in the API")

			Expect(k8sClient.Create(ctx, fixtures.Resource)).Should(Succeed())

			apiDefinitionFixture := fixtures.Api
			managementContextFixture := fixtures.Context

			apiLookupKey := types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
			contextLookupKey := types.NamespacedName{Name: managementContextFixture.Name, Namespace: namespace}

			By("Creating a management context to synchronize with the REST API")
			Expect(k8sClient.Create(ctx, managementContextFixture)).Should(Succeed())

			managementContext := new(gio.ManagementContext)
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContext)
			}, timeout, interval).Should(Succeed())

			By("Creating an API definition resource referencing the management context")
			Expect(k8sClient.Create(ctx, apiDefinitionFixture)).Should(Succeed())

			apiDefinition := new(gio.ApiDefinition)
			Eventually(func() error {
				if err = k8sClient.Get(ctx, apiLookupKey, apiDefinition); err != nil {
					return err
				}
				return internal.AssertStatusIsSet(apiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Calling gateway endpoint, expecting the API to be available")

			var endpoint = internal.GatewayUrl + apiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() error {
				res, callErr := httpClient.Get(endpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, expectedGatewayStatusCode)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Calling rest API, expecting one API to match status cross ID")

			apim, err := internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() error {
				api, apiErr := apim.APIs.GetByID(apiDefinition.Status.ID)
				if apiErr != nil {
					return apiErr
				}
				return internal.AssertApiEntityMatchesStatus(api, apiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())
		},
		Entry(
			"should import with cache resource ref",
			internal.ApiResourceCacheFile,
			internal.ApiWithCacheResourceRefFile,
			200,
		),
		Entry(
			"should import with cache redis resource ref",
			internal.ApiResourceCacheRedisFile,
			internal.ApiWithCacheRedisResourceRefFile,
			200,
		),
		Entry(
			"should import with oauth2 generic resource ref",
			internal.ApiResourceOauth2GenericFile,
			internal.ApiWithOAuth2GenericResourceRefFile,
			200,
		),
		Entry(
			"should import with oauth2 am resource ref",
			internal.ApiResourceOauth2AMFile,
			internal.ApiWithOauth2AmResourceRefFile,
			200,
		),
		Entry(
			"should import with keycloak adapter resource ref",
			internal.ApiResourceKeycloakAdapterFile,
			internal.ApiWithKeycloakAdapterRefFile,
			200,
		),
		Entry(
			"should import with LDAP auth provider ref",
			internal.ApiResourceLDAPAuthProviderFile,
			internal.ApiWithLDAPAuthProviderRefFile,
			401,
		),
		Entry(
			"should import with inline auth provider ref",
			internal.ApiResourceInlineAuthProviderFile,
			internal.ApiWithInlineAuthProviderRefFile,
			401,
		),
		Entry(
			"should import with HTTP auth provider ref",
			internal.ApiResourceHTTPAuthProviderFile,
			internal.ApiWithHTTPAuthProviderRefFile,
			401,
		),
	)
})
