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

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
)

var _ = Describe("Create with portalNavigation", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should persist the portalNavigation tree in list order", func() {
		fixtures := fixture.Builder().
			WithAPIv4(constants.ApiV4WithContextFile).
			WithContext(constants.ContextWithSecretFile).
			Build()

		navigation := []*v4.NavigationPath{
			{Path: "/projects/alpha", DisplayName: utils.ToReference("Alpha"), Order: utils.ToReference(int32(1))},
			{Path: "/projects/alpha/docs"},
			{Path: "/projects/alpha/src"},
			{Path: "/projects/beta"},
			{Path: "/archive/2024"},
		}
		fixtures.APIv4.Spec.PortalNavigation = navigation

		fixtures.Apply()

		By("expecting API V4 status to be completed")

		Expect(assert.ApiV4Completed(fixtures.APIv4)).To(Succeed())
		Expect(assert.ApiV4Accepted(fixtures.APIv4)).To(Succeed())
		Expect(assert.ManagedByAutomationAPI(fixtures.APIv4)).To(Succeed())

		By("calling rest API, expecting portalNavigation to round-trip in the same order")

		client := apim.NewClient(ctx)
		hrid := refs.NewNamespacedNameFromObject(fixtures.APIv4).HRID()
		Eventually(func() error {
			api, apiErr := client.APIs.GetV4ByHRID(hrid)
			if apiErr != nil {
				return apiErr
			}
			paths := make([]string, 0, len(api.PortalNavigation))
			for _, nav := range api.PortalNavigation {
				paths = append(paths, nav.Path)
			}
			return assert.Equals("API V4 portalNavigation", []string{
				"/projects/alpha",
				"/projects/alpha/docs",
				"/projects/alpha/src",
				"/projects/beta",
				"/archive/2024",
			}, paths)
		}, timeout, interval).Should(Succeed(), fixtures.APIv4.Name)
	})
})
