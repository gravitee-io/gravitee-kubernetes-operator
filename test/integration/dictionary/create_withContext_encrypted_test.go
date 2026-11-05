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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create an encrypted dictionary", labels.WithContext, func() {
	ctx := context.Background()

	It("should resolve the secret and read the value back masked", func() {
		fixtures := fixture.
			Builder().
			AddSecret(constants.DictionarySecretFile).
			WithContext(constants.ContextWithCredentialsFile).
			WithDictionary(constants.DictionaryEncryptedFile).
			Build().
			Apply()

		apimClient, err := apim.FromContextRef(ctx, fixtures.Dictionary.ContextRef(), fixtures.Dictionary.GetNamespace())
		Expect(err).ToNot(HaveOccurred())

		hrid := refs.NewNamespacedNameFromObject(fixtures.Dictionary).HRID()
		remote, err := apimClient.Dictionaries.GetByHRID(hrid)

		Expect(err).ToNot(HaveOccurred())
		Expect(remote.Manual).ToNot(BeNil())
		Expect(remote.Manual.EncryptedProperties).To(HaveKey("apiKey"))
		Expect(remote.Manual.EncryptedProperties["apiKey"].Value).ToNot(BeEmpty())
		Expect(remote.Manual.EncryptedProperties["apiKey"].Value).ToNot(Equal("a secret value APIM stores encrypted"))
	})
})
