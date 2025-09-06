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
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should create an API with a single category", func() {
		Skip(`
			This test was migrated and moved to e2e test suite
		`)
		fixtures := fixture.
			Builder().
			WithAPI(constants.Api).
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		By("create a category in APIM")
		categoryName := random.GetName()
		apim := apim.NewClient(ctx)
		Expect(apim.Env.CreateCategory(&model.Category{Name: categoryName})).To(Succeed())

		By("applying the API with a category")
		fixtures.API.Spec.Categories = []string{categoryName}
		fixtures = fixtures.Apply()

		By("verifying that the category was added")
		expectedCategories := []string{categoryName}

		Eventually(func() error {
			apiExport, err := apim.Export.V2Api(fixtures.API.Status.ID)
			if err != nil {
				return err
			}

			exportedCategories := apiExport.Spec.Categories

			return assert.Equals("categories", expectedCategories, exportedCategories)
		}, timeout, interval).Should(Succeed(), fixtures.API.Name)
	})
})
