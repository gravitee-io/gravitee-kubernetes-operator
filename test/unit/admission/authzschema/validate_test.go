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

package authzschema_test

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
	admission "github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/authzschema"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

var _ = Describe("Validating an AuthzSchema target", func() {
	const ns = "travel"
	ctx := context.Background()
	ctrl := admission.AdmissionCtrl{}
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
	byRef := func(ref refs.NamespacedName) *authzpdp.Target {
		return &authzpdp.Target{PdpRef: &ref}
	}
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
	cluster := func(objects ...client.Object) {
		scheme := runtime.NewScheme()
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		k8s.RegisterClient(fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build())
	}

	BeforeEach(func() {
		platformCalls.Store(0)
	})

	It("warns and skips the dry run when the AuthzPdp does not exist yet", func() {
		cluster()

		warnings, err := ctrl.ValidateCreate(ctx, schema(byRef(refs.NamespacedName{Name: "orders-eu"})))

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ConsistOf(
			"platform checks deferred (authz pdp not found: [travel/orders-eu])"))
		Expect(platformCalls.Load()).To(BeZero())
	})

	It("warns and skips the dry run when the AuthzPdp has not synced yet", func() {
		cluster(pdp(devCtx, ""))

		warnings, err := ctrl.ValidateCreate(ctx, schema(byRef(refs.NamespacedName{Name: "orders-eu"})))

		Expect(err).ToNot(HaveOccurred())
		Expect(warnings).To(ConsistOf(ContainSubstring("authz pdp not synced: [travel/orders-eu]")))
		Expect(platformCalls.Load()).To(BeZero())
	})

	// The controller's own finalizer and annotation writes go through the webhook: refusing them
	// while the AuthzPdp is pending would stop the controller from saving its finalizer.
	It("admits a metadata-only update while the AuthzPdp is pending", func() {
		cluster()
		oldObj := schema(byRef(refs.NamespacedName{Name: "orders-eu"}))
		oldObj.Status.ID = "7a1e0c2b"
		newObj := oldObj.DeepCopy()
		newObj.Finalizers = []string{"finalizers.gravitee.io/authzschemas"}

		_, err := ctrl.ValidateUpdate(ctx, oldObj, newObj)

		Expect(err).ToNot(HaveOccurred())
		Expect(platformCalls.Load()).To(BeZero())
	})

	It("refuses a pdpRef naming another kind", func() {
		cluster()

		_, err := ctrl.ValidateCreate(ctx, schema(byRef(refs.NamespacedName{Name: "orders-eu", Kind: "McpProxy"})))

		Expect(err).To(MatchError(ContainSubstring("target.pdpRef [travel/orders-eu]: kind must be AuthzPdp, got [McpProxy]")))
		Expect(platformCalls.Load()).To(BeZero())
	})

	It("refuses an AuthzPdp synced through another management context", func() {
		cluster(pdp(refs.NamespacedName{Name: "prod-ctx"}, "orders@eu"))

		_, err := ctrl.ValidateCreate(ctx, schema(byRef(refs.NamespacedName{Name: "orders-eu"})))

		Expect(err).To(MatchError(ContainSubstring(
			"authz pdp [travel/orders-eu] belongs to management context [travel/prod-ctx], not [travel/dev-ctx]")))
		Expect(platformCalls.Load()).To(BeZero())
	})

	// A warning would admit a schema that waits on ResolvedRefs=False forever: syncing the
	// AuthzPdp does not change its context.
	It("refuses an AuthzPdp of another management context that has not synced yet", func() {
		cluster(pdp(refs.NamespacedName{Name: "prod-ctx"}, ""))

		warnings, err := ctrl.ValidateCreate(ctx, schema(byRef(refs.NamespacedName{Name: "orders-eu"})))

		Expect(err).To(MatchError(ContainSubstring("belongs to management context [travel/prod-ctx]")))
		Expect(warnings).To(BeEmpty())
		Expect(platformCalls.Load()).To(BeZero())
	})

	// The stub answers 500: reaching it proves the dry run runs once the target resolves.
	It("dry-runs a resolved target against the platform", func() {
		cluster(pdp(devCtx, "orders@eu"))

		_, err := ctrl.ValidateCreate(ctx, schema(byRef(refs.NamespacedName{Name: "orders-eu"})))

		Expect(err).To(HaveOccurred())
		Expect(platformCalls.Load()).To(BeNumerically(">", 0))
	})
})
