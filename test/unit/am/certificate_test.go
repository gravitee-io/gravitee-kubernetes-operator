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
package am_test

import (
	"context"
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am"
	certificate "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/certificate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

const keystoreFile = `{"name":"keystore.p12","content":"a2V5c3RvcmU="}`

// certWithReadyDomain registers a fake cluster holding a domain AM has created, and returns a PKCS#12
// certificate under it.
func certWithReadyDomain() *v1alpha1.AMCertificate {
	d := &v1alpha1.AMSecurityDomain{ObjectMeta: metav1.ObjectMeta{Name: "domain", Namespace: "ns"}}
	d.Status.Key = "ns-domain"
	scheme := runtime.NewScheme()
	Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
	k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).WithObjects(d).Build())

	c := &v1alpha1.AMCertificate{ObjectMeta: metav1.ObjectMeta{Name: "cert", Namespace: "ns"}}
	c.Spec.DomainRef = refs.NamespacedName{Name: "domain"}
	c.Spec.Name = new("Token signing")
	c.Spec.Type = new("pkcs12-am-certificate")
	c.Spec.Configuration = utils.NewGenericStringMap().
		Put("content", keystoreFile).
		Put("alias", "signing").
		Put("storepass", "changeme")
	return c
}

var _ = Describe("AMCertificate DTO mapping", func() {
	It("keys the certificate and its domain by their HRIDs", func() {
		cert := certWithReadyDomain()
		cert.Spec.DomainRef.Namespace = "ignored"

		dto, err := certificate.ToCertificateDTO(cert)
		Expect(err).ToNot(HaveOccurred())

		Expect(dto.Key).To(Equal("ns-cert"))
		Expect(dto.DomainKey).To(Equal("ns-domain"))
	})

	It("sends configuration as a JSON string, the keystore file as given", func() {
		dto, err := certificate.ToCertificateDTO(certWithReadyDomain())
		Expect(err).ToNot(HaveOccurred())

		payload, err := json.Marshal(dto.Certificate)
		Expect(err).ToNot(HaveOccurred())
		var sent map[string]any
		Expect(json.Unmarshal(payload, &sent)).To(Succeed())
		sentConfiguration, ok := sent["configuration"].(string)
		Expect(ok).To(BeTrue(), "configuration sent as %T", sent["configuration"])
		var configuration map[string]any
		Expect(json.Unmarshal([]byte(sentConfiguration), &configuration)).To(Succeed())
		Expect(configuration).To(Equal(map[string]any{
			"content":   keystoreFile,
			"alias":     "signing",
			"storepass": "changeme",
		}))
		Expect(sent["name"]).To(Equal("Token signing"))
		Expect(sent["type"]).To(Equal("pkcs12-am-certificate"))
	})

	It("defaults system to false", func() {
		dto, err := certificate.ToCertificateDTO(certWithReadyDomain())
		Expect(err).ToNot(HaveOccurred())

		Expect(dto.System).ToNot(BeNil())
		Expect(*dto.System).To(BeFalse())
	})
})

var _ = Describe("AMCertificate drift mapping", func() {
	It("blanks name, type and configuration of a system certificate: AM supplies them", func() {
		cert := certWithReadyDomain()
		cert.Spec.System = new(true)

		dto, err := certificate.ToCertificateDTOForDrift(cert)
		Expect(err).ToNot(HaveOccurred())

		Expect(dto.Name).To(BeNil())
		Expect(dto.Type).To(BeNil())
		Expect(dto.Configuration).To(BeNil())
		Expect(dto.Key).To(Equal("ns-cert"))
	})

	It("keeps name, type and configuration of a regular certificate", func() {
		cert := certWithReadyDomain()

		dto, err := certificate.ToCertificateDTOForDrift(cert)
		Expect(err).ToNot(HaveOccurred())
		expected, err := certificate.ToCertificateDTO(cert)
		Expect(err).ToNot(HaveOccurred())

		Expect(dto).To(Equal(expected))
	})
})

var _ = Describe("AMCertificate PreCheck", func() {
	ctx := context.Background()

	It("admits a valid certificate without warnings", func() {
		errs := certificate.PreCheck(ctx, certWithReadyDomain())

		Expect(errs.IsSevere()).To(BeFalse())
		Expect(warnings(errs)).To(BeEmpty())
	})

	It("rejects a key AM would refuse", func() {
		cert := certWithReadyDomain()
		cert.Name = "cert.with.dots"

		Expect(certificate.PreCheck(ctx, cert).IsSevere()).To(BeTrue())
	})

	It("rejects a domainRef in another namespace", func() {
		cert := certWithReadyDomain()
		cert.Spec.DomainRef.Namespace = "other"

		Expect(certificate.PreCheck(ctx, cert).IsSevere()).To(BeTrue())
	})

	It("warns when the domain is missing", func() {
		cert := certWithReadyDomain()
		cert.Spec.DomainRef.Name = "missing"

		errs := certificate.PreCheck(ctx, cert)

		Expect(errs.IsSevere()).To(BeFalse())
		Expect(warnings(errs)).To(ContainElement(ContainSubstring("not found or not yet created in AM")))
	})

	It("warns that name, type and configuration are ignored for the system certificate", func() {
		cert := certWithReadyDomain()
		cert.Spec.System = new(true)

		errs := certificate.PreCheck(ctx, cert)

		Expect(errs.IsSevere()).To(BeFalse())
		Expect(warnings(errs)).To(ConsistOf(
			"'configuration', 'name', 'type' will be ignored when 'system' is 'true'."))
	})
})

var _ = Describe("AMCertificate UpdateStatus", func() {
	response := func(expiresAt *time.Time) certificate.Response {
		return certificate.Response{
			DomainSubResourceResponse: am.DomainSubResourceResponse{
				DomainKey: "ns-domain",
				BaseResponse: am.BaseResponse{
					Key:    "ns-cert",
					OrgEnv: am.OrgEnv{OrgID: "DEFAULT", EnvID: "DEFAULT"},
				},
			},
			ExpiresAt: expiresAt,
		}
	}

	It("copies the keys and expiresAt", func() {
		expiresAt := time.Date(2027, 10, 1, 0, 0, 0, 0, time.UTC)
		cert := &v1alpha1.AMCertificate{}

		Expect(certificate.UpdateStatus(context.Background(), cert, response(&expiresAt))).To(Succeed())

		Expect(cert.Status.Key).To(Equal("ns-cert"))
		Expect(cert.Status.DomainKey).To(Equal("ns-domain"))
		Expect(cert.Status.OrgID).To(Equal("DEFAULT"))
		Expect(cert.Status.EnvID).To(Equal("DEFAULT"))
		Expect(cert.Status.ExpiresAt).ToNot(BeNil())
		Expect(cert.Status.ExpiresAt.Time).To(BeTemporally("==", expiresAt))
	})

	It("keeps expiresAt nil when AM does not report it", func() {
		cert := &v1alpha1.AMCertificate{}

		Expect(certificate.UpdateStatus(context.Background(), cert, response(nil))).To(Succeed())

		Expect(cert.Status.ExpiresAt).To(BeNil())
	})
})
