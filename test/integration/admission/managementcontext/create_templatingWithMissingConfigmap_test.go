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

package managementcontext

import (
	"context"

	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/mctx"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate create", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := adm.AdmissionCtrl{}

	It("should return error on mCtx creation with missing configmap for templating", func() {
		fixtures := fixture.
			Builder().
			WithContext(constants.ContextWithCredentialsFile).
			Build()

		By("updating mCtx baseUrl")

		fixtures.Context.Spec.BaseUrl = "[[ configmap `missing-configmap/url` ]]"

		By("checking that mCtx creation does not pass validation")

		_, err := admissionCtrl.ValidateCreate(ctx, fixtures.Context)

		Expect(err).To(Equal(
			errors.NewSeveref(
				"configmap [default/missing-configmap] not found",
			),
		))
	})
})
