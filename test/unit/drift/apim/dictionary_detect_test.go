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
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/dictionary"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const dynamicDictHRID = "dynamic-dict"

const manualDictHRID = "manual-dict"

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
		Entry("equal manual struct from CRD mapping, encrypted properties carried through",
			completeManualDictionaryDTO(),
			model.ToDictionaryDTO(completeManualDictionarySpec(), manualDictHRID),
		),
		Entry("a masked remote secret matches the value in the manifest",
			manualDictionaryWithSecret("my secret API key"),
			manualDictionaryWithSecret(apimMask),
		),
		Entry("omitted encryption equivalent to APIM's default of false",
			model.DictionaryDTO{
				Dynamic: &model.DynamicSpec{Provider: &model.Provider{}},
			},
			model.DictionaryDTO{
				Dynamic: &model.DynamicSpec{
					Provider:   &model.Provider{},
					Encryption: &dictionary.EncryptionPolicy{EncryptOnFetch: false},
				},
			},
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

	DescribeTable("drifting values",
		func(crd, remote any) {
			result := drift.DetectWithNamespace(crd, remote, "")

			Expect(result.DriftDetected()).To(BeTrue())
		},
		Entry("the key is encrypted in APIM but plain in the manifest",
			model.DictionaryDTO{
				Manual: &model.ManualSpec{Properties: map[string]string{"apiKey": "my secret API key"}},
			},
			manualDictionaryWithSecret(apimMask),
		),
		Entry("the key is plain in APIM but encrypted in the manifest",
			manualDictionaryWithSecret("my secret API key"),
			model.DictionaryDTO{
				Manual: &model.ManualSpec{Properties: map[string]string{"apiKey": "my secret API key"}},
			},
		),
		Entry("encrypt on fetch turned on outside Git",
			model.DictionaryDTO{
				Dynamic: &model.DynamicSpec{Provider: &model.Provider{}},
			},
			model.DictionaryDTO{
				Dynamic: &model.DynamicSpec{
					Provider:   &model.Provider{},
					Encryption: &dictionary.EncryptionPolicy{EncryptOnFetch: true},
				},
			},
		),
	)

	Describe("Redaction through the encryptedProperties map", func() {
		It("names the key and hides the manifest's secret", func() {
			result := drift.DetectWithNamespace(
				manualDictionaryWithSecret("my secret API key"),
				model.DictionaryDTO{
					Manual: &model.ManualSpec{Properties: map[string]string{"apiKey": "now-plain-in-apim"}},
				},
				"",
			)
			report := result.String()

			Expect(result.DriftDetected()).To(BeTrue())
			Expect(report).To(ContainSubstring("apiKey"))
			Expect(report).NotTo(ContainSubstring("my secret API key"))
		})
	})

	Describe("Redaction survives a merge", func() {
		It("names the key and prints neither value on the admission path", func() {
			crd := manualDictionaryWithSecret("my secret API key")
			remote := model.DictionaryDTO{
				Manual: &model.ManualSpec{Properties: map[string]string{"apiKey": "now-plain-in-apim"}},
			}

			// Admission compares the old and the new CRD against the remote, then merges.
			merged := drift.Merge(
				drift.DetectWithNamespace(crd, remote, ""),
				drift.DetectWithNamespace(crd, remote, ""),
			)
			report := merged.String()

			Expect(merged.DriftDetected()).To(BeTrue())
			Expect(report).To(ContainSubstring("apiKey"))
			Expect(report).NotTo(ContainSubstring("my secret API key"))
			Expect(report).NotTo(ContainSubstring("now-plain-in-apim"))
		})
	})

	Describe("A key the manifest never declared", func() {
		It("is tolerated under properties", func() {
			merged := mergedDrift(
				model.DictionaryDTO{Manual: &model.ManualSpec{Properties: map[string]string{"env": "test"}}},
				model.DictionaryDTO{Manual: &model.ManualSpec{
					Properties: map[string]string{"env": "test", "addedInConsole": "hello"},
				}},
			)

			Expect(merged.DriftDetected()).To(BeFalse())
		})

		It("is tolerated under encryptedProperties", func() {
			merged := mergedDrift(
				model.DictionaryDTO{Manual: &model.ManualSpec{Properties: map[string]string{"env": "test"}}},
				model.DictionaryDTO{Manual: &model.ManualSpec{
					Properties:          map[string]string{"env": "test"},
					EncryptedProperties: map[string]dictionary.EncryptableValue{"addedInConsole": {Value: apimMask}},
				}},
			)

			Expect(merged.DriftDetected()).To(BeFalse())
		})
	})

	Describe("CRD mapping", func() {
		It("sends an empty properties map so a dictionary can still be cleared", func() {
			dto := model.ToDictionaryDTO(
				dictionary.Type{
					Name:           "manual-dict",
					DictionaryType: dictionary.ManualType,
					Manual:         &dictionary.ManualSpec{Properties: map[string]string{}},
				},
				manualDictHRID,
			)

			payload, err := json.Marshal(dto)

			Expect(err).ToNot(HaveOccurred())
			Expect(string(payload)).To(ContainSubstring(`"properties":{}`))
		})

		It("carries encrypt-on-fetch into the DTO", func() {
			dto := model.ToDictionaryDTO(completeDynamicDictionarySpec(), dynamicDictHRID)

			Expect(dto.Dynamic.Encryption).To(Equal(&dictionary.EncryptionPolicy{EncryptOnFetch: true}))
		})

		It("carries encrypted properties into the DTO", func() {
			dto := model.ToDictionaryDTO(completeManualDictionarySpec(), manualDictHRID)

			Expect(dto.Manual.EncryptedProperties).To(
				HaveKeyWithValue("apiKey", dictionary.EncryptableValue{Value: "my secret API key"}),
			)
		})
	})

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

const apimMask = "\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022"

// mergedDrift runs the comparison the admission path runs: old and new CRD against the remote.
func mergedDrift(crd, remote model.DictionaryDTO) drift.Result {
	return drift.Merge(
		drift.DetectWithNamespace(crd, remote, ""),
		drift.DetectWithNamespace(crd, remote, ""),
	)
}

func manualDictionaryWithSecret(value string) model.DictionaryDTO {
	return model.DictionaryDTO{
		Manual: &model.ManualSpec{
			EncryptedProperties: map[string]dictionary.EncryptableValue{"apiKey": {Value: value}},
		},
	}
}

func completeManualDictionarySpec() dictionary.Type {
	GinkgoHelper()
	return loadFixture[dictionary.Type]("dictionary_manual_spec.json")
}

func completeDynamicDictionarySpec() dictionary.Type {
	GinkgoHelper()
	return loadFixture[dictionary.Type]("dictionary_dynamic_spec.json")
}

func completeDynamicDictionaryDTO() model.DictionaryDTO {
	GinkgoHelper()
	return loadFixture[model.DictionaryDTO]("dictionary_dynamic_dto.json")
}
