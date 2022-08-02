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
// +kubebuilder:docs-gen:collapse=Apache License
package controllers

import (
	"context"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test"
)

// +kubebuilder:docs-gen:collapse=Imports

var _ = Describe("API Definition Controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		namespace = "default"

		timeout  = time.Second * 10
		interval = time.Millisecond * 500
	)

	ctx := context.Background()
	httpClient := http.Client{Timeout: 5 * time.Second}

	AfterEach(func() {
		// Delete the API definition
		Eventually(func() error {
			return k8sClient.DeleteAllOf(ctx, new(gio.ApiDefinition), &client.DeleteAllOfOptions{
				ListOptions:   client.ListOptions{Namespace: namespace},
				DeleteOptions: client.DeleteOptions{},
			})
		}).ShouldNot(HaveOccurred())

		// Delete the ManagementContext
		Eventually(func() error {
			return k8sClient.DeleteAllOf(ctx, new(gio.ManagementContext), &client.DeleteAllOfOptions{
				ListOptions:   client.ListOptions{Namespace: namespace},
				DeleteOptions: client.DeleteOptions{},
			})
		}).ShouldNot(HaveOccurred())
	})

	Context("API definition Resource", func() {
		It("Should create an API Definition", func() {
			By("Without a management context")

			const sample = "../config/samples/apim/basic-example.yml"

			api, err := test.NewApiDefinition(sample)
			Expect(err).ToNot(HaveOccurred())

			expectedApiName := api.Spec.Name

			By("Create an API definition resource without a management context")

			Expect(k8sClient.Create(ctx, api)).Should(Succeed())

			By("Get created resource and expect to find it")

			apiLookupKey := types.NamespacedName{Name: api.Name, Namespace: namespace}
			createdApi := new(gio.ApiDefinition)
			Eventually(func() bool {
				err = k8sClient.Get(ctx, apiLookupKey, createdApi)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			var endpoint = test.GatewayUrl + createdApi.Spec.Proxy.VirtualHosts[0].Path

			Expect(createdApi.Spec.Name).Should(Equal(expectedApiName))

			By("Call gateway endpoint and expect the API to be available")

			Eventually(func() bool {
				res, callErr := httpClient.Get(endpoint)
				return callErr == nil && res.StatusCode == 200
			}, timeout, interval).Should(BeTrue())
		})

		It("Should create an API Definition", func() {
			By("With a management context")

			const mgmtContextSample = "../config/samples/context/dev/managementcontext_credentials.yaml"
			const apiSample = "../config/samples/apim/basic-example-with-ctx.yml"

			By("Create a management context to synchronize with the REST API")

			mgmtContext, err := test.NewManagementContext(mgmtContextSample)
			Expect(err).ToNot(HaveOccurred())
			Expect(k8sClient.Create(ctx, mgmtContext)).Should(Succeed())

			By("Create an API definition resource referencing the management context")

			apiDefinition, err := test.NewApiDefinition(apiSample)
			Expect(err).ToNot(HaveOccurred())

			expectedApiName := apiDefinition.Spec.Name

			apiDefinition.Namespace = namespace

			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			By("Get created resource and expect to find it")

			apiLookupKey := types.NamespacedName{Name: apiDefinition.Name, Namespace: namespace}
			createdApi := new(gio.ApiDefinition)
			Eventually(func() bool {
				err = k8sClient.Get(ctx, apiLookupKey, createdApi)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(createdApi.Spec.Name).Should(Equal(expectedApiName))

			By("Call gateway endpoint and expect the API to be available")

			var endpoint = test.GatewayUrl + createdApi.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() bool {
				res, callErr := httpClient.Get(endpoint)
				return callErr == nil && res.StatusCode == 200
			}, timeout, interval).Should(BeTrue())

			By("Call rest API and expect one API matching status cross ID")

			apimClient := apim.NewClient(ctx, mgmtContext, httpClient)
			Eventually(func() bool {
				apis, apisErr := apimClient.FindByCrossId(createdApi.Status.CrossID)
				return apisErr == nil && len(apis) == 1
			}, timeout, interval).Should(BeTrue())

			apis, err := apimClient.FindByCrossId(createdApi.Status.CrossID)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(apis)).To(Equal(1))
		})
	})
})
