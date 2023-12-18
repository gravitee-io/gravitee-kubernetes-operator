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

	"bytes"
	"fmt"

	apimErrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	xhttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal"
	"github.com/pavlo-v-chernykh/keystore-go/v4"
	coreV1 "k8s.io/api/core/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	netV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
)

const TLSCN = "httpbin.example.com"

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

			createdAPIDefinition := &v1alpha1.ApiDefinition{}
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdAPIDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			Expect(createdAPIDefinition.Name).Should(Equal(ingressFixture.Name))

			Expect(createdAPIDefinition.Spec.Proxy.VirtualHosts[0].Path).Should(Equal("/get" + fixtureGenerator.Suffix))
			Expect(createdAPIDefinition.Spec.Proxy.Groups[0].Endpoints).Should(Equal(
				[]*v2.Endpoint{
					{
						Name:   "rule01-path01",
						Target: "http://httpbin.default.svc.cluster.local:8000",
						Type:   "http",
					},
				},
			))

			By("Checking events")
			Eventually(
				getEventReasons(ingressFixture), timeout, interval,
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

			createdAPIDefinition := &v1alpha1.ApiDefinition{}
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdAPIDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			Expect(createdAPIDefinition.Spec.Proxy.VirtualHosts).Should(HaveLen(3))

			Expect(createdAPIDefinition.Spec.Proxy.VirtualHosts).Should(
				Equal([]*v2.VirtualHost{
					{
						Host: "foo.example.com",
						Path: "/ingress/foo" + fixtureGenerator.Suffix,
					},
					{
						Host: "bar.example.com",
						Path: "/ingress/bar" + fixtureGenerator.Suffix,
					},
					{
						Host: "",
						Path: "/ingress/baz" + fixtureGenerator.Suffix,
					},
				}),
			)

			Expect(createdAPIDefinition.Spec.Proxy.Groups[0].Endpoints).Should(Equal(
				[]*v2.Endpoint{
					{
						Name:   "rule01-path01",
						Target: "http://httpbin-1.default.svc.cluster.local:8080",
						Type:   "http",
					},
					{
						Name:   "rule02-path01",
						Target: "http://httpbin-2.default.svc.cluster.local:8080",
						Type:   "http",
					},
					{
						Name:   "rule03-path01",
						Target: "http://httpbin-3.default.svc.cluster.local:8080",
						Type:   "http",
					},
				},
			))

			cli := xhttp.NewClient(ctx, nil)

			By("Checking that rule with host foo.example.com is working")

			host := new(internal.Host)

			Eventually(func() error {
				base, xErr := xhttp.NewURL(internal.GatewayUrl)
				Expect(xErr).ToNot(HaveOccurred())

				rulePath := "/foo" + fixtureGenerator.Suffix

				url := base.WithPath("/ingress").WithPath(rulePath).WithPath("/hostname")

				return cli.Get(url, host, xhttp.WithHost("foo.example.com"))
			}, timeout, interval).ShouldNot(HaveOccurred())

			Expect(internal.AssertHostPrefix(host, "httpbin-1")).ToNot(HaveOccurred())

			By("Checking that rule with host bar.example.com is working")

			host = new(internal.Host)

			Eventually(func() error {

				base, xErr := xhttp.NewURL(internal.GatewayUrl)
				Expect(xErr).ToNot(HaveOccurred())

				rulePath := "/bar" + fixtureGenerator.Suffix

				url := base.WithPath("/ingress").WithPath(rulePath).WithPath("/hostname")

				return cli.Get(url, host, xhttp.WithHost("bar.example.com"))
			}, timeout, interval).ShouldNot(HaveOccurred())

			Expect(internal.AssertHostPrefix(host, "httpbin-2")).ToNot(HaveOccurred())

			By("Checking that rule with no host is working with no host")

			host = new(internal.Host)

			Eventually(func() error {
				base, xErr := xhttp.NewURL(internal.GatewayUrl)
				Expect(xErr).ToNot(HaveOccurred())

				rulePath := "/baz" + fixtureGenerator.Suffix

				url := base.WithPath("/ingress").WithPath(rulePath).WithPath("/hostname")

				return cli.Get(url, host)
			}, timeout, interval).ShouldNot(HaveOccurred())

			Expect(internal.AssertHostPrefix(host, "httpbin-3")).ToNot(HaveOccurred())

			By("Checking that rule with no host is working with unknown host")

			host = new(internal.Host)

			Eventually(func() error {
				base, xErr := xhttp.NewURL(internal.GatewayUrl)
				Expect(xErr).ToNot(HaveOccurred())

				rulePath := "/baz" + fixtureGenerator.Suffix

				url := base.WithPath("/ingress").WithPath(rulePath).WithPath("/hostname")

				return cli.Get(url, host, xhttp.WithHost("unknown.example.com"))
			}, timeout, interval).ShouldNot(HaveOccurred())

			Expect(internal.AssertHostPrefix(host, "httpbin-3")).ToNot(HaveOccurred())

			By("Checking that a 404 is returned if no rule matches, using a custom template")

			Eventually(func() error {
				base, xErr := xhttp.NewURL(internal.GatewayUrl)
				Expect(xErr).ToNot(HaveOccurred())

				rulePath := "/baz" + fixtureGenerator.Suffix

				url := base.WithPath("/ingress").WithPath(rulePath).WithPath("/hostname")

				callErr := cli.Get(url, host, xhttp.WithHost("foo.example.com"))
				nfErr := new(apimErrors.ServerError)
				if !errors.As(callErr, nfErr) {
					return internal.NewAssertionError("error", nfErr, callErr)
				}
				return internal.AssertEquals("message", "not-found", nfErr.Message)
			}, timeout, interval).ShouldNot(HaveOccurred())
		})

		It("Should create the ingress, default ApiDefinition and update GW keystore", func() {
			By("Initializing the Ingress fixture")
			fixtureGenerator := internal.NewFixtureGenerator()
			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Ingress: internal.IngressWithTLS,
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

			createdApiDefinition := &v1alpha1.ApiDefinition{}
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			expectedApiName := ingressFixture.Name
			Expect(createdApiDefinition.Name).Should(Equal(expectedApiName))

			Expect(createdApiDefinition.Spec.Proxy.VirtualHosts[0].Path).Should(Equal("/get-tls" + fixtureGenerator.Suffix))
			Expect(createdApiDefinition.Spec.Proxy.Groups[0].Endpoints).Should(Equal(
				[]*v2.Endpoint{
					{
						Name:   "rule01-path01",
						Target: "http://httpbin.default.svc.cluster.local:8000",
						Type:   "http",
					},
				},
			))

			By("Checking events")
			Eventually(
				getEventReasons(ingressFixture),
				timeout, interval,
			).Should(
				ContainElements([]string{"UpdateSucceeded", "UpdateStarted"}),
			)

			ksCredentials := &coreV1.Secret{}
			Eventually(func() error {
				ksObjectKey := types.NamespacedName{
					Namespace: namespace,
					Name:      "gw-keystore-credentials",
				}
				return k8sClient.Get(ctx, ksObjectKey, ksCredentials)
			}, timeout, interval).ShouldNot(HaveOccurred())

			ksSecret := &coreV1.Secret{}
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
						return nil
					}
				}

				return fmt.Errorf("no keyair found for %s in the keystore", TLSCN)
			}, timeout, interval).ShouldNot(HaveOccurred())

			Expect(k8sClient.Delete(ctx, ingressFixture)).Should(Succeed())
		})
	})

	Context("With API definition template", func() {

		It("Should create the ingress and use the ApiDefinition Template", func() {
			By("Initializing the Ingress fixture")
			fixtureGenerator := internal.NewFixtureGenerator()
			fixtures, err := fixtureGenerator.NewFixtures(internal.FixtureFiles{
				Api:     internal.ApiTemplateWithApiKeyPlanFile,
				Ingress: internal.IngressWithTemplateFile,
			})
			Expect(err).ToNot(HaveOccurred())

			By("Create an API definition template")

			apiDefinitionTemplate := fixtures.Api
			Expect(k8sClient.Create(ctx, apiDefinitionTemplate)).Should(Succeed())

			apiTemplateLookupKey := types.NamespacedName{Name: apiDefinitionTemplate.Name, Namespace: namespace}
			By("Expect the API Template to be ready")

			savedAPITemplate := new(v1alpha1.ApiDefinition)
			Eventually(func() error {
				return k8sClient.Get(ctx, apiTemplateLookupKey, savedAPITemplate)
			}, timeout, interval).Should(Succeed())

			ingressFixture := fixtures.Ingress
			ingressFixture.Annotations[keys.IngressTemplateAnnotation] = apiTemplateLookupKey.Name
			ingressLookupKey := types.NamespacedName{Name: ingressFixture.Name, Namespace: namespace}

			By("Creating an Ingress and the ApiDefinition")
			Expect(k8sClient.Create(ctx, ingressFixture)).Should(Succeed())

			By("Getting created resource and expect to find it")
			createdIngress := &netV1.Ingress{}
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdIngress)
			}, timeout, interval).ShouldNot(HaveOccurred())

			createdApiDefinition := &v1alpha1.ApiDefinition{}
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressLookupKey, createdApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			expectedApiName := ingressFixture.Name
			Expect(createdApiDefinition.Name).Should(Equal(expectedApiName))

			Expect(len(createdApiDefinition.Spec.Plans)).Should(Equal(1))
			Expect(createdApiDefinition.Spec.Plans[0].Security).Should(Equal("API_KEY"))

			By("Checking events")
			Eventually(
				getEventReasons(ingressFixture), timeout, interval,
			).Should(
				ContainElements([]string{"UpdateSucceeded", "UpdateStarted"}),
			)
		})

	})
})
