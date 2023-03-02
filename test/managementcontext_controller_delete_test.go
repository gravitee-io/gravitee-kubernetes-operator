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
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Deleting a management context", func() {
	Context("Not linked to an api definition", func() {
		var contextFixture *gio.ManagementContext
		var contextLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Initializing the management context fixture")
			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Context: internal.ContextWithCredentialsFile,
			})

			Expect(err).ToNot(HaveOccurred())

			contextFixture = fixtures.Context
			contextLookupKey = types.NamespacedName{Name: contextFixture.Name, Namespace: namespace}
		})

		It("Should delete the management context", func() {
			By("Creating a new management context")

			Expect(k8sClient.Create(ctx, contextFixture)).Should(Succeed())

			By("Getting created resource and expect to find it")

			createdContext := new(gio.ManagementContext)
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, createdContext)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Deleting the management context")
			Expect(k8sClient.Delete(ctx, createdContext)).ToNot(HaveOccurred())

			By("Checking the management context has been deleted")
			context := &gio.ManagementContext{}
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, context)
			}, timeout, interval).ShouldNot(Succeed())
		})
	})

	Context("Linked to an api definition", func() {
		var contextFixture *gio.ManagementContext
		var contextLookupKey types.NamespacedName
		var apiFixture *gio.ApiDefinition
		var apiLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Initializing the management context fixture and api definition")
			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:     internal.ApiWithContextFile,
				Context: internal.ContextWithCredentialsFile,
			})

			Expect(err).ToNot(HaveOccurred())

			contextFixture = fixtures.Context
			contextLookupKey = types.NamespacedName{Name: contextFixture.Name, Namespace: namespace}

			apiFixture = fixtures.Api
			apiLookupKey = types.NamespacedName{Name: apiFixture.Name, Namespace: namespace}
		})

		AfterEach(func() {
			Expect(k8sClient.Delete(ctx, apiFixture)).Should(Succeed())
		})

		It("Should not delete the management context", func() {
			By("Creating a new management context")
			Expect(k8sClient.Create(ctx, contextFixture)).Should(Succeed())

			By("Getting created resource and expect to find it")
			createdContext := new(gio.ManagementContext)
			Eventually(func() error {
				if err := k8sClient.Get(ctx, contextLookupKey, createdContext); err != nil {
					return err
				}
				return internal.AssertEquals("finalizer.length", 1, len(createdContext.ObjectMeta.Finalizers))
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Creating a api definition")
			Expect(k8sClient.Create(ctx, apiFixture)).Should(Succeed())

			By("Getting created resource and expect to find it")
			createdApi := new(gio.ApiDefinition)
			Consistently(func() error { // just to let gko have time to configure API definition
				return k8sClient.Get(ctx, apiLookupKey, createdApi)
			}, timeout/10, interval).Should(Succeed())

			By("Trying to delete the management context")
			Expect(k8sClient.Delete(ctx, createdContext)).ToNot(HaveOccurred())

			By("Checking the management context has not been deleted")
			context := &gio.ManagementContext{}
			Eventually(func() error {
				return k8sClient.Get(ctx, contextLookupKey, context)
			}, timeout, interval).Should(Succeed())
		})
	})
})
