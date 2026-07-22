package sharedpolicygroups

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/policygroups"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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

var _ = Describe("Validate drift", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := admission.AdmissionCtrl{}

	It("should not drift on a simple update", func() {
		fixtures := fixture.
			Builder().
			WithSharedPolicyGroups(constants.SharedPolicyGroupsFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the shared policy group description")
		newSpg := fixtures.SharedPolicyGroup.DeepCopy()
		newSpg.Spec.Description = utils.ToReference("updated description")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.SharedPolicyGroup, newSpg)
			return err
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should detect drift", func() {
		fixtures := fixture.
			Builder().
			WithSharedPolicyGroups(constants.SharedPolicyGroupsFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the remote shared policy group description")

		newSpg := fixtures.SharedPolicyGroup.DeepCopy()
		newSpg.PopulateIDs(nil, true)

		newSpg.Spec.Description = utils.ToReference("remote updated description")

		apim := apim.NewClient(ctx)

		_, err := apim.SharedPolicyGroup.CreateOrUpdate(newSpg)
		Expect(err).ToNot(HaveOccurred())

		By("changing the CRD shared policy group description")

		newSpg.Spec.Description = utils.ToReference("local CRD description")

		By("checking that shared policy group validation returns error")

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.SharedPolicyGroup, newSpg)
			return assert.DriftDetected("description: \"local CRD description\" != \"remote updated description\"", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
