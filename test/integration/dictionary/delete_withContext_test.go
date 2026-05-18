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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Delete", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval
	ctx := context.Background()

	It("should delete manual dictionary in APIM", func() {
		fixtures := fixture.Builder().
			AddSecret(constants.ContextSecretFile).
			WithDictionary(constants.DictionaryManualFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting dictionary status to be completed")

		Expect(assert.DictionaryAccepted(fixtures.Dictionary)).To(Succeed())

		By("calling rest API, expecting to find dictionary")

		apim := apim.NewClient(ctx)
		hrid := refs.NewNamespacedNameFromObject(fixtures.Dictionary).HRID()

		Eventually(func() error {
			dict, dictErr := apim.Dictionaries.GetByHRID(hrid)
			if dictErr != nil {
				return dictErr
			}
			return assert.NotEmptyString("id", dict.ID)
		}, timeout, interval).Should(Succeed(), fixtures.Dictionary.Name)

		By("deleting dictionary")

		Expect(manager.Client().Delete(ctx, fixtures.Dictionary.DeepCopy())).To(Succeed())

		By("expecting dictionary to be deleted from k8s")

		Eventually(func() error {
			return assert.Deleted(ctx, "Dictionary", fixtures.Dictionary)
		}, timeout, interval).Should(Succeed(), fixtures.Dictionary.Name)

		By("calling rest API, expecting not to find dictionary")

		Eventually(func() error {
			_, dictErr := apim.Dictionaries.GetByHRID(hrid)
			return assert.NotFoundError(dictErr)
		}, timeout, interval).Should(Succeed(), fixtures.Dictionary.Name)
	})
})
