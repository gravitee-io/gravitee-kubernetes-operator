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

	apiV4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate update", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := v4.AdmissionCtrl{}

	DescribeTable("should return error on api update with invalid Api Type",
		func(apiName string, apiType apiV4.ApiType, errorMsg string) {
			fixtures := fixture.
				Builder().
				WithAPIv4(apiName).
				WithContext(constants.ContextWithCredentialsFile).
				Build().
				Apply()

			By("Updating API Type")
			api := fixtures.APIv4.DeepCopy()
			api.Spec.Type = apiType

			By("checking that api does not pass validation")

			Eventually(func() error {
				_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.APIv4, api)
				return assert.Equals(
					"error",
					errorMsg,
					err.Error(),
				)
			}, constants.ConsistentTimeout, interval).Should(Succeed())
		},
		Entry("should throw errors with wrong API Type",
			constants.NativeApiV4,
			apiV4.ApiType("PROXY"),
			"it is not possible to change the API type 'NATIVE' to something else [PROXY]",
		),
		Entry("should throw errors with with wrong API Type",
			constants.ApiV4,
			apiV4.ApiType("NATIVE"),
			"it is not possible to convert a none NATIVE API to a NATIVE API",
		),
	)
})
