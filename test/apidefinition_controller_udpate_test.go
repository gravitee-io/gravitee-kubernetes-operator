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

/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package test

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"
)

var _ = Describe("API Definition Controller", func() {

	httpClient := http.Client{Timeout: 5 * time.Second}

	Context("With basic ApiDefinition", func() {
		var apiDefinitionFixture *gio.ApiDefinition
		var apiLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Create an API definition resource without a management context")

			fixtureGenerator := internal.NewFixtureGenerator()

			fixture, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api: internal.BasicApiFile,
			})

			Expect(err).ToNot(HaveOccurred())

			apiDefinition := fixture.Api

			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			apiLookupKey = types.NamespacedName{Name: apiDefinition.Name, Namespace: namespace}
		})

		It("Should update an API Definition", func() {
			createdApiDefinition := new(gio.ApiDefinition)

			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, createdApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Call initial API definition URL and expect no error")

			// Check created api is callable
			var endpointInitial = internal.GatewayUrl + createdApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() error {
				res, callErr := httpClient.Get(endpointInitial)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Update the context path in API definition and expect no error")

			updatedApiDefinition := createdApiDefinition.DeepCopy()

			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, updatedApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			expectedPath := updatedApiDefinition.Spec.Proxy.VirtualHosts[0].Path + "-updated"
			updatedApiDefinition.Spec.Proxy.VirtualHosts[0].Path = expectedPath

			Eventually(func() error {
				update := new(gio.ApiDefinition)
				if err := k8sClient.Get(ctx, apiLookupKey, update); err != nil {
					return err
				}
				updatedApiDefinition.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Call updated API definition URL and expect no error")

			var endpointUpdated = internal.GatewayUrl + expectedPath

			Eventually(func() error {
				res, callErr := httpClient.Get(endpointUpdated)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Check events")
			Expect(
				getEventsReason(apiDefinitionFixture.GetNamespace(), apiDefinitionFixture.GetName()),
			).Should(
				ContainElements([]string{"UpdateSucceeded", "UpdateStarted"}),
			)
		})
	})

	Context("With basic ApiDefinition & ManagementContext", func() {
		var apiDefinitionFixture *gio.ApiDefinition
		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Create a management context to synchronize with the REST API")

			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:     internal.BasicApiFile,
				Context: internal.ContextWithSecretFile,
			})

			Expect(err).ToNot(HaveOccurred())

			managementContext := fixtures.Context
			Expect(k8sClient.Create(ctx, managementContext)).Should(Succeed())

			contextLookupKey = types.NamespacedName{Name: managementContext.Name, Namespace: namespace}
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContext)
			}, timeout, interval).Should(Succeed())

			By("Create an API definition resource stared by default")

			apiDefinition := fixtures.Api
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
		})

		It("Should update an API Definition", func() {
			createdApiDefinition := new(gio.ApiDefinition)

			Eventually(func() error {
				if err := k8sClient.Get(ctx, apiLookupKey, createdApiDefinition); err != nil {
					return err
				}
				return internal.AssertStatusIsSet(createdApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Call initial API definition URL and expect no error")

			// Check created api is callable
			var endpointInitial = internal.GatewayUrl + createdApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() error {
				res, callErr := httpClient.Get(endpointInitial)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Update the context path in API definition and expect no error")

			updatedApiDefinition := createdApiDefinition.DeepCopy()

			expectedPath := updatedApiDefinition.Spec.Proxy.VirtualHosts[0].Path + "-updated"
			expectedName := updatedApiDefinition.Spec.Name + "-updated"
			updatedApiDefinition.Spec.Proxy.VirtualHosts[0].Path = expectedPath
			updatedApiDefinition.Spec.Name = expectedName

			Eventually(func() error {
				update := new(gio.ApiDefinition)
				if err := k8sClient.Get(ctx, apiLookupKey, update); err != nil {
					return err
				}
				updatedApiDefinition.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Call updated API definition URL and expect no error")

			var endpointUpdated = internal.GatewayUrl + expectedPath

			Eventually(func() bool {
				res, callErr := httpClient.Get(endpointUpdated)
				return callErr == nil && res.StatusCode == 200
			}, timeout, interval).Should(BeTrue())

			By("Call rest API and expect one API matching status cross ID & updated name")

			apimClient, err := internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() error {
				api, cliErr := apimClient.APIs.GetByID(updatedApiDefinition.Status.ID)
				if cliErr != nil {
					return cliErr
				}

				if api.Name != expectedName {
					return fmt.Errorf("API name mismatch: %s != %s", api.Name, expectedName)
				}

				return nil
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Calling rest API, expecting one API matching status ID and kubernetes context")

			Eventually(func() error {
				api, cliErr := apimClient.APIs.GetByID(updatedApiDefinition.Status.ID)
				if cliErr != nil {
					return cliErr
				}

				if api.DefinitionContext.Mode != "fully_managed" {
					return fmt.Errorf("API mode mismatch: %s != %s", api.DefinitionContext.Mode, "fully_managed")
				}

				if api.DefinitionContext.Origin != "kubernetes" {
					return fmt.Errorf("API origin mismatch: %s != %s", api.DefinitionContext.Mode, "kubernetes")
				}

				return nil
			}, timeout, interval).ShouldNot(HaveOccurred())
		})
	})

	Context("With basic ApiDefinition & ManagementContext adding context ref on update", func() {
		var contextFixture *gio.ManagementContext
		var apiDefinitionFixture *gio.ApiDefinition
		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Create a management context to synchronize with the REST API")

			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:     internal.BasicApiFile,
				Context: internal.ContextWithSecretFile,
			})

			Expect(err).ToNot(HaveOccurred())

			managementContext := fixtures.Context
			Expect(k8sClient.Create(ctx, managementContext)).Should(Succeed())

			contextLookupKey = types.NamespacedName{Name: managementContext.Name, Namespace: namespace}
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContext)
			}, timeout, interval).Should(Succeed())

			By("Create an API definition resource stared by default")

			apiDefinition := fixtures.Api
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			contextFixture = managementContext
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
		})

		It("Should update an API Definition, adding a management context", func() {
			createdApiDefinition := new(gio.ApiDefinition)

			Eventually(func() error {
				if err := k8sClient.Get(ctx, apiLookupKey, createdApiDefinition); err != nil {
					return err
				}
				return internal.AssertStatusIsSet(createdApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Updating the context ref in API definition, expecting no error")

			updatedApiDefinition := createdApiDefinition.DeepCopy()

			updatedApiDefinition.Spec.Context = &model.NamespacedName{
				Name:      contextFixture.Name,
				Namespace: contextFixture.Namespace,
			}

			Eventually(func() error {
				update := new(gio.ApiDefinition)
				if err := k8sClient.Get(ctx, apiLookupKey, update); err != nil {
					return err
				}
				updatedApiDefinition.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Calling rest API, expecting one API matching status ID")

			apimClient, err := internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() error {
				_, err = apimClient.APIs.GetByID(updatedApiDefinition.Status.ID)
				return err
			}, timeout, interval).ShouldNot(HaveOccurred())
		})
	})
})
