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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stop", labels.WithContext, func() {
	httpClient := http.Client{Timeout: 5 * time.Second}

	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should start API", func() {
		fixtures := fixture.Builder().
			WithAPI(constants.BasicApiFile).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		endpoint := constants.BuildAPIEndpoint(fixtures.API)

		By("calling gateway endpoint, expecting status 200")

		Eventually(func() error {
			res, callErr := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusOK)
		}, timeout, interval).ShouldNot(HaveOccurred())

		By("calling rest API, expecting state 'STARTED'")

		apim := apim.NewClient(ctx)
		Eventually(func() error {
			api, apiErr := apim.APIs.GetByID(fixtures.API.Status.ID)
			if apiErr != nil {
				return apiErr
			}
			return assert.Equals("state", "STARTED", api.State)
		}, timeout, interval).Should(Succeed())

		By("updating the API, setting state to 'STOPPED'")

		updated := fixtures.API.DeepCopy()
		updated.Spec.State = base.StateStopped

		Eventually(func() error {
			return manager.UpdateSafely(updated)
		}, timeout, interval).ShouldNot(HaveOccurred())

		By("calling gateway endpoint, expecting status 404")

		Eventually(func() error {
			res, callErr := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusNotFound)
		}, timeout, interval).ShouldNot(HaveOccurred())

		By("calling rest API, expecting state 'STOPPED'")

		Eventually(func() error {
			api, apiErr := apim.APIs.GetByID(fixtures.API.Status.ID)
			if apiErr != nil {
				return apiErr
			}
			return assert.Equals("state", "STOPPED", api.State)
		}, timeout, interval).Should(Succeed())

	})
})
