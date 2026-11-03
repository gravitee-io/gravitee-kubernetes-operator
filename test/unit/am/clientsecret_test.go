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

package am_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/am/clientsecret"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

const clientSecretValue = "jFXuh6wb1_dhEu5h8Yi4Dp0fC3ZboI0I-SDoaGE22eI"

// newClientSecret returns a client secret of the mobile application, read from the client-secret entry of the
// mobile-secret Secret.
func newClientSecret() *v1alpha1.AMClientSecret {
	cs := &v1alpha1.AMClientSecret{ObjectMeta: metav1.ObjectMeta{Name: "mobile-2026", Namespace: "ns"}}
	cs.Spec.ApplicationRef = refs.NamespacedName{Name: "mobile"}
	cs.Spec.SecretRef = am.SecretKeyRef{Name: "mobile-secret", Key: "client-secret"}
	return cs
}

// readyApplication is the mobile application of the domain domain, accepted by AM.
func readyApplication() *v1alpha1.AMApplication {
	app := &v1alpha1.AMApplication{ObjectMeta: metav1.ObjectMeta{Name: "mobile", Namespace: "ns"}}
	app.Spec.DomainRef = refs.NamespacedName{Name: "domain"}
	app.Status.Key = "ns-mobile"
	return app
}

// mobileDomain is the AMSecurityDomain of the mobile application.
func mobileDomain() *v1alpha1.AMSecurityDomain {
	return &v1alpha1.AMSecurityDomain{ObjectMeta: metav1.ObjectMeta{Name: "domain", Namespace: "ns"}}
}

func valueSecret(data map[string][]byte) *corev1.Secret {
	return &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "mobile-secret", Namespace: "ns"}, Data: data}
}

// registerCluster registers a fake cluster holding objects.
func registerCluster(objects ...client.Object) {
	scheme := runtime.NewScheme()
	Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
	Expect(corev1.AddToScheme(scheme)).To(Succeed())
	k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build())
}

var _ = Describe("AMClientSecret DTO mapping", func() {
	It("keys the secret and its application by their HRIDs, with the domain key and value resolved", func() {
		cs := newClientSecret()
		cs.Spec.ApplicationRef.Namespace = "ignored"
		cs.Resolved = v1alpha1.ClientSecretResolved{DomainKey: "ns-domain", Value: clientSecretValue}

		dto, err := clientsecret.ToClientSecretDTO(cs)
		Expect(err).ToNot(HaveOccurred())

		Expect(dto.Key).To(Equal("ns-mobile-2026"))
		Expect(dto.ApplicationKey).To(Equal("ns-mobile"))
		Expect(dto.DomainKey).To(Equal("ns-domain"))
		Expect(dto.Value).To(Equal(new(clientSecretValue)))
		Expect(dto.Accepted).To(BeFalse())
	})

	It("leaves the value out when it is not resolved, and marks an accepted secret", func() {
		cs := newClientSecret()
		cs.Status.Key = "ns-mobile-2026"

		dto, err := clientsecret.ToClientSecretDTO(cs)
		Expect(err).ToNot(HaveOccurred())

		Expect(dto.Value).To(BeNil())
		Expect(dto.Accepted).To(BeTrue())
	})
})

var _ = Describe("AMClientSecret reference resolution", func() {
	ctx := context.Background()

	It("resolves the domain key of the application and the value of the Secret entry", func() {
		registerCluster(readyApplication(), valueSecret(map[string][]byte{"client-secret": []byte(clientSecretValue)}))
		cs := newClientSecret()

		Expect(clientsecret.ResolveRefs(ctx, cs, "ns")).To(Succeed())

		Expect(cs.Resolved.DomainKey).To(Equal("ns-domain"))
		Expect(cs.Resolved.Value).To(Equal(clientSecretValue))
	})

	It("returns a NotFound when the application is missing", func() {
		registerCluster()

		err := clientsecret.ResolveRefs(ctx, newClientSecret(), "ns")

		Expect(apierrors.IsNotFound(err)).To(BeTrue())
		Expect(err).To(MatchError(ContainSubstring("AMApplication [mobile] not found")))
	})

	It("waits for an application AM has not accepted yet", func() {
		app := readyApplication()
		app.Status.Key = ""
		registerCluster(app)

		err := clientsecret.ResolveRefs(ctx, newClientSecret(), "ns")

		Expect(errors.Is(err, clientsecret.ErrApplicationNotReady)).To(BeTrue())
	})

	It("names the Secret when it is missing", func() {
		registerCluster(readyApplication())

		err := clientsecret.ResolveRefs(ctx, newClientSecret(), "ns")

		Expect(apierrors.IsNotFound(err)).To(BeTrue())
		Expect(err).To(MatchError(ContainSubstring("Secret [mobile-secret] not found")))
	})

	DescribeTable("names the Secret and the key, never a value, when the entry is missing or empty",
		func(data map[string][]byte) {
			registerCluster(readyApplication(), valueSecret(data))

			err := clientsecret.ResolveRefs(ctx, newClientSecret(), "ns")

			Expect(errors.Is(err, clientsecret.ErrSecretEntryMissing)).To(BeTrue())
			Expect(err).To(MatchError(ContainSubstring("key [client-secret] of Secret [mobile-secret]")))
			Expect(err.Error()).ToNot(ContainSubstring(clientSecretValue))
		},
		Entry("missing", map[string][]byte{"other": []byte(clientSecretValue)}),
		Entry("empty", map[string][]byte{"client-secret": {}, "other": []byte(clientSecretValue)}),
	)

	It("names the application when its owner is missing", func() {
		registerCluster()

		_, _, err := clientsecret.GetOwner(ctx, newClientSecret())

		Expect(apierrors.IsNotFound(err)).To(BeTrue())
		Expect(err).To(MatchError(ContainSubstring("AMApplication [mobile] not found")))
	})

	It("does not read the Secret, nor wait for the application, on delete", func() {
		app := readyApplication()
		app.Status.Key = ""
		registerCluster(app, mobileDomain())
		cs := newClientSecret()
		cs.DeletionTimestamp = new(metav1.Now())

		Expect(clientsecret.ResolveRefs(ctx, cs, "ns")).To(Succeed())

		Expect(cs.Resolved.DomainKey).To(Equal("ns-domain"))
		Expect(cs.Resolved.Value).To(BeEmpty())
	})

	It("returns a NotFound on delete when the domain is gone: AM deleted the secret with it", func() {
		registerCluster(readyApplication())
		cs := newClientSecret()
		cs.DeletionTimestamp = new(metav1.Now())

		Expect(apierrors.IsNotFound(clientsecret.ResolveRefs(ctx, cs, "ns"))).To(BeTrue())
	})

	DescribeTable("admits a missing dependency in admission, leaving the value empty",
		func(objects ...client.Object) {
			registerCluster(objects...)
			cs := newClientSecret()

			Expect(clientsecret.ResolveAdmissionRefs(ctx, cs, "ns")).To(Succeed())

			Expect(cs.Resolved.Value).To(BeEmpty())
		},
		Entry("no application"),
		Entry("application not accepted", &v1alpha1.AMApplication{
			ObjectMeta: metav1.ObjectMeta{Name: "mobile", Namespace: "ns"},
		}),
		Entry("no Secret", readyApplication()),
		Entry("no entry", readyApplication(), valueSecret(map[string][]byte{})),
		Entry("empty entry", readyApplication(), valueSecret(map[string][]byte{"client-secret": {}})),
	)
})

var _ = Describe("AMClientSecret admission checks", func() {
	ctx := context.Background()

	It("rejects a key longer than the 64 characters AM allows for a client secret", func() {
		registerCluster(readyApplication(), valueSecret(map[string][]byte{"client-secret": []byte(clientSecretValue)}))
		cs := newClientSecret()
		cs.Name = strings.Repeat("a", 62)

		_, err := clientsecret.PreCheck(ctx, cs).Map()

		Expect(err).To(MatchError(ContainSubstring("longer than the 64 characters")))
	})

	It("accepts a key of exactly 64 characters", func() {
		registerCluster(readyApplication(), valueSecret(map[string][]byte{"client-secret": []byte(clientSecretValue)}))
		cs := newClientSecret()
		cs.Name = strings.Repeat("a", 61)

		warnings, err := clientsecret.PreCheck(ctx, cs).Map()

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(BeEmpty())
	})

	It("rejects a name that does not make a valid AM key", func() {
		registerCluster()
		cs := newClientSecret()
		cs.Name = "secret.with.dots"

		_, err := clientsecret.PreCheck(ctx, cs).Map()

		Expect(err).To(MatchError(ContainSubstring("is invalid: AM requires lowercase alphanumerics and hyphens")))
	})

	It("rejects an applicationRef in another namespace", func() {
		registerCluster()
		cs := newClientSecret()
		cs.Spec.ApplicationRef.Namespace = "other"

		_, err := clientsecret.PreCheck(ctx, cs).Map()

		Expect(err).To(MatchError(ContainSubstring("cross-namespace is not supported")))
	})

	DescribeTable("warns about a missing dependency, never with the value",
		func(want string, objects ...client.Object) {
			registerCluster(objects...)

			warnings, err := clientsecret.PreCheck(ctx, newClientSecret()).Map()

			Expect(err).ToNot(HaveOccurred())
			Expect(warnings).To(ConsistOf(And(
				ContainSubstring(want),
				ContainSubstring("will be created in AM once it exists"),
				Not(ContainSubstring(clientSecretValue)),
			)))
		},
		Entry("no application", "AMApplication [mobile] not found"),
		Entry("no Secret", "Secret [mobile-secret] not found", readyApplication()),
		Entry("no entry", "key [client-secret] of Secret [mobile-secret]",
			readyApplication(), valueSecret(map[string][]byte{"other": []byte(clientSecretValue)})),
		Entry("empty entry", "key [client-secret] of Secret [mobile-secret]",
			readyApplication(), valueSecret(map[string][]byte{"client-secret": {}})),
	)
})

var _ = Describe("AMClientSecret status", func() {
	ctx := context.Background()

	It("copies the keys and the expiration AM set", func() {
		expiresAt := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
		cs := newClientSecret()
		resp := clientsecret.Response{ApplicationKey: "ns-mobile", ExpiresAt: &expiresAt}
		resp.Key, resp.DomainKey = "ns-mobile-2026", "ns-domain"

		Expect(clientsecret.UpdateStatus(ctx, cs, resp)).To(Succeed())

		Expect(cs.Status.Key).To(Equal("ns-mobile-2026"))
		Expect(cs.Status.DomainKey).To(Equal("ns-domain"))
		Expect(cs.Status.ApplicationKey).To(Equal("ns-mobile"))
		Expect(cs.Status.ExpiresAt.Time).To(BeTemporally("==", expiresAt))
	})

	It("leaves expiresAt empty when the secret does not expire", func() {
		cs := newClientSecret()
		cs.Status.ExpiresAt = new(metav1.Now())

		Expect(clientsecret.UpdateStatus(ctx, cs, clientsecret.Response{})).To(Succeed())

		Expect(cs.Status.ExpiresAt).To(BeNil())
	})
})

var _ = Describe("AMClientSecret AM calls", func() {
	ctx := context.Background()

	DescribeTable("keeps AM's 400 message, with the rotation hint once the secret is accepted",
		func(accepted, wantHint bool) {
			s := newAMServer(http.StatusBadRequest, "application/json", `{"message":"client secret is immutable"}`)
			dto := clientsecret.ClientSecret{DomainKey: "ns-domain", ApplicationKey: "ns-mobile", Accepted: accepted}
			dto.Key, dto.Value = "ns-mobile-2026", new(clientSecretValue)

			_, err := clientsecret.Upsert(ctx, clientFor(s), dto)

			Expect(err).To(MatchError(ContainSubstring("client secret is immutable")))
			if wantHint {
				Expect(err).To(MatchError(ContainSubstring("rotate it with a new AMClientSecret")))
			} else {
				Expect(err).ToNot(MatchError(ContainSubstring("rotate it")))
			}
		},
		Entry("never accepted", false, false),
		Entry("accepted", true, true),
	)

	It("does not fail the delete of a secret whose application is gone", func() {
		s := newAMServer(http.StatusNotFound, "application/json", `{"message":"application not found"}`)
		dto := clientsecret.ClientSecret{DomainKey: "ns-domain", ApplicationKey: "ns-mobile"}
		dto.Key = "ns-mobile-2026"

		Expect(clientsecret.Delete(ctx, clientFor(s), dto)).To(Succeed())
	})
})
