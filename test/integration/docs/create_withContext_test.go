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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
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

	It("should create a portal-attached documentation page in APIM", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithPortal(constants.PortalFile).
			WithDocumentation(constants.DocumentationPortalFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting documentation status to be completed")

		Expect(assert.DocumentationAccepted(fixtures.Documentation)).To(Succeed())
		Expect(assert.ManagedByAutomationAPI(fixtures.Documentation)).To(Succeed())

		By("calling rest API, expecting to find the documentation page")

		apimClient := apim.NewClient(ctx)
		portalHrid := refs.NewNamespacedNameFromObject(fixtures.Portal).HRID()
		docHrid := refs.NewNamespacedNameFromObject(fixtures.Documentation).HRID()

		Eventually(func() error {
			doc, docErr := apimClient.Documentations.GetByHRID(
				service.DocumentationParent{Portal: portalHrid}, docHrid,
			)
			if docErr != nil {
				return docErr
			}
			return assert.NotEmptyString("id", doc.ID)
		}, timeout, interval).Should(Succeed(), fixtures.Documentation.Name)
	})

	It("should create an API-attached documentation page in APIM", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithAPIv4(constants.ApiV4WithContextFile).
			WithDocumentation(constants.DocumentationApiFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting documentation status to be completed")

		Expect(assert.DocumentationAccepted(fixtures.Documentation)).To(Succeed())
		Expect(assert.ManagedByAutomationAPI(fixtures.Documentation)).To(Succeed())

		By("calling rest API, expecting to find the documentation page")

		apimClient := apim.NewClient(ctx)
		apiHrid := refs.NewNamespacedNameFromObject(fixtures.APIv4).HRID()
		docHrid := refs.NewNamespacedNameFromObject(fixtures.Documentation).HRID()

		Eventually(func() error {
			doc, docErr := apimClient.Documentations.GetByHRID(
				service.DocumentationParent{API: apiHrid}, docHrid,
			)
			if docErr != nil {
				return docErr
			}
			return assert.NotEmptyString("id", doc.ID)
		}, timeout, interval).Should(Succeed(), fixtures.Documentation.Name)
	})
})
