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

	nav "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/navigation"
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

	It("should create portal link in APIM with its declared visibility", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithPortal(constants.PortalFile).
			WithPortalLink(constants.PortalLinkFile).
			WithContext(constants.ContextWithSecretFile).
			Build()

		fixtures.PortalLink.Spec.Visibility = new(nav.Private)

		fixtures.Apply()

		By("expecting portal link status to be completed")

		Expect(assert.PortalLinkAccepted(fixtures.PortalLink)).To(Succeed())

		By("calling rest API, expecting the declared visibility to be applied to the link")

		apim := apim.NewClient(ctx)
		portalHrid := refs.NewNamespacedNameFromObject(fixtures.Portal).HRID()
		linkHrid := refs.NewNamespacedNameFromObject(fixtures.PortalLink).HRID()

		Eventually(func() error {
			link, linkErr := apim.Links.GetByHRID(portalHrid, linkHrid)
			if linkErr != nil {
				return linkErr
			}
			return assert.Equals("Portal link visibility", nav.Private, link.Visibility)
		}, timeout, interval).Should(Succeed(), fixtures.PortalLink.Name)
	})

	It("should create portal link in APIM", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithPortal(constants.PortalFile).
			WithPortalLink(constants.PortalLinkFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting portal link status to be completed")

		Expect(assert.PortalLinkAccepted(fixtures.PortalLink)).To(Succeed())
		Expect(assert.ManagedByAutomationAPI(fixtures.PortalLink)).To(Succeed())

		By("calling rest API, expecting the link to round-trip")

		apim := apim.NewClient(ctx)
		portalHrid := refs.NewNamespacedNameFromObject(fixtures.Portal).HRID()
		linkHrid := refs.NewNamespacedNameFromObject(fixtures.PortalLink).HRID()

		Eventually(func() error {
			link, linkErr := apim.Links.GetByHRID(portalHrid, linkHrid)
			if linkErr != nil {
				return linkErr
			}
			if err := assert.NotEmptyString("id", link.ID); err != nil {
				return err
			}
			if err := assert.Equals("Portal link href", fixtures.PortalLink.Spec.Href, link.Href); err != nil {
				return err
			}
			return assert.Equals("Portal link location", *fixtures.PortalLink.Spec.Location, link.Location)
		}, timeout, interval).Should(Succeed(), fixtures.PortalLink.Name)
	})
})
