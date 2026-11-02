// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package amapplication

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amapplication"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Validate create", func() {
	ctx := context.Background()
	admissionCtrl := amapplication.NewAdmissionCtrl()

	It("should pass dry-run validation", func() {
		app := applyAllButApplication(constants.AMApplicationFullFile).AMApplication

		warnings, err := admissionCtrl.ValidateCreate(ctx, app)
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(BeEmpty())
	})

	It("should fail dry-run validation when AM rejects it", func() {
		By("restarting the mock AM with dryRunReject=true")
		mockServer.Stop()
		rejectServer := am.Start(true)
		defer func() {
			rejectServer.Stop()
			mockServer = am.Start(false)
		}()

		app := applyAllButApplication(constants.AMApplicationMinimalFile).AMApplication

		_, err := admissionCtrl.ValidateCreate(ctx, app)
		Expect(err).To(HaveOccurred())
	})

	It("should admit with a warning, without dry-run, when the domain does not exist", func() {
		fixtures := fixture.Builder().
			WithAMContext(constants.AMContextFile).
			WithAMApplication(constants.AMApplicationKeysFile).
			Build()
		app := fixtures.AMApplication
		fixtures.AMApplication = nil
		fixtures.Apply()

		warnings, err := admissionCtrl.ValidateCreate(ctx, app)
		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ContainElement(ContainSubstring("not found or not yet created in AM")))
	})

	It("should reject a domainRef in another namespace", func() {
		app := withDomain(constants.AMApplicationKeysFile).AMApplication
		app.Spec.DomainRef.Namespace = "other"

		_, err := admissionCtrl.ValidateCreate(ctx, app)
		Expect(err).To(MatchError(ContainSubstring("domainRef.namespace")))
	})

	It("should reject a name that does not make a valid AM key", func() {
		app := withDomain(constants.AMApplicationKeysFile).AMApplication
		app.Name = "application.with.dots"

		_, err := admissionCtrl.ValidateCreate(ctx, app)
		Expect(err).To(HaveOccurred())
	})

	It("should reject a clientId other than the one AM generated", func() {
		app := applyAllButApplication(constants.AMApplicationMinimalFile).AMApplication
		app.Status.ClientID = "generated-by-am"
		changed := app.DeepCopy()
		changed.Spec.ClientId = new("chosen-later")

		_, err := admissionCtrl.ValidateUpdate(ctx, app, changed)
		Expect(err).To(MatchError(ContainSubstring(
			"clientId [chosen-later] differs from the one AM holds, [generated-by-am]")))
	})

	It("should admit setting the clientId AM generated", func() {
		app := applyAllButApplication(constants.AMApplicationMinimalFile).AMApplication
		app.Status.ClientID = "generated-by-am"
		same := app.DeepCopy()
		same.Spec.ClientId = new("generated-by-am")

		_, err := admissionCtrl.ValidateUpdate(ctx, app, same)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("certificateRef and identityRef", func() {
		It("should admit with a warning a certificateRef to a certificate that does not exist yet", func() {
			app := applyAllButApplication(constants.AMApplicationMinimalFile).AMApplication
			app.Spec.CertificateRef.Name = "missing"

			warnings, err := admissionCtrl.ValidateCreate(ctx, app)
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(ContainElement(
				"AMCertificate [missing] not found; [" + app.Name + "] will be created in AM once it exists"))
		})

		It("should admit with a warning an identityRef to an identity provider that does not exist yet", func() {
			app := applyAllButApplication(constants.AMApplicationFullFile).AMApplication
			app.Spec.IdentityProviders[0].IdentityRef.Name = "missing"

			warnings, err := admissionCtrl.ValidateCreate(ctx, app)
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(ContainElement(
				"AMIdentityProvider [missing] not found; [" + app.Name + "] will be created in AM once it exists"))
		})

		It("should reject a certificateRef in another namespace", func() {
			app := applyAllButApplication(constants.AMApplicationMinimalFile).AMApplication
			app.Spec.CertificateRef.Namespace = "other"

			_, err := admissionCtrl.ValidateCreate(ctx, app)
			Expect(err).To(MatchError(ContainSubstring("cross-namespace is not supported")))
		})

		It("should reject an identityRef in another namespace", func() {
			app := applyAllButApplication(constants.AMApplicationFullFile).AMApplication
			app.Spec.IdentityProviders[0].IdentityRef.Namespace = "other"

			_, err := admissionCtrl.ValidateCreate(ctx, app)
			Expect(err).To(MatchError(ContainSubstring("cross-namespace is not supported")))
		})

		It("should reject a certificateRef to a certificate of another domain", func() {
			fixtures := applyAllButApplication(constants.AMApplicationMinimalFile)
			app := fixtures.AMApplication
			otherCert := inOtherDomain(ctx, fixtures.AMCertificate)
			app.Spec.CertificateRef.Name = otherCert.Name

			_, err := admissionCtrl.ValidateCreate(ctx, app)
			Expect(err).To(MatchError(ContainSubstring(
				"AMCertificate [" + otherCert.Name + "] belongs to domain [other-domain], it should belong to [" +
					fixtures.AMSecurityDomain.Name + "]")))
		})

		It("should reject an identityRef to an identity provider of another domain", func() {
			fixtures := applyAllButApplication(constants.AMApplicationFullFile)
			app := fixtures.AMApplication
			otherIdp := &v1alpha1.AMIdentityProvider{
				ObjectMeta: metaV1.ObjectMeta{Name: fixtures.AMIdentityProvider.Name + "-other", Namespace: constants.Namespace},
				Spec:       *fixtures.AMIdentityProvider.Spec.DeepCopy(),
			}
			otherIdp.Spec.DomainRef = refs.NamespacedName{Name: "other-domain"}
			Expect(manager.Client().Create(ctx, otherIdp)).To(Succeed())
			app.Spec.IdentityProviders[0].IdentityRef.Name = otherIdp.Name

			_, err := admissionCtrl.ValidateCreate(ctx, app)
			Expect(err).To(MatchError(ContainSubstring(
				"AMIdentityProvider [" + otherIdp.Name + "] belongs to domain [other-domain]")))
		})

		It("should not check them while the application is being deleted", func() {
			fixtures := applyAllButApplication(constants.AMApplicationMinimalFile)
			app := fixtures.AMApplication
			otherCert := inOtherDomain(ctx, fixtures.AMCertificate)
			deleted := app.DeepCopy()
			deleted.Spec.CertificateRef.Name = otherCert.Name
			deleted.DeletionTimestamp = new(metaV1.Now())

			warnings, err := admissionCtrl.ValidateUpdate(ctx, app, deleted)
			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(BeEmpty())
		})
	})
})

// inOtherDomain creates a copy of cert under other-domain, a domain that does not need to exist.
func inOtherDomain(ctx context.Context, cert *v1alpha1.AMCertificate) *v1alpha1.AMCertificate {
	GinkgoHelper()
	other := &v1alpha1.AMCertificate{
		ObjectMeta: metaV1.ObjectMeta{Name: cert.Name + "-other", Namespace: constants.Namespace},
		Spec:       *cert.Spec.DeepCopy(),
	}
	other.Spec.DomainRef = refs.NamespacedName{Name: "other-domain"}
	Expect(manager.Client().Create(ctx, other)).To(Succeed())
	return other
}
