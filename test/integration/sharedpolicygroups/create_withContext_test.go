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

package sharedpolicygroups

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", labels.WithContext, func() {
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("with a Shared Group Policy", func() {
		fixtures := fixture.Builder().
			WithContext(constants.ContextWithCredentialsFile).
			WithSharedPolicyGroups(constants.SharedPolicyGroupsFile).
			Build().
			Apply()

		By("expecting SPG status to be completed")

		Expect(assert.SharedPolicyGroupCompleted(fixtures.SharedPolicyGroup)).To(Succeed())

		By("calling rest API, expecting SPG to match status cross ID")

		apim := apim.NewClient(ctx)

		Eventually(func() error {
			spg, apiErr := apim.SharedPolicyGroup.GetByID(fixtures.SharedPolicyGroup.Status.ID)
			if apiErr != nil {
				return apiErr
			}
			return assert.Equals("SPG entity crossId", fixtures.SharedPolicyGroup.Status.CrossID, spg.CrossID)
		}, timeout, interval).ShouldNot(HaveOccurred(), fixtures.SharedPolicyGroup.Name)

		By("expecting SPG event to have been emitted")

		assert.EventsEmitted(fixtures.SharedPolicyGroup, "UpdateStarted", "UpdateSucceeded")
	},
	)
})
