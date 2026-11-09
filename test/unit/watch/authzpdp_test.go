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

package watch_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
)

// The dependent is an entity that has moved off the PDP: once the platform holds its new targets,
// its update names the PDP no more, so the handler does not read it.
var _ = Describe("WatchAuthzPdpReleases", func() {
	ctx := context.Background()

	pdp := func(namespace, name string, deleting bool) *v1alpha1.AuthzPdp {
		p := &v1alpha1.AuthzPdp{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace}}
		if deleting {
			p.Finalizers = []string{core.AuthzPdpFinalizer}
			p.DeletionTimestamp = &metav1.Time{Time: time.Now()}
		}
		return p
	}
	dependent := &v1alpha1.AuthzEntity{
		ObjectMeta: metav1.ObjectMeta{Name: "agents-group", Namespace: "orders"},
		Spec: v1alpha1.AuthzEntitySpec{Type: authzentity.Type{
			EntityID: "agents", Kind: "PRINCIPAL", EntityType: "Group",
			Targets: []authzpdp.Target{{Scope: "default"}},
		}},
	}

	watcher := func(objects ...client.Object) *watch.Type {
		scheme := runtime.NewScheme()
		Expect(v1alpha1.AddToScheme(scheme)).To(Succeed())
		cli := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build()
		return watch.New(ctx, cli, &v1alpha1.AuthzPdpList{})
	}
	queued := func(q workqueue.TypedRateLimitingInterface[reconcile.Request]) []types.NamespacedName {
		names := []types.NamespacedName{}
		for q.Len() > 0 {
			req, _ := q.Get()
			names = append(names, req.NamespacedName)
			q.Done(req)
		}
		return names
	}
	newQueue := func() workqueue.TypedRateLimitingInterface[reconcile.Request] {
		return workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[reconcile.Request]())
	}
	beingDeleted := []types.NamespacedName{
		{Namespace: "billing", Name: "billing-eu"},
		{Namespace: "orders", Name: "orders-eu"},
	}
	cluster := []client.Object{
		pdp("orders", "orders-eu", true),
		pdp("orders", "orders-us", false),
		pdp("billing", "billing-eu", true),
	}

	It("retries every PDP being deleted, in any namespace, when a dependent is updated", func() {
		q := newQueue()
		defer q.ShutDown()

		watcher(cluster...).WatchAuthzPdpReleases().Update(ctx,
			event.UpdateEvent{ObjectOld: dependent, ObjectNew: dependent}, q)

		Expect(queued(q)).To(ConsistOf(beingDeleted))
	})

	It("retries every PDP being deleted when a dependent is deleted", func() {
		q := newQueue()
		defer q.ShutDown()

		watcher(cluster...).WatchAuthzPdpReleases().Delete(ctx, event.DeleteEvent{Object: dependent}, q)

		Expect(queued(q)).To(ConsistOf(beingDeleted))
	})

	It("queues nothing while no PDP is being deleted", func() {
		q := newQueue()
		defer q.ShutDown()

		watcher(pdp("orders", "orders-us", false)).WatchAuthzPdpReleases().Update(ctx,
			event.UpdateEvent{ObjectOld: dependent, ObjectNew: dependent}, q)

		Expect(q.Len()).To(BeZero())
	})
})
