// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Delete", labels.WithContext, func() {

	k8sClient := manager.Client()
	httpClient := http.Client{Timeout: 5 * time.Second}

	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should delete API V4 in cluster when not found in APIM", func() {
		fixtures := fixture.Builder().
			WithContext(constants.ContextWithSecretFile).
			WithAPIv4(constants.ApiV4WithContextFile).
			Build().Apply()

		By("deleting the API v4 in APIM")

		apim := apim.NewClient(ctx)

		err := apim.APIs.DeleteV4(fixtures.APIv4.Status.ID)
		Expect(err).ToNot(HaveOccurred())

		By("calling API v4 endpoint, expecting status 200")

		var endpoint = constants.BuildAPIV4Endpoint(fixtures.APIv4.Spec.Listeners[0])
		Eventually(func() error {
			res, callErr := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusOK)
		}, timeout, interval).Should(Succeed())

		By("deleting the API V4 Definition")

		err = k8sClient.Delete(ctx, fixtures.APIv4.DeepCopy())
		Expect(err).ToNot(HaveOccurred())

		By("expecting config map to be deleted")

		cm := &v1.ConfigMap{}
		Eventually(func() error {
			return k8sClient.Get(ctx, types.NamespacedName{
				Name:      fixtures.APIv4.Name,
				Namespace: fixtures.APIv4.Namespace,
			}, cm)
		}, timeout, interval).ShouldNot(Succeed())

		By("calling gateway endpoint, expecting status 404")

		Eventually(func() error {
			res, callErr := httpClient.Get(endpoint)
			return assert.NoErrorAndHTTPStatus(callErr, res, http.StatusNotFound)
		}, timeout, interval).Should(Succeed())

		By("expecting API V4 events to have been emitted")

		assert.EventsEmitted(fixtures.APIv4, "DeleteSucceeded", "DeleteStarted")
	})
})
