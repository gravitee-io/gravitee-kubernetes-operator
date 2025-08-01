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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/sharedpolicygroups"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	spg "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/policygroups"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate update", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := spg.AdmissionCtrl{}

	It("should return no error", func() {
		fixtures := fixture.
			Builder().
			WithSharedPolicyGroups(constants.SharedPolicyGroupsFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("having an invalid step")

		updated := *fixtures.SharedPolicyGroup
		configuration := utils.NewGenericStringMap()
		configuration.Put("key", "value")
		updated.Spec.Steps = []*sharedpolicygroups.Step{
			{
				Enabled:       true,
				Policy:        utils.ToReference("policy_throw_unexpected_policy_exception"),
				Configuration: configuration,
			},
		}

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.SharedPolicyGroup, &updated)
			return assert.Equals(
				"error",
				errors.NewSevere(
					"Plugin [policy_throw_unexpected_policy_exception] cannot be found.",
				),
				err,
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
