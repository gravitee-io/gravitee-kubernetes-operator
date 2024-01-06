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

package apiresource

import (
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("API Resource Controller", func() {
	Context("Delete an API Resource with no reference to an API", func() {
		var apiLookupKey types.NamespacedName
		var resourceLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Creating an API and an API resource")

			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:      internal.BasicApiFile,
				Resource: internal.ApiResourceCacheFile,
			})

			Expect(err).ToNot(HaveOccurred())

			apiLookupKey = types.NamespacedName{Name: fixtures.Api.Name, Namespace: internal.Namespace}
			resourceLookupKey = types.NamespacedName{Name: fixtures.Resource.Name, Namespace: internal.Namespace}

			Expect(k8sClient.Create(ctx, fixtures.Resource)).Should(Succeed())
			createdResource := new(v1beta1.ApiResource)
			Eventually(func() error {
				err = k8sClient.Get(ctx, resourceLookupKey, createdResource)
				if err != nil {
					return err
				}
				if len(createdResource.Finalizers) == 0 {
					return fmt.Errorf("Resource does not have any finalizer: %s", resourceLookupKey)
				}
				return nil
			}, timeout, interval).Should(Succeed())

			Expect(k8sClient.Create(ctx, fixtures.Api)).Should(Succeed())
			createdApi := new(v1alpha1.ApiDefinition)
			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, createdApi)
			}, timeout, interval).Should(Succeed())
		})

		It("Should delete the resource once not referenced", func() {
			createdResource := new(v1beta1.ApiResource)
			Eventually(func() error {
				return k8sClient.Get(ctx, resourceLookupKey, createdResource)
			}, timeout, interval).Should(Succeed())
			Expect(k8sClient.Delete(ctx, createdResource)).Should(Succeed())

			Eventually(func() error {
				return k8sClient.Get(ctx, resourceLookupKey, createdResource)
			}, timeout, interval).Should(Succeed())

			createdApi := new(v1alpha1.ApiDefinition)
			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, createdApi)
			}, timeout, interval).Should(Succeed())

			Expect(k8sClient.Delete(ctx, createdApi)).Should(Succeed())
			Eventually(func() error {
				err := k8sClient.Get(ctx, apiLookupKey, createdApi)
				if errors.IsNotFound(err) {
					return nil
				}
				return fmt.Errorf("Should not find api: %s", apiLookupKey)
			}, timeout, interval).ShouldNot(HaveOccurred())

			Eventually(func() error {
				err := k8sClient.Get(ctx, resourceLookupKey, createdResource)
				if errors.IsNotFound(err) {
					return nil
				}
				return fmt.Errorf("Should not find resource: %s", resourceLookupKey)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})
	})
})
