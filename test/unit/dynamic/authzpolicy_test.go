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

package dynamic

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpolicy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

var _ = Describe("Resolving the refs of an AuthzPolicy", func() {
	const ns = "travel"
	ctx := context.Background()
	devCtx := refs.NamespacedName{Name: "dev-ctx"}

	policy := func(targets ...authzpdp.Target) *v1alpha1.AuthzPolicy {
		return &v1alpha1.AuthzPolicy{
			ObjectMeta: metav1.ObjectMeta{Name: "hotel-booking-read", Namespace: ns},
			Spec: v1alpha1.AuthzPolicySpec{
				Context: &devCtx,
				Type:    authzpolicy.Type{EntityID: "policy.hotel-booking-read", Targets: targets},
			},
		}
	}
	cluster := func(objects ...client.Object) {
		scheme := runtime.NewScheme()
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build())
	}
	byRef := func(name string) authzpdp.Target {
		return authzpdp.Target{PdpRef: &refs.NamespacedName{Name: name}}
	}
	synced := &v1alpha1.AuthzPdp{
		ObjectMeta: metav1.ObjectMeta{Name: "orders-eu", Namespace: ns},
		Spec: v1alpha1.AuthzPdpSpec{
			Context: &devCtx,
			Type:    authzpdp.Type{Name: "orders-eu", TargetPdpID: "orders-eu"},
		},
		Status: v1alpha1.AuthzPdpStatus{Status: authzpdp.Status{RoutingScope: "orders@eu"}},
	}

	It("replaces each target with the scope it sends, in declaration order", func() {
		cluster(synced)
		p := policy(authzpdp.Target{Scope: "default"}, byRef("orders-eu"), authzpdp.Target{Scope: "*"})

		Expect(dynamic.ResolveAuthzPolicyRefs(ctx, p)).To(Succeed())

		Expect(p.Spec.Targets).To(Equal([]authzpdp.Target{{Scope: "default"}, {Scope: "orders@eu"}, {Scope: "*"}}))
	})

	It("leaves a policy without targets without any", func() {
		cluster()
		p := policy()

		Expect(dynamic.ResolveAuthzPolicyRefs(ctx, p)).To(Succeed())

		Expect(p.Spec.Targets).To(BeNil())
	})

	It("leaves the targets as they were while an AuthzPdp is pending", func() {
		cluster()
		p := policy(byRef("orders-eu"), authzpdp.Target{Scope: "default"})

		err := dynamic.ResolveAuthzPolicyRefs(ctx, p)

		Expect(dynamic.IsAuthzPdpPending(err)).To(BeTrue())
		Expect(p.Spec.Targets).To(Equal([]authzpdp.Target{byRef("orders-eu"), {Scope: "default"}}))
	})
})
