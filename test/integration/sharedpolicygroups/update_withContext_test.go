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
	timeout := constants.EventualTimeout
	interval := constants.Interval

	ctx := context.Background()

	It("should update api definition", func() {
		fixtures := fixture.Builder().
			WithSharedPolicyGroups(constants.SharedPolicyGroupsFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("updating SPG description")

		updated := fixtures.SharedPolicyGroup.DeepCopy()
		updated.Spec.Name += "-updated"

		Eventually(func() error {
			return manager.UpdateSafely(ctx, updated)
		}, timeout, interval).Should(Succeed())

		By("calling rest API, expecting SPG to have been updated")

		apim := apim.NewClient(ctx)
		Eventually(func() error {
			spg, cliErr := apim.SharedPolicyGroup.GetByID(fixtures.SharedPolicyGroup.Status.ID)
			if cliErr != nil {
				return cliErr
			}
			return assert.Equals("SPG name", updated.Spec.Name, spg.Name)
		}, timeout, interval).Should(Succeed())

		By("expecting SPG event to have been emitted")

		assert.EventsEmitted(fixtures.SharedPolicyGroup, "UpdateStarted", "UpdateSucceeded")
	})
})
