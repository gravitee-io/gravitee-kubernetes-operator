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

package apim_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/controllers/apim/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

var _ = Describe("Authz entity reconcile", func() {
	key := types.NamespacedName{Namespace: "gravitee", Name: "agents-group"}
	req := ctrl.Request{NamespacedName: key}

	setup := func(objects ...client.Object) (*authzentity.Reconciler, client.Client) {
		scheme := runtime.NewScheme()
		Expect(clientgoscheme.AddToScheme(scheme)).To(Succeed())
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		// The manager's cache sets the kind on what Get returns, and an Update clears it when it
		// decodes the API server's answer. The fake client does neither.
		cli := fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(objects...).
			WithStatusSubresource(&v1alpha1.AuthzEntity{}).
			WithInterceptorFuncs(interceptor.Funcs{
				Get: func(ctx context.Context, c client.WithWatch, key client.ObjectKey, obj client.Object,
					opts ...client.GetOption) error {
					if err := c.Get(ctx, key, obj, opts...); err != nil {
						return err
					}
					gvk, err := apiutil.GVKForObject(obj, scheme)
					obj.GetObjectKind().SetGroupVersionKind(gvk)
					return err
				},
				Update: func(ctx context.Context, c client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
					err := c.Update(ctx, obj, opts...)
					obj.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{})
					return err
				},
			}).
			Build()
		k8s.RegisterClient(cli)
		return &authzentity.Reconciler{Scheme: scheme, Client: cli, Recorder: record.NewFakeRecorder(20)}, cli
	}

	stored := func(cli client.Client) *v1alpha1.AuthzEntity {
		e := &v1alpha1.AuthzEntity{}
		Expect(cli.Get(context.Background(), key, e)).To(Succeed())
		return e
	}

	condition := func(cli client.Client, conditionType string) *metav1.Condition {
		c := k8s.GetCondition(stored(cli), conditionType)
		Expect(c).ToNot(BeNil())
		return c
	}

	twoPdps := func() *v1alpha1.AuthzEntity {
		entity := newAuthzEntity()
		entity.Spec.Targets = []authzpdp.Target{
			{PdpRef: &refs.NamespacedName{Name: "orders-eu"}},
			{Scope: "default"},
			{PdpRef: &refs.NamespacedName{Name: "orders-us"}},
		}
		return entity
	}

	unsynced := func(name string) *v1alpha1.AuthzPdp {
		pdp := newAuthzPdp()
		pdp.Name = name
		return pdp
	}

	It("waits on ResolvedRefs=False naming every pending AuthzPdp", func() {
		r, cli := setup(twoPdps(), unsynced("orders-us"))

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).To(HaveOccurred())

		e := stored(cli)
		Expect(e.Finalizers).To(ContainElement(core.AuthzEntityFinalizer))
		resolved := condition(cli, k8s.ConditionResolvedRefs)
		Expect(resolved.Status).To(Equal(metav1.ConditionFalse))
		Expect(resolved.Message).To(ContainSubstring("authz pdp not found: [gravitee/orders-eu]; " +
			"authz pdp not synced: [gravitee/orders-us] has no routing scope yet"))
		Expect(e.Status.ID).To(BeEmpty())
	})

	It("reports an AuthzPdp of another context on Accepted, even beside a pending one", func() {
		prod := newAuthzPdp()
		prod.Name = "orders-us"
		prod.Spec.Context = &refs.NamespacedName{Name: "prod-ctx"}
		prod.Status.RoutingScope = "orders@us"
		r, cli := setup(twoPdps(), prod)

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).To(HaveOccurred())

		accepted := condition(cli, k8s.ConditionAccepted)
		Expect(accepted.Status).To(Equal(metav1.ConditionFalse))
		Expect(accepted.Message).To(ContainSubstring(
			"authz pdp [gravitee/orders-us] belongs to management context [gravitee/prod-ctx], not [gravitee/dev-ctx]"))
		Expect(k8s.GetCondition(stored(cli), k8s.ConditionResolvedRefs).Status).ToNot(Equal(metav1.ConditionFalse))
	})

	// Admission had no say, a pdpRef pending at apply time for instance: the reconcile's 400 is
	// where the ownership finding shows, and the entity is retried with backoff.
	It("keeps the platform's refusal in status.errors, naming the source, and requeues", func() {
		entity := newAuthzEntity()
		entity.Spec.Context = &refs.NamespacedName{Name: "refusing-ctx"}
		entity.Spec.Targets = []authzpdp.Target{{Scope: "default"}}
		r, cli := setup(entity)

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).To(HaveOccurred())

		e := stored(cli)
		Expect(e.Status.Errors.Severe).To(ConsistOf(ownershipFinding))
		Expect(e.Status.Errors.Severe[0]).To(ContainSubstring("(source: gravitee_am)"))
		Expect(condition(cli, k8s.ConditionAccepted).Status).To(Equal(metav1.ConditionFalse))
		Expect(e.Status.ID).To(BeEmpty())
	})

	// The entity controller applies each child of a deleted parent again (watch.WatchAuthzEntityParents).
	// The platform has detached a PRINCIPAL child: it takes the declaration back as it is.
	It("applies a PRINCIPAL again after its parent's delete", func() {
		entity := newAuthzEntity()
		entity.Spec.Context = &refs.NamespacedName{Name: "cascade-ctx"}
		entity.Spec.Targets = []authzpdp.Target{{Scope: "default"}}
		r, cli := setup(entity)

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())

		Expect(condition(cli, k8s.ConditionAccepted).Status).To(Equal(metav1.ConditionTrue))
		Expect(stored(cli).Status.ID).To(Equal("5e7b"))
	})

	// The platform has deleted a RESOURCE child with its parent, and refuses to declare it again
	// under a parent that does not exist.
	It("refuses a RESOURCE whose parent was deleted, naming that parent", func() {
		entity := newAuthzEntity()
		entity.Spec.Context = &refs.NamespacedName{Name: "cascade-ctx"}
		entity.Spec.Kind, entity.Spec.EntityType = "RESOURCE", "Hotel"
		entity.Spec.Targets = []authzpdp.Target{{Scope: "default"}}
		r, cli := setup(entity)

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).To(HaveOccurred())

		e := stored(cli)
		Expect(e.Status.Errors.Severe).To(ConsistOf(missingParentFinding(`Group::"travel"`)))
		accepted := condition(cli, k8s.ConditionAccepted)
		Expect(accepted.Status).To(Equal(metav1.ConditionFalse))
		// The message carries the platform's answer as received, its quotes escaped.
		Expect(accepted.Message).To(ContainSubstring(`parent Group::\"travel\" does not exist`))
		Expect(e.Status.ID).To(BeEmpty())
	})

	It("deletes an entity still waiting on its AuthzPdp", func() {
		r, cli := setup(twoPdps())

		_, err := r.Reconcile(context.Background(), req)
		Expect(err).To(HaveOccurred())
		Expect(cli.Delete(context.Background(), stored(cli))).To(Succeed())

		// The ManagementContext is missing, so there is no platform to delete from: the finalizer
		// goes, and the entity with it. The delete needs no target.
		_, err = r.Reconcile(context.Background(), req)
		Expect(err).ToNot(HaveOccurred())
		Expect(apierrors.IsNotFound(cli.Get(context.Background(), key, &v1alpha1.AuthzEntity{}))).To(BeTrue())
	})
})
