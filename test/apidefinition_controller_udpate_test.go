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

			apiDefinition, err := internal.NewApiDefinition(internal.BasicApiFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			apiLookupKey = types.NamespacedName{Name: apiDefinition.Name, Namespace: namespace}
		})

		It("Should update an API Definition", func() {
			createdApiDefinition := new(gio.ApiDefinition)

			Eventually(func() bool {
				err := k8sClient.Get(ctx, apiLookupKey, createdApiDefinition)
				return err == nil && createdApiDefinition.Status.CrossID != ""
			}, timeout, interval).Should(BeTrue())

			By("Call initial API definition URL and expect no error")

			// Check created api is callable
			var endpointInitial = internal.GatewayUrl + createdApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() bool {
				res, callErr := httpClient.Get(endpointInitial)
				return callErr == nil && res.StatusCode == 200
			}, timeout, interval).Should(BeTrue())

			By("Update the context path in API definition and expect no error")

			updatedApiDefinition := createdApiDefinition.DeepCopy()

			expectedPath := updatedApiDefinition.Spec.Proxy.VirtualHosts[0].Path + "-updated"
			updatedApiDefinition.Spec.Proxy.VirtualHosts[0].Path = expectedPath

			err := k8sClient.Update(ctx, updatedApiDefinition)
			Expect(err).ToNot(HaveOccurred())

			By("Call updated API definition URL and expect no error")

			var endpointUpdated = internal.GatewayUrl + expectedPath

			Eventually(func() bool {
				res, callErr := httpClient.Get(endpointUpdated)
				return callErr == nil && res.StatusCode == 200
			}, timeout, interval).Should(BeTrue())

			By("Check events")
			Expect(
				getEventsReason(apiDefinitionFixture),
			).Should(
				ContainElements([]string{"Updated", "Updating"}),
			)
		})
	})

	Context("With basic ApiDefinition & ManagementContext", func() {
		var apiDefinitionFixture *gio.ApiDefinition
		var apiLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Create a management context to synchronize with the REST API")

			apiWithContext, err := internal.NewApiWithRandomContext(
				internal.BasicApiWithContextFile, internal.ContextWithSecretFile,
			)

			Expect(err).ToNot(HaveOccurred())

			managementContext := apiWithContext.Context
			Expect(k8sClient.Create(ctx, managementContext)).Should(Succeed())

			By("Create an API definition resource stared by default")

			apiDefinition := apiWithContext.Api
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
		})

		It("Should update an API Definition", func() {
			createdApiDefinition := new(gio.ApiDefinition)

			Eventually(func() bool {
				err := k8sClient.Get(ctx, apiLookupKey, createdApiDefinition)
				return err == nil && createdApiDefinition.Status.CrossID != ""
			}, timeout, interval).Should(BeTrue())

			By("Call initial API definition URL and expect no error")

			// Check created api is callable
			var endpointInitial = internal.GatewayUrl + createdApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() bool {
				res, callErr := httpClient.Get(endpointInitial)
				return callErr == nil && res.StatusCode == 200
			}, timeout, interval).Should(BeTrue())

			By("Update the context path in API definition and expect no error")

			updatedApiDefinition := createdApiDefinition.DeepCopy()

			expectedPath := updatedApiDefinition.Spec.Proxy.VirtualHosts[0].Path + "-updated"
			expectedName := updatedApiDefinition.Spec.Name + "-updated"
			updatedApiDefinition.Spec.Proxy.VirtualHosts[0].Path = expectedPath
			updatedApiDefinition.Spec.Name = expectedName

			err := k8sClient.Update(ctx, updatedApiDefinition)
			Expect(err).ToNot(HaveOccurred())

			By("Call updated API definition URL and expect no error")

			var endpointUpdated = internal.GatewayUrl + expectedPath

			Eventually(func() bool {
				res, callErr := httpClient.Get(endpointUpdated)
				return callErr == nil && res.StatusCode == 200
			}, timeout, interval).Should(BeTrue())

			By("Call rest API and expect one API matching status cross ID & updated name")

			apimClient, err := internal.NewApimClient(ctx)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				api, apiErr := apimClient.GetByCrossId(updatedApiDefinition.Status.CrossID)
				return apiErr == nil &&
					api.Id == updatedApiDefinition.Status.ID &&
					api.Name == expectedName
			}, timeout, interval).Should(BeTrue())
		})
	})

	Context("With basic ApiDefinition & ManagementContext adding context ref on update", func() {
		var managementContextFixture *gio.ManagementContext
		var apiDefinitionFixture *gio.ApiDefinition
		var apiLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Create a management context to synchronize with the REST API")

			apiWithContext, err := internal.NewApiWithRandomContext(
				internal.BasicApiWithContextFile, internal.ContextWithSecretFile,
			)

			Expect(err).ToNot(HaveOccurred())

			managementContext := apiWithContext.Context
			Expect(k8sClient.Create(ctx, managementContext)).Should(Succeed())

			By("Create an API definition resource stared by default")

			apiDefinition := apiWithContext.Api
			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiDefinitionFixture = apiDefinition
			managementContextFixture = managementContext
			apiLookupKey = types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
		})

		It("Should update an API Definition, adding a management context", func() {
			createdApiDefinition := new(gio.ApiDefinition)

			Eventually(func() bool {
				err := k8sClient.Get(ctx, apiLookupKey, createdApiDefinition)
				return err == nil && createdApiDefinition.Status.CrossID != ""
			}, timeout, interval).Should(BeTrue())

			By("Updating the context ref in API definition, expecting no error")

			updatedApiDefinition := createdApiDefinition.DeepCopy()

			updatedApiDefinition.Spec.Context = &model.ContextRef{
				Name:      managementContextFixture.Name,
				Namespace: managementContextFixture.Namespace,
			}

			err := k8sClient.Update(ctx, updatedApiDefinition)
			Expect(err).ToNot(HaveOccurred())

			By("Calling rest API, expecting one API matching status ID")

			apimClient, err := internal.NewApimClient(ctx)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() bool {
				api, apiErr := apimClient.GetByCrossId(updatedApiDefinition.Status.CrossID)
				return apiErr == nil && api.Id == updatedApiDefinition.Status.ID
			}, timeout, interval).Should(BeTrue())
		})
	})
})
