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

package portaltheme

import (
	"context"
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portaltheme"
	adm "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/portaltheme"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate update", labels.WithContext, func() {
	interval := constants.Interval
	ctx := context.Background()
	admissionCtrl := adm.AdmissionCtrl{}

	It("should accept an update of an existing theme", func() {
		fixtures := fixture.
			Builder().
			WithPortalTheme(constants.PortalThemeFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		updated := fixtures.PortalTheme.DeepCopy()
		setName(updated, updatedName)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.PortalTheme, updated)

		Expect(err).ToNot(HaveOccurred())
	})

	It("should return severe error when contextRef cannot be resolved", func() {
		fixtures := fixture.
			Builder().
			WithPortalTheme(constants.PortalThemeFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		updated := fixtures.PortalTheme.DeepCopy()
		updated.Spec.Context.Name = "unresolved"

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.PortalTheme, updated)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should return severe error when the update is rejected by the dry run", func() {
		fixtures := fixture.
			Builder().
			WithPortalTheme(constants.PortalThemeFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		updated := fixtures.PortalTheme.DeepCopy()
		updated.Spec.Name = ""
		updated.Spec.Definition = portaltheme.Definition{}

		Eventually(func() error {
			_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.PortalTheme, updated)
			return assert.NotNil("admission error", err)
		}, constants.EventualTimeout, interval).Should(Succeed())
	})

	It("should skip validation when the theme is being deleted", func() {
		fixtures := fixture.
			Builder().
			WithPortalTheme(constants.PortalThemeFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		// A theme on its way out is rejected by the dry run, but the update carrying
		// the deletion timestamp must still go through so the finalizer can be removed.
		deleting := fixtures.PortalTheme.DeepCopy()
		deleting.DeletionTimestamp = &metav1.Time{Time: time.Now()}
		deleting.Spec.Name = ""
		deleting.Spec.Definition = portaltheme.Definition{}

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.PortalTheme, deleting)

		Expect(err).ToNot(HaveOccurred())
	})
})
