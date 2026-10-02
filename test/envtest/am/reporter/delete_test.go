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
package amreporter

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

	It("should delete the reporter from AM (D1)", func() {
		reporter := fixtures().Apply().AMReporter

		Expect(manager.Client().Delete(ctx, reporter.DeepCopy())).To(Succeed())

		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMReporter", reporter)
		}, timeout, interval).Should(Succeed(), reporter.Name)
		status, _ := remoteReporter(ctx, reporter)
		Expect(status).To(Equal(http.StatusNotFound))
	})

	It("should release the reporter once its domain is deleted (D3)", func() {
		fixtures := fixtures().Apply()
		domain, reporter := fixtures.AMSecurityDomain, fixtures.AMReporter

		By("deleting the domain: AM deletes its reporters")

		Expect(manager.Client().Delete(ctx, domain.DeepCopy())).To(Succeed())
		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMSecurityDomain", domain)
		}, timeout, interval).Should(Succeed(), domain.Name)
		status, _ := remoteReporter(ctx, reporter)
		Expect(status).To(Equal(http.StatusNotFound))

		By("deleting the reporter by hand (envtest has no garbage collector)")

		Expect(manager.Client().Delete(ctx, reporter.DeepCopy())).To(Succeed())
		Eventually(ctx, func() error {
			return assert.Deleted(ctx, "AMReporter", reporter)
		}, timeout, interval).Should(Succeed(), reporter.Name)
	})
})
