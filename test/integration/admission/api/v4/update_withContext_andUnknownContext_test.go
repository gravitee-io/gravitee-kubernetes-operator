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
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate update", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := v4.AdmissionCtrl{}

	It("should fail on API update with unknown management context", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4WithContextFile).
			Build().
			Apply()

		By("checking that API update does not pass validation")

		Consistently(func() error {
			newApi := fixtures.APIv4.DeepCopy()
			unknownContext := refs.NewNamespacedName("", "unknown")
			newApi.Spec.Context = &unknownContext

			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.APIv4, newApi)
			return err
		}, constants.ConsistentTimeout, interval).ShouldNot(Succeed())
	})
})
