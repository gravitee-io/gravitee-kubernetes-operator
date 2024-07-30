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

package admission

import (
	"context"

	apiV4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate update", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := v4.AdmissionCtrl{}

	It("should return warning on API creation with unknown category", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("preparing API for import")

		fixtures.APIv4.Spec.DefinitionContext = apiV4.NewDefaultKubernetesContext()
		fixtures.APIv4.PopulateIDs(fixtures.Context)

		By("adding an unknown category to the API")

		unknownCategory := random.GetName()

		fixtures.APIv4.Spec.Categories = []string{unknownCategory}

		By("checking that API validation returns warnings")

		Eventually(func() error {
			warnings, err := admissionCtrl.ValidateUpdate(ctx, fixtures.APIv4, fixtures.APIv4)
			if err != nil {
				return err
			}
			if err = assert.SliceOfSize("warnings", warnings, 1); err != nil {
				return err
			}
			return assert.Equals(
				"warning",
				errors.NewWarning("category [%s] is not defined in environment [DEFAULT]", unknownCategory).Error(),
				warnings[0],
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
