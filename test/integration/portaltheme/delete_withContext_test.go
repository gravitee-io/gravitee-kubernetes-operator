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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Delete", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should delete theme in APIM", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithPortalTheme(constants.PortalThemeFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting theme status to be completed")

		Expect(assert.PortalThemeAccepted(fixtures.PortalTheme)).To(Succeed())

		By("calling rest API, expecting to find theme")

		apim := apim.NewClient(ctx)
		hrid := refs.NewNamespacedNameFromObject(fixtures.PortalTheme).HRID()

		Eventually(func() error {
			thm, thmErr := apim.PortalThemes.GetByHRID(hrid)
			if thmErr != nil {
				return thmErr
			}
			return assert.NotEmptyString("id", thm.ID)
		}, timeout, interval).Should(Succeed(), fixtures.PortalTheme.Name)

		By("deleting theme")

		Expect(manager.Client().Delete(ctx, fixtures.PortalTheme.DeepCopy())).To(Succeed())

		By("expecting theme to be deleted from k8s")

		Eventually(func() error {
			return assert.Deleted(ctx, "PortalTheme", fixtures.PortalTheme)
		}, timeout, interval).Should(Succeed(), fixtures.PortalTheme.Name)

		By("calling rest API, expecting not to find theme")

		Eventually(func() error {
			_, thmErr := apim.PortalThemes.GetByHRID(hrid)
			return assert.NotFoundError(thmErr)
		}, timeout, interval).Should(Succeed(), fixtures.PortalTheme.Name)
	})
})
