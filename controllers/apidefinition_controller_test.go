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
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/apim"
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
	gvk := gio.GroupVersion.WithKind("ApiDefinition")
	decode := scheme.Codecs.UniversalDecoder().Decode

	AfterEach(func() {
		// Delete the API definition
		Eventually(func() error {
			return k8sClient.DeleteAllOf(ctx, new(gio.ApiDefinition), &client.DeleteAllOfOptions{
				ListOptions:   client.ListOptions{Namespace: namespace},
				DeleteOptions: client.DeleteOptions{},
			})
		}).ShouldNot(HaveOccurred())

		// Delete the ManagementContext})
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
			const apiName = "K8s Basic Example"
			const endpoint = "http://localhost:9000/gateway/k8s-basic"

			crd, err := ioutil.ReadFile(sample)
			Expect(err).ToNot(HaveOccurred())

			decoded, _, err := decode(crd, &gvk, new(gio.ApiDefinition))
			Expect(err).ToNot(HaveOccurred())

			api, ok := decoded.(*gio.ApiDefinition)
			Expect(ok).To(BeTrue())

			By("Create an API definition resource referencing the management context")

			Expect(k8sClient.Create(ctx, api)).Should(Succeed())

			By("Get created resource and expect to find it")

			apiLookupKey := types.NamespacedName{Name: api.Name, Namespace: namespace}
			createdApi := new(gio.ApiDefinition)
			Eventually(func() bool {
				err = k8sClient.Get(ctx, apiLookupKey, createdApi)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(createdApi.Spec.Name).Should(Equal(apiName))

			By("Call gateway endpoint and expect the API to be available")

			Eventually(func() bool {
				res, callErr := httpClient.Get(endpoint)
				return callErr == nil && res.StatusCode == 200
			}, timeout, interval).Should(BeTrue())
		})

		It("Should create an API Definition", func() {
			By("With a management context")

			const apimCtxSample = "../config/samples/context/dev/managementcontext_credentials.yaml"
			const apiSample = "../config/samples/apim/basic-example-with-ctx.yml"
			const apiName = "K8s Basic Example With Management Context"
			const endpoint = "http://localhost:9000/gateway/k8s-basic-with-ctx"

			By("Create a management context to synchronize with the REST API")

			crdMgmtContext, err := ioutil.ReadFile(apimCtxSample)
			Expect(err).ToNot(HaveOccurred())

			decodedMgmtContext, _, err := decode(crdMgmtContext, &gvk, new(gio.ManagementContext))
			Expect(err).ToNot(HaveOccurred())

			mgmtContext, ok := decodedMgmtContext.(*gio.ManagementContext)
			Expect(ok).To(BeTrue())

			mgmtContext.Namespace = namespace
			Expect(k8sClient.Create(ctx, mgmtContext)).Should(Succeed())

			By("Create an API definition resource referencing the management context")

			crdApiDefinition, err := ioutil.ReadFile(apiSample)
			Expect(err).ToNot(HaveOccurred())

			decodedApiDefinition, _, err := decode(crdApiDefinition, &gvk, new(gio.ApiDefinition))
			Expect(err).ToNot(HaveOccurred())

			apiDefinition, ok := decodedApiDefinition.(*gio.ApiDefinition)
			Expect(ok).To(BeTrue())

			apiDefinition.Namespace = namespace

			Expect(k8sClient.Create(ctx, apiDefinition)).Should(Succeed())

			By("Get created resource and expect to find it")

			apiLookupKey := types.NamespacedName{Name: apiDefinition.Name, Namespace: namespace}
			createdApi := new(gio.ApiDefinition)
			Eventually(func() bool {
				err = k8sClient.Get(ctx, apiLookupKey, createdApi)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(createdApi.Spec.Name).Should(Equal(apiName))

			By("Call gateway endpoint and expect the API to be available")

			Eventually(func() bool {
				res, callErr := httpClient.Get(endpoint)
				return callErr == nil && res.StatusCode == 200
			}, timeout, interval).Should(BeTrue())

			By("Call rest API and expect one API matching status cross ID")

			apimClient := apim.NewClient(ctx, mgmtContext, httpClient)
			Eventually(func() bool {
				apis, apisErr := apimClient.FindByCrossId(createdApi.Status.ApiID)
				return apisErr == nil && len(apis) == 1
			}, timeout, interval).Should(BeTrue())

			apis, err := apimClient.FindByCrossId(createdApi.Status.ApiID)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(apis)).To(Equal(1))
		})
	})
})
