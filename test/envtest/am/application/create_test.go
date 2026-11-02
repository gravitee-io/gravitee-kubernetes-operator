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
	"fmt"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/assert"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
)

var _ = Describe("Create", func() {
	ctx := context.Background()

	It("should create the application under its domain, owned by the domain, with the clientId AM generated", func() {
		fixtures := fixtures(constants.AMApplicationMinimalFile).Apply()
		app := fixtures.AMApplication

		Expect(app.Status.Key).To(Equal(refs.NewNamespacedNameFromObject(app).HRID()))
		Expect(app.Status.DomainKey).To(Equal(fixtures.AMSecurityDomain.Status.Key))

		By("expecting the clientId AM generated in the status")

		status, remote := remoteApplication(ctx, app)
		Expect(status).To(Equal(http.StatusOK))
		Expect(remote.ClientId).ToNot(BeNil())
		Expect(app.Status.ClientID).To(Equal(*remote.ClientId))

		By("expecting the domain to own the application")

		Expect(app.OwnerReferences).To(HaveLen(1))
		owner := app.OwnerReferences[0]
		Expect(owner.Kind).To(Equal("AMSecurityDomain"))
		Expect(owner.Name).To(Equal(fixtures.AMSecurityDomain.Name))
		Expect(owner.UID).To(Equal(fixtures.AMSecurityDomain.UID))
	})

	It("should wait for its domain when applied before it", func() {
		fixtures := fixtures(constants.AMApplicationKeysFile)
		domain, app := fixtures.AMSecurityDomain, fixtures.AMApplication
		fixtures.AMSecurityDomain, fixtures.AMCertificate, fixtures.AMIdentityProvider, fixtures.AMApplication =
			nil, nil, nil, nil
		fixtures.Apply()

		By("creating the application before its domain")

		Expect(manager.Client().Create(ctx, app)).To(Succeed())
		expectUnresolvedRefs(ctx, app, "AMSecurityDomain")

		By("creating the domain")

		Expect(manager.Client().Create(ctx, domain)).To(Succeed())
		expectAccepted(ctx, app)
	})

	It("should wait for its certificate and identity provider when applied before them", func() {
		fixtures := fixtures(constants.AMApplicationFullFile)
		cert, idp, app := fixtures.AMCertificate, fixtures.AMIdentityProvider, fixtures.AMApplication
		fixtures.AMCertificate, fixtures.AMIdentityProvider, fixtures.AMApplication = nil, nil, nil
		fixtures.Apply()

		By("creating the application before its certificate and identity provider")

		Expect(manager.Client().Create(ctx, app)).To(Succeed())
		expectUnresolvedRefs(ctx, app, "AMCertificate")
		status, _ := remoteApplication(ctx, app)
		Expect(status).To(Equal(http.StatusNotFound))

		By("creating the certificate: the application now waits for its identity provider")

		Expect(manager.Client().Create(ctx, cert)).To(Succeed())
		expectUnresolvedRefs(ctx, app, "AMIdentityProvider")

		By("creating the identity provider: the watches pick the application up")

		Expect(manager.Client().Create(ctx, idp)).To(Succeed())
		expectAccepted(ctx, app)

		status, remote := remoteApplication(ctx, app)
		Expect(status).To(Equal(http.StatusOK))
		Expect(remote.Certificate).To(Equal(refs.NewNamespacedNameFromObject(cert).HRID()))
		Expect(remote.IdentityProviders).To(HaveLen(1))
		Expect(remote.IdentityProviders[0].Identity).To(Equal(refs.NewNamespacedNameFromObject(idp).HRID()))
	})

	It("should be rejected by AM with a clientId another application of the domain uses (AC4)", func() {
		fixtures := fixtures(constants.AMApplicationClientIDFile).Apply()
		first := fixtures.AMApplication

		duplicate := &v1alpha1.AMApplication{
			ObjectMeta: metaV1.ObjectMeta{Name: first.Name + "-duplicate", Namespace: constants.Namespace},
			Spec:       *first.Spec.DeepCopy(),
		}
		Expect(manager.Client().Create(ctx, duplicate)).To(Succeed())

		expectRejected(ctx, duplicate, "clientId ["+*first.Spec.ClientId+"] is already used by application ["+
			first.Status.Key+"]")
		status, _ := remoteApplication(ctx, duplicate)
		Expect(status).To(Equal(http.StatusNotFound))
	})

	It("should not create an application referencing a certificate of another domain", func() {
		fixtures := fixtures(constants.AMApplicationMinimalFile)
		app := fixtures.AMApplication
		fixtures.AMApplication = nil
		fixtures.Apply()

		By("creating another domain with its own certificate")

		otherDomain := &v1alpha1.AMSecurityDomain{
			ObjectMeta: metaV1.ObjectMeta{Name: fixtures.AMSecurityDomain.Name + "-other", Namespace: constants.Namespace},
			Spec:       *fixtures.AMSecurityDomain.Spec.DeepCopy(),
		}
		otherDomain.Spec.Path += "-other"
		Expect(manager.Client().Create(ctx, otherDomain)).To(Succeed())
		otherCert := &v1alpha1.AMCertificate{
			ObjectMeta: metaV1.ObjectMeta{Name: fixtures.AMCertificate.Name + "-other", Namespace: constants.Namespace},
			Spec:       *fixtures.AMCertificate.Spec.DeepCopy(),
		}
		otherCert.Spec.DomainRef = refs.NamespacedName{Name: otherDomain.Name}
		Expect(manager.Client().Create(ctx, otherCert)).To(Succeed())

		By("creating the application with a certificateRef to that certificate")

		app.Spec.CertificateRef = &refs.NamespacedName{Name: otherCert.Name}
		Expect(manager.Client().Create(ctx, app)).To(Succeed())

		expectUnresolvedRefs(ctx, app, fmt.Sprintf("belongs to domain [%s], it should belong to [%s]",
			otherDomain.Name, fixtures.AMSecurityDomain.Name))
		status, _ := remoteApplication(ctx, app)
		Expect(status).To(Equal(http.StatusNotFound))
	})
})

// expectRejected waits until the application reports Accepted=False with a message containing part.
func expectRejected(ctx context.Context, app *v1alpha1.AMApplication, part string) {
	GinkgoHelper()
	Eventually(ctx, func() error {
		if err := manager.GetLatest(ctx, app); err != nil {
			return err
		}
		accepted := k8s.MapConditions(app.Status.Conditions)[k8s.ConditionAccepted]
		if accepted.Status != metaV1.ConditionFalse || !strings.Contains(accepted.Message, part) {
			return fmt.Errorf("Accepted %s: %s", accepted.Status, accepted.Message)
		}
		return nil
	}, constants.EventualTimeout, constants.Interval).Should(Succeed(), app.Name)
}

// expectUnresolvedRefs waits until the application reports ResolvedRefs=False with a message containing part.
func expectUnresolvedRefs(ctx context.Context, app *v1alpha1.AMApplication, part string) {
	GinkgoHelper()
	Eventually(ctx, func() error {
		if err := manager.GetLatest(ctx, app); err != nil {
			return err
		}
		refs := k8s.MapConditions(app.Status.Conditions)[k8s.ConditionResolvedRefs]
		if refs.Status != metaV1.ConditionFalse || !strings.Contains(refs.Message, part) {
			return fmt.Errorf("ResolvedRefs %s: %s", refs.Status, refs.Message)
		}
		return nil
	}, constants.EventualTimeout, constants.Interval).Should(Succeed(), app.Name)
}

// expectAccepted waits until the application is accepted.
func expectAccepted(ctx context.Context, app *v1alpha1.AMApplication) {
	GinkgoHelper()
	Eventually(ctx, func() error {
		if err := manager.GetLatest(ctx, app); err != nil {
			return err
		}
		return assert.AMApplicationAccepted(app)
	}, constants.EventualTimeout, constants.Interval).Should(Succeed(), app.Name)
}
