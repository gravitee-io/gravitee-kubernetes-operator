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

package apidefinition_test

import (
	"context"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
)

var _ = Describe("Create", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()
	httpClient := http.Client{Timeout: 5 * time.Second}

	It("should update existing api in management API", func() {
		fixtures := fixture.Builder().
			WithContext(constants.ContextWithSecretFile).
			WithAPI(constants.ApiWithIds).
			Build()

		By("creating API in management api")

		apim := apim.NewClient(ctx)

		err := apim.APIs.DeleteV2(fixtures.API.Spec.ID)
		Expect(errors.IgnoreNotFound(err)).ToNot(HaveOccurred())

		Eventually(func() error {
			_, err = apim.APIs.ImportV2(http.MethodPost, &fixtures.API.Spec.Api)
			return err
		}, timeout, interval).Should(Succeed())

		Expect(err).ToNot(HaveOccurred())

		By("creating API in cluster")

		fixtures.Apply()

		By("calling management API and expecting API origin to be kubernetes")

		Eventually(func() error {
			api, apiErr := apim.APIs.GetByID(fixtures.API.Spec.ID)
			if apiErr != nil {
				return apiErr
			}
			return assert.Equals("API origin", "kubernetes", api.DefinitionContext.Origin)
		}, timeout, interval).Should(Succeed())
		By("calling gateway endpoint, expecting status 200")

		endpoint := constants.BuildAPIEndpoint(fixtures.API)
		Eventually(func() error {
			res, callErr := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusOK)
		}, timeout, interval).Should(Succeed())
	})
})
