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
package controllers

import (
	"context"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	managementapi "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test"
)

var _ = Describe("Checking ApiKey plan and subscription", Ordered, func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		namespace = "default"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	ctx := context.Background()
	httpClient := http.Client{Timeout: 5 * time.Second}

	Context("Checking ApiKey plan and subscription", Ordered, func() {
		var managementContextFixture *gio.ManagementContext
		var apiDefinitionFixture *gio.ApiDefinition

		var createdApiDefinition *gio.ApiDefinition

		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName

		var mgmtClient *managementapi.Client

		BeforeAll(func() {
			By("Create a management context to synchronize with the REST API")
			managementContext, err := test.NewManagementContext(
				"../config/samples/context/dev/managementcontext_credentials.yaml")
			Expect(err).ToNot(HaveOccurred())
			Expect(k8sClient.Create(ctx, managementContext)).Should(Succeed())

			By("Create an API definition resource stared by default")
			apiDefinition, err := test.NewApiDefinition("../config/samples/apim/apikey-example-with-ctx.yml")
			Expect(err).ToNot(HaveOccurred())
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			managementContextFixture = managementContext
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
			contextLookupKey = types.NamespacedName{Name: managementContextFixture.Name, Namespace: namespace}

			By("Expect the API Definition is Ready")
			createdApiDefinition = new(gio.ApiDefinition)
			Eventually(func() bool {
				k8sErr := k8sClient.Get(ctx, apiLookupKey, createdApiDefinition)
				return k8sErr == nil && createdApiDefinition.Status.CrossID != ""
			}, timeout, interval).Should(BeTrue())

			mgmtClient = managementapi.NewClient(ctx, managementContextFixture, httpClient)
		})

		AfterAll(func() {
			Expect(k8sClient.Delete(ctx, apiDefinitionFixture)).Should(Succeed())

			Expect(k8sClient.Delete(ctx, managementContextFixture)).Should(Succeed())

			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, apiDefinitionFixture)
			}, timeout, interval).ShouldNot(Succeed())

			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContextFixture)
			}, timeout, interval).ShouldNot(Succeed())
		})

		It("Should return unauthorize without subscription", func() {
			var gatewayEndpoint = test.GatewayUrl + createdApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() bool {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return callErr == nil && res.StatusCode == 401
			}, timeout, interval).Should(BeTrue())
		})

		It("Should return success with subscription", func() {
			// Check created api is callable
			var gatewayEndpoint = test.GatewayUrl + createdApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			By("Create a subscription and get api key")
			// Get first active application
			mgmtApplications, mgmtErr := mgmtClient.SearchApplications("", "ACTIVE")
			Expect(mgmtErr).ToNot(HaveOccurred())
			defaultApplication := mgmtApplications[0]

			// Get Api description with plan
			mgmtApi, mgmtErr := mgmtClient.GetApiById(createdApiDefinition.Status.ID)
			Expect(mgmtErr).ToNot(HaveOccurred())

			// Create subscription
			mgmtSubscription, mgmtErr := mgmtClient.SubscribeToPlan(mgmtApi.Id, defaultApplication.Id, mgmtApi.Plans[0].Id)
			Expect(mgmtErr).ToNot(HaveOccurred())

			// Get subscription api keys
			mgmtSubscriptionApiKeys, mgmtErr := mgmtClient.GetSubscriptionApiKey(mgmtApi.Id, mgmtSubscription.Id)
			Expect(mgmtErr).ToNot(HaveOccurred())

			By("Call gateway with subscription api key")
			Eventually(func() bool {
				req, callErr := http.NewRequest("GET", gatewayEndpoint, nil)
				if callErr != nil {
					return false
				}
				req.Header.Set("X-Gravitee-Api-Key", mgmtSubscriptionApiKeys[0].Key)

				res, callErr := httpClient.Do(req)
				if callErr != nil {
					return false
				}

				if res.Body != nil {
					defer res.Body.Close()
				}

				return callErr == nil && res.StatusCode == 200
			}, timeout, interval).Should(BeTrue())
		})
	})
})
