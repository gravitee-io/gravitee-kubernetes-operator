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

// A portal whose theme is missing is not synced to APIM, so its listing waits for it (GKO-3480).
var _ = Describe("Create", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should wait for its portal to be synced before creating the listing in APIM", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithAPIv4(constants.ApiV4WithContextFile).
			WithPortalTheme(constants.PortalThemeFile).
			WithPortal(constants.PortalWithThemeFile).
			WithPortalListing(constants.PortalListingFile).
			WithContext(constants.ContextWithSecretFile).
			Build()

		By("applying the context and the API, holding back the theme the portal activates")

		themeCRD, portalCRD, listingCRD := fixtures.PortalTheme, fixtures.Portal, fixtures.PortalListing
		fixtures.PortalTheme, fixtures.Portal, fixtures.PortalListing = nil, nil, nil
		fixtures.Apply()

		Expect(manager.Client().Create(ctx, portalCRD)).To(Succeed())
		Expect(manager.Client().Create(ctx, listingCRD)).To(Succeed())

		By("expecting the listing to report its portal ref as unresolved, with no APIM id")

		Eventually(func() error {
			if err := manager.GetLatest(ctx, listingCRD); err != nil {
				return err
			}
			if err := assert.IsUnresolvedWith(listingCRD, "portal not synced"); err != nil {
				return err
			}
			return assert.Equals("listing status id", "", listingCRD.Status.ID)
		}, timeout, interval).Should(Succeed(), listingCRD.Name)

		By("creating the missing theme, expecting the portal and then the listing to be synced")

		Expect(manager.Client().Create(ctx, themeCRD)).To(Succeed())

		Eventually(func() error {
			if err := manager.GetLatest(ctx, listingCRD); err != nil {
				return err
			}
			return assert.PortalListingAccepted(listingCRD)
		}, timeout, interval).Should(Succeed(), listingCRD.Name)

		apim := apim.NewClient(ctx)
		portalHrid := refs.NewNamespacedNameFromObject(portalCRD).HRID()
		listingHrid := refs.NewNamespacedNameFromObject(listingCRD).HRID()

		Eventually(func() error {
			_, err := apim.Listings.GetByHRID(portalHrid, listingHrid)
			return err
		}, timeout, interval).Should(Succeed(), listingCRD.Name)
	})
})
