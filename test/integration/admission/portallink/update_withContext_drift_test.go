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

package portallink

import (
	"context"
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/portallink"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/labels"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	updatedName  = "updated link name"
	remoteName   = "remote updated link name"
	localCRDName = "local CRD link name"
	nameAssert   = `name: "%s" != "%s"`
)

var _ = Describe("Validate drift", labels.WithContext, func() {
	ctx := context.Background()
	admissionCtrl := admission.AdmissionCtrl{}

	It("should not drift with minimal fields", func() {
		fixtures := fixture.
			Builder().
			WithPortal(constants.PortalFile).
			WithPortalLink(constants.PortalLinkDriftMinimalFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the portal link name")
		newLink := fixtures.PortalLink.DeepCopy()
		setLinkName(newLink, updatedName)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.PortalLink, newLink)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should detect drift with minimal fields", func() {
		fixtures := fixture.
			Builder().
			WithPortal(constants.PortalFile).
			WithPortalLink(constants.PortalLinkDriftMinimalFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		validateNameDrift(ctx, admissionCtrl, fixtures)
	})

	It("should not drift with all fields", func() {
		fixtures := fixture.
			Builder().
			WithPortal(constants.PortalFile).
			WithPortalLink(constants.PortalLinkDriftFullFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the portal link name")
		newLink := fixtures.PortalLink.DeepCopy()
		setLinkName(newLink, updatedName)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.PortalLink, newLink)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should detect drift with all fields", func() {
		fixtures := fixture.
			Builder().
			WithPortal(constants.PortalFile).
			WithPortalLink(constants.PortalLinkDriftFullFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		validateNameDrift(ctx, admissionCtrl, fixtures)
	})

})

func setLinkName(link *v1alpha1.PortalLink, name string) {
	GinkgoHelper()
	link.Spec.Name = name
}

func validateNameDrift(
	ctx context.Context,
	admissionCtrl admission.AdmissionCtrl,
	fixtures *fixture.Objects,
) {
	GinkgoHelper()
	apimClient := apim.NewClient(ctx)

	oldLink := fixtures.PortalLink
	newLink := fixtures.PortalLink.DeepCopy()

	setLinkName(newLink, remoteName)

	_, err := apimClient.Links.CreateOrUpdate(newLink, fixtures.Portal)
	Expect(err).ToNot(HaveOccurred())

	setLinkName(newLink, localCRDName)

	Eventually(func() error {
		_, err := admissionCtrl.ValidateUpdate(ctx, oldLink, newLink)
		return assert.DriftDetected(
			fmt.Sprintf(nameAssert, localCRDName, remoteName),
			err,
		)
	}, constants.EventualTimeout, constants.Interval).Should(Succeed())
}
