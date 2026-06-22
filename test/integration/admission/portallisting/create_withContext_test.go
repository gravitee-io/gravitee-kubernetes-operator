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
	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/portallisting"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate create", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := adm.AdmissionCtrl{}

	It("should return severe error when portalRef cannot be resolved", func() {
		listing := fixture.
			Builder().
			WithPortalListing(constants.PortalListingFile).
			Build()

		listing.PortalListing.Spec.Portal.Name = "unresolved"

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, listing.PortalListing)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when an API ref is not a v4 API kind", func() {
		listing := fixture.
			Builder().
			WithPortalListing(constants.PortalListingFile).
			Build()

		listing.PortalListing.Spec.APIs[0].Ref.Kind = "ApiDefinition"

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, listing.PortalListing)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when an API context differs from the portal context", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithAPIv4(constants.ApiV4WithContextFile).
			WithPortal(constants.PortalFile).
			WithPortalListing(constants.PortalListingFile).
			WithContext(constants.ContextWithSecretFile).
			Build()

		// point the API at a different management context than the portal
		fixtures.APIv4.Spec.Context = &refs.NamespacedName{
			Name:      "other-context",
			Namespace: constants.Namespace,
		}

		Expect(manager.Client().Create(ctx, fixtures.Context)).To(Succeed())
		Expect(manager.Client().Create(ctx, fixtures.Portal)).To(Succeed())
		Expect(manager.Client().Create(ctx, fixtures.APIv4)).To(Succeed())

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.PortalListing)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error (no panic) when a listed API has no management context", func() {
		fixtures := fixture.Builder().
			WithAPIv4(constants.ApiV4WithContextFile).
			WithPortal(constants.PortalFile).
			WithPortalListing(constants.PortalListingFile).
			WithContext(constants.ContextWithSecretFile).
			Build()

		// the API has no management context (dbless mode); the portal keeps its context
		fixtures.APIv4.Spec.Context = nil

		Expect(manager.Client().Create(ctx, fixtures.Portal)).To(Succeed())
		Expect(manager.Client().Create(ctx, fixtures.APIv4)).To(Succeed())

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.PortalListing)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when the referenced portal has no management context", func() {
		fixtures := fixture.Builder().
			WithPortal(constants.PortalFile).
			WithPortalListing(constants.PortalListingFile).
			Build()

		// portal without a management context
		fixtures.Portal.Spec.Context = nil

		Expect(manager.Client().Create(ctx, fixtures.Portal)).To(Succeed())

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.PortalListing)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
