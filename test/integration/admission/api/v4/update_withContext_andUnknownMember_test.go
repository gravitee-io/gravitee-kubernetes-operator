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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apidefinition"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate create", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := v4.AdmissionCtrl{}

	It("should return warning on API creation with unknown member", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		newAPI := fixtures.APIv4.DeepCopy()

		By("adding an unknown member to the API")
		unknownMemberName := random.GetName()
		newAPI.Spec.Members = []*base.Member{
			base.NewGraviteeMember(unknownMemberName, "REVIEWER"),
		}

		By("preparing API for import")
		err := apidefinition.PrepareV4SpecForAutomation(ctx, newAPI)
		Expect(err).ToNot(HaveOccurred())

		By("checking that API validation returns warnings")

		Eventually(func() error {
			warnings, err := admissionCtrl.ValidateUpdate(ctx, fixtures.APIv4, newAPI)
			if err != nil {
				return err
			}
			if err = assert.SliceOfSize("warnings", warnings, 1); err != nil {
				return err
			}
			return assert.Equals(
				"warning",
				errors.NewWarningf(
					"member [%s] of source [gravitee] could not be found in organization [DEFAULT]",
					unknownMemberName,
				).Error(),
				warnings[0],
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
