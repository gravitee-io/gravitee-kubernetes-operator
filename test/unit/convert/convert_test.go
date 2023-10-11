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

package convert

import (
	"testing"

	convertV2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/convert/v2"
	convertV4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/convert/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/unit/convert/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConvert(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core model conversion")
}

var _ = Context("When converting", func() {

	DescribeTable("from v2 to v4",
		func(path string) {
			in := internal.UnmarshalV2(path)
			converted := convertV4.FromV2(in)
			given := internal.Marshal(converted)
			expected := internal.ReadJSON(path + "/v4.json")
			Expect(given).To(MatchJSON(expected))
		},
		Entry("should convert virtual hosts to listeners", "entrypoints"),
		Entry("should convert endpoints", "endpoints"),
		Entry("should convert plans", "plans"),
		Entry("should convert flows", "flows"),
		Entry("should convert analytics", "analytics"),
		Entry("should convert services", "services"),
	)

	DescribeTable("from v4 to v2",
		func(path string) {
			in := internal.UnmarshalV4(path)
			converted := convertV2.FromV4(in)
			given := internal.Marshal(converted)
			expected := internal.ReadJSON(path + "/v2.json")
			Expect(given).To(MatchJSON(expected))
		},
		Entry("should convert listeners to virtual hosts", "entrypoints"),
		Entry("should convert endpoints", "endpoints"),
		Entry("should convert plans", "plans"),
		Entry("should convert flows", "flows"),
		Entry("should convert analytics", "analytics"),
		Entry("should convert services", "services"),
	)
})
