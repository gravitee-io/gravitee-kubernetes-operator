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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/application"
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
	admissionCtrl := application.AdmissionCtrl{}

	It("should return error on application creation with missing management context", func() {
		fixtures := fixture.
			Builder().
			WithApplication(constants.Application).
			Build()

		By("checking that application does not pass validation")

		Consistently(func() error {
			unknownContext := refs.NewNamespacedName("", "unknown")
			fixtures.Application.Spec.Context = &unknownContext

			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.Application)
			return assert.Equals(
				"error",
				fmt.Errorf(
					"resource [%s] references management context [default/dev-ctx] that doesn't exist in the cluster",
					fixtures.Application.Name,
				),
				err,
			)
		}, constants.ConsistentTimeout, interval).ShouldNot(Succeed())
	})
})
