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

package v4

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should update inline markdown in APIM", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4WithMarkdownPage).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		apim := apim.NewClient(ctx)
		apiID := fixtures.APIv4.Status.ID

		By("updating markdown content")

		markdown := fixtures.APIv4.Spec.Pages["markdown"]
		updatedContent := markdown.Content + "\n" + "This is an update."
		markdown.Content = updatedContent

		Eventually(func() error {
			return manager.UpdateSafely(ctx, fixtures.APIv4)
		}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)

		By("checking pages number in APIM")

		Eventually(func() error {
			pages, err := apim.Pages.FindByAPIV4(apiID)
			if err != nil {
				return err
			}
			return assert.SliceOfSize("pages", pages, 2)
		}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)

		By("checking swagger content in APIM")

		Expect(markdown).ToNot(BeNil())

		Eventually(func() error {
			pages, err := apim.Pages.FindByAPI(apiID, model.NewPageQuery().WithType("MARKDOWN"))
			if err != nil {
				return err
			}
			// Query string doesn't work for V4 APIs so for the moment
			// We get both pages
			var page model.Page
			if pages[0].Type == "MARKDOWN" {
				page = pages[0]
			} else {
				page = pages[1]
			}
			return assert.Equals("markdown content", updatedContent, page.Content)
		}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)
	})
})
