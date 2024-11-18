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

	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate update", labels.WithoutContext, func() {
	ctx := context.Background()
	admissionCtrl := adm.AdmissionCtrl{}

	It("should return error on API V4 update with missing configmap for templating", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4WithTemplatingFile).
			Build().
			Apply()

		By("updating plan security")

		(*fixtures.APIv4.Spec.Plans)["KeyLess"].Security.Type = "[[ configmap `graviteeio-templating/unknown` ]]"

		By("checking that API V4 update does not pass validation")

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.APIv4, fixtures.APIv4)

		Expect(err).To(Equal(
			errors.NewSeveref(
				"key [unknown] not found in configmap [default/graviteeio-templating]",
			),
		))
	})
})