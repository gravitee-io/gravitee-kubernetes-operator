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

package authzpdp_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpolicy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzschema"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/search"
)

var _ = Describe("Deleting an AuthzPdp", func() {
	const ns = "travel-team"
	const refused = "The platform refuses the delete (409) while one of them targets no other gateway: " +
		"the AuthzPdp stays Terminating until it is retargeted"

	ctx := context.Background()
	ctrl := admission.AdmissionCtrl{}
	devCtx := refs.NamespacedName{Name: "dev-ctx"}
	target := authzpdp.Target{PdpRef: &refs.NamespacedName{Name: "orders-eu"}}

	pdp := &v1alpha1.AuthzPdp{
		ObjectMeta: metav1.ObjectMeta{Name: "orders-eu", Namespace: ns},
		Spec:       v1alpha1.AuthzPdpSpec{Context: &devCtx},
		Status:     v1alpha1.AuthzPdpStatus{Status: authzpdp.Status{RoutingScope: "orders@eu"}},
	}
	schema := &v1alpha1.AuthzSchema{
		ObjectMeta: metav1.ObjectMeta{Name: "travel-schema", Namespace: ns},
		Spec: v1alpha1.AuthzSchemaSpec{
			Context: &devCtx,
			Type:    authzschema.Type{Name: "Travel schema", Target: &target, Schema: "entity User;"},
		},
	}
	entity := &v1alpha1.AuthzEntity{
		ObjectMeta: metav1.ObjectMeta{Name: "agents-group", Namespace: ns},
		Spec: v1alpha1.AuthzEntitySpec{
			Context: &devCtx,
			Type: authzentity.Type{
				EntityID: "agents", Kind: "PRINCIPAL", EntityType: "Group",
				Targets: []authzpdp.Target{{Scope: "default"}, target},
			},
		},
	}

	policy := &v1alpha1.AuthzPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "hotel-booking-read", Namespace: ns},
		Spec: v1alpha1.AuthzPolicySpec{
			Context: &devCtx,
			Type: authzpolicy.Type{
				EntityID: "policy.hotel-booking-read", Name: "Hotel booking", PolicyType: "API",
				PolicyText: "permit(principal, action, resource);", Targets: []authzpdp.Target{target},
			},
		},
	}

	// cluster indexes the dependents with the functions InitCache registers, as the manager's
	// cache does.
	cluster := func(objects ...client.Object) {
		c := &recordingCache{indexers: map[string]client.IndexerFunc{}}
		Expect(search.InitCache(ctx, c)).To(Succeed())
		scheme := runtime.NewScheme()
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		schemaField, entityField := search.AuthzSchemaAuthzPdpField.String(), search.AuthzEntityAuthzPdpField.String()
		policyField := search.AuthzPolicyAuthzPdpField.String()
		k8s.RegisterClient(fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(objects...).
			WithIndex(&v1alpha1.AuthzSchema{}, schemaField, c.indexers[schemaField]).
			WithIndex(&v1alpha1.AuthzEntity{}, entityField, c.indexers[entityField]).
			WithIndex(&v1alpha1.AuthzPolicy{}, policyField, c.indexers[policyField]).
			Build())
	}

	It("warns about a schema without the 409 sentence: the platform lets a schema go untargeted", func() {
		cluster(schema)

		warnings, err := ctrl.ValidateDelete(ctx, pdp)

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ConsistOf("authz pdp [travel-team/orders-eu] is still targeted by AuthzSchema " +
			"[travel-team/travel-schema]: deleting it removes orders@eu from their targets, and each waits on " +
			"ResolvedRefs=False until its pdpRef names another AuthzPdp"))
	})

	It("warns about an entity with the 409 sentence", func() {
		cluster(entity)

		warnings, err := ctrl.ValidateDelete(ctx, pdp)

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ConsistOf("authz pdp [travel-team/orders-eu] is still targeted by AuthzEntity " +
			"[travel-team/agents-group]: deleting it removes orders@eu from their targets, and each waits on " +
			"ResolvedRefs=False until its pdpRef names another AuthzPdp. " + refused))
	})

	It("warns about a policy with the 409 sentence", func() {
		cluster(policy)

		warnings, err := ctrl.ValidateDelete(ctx, pdp)

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ConsistOf("authz pdp [travel-team/orders-eu] is still targeted by AuthzPolicy " +
			"[travel-team/hotel-booking-read]: deleting it removes orders@eu from their targets, and each waits on " +
			"ResolvedRefs=False until its pdpRef names another AuthzPdp. " + refused))
	})

	It("warns once per kind, the schema first", func() {
		cluster(schema, entity, policy)

		warnings, err := ctrl.ValidateDelete(ctx, pdp)

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(HaveLen(3))
		Expect(warnings[0]).To(ContainSubstring("AuthzSchema [travel-team/travel-schema]"))
		Expect(warnings[0]).ToNot(ContainSubstring("409"))
		Expect(warnings[1]).To(ContainSubstring("AuthzEntity [travel-team/agents-group]"))
		Expect(warnings[1]).To(HaveSuffix(refused))
		Expect(warnings[2]).To(ContainSubstring("AuthzPolicy [travel-team/hotel-booking-read]"))
		Expect(warnings[2]).To(HaveSuffix(refused))
	})
})
