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

package apim_test

import (
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("API v4 payload", func() {
	It("sends the failover failure condition and forced next endpoint to APIM", func() {
		api := &v4.Api{V4BaseApi: &v4.V4BaseApi{
			ApiBase: &base.ApiBase{},
			Failover: &v4.Failover{
				Enabled:                    new(true),
				FailureCondition:           new("{#response.status >= 500}"),
				ForceNextEndpointOnFailure: new(true),
			},
		}}

		payload, err := json.Marshal(model.ToAPIV4DTO(api))

		Expect(err).ToNot(HaveOccurred())
		var sent struct {
			Failover map[string]any `json:"failover"`
		}
		Expect(json.Unmarshal(payload, &sent)).To(Succeed())
		Expect(sent.Failover).To(HaveKeyWithValue("failureCondition", "{#response.status >= 500}"))
		Expect(sent.Failover).To(HaveKeyWithValue("forceNextEndpointOnFailure", true))
	})

	DescribeTable("sends flowExecution.mode in the Automation API spelling",
		func(mode v4.FlowMode, expected v4.FlowMode) {
			api := &v4.Api{V4BaseApi: &v4.V4BaseApi{
				ApiBase:       &base.ApiBase{},
				FlowExecution: &v4.FlowExecution{Mode: mode},
			}}

			dto := model.ToAPIV4DTO(api)

			Expect(dto.FlowExecution.Mode).To(Equal(expected))
			Expect(api.FlowExecution.Mode).To(Equal(mode))
		},
		Entry("best-match", v4.FlowMode("best-match"), v4.FlowModeBestMatch),
		Entry("BEST-MATCH", v4.FlowMode("BEST-MATCH"), v4.FlowModeBestMatch),
		Entry("Best-Match", v4.FlowMode("Best-Match"), v4.FlowModeBestMatch),
		Entry("best_match", v4.FlowMode("best_match"), v4.FlowModeBestMatch),
		Entry("BEST_MATCH", v4.FlowModeBestMatch, v4.FlowModeBestMatch),
		Entry("default", v4.FlowMode("default"), v4.FlowModeDefault),
		Entry("DEFAULT", v4.FlowModeDefault, v4.FlowModeDefault),
		Entry("unset", v4.FlowMode(""), v4.FlowMode("")),
	)
})
