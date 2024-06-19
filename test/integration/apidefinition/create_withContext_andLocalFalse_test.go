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

package apidefinition

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	xhttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/endpoint"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Create", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()
	httpCli := xhttp.NewNoAuthClient(ctx)

	It("should not create a config map and sync from management API", func() {
		fixtures := fixture.Builder().
			WithAPI(constants.ApiWithSyncFromAPIM).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting API status to be completed")

		Expect(assert.ApiCompleted(fixtures.API)).To(Succeed())

		By("expecting not to find config map")

		Eventually(func() error {
			cm := &v1.ConfigMap{}
			return manager.Client().Get(ctx, types.NamespacedName{
				Name:      fixtures.API.Name,
				Namespace: fixtures.API.Namespace,
			}, cm)
		}, timeout, interval).ShouldNot(Succeed())

		By("calling gateway endpoint, expecting status 200")

		Eventually(func() error {
			url := endpoint.ForV2(fixtures.API)
			return httpCli.Get(url, nil)
		}, timeout, interval).Should(Succeed())
	})
})
