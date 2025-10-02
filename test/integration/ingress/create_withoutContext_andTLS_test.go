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

package ingress

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Create", labels.WithoutContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should update pem registry", func() {
		fixtures := fixture.Builder().
			AddConfigMap(constants.IngressPEMRegistry).
			AddSecret(constants.IngressWithTLSSecretFile).
			WithIngress(constants.IngressWithTLS).
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
		Expect(proxy.Groups).ToNot(BeEmpty())
		Expect(proxy.Groups[0].Endpoints).ToNot(BeEmpty())

		path := proxy.VirtualHosts[0].Path
		expectedPath := "/secured" + fixtures.GetGeneratedSuffix()
		Expect(path).To(Equal(expectedPath))

		endpoints := proxy.Groups[0].Endpoints
		Expect(endpoints).To(Equal(
			[]*v2.Endpoint{
				{
					Name:   toPointer("rule01-path01"),
					Target: toPointer("http://httpbin-1.default.svc.cluster.local:8080"),
				},
			},
		))

		By("checking pem registry")

		Eventually(func() error {
			cm := fixtures.ConfigMaps[0]
			err := manager.GetLatest(ctx, cm)
			if err != nil {
				return err
			}
			return assert.MapContaining(
				cm.Data,
				fixtures.GetIngressPEMRegistryKey(),
				fmt.Sprintf(`["%s/httpbin.example.com"]`, constants.Namespace),
			)
		}, timeout, interval).Should(Succeed())
	})
})
