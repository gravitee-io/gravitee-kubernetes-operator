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
)

var _ = Describe("Create", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should create portal listing in APIM and persist apis in list order", func() {
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
		Expect(assert.ManagedByAutomationAPI(fixtures.PortalListing)).To(Succeed())

		By("calling rest API, expecting apis to round-trip in the same order")

		apim := apim.NewClient(ctx)
		portalHrid := refs.NewNamespacedNameFromObject(fixtures.Portal).HRID()
		listingHrid := refs.NewNamespacedNameFromObject(fixtures.PortalListing).HRID()

		expectedLocations := make([]string, 0, len(fixtures.PortalListing.Spec.APIs))
		for _, entry := range fixtures.PortalListing.Spec.APIs {
			expectedLocations = append(expectedLocations, entry.Location)
		}

		Eventually(func() error {
			listing, listingErr := apim.Listings.GetByHRID(portalHrid, listingHrid)
			if listingErr != nil {
				return listingErr
			}
			if err := assert.NotEmptyString("id", listing.ID); err != nil {
				return err
			}
			locations := make([]string, 0, len(listing.APIs))
			for _, entry := range listing.APIs {
				locations = append(locations, entry.Location)
			}
			return assert.Equals("Portal listing locations", expectedLocations, locations)
		}, timeout, interval).Should(Succeed(), fixtures.PortalListing.Name)
	})
})
