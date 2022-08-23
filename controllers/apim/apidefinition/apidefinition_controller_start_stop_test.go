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
package apidefinition

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

var _ = Describe("API Definition Controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		namespace = "default"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	ctx := context.Background()
	httpClient := http.Client{Timeout: 5 * time.Second}

	Context("With Started basic ApiDefinition & ManagementContext", func() {
		var managementContextFixture *gio.ManagementContext
		var apiDefinitionFixture *gio.ApiDefinition
		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Create a management context to synchronize with the REST API")
			managementContext, err := test.NewManagementContext(
				"../../../config/samples/context/dev/managementcontext_credentials.yaml")
			Expect(err).ToNot(HaveOccurred())
			Expect(k8sClient.Create(ctx, managementContext)).Should(Succeed())

			By("Create an API definition resource stared by default")
			apiDefinition, err := test.NewApiDefinition("../../../config/samples/apim/basic-example-with-ctx.yml")
			Expect(err).ToNot(HaveOccurred())
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			managementContextFixture = managementContext
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
			contextLookupKey = types.NamespacedName{Name: managementContextFixture.Name, Namespace: namespace}
		})

		AfterEach(func() {
			Expect(k8sClient.Delete(ctx, apiDefinitionFixture)).Should(Succeed())

			Expect(k8sClient.Delete(ctx, managementContextFixture)).Should(Succeed())

			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, apiDefinitionFixture)
			}, timeout, interval).ShouldNot(Succeed())

			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContextFixture)
			}, timeout, interval).ShouldNot(Succeed())
		})

		It("Should Stop an API Definition", func() {
			createdApiDefinition := new(gio.ApiDefinition)

			// Expect the API Definition is Ready
			Eventually(func() bool {
				err := k8sClient.Get(ctx, apiLookupKey, createdApiDefinition)
				return err == nil && createdApiDefinition.Status.CrossID != ""
			}, timeout, interval).Should(BeTrue())

			By("Call initial API definition URL and expect no error")

			// Check created api is callable
			var gatewayEndpoint = test.GatewayUrl + createdApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() bool {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return callErr == nil && res.StatusCode == 200
			}, timeout, interval).Should(BeTrue())

			By("Stop the API by define state to STOPPED")

			updatedApiDefinition := createdApiDefinition.DeepCopy()

			updatedApiDefinition.Spec.State = "STOPPED"

			err := k8sClient.Update(ctx, updatedApiDefinition)
			Expect(err).ToNot(HaveOccurred())

			By("Call updated API definition URL and expect 404")
			Eventually(func() bool {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return callErr == nil && res.StatusCode == 404
			}, timeout, interval).Should(BeTrue())

			By("Call rest API and expect STOPPED state")
			apimClient := managementapi.NewClient(ctx, managementContextFixture, httpClient)
			Eventually(func() bool {
				api, apiErr := apimClient.GetByCrossId(updatedApiDefinition.Status.CrossID)

				return apiErr == nil &&
					api.Id == updatedApiDefinition.Status.ID &&
					api.State == "STOPPED"
			}, timeout, interval).Should(BeTrue())
		})
	})
})
