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

package v2

import (
	"context"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Update", labels.WithoutContext, func() {
	httpClient := http.Client{Timeout: 5 * time.Second}

	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should update api definition", func() {
		fixtures := fixture.Builder().
			WithAPI(constants.Api).
			Build().
			Apply()

		By("calling gateway endpoint, expecting status 200")

		endpoint := constants.BuildAPIEndpoint(fixtures.API)
		Eventually(func() error {
			res, callErr := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusOK)
		}, timeout, interval).Should(Succeed())

		By("updating api context path")

		updated := fixtures.API.DeepCopy()
		updated.Spec.Proxy.VirtualHosts[0].Path += "-updated"
		updatedEndpoint := constants.BuildAPIEndpoint(updated)

		Eventually(func() error {
			return manager.UpdateSafely(ctx, updated)
		}, timeout, interval).Should(Succeed())

		By("calling updated endpoint, expecting status 200")

		Eventually(func() error {
			res, callErr := httpClient.Get(updatedEndpoint)
			return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusOK)
		}, timeout, interval).Should(Succeed())

		By("expecting API event to have been emitted")

		assert.EventsEmitted(fixtures.API, "UpdateStarted", "UpdateSucceeded")
	})
})
