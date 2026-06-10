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
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Delete", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should delete portal in APIM", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithPortal(constants.PortalFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting portal status to be completed")

		Expect(assert.PortalAccepted(fixtures.Portal)).To(Succeed())

		By("calling rest API, expecting to find portal")

		apim := apim.NewClient(ctx)
		hrid := refs.NewNamespacedNameFromObject(fixtures.Portal).HRID()

		Eventually(func() error {
			prtl, prtlErr := apim.Portals.GetByHRID(hrid)
			if prtlErr != nil {
				return prtlErr
			}
			return assert.NotEmptyString("id", prtl.ID)
		}, timeout, interval).Should(Succeed(), fixtures.Portal.Name)

		By("deleting portal")

		Expect(manager.Client().Delete(ctx, fixtures.Portal.DeepCopy())).To(Succeed())

		By("expecting portal to be deleted from k8s")

		Eventually(func() error {
			return assert.Deleted(ctx, "Portal", fixtures.Portal)
		}, timeout, interval).Should(Succeed(), fixtures.Portal.Name)

		By("calling rest API, expecting not to find portal")

		Eventually(func() error {
			_, prtlErr := apim.Portals.GetByHRID(hrid)
			return assert.NotFoundError(prtlErr)
		}, timeout, interval).Should(Succeed(), fixtures.Portal.Name)
	})
})
