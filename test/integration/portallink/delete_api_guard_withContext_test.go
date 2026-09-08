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

package portallink

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Delete API guard", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should block API deletion while a portal link references it", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithAPIv4(constants.ApiV4WithContextFile).
			WithPortalLink(constants.PortalLinkApiFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting portal link status to be completed")

		Expect(assert.PortalLinkAccepted(fixtures.PortalLink)).To(Succeed())

		By("deleting the API")

		Expect(manager.Client().Delete(ctx, fixtures.APIv4)).To(Succeed())

		By("expecting to still find the API while the link references it")

		checkUntil := constants.ConsistentTimeout
		Consistently(func() error {
			return manager.GetLatest(ctx, fixtures.APIv4)
		}, checkUntil, interval).Should(Succeed())

		By("deleting the portal link")

		Expect(manager.Client().Delete(ctx, fixtures.PortalLink)).To(Succeed())

		By("expecting portal link to be deleted from k8s")

		Eventually(func() error {
			return assert.Deleted(ctx, "PortalLink", fixtures.PortalLink)
		}, timeout, interval).Should(Succeed(), fixtures.PortalLink.Name)

		By("expecting the API to have been deleted once unreferenced")

		Eventually(func() error {
			return assert.Deleted(ctx, "ApiV4Definition", fixtures.APIv4)
		}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)
	})
})
