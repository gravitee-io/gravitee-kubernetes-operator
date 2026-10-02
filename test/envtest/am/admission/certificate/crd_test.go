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
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/constants"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/manager"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test/internal/integration/random"
)

// The CRD rules (schema and CEL) are evaluated by the API server, not by the webhook handler.
var _ = Describe("CRD validation", func() {
	ctx := context.Background()

	create := func(obj client.Object) error {
		err := manager.Client().Create(ctx, obj)
		if err == nil {
			DeferCleanup(func() { Expect(client.IgnoreNotFound(manager.Client().Delete(ctx, obj))).To(Succeed()) })
		}
		return err
	}
	newCertificate := func() *v1alpha1.AMCertificate {
		cert := &v1alpha1.AMCertificate{
			ObjectMeta: metaV1.ObjectMeta{Name: random.GetName(), Namespace: constants.Namespace},
		}
		cert.Spec.DomainRef = refs.NamespacedName{Name: "some-domain"}
		return cert
	}

	It("should require name, type and configuration unless system is true", func() {
		Expect(create(newCertificate())).
			To(MatchError(ContainSubstring("name, type and configuration are required unless system is true")))

		system := newCertificate()
		system.Spec.System = new(true)
		Expect(create(system)).To(Succeed())
	})

	It("should require domainRef", func() {
		cert := &unstructured.Unstructured{Object: map[string]any{
			"apiVersion": "gravitee.io/v1alpha1",
			"kind":       "AMCertificate",
			"metadata":   map[string]any{"name": random.GetName(), "namespace": constants.Namespace},
			"spec":       map[string]any{"system": true},
		}}
		Expect(create(cert)).To(MatchError(ContainSubstring("domainRef")))
	})

	It("should reject an empty domainRef name", func() {
		cert := newCertificate()
		cert.Spec.System = new(true)
		cert.Spec.DomainRef.Name = ""
		Expect(create(cert)).To(MatchError(ContainSubstring("domainRef.name must not be empty")))
	})

	Describe("immutable fields", func() {
		var cert *v1alpha1.AMCertificate

		BeforeEach(func() {
			cert = withDomain().Apply().AMCertificate
		})

		update := func(mutate func(*v1alpha1.AMCertificate)) error {
			latest := cert.DeepCopy()
			Expect(manager.GetLatest(ctx, latest)).To(Succeed())
			mutate(latest)
			return manager.Client().Update(ctx, latest)
		}

		It("should reject a type change", func() {
			Expect(update(func(c *v1alpha1.AMCertificate) { c.Spec.Type = new("other-type") })).
				To(MatchError(ContainSubstring("type is immutable")))
			Expect(update(func(c *v1alpha1.AMCertificate) { c.Spec.Type = nil })).
				To(MatchError(ContainSubstring("type is immutable")))
		})

		It("should reject a system change", func() {
			Expect(update(func(c *v1alpha1.AMCertificate) { c.Spec.System = new(true) })).
				To(MatchError(ContainSubstring("system is immutable")))
		})

		It("should reject a domainRef change", func() {
			Expect(update(func(c *v1alpha1.AMCertificate) { c.Spec.DomainRef.Name = "other-domain" })).
				To(MatchError(ContainSubstring("domainRef is immutable")))
		})

		It("should accept a name change", func() {
			Expect(update(func(c *v1alpha1.AMCertificate) { c.Spec.Name = new("Renamed") })).To(Succeed())
		})
	})
})
