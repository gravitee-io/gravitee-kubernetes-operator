// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package amapplication

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Delete", func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should delete the application from AM", func() {
		app := fixtures(constants.AMApplicationMinimalFile).Apply().AMApplication
		status, _ := remoteApplication(ctx, app)
		Expect(status).To(Equal(http.StatusOK))

		Expect(manager.Client().Delete(ctx, app.DeepCopy())).To(Succeed())

		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMApplication", app)
		}, timeout, interval).Should(Succeed(), app.Name)
		status, _ = remoteApplication(ctx, app)
		Expect(status).To(Equal(http.StatusNotFound))
	})

	It("should release the application once its domain is deleted", func() {
		fixtures := fixtures(constants.AMApplicationMinimalFile).Apply()
		domain, app := fixtures.AMSecurityDomain, fixtures.AMApplication

		By("deleting the domain: AM deletes its applications")

		Expect(manager.Client().Delete(ctx, domain.DeepCopy())).To(Succeed())
		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMSecurityDomain", domain)
		}, timeout, interval).Should(Succeed(), domain.Name)
		status, _ := remoteApplication(ctx, app)
		Expect(status).To(Equal(http.StatusNotFound))

		By("deleting the application by hand (envtest has no garbage collector)")

		Expect(manager.Client().Delete(ctx, app.DeepCopy())).To(Succeed())
		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMApplication", app)
		}, timeout, interval).Should(Succeed(), app.Name)
	})

	It("should delete the application from AM once its certificate is gone", func() {
		fixtures := fixtures(constants.AMApplicationMinimalFile).Apply()
		cert, app := fixtures.AMCertificate, fixtures.AMApplication

		By("deleting the certificate the application references (no admission webhook in this suite)")

		Expect(manager.Client().Delete(ctx, cert.DeepCopy())).To(Succeed())
		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMCertificate", cert)
		}, timeout, interval).Should(Succeed(), cert.Name)

		By("deleting the application: its finalizer does not wait for the certificate")

		Expect(manager.Client().Delete(ctx, app.DeepCopy())).To(Succeed())
		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMApplication", app)
		}, timeout, interval).Should(Succeed(), app.Name)
		status, _ := remoteApplication(ctx, app)
		Expect(status).To(Equal(http.StatusNotFound))
	})
})
