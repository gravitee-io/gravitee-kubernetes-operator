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

package portal

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

// forceDeleteContext deletes a management context the way a user removing its finalizer by hand does: its
// own finalizer keeps it while portals reference it.
func forceDeleteContext(ctx context.Context, mCtx *v1alpha1.ManagementContext) {
	GinkgoHelper()
	Expect(manager.Client().Delete(ctx, mCtx.DeepCopy())).To(Succeed())
	Eventually(func() error {
		latest := mCtx.DeepCopy()
		if err := manager.GetLatest(ctx, latest); err != nil {
			return err
		}
		latest.SetFinalizers(nil)
		return manager.Client().Update(ctx, latest)
	}, constants.EventualTimeout, constants.Interval).Should(Succeed())
	Eventually(func() error {
		return assert.Deleted(ctx, "ManagementContext", mCtx)
	}, constants.EventualTimeout, constants.Interval).Should(Succeed(), mCtx.Name)
}

var _ = Describe("Delete, management context gone", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should release the portal when its management context was deleted first", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithPortal(constants.PortalFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		Expect(assert.PortalAccepted(fixtures.Portal)).To(Succeed())

		apimClient := apim.NewClient(ctx)
		hrid := refs.NewNamespacedNameFromObject(fixtures.Portal).HRID()
		DeferCleanup(func() {
			Expect(apimClient.DeletePortal(ctx, hrid)).To(Succeed())
		})

		By("deleting the management context first")

		forceDeleteContext(ctx, fixtures.Context)

		By("deleting the portal")

		Expect(manager.Client().Delete(ctx, fixtures.Portal.DeepCopy())).To(Succeed())

		By("expecting the portal to be deleted from k8s")

		Eventually(func() error {
			return assert.Deleted(ctx, "Portal", fixtures.Portal)
		}, timeout, interval).Should(Succeed(), fixtures.Portal.Name)
	})

	It("should keep a portal a link references when its management context was deleted first", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithPortal(constants.PortalFile).
			WithPortalLink(constants.PortalLinkFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		Expect(assert.PortalLinkAccepted(fixtures.PortalLink)).To(Succeed())

		apimClient := apim.NewClient(ctx)
		portalHrid := refs.NewNamespacedNameFromObject(fixtures.Portal).HRID()
		linkHrid := refs.NewNamespacedNameFromObject(fixtures.PortalLink).HRID()
		DeferCleanup(func() {
			Expect(apimClient.DeletePortalLink(ctx, navigation.Parent{PortalHrid: portalHrid}, linkHrid)).To(Succeed())
			Expect(apimClient.DeletePortal(ctx, portalHrid)).To(Succeed())
		})

		By("deleting the management context first")

		forceDeleteContext(ctx, fixtures.Context)

		By("deleting the portal")

		Expect(manager.Client().Delete(ctx, fixtures.Portal.DeepCopy())).To(Succeed())

		By("expecting to still find the portal while the link references it")

		Consistently(func() error {
			return manager.GetLatest(ctx, fixtures.Portal.DeepCopy())
		}, constants.ConsistentTimeout, interval).Should(Succeed())

		By("deleting the link")

		Expect(manager.Client().Delete(ctx, fixtures.PortalLink.DeepCopy())).To(Succeed())

		By("expecting the link, then the portal, to be deleted from k8s")

		Eventually(func() error {
			return assert.Deleted(ctx, "PortalLink", fixtures.PortalLink)
		}, timeout, interval).Should(Succeed(), fixtures.PortalLink.Name)
		Eventually(func() error {
			return assert.Deleted(ctx, "Portal", fixtures.Portal)
		}, timeout, interval).Should(Succeed(), fixtures.Portal.Name)
	})
})
