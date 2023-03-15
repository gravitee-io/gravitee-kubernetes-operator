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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	netV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Updating an ingress", func() {
	Context("Without api definition template", func() {
		var ingressFixture *netV1.Ingress
		var ingressLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Initializing the Ingress fixture")
			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Ingress: internal.IngressWithoutTemplateFile,
			})

			Expect(err).ToNot(HaveOccurred())

			ingressFixture = fixtures.Ingress
			ingressLookupKey = types.NamespacedName{Name: ingressFixture.Name, Namespace: namespace}
		})

		It("Should update an Ingress and the default ApiDefinition", func() {
			By("Creating an Ingress and the default ApiDefinition")

			Expect(k8sClient.Create(ctx, ingressFixture)).Should(Succeed())

			By("Getting created resource and expect to find it")
			createdIngress := new(netV1.Ingress)
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdIngress)
			}, timeout, interval).ShouldNot(HaveOccurred())

			createdAPIDefinition := new(gio.ApiDefinition)
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdAPIDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Updating the Ingress")
			fooPath := "/foo"

			updatedIngress := createdIngress.DeepCopy()
			updatedIngress.Spec.Rules[0].HTTP.Paths[0].Path = fooPath

			Eventually(func() error {
				update := new(netV1.Ingress)
				if err := k8sClient.Get(ctx, ingressLookupKey, update); err != nil {
					return err
				}
				updatedIngress.Spec.DeepCopyInto(&update.Spec)
				return k8sClient.Update(ctx, update)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Checking the Ingress and ApiDefinition values")
			ingressWithUpdatedPath := new(netV1.Ingress)
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, ingressWithUpdatedPath)
			}, timeout, interval).ShouldNot(HaveOccurred())
			Expect(ingressWithUpdatedPath.Spec.Rules[0].HTTP.Paths[0].Path).To(Equal(fooPath))

			Eventually(func() bool {
				apiDefinitionWithUpdatedPath := new(gio.ApiDefinition)
				Eventually(func() error {
					return k8sClient.Get(ctx, ingressLookupKey, apiDefinitionWithUpdatedPath)
				}, timeout, interval).ShouldNot(HaveOccurred())
				return apiDefinitionWithUpdatedPath.Spec.Proxy.VirtualHosts[0].Path == fooPath
			}).ShouldNot(Equal(false))

			By("Checking events")
			Expect(
				getEventsReason(ingressFixture.GetNamespace(), ingressFixture.GetName()),
			).Should(
				ContainElements([]string{"UpdateSucceeded", "UpdateStarted"}),
			)
		})
	})

	Context("With API definition template", func() {
		var apiDefinitionTemplate *gio.ApiDefinition
		var ingressFixture *netV1.Ingress
		var ingressLookupKey types.NamespacedName
		var apiTemplateLookupKey types.NamespacedName

		BeforeEach(func() {
			By("Initializing the Ingress fixture")
			fixtureGenerator := internal.NewFixtureGenerator()
			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:     internal.ApiTemplateWithApiKeyPlanFile,
				Ingress: internal.IngressWithTemplateFile,
			})
			Expect(err).ToNot(HaveOccurred())

			By("Create an API definition template")

			apiDefinitionTemplate = fixtures.Api
			Expect(k8sClient.Create(ctx, apiDefinitionTemplate)).Should(Succeed())

			apiTemplateLookupKey = types.NamespacedName{Namespace: namespace, Name: apiDefinitionTemplate.Name}

			ingressFixture = fixtures.Ingress
			ingressFixture.Annotations[keys.IngressTemplateAnnotation] = apiDefinitionTemplate.Name
			ingressLookupKey = types.NamespacedName{Name: ingressFixture.Name, Namespace: namespace}

			By("Creating an Ingress and the default ApiDefinition")
			Expect(k8sClient.Create(ctx, ingressFixture)).Should(Succeed())
		})

		When("Updating API definition template", func() {
			It("it should update the ingress and final api definition", func() {
				By("Getting created resource and expect to find it")

				createdAPIDefinition := &gio.ApiDefinition{}
				Eventually(func() error {
					return k8sClient.Get(ctx, ingressLookupKey, createdAPIDefinition)
				}, timeout, interval).ShouldNot(HaveOccurred())

				Expect(len(createdAPIDefinition.Spec.Plans)).Should(Equal(1))
				Expect(createdAPIDefinition.Spec.Plans[0].Security).Should(Equal("API_KEY"))

				currentAPITemplate := new(gio.ApiDefinition)
				Eventually(func() error {
					return k8sClient.Get(ctx, apiTemplateLookupKey, currentAPITemplate)
				}).Should(Succeed())

				By("update api template")

				updatedAPITemplate := currentAPITemplate.DeepCopy()
				updatedAPITemplate.Spec.Plans = append(updatedAPITemplate.Spec.Plans, &model.Plan{
					Name:     "Default keyless plan",
					Security: "KEY_LESS",
					Status:   "PUBLISHED",
				})

				Eventually(func() error {
					update := new(gio.ApiDefinition)
					if err := k8sClient.Get(ctx, apiTemplateLookupKey, update); err != nil {
						return err
					}
					updatedAPITemplate.Spec.DeepCopyInto(&update.Spec)
					return k8sClient.Update(ctx, update)
				}).ShouldNot(HaveOccurred())

				updateAPIDefinition := &gio.ApiDefinition{}
				Consistently(func() error {
					return k8sClient.Get(ctx, ingressLookupKey, updateAPIDefinition)
				}, timeout/5, interval).ShouldNot(HaveOccurred())

				Expect(len(updateAPIDefinition.Spec.Plans)).Should(Equal(2))
				Expect(updateAPIDefinition.Spec.Plans[0].Security).Should(Equal("API_KEY"))
				Expect(updateAPIDefinition.Spec.Plans[1].Security).Should(Equal("KEY_LESS"))

				By("Checking events")
				Expect(
					getEventsReason(ingressLookupKey.Namespace, ingressLookupKey.Name),
				).Should(
					ContainElements([]string{"UpdateSucceeded", "UpdateStarted"}),
				)
			})
		})

		When("Updating the ingress", func() {
			It("it should update the api definition", func() {
				By("Getting created resource and expect to find it")
				createdIngress := new(netV1.Ingress)
				Eventually(func() error {
					return k8sClient.Get(ctx, ingressLookupKey, createdIngress)
				}, timeout, interval).ShouldNot(HaveOccurred())

				createdAPIDefinition := new(gio.ApiDefinition)
				Eventually(func() error {
					return k8sClient.Get(ctx, ingressLookupKey, createdAPIDefinition)
				}, timeout, interval).ShouldNot(HaveOccurred())

				By("Updating the Ingress")
				fooPath := "/foo-tls"

				updatedIngress := createdIngress.DeepCopy()
				updatedIngress.Spec.Rules[0].HTTP.Paths[0].Path = fooPath

				Eventually(func() error {
					update := new(netV1.Ingress)
					if err := k8sClient.Get(ctx, ingressLookupKey, update); err != nil {
						return err
					}
					updatedIngress.Spec.DeepCopyInto(&update.Spec)
					return k8sClient.Update(ctx, update)
				}, timeout, interval).ShouldNot(HaveOccurred())

				By("Checking the Ingress and ApiDefinition values")
				ingressWithUpdatedPath := new(netV1.Ingress)
				Eventually(func() error {
					if err := k8sClient.Get(ctx, ingressLookupKey, ingressWithUpdatedPath); err != nil {
						return err
					}

					if ingressWithUpdatedPath.Spec.Rules[0].HTTP.Paths[0].Path != fooPath {
						return fmt.Errorf("the ingress path hasn't been updated")
					}

					return nil
				}, timeout, interval).ShouldNot(HaveOccurred())

				Eventually(func() bool {
					apiDefinitionWithUpdatedPath := new(gio.ApiDefinition)
					Eventually(func() error {
						return k8sClient.Get(ctx, ingressLookupKey, apiDefinitionWithUpdatedPath)
					}, timeout, interval).ShouldNot(HaveOccurred())
					return apiDefinitionWithUpdatedPath.Spec.Proxy.VirtualHosts[0].Path == fooPath
				}).ShouldNot(Equal(false))

				By("Checking events")
				Expect(
					getEventsReason(ingressLookupKey.Namespace, ingressLookupKey.Name),
				).Should(
					ContainElements([]string{"UpdateSucceeded", "UpdateStarted"}),
				)
			})
		})
	})
})
