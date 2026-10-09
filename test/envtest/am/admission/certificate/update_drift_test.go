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

package amcertificate

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-automation-sdk/common/pkg/unstructured"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/amcertificate"
	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/certificate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
)

var _ = Describe("Validate drift", func() {
	ctx := context.Background()
	admissionCtrl := amcertificate.NewAdmissionCtrl()

	// inAM applies the domain only, edits the certificate spec with local, and stores the certificate in the
	// AM mock as AM returns it: the spec, with the fields remote overrides.
	inAM := func(local func(*v1alpha1.AMCertificateSpec), remote func(*internal.Certificate)) *v1alpha1.AMCertificate {
		fixtures := withDomain()
		cert := fixtures.AMCertificate
		local(&cert.Spec)
		fixtures.AMCertificate = nil
		fixtures.Apply()

		dto, err := internal.ToCertificateDTO(cert)
		Expect(err).ToNot(HaveOccurred())
		remote(&dto)
		resp, err := am.NewSDKClient().UpsertCertificateWithResponse(ctx, dto.DomainKey, nil, dto.Certificate)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.JSON200).ToNot(BeNil(), string(resp.Body))
		// synced as the controller leaves it: drift only runs for a resource AM accepted
		cert.Status.Key = resp.JSON200.Key
		return cert
	}
	validateUpdate := func(cert *v1alpha1.AMCertificate) func() error {
		return func() error {
			updated := cert.DeepCopy()
			updated.Annotations = map[string]string{"updated": "true"}
			_, err := admissionCtrl.ValidateUpdate(ctx, cert, updated)
			return err
		}
	}

	// Left unset, name, type and configuration never drift: a field the CRD does not set is not compared.
	// Set, they are ignored by AM for a system certificate, which returns its own.
	It("should not drift on the name, type and configuration AM replaces for a system certificate", func() {
		cert := inAM(func(spec *v1alpha1.AMCertificateSpec) {
			spec.System = new(true)
		}, func(remote *internal.Certificate) {
			remote.Name = new("Default")
			remote.Type = new("javakeystore-am-certificate")
			remote.Configuration = unstructured.StringifiedFrom(map[string]any{"alias": "default", "storepass": "********"})
		})

		Consistently(validateUpdate(cert), constants.ConsistentTimeout, constants.Interval).Should(Succeed())
	})

	It("should detect drift on the name of a regular certificate", func() {
		cert := inAM(func(*v1alpha1.AMCertificateSpec) {}, func(remote *internal.Certificate) {
			remote.Name = new("Renamed in AM")
		})

		Eventually(validateUpdate(cert), constants.EventualTimeout, constants.Interval).
			Should(MatchError(And(ContainSubstring("drift detected"), ContainSubstring("Renamed in AM"))))
	})

	It("should not drift on the passwords and keystore file AM returns masked", func() {
		cert := inAM(func(*v1alpha1.AMCertificateSpec) {}, func(remote *internal.Certificate) {
			remote.Configuration.Object["content"] = "********"
			remote.Configuration.Object["storepass"] = "********"
			remote.Configuration.Object["keypass"] = "********"
		})

		Consistently(validateUpdate(cert), constants.ConsistentTimeout, constants.Interval).Should(Succeed())
	})
})
