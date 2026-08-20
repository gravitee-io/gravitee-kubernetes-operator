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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/portallisting"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	updatedLocationSuffix = "/projects/beta"
	remoteLocationSuffix  = "/archive/2024"
	localLocationSuffix   = "/projects/alpha/docs"
)

var _ = Describe("Validate drift", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := admission.AdmissionCtrl{}

	It("should not drift on with minimal fields", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4WithContextFile).
			WithPortal(constants.PortalDriftFullFile).
			WithPortalListing(constants.PortalListingDriftMinimalFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the portal listing location")
		newListing := fixtures.PortalListing.DeepCopy()
		setLocation(newListing, fixtures.GetNavigationRoot()+updatedLocationSuffix)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.PortalListing, newListing)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should detect drift with minimal fields", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4WithContextFile).
			WithPortal(constants.PortalDriftFullFile).
			WithPortalListing(constants.PortalListingDriftMinimalFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the remote portal listing location")
		validateLocationDrift(ctx, admissionCtrl, fixtures)
	})

	It("should not drift on with all fields", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4WithContextFile).
			WithPortal(constants.PortalDriftFullFile).
			WithPortalListing(constants.PortalListingDriftFullFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the portal listing location")
		newListing := fixtures.PortalListing.DeepCopy()
		setLocation(newListing, fixtures.GetNavigationRoot()+updatedLocationSuffix)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.PortalListing, newListing)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should detect drift with all fields", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4WithContextFile).
			WithPortal(constants.PortalDriftFullFile).
			WithPortalListing(constants.PortalListingDriftFullFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the remote portal listing location")
		validateLocationDrift(ctx, admissionCtrl, fixtures)
	})
})

func setLocation(listing *v1alpha1.PortalListing, location string) {
	GinkgoHelper()
	listing.Spec.APIs[0].Location = location
}

func validateLocationDrift(
	ctx context.Context,
	admissionCtrl admission.AdmissionCtrl,
	fixtures *fixture.Objects,
) {
	GinkgoHelper()
	apimClient := apim.NewClient(ctx)

	oldListing := fixtures.PortalListing
	newListing := fixtures.PortalListing.DeepCopy()

	remoteLocation := fixtures.GetNavigationRoot() + remoteLocationSuffix
	localLocation := fixtures.GetNavigationRoot() + localLocationSuffix

	setLocation(newListing, remoteLocation)

	_, err := apimClient.Listings.CreateOrUpdate(newListing, fixtures.Portal)
	Expect(err).ToNot(HaveOccurred())

	setLocation(newListing, localLocation)

	Eventually(func() error {
		_, err := admissionCtrl.ValidateUpdate(ctx, oldListing, newListing)
		return assert.DriftDetected(
			fmt.Sprintf(`apis[0]:
  location: "%s" != "%s"`, localLocation, remoteLocation),
			err,
		)
	}, constants.EventualTimeout, constants.Interval).Should(Succeed())
}
