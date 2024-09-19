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

package application

import (
	"context"

	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/application"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate update", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := adm.AdmissionCtrl{}

	It("should return error on Application update with missing secrete for templating", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.Application).
			Build().
			Apply()

		By("updating app domain")

		fixtures.Application.Spec.Domain = "[[ secret `missing-secret/domain` ]]"

		By("checking that Application update does not pass validation")

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Application, fixtures.Application)

		Expect(err).To(Equal(
			errors.NewSeveref(
				"secret [default/missing-secret] not found",
			),
		))
	})
})
