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

package test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1beta1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"
)

var _ = Describe("API Resource Controller", func() {

	Context("Update an API Resource", func() {

		var apiLookupKey types.NamespacedName
		var contextLookupKey types.NamespacedName
		var resourceLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Creating an API and an API resource with a management context")

			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:      internal.BasicApiFile,
				Context:  internal.ClusterContextFile,
				Resource: internal.ApiResourceCacheFile,
			})

			Expect(err).ToNot(HaveOccurred())

			managementContext := fixtures.Context
			Expect(k8sClient.Create(ctx, managementContext)).Should(Succeed())

			contextLookupKey = types.NamespacedName{Name: managementContext.Name, Namespace: namespace}
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, managementContext)
			}, timeout, interval).Should(Succeed())

			Expect(k8sClient.Create(ctx, fixtures.Resource)).Should(Succeed())

			Expect(k8sClient.Create(ctx, fixtures.Api)).Should(Succeed())

			apiLookupKey = types.NamespacedName{Name: fixtures.Api.Name, Namespace: namespace}
			resourceLookupKey = types.NamespacedName{Name: fixtures.Resource.Name, Namespace: namespace}
		})

		It("Should update the API definition on resource update", func() {
			createdResource := new(v1beta1.ApiResource)

			Eventually(func() error {
				return k8sClient.Get(ctx, resourceLookupKey, createdResource)
			}, timeout, interval).Should(Succeed())

			Expect(createdResource.Spec.Enabled).To(BeTrue())

			createdApi := new(v1alpha1.ApiDefinition)

			Eventually(func() error {
				err := k8sClient.Get(ctx, apiLookupKey, createdApi)
				if err != nil {
					return err
				}

				return internal.AssertStatusMatches(createdApi, v1alpha1.ApiDefinitionStatus{
					EnvID:                        "DEFAULT",
					OrgID:                        "DEFAULT",
					CrossID:                      createdApi.GetOrGenerateCrossID(),
					ID:                           createdApi.PickID(),
					Status:                       v1alpha1.ProcessingStatusCompleted,
					DeprecatedStatus:             v1alpha1.ProcessingStatusCompleted,
					State:                        "STARTED",
					ObservedGeneration:           1,
					DeprecatedObservedGeneration: 1,
				})
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Updating the API resource, expecting the API definition resource to be updated")

			updatedResource := createdResource.DeepCopy()
			updatedResource.Spec.Enabled = false

			Eventually(func() error {
				update := new(v1beta1.ApiResource)
				if err := k8sClient.Get(ctx, resourceLookupKey, update); err != nil {
					return err
				}
				updatedResource.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

			apimClient, err := internal.NewAPIM(ctx)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() error {
				api, apiErr := apimClient.APIs.GetByID(createdApi.Status.ID)
				if apiErr != nil {
					return apiErr
				}

				if api.Resources == nil || len(api.Resources) == 0 {
					return errors.New("no resources found")
				}

				if api.Resources[0].Enabled {
					return errors.New("api not updated")
				}

				return nil
			}, timeout, interval).ShouldNot(HaveOccurred())
		})
	})
})
