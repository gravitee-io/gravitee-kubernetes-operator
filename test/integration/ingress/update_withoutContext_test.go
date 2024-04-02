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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"

	xhttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
)

var _ = Describe("Update", labels.WithoutContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	httpCli := xhttp.NewNoAuthClient(ctx)

	It("should update backing api definition", func() {
		fixtures := fixture.Builder().
			WithIngress(constants.IngressWithoutTemplateFile).
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
		Expect(proxy.VirtualHosts).ToNot(BeEmpty())
		endpoint := constants.BuildAPIEndpoint(apiDef)
		Expect(assert.PathEquals(endpoint, "/httpbin"+fixtures.GetGeneratedSuffix())).To(Succeed())

		By("updating ingress path")

		updated := fixtures.Ingress.DeepCopy()
		updated.Spec.Rules[0].HTTP.Paths[0].Path = "/" + fixtures.GetGeneratedSuffix()[1:]

		Expect(manager.UpdateSafely(updated)).To(Succeed())

		By("expecting API proxy to match updated rule")

		Eventually(func() error {
			latest, err := manager.GetLatest(apiDef)
			if err != nil {
				return err
			}
			endpoint = constants.BuildAPIEndpoint(latest)
			return assert.PathEquals(endpoint, "/"+fixtures.GetGeneratedSuffix()[1:])
		}, timeout, interval).Should(Succeed())

		By("checking access to backend service")

		Eventually(func() error {
			url, err := xhttp.NewURL(endpoint)
			Expect(err).ToNot(HaveOccurred())
			return httpCli.Get(url.WithPath("hostname").String(), new(Host), xhttp.WithHost("httpbin.example.com"))
		}, timeout, interval).ShouldNot(HaveOccurred())
	})
})
