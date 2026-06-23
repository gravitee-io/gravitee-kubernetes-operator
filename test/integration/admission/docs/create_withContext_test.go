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

package docs

import (
	"context"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/docs"
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
	admissionCtrl := adm.AdmissionCtrl{}

	It("should return severe error when neither portalRef nor apiRef is set", func() {
		doc := fixture.
			Builder().
			WithDocumentation(constants.DocumentationPortalFile).
			Build()

		doc.Documentation.Spec.Portal = nil
		doc.Documentation.Spec.API = nil

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, doc.Documentation)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when both portalRef and apiRef are set", func() {
		doc := fixture.
			Builder().
			WithDocumentation(constants.DocumentationPortalFile).
			Build()

		doc.Documentation.Spec.API = &refs.NamespacedName{Name: "some-api"}

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, doc.Documentation)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when apiRef is not a v4 API kind", func() {
		doc := fixture.
			Builder().
			WithDocumentation(constants.DocumentationApiFile).
			Build()

		doc.Documentation.Spec.API.Kind = "ApiDefinition"

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, doc.Documentation)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when portalRef cannot be resolved", func() {
		doc := fixture.
			Builder().
			WithDocumentation(constants.DocumentationPortalFile).
			Build()

		doc.Documentation.Spec.Portal.Name = "unresolved"

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, doc.Documentation)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})
})
