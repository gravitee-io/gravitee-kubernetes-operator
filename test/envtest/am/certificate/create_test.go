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
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	internal "github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/certificate"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Create", func() {
	ctx := context.Background()

	It("should create the certificate under its domain, owned by the domain", func() {
		fixtures := fixtures().Apply()
		cert := fixtures.AMCertificate

		By("expecting the certificate to be accepted")

		Expect(assert.AMCertificateAccepted(cert)).To(Succeed())
		dto, err := internal.ToCertificateDTO(cert)
		Expect(err).ToNot(HaveOccurred())
		Expect(cert.Status.Key).To(Equal(dto.Key))
		Expect(cert.Status.DomainKey).To(Equal(dto.DomainKey))

		By("expecting the certificate in AM, under the domain")

		status, remote := remoteCertificate(ctx, cert)
		Expect(status).To(Equal(http.StatusOK))
		Expect(remote.Name).To(Equal(cert.Spec.Name))
		Expect(remote.Configuration.Object).To(HaveKeyWithValue("alias", "signing"))

		By("expecting the domain to own the certificate")

		Expect(cert.OwnerReferences).To(HaveLen(1))
		owner := cert.OwnerReferences[0]
		Expect(owner.Kind).To(Equal("AMSecurityDomain"))
		Expect(owner.Name).To(Equal(fixtures.AMSecurityDomain.Name))
		Expect(owner.UID).To(Equal(fixtures.AMSecurityDomain.UID))
	})

	It("should wait for its domain when applied before it", func() {
		fixtures := fixtures()
		domain, cert := fixtures.AMSecurityDomain, fixtures.AMCertificate
		fixtures.AMSecurityDomain, fixtures.AMCertificate = nil, nil
		fixtures.Apply()

		By("creating the certificate before its domain")

		Expect(manager.Client().Create(ctx, cert)).To(Succeed())

		Eventually(ctx, func() error {
			if err := manager.GetLatest(ctx, cert); err != nil {
				return err
			}
			refs := k8s.MapConditions(cert.Status.Conditions)[k8s.ConditionResolvedRefs]
			return assert.Equals("ResolvedRefs", metaV1.ConditionFalse, refs.Status)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), cert.Name)

		By("creating the domain")

		Expect(manager.Client().Create(ctx, domain)).To(Succeed())

		By("expecting the certificate to be accepted once the domain is in AM")

		Eventually(ctx, func() error {
			if err := manager.GetLatest(ctx, cert); err != nil {
				return err
			}
			return assert.AMCertificateAccepted(cert)
		}, constants.EventualTimeout, constants.Interval).Should(Succeed(), cert.Name)

		status, _ := remoteCertificate(ctx, cert)
		Expect(status).To(Equal(http.StatusOK))
	})
})
