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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v2"
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
	admissionCtrl := v2.AdmissionCtrl{}

	It("should return error on api creation with missing parent page", func() {
		fixtures1 := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		fixtures2 := fixture.
			Builder().
			WithAPI(constants.ApiWithPagesFile).
			Build()

		fixtures2.Context = fixtures1.Context

		By("checking that API does not pass validation")

		p := "missing-parent-page"
		(*fixtures2.API.Spec.Pages)["markdown"].Parent = &p

		Consistently(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures2.API)
			return assert.Equals(
				"error",
				errors.NewSeveref(
					"can not apply API [%s]. Parent page [missing-parent-page] can not be found for page [markdown]",
					fixtures2.API.Name,
				),
				err,
			)
		}, constants.ConsistentTimeout, interval).Should(Succeed())
	})
})
