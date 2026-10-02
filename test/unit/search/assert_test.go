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

package search_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	amdomain "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am/domain"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

var _ = Describe("AssertNoAMContextRef", func() {
	const ns = "ns"

	amContext := &v1alpha1.AMContext{ObjectMeta: metav1.ObjectMeta{Name: "am-ctx", Namespace: ns}}
	domain := func(name, contextName string) *v1alpha1.AMSecurityDomain {
		return &v1alpha1.AMSecurityDomain{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
			Spec:       v1alpha1.AMSecurityDomainSpec{Context: &refs.NamespacedName{Name: contextName}},
		}
	}
	// cluster registers a fake client indexing domains with the index function InitCache registers.
	cluster := func(objects ...client.Object) {
		c := &recordingCache{indexers: map[string]client.IndexerFunc{}}
		Expect(search.InitCache(context.Background(), c)).To(Succeed())
		scheme := runtime.NewScheme()
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		field := search.AMSecurityContextField.String()
		k8s.RegisterClient(fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(objects...).
			WithIndex(&v1alpha1.AMSecurityDomain{}, field, c.indexers[field]).
			Build())
	}

	It("fails with the number of domains that reference the context", func() {
		cluster(domain("d1", "am-ctx"), domain("d2", "other-ctx"), domain("d3", "am-ctx"))

		err := search.AssertNoAMContextRef(context.Background(), amContext)

		Expect(err).To(MatchError(ContainSubstring("[am-ctx] cannot be deleted because 2 AM security domains")))
	})

	It("counts a domain being deleted", func() {
		deleting := domain("d1", "am-ctx")
		deleting.Finalizers = []string{"finalizers.gravitee.io/amsecuritydomains"}
		deleting.DeletionTimestamp = &metav1.Time{Time: time.Now()}
		cluster(deleting)

		Expect(search.AssertNoAMContextRef(context.Background(), amContext)).
			To(MatchError(ContainSubstring("because 1 AM security domains")))
	})

	It("succeeds when no domain references the context", func() {
		cluster(domain("d1", "other-ctx"))

		Expect(search.AssertNoAMContextRef(context.Background(), amContext)).To(Succeed())
	})
})

var _ = Describe("AssertNoAMCertificateRef", func() {
	const ns = "ns"

	// key in AM: ns-cert
	cert := &v1alpha1.AMCertificate{ObjectMeta: metav1.ObjectMeta{Name: "cert", Namespace: ns}}
	domain := func(name, namespace, fallback string) *v1alpha1.AMSecurityDomain {
		d := &v1alpha1.AMSecurityDomain{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace}}
		d.Spec.CertificateSettings = &amdomain.CertificateSettings{FallbackCertificate: new(fallback)}
		return d
	}
	// cluster registers a fake client indexing domains with the index function InitCache registers.
	cluster := func(objects ...client.Object) {
		c := &recordingCache{indexers: map[string]client.IndexerFunc{}}
		Expect(search.InitCache(context.Background(), c)).To(Succeed())
		scheme := runtime.NewScheme()
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		field := search.AMSecurityDomainFallbackCertificateField.String()
		k8s.RegisterClient(fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(objects...).
			WithIndex(&v1alpha1.AMSecurityDomain{}, field, c.indexers[field]).
			Build())
	}

	It("fails with the number of domains that use the certificate as fallback", func() {
		cluster(domain("d1", ns, "ns-cert"), domain("d2", ns, "ns-other"), domain("d3", ns, "ns-cert"))

		err := search.AssertNoAMCertificateRef(context.Background(), cert)

		Expect(err).To(MatchError(ContainSubstring(
			"[cert] cannot be deleted because 2 AM security domains use it as fallback certificate")))
	})

	It("ignores a domain being deleted", func() {
		deleting := domain("d1", ns, "ns-cert")
		deleting.Finalizers = []string{"finalizers.gravitee.io/amsecuritydomains"}
		deleting.DeletionTimestamp = &metav1.Time{Time: time.Now()}
		cluster(deleting)

		Expect(search.AssertNoAMCertificateRef(context.Background(), cert)).To(Succeed())
	})

	It("ignores a domain of another namespace", func() {
		cluster(domain("d1", "other", "ns-cert"))

		Expect(search.AssertNoAMCertificateRef(context.Background(), cert)).To(Succeed())
	})
})
