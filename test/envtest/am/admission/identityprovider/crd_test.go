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

package amidentityprovider

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
	newIdentityProvider := func() *v1alpha1.AMIdentityProvider {
		idp := &v1alpha1.AMIdentityProvider{
			ObjectMeta: metaV1.ObjectMeta{Name: random.GetName(), Namespace: constants.Namespace},
		}
		idp.Spec.DomainRef = refs.NamespacedName{Name: "some-domain"}
		return idp
	}

	It("should require name, type and configuration unless system is true", func() {
		Expect(create(newIdentityProvider())).
			To(MatchError(ContainSubstring("name, type and configuration are required unless system is true")))

		system := newIdentityProvider()
		system.Spec.System = new(true)
		Expect(create(system)).To(Succeed())
	})

	It("should require domainRef", func() {
		idp := &unstructured.Unstructured{Object: map[string]any{
			"apiVersion": "gravitee.io/v1alpha1",
			"kind":       "AMIdentityProvider",
			"metadata":   map[string]any{"name": random.GetName(), "namespace": constants.Namespace},
			"spec":       map[string]any{"system": true},
		}}
		Expect(create(idp)).To(MatchError(ContainSubstring("domainRef")))
	})

	It("should reject an empty domainRef name", func() {
		idp := newIdentityProvider()
		idp.Spec.System = new(true)
		idp.Spec.DomainRef.Name = ""
		Expect(create(idp)).To(MatchError(ContainSubstring("domainRef.name must not be empty")))
	})

	Describe("immutable fields", func() {
		var idp *v1alpha1.AMIdentityProvider

		BeforeEach(func() {
			idp = withDomain().Apply().AMIdentityProvider
		})

		update := func(mutate func(*v1alpha1.AMIdentityProvider)) error {
			latest := idp.DeepCopy()
			Expect(manager.GetLatest(ctx, latest)).To(Succeed())
			mutate(latest)
			return manager.Client().Update(ctx, latest)
		}

		It("should reject a type change", func() {
			Expect(update(func(p *v1alpha1.AMIdentityProvider) { p.Spec.Type = new("other-type") })).
				To(MatchError(ContainSubstring("type is immutable")))
			Expect(update(func(p *v1alpha1.AMIdentityProvider) { p.Spec.Type = nil })).
				To(MatchError(ContainSubstring("type is immutable")))
		})

		It("should reject a system change", func() {
			Expect(update(func(p *v1alpha1.AMIdentityProvider) { p.Spec.System = new(true) })).
				To(MatchError(ContainSubstring("system is immutable")))
		})

		It("should reject a domainRef change", func() {
			Expect(update(func(p *v1alpha1.AMIdentityProvider) { p.Spec.DomainRef.Name = "other-domain" })).
				To(MatchError(ContainSubstring("domainRef is immutable")))
		})

		It("should accept a name change", func() {
			Expect(update(func(p *v1alpha1.AMIdentityProvider) { p.Spec.Name = new("Renamed") })).To(Succeed())
		})
	})
})
