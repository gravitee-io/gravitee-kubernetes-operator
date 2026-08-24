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

	It("should create theme in APIM and persist its definition", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithPortalTheme(constants.PortalThemeFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting theme status to be completed")

		Expect(assert.PortalThemeAccepted(fixtures.PortalTheme)).To(Succeed())
		Expect(assert.ManagedByAutomationAPI(fixtures.PortalTheme)).To(Succeed())

		By("calling rest API, expecting the definition to round-trip")

		apim := apim.NewClient(ctx)
		hrid := refs.NewNamespacedNameFromObject(fixtures.PortalTheme).HRID()
		expected := fixtures.PortalTheme.Spec.Definition

		Eventually(func() error {
			thm, thmErr := apim.PortalThemes.GetByHRID(hrid)
			if thmErr != nil {
				return thmErr
			}
			if err := assert.NotEmptyString("id", thm.ID); err != nil {
				return err
			}
			if err := assert.Equals("Theme name", fixtures.PortalTheme.Spec.Name, thm.Name); err != nil {
				return err
			}
			if err := assert.Equals(
				"Theme font family",
				*expected.Font.FontFamily,
				*thm.Definition.Font.FontFamily,
			); err != nil {
				return err
			}
			return assert.Equals(
				"Theme primary color",
				*expected.Color.Primary,
				*thm.Definition.Color.Primary,
			)
		}, timeout, interval).Should(Succeed(), fixtures.PortalTheme.Name)
	})
})
