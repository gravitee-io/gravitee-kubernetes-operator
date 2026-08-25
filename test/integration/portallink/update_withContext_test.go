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

package portallink

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
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

	It("should update portal link location in APIM", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithPortal(constants.PortalFile).
			WithPortalLink(constants.PortalLinkFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting portal link status to be completed")

		Expect(assert.PortalLinkAccepted(fixtures.PortalLink)).To(Succeed())

		By("calling rest API, expecting to find portal link")

		apim := apim.NewClient(ctx)
		parent := service.LinkParent{Portal: fixtures.Portal}
		linkHrid := refs.NewNamespacedNameFromObject(fixtures.PortalLink).HRID()

		Eventually(func() error {
			link, linkErr := apim.Links.GetByHRID(parent, linkHrid)
			if linkErr != nil {
				return linkErr
			}
			return assert.NotEmptyString("id", link.ID)
		}, timeout, interval).Should(Succeed(), fixtures.PortalLink.Name)

		By("updating the link location to another portal navigation path")

		updated := fixtures.PortalLink.DeepCopy()
		newLocation := fixtures.GetNavigationRoot() + "/projects/beta"
		updated.Spec.Location = &newLocation

		Expect(manager.UpdateSafely(ctx, updated)).To(Succeed())

		By("calling rest API, expecting link location to be up to date")

		Eventually(func() error {
			link, linkErr := apim.Links.GetByHRID(parent, linkHrid)
			if linkErr != nil {
				return linkErr
			}
			return assert.Equals("Portal link location", newLocation, link.Location)
		}, timeout, interval).Should(Succeed(), fixtures.PortalLink.Name)
	})
})
