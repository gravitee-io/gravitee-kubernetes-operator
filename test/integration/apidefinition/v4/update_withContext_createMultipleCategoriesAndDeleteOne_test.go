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
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should create a v4 API with multiple categories and then remove one category", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4).
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		By("create categories in APIM")
		categoryName1 := random.GetName()
		categoryName2 := random.GetName()
		categoryName3 := random.GetName()
		apim := apim.NewClient(ctx)
		Expect(apim.Env.CreateCategory(&model.Category{Name: categoryName1})).To(Succeed())
		Expect(apim.Env.CreateCategory(&model.Category{Name: categoryName2})).To(Succeed())
		Expect(apim.Env.CreateCategory(&model.Category{Name: categoryName3})).To(Succeed())

		By("applying the API with created category")
		fixtures.APIv4.Spec.Categories = []string{categoryName1, categoryName2, categoryName3}
		fixtures = fixtures.Apply()

		By("verifying that categories were created")
		expected3Categories := []string{categoryName1, categoryName2, categoryName3}
		Eventually(func() error {
			apiExport, err := apim.Export.V4Api(fixtures.APIv4.Status.ID)
			if err != nil {
				return err
			}

			exported3Categories := apiExport.Spec.Categories
			return assert.SliceEqualsSorted("categories", expected3Categories, exported3Categories, strings.Compare)
		}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)

		By("removing one category")
		fixtures.APIv4.Spec.Categories = []string{categoryName1, categoryName3}
		Eventually(func() error {
			return manager.UpdateSafely(ctx, fixtures.APIv4)
		}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)

		By("verifying that categories were updated")
		expected2Categories := []string{categoryName1, categoryName3}
		Eventually(func() error {
			updatedApiExport, err := apim.Export.V4Api(fixtures.APIv4.Status.ID)
			if err != nil {
				return err
			}

			exported2categories := updatedApiExport.Spec.Categories
			return assert.SliceEqualsSorted("categories", expected2Categories, exported2categories, strings.Compare)
		}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)
	})
})
