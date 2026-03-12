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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
)

var _ = Describe("Create", labels.WithContext, func() {
	DescribeTable("AI type APIs with a management context",
		func(builder *fixture.FSBuilder) {
			fixtures := builder.Build()
			fixtures.Apply()

			By("expecting API V4 status to be completed")
			Expect(assert.ApiV4Completed(fixtures.APIv4)).To(Succeed())

			By("expecting API V4 to be accepted")
			Expect(assert.ApiV4Accepted(fixtures.APIv4)).To(Succeed())

			By("expecting API V4 event to have been emitted")
			assert.EventsEmitted(fixtures.APIv4, "UpdateStarted", "UpdateSucceeded")
		},
		Entry("should accept LLM_PROXY type",
			fixture.Builder().
				WithAPIv4(constants.ApiV4LLMProxy).
				WithContext(constants.ContextWithSecretFile),
		),
		Entry("should accept MCP_PROXY type",
			fixture.Builder().
				WithAPIv4(constants.ApiV4MCPProxy).
				WithContext(constants.ContextWithSecretFile),
		),
		Entry("should accept A2A_PROXY type",
			fixture.Builder().
				WithAPIv4(constants.ApiV4A2AProxy).
				WithContext(constants.ContextWithSecretFile),
		),
	)
})
