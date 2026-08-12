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

package portal

import (
	"context"

	nav "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/navigation"
	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/portal"
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
	admissionCtrl := adm.AdmissionCtrl{}

	It("should return severe error when contextRef cannot be resolved", func() {
		prtl := fixture.
			Builder().
			WithPortal(constants.PortalFile).
			Build()

		prtl.Portal.Spec.Context.Name = "unresolved"

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, prtl.Portal)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when contextRef is missing", func() {
		prtl := fixture.
			Builder().
			WithPortal(constants.PortalFile).
			Build()

		prtl.Portal.Spec.Context = nil

		Eventually(func() error {
			_, err := admissionCtrl.ValidateCreate(ctx, prtl.Portal)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should warn when the deprecated navigation is used", func() {
		prtl := fixture.
			Builder().
			WithPortal(constants.PortalFile).
			Build()

		warnings, _ := admissionCtrl.ValidateCreate(ctx, prtl.Portal)

		Expect(warnings).To(ContainElement(ContainSubstring("spec.navigation is deprecated")))
	})

	It("should return severe error when navigation and structure are both set", func() {
		prtl := fixture.
			Builder().
			WithPortal(constants.PortalStructureFile).
			Build()

		prtl.Portal.Spec.Navigation = []*nav.NavigationPath{{Path: "/legacy"}}

		_, err := admissionCtrl.ValidateCreate(ctx, prtl.Portal)

		Expect(assert.Equals(
			"error",
			errors.NewSevere("navigation and structure cannot be used at the same time"),
			err,
		)).To(Succeed())
	})

	It("should not warn when only the structure is set", func() {
		prtl := fixture.
			Builder().
			WithPortal(constants.PortalStructureFile).
			Build()

		warnings, _ := admissionCtrl.ValidateCreate(ctx, prtl.Portal)

		Expect(warnings).ToNot(ContainElement(ContainSubstring("deprecated")))
	})
})
