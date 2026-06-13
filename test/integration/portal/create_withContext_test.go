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

package portal

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

	It("should create portal in APIM and persist navigation in list order", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithPortal(constants.PortalFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting portal status to be completed")

		Expect(assert.PortalAccepted(fixtures.Portal)).To(Succeed())
		Expect(assert.ManagedByAutomationAPI(fixtures.Portal)).To(Succeed())

		By("calling rest API, expecting navigation to round-trip in the same order")

		apim := apim.NewClient(ctx)
		hrid := refs.NewNamespacedNameFromObject(fixtures.Portal).HRID()

		expectedPaths := make([]string, 0, len(fixtures.Portal.Spec.Navigation))
		for _, nav := range fixtures.Portal.Spec.Navigation {
			expectedPaths = append(expectedPaths, nav.Path)
		}

		Eventually(func() error {
			prtl, prtlErr := apim.Portals.GetByHRID(hrid)
			if prtlErr != nil {
				return prtlErr
			}
			if err := assert.NotEmptyString("id", prtl.ID); err != nil {
				return err
			}
			paths := make([]string, 0, len(prtl.Navigation))
			for _, nav := range prtl.Navigation {
				paths = append(paths, nav.Path)
			}
			// APIM materialises implicit intermediate parent folders in the GET
			// /portals/{hrid} response (e.g. /projects, /archive), so the round-trip
			// is not lossless. The authored entries are asserted as an ordered
			// subsequence: they must appear in list order, with the implicit folders
			// tolerated in between.
			return assert.ContainsInOrder("Portal navigation", expectedPaths, paths)
		}, timeout, interval).Should(Succeed(), fixtures.Portal.Name)
	})
})
