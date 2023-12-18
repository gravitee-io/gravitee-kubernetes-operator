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
	"bytes"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pavlo-v-chernykh/keystore-go/v4"
	core "k8s.io/api/core/v1"
	netV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Deleting an ingress", func() {
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

		It("Should delete the Ingress and the default ApiDefinition", func() {
			By("Creating an Ingress and the default ApiDefinition")

			Expect(k8sClient.Create(ctx, ingressFixture)).Should(Succeed())

			By("Getting created resource and expect to find it")

			createdIngress := new(netV1.Ingress)
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdIngress)
			}, timeout, interval).ShouldNot(HaveOccurred())

			createdAPIDefinition := new(v1alpha1.ApiDefinition)
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdAPIDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Delete the Ingress")
			Expect(k8sClient.Delete(ctx, createdIngress)).ToNot(HaveOccurred())

			By("Checking the Ingress has been deleted")
			ing := &netV1.Ingress{}
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, ing)
			}, timeout, interval).ShouldNot(Succeed())

			api := &v1alpha1.ApiDefinition{}
			Eventually(func() error {
				err := k8sClient.Get(ctx, ingressLookupKey, api)
				return err
			}, timeout, interval).ShouldNot(Succeed())

			By("Checking events")
			Eventually(
				getEventReasons(ingressFixture), timeout, interval,
			).Should(
				ContainElements([]string{"DeleteSucceeded", "DeleteStarted"}),
			)
		})
	})

	Context("With api definition template", func() {
		var apiDefinitionTemplate *v1alpha1.ApiDefinition
		var ingressFixture *netV1.Ingress
		var createdIngress *netV1.Ingress
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

			apiTemplateLookupKey = types.NamespacedName{Name: apiDefinitionTemplate.Name, Namespace: namespace}
			By("Expect the API Template to be ready")

			savedAPITemplate := new(v1alpha1.ApiDefinition)
			Eventually(func() error {
				return k8sClient.Get(ctx, apiTemplateLookupKey, savedAPITemplate)
			}, timeout, interval).Should(Succeed())

			ingressFixture = fixtures.Ingress
			ingressFixture.Annotations[keys.IngressTemplateAnnotation] = apiTemplateLookupKey.Name
			ingressLookupKey = types.NamespacedName{Name: ingressFixture.Name, Namespace: namespace}

			By("Creating an Ingress and the ApiDefinition")
			Expect(k8sClient.Create(ctx, ingressFixture)).Should(Succeed())

			By("Getting created ingress resource and expected to find it")
			createdIngress = &netV1.Ingress{}
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdIngress)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Getting created api definition and expected to find it")
			createdAPIDefinition := &v1alpha1.ApiDefinition{}
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdAPIDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			expectedAPIName := ingressFixture.Name
			Expect(createdAPIDefinition.Name).Should(Equal(expectedAPIName))

			Expect(len(createdAPIDefinition.Spec.Plans)).Should(Equal(1))
			Expect(createdAPIDefinition.Spec.Plans[0].Security).Should(Equal("API_KEY"))
		})

		When("API template has a reference to an exiting ingress", func() {
			It("Should NOT delete the API Template", func() {
				By("Deleting the API definition template")
				Expect(k8sClient.Delete(ctx, apiDefinitionTemplate)).ToNot(HaveOccurred())
				Consistently(func() error {
					By("Expect the API Template to be still available")

					savedAPITemplate := new(v1alpha1.ApiDefinition)
					return k8sClient.Get(ctx, apiTemplateLookupKey, savedAPITemplate)
				}, timeout).Should(Succeed())
			})
		})

		When("API template does not have a reference to an exiting ingress", func() {
			It("Should delete the Ingress and the ApiDefinition", func() {

				By("Deleting the Ingress")

				Expect(k8sClient.Delete(ctx, createdIngress)).ToNot(HaveOccurred())

				By("Deleting the API definition template")

				Eventually(func() error {
					savedAPITemplate := new(v1alpha1.ApiDefinition)
					if err := k8sClient.Get(ctx, apiTemplateLookupKey, savedAPITemplate); err != nil {
						return err
					}

					return k8sClient.Delete(ctx, savedAPITemplate)
				}, timeout, interval).Should(Succeed())

				By("Checking events")
				Eventually(
					getEventReasons(ingressFixture), timeout, interval,
				).Should(
					ContainElements([]string{"DeleteSucceeded", "DeleteStarted"}),
				)
			})
		})
	})

	Context("With API definition template", func() {
		It("Should delete the Ingress, the default ApiDefinition and keypair in GW keystore", func() {
			By("Initializing the Ingress fixture")
			fixtureGenerator := internal.NewFixtureGenerator()

			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Ingress: internal.IngressWithTLS,
			})

			Expect(err).ToNot(HaveOccurred())

			ingressFixture := fixtures.Ingress
			ingressLookupKey := types.NamespacedName{Name: ingressFixture.Name, Namespace: namespace}

			By("Creating an Ingress and the default ApiDefinition")

			Expect(k8sClient.Create(ctx, ingressFixture)).Should(Succeed())

			By("Getting created resource and expect to find it")

			createdIngress := new(netV1.Ingress)
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdIngress)
			}, timeout, interval).ShouldNot(HaveOccurred())

			createdApiDefinition := new(v1alpha1.ApiDefinition)
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Delete the Ingress")
			Expect(k8sClient.Delete(ctx, createdIngress)).ToNot(HaveOccurred())

			By("Checking the Ingress has been deleted")
			ing := &netV1.Ingress{}
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, ing)
			}, timeout, interval).ShouldNot(Succeed())

			api := &v1alpha1.ApiDefinition{}
			Eventually(func() error {
				err = k8sClient.Get(ctx, ingressLookupKey, api)
				return err
			}, timeout, interval).ShouldNot(Succeed())

			By("Checking events")
			Eventually(
				getEventReasons(ingressFixture), timeout, interval,
			).Should(
				ContainElements([]string{"DeleteSucceeded", "DeleteStarted"}),
			)

			ksCredentials := &core.Secret{}
			Eventually(func() error {
				ksObjectKey := types.NamespacedName{
					Namespace: namespace,
					Name:      "gw-keystore-credentials",
				}
				return k8sClient.Get(ctx, ksObjectKey, ksCredentials)
			}, timeout, interval).ShouldNot(HaveOccurred())

			ksSecret := &core.Secret{}
			Eventually(func() error {
				ksObjectKey := types.NamespacedName{
					Namespace: namespace,
					Name:      string(ksCredentials.Data["name"]),
				}
				return k8sClient.Get(ctx, ksObjectKey, ksSecret)
			}, timeout, interval).ShouldNot(HaveOccurred())

			Eventually(func() error {
				data := ksSecret.Data["keystore"]
				if data == nil {
					return fmt.Errorf("gateway keystore not found")
				}

				ks := keystore.New()
				err = ks.Load(bytes.NewReader(data), []byte("changeme"))
				if err != nil {
					return fmt.Errorf("can't load the gateway keystore")
				}

				for _, a := range ks.Aliases() {
					if a == TLSCN {
						return fmt.Errorf("tls keypair shouldn't be in the gateway keystore")
					}
				}

				return nil
			}, timeout, interval).ShouldNot(HaveOccurred())
		})

	})

})
