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

package drift

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/dictionary"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

const dynamicDictHRID = "dynamic-dict"

var _ = Describe("Dictionary Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			model.DictionaryDTO{},
			model.DictionaryDTO{},
		),
		Entry("equal manual struct",
			completeManualDictionaryDTO(),
			completeManualDictionaryDTO(),
		),
		Entry("equal dynamic struct from CRD mapping",
			completeDynamicDictionaryDTO(),
			model.ToDictionaryDTO(completeDynamicDictionarySpec(), dynamicDictHRID),
		),
		Entry("provider headers empty equivalent to nil",
			model.DictionaryDTO{
				Dynamic: &model.DynamicSpec{
					Provider: &model.Provider{},
				},
			},
			model.DictionaryDTO{
				Dynamic: &model.DynamicSpec{
					Provider: &model.Provider{
						Headers: []model.ProviderHeader{},
					},
				},
			},
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeDynamicDictionaryDTO(), completeDynamicDictionaryDTO())
		})
	})
})

func completeManualDictionaryDTO() model.DictionaryDTO {
	GinkgoHelper()
	return loadFixture[model.DictionaryDTO]("dictionary_manual_dto.json")
}

func completeDynamicDictionarySpec() dictionary.Type {
	GinkgoHelper()
	return loadFixture[dictionary.Type]("dictionary_dynamic_spec.json")
}

func completeDynamicDictionaryDTO() model.DictionaryDTO {
	GinkgoHelper()
	return loadFixture[model.DictionaryDTO]("dictionary_dynamic_dto.json")
}
