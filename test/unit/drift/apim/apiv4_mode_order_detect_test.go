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

package apim

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

// GKO-3433: values APIM stores differently from what the CRD sends are not drift.
var _ = Describe("API v4 Drift detection on what APIM normalises", func() {
	DescribeTable("no drift",
		func(crd, remote string) {
			crdDTO := crdDriftAPIV4DTO(crd)
			remoteDTO := parseJSON[model.APIV4DTO](remote).Canonicalize()
			expectNoDrift(drift.DetectWithNamespace(crdDTO, remoteDTO, ""))
		},
		Entry("lowercase default flow mode against the uppercase APIM stores",
			`{"name": "api", "flowExecution": {"mode": "default"}}`,
			`{"name": "api", "flowExecution": {"mode": "DEFAULT"}}`,
		),
		Entry("lowercase best_match flow mode against the uppercase APIM stores",
			`{"name": "api", "flowExecution": {"mode": "best_match"}}`,
			`{"name": "api", "flowExecution": {"mode": "BEST_MATCH"}}`,
		),
		Entry("plan order 1 against the only plan APIM lists",
			`{"name": "api", "plans": {"JWTPlan": `+apiV4JWTPlanWithOrder+`}}`,
			`{"name": "api", "plans": [`+apiV4JWTPlan+`]}`,
		),
	)
})

const apiV4JWTPlanWithOrder = `{"hrid": "JWTPlan", "name": "JWT plan", "status": "PUBLISHED", "order": 1, ` +
	`"security": {"type": "JWT"}}`
