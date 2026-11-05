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

package dictionary

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/dictionary"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate manual properties", labels.WithContext, func() {
	ctx := context.Background()

	It("should reject a key declared in both property maps", func() {
		dict := fixture.Builder().WithDictionary(constants.DictionaryManualFile).Build()
		dict.Dictionary.Spec.Manual = &dictionary.ManualSpec{
			Properties:          map[string]string{"apiKey": "plain"},
			EncryptedProperties: map[string]dictionary.EncryptableValue{"apiKey": {Value: "secret"}},
		}

		err := manager.Client().Create(ctx, dict.Dictionary)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("a key must not appear in both properties and encryptedProperties"))
	})

	It("should reject a manual dictionary declaring neither map", func() {
		dict := fixture.Builder().WithDictionary(constants.DictionaryManualFile).Build()
		dict.Dictionary.Spec.Manual = &dictionary.ManualSpec{}

		err := manager.Client().Create(ctx, dict.Dictionary)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("a manual dictionary must declare properties or encryptedProperties"))
	})

	It("should reject an empty encrypted value", func() {
		dict := fixture.Builder().WithDictionary(constants.DictionaryManualFile).Build()
		dict.Dictionary.Spec.Manual = &dictionary.ManualSpec{
			EncryptedProperties: map[string]dictionary.EncryptableValue{"apiKey": {Value: ""}},
		}

		err := manager.Client().Create(ctx, dict.Dictionary)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("should be at least 1 chars long"))
	})
})
