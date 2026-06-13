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

package portallisting

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

var _ = Describe("Delete Portal guard", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should block Portal deletion while a portal listing references it", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithAPIv4(constants.ApiV4WithContextFile).
			WithPortal(constants.PortalFile).
			WithPortalListing(constants.PortalListingFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting portal listing status to be completed")

		Expect(assert.PortalListingAccepted(fixtures.PortalListing)).To(Succeed())

		By("deleting the portal")

		Expect(manager.Client().Delete(ctx, fixtures.Portal)).To(Succeed())

		By("expecting to still find the portal while the listing references it")

		checkUntil := constants.ConsistentTimeout
		Consistently(func() error {
			return manager.GetLatest(ctx, fixtures.Portal)
		}, checkUntil, interval).Should(Succeed())

		By("deleting the portal listing")

		Expect(manager.Client().Delete(ctx, fixtures.PortalListing)).To(Succeed())

		By("expecting portal listing to be deleted from k8s")

		Eventually(func() error {
			return assert.Deleted(ctx, "PortalListing", fixtures.PortalListing)
		}, timeout, interval).Should(Succeed(), fixtures.PortalListing.Name)

		By("expecting the portal to have been deleted once unreferenced")

		Eventually(func() error {
			return assert.Deleted(ctx, "Portal", fixtures.Portal)
		}, timeout, interval).Should(Succeed(), fixtures.Portal.Name)
	})
})
