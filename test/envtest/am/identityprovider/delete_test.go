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

package amidentityprovider

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Delete", func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should delete the identity provider from AM (D1)", func() {
		idp := fixtures().Apply().AMIdentityProvider

		Expect(manager.Client().Delete(ctx, idp.DeepCopy())).To(Succeed())

		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMIdentityProvider", idp)
		}, timeout, interval).Should(Succeed(), idp.Name)
		status, _ := remoteIdentityProvider(ctx, idp)
		Expect(status).To(Equal(http.StatusNotFound))
	})

	It("should release the identity provider once its domain is deleted (D3)", func() {
		fixtures := fixtures().Apply()
		domain, idp := fixtures.AMSecurityDomain, fixtures.AMIdentityProvider

		By("deleting the domain: AM deletes its identity providers")

		Expect(manager.Client().Delete(ctx, domain.DeepCopy())).To(Succeed())
		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMSecurityDomain", domain)
		}, timeout, interval).Should(Succeed(), domain.Name)
		status, _ := remoteIdentityProvider(ctx, idp)
		Expect(status).To(Equal(http.StatusNotFound))

		By("deleting the identity provider by hand (envtest has no garbage collector)")

		Expect(manager.Client().Delete(ctx, idp.DeepCopy())).To(Succeed())
		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMIdentityProvider", idp)
		}, timeout, interval).Should(Succeed(), idp.Name)
	})

	It("should keep the AMContext while a domain references it (D7)", func() {
		fixtures := fixtures().Apply()
		amContext, domain := fixtures.AMContext, fixtures.AMSecurityDomain

		By("deleting the AMContext: its finalizer stays")

		Expect(manager.Client().Delete(ctx, amContext.DeepCopy())).To(Succeed())
		Consistently(ctx, func() error {
			if err := manager.GetLatest(ctx, amContext); err != nil {
				return err
			}
			return assert.HasFinalizer(amContext, core.AMContextFinalizer)
		}, constants.ConsistentTimeout, interval).Should(Succeed(), amContext.Name)
		Expect(amContext.IsBeingDeleted()).To(BeTrue())

		By("deleting the domain")

		Expect(manager.Client().Delete(ctx, domain.DeepCopy())).To(Succeed())
		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMSecurityDomain", domain)
		}, timeout, interval).Should(Succeed(), domain.Name)

		By("touching the AMContext: a refused delete is not retried, as for a ManagementContext")

		Expect(manager.GetLatest(ctx, amContext)).To(Succeed())
		amContext.Annotations["test.gravitee.io/touch"] = "true"
		Expect(manager.Client().Update(ctx, amContext)).To(Succeed())

		By("expecting the AMContext to go")

		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMContext", amContext)
		}, timeout, interval).Should(Succeed(), amContext.Name)
	})
})
