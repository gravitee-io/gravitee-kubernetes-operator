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

			createdApiDefinition := new(gio.ApiDefinition)
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

			api := &gio.ApiDefinition{}
			Eventually(func() error {
				err := k8sClient.Get(ctx, ingressLookupKey, api)
				return err
			}, timeout, interval).ShouldNot(Succeed())

			By("Checking events")
			Expect(
				getEventsReason(ingressFixture.GetNamespace(), ingressFixture.GetName()),
			).Should(
				ContainElements([]string{"DeleteSucceeded", "DeleteStarted"}),
			)
		})
	})
})
