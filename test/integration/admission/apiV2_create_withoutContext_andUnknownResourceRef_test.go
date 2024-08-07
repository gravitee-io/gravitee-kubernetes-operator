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
	"fmt"

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate create", labels.WithoutContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := v2.AdmissionCtrl{}

	It("should return error on api creation with missing resource as ref", func() {
		fixtures := fixture.
			Builder().
			WithAPI(constants.ApiWithCacheRedisResourceRefFile).
			Build()

		By("checking that api does not pass validation")

		Consistently(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.API)
			return assert.Equals(
				"error",
				fmt.Sprintf(
					"api references resource [%s] that does not exist in the cluster",
					fixtures.API.Spec.Resources[0].Ref,
				),
				err.Error(),
			)
		}, constants.ConsistentTimeout, interval).Should(Succeed())
	})
})
