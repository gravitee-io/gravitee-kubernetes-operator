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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	netV1 "k8s.io/api/networking/v1"

	"k8s.io/apimachinery/pkg/api/errors"
)

var _ = Describe("Delete", labels.WithoutContext, func() {

	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should delete ingress and update PEM registry", func() {
		fixtures := fixture.Builder().
			AddConfigMap(constants.IngressPEMRegistry).
			AddSecret(constants.IngressWithTLSSecretFile).
			WithIngress(constants.IngressWithTLS).
			Build().
			Apply()

		By("deleting the ingress")

		Expect(manager.Delete(fixtures.Ingress.DeepCopy())).To(Succeed())

		By("expecting ingress and API definition to have been deleted")

		Eventually(func() error {
			err := manager.Client().Get(ctx, types.NamespacedName{
				Namespace: fixtures.Ingress.Namespace,
				Name:      fixtures.Ingress.Name,
			}, new(netV1.Ingress))
			if !errors.IsNotFound(err) {
				return assert.Equals("error", "[NOT FOUND]", err)
			}
			return nil
		}, timeout, interval).Should(Succeed())

		Eventually(func() error {
			err := manager.Client().Get(ctx, types.NamespacedName{
				Namespace: fixtures.Ingress.Namespace,
				Name:      fixtures.Ingress.Name,
			}, new(v1alpha1.ApiDefinition))
			if !errors.IsNotFound(err) {
				return assert.Equals("error", "[NOT FOUND]", err)
			}
			return nil
		}, timeout, interval).Should(Succeed())

		By("checking pem registry")

		Eventually(func() error {
			cm, err := manager.GetLatest(fixtures.ConfigMaps[0])
			if err != nil {
				return err
			}
			return assert.MapNotContaining(cm.Data, fixtures.GetIngressPEMRegistryKey())
		}, timeout, interval).Should(Succeed())

		By("expecting ingress events to have been emitted")

		assert.EventsEmitted(fixtures.Ingress, "DeleteSucceeded", "DeleteStarted")
	})
})
