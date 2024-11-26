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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	tHTTP "github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Start", labels.WithoutContext, func() {
	httpClient := tHTTP.NewClient()

	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should start API", func() {
		fixtures := fixture.Builder().
			WithAPI(constants.ApiWithStateStopped).
			Build().
			Apply()

		endpoint := constants.BuildAPIEndpoint(fixtures.API)

		By("calling gateway endpoint, expecting status 404")

		Eventually(func() error {
			res, callErr := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusNotFound)
		}, timeout, interval).Should(Succeed())

		By("updating the API, setting state to `STARTED`")

		updated := fixtures.API.DeepCopy()
		updated.Spec.State = base.StateStarted

		Eventually(func() error {
			return manager.UpdateSafely(ctx, updated)
		}, timeout, interval).Should(Succeed())

		By("calling gateway endpoint, expecting status 200")

		Eventually(func() error {
			res, callErr := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusOK)
		}, timeout, interval).Should(Succeed())

	})
})
