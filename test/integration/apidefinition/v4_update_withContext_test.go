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

	"github.com/pkg/errors"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Update", labels.WithContext, func() {
	httpClient := http.Client{Timeout: 5 * time.Second}

	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should update api definition", func() {
		fixtures := fixture.Builder().
			WithAPIv4(constants.ApiV4).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("calling gateway endpoint, expecting status 200")

		endpoint := constants.BuildAPIV4Endpoint(fixtures.APIv4.Spec.Listeners[0])
		Eventually(func() error {
			res, callErr := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusOK)
		}, timeout, interval).Should(Succeed())

		By("updating api context path")

		updated := fixtures.APIv4.DeepCopy()
		Eventually(func() error {
			listener, ok := updated.Spec.Listeners[0].ToListener().(*v4.HttpListener)
			if !ok {
				return errors.Errorf("listener not of type *v4.HttpListener")
			}
			listener.Paths[0].Path += "-updated"
			updated.Spec.Listeners[0] = v4.ToGenericListener(listener)
			updated.Spec.Name += "-updated"
			return nil
		}, timeout, interval).Should(Succeed())

		updatedEndpoint := constants.BuildAPIV4Endpoint(updated.Spec.Listeners[0])

		Eventually(func() error {
			return manager.UpdateSafely(updated)
		}, timeout, interval).Should(Succeed())

		By("calling updated endpoint, expecting status 200")

		Eventually(func() error {
			res, callErr := httpClient.Get(updatedEndpoint)
			return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusOK)
		}, timeout, interval).Should(Succeed())

		By("calling rest API V4, expecting API to have been updated")

		apim := apim.NewClient(ctx)
		Eventually(func() error {
			api, cliErr := apim.APIs.GetV4ByID(fixtures.APIv4.Status.ID)
			if cliErr != nil {
				return cliErr
			}
			return assert.Equals("API name", updated.Spec.Name, api.Name)
		}, timeout, interval).Should(Succeed())

		By("expecting API event to have been emitted")

		assert.EventsEmitted(fixtures.APIv4, "UpdateStarted", "UpdateSucceeded")
	})
})
