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
	"net/http"
	"time"

	v1 "k8s.io/api/core/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"
)

var _ = Describe("Create a basic API", func() {
	httpClient := http.Client{Timeout: 5 * time.Second}

	Context("with a management context", func() {
		var apiDefinitionFixture *gio.ApiDefinition
		var managementContextFixture *gio.ManagementContext
		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName

		BeforeEach(func() {
			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:     internal.ApiWithContextFile,
				Context: internal.ClusterContextFile,
			})

			Expect(err).ToNot(HaveOccurred())

			apiDefinitionFixture = fixtures.Api
			managementContextFixture = fixtures.Context

			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
			contextLookupKey = types.NamespacedName{Name: managementContextFixture.Name, Namespace: namespace}
		})

		It("should create an API Definition and a ConfigMap when `local` attribute equals true", func() {
			By("Creating a management context to synchronize with the REST API")
			Expect(k8sClient.Create(ctx, managementContextFixture)).Should(Succeed())

			managementContext := new(gio.ManagementContext)
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContext)
			}, timeout, interval).Should(Succeed())

			By("Creating an API definition resource referencing the management context")
			apiDefinitionFixture.Spec.IsLocal = true
			Expect(k8sClient.Create(ctx, apiDefinitionFixture)).Should(Succeed())

			apiDefinition := new(gio.ApiDefinition)
			Eventually(func() error {
				err := k8sClient.Get(ctx, apiLookupKey, apiDefinition)
				return internal.AssertNoErrorAndStatusCompleted(err, apiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Expecting the ConfigMap has been created")
			cm := &v1.ConfigMap{}
			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, cm)
			}, timeout, interval).Should(Succeed())

			By("Calling the gateway endpoint and expect the API to be available")
			var endpoint = internal.GatewayUrl + apiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() error {
				res, callErr := httpClient.Get(endpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})

		It("should only create an API Definition when `local` attribute equals false", func() {
			By("Creating a management context to synchronize with the REST API")
			Expect(k8sClient.Create(ctx, managementContextFixture)).Should(Succeed())

			managementContext := new(gio.ManagementContext)
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContext)
			}, timeout, interval).Should(Succeed())

			By("Creating an API definition resource referencing the management context")
			apiDefinitionFixture.Spec.IsLocal = false
			Expect(k8sClient.Create(ctx, apiDefinitionFixture)).Should(Succeed())

			apiDefinition := new(gio.ApiDefinition)
			Eventually(func() error {
				if err := k8sClient.Get(ctx, apiLookupKey, apiDefinition); err != nil {
					return err
				}
				return internal.AssertEquals("local", apiDefinition.Spec.IsLocal, false)
			}, timeout, interval).Should(Succeed())

			By("Expecting the ConfigMap has not been created")
			cm := &v1.ConfigMap{}
			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, cm)
			}, timeout, interval).ShouldNot(Succeed())

			By("Calling the gateway endpoint and expect the API to be available")
			var endpoint = internal.GatewayUrl + apiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() error {
				res, callErr := httpClient.Get(endpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})

		It("should remove the ConfigMap when user switch from `local` equals true to false ", func() {
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
				return k8sClient.Get(ctx, apiLookupKey, apiDefinition)
			}, timeout, interval).Should(Succeed())
			Expect(apiDefinition.Spec.IsLocal).To(BeTrue())

			By("Expecting the ConfigMap has not been created")
			cm := &v1.ConfigMap{}
			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, cm)
			}, timeout, interval).ShouldNot(Succeed())

			By("Updating the ManagementContext to switch from `local` equals true to false")
			Eventually(func() error {
				createdApiDefinition := new(gio.ApiDefinition)
				if err := k8sClient.Get(ctx, apiLookupKey, createdApiDefinition); err != nil {
					return err
				}
				createdApiDefinition.Spec.IsLocal = false
				return k8sClient.Update(ctx, createdApiDefinition)
			}).Should(Succeed())

			updatedApiDefinition := new(gio.ApiDefinition)
			Eventually(func() error {
				if err := k8sClient.Get(ctx, apiLookupKey, updatedApiDefinition); err != nil {
					return err
				}
				return internal.AssertEquals("local", updatedApiDefinition.Spec.IsLocal, false)
			}, timeout, interval).Should(Succeed())

			By("Expecting the ConfigMap has been removed")
			removedConfigMap := &v1.ConfigMap{}
			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, removedConfigMap)
			}, timeout, interval).ShouldNot(Succeed())
		})
	})
})
