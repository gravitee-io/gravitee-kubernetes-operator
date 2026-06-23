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

	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/portallisting"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate update", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := adm.AdmissionCtrl{}

	It("should return severe error when portalRef is changed", func() {
		listing := fixture.
			Builder().
			WithPortalListing(constants.PortalListingFile).
			Build()

		updated := listing.PortalListing.DeepCopy()
		updated.Spec.Portal.Name += "-repointed"

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, listing.PortalListing, updated)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
