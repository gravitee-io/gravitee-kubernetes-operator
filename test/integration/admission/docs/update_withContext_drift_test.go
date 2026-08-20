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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/docs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/service"
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
			WithPortal(constants.PortalDriftFullFile).
			WithDocumentation(constants.DocumentationDriftMinimalFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the documentation name")
		newDoc := fixtures.Documentation.DeepCopy()
		setName(newDoc, updatedName)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Documentation, newDoc)
		Expect(err).ToNot(HaveOccurred())

	})

	It("should detect drift with minimal fields", func() {
		fixtures := fixture.
			Builder().
			WithPortal(constants.PortalDriftFullFile).
			WithDocumentation(constants.DocumentationDriftMinimalFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the remote documentation name")
		validateNameDrift(ctx, admissionCtrl, fixtures)
	})

	It("should not drift on with all fields", func() {
		fixtures := fixture.
			Builder().
			WithPortal(constants.PortalDriftFullFile).
			WithDocumentation(constants.DocumentationDriftFullFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the documentation name")
		newDoc := fixtures.Documentation.DeepCopy()
		setName(newDoc, updatedName)

		_, err := admissionCtrl.ValidateUpdate(ctx, fixtures.Documentation, newDoc)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should detect drift with all fields", func() {
		fixtures := fixture.
			Builder().
			WithPortal(constants.PortalDriftFullFile).
			WithDocumentation(constants.DocumentationDriftFullFile).
			WithContext(constants.ContextWithCredentialsFile).
			Build().
			Apply()

		By("changing the remote documentation name")
		validateNameDrift(ctx, admissionCtrl, fixtures)
	})
})

func setName(doc *v1alpha1.Documentation, name string) {
	GinkgoHelper()
	doc.Spec.Name = name
}

func validateNameDrift(
	ctx context.Context,
	admissionCtrl admission.AdmissionCtrl,
	fixtures *fixture.Objects,
) {
	GinkgoHelper()
	apimClient := apim.NewClient(ctx)

	oldDoc := fixtures.Documentation.DeepCopy()
	newDoc := fixtures.Documentation.DeepCopy()
	parent := service.DocumentationParent{Portal: fixtures.Portal}

	setName(newDoc, remoteName)

	_, err := apimClient.Documentations.CreateOrUpdate(newDoc, parent)
	Expect(err).ToNot(HaveOccurred())

	// revert to the way it was in the CRD
	newDoc.Spec.Portal.Namespace = ""

	setName(newDoc, localCRDName)

	Eventually(func() error {
		_, err := admissionCtrl.ValidateUpdate(ctx, oldDoc, newDoc)
		return assert.DriftDetected(driftNameAssert, err)
	}, constants.ConsistentTimeout, constants.Interval).Should(Succeed())
}
