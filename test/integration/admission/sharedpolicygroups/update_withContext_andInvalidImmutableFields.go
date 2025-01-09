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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/policygroups"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/uuid"

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

	DescribeTable("should return error modifying immutable fields",
		func(apply func(spg *v1alpha1.SharedPolicyGroup) *v1alpha1.SharedPolicyGroup, msg string) {
			fixtures := fixture.
				Builder().
				WithSharedPolicyGroups(constants.SharedPolicyGroupsFile).
				WithContext(constants.ContextWithCredentialsFile).
				Build().
				Apply()

			By("having an modified immutable field")

			group, _ := fixtures.SharedPolicyGroup.DeepCopyObject().(*v1alpha1.SharedPolicyGroup)
			updated := apply(group)

			Eventually(func() error {
				_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.SharedPolicyGroup, updated)
				return assert.Equals(
					"error",
					errors.NewSevere(fmt.Sprintf(msg, fixtures.SharedPolicyGroup.Status.CrossID)),
					err,
				)
			}, constants.EventualTimeout, interval).Should(Succeed())
		},
		Entry("should throw errors with modified ApiType",
			func(spg *v1alpha1.SharedPolicyGroup) *v1alpha1.SharedPolicyGroup {
				spg.Spec.ApiType = "MESSAGE"
				return spg
			},
			"can not change Shared Policy Group [%s] ApiType [PROXY], once it is created",
		),
		Entry("should throw errors with modified CrossID",
			func(spg *v1alpha1.SharedPolicyGroup) *v1alpha1.SharedPolicyGroup {
				id := uuid.NewV4String()
				spg.Spec.CrossID = &id
				return spg
			},
			"can not change Shared Policy Group [%s] CrossID, once it is created",
		),
		Entry("should throw errors with modified ApiType",
			func(spg *v1alpha1.SharedPolicyGroup) *v1alpha1.SharedPolicyGroup {
				spg.Spec.ApiType = "MESSAGE"
				return spg
			},
			"can not change Shared Policy Group [%s] ApiType [PROXY], once it is created",
		),
		Entry("should throw errors with modified Phase",
			func(spg *v1alpha1.SharedPolicyGroup) *v1alpha1.SharedPolicyGroup {
				phase := policygroups.FlowPhase("INTERACT")
				spg.Spec.Phase = &phase
				return spg
			},
			"can not change Shared Policy Group [%s] Phase [REQUEST], once it is created",
		),
	)
})
