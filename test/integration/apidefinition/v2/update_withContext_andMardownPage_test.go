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

package v2

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
			WithAPI(constants.ApiWithMarkdownPage).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		apim := apim.NewClient(ctx)
		apiID := fixtures.API.Status.ID

		By("updating markdown content")

		markdown := (*fixtures.API.Spec.Pages)["markdown"]
		updatedContent := *markdown.Content + "\n" + "This is an update."
		markdown.Content = &updatedContent

		Eventually(func() error {
			return manager.UpdateSafely(ctx, fixtures.API)
		}, timeout, interval).Should(Succeed(), fixtures.API.Name)

		By("checking pages number in APIM")

		Eventually(func() error {
			pages, err := apim.Pages.FindByAPI(apiID)
			if err != nil {
				return err
			}
			return assert.SliceOfSize("pages", pages, 2)
		}, timeout, interval).Should(Succeed(), fixtures.API.Name)

		By("checking swagger content in APIM")

		Expect(markdown).ToNot(BeNil())

		Eventually(func() error {
			pages, err := apim.Pages.FindByAPI(apiID, model.NewPageQuery().WithType("MARKDOWN"))
			if err != nil {
				return err
			}
			if err = assert.SliceOfSize("pages", pages, 1); err != nil {
				return err
			}
			page := pages[0]
			return assert.Equals("markdown content", updatedContent, page.Content)
		}, timeout, interval).Should(Succeed(), fixtures.API.Name)
	})
})
