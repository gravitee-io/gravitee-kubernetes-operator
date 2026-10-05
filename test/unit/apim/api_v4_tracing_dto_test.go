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

var _ = Describe("API v4 tracing payload", func() {
	It("sends the span attribute redaction rules to APIM", func() {
		api := &v4.Api{V4BaseApi: &v4.V4BaseApi{
			ApiBase: &base.ApiBase{},
			Analytics: &v4.Analytics{Enabled: true, Tracing: &v4.Tracing{
				Enabled: new(true),
				Redaction: &v4.TracingRedaction{
					DefaultReplacement: new("[MASKED]"),
					Rules: []v4.TracingRedactionRule{{
						AttributeNamePattern: "http.request.header.authorization",
						ValuePattern:         new("^Bearer "),
						MaskingStrategy: &v4.TracingMaskingStrategy{
							Type:         v4.TracingMaskingTypePartial,
							Replacement:  new("#"),
							PrefixLength: new(7),
							SuffixLength: new(2),
						},
					}},
				},
			}},
		}}

		payload, err := json.Marshal(model.ToAPIV4DTO(api))

		Expect(err).ToNot(HaveOccurred())
		var sent struct {
			Analytics struct {
				Tracing struct {
					Redaction json.RawMessage `json:"redaction"`
				} `json:"tracing"`
			} `json:"analytics"`
		}
		Expect(json.Unmarshal(payload, &sent)).To(Succeed())
		Expect(sent.Analytics.Tracing.Redaction).To(MatchJSON(`{
			"defaultReplacement": "[MASKED]",
			"rules": [{
				"attributeNamePattern": "http.request.header.authorization",
				"valuePattern": "^Bearer ",
				"maskingStrategy": {"type": "PARTIAL", "replacement": "#", "prefixLength": 7, "suffixLength": 2}
			}]
		}`))
	})
})
