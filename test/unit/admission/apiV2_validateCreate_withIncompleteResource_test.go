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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate create", func() {

	ctx := context.Background()
	admissionCtrl := v2.AdmissionCtrl{}

	DescribeTable("with missing API resource property", func(
		mutate func(*v1alpha1.ApiDefinition), expectedMessage string,
	) {
		fixtures := fixture.
			Builder().
			WithAPI(constants.ApiWithOauth2AmResourceFile).
			Build()

		By("removing required property from API resource")

		mutate(fixtures.API)

		By("checking that API validation returns errors")

		warn, err := admissionCtrl.ValidateCreate(ctx, fixtures.API)

		Expect(warn).To(BeEmpty())
		Expect(err.Error()).To(Equal(expectedMessage))
	},
		Entry("type",
			func(api *v1alpha1.ApiDefinition) {
				api.Spec.Resources[0].Type = ""
			},
			"missing required string value in API resource property [type]",
		),
		Entry("name",
			func(api *v1alpha1.ApiDefinition) {
				api.Spec.Resources[0].Name = ""
			},
			"missing required string value in API resource property [name]",
		),
		Entry("configuration",
			func(api *v1alpha1.ApiDefinition) {
				api.Spec.Resources[0].Configuration = nil
			},
			"missing required object value in API resource property [configuration]",
		),
	)
})
