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

package application

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Update", labels.WithContext, func() {
	timeout := constants.EventualTimeout

	interval := constants.Interval

	ctx := context.Background()

	It("should update application in APIM", func() {
		fixtures := fixture.Builder().
			WithApplication(constants.Application).
			WithContext(constants.ContextWithSecretFile).
			Build().
			Apply()

		By("expecting application status to be completed")

		Expect(assert.ApplicationCompleted(fixtures.Application)).To(Succeed())

		By("calling rest API, expecting to find application")

		apim := apim.NewClient(ctx)

		Eventually(func() error {
			app, appErr := apim.Applications.GetByID(fixtures.Application.Status.ID)
			if appErr != nil {
				return appErr
			}
			return assert.Equals("name", fixtures.Application.Spec.Name, app.Name)
		}, timeout, interval).Should(Succeed(), fixtures.Application.Name)

		By("calling rest API, expecting to find application metadata")
		Eventually(func() error {
			metadata, appErr := apim.Applications.GetMetadataByApplicationID(fixtures.Application.Status.ID)
			if appErr != nil {
				return appErr
			}
			return assert.SliceOfSize("metadata", *metadata, 2)
		}, timeout, interval).Should(Succeed(), fixtures.Application.Name)

		By("updating application name")

		updated := fixtures.Application.DeepCopy()
		updated.Spec.Name += "-updated"
		(*updated.Spec.MetaData)[0].Name = "test metadata update"

		Expect(manager.UpdateSafely(ctx, updated)).To(Succeed())

		By("calling rest API, expecting application to be up to date")

		Eventually(func() error {
			app, appErr := apim.Applications.GetByID(fixtures.Application.Status.ID)
			if appErr != nil {
				return appErr
			}
			return assert.Equals("name", updated.Spec.Name, app.Name)
		}, timeout, interval).Should(Succeed(), fixtures.Application.Name)

		By("calling rest API, expecting to find updated application metadata")
		Eventually(func() error {
			metadata, appErr := apim.Applications.GetMetadataByApplicationID(fixtures.Application.Status.ID)
			if appErr != nil {
				return appErr
			}
			return assert.SliceOfSize("metadata", *metadata, 2)
		}, timeout, interval).Should(Succeed(), fixtures.Application.Name)

	})
})
