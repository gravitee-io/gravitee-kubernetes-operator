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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/policygroups"
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

var _ = Describe("Validate create", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := spg.AdmissionCtrl{}

	It("should return no error", func() {
		fixtures1 := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		fixtures2 := fixture.
			Builder().
			WithSharedPolicyGroups(constants.SharedPolicyGroupsFile).
			Build()

		By("having an invalid phase")

		fixtures2.SharedPolicyGroup.Spec.Context = fixtures1.Context.GetNamespacedName()

		uuid := uuid.NewV4String()
		fixtures2.SharedPolicyGroup.Spec.CrossID = &uuid
		fixtures2.SharedPolicyGroup.Spec.ApiType = "PROXY"
		fixtures2.SharedPolicyGroup.Spec.Phase = (*policygroups.FlowPhase)(utils.ToReference("SUBSCRIBE"))

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures2.SharedPolicyGroup)
			return assert.Equals(
				"error",
				errors.NewSevere(
					"Invalid phase SUBSCRIBE for API type PROXY",
				),
				err,
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
