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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzschema"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s/dynamic"
)

var _ = Describe("Resolving the refs of an AuthzSchema", func() {
	const ns = "travel"
	ctx := context.Background()
	devCtx := refs.NamespacedName{Name: "dev-ctx"}

	schema := func(target *authzpdp.Target) *v1alpha1.AuthzSchema {
		return &v1alpha1.AuthzSchema{
			ObjectMeta: metav1.ObjectMeta{Name: "travel-schema", Namespace: ns},
			Spec: v1alpha1.AuthzSchemaSpec{
				Context: &devCtx,
				Type:    authzschema.Type{Name: "Travel schema", Target: target, Schema: "entity User;\n"},
			},
		}
	}
	cluster := func(objects ...client.Object) {
		scheme := runtime.NewScheme()
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build())
	}
	byRef := func() *authzpdp.Target {
		return &authzpdp.Target{PdpRef: &refs.NamespacedName{Name: "orders-eu"}}
	}
	synced := &v1alpha1.AuthzPdp{
		ObjectMeta: metav1.ObjectMeta{Name: "orders-eu", Namespace: ns},
		Spec: v1alpha1.AuthzPdpSpec{
			Context: &devCtx,
			Type:    authzpdp.Type{Name: "Orders PDP (EU)", TargetPdpID: "orders", Tag: "eu"},
		},
		Status: v1alpha1.AuthzPdpStatus{Status: authzpdp.Status{RoutingScope: "orders@eu"}},
	}

	It("replaces a pdpRef with the routing scope of its AuthzPdp", func() {
		cluster(synced)
		s := schema(byRef())

		Expect(dynamic.ResolveAuthzSchemaRefs(ctx, s)).To(Succeed())

		Expect(s.Spec.Target).To(Equal(&authzpdp.Target{Scope: "orders@eu"}))
	})

	It("keeps a scope as written", func() {
		cluster()
		s := schema(&authzpdp.Target{Scope: "default"})

		Expect(dynamic.ResolveAuthzSchemaRefs(ctx, s)).To(Succeed())

		Expect(s.Spec.Target).To(Equal(&authzpdp.Target{Scope: "default"}))
	})

	It("leaves a schema without a target without one", func() {
		cluster()
		s := schema(nil)

		Expect(dynamic.ResolveAuthzSchemaRefs(ctx, s)).To(Succeed())

		Expect(s.Spec.Target).To(BeNil())
	})

	It("leaves the target as it was while the AuthzPdp is pending", func() {
		cluster()
		s := schema(byRef())

		err := dynamic.ResolveAuthzSchemaRefs(ctx, s)

		Expect(dynamic.IsAuthzPdpPending(err)).To(BeTrue())
		Expect(s.Spec.Target).To(Equal(byRef()))
	})
})
