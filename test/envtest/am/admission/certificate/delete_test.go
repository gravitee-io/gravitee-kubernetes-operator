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
package amcertificate

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am/domain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amcertificate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/fixture"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

const holdFinalizer = "test.gravitee.io/hold"

var _ = Describe("Validate delete (D2a)", func() {
	ctx := context.Background()
	admissionCtrl := amcertificate.NewAdmissionCtrl()

	// withFallback applies a domain using the certificate as fallback, and the certificate.
	withFallback := func() *fixture.Objects {
		fixtures := withDomain()
		fallback := refs.NewNamespacedNameFromObject(fixtures.AMCertificate).HRID()
		fixtures.AMSecurityDomain.Spec.CertificateSettings = &domain.CertificateSettings{FallbackCertificate: &fallback}
		return fixtures.Apply()
	}
	validateDelete := func(cert *v1alpha1.AMCertificate) func() error {
		return func() error {
			_, err := admissionCtrl.ValidateDelete(ctx, cert)
			return err
		}
	}

	It("should reject deleting a certificate a domain uses as fallback, until the domain changes it", func() {
		fixtures := withFallback()
		domain, cert := fixtures.AMSecurityDomain, fixtures.AMCertificate

		Eventually(validateDelete(cert), constants.EventualTimeout, constants.Interval).
			Should(MatchError(ContainSubstring("cannot be deleted because 1 AM security domains use it as fallback")))

		By("changing the domain's fallback certificate")

		Expect(manager.GetLatest(ctx, domain)).To(Succeed())
		domain.Spec.CertificateSettings.FallbackCertificate = new("another-certificate")
		Expect(manager.Client().Update(ctx, domain)).To(Succeed())

		Eventually(validateDelete(cert), constants.EventualTimeout, constants.Interval).Should(Succeed())
	})

	It("should admit deleting a certificate whose fallback domain is being deleted", func() {
		fixtures := withFallback()
		domain, cert := fixtures.AMSecurityDomain, fixtures.AMCertificate
		Eventually(validateDelete(cert), constants.EventualTimeout, constants.Interval).Should(HaveOccurred())

		By("deleting the domain, held by a finalizer so it stays being deleted")

		Expect(manager.GetLatest(ctx, domain)).To(Succeed())
		controllerutil.AddFinalizer(domain, holdFinalizer)
		Expect(manager.Client().Update(ctx, domain)).To(Succeed())
		DeferCleanup(func() {
			Expect(manager.GetLatest(ctx, domain)).To(Succeed())
			controllerutil.RemoveFinalizer(domain, holdFinalizer)
			Expect(manager.Client().Update(ctx, domain)).To(Succeed())
		})
		Expect(manager.Client().Delete(ctx, domain.DeepCopy())).To(Succeed())

		Eventually(validateDelete(cert), constants.EventualTimeout, constants.Interval).Should(Succeed())
	})
})
