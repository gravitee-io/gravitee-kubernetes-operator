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
	admissionv4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"
)

var _ = Describe("Validate create", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := admissionv4.AdmissionCtrl{}

	It("should return error on api creation with missing Shared Policy Group", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4WithContextFile).
			Build()

		fixtures.APIv4.Spec.Flows = []*v4.Flow{
			{
				Enabled: ptr.To(true),
				Request: []*v4.FlowStep{
					{
						SharedPolicyGroup: &refs.NamespacedName{
							Name: "missing-shared-policy-group",
						},
					},
				},
			},
		}

		By("checking that API does not pass validation")

		Consistently(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.APIv4)
			return assert.Equals(
				"error",
				errors.NewSeveref(
					"unable to get Shared Policy Group [missing-shared-policy-group] in namespace [%s]",
					fixtures.APIv4.Namespace,
				),
				err,
			)
		}, constants.ConsistentTimeout, interval).Should(Succeed())
	})
})
