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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("API v4 automation payload", func() {
	It("sends the failover failure condition and forced next endpoint to APIM", func() {
		spec := v1alpha1.ApiV4DefinitionSpec{Api: v4.Api{V4BaseApi: &v4.V4BaseApi{
			ApiBase: &base.ApiBase{},
			Failover: &v4.Failover{
				Enabled:                    new(true),
				FailureCondition:           new("{#response.status >= 500}"),
				ForceNextEndpointOnFailure: new(true),
			},
		}}}

		payload, err := json.Marshal(model.ToAutomation(spec))

		Expect(err).ToNot(HaveOccurred())
		var sent struct {
			Failover map[string]any `json:"failover"`
		}
		Expect(json.Unmarshal(payload, &sent)).To(Succeed())
		Expect(sent.Failover).To(HaveKeyWithValue("failureCondition", "{#response.status >= 500}"))
		Expect(sent.Failover).To(HaveKeyWithValue("forceNextEndpointOnFailure", true))
	})
})
