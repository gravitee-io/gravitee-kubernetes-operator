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

package portaltheme

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

	It("should not delete theme while a portal is using it", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithPortalTheme(constants.PortalThemeFile).
			WithPortal(constants.PortalWithThemeFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting the portal to reference the theme")

		Expect(assert.Equals(
			"Portal active theme",
			fixtures.PortalTheme.Name,
			fixtures.Portal.Spec.Theme.Name,
		)).To(Succeed())

		By("deleting theme")

		Expect(manager.Client().Delete(ctx, fixtures.PortalTheme.DeepCopy())).To(Succeed())

		By("expecting the theme to be kept while the portal still uses it")

		Consistently(func() error {
			return manager.GetLatest(ctx, fixtures.PortalTheme)
		}, constants.ConsistentTimeout, interval).Should(Succeed(), fixtures.PortalTheme.Name)

		By("deleting the portal that uses the theme")

		Expect(manager.Client().Delete(ctx, fixtures.Portal.DeepCopy())).To(Succeed())

		By("expecting the theme to be deleted once nothing uses it")

		Eventually(func() error {
			return assert.Deleted(ctx, "PortalTheme", fixtures.PortalTheme)
		}, timeout, interval).Should(Succeed(), fixtures.PortalTheme.Name)
	})
})
