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
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Delete, active theme gone", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should delete the portal from APIM when the theme it activates was deleted first", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithPortalTheme(constants.PortalThemeFile).
			WithPortal(constants.PortalWithThemeFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		Expect(assert.PortalAccepted(fixtures.Portal)).To(Succeed())

		apimClient := apim.NewClient(ctx)
		hrid := refs.NewNamespacedNameFromObject(fixtures.Portal).HRID()
		themeHrid := refs.NewNamespacedNameFromObject(fixtures.PortalTheme).HRID()
		DeferCleanup(func() {
			Expect(apimClient.DeletePortal(ctx, hrid)).To(Succeed())
			Expect(apimClient.DeletePortalTheme(ctx, themeHrid)).To(Succeed())
		})

		By("deleting the theme first, removing the finalizer its portal keeps")

		Expect(manager.Client().Delete(ctx, fixtures.PortalTheme.DeepCopy())).To(Succeed())
		Eventually(func() error {
			latest := fixtures.PortalTheme.DeepCopy()
			if err := manager.GetLatest(ctx, latest); err != nil {
				return err
			}
			latest.SetFinalizers(nil)
			return manager.Client().Update(ctx, latest)
		}, timeout, interval).Should(Succeed())
		Eventually(func() error {
			return assert.Deleted(ctx, "PortalTheme", fixtures.PortalTheme)
		}, timeout, interval).Should(Succeed(), fixtures.PortalTheme.Name)

		By("deleting the portal")

		Expect(manager.Client().Delete(ctx, fixtures.Portal.DeepCopy())).To(Succeed())

		By("expecting the portal to be deleted from k8s")

		Eventually(func() error {
			return assert.Deleted(ctx, "Portal", fixtures.Portal)
		}, timeout, interval).Should(Succeed(), fixtures.Portal.Name)

		By("calling rest API, expecting the portal to be deleted from APIM")

		_, err := apimClient.GetPortal(ctx, hrid)
		Expect(assert.NotFoundError(err)).To(Succeed())
	})
})
