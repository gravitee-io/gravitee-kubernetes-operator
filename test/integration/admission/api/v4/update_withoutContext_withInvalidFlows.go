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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"

	apiV4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate update", labels.WithoutContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := v4.AdmissionCtrl{}

	DescribeTable("should return error on api update with invalid Flow",
		func(apiName string, flow func() *apiV4.Flow, errorMsg string) {
			fixtures := fixture.
				Builder().
				WithAPIv4(apiName).
				Build().
				Apply()

			By("Updating API Plan")
			fixtures.APIv4.Spec.Flows = []*apiV4.Flow{flow()}

			By("checking that api does not pass validation")

			Eventually(func() error {
				_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.APIv4, fixtures.APIv4)
				return assert.Equals(
					"error",
					fmt.Sprintf(
						"%s [%s]",
						errorMsg,
						fixtures.APIv4.Name,
					),
					err.Error(),
				)
			}, constants.ConsistentTimeout, interval).Should(Succeed())
		},
		Entry("should throw errors with Connect flow step",
			constants.ApiV4,
			func() *apiV4.Flow {
				flow := apiV4.NewFlow("wrong")
				step := new(base.FlowStep)
				flow.Connect = []*apiV4.FlowStep{apiV4.NewFlowStep(*step)}
				return flow
			},
			"Connect Flow is not supported in V4 API",
		),
		Entry("should throw errors with Interact flow step",
			constants.ApiV4,
			func() *apiV4.Flow {
				flow := apiV4.NewFlow("wrong")
				step := new(base.FlowStep)
				flow.Interact = []*apiV4.FlowStep{apiV4.NewFlowStep(*step)}
				return flow
			},
			"Interact Flow is not supported in V4 API",
		),
		Entry("should throw errors with Request flow step",
			constants.NativeApiV4,
			func() *apiV4.Flow {
				flow := apiV4.NewFlow("wrong")
				step := new(base.FlowStep)
				flow.Request = []*apiV4.FlowStep{apiV4.NewFlowStep(*step)}
				return flow
			},
			"Request Flow is not supported in Native API",
		),
		Entry("should throw errors with Response flow step",
			constants.NativeApiV4,
			func() *apiV4.Flow {
				flow := apiV4.NewFlow("wrong")
				step := new(base.FlowStep)
				flow.Response = []*apiV4.FlowStep{apiV4.NewFlowStep(*step)}
				return flow
			},
			"Response Flow is not supported in Native API",
		),
	)
})
