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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("API Definition Controller", func() {

	httpClient := http.Client{Timeout: 5 * time.Second}

	Context("With STARTED ApiDefinition & ManagementContext", func() {
		var apiDefinitionFixture *v1alpha1.ApiDefinition
		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Create a management context to synchronize with the REST API")

			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:     internal.ApiWithContextFile,
				Context: internal.ClusterContextFile,
			})

			Expect(err).ToNot(HaveOccurred())

			apiDefinition := fixtures.Api
			managementContext := fixtures.Context
			contextLookupKey = types.NamespacedName{Name: managementContext.Name, Namespace: namespace}
			apiLookupKey = types.NamespacedName{Name: apiDefinition.Name, Namespace: namespace}

			Expect(err).ToNot(HaveOccurred())
			Expect(k8sClient.Create(ctx, managementContext)).Should(Succeed())

			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContext)
			}, timeout, interval).Should(Succeed())

			By("Create an API definition resource in the STARTED state")

			Expect(err).ToNot(HaveOccurred())
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, apiDefinition)
			}, timeout, interval).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
		})

		It("Should Delete an API Definition", func() {

			// Expect the API Definition is Ready
			createdApiDefinition := new(v1alpha1.ApiDefinition)
			Eventually(func() error {
				err := k8sClient.Get(ctx, apiLookupKey, createdApiDefinition)
				return internal.AssertNoErrorAndStatusCompleted(err, createdApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Call initial API definition URL and expect no error")

			// Check created api is callable
			var gatewayEndpoint = internal.GatewayUrl + createdApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() error {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Delete the API Definition")

			err := k8sClient.Delete(ctx, apiDefinitionFixture)
			Expect(err).ToNot(HaveOccurred())

			By("Call deleted API definition URL and expect 404")
			Eventually(func() error {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusNotFound)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Call rest API and expect DELETED api")

			apim, err := internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() error {
				_, apiErr := apim.APIs.GetByID(createdApiDefinition.Status.ID)
				return errors.IgnoreNotFound(apiErr)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Expect that the ConfigMap has been deleted")
			cm := &v1.ConfigMap{}
			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{
					Name:      createdApiDefinition.Name,
					Namespace: createdApiDefinition.Namespace,
				}, cm)
			}, timeout, interval).ShouldNot(Succeed())

			By("Check events")
			Eventually(
				getEventReasons(apiDefinitionFixture),
				timeout, interval,
			).Should(
				ContainElements([]string{"DeleteSucceeded", "DeleteStarted"}),
			)
		})

		It("Should detect when API has already been deleted", func() {
			apim, err := internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())

			createdApiDefinition := new(v1alpha1.ApiDefinition)
			Eventually(func() error {
				err = k8sClient.Get(ctx, apiLookupKey, createdApiDefinition)
				return internal.AssertNoErrorAndStatusCompleted(err, createdApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Call initial API definition URL and expect no error")

			var gatewayEndpoint = internal.GatewayUrl + createdApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() error {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Delete the API calling directly the REST API")

			Eventually(func() error {
				return apim.APIs.Delete(createdApiDefinition.Status.ID)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Delete the API Definition")
			err = k8sClient.Delete(ctx, apiDefinitionFixture)
			Expect(err).ToNot(HaveOccurred())

			By("Expect that the ConfigMap has been deleted")
			cm := &v1.ConfigMap{}
			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{
					Name:      createdApiDefinition.Name,
					Namespace: createdApiDefinition.Namespace,
				}, cm)
			}, timeout, interval).ShouldNot(Succeed())
		})
	})
})
