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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/portal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	updatedName     = "updated name"
	remoteName      = "remote updated name"
	localCRDName    = "local CRD name"
	driftNameAssert = `name: "local CRD name" != "remote updated name"`
)

var _ = Describe("Validate drift", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := admission.AdmissionCtrl{}

	It("should not drift on with minimal fields", func() {
		fixtures := fixture.
			Builder().
			WithPortal(constants.PortalDriftMinimalFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the portal name")
		newPortal := fixtures.Portal.DeepCopy()
		setName(newPortal, updatedName)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Portal, newPortal)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should detect drift with minimal fields", func() {
		fixtures := fixture.
			Builder().
			WithPortal(constants.PortalDriftMinimalFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the remote portal name")
		newPortal := fixtures.Portal.DeepCopy()
		validateNameDrift(ctx, admissionCtrl, fixtures.Portal, newPortal)
	})

	It("should not drift on with all fields", func() {
		fixtures := fixture.
			Builder().
			WithPortal(constants.PortalDriftFullFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the portal name")
		newPortal := fixtures.Portal.DeepCopy()
		setName(newPortal, updatedName)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Portal, newPortal)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should detect drift with all fields", func() {
		fixtures := fixture.
			Builder().
			WithPortal(constants.PortalDriftFullFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the remote portal name")
		newPortal := fixtures.Portal.DeepCopy()
		validateNameDrift(ctx, admissionCtrl, fixtures.Portal, newPortal)
	})

})

func setName(prtl *v1alpha1.Portal, name string) {
	GinkgoHelper()
	prtl.Spec.Name = name
}

func validateNameDrift(
	ctx context.Context,
	admissionCtrl admission.AdmissionCtrl,
	oldPortal, newPortal *v1alpha1.Portal,
) {
	GinkgoHelper()
	apimClient := apim.NewClient(ctx)

	setName(newPortal, remoteName)

	_, err := apimClient.Portals.CreateOrUpdate(newPortal)
	Expect(err).ToNot(HaveOccurred())

	setName(newPortal, localCRDName)

	Eventually(func() error {
		_, err := admissionCtrl.ValidateUpdate(ctx, oldPortal, newPortal)
		return assert.DriftDetected(driftNameAssert, err)
	}, constants.EventualTimeout, constants.Interval).Should(Succeed())
}
