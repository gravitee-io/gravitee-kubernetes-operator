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
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

// +kubebuilder:docs-gen:collapse=Imports

var _ = Describe("API Definition Controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		namespace = "apim-dev"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	var ctx = context.Background()

	gvk := gio.GroupVersion.WithKind("ApiDefinition")
	decode := scheme.Codecs.UniversalDecoder().Decode

	httpClient := retryablehttp.NewClient()
	httpClient.RetryWaitMin = 2 * time.Second
	httpClient.RetryMax = 5
	httpClient.CheckRetry = func(ctx context.Context, res *http.Response, err error) (bool, error) {
		if err != nil {
			return true, err
		}
		if res.StatusCode != 200 {
			return true, errors.New(http.StatusText(http.StatusNotFound))
		}
		return false, nil
	}

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

		var apiDefinition *gio.ApiDefinition
		BeforeEach(func() {
			// Create the API definition
			const apiDefinitionSample = "../config/samples/apim/basic-example.yml"

			apiDefinitionCrd, err := ioutil.ReadFile(apiDefinitionSample)
			Expect(err).ToNot(HaveOccurred())

			apiDefinitionDecoded, _, err := decode(apiDefinitionCrd, &gvk, new(gio.ApiDefinition))
			Expect(err).ToNot(HaveOccurred())

			apiDefinition, ok := apiDefinitionDecoded.(*gio.ApiDefinition)
			Expect(ok).To(BeTrue())

			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			apiLookupKey := types.NamespacedName{Name: apiDefinition.Name, Namespace: namespace}
			createdApi := new(gio.ApiDefinition)
			Eventually(func() bool {
				err = k8sClient.Get(ctx, apiLookupKey, createdApi)
				return err == nil
			}, timeout, interval).Should(BeTrue())
		})

		It("Should update an API Definition", func() {
			By("Without a management context")

			// Check created api is callable
			const endpointInitial = "http://localhost:9000/gateway/k8s-basic"
			response, err := httpClient.Get(endpointInitial)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			// Update the API definition context path

			apiDefinition.Spec.Proxy.VirtualHosts[0].Path = "/k8s-basic-updated"
			apiDefinition.Spec.Proxy.Groups[0].Endpoints[0].Target = "https://api.gravitee.io/whattimeisit"

			Expect(k8sClient.Update(ctx, apiDefinition)).Should(Succeed())

			var apiDefinitionUpdated = gio.ApiDefinition{}
			apiLookupKey := types.NamespacedName{Name: apiDefinition.Name, Namespace: namespace}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, apiLookupKey, &apiDefinitionUpdated)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(apiDefinitionUpdated.Spec.Proxy.VirtualHosts[0].Path).To(Equal("/k8s-basic-updated"))

			const endpointUpdated = "http://localhost:9000/gateway/k8s-basic-updated"
			response, err = httpClient.Get(endpointUpdated)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
		})
	})
})
