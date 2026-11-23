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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/navigation"
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

	It("should create a documentation page with an empty content", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithPortal(constants.PortalFileDeprecated).
			WithDocumentation(constants.DocumentationPortalFile).
			WithContext(constants.ContextWithSecretFile).
			Build()

		// The APIM SDK sends "content":"" where GKO used to omit an empty content.
		fixtures.Documentation.Spec.Content = ""

		fixtures.Apply()

		By("expecting documentation status to be completed")

		Expect(assert.DocumentationAccepted(fixtures.Documentation)).To(Succeed())

		By("calling rest API, expecting the documentation page to have an empty content")

		apimClient := apim.NewClient(ctx)
		docHrid := refs.NewNamespacedNameFromObject(fixtures.Documentation).HRID()

		Eventually(func() error {
			doc, docErr := apimClient.GetDocumentation(ctx,
				navigation.Parent{PortalHrid: refs.NewNamespacedNameFromObject(fixtures.Portal).HRID()}, docHrid,
			)
			if docErr != nil {
				return docErr
			}
			return assert.Equals("Documentation content", "", doc.Content)
		}, timeout, interval).Should(Succeed(), fixtures.Documentation.Name)
	})
})
