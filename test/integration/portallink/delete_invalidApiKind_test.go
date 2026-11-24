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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

// The CRD itself rejects a link with both or neither of portalRef / apiRef. An apiRef of another kind
// than ApiV4Definition is refused by the webhook only: without it, the link reaches the reconcile.
var _ = Describe("Delete, apiRef of a v2 API", func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should release a link that references a v2 API, which was never synced", func() {
		link := fixture.Builder().
			WithPortalLink(constants.PortalLinkApiFile).
			Build().
			PortalLink
		link.Spec.API = &refs.NamespacedName{Name: "legacy-api", Kind: "ApiDefinition"}

		By("creating the link without the webhook, expecting the reconcile to save its finalizer")

		Expect(manager.Client().Create(ctx, link)).To(Succeed())
		Eventually(func() error {
			if err := manager.GetLatest(ctx, link); err != nil {
				return err
			}
			if !controllerutil.ContainsFinalizer(link, core.PortalLinkFinalizer) {
				return assert.Equals("finalizer", core.PortalLinkFinalizer, link.GetFinalizers())
			}
			return assert.IsUnresolved(link)
		}, timeout, interval).Should(Succeed(), link.Name)

		By("deleting the link, expecting it to be deleted from k8s")

		Expect(manager.Client().Delete(ctx, link.DeepCopy())).To(Succeed())
		Eventually(func() error {
			return assert.Deleted(ctx, "PortalLink", link)
		}, timeout, interval).Should(Succeed(), link.Name)
	})
})
