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

package portaltheme

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portaltheme"
	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/portaltheme"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate create", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := adm.AdmissionCtrl{}

	It("should return severe error when contextRef cannot be resolved", func() {
		thm := fixture.
			Builder().
			WithPortalTheme(constants.PortalThemeFile).
			Build()

		thm.PortalTheme.Spec.Context.Name = "unresolved"

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, thm.PortalTheme)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when contextRef is missing", func() {
		thm := fixture.
			Builder().
			WithPortalTheme(constants.PortalThemeFile).
			Build()

		thm.PortalTheme.Spec.Context = nil

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, thm.PortalTheme)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when the theme is rejected by the dry run", func() {
		thm := fixture.
			Builder().
			WithPortalTheme(constants.PortalThemeFile).
			WithContext(constants.ContextWithSecretFile).
			Build()

		thm.PortalTheme.Spec.Name = ""
		thm.PortalTheme.Spec.Definition = portaltheme.Definition{}

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, thm.PortalTheme)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
