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

	"errors"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi"
	clientError "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi/clienterror"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
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
			err := k8sClient.Delete(ctx, apiDefinitionFixture)
			if err != nil {
				// wait deleted only if not already deleted
				Eventually(func() error {
					return k8sClient.Get(ctx, apiLookupKey, apiDefinitionFixture)
				}, timeout, interval).ShouldNot(Succeed())
			}

			err = k8sClient.Delete(ctx, managementContextFixture)
			if err != nil {
				// wait deleted only if not already deleted
				Eventually(func() error {
					return k8sClient.Get(ctx, contextLookupKey, managementContextFixture)
				}, timeout, interval).ShouldNot(Succeed())
			}
		})

		It("Should Delete an API Definition", func() {
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

			By("Delete the API Definition")

			err := k8sClient.Delete(ctx, apiDefinitionFixture)
			Expect(err).ToNot(HaveOccurred())

			By("Call deleted API definition URL and expect 404")
			Eventually(func() bool {
				res, callErr := httpClient.Get(gatewayEndpoint)
				return callErr == nil && res.StatusCode == 404
			}, timeout, interval).Should(BeTrue())

			By("Call rest API and expect DELETED api")
			apimClient := managementapi.NewClient(ctx, managementContextFixture, httpClient)
			Eventually(func() bool {
				_, apiErr := apimClient.GetApiById(createdApiDefinition.Status.CrossID)
				var apiNotFoundError *clientError.ApiNotFoundError
				return apiErr != nil && errors.As(apiErr, &apiNotFoundError)
			}, timeout, interval).Should(BeTrue())

			By("Expect that the ConfigMap has been deleted")
			cm := &v1.ConfigMap{}
			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{
					Name:      createdApiDefinition.Name,
					Namespace: createdApiDefinition.Namespace,
				}, cm)
			}, timeout, interval).ShouldNot(Succeed())
		})

		It("Should detect when API has already been deleted", func() {
			createdApiDefinition := new(gio.ApiDefinition)
			managementClient := managementapi.NewClient(ctx, managementContextFixture, httpClient)

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

			By("Delete the API calling directly the REST API")
			deleteErr := managementClient.DeleteApi(createdApiDefinition.Status.ID)
			Expect(deleteErr).ToNot(HaveOccurred())

			By("Delete the API Definition")
			err := k8sClient.Delete(ctx, apiDefinitionFixture)
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
