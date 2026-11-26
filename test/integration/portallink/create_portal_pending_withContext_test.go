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

// A portal whose theme is missing is not synced to APIM, so its link waits for it (GKO-3480).
var _ = Describe("Create", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should wait for its portal to be synced before creating the link in APIM", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithPortalTheme(constants.PortalThemeFile).
			WithPortal(constants.PortalWithThemeFile).
			WithPortalLink(constants.PortalLinkFile).
			WithContext(constants.ContextWithSecretFile).
			Build()

		By("applying the context, holding back the theme the portal activates")

		themeCRD, portalCRD, linkCRD := fixtures.PortalTheme, fixtures.Portal, fixtures.PortalLink
		fixtures.PortalTheme, fixtures.Portal, fixtures.PortalLink = nil, nil, nil
		fixtures.Apply()

		Expect(manager.Client().Create(ctx, portalCRD)).To(Succeed())
		Expect(manager.Client().Create(ctx, linkCRD)).To(Succeed())

		By("expecting the link to report its portal ref as unresolved, with no APIM id")

		Eventually(func() error {
			if err := manager.GetLatest(ctx, linkCRD); err != nil {
				return err
			}
			if err := assert.IsUnresolvedWith(linkCRD, "portal not synced"); err != nil {
				return err
			}
			return assert.Equals("link status id", "", linkCRD.Status.ID)
		}, timeout, interval).Should(Succeed(), linkCRD.Name)

		By("creating the missing theme, expecting the portal and then the link to be synced")

		Expect(manager.Client().Create(ctx, themeCRD)).To(Succeed())

		Eventually(func() error {
			if err := manager.GetLatest(ctx, linkCRD); err != nil {
				return err
			}
			return assert.PortalLinkAccepted(linkCRD)
		}, timeout, interval).Should(Succeed(), linkCRD.Name)

		apimClient := apim.NewClient(ctx)
		linkHrid := refs.NewNamespacedNameFromObject(linkCRD).HRID()

		Eventually(func() error {
			_, err := apimClient.Links.GetByHRID(service.LinkParent{Portal: portalCRD}, linkHrid)
			return err
		}, timeout, interval).Should(Succeed(), linkCRD.Name)
	})
})
