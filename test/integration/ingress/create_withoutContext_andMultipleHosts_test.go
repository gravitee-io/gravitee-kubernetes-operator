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

package ingress_test

import (
	"context"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	iErrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"

	xhttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
)

var _ = Describe("Create", labels.WithoutContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	httpCli := xhttp.NewNoAuthClient(ctx)

	It("should expose backend services", func() {
		env.Config.CMTemplate404NS = constants.Namespace
		env.Config.CMTemplate404Name = "template-404"

		fixtures := fixture.Builder().
			AddConfigMap(constants.Ingress404ResponseTemplate).
			WithIngress(constants.IngressWithMultipleHosts).
			Build().
			Apply()

		By("expecting to find an API definition created for ingress")

		apiDef := new(v1alpha1.ApiDefinition)

		Eventually(func() error {
			return manager.Client().Get(ctx, types.NamespacedName{
				Namespace: fixtures.Ingress.Namespace,
				Name:      fixtures.Ingress.Name,
			}, apiDef)
		}, timeout, interval).Should(Succeed())

		By("expecting API proxy to match ingress rule")

		proxy := apiDef.Spec.Proxy
		Expect(proxy.VirtualHosts).To(HaveLen(3))
		Expect(proxy.Groups).ToNot(BeEmpty())
		Expect(proxy.Groups[0].Endpoints).To(HaveLen(3))

		Expect(proxy.VirtualHosts).To(Equal(
			[]*v2.VirtualHost{
				{
					Host: "foo.example.com",
					Path: "/ingress/foo" + fixtures.GetGeneratedSuffix(),
				},
				{
					Host: "bar.example.com",
					Path: "/ingress/bar" + fixtures.GetGeneratedSuffix(),
				},
				{
					Host: "",
					Path: "/ingress/baz" + fixtures.GetGeneratedSuffix(),
				},
			},
		))

		endpoints := proxy.Groups[0].Endpoints
		Expect(endpoints).To(Equal(
			[]*v2.Endpoint{
				{
					Name:   "rule01-path01",
					Target: "http://httpbin-1.default.svc.cluster.local:8080",
				},
				{
					Name:   "rule02-path01",
					Target: "http://httpbin-2.default.svc.cluster.local:8080",
				},
				{
					Name:   "rule03-path01",
					Target: "http://httpbin-3.default.svc.cluster.local:8080",
				},
			},
		))

		By("checking routing to httpbin-1")

		host := new(Host)
		url := fmt.Sprintf("%s/ingress/foo%s/hostname", constants.GatewayUrl, fixtures.GetGeneratedSuffix())

		Eventually(func() error {
			return httpCli.Get(url, host, xhttp.WithHost("foo.example.com"))
		}, timeout, interval).ShouldNot(HaveOccurred())

		Expect(assert.StrStartsWith(host.Name, "httpbin-1")).To(Succeed())

		By("checking routing to httpbin-2")

		host = new(Host)
		url = fmt.Sprintf("%s/ingress/bar%s/hostname", constants.GatewayUrl, fixtures.GetGeneratedSuffix())

		Eventually(func() error {
			return httpCli.Get(url, host, xhttp.WithHost("bar.example.com"))
		}, timeout, interval).ShouldNot(HaveOccurred())

		Expect(assert.StrStartsWith(host.Name, "httpbin-2")).To(Succeed())

		By("checking routing to httpbin-3")

		host = new(Host)
		url = fmt.Sprintf("%s/ingress/baz%s/hostname", constants.GatewayUrl, fixtures.GetGeneratedSuffix())

		Eventually(func() error {
			return httpCli.Get(url, host)
		}, timeout, interval).ShouldNot(HaveOccurred())

		Expect(assert.StrStartsWith(host.Name, "httpbin-3")).To(Succeed())

		By("checking error routing using response template")

		Eventually(func() error {
			callErr := httpCli.Get(url, host, xhttp.WithHost("foo.example.com"))
			nfErr := new(iErrors.ServerError)
			if !errors.As(callErr, nfErr) {
				return assert.Equals("error", nfErr, callErr)
			}
			return assert.Equals("message", "not-found-test", nfErr.Message)
		}, timeout, interval).ShouldNot(HaveOccurred())
	})
})
