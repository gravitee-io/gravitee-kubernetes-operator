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

package docs

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/docs"
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

	It("should return severe error when a portal documentation is reassigned to an API", func() {
		fixtures := fixture.
			Builder().
			WithDocumentation(constants.DocumentationPortalFile).
			Build()

		updated := fixtures.Documentation.DeepCopy()
		updated.Spec.Portal = nil
		updated.Spec.API = &refs.NamespacedName{Name: "some-api"}

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Documentation, updated)
			return assert.SevereError(
				errors.NewSevere("a documentation cannot be reassigned from a portal to an API"),
				err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when an API documentation is reassigned to a portal", func() {
		fixtures := fixture.
			Builder().
			WithDocumentation(constants.DocumentationApiFile).
			Build()

		updated := fixtures.Documentation.DeepCopy()
		updated.Spec.API = nil
		updated.Spec.Portal = &refs.NamespacedName{Name: "some-portal"}

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Documentation, updated)
			return assert.SevereError(
				errors.NewSevere("a documentation cannot be reassigned from an API to a portal"), err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when portalRef is changed", func() {
		fixtures := fixture.
			Builder().
			WithDocumentation(constants.DocumentationPortalFile).
			Build()

		updated := fixtures.Documentation.DeepCopy()
		updated.Spec.Portal.Name += "-repointed"

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Documentation, updated)
			return assert.SevereError(
				errors.NewSeveref("portalRef is immutable; documentation cannot be moved to a different portal (from [%s] to [%s])",
					fixtures.Documentation.Spec.GetPortalRef(),
					updated.Spec.GetPortalRef()),
				err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when apiRef is changed", func() {
		fixtures := fixture.
			Builder().
			WithDocumentation(constants.DocumentationApiFile).
			Build()

		updated := fixtures.Documentation.DeepCopy()
		updated.Spec.API.Name += "-repointed"

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Documentation, updated)
			return assert.SevereError(
				errors.NewSeveref("apiRef is immutable; documentation cannot be moved to a different API (from [%s] to [%s])",
					fixtures.Documentation.Spec.GetApiRef(),
					updated.Spec.GetApiRef()),
				err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
