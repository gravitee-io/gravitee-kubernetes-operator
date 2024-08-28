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

package v2

import (
	"context"

	apiv2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v2"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
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
	admissionCtrl := v2.AdmissionCtrl{}

	It("should return warning on API creation with unknown member", func() {
		fixtures := fixture.
			Builder().
			WithAPI(constants.Api).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("preparing API for import")

		fixtures.API.Spec.DefinitionContext = &apiv2.DefinitionContext{
			Origin:   apiv2.OriginKubernetes,
			Mode:     apiv2.ModeFullyManaged,
			SyncFrom: apiv2.OriginKubernetes,
		}
		fixtures.API.PopulateIDs(fixtures.Context)

		By("adding an unknown member to the API")

		unknownMemberName := random.GetName()

		fixtures.API.Spec.Members = []*base.Member{
			base.NewGraviteeMember(unknownMemberName, "REVIEWER"),
		}

		By("checking that API validation returns warnings")

		Eventually(func() error {
			warnings, err := admissionCtrl.ValidateUpdate(ctx, fixtures.API, fixtures.API)
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
