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

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

// +kubebuilder:docs-gen:collapse=Imports

var _ = Describe("API Definition Controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		namespace = "default"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	gvk := gio.GroupVersion.WithKind("ApiDefinition")
	decode := scheme.Codecs.UniversalDecoder().Decode

	cli := retryablehttp.NewClient()
	cli.RetryWaitMin = 2 * time.Second
	cli.RetryMax = 5
	cli.CheckRetry = func(ctx context.Context, res *http.Response, err error) (bool, error) {
		if err != nil {
			return true, err
		}
		if res.StatusCode != 200 {
			return true, errors.New(http.StatusText(http.StatusNotFound))
		}
		return false, nil
	}

	Context("API definition Resource", func() {
		It("Should create an API Definition", func() {
			By("Without a management context")

			const sample = "../config/samples/apim/basic-example.yml"
			const apiName = "K8s Basic Example"
			const apiUrl = "http://localhost:9000/gateway/k8s-basic"

			ctx := context.Background()

			crd, err := ioutil.ReadFile(sample)
			Expect(err).ToNot(HaveOccurred())

			decoded, _, err := decode(crd, &gvk, new(gio.ApiDefinition))
			Expect(err).ToNot(HaveOccurred())

			api, ok := decoded.(*gio.ApiDefinition)
			Expect(ok).To(BeTrue())

			api.Namespace = namespace

			Expect(err).ToNot(HaveOccurred())
			Expect(k8sClient.Create(ctx, api)).Should(Succeed())

			apiLookupKey := types.NamespacedName{Name: api.Name, Namespace: namespace}
			createdApi := new(gio.ApiDefinition)

			Eventually(func() bool {
				err = k8sClient.Get(ctx, apiLookupKey, createdApi)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(createdApi.Spec.Name).Should(Equal(apiName))

			Expect(k8sClient.Delete(ctx, api)).Should(Succeed())

			res, err := cli.Get(apiUrl)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(200))

		})
	})
})
