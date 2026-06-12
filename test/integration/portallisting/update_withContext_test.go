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

var _ = Describe("Update", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should update portal listing location in APIM", func() {
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

		By("calling rest API, expecting to find portal listing")

		apim := apim.NewClient(ctx)
		portalHrid := refs.NewNamespacedNameFromObject(fixtures.Portal).HRID()
		listingHrid := refs.NewNamespacedNameFromObject(fixtures.PortalListing).HRID()

		Eventually(func() error {
			listing, listingErr := apim.Listings.GetByHRID(portalHrid, listingHrid)
			if listingErr != nil {
				return listingErr
			}
			return assert.NotEmptyString("id", listing.ID)
		}, timeout, interval).Should(Succeed(), fixtures.PortalListing.Name)

		By("updating the listing location to another portal navigation path")

		updated := fixtures.PortalListing.DeepCopy()
		updated.Spec.APIs[0].Location = "/projects/beta"

		Expect(manager.UpdateSafely(ctx, updated)).To(Succeed())

		By("calling rest API, expecting listing location to be up to date")

		Eventually(func() error {
			listing, listingErr := apim.Listings.GetByHRID(portalHrid, listingHrid)
			if listingErr != nil {
				return listingErr
			}
			if len(listing.APIs) == 0 {
				return assert.NotEmptyString("apis", "")
			}
			return assert.Equals("Portal listing location", "/projects/beta", listing.APIs[0].Location)
		}, timeout, interval).Should(Succeed(), fixtures.PortalListing.Name)
	})
})
