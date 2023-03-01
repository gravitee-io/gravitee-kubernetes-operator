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
	xhttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	netV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Creating an ingress", func() {
	Context("Without api definition template", func() {
		var ingressFixture *netV1.Ingress
		var ingressLookupKey types.NamespacedName

		It("Should create the ingress and use the default ApiDefinition", func() {
			By("Initializing the Ingress fixture")
			fixtureGenerator := internal.NewFixtureGenerator()
			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Ingress: internal.IngressWithoutTemplateFile,
			})
			Expect(err).ToNot(HaveOccurred())
			ingressFixture = fixtures.Ingress
			ingressLookupKey = types.NamespacedName{Name: ingressFixture.Name, Namespace: namespace}

			By("Creating an Ingress and the default ApiDefinition")
			Expect(k8sClient.Create(ctx, ingressFixture)).Should(Succeed())

			By("Getting created resource and expect to find it")
			createdIngress := &netV1.Ingress{}
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdIngress)
			}, timeout, interval).ShouldNot(HaveOccurred())

			createdApiDefinition := &gio.ApiDefinition{}
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			expectedApiName := ingressFixture.Name
			Expect(createdApiDefinition.Name).Should(Equal(expectedApiName))

			Expect(createdApiDefinition.Spec.Proxy.VirtualHosts[0].Path).Should(Equal("/get"))
			Expect(createdApiDefinition.Spec.Proxy.Groups[0].Endpoints).Should(Equal(
				[]*model.HttpEndpoint{
					{
						Name:   "httpbin.example.com/get",
						Target: "http://httpbin.default.svc.cluster.local:8000",
					},
				},
			))

			By("Checking events")
			Expect(
				getEventsReason(ingressFixture.GetNamespace(), ingressFixture.GetName()),
			).Should(
				ContainElements([]string{"UpdateSucceeded", "UpdateStarted"}),
			)
		})

		It("Should create the ingress and the api definition with multiple hosts", func() {
			By("Initializing the Ingress fixture")
			fixtureGenerator := internal.NewFixtureGenerator()
			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Ingress: internal.IngressWithMultipleHosts,
			})
			Expect(err).ToNot(HaveOccurred())
			ingressFixture = fixtures.Ingress
			ingressLookupKey = types.NamespacedName{Name: ingressFixture.Name, Namespace: namespace}

			By("Creating the ingress with multiple hosts")
			Expect(k8sClient.Create(ctx, ingressFixture)).Should(Succeed())

			By("Getting created resource and expect to find it")
			createdIngress := &netV1.Ingress{}
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdIngress)
			}, timeout, interval).ShouldNot(HaveOccurred())

			createdApiDefinition := &gio.ApiDefinition{}
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			Expect(createdApiDefinition.Spec.Proxy.VirtualHosts).Should(HaveLen(3))

			Expect(createdApiDefinition.Spec.Proxy.VirtualHosts).Should(
				Equal([]*model.VirtualHost{
					{
						Host: "foo.example.com",
						Path: "/ingress/foo",
					},
					{
						Host: "bar.example.com",
						Path: "/ingress/bar",
					},
					{
						Host: "",
						Path: "/ingress/baz",
					},
				}),
			)

			Expect(createdApiDefinition.Spec.Proxy.Groups[0].Endpoints).Should(Equal(
				[]*model.HttpEndpoint{
					{
						Name:   "foo.example.com/ingress/foo",
						Target: "http://svc-1.default.svc.cluster.local:8080",
					},
					{
						Name:   "bar.example.com/ingress/bar",
						Target: "http://svc-2.default.svc.cluster.local:8080",
					},
					{
						Name:   "ingress/baz",
						Target: "http://svc-3.default.svc.cluster.local:8080",
					},
				},
			))

			cli := xhttp.NewClient(ctx, nil)

			By("Checking that rule with host foo.example.com is working")

			host := new(internal.Host)

			Eventually(func() error {
				url := internal.GatewayUrl + "/ingress/foo/hostname"
				return cli.Get(url, host, xhttp.WithHost("foo.example.com"))
			}, timeout, interval).ShouldNot(HaveOccurred())

			Expect(internal.AssertHostPrefix(host, "svc-1")).ToNot(HaveOccurred())

			By("Checking that rule with host bar.example.com is working")

			host = new(internal.Host)

			Eventually(func() error {
				url := internal.GatewayUrl + "/ingress/bar/hostname"
				return cli.Get(url, host, xhttp.WithHost("bar.example.com"))
			}, timeout, interval).ShouldNot(HaveOccurred())

			Expect(internal.AssertHostPrefix(host, "svc-2")).ToNot(HaveOccurred())

			By("Checking that rule with no host is working")

			host = new(internal.Host)

			Eventually(func() error {
				url := internal.GatewayUrl + "/ingress/baz/hostname"
				return cli.Get(url, host)
			}, timeout, interval).ShouldNot(HaveOccurred())

			Expect(internal.AssertHostPrefix(host, "svc-3")).ToNot(HaveOccurred())
		})
	})
})
