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

	Context("With Started basic ApiDefinition & ManagementContext", func() {
		var managementContextFixture *gio.ManagementContext
		var apiDefinitionFixture *gio.ApiDefinition
		var apiLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Create a management context to synchronize with the REST API")
			managementContext, err := internal.NewManagementContext(internal.ContextWithSecretFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(k8sClient.Create(ctx, managementContext)).Should(Succeed())

			By("Create an API definition resource stared by default")
			apiDefinition, err := internal.NewApiDefinition(internal.BasicApiWithContextFile)

			Expect(err).ToNot(HaveOccurred())
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			managementContextFixture = managementContext
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
		})

		AfterEach(func() {
			cleanupApiDefinitionAndManagementContext(apiDefinitionFixture, managementContextFixture)
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
			var gatewayEndpoint = internal.GatewayUrl + createdApiDefinition.Spec.Proxy.VirtualHosts[0].Path

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

			apimClient, err := internal.NewApimClient(ctx)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				api, apiErr := apimClient.GetByCrossId(updatedApiDefinition.Status.CrossID)

				return apiErr == nil &&
					api.Id == updatedApiDefinition.Status.ID &&
					api.State == "STOPPED"
			}, timeout, interval).Should(BeTrue())
		})
	})
})
