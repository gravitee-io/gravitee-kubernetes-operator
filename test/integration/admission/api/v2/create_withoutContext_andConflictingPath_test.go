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

	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate create", labels.WithoutContext, func() {
	ctx := context.Background()
	admissionCtrl := adm.AdmissionCtrl{}

	It("should return error on API creation with conflicting path", func() {
		fixtures := fixture.
			Builder().
			WithAPIv4(constants.ApiV4).
			Build().
			Apply()

		By("building an API with path conflict")

		existingPath := fixtures.APIv4.GetContextPaths()[0]
		conflictingPath := existingPath

		api := fixture.Builder().WithAPI(constants.Api).Build().API
		api.Spec.Proxy.VirtualHosts[0].Path = conflictingPath

		By("checking that API creation does not pass validation")

		_, err := admissionCtrl.ValidateCreate(ctx, api)

		Expect(err).To(Equal(
			errors.NewSeveref(
				"invalid API context path [%s]. API [%s/%s] is already defined with the same path",
				conflictingPath, fixtures.APIv4.Namespace, fixtures.APIv4.Name,
			),
		))
	})
})
