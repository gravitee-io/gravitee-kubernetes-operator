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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	tHTTP "github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Update", labels.WithContext, func() {
	httpClient := tHTTP.NewClient()

	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should delete config map and sync from management API", func() {
		fixtures := fixture.Builder().
			WithAPI(constants.Api).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting API status to be completed")

		Expect(assert.ApiCompleted(fixtures.API)).To(Succeed())

		By("expecting to find config map")

		cm := &v1.ConfigMap{}
		Eventually(func() error {
			return manager.Client().Get(ctx, types.NamespacedName{
				Name:      fixtures.API.Name,
				Namespace: fixtures.API.Namespace,
			}, cm)
		}, timeout, interval).Should(Succeed())

		By("calling gateway endpoint, expecting status 200")

		endpoint := constants.BuildAPIEndpoint(fixtures.API)
		Eventually(func() error {
			res, callErr := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusOK)
		}, timeout, interval).Should(Succeed())

		By("updating the API, setting `local` to `false`")

		updated := fixtures.API.DeepCopy()
		updated.Spec.IsLocal = false

		// we unfortunately rely on a side effect in admission controller
		// to delete the config map
		admCtrl := v2.AdmissionCtrl{}
		_, err := admCtrl.ValidateUpdate(ctx, fixtures.API, updated)
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() error {
			return manager.UpdateSafely(ctx, updated)
		}, timeout, interval).Should(Succeed())

		By("expecting config map to be deleted")

		Eventually(func() error {
			return manager.Client().Get(ctx, types.NamespacedName{
				Name:      fixtures.API.Name,
				Namespace: fixtures.API.Namespace,
			}, cm)
		}, timeout, interval).ShouldNot(Succeed())

		By("calling gateway endpoint, expecting status 200")

		Eventually(func() error {
			res, callErr := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusOK)
		}, timeout, interval).Should(Succeed())
	})
})
