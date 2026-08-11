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

	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/portallink"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
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
		link := fixture.
			Builder().
			WithPortalLink(constants.PortalLinkFile).
			Build()

		updated := link.PortalLink.DeepCopy()
		updated.Spec.Portal.Name += "-repointed"

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, link.PortalLink, updated)
			return assert.SevereError(
				errors.NewSeveref(
					"portalRef is immutable. Detected change from [%s] to [%s]",
					link.PortalLink.Spec.Portal.String(), updated.Spec.Portal.String(),
				),
				err,
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
