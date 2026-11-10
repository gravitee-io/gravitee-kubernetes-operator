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

package search_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzschema"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

var _ = Describe("FindAuthzPdpDependents", func() {
	const ns = "orders"

	pdpRef := refs.NamespacedName{Name: "orders-eu", Namespace: ns}
	devCtx := refs.NamespacedName{Name: "dev-ctx"}
	pdp := &v1alpha1.AuthzPdp{
		ObjectMeta: metav1.ObjectMeta{Name: "orders-eu", Namespace: ns},
		Spec:       v1alpha1.AuthzPdpSpec{Context: &devCtx},
	}
	schemaThrough := func(name string, target *authzpdp.Target, contextRef refs.NamespacedName) *v1alpha1.AuthzSchema {
		return &v1alpha1.AuthzSchema{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
			Spec: v1alpha1.AuthzSchemaSpec{
				Context: &contextRef,
				Type:    authzschema.Type{Name: name, Target: target, Schema: "entity User;"},
			},
		}
	}
	targeting := func(name string, ref refs.NamespacedName) *v1alpha1.AuthzSchema {
		return schemaThrough(name, &authzpdp.Target{PdpRef: &ref}, devCtx)
	}
	entityTargeting := func(name string, targets ...authzpdp.Target) *v1alpha1.AuthzEntity {
		return &v1alpha1.AuthzEntity{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
			Spec: v1alpha1.AuthzEntitySpec{
				Context: &devCtx,
				Type:    authzentity.Type{EntityID: name, Kind: "PRINCIPAL", EntityType: "Group", Targets: targets},
			},
		}
	}
	// cluster indexes schemas and entities with the functions InitCache registers, as the
	// manager's cache does.
	cluster := func(objects ...client.Object) {
		c := &recordingCache{indexers: map[string]client.IndexerFunc{}}
		Expect(search.InitCache(context.Background(), c)).To(Succeed())
		scheme := runtime.NewScheme()
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		schemaField, entityField := search.AuthzSchemaAuthzPdpField.String(), search.AuthzEntityAuthzPdpField.String()
		k8s.RegisterClient(fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(objects...).
			WithIndex(&v1alpha1.AuthzSchema{}, schemaField, c.indexers[schemaField]).
			WithIndex(&v1alpha1.AuthzEntity{}, entityField, c.indexers[entityField]).
			Build())
	}

	It("registers AuthzSchema, then AuthzEntity, which needs a target", func() {
		kinds := search.AuthzPdpDependents()

		Expect(kinds).To(HaveLen(2))
		Expect(kinds[0].Kind).To(Equal("AuthzSchema"))
		Expect(kinds[0].Field).To(Equal(search.AuthzSchemaAuthzPdpField))
		Expect(kinds[0].NewList()).To(BeAssignableToTypeOf(&v1alpha1.AuthzSchemaList{}))
		Expect(kinds[0].Targets).To(BeNil())
		Expect(kinds[1].Kind).To(Equal("AuthzEntity"))
		Expect(kinds[1].Field).To(Equal(search.AuthzEntityAuthzPdpField))
		Expect(kinds[1].NewList()).To(BeAssignableToTypeOf(&v1alpha1.AuthzEntityList{}))
		Expect(kinds[1].Targets).ToNot(BeNil())
	})

	It("finds an entity from either of the PDPs it targets, and leaves out a scope-only entity", func() {
		usRef := refs.NamespacedName{Name: "orders-us", Namespace: ns}
		usPdp := &v1alpha1.AuthzPdp{
			ObjectMeta: metav1.ObjectMeta{Name: "orders-us", Namespace: ns},
			Spec:       v1alpha1.AuthzPdpSpec{Context: &devCtx},
		}
		cluster(
			entityTargeting("agents", authzpdp.Target{PdpRef: &refs.NamespacedName{Name: "orders-eu"}},
				authzpdp.Target{Scope: "default"}, authzpdp.Target{PdpRef: &usRef}),
			entityTargeting("scoped", authzpdp.Target{Scope: "orders@eu"}),
		)

		for _, target := range []*v1alpha1.AuthzPdp{pdp, usPdp} {
			found, err := search.FindAuthzPdpDependents(context.Background(), target, search.AuthzPdpDependents())

			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(Equal([]search.Dependents{
				{Kind: "AuthzEntity", Names: []string{"orders/agents"}},
			}))
		}
	})

	It("names apart the entities whose every target is the PDP, its namespace defaulted", func() {
		cluster(
			entityTargeting("alice", authzpdp.Target{PdpRef: &refs.NamespacedName{Name: "orders-eu"}}),
			entityTargeting("bob", authzpdp.Target{PdpRef: &pdpRef}, authzpdp.Target{PdpRef: &pdpRef}),
			entityTargeting("agents", authzpdp.Target{PdpRef: &pdpRef}, authzpdp.Target{Scope: "default"}),
		)

		found, err := search.FindAuthzPdpDependents(context.Background(), pdp, search.AuthzPdpDependents())

		Expect(err).ToNot(HaveOccurred())
		Expect(found).To(Equal([]search.Dependents{{
			Kind:  "AuthzEntity",
			Names: []string{"orders/agents", "orders/alice", "orders/bob"},
			Sole:  []string{"orders/alice", "orders/bob"},
		}}))
	})

	It("names the schemas whose pdpRef names the PDP, sorted, its namespace defaulted", func() {
		cluster(
			targeting("orders-schema", refs.NamespacedName{Name: "orders-eu"}),
			targeting("billing-schema", refs.NamespacedName{Name: "billing"}),
			targeting("audit-schema", pdpRef),
		)

		found, err := search.FindAuthzPdpDependents(context.Background(), pdp, search.AuthzPdpDependents())

		Expect(err).ToNot(HaveOccurred())
		Expect(found).To(Equal([]search.Dependents{{
			Kind:  "AuthzSchema",
			Names: []string{"orders/audit-schema", "orders/orders-schema"},
		}}))
	})

	It("leaves out a scope target, even one naming the PDP's routing scope, and a schema without target", func() {
		cluster(
			targeting("orders-schema", pdpRef),
			schemaThrough("scoped-schema", &authzpdp.Target{Scope: "orders@eu"}, devCtx),
			schemaThrough("untargeted-schema", nil, devCtx),
		)

		found, err := search.FindAuthzPdpDependents(context.Background(), pdp, search.AuthzPdpDependents())

		Expect(err).ToNot(HaveOccurred())
		Expect(found).To(Equal([]search.Dependents{{Kind: "AuthzSchema", Names: []string{"orders/orders-schema"}}}))
	})

	It("leaves out a schema synced through another context, whose pdpRef does not resolve", func() {
		cluster(
			targeting("orders-schema", pdpRef),
			schemaThrough("prod-schema", &authzpdp.Target{PdpRef: &pdpRef}, refs.NamespacedName{Name: "prod-ctx"}),
			schemaThrough("same-ctx-schema", &authzpdp.Target{PdpRef: &pdpRef},
				refs.NamespacedName{Name: "dev-ctx", Namespace: ns}),
		)

		found, err := search.FindAuthzPdpDependents(context.Background(), pdp, search.AuthzPdpDependents())

		Expect(err).ToNot(HaveOccurred())
		Expect(found).To(Equal([]search.Dependents{{
			Kind:  "AuthzSchema",
			Names: []string{"orders/orders-schema", "orders/same-ctx-schema"},
		}}))
	})

	It("leaves out a kind with no resource targeting the PDP", func() {
		cluster(targeting("billing-schema", refs.NamespacedName{Name: "billing"}))

		found, err := search.FindAuthzPdpDependents(context.Background(), pdp, search.AuthzPdpDependents())

		Expect(err).ToNot(HaveOccurred())
		Expect(found).To(BeEmpty())
	})

	It("refuses a kind whose resources carry no contextRef", func() {
		scheme := runtime.NewScheme()
		Expect(corev1.AddToScheme(scheme)).To(Succeed())
		field := search.AuthzSchemaAuthzPdpField.String()
		k8s.RegisterClient(fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "orders-schema", Namespace: ns}}).
			WithIndex(&corev1.ConfigMap{}, field, func(client.Object) []string { return []string{pdpRef.String()} }).
			Build())
		configMaps := search.AuthzPdpDependent{
			Kind:    "ConfigMap",
			Field:   search.AuthzSchemaAuthzPdpField,
			NewList: func() client.ObjectList { return &corev1.ConfigMapList{} },
		}

		_, err := search.FindAuthzPdpDependents(context.Background(), pdp,
			[]search.AuthzPdpDependent{configMaps})

		Expect(err).To(MatchError(ContainSubstring("unable to read the contextRef of ConfigMap")))
	})
})
