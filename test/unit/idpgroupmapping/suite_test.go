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

package idpgroupmapping

import (
	"testing"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/idpgroupmapping"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIDPGroupMapping(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "package idpgroupmapping")
}

var _ = Describe("IsGroupExpression", func() {
	DescribeTable("detects values the Gravitee EL engine evaluates",
		func(value string, expected bool) {
			Expect(idpgroupmapping.IsGroupExpression(value)).To(Equal(expected))
		},
		Entry("jsonPath on the profile", "{#jsonPath(#profile, '$.groups')}", true),
		Entry("spaces before the hash", "{  #profile['groups']}", true),
		Entry("type reference", "{T(java.lang.String).valueOf(#profile['team'])}", true),
		Entry("parenthesis", "{(#profile['team'])}", true),
		Entry("expression after a prefix", "team-{#profile['team']}", true),
		Entry("plain group name", "developers", false),
		Entry("group name with dashes", "platform-admin", false),
		Entry("group id", "7b3e4a6c-1f2d-4c5b-9e8a-0d1c2b3a4f5e", false),
		Entry("empty value", "", false),
		Entry("brace not followed by an expression", "{team}", false),
		Entry("lowercase t is not a type reference", "{t(x)}", false),
	)
})
