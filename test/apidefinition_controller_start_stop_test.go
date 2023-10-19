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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"
)

var _ = Describe("API Definition Controller", func() {

	httpClient := http.Client{Timeout: 5 * time.Second}

	Context("With Started basic ApiDefinition", func() {
		var apiDefinitionFixture *gio.ApiDefinition
		var apiLookupKey types.NamespacedName

		BeforeEach(func() {
			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:     internal.BasicApiFile,
				Context: internal.ClusterContextFile,
			})

			Expect(err).ToNot(HaveOccurred())

			By("Creating a management context")
			Expect(k8sClient.Create(ctx, fixtures.Context)).Should(Succeed())

			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{Name: fixtures.Context.Name, Namespace: namespace}, fixtures.Context)
			}, timeout, interval).Should(Succeed())

			By("Creating an API definition stared by default")
			apiDefinition := fixtures.Api
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
		})

		It("Should Stop an API Definition", func() {
			createdApiDefinition := new(gio.ApiDefinition)

			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, createdApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())
			Expect(createdApiDefinition.Spec.State).Should(Equal("STARTED"))

			By("Call initial API definition URL and expect 200")
			var gatewayEndpoint = internal.GatewayUrl + createdApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() error {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusOK)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Stop the API by define state to STOPPED")
			updatedApiDefinition := createdApiDefinition.DeepCopy()
			updatedApiDefinition.Spec.State = "STOPPED"
			Eventually(func() error {
				return internal.UpdateSafely(updatedApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Call updated API definition URL and expect 404")
			Eventually(func() error {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return internal.AssertNoErrorAndHTTPStatus(callErr, res, http.StatusNotFound)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Calling rest API and expecting state to be 'STOPPED'")
			apimClient, err := internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() error {
				err = k8sClient.Get(ctx, apiLookupKey, updatedApiDefinition)
				return internal.AssertEquals("API state", "STOPPED", updatedApiDefinition.Status.State)
			}, timeout, interval)

			Eventually(func() error {
				api, cliErr := apimClient.APIs.GetByID(updatedApiDefinition.Status.ID)
				if cliErr != nil {
					return cliErr
				}
				return internal.AssertEquals("API state", "STOPPED", api.State)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})
	})
})
