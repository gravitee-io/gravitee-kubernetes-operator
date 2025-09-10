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

package managementcontext

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"k8s.io/apimachinery/pkg/api/errors"
)

var _ = Describe("Delete", labels.WithContext, func() {
	timeout := constants.EventualTimeout

	interval := constants.Interval

	ctx := context.Background()

	It("should delete secret finalizer when there is no reference anymore", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithContext(constants.ContextWithSecretFile).
			Build()

		fixtures.Secrets[0].Name += random.GetSuffix()
		namespacedName := refs.NewNamespacedName(fixtures.Secrets[0].Namespace, fixtures.Secrets[0].Name)
		fixtures.Context.Spec.Auth.SecretRef = &namespacedName

		fixtures.Apply()

		Eventually(func() error {
			return assert.Equals("annotation", fmt.Sprintf("%s/%s",
				fixtures.Secrets[0].Namespace, fixtures.Secrets[0].Name),
				fixtures.Context.Annotations["gravitee.io/last-secret-ref"])
		}, timeout, interval).Should(Succeed())

		By("updating management context and removing secret reference")
		fixtures.Context.Spec.Auth = nil

		Expect(manager.Client().Update(ctx, fixtures.Context)).To(Succeed())

		By("deleting secret")

		Expect(manager.Client().Delete(ctx, fixtures.Secrets[0])).To(Succeed())

		By("expecting secret to have been deleted")
		Eventually(func() error {
			kErr := manager.GetLatest(ctx, fixtures.Secrets[0])
			if kErr == nil || errors.IsNotFound(kErr) {
				return nil
			}
			return assert.Equals("error", "[NOT FOUND]", kErr)
		}, timeout, interval).Should(Succeed())

		By("deleting management context")

		Expect(manager.Client().Delete(ctx, fixtures.Context)).To(Succeed())

		By("expecting management context to have been deleted")

		Eventually(func() error {
			kErr := manager.GetLatest(ctx, fixtures.Context)
			if kErr == nil || errors.IsNotFound(kErr) {
				return nil
			}
			return assert.Equals("error", "[NOT FOUND]", kErr)
		}, timeout, interval).Should(Succeed())
	})
})
