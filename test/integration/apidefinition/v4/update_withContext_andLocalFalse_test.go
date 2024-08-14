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

package v4

import (
	"context"
	"net/http"
	"time"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

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

	It("should delete config map and sync from management API", func() {
		fixtures := fixture.Builder().
			WithAPIv4(constants.ApiV4).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting API V4 status to be completed")

		Expect(assert.ApiV4Completed(fixtures.APIv4)).To(Succeed())

		By("expecting to find config map")

		cm := &v1.ConfigMap{}
		Eventually(func() error {
			return manager.Client().Get(ctx, types.NamespacedName{
				Name:      fixtures.APIv4.Name,
				Namespace: fixtures.APIv4.Namespace,
			}, cm)
		}, timeout, interval).Should(Succeed())

		By("calling gateway endpoint, expecting status 200")

		endpoint := constants.BuildAPIV4Endpoint(fixtures.APIv4.Spec.Listeners[0])
		Eventually(func() error {
			res, callErr := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusOK)
		}, timeout, interval).Should(Succeed())

		By("updating the API V4, setting `local` to `false`")

		updated := fixtures.APIv4.DeepCopy()
		updated.Spec.DefinitionContext = &v4.DefinitionContext{
			Origin:   "KUBERNETES",
			SyncFrom: "MANAGEMENT",
		}

		Eventually(func() error {
			return manager.UpdateSafely(ctx, updated)
		}, timeout, interval).Should(Succeed())

		By("expecting config map to be deleted")

		Eventually(func() error {
			return manager.Client().Get(ctx, types.NamespacedName{
				Name:      fixtures.APIv4.Name,
				Namespace: fixtures.APIv4.Namespace,
			}, cm)
		}, timeout, interval).ShouldNot(Succeed())

		By("calling gateway endpoint, expecting status 200")

		Eventually(func() error {
			res, callErr := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusOK)
		}, timeout, interval).Should(Succeed())
	})
})
