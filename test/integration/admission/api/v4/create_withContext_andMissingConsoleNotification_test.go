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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v4"
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
	admissionCtrl := v4.AdmissionCtrl{}

	It("should return error on API creation with missing notification ref", func() {

		fixture.Builder().WithContext(constants.ContextWithCredentialsFile).Build().Apply()

		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4WithContextFile).Build()

		fixtures.APIv4.Spec.NotificationsRefs = []refs.NamespacedName{{Name: "missing", Namespace: constants.Namespace}}

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, fixtures.APIv4)

			return assert.Equals(
				"severe",
				errors.NewSeveref(
					"api references notification [%s] that does not exist in the cluster",
					fixtures.APIv4.Spec.NotificationsRefs[0].String(),
				).Error(),
				err.Error(),
			)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
