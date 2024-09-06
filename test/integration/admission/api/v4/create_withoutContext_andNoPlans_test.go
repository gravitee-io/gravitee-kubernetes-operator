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

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate create", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := v4.AdmissionCtrl{}

	It("should return error on API creation with no plans and state STARTED", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4).
			Build()

		By("removing existing plans")

		clear(fixtures.APIv4.Spec.Plans)

		By("checking that API validation returns errors")

		_, err := admissionCtrl.ValidateCreate(ctx, fixtures.APIv4)

		Expect(err).To(Equal(
			errors.NewSeveref("cannot apply API [%s]. Its state is set to STARTED, "+
				"but the API has no plans. APIs must have at least one plan in order to "+
				"be deployed.",
				fixtures.APIv4.Name,
			),
		))
	})
})
