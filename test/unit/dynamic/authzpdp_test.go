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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

var _ = Describe("Resolving a pdpRef", func() {
	const ns = "orders"
	ctx := context.Background()

	pdp := func(contextRef refs.NamespacedName, scope string) *v1alpha1.AuthzPdp {
		return &v1alpha1.AuthzPdp{
			ObjectMeta: metav1.ObjectMeta{Name: "orders-eu", Namespace: ns},
			Spec: v1alpha1.AuthzPdpSpec{
				Context: &contextRef,
				Type:    authzpdp.Type{Name: "Orders PDP (EU)", TargetPdpID: "orders", Tag: "eu"},
			},
			Status: v1alpha1.AuthzPdpStatus{Status: authzpdp.Status{RoutingScope: scope}},
		}
	}
	// dependent stands in for a schema, an entity or a policy: any resource with a context.
	dependent := func(contextRef refs.NamespacedName) core.ContextAwareObject {
		return &v1alpha1.McpProxy{
			ObjectMeta: metav1.ObjectMeta{Name: "orders-schema", Namespace: ns},
			Spec:       v1alpha1.McpProxySpec{Context: &contextRef},
		}
	}
	cluster := func(objects ...client.Object) {
		scheme := runtime.NewScheme()
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build())
	}
	devCtx := refs.NamespacedName{Name: "dev-ctx"}
	ref := refs.NamespacedName{Name: "orders-eu"}

	It("returns the routing scope of a synced PDP, the namespace defaulted to the dependent's", func() {
		cluster(pdp(devCtx, "orders@eu"))

		scope, err := dynamic.ResolveAuthzPdpRoutingScope(ctx, dependent(devCtx), ref)

		Expect(err).ToNot(HaveOccurred())
		Expect(scope).To(Equal("orders@eu"))
	})

	It("matches a context named with and without its namespace", func() {
		cluster(pdp(refs.NamespacedName{Name: "dev-ctx", Namespace: ns}, "orders@eu"))

		scope, err := dynamic.ResolveAuthzPdpRoutingScope(ctx, dependent(devCtx), ref)

		Expect(err).ToNot(HaveOccurred())
		Expect(scope).To(Equal("orders@eu"))
	})

	It("waits for a PDP that does not exist", func() {
		cluster()

		_, err := dynamic.ResolveAuthzPdpRoutingScope(ctx, dependent(devCtx), ref)

		Expect(err).To(MatchError(dynamic.ErrAuthzPdpNotFound))
		Expect(err).To(MatchError(ContainSubstring("[orders/orders-eu]")))
		Expect(dynamic.IsAuthzPdpPending(err)).To(BeTrue())
	})

	It("waits for a PDP in another namespace that does not exist", func() {
		cluster(pdp(devCtx, "orders@eu"))

		_, err := dynamic.ResolveAuthzPdpRoutingScope(ctx, dependent(devCtx),
			refs.NamespacedName{Name: "orders-eu", Namespace: "billing"})

		Expect(err).To(MatchError(dynamic.ErrAuthzPdpNotFound))
	})

	It("treats a PDP being deleted as missing", func() {
		deleting := pdp(devCtx, "orders@eu")
		deleting.Finalizers = []string{core.AuthzPdpFinalizer}
		deleting.DeletionTimestamp = &metav1.Time{Time: time.Now()}
		cluster(deleting)

		_, err := dynamic.ResolveAuthzPdpRoutingScope(ctx, dependent(devCtx), ref)

		Expect(err).To(MatchError(dynamic.ErrAuthzPdpNotFound))
		Expect(err).To(MatchError(ContainSubstring("is being deleted")))
		Expect(dynamic.IsAuthzPdpPending(err)).To(BeTrue())
	})

	It("waits for a PDP the platform has not synced yet", func() {
		cluster(pdp(devCtx, ""))

		_, err := dynamic.ResolveAuthzPdpRoutingScope(ctx, dependent(devCtx), ref)

		Expect(err).To(MatchError(dynamic.ErrAuthzPdpNotSynced))
		Expect(dynamic.IsAuthzPdpPending(err)).To(BeTrue())
	})

	It("refuses a PDP synced through another management context", func() {
		cluster(pdp(refs.NamespacedName{Name: "prod-ctx"}, "orders@eu"))

		_, err := dynamic.ResolveAuthzPdpRoutingScope(ctx, dependent(devCtx), ref)

		Expect(err).To(MatchError(
			"authz pdp [orders/orders-eu] is synced through management context [orders/prod-ctx], not [orders/dev-ctx]"))
		Expect(dynamic.IsAuthzPdpPending(err)).To(BeFalse())
	})
})
