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

package watch_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
)

const ns = "travel"

var devCtx = refs.NamespacedName{Name: "dev-ctx"}

func newEntity(name, kind, entityType string, contextRef refs.NamespacedName, parents ...authzentity.Parent) *v1alpha1.AuthzEntity {
	return &v1alpha1.AuthzEntity{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
		Spec: v1alpha1.AuthzEntitySpec{
			Context: &contextRef,
			Type:    authzentity.Type{EntityID: name, Kind: kind, EntityType: entityType, Parents: parents},
		},
	}
}

func request(name string) reconcile.Request {
	return reconcile.Request{NamespacedName: types.NamespacedName{Namespace: ns, Name: name}}
}

var _ = Describe("Watching AuthzEntity parents", func() {
	ctx := context.Background()
	agents := authzentity.Parent{EntityType: "Group", EntityID: "agents"}

	It("queues, on a parent's delete, each entity of its context that declares it", func() {
		parent := newEntity("agents", "PRINCIPAL", "Group", devCtx)
		cli := cluster(
			newEntity("alice", "PRINCIPAL", "User", devCtx, agents),
			newEntity("ritz", "RESOURCE", "Hotel", refs.NamespacedName{Name: "dev-ctx", Namespace: ns},
				authzentity.Parent{EntityType: "Group", EntityID: "travel"}, agents),
			newEntity("savoy", "RESOURCE", "Hotel", refs.NamespacedName{Name: "prod-ctx"}, agents),
			newEntity("negresco", "RESOURCE", "Hotel", devCtx, authzentity.Parent{EntityType: "Team", EntityID: "agents"}),
		)
		q := newQueue()

		watch.New(ctx, cli, &v1alpha1.AuthzEntityList{}).WatchAuthzEntityParents().
			Delete(ctx, event.DeleteEvent{Object: parent}, q)

		Expect(queued(q)).To(ConsistOf(request("alice"), request("ritz")))
	})

	// The platform deletes a RESOURCE child with its parent, and that child's children in turn: a
	// RESOURCE child is deleted, a PRINCIPAL child detached. It does not go through a PRINCIPAL.
	It("queues, on a parent's delete, the entities its RESOURCE children cascade to", func() {
		parent := newEntity("agents", "PRINCIPAL", "Group", devCtx)
		ritz := authzentity.Parent{EntityType: "Hotel", EntityID: "ritz"}
		alice := authzentity.Parent{EntityType: "User", EntityID: "alice"}
		cli := cluster(
			newEntity("ritz", "RESOURCE", "Hotel", devCtx, agents),
			newEntity("ritz-suite", "RESOURCE", "Room", devCtx, ritz),
			newEntity("ritz-suite-minibar", "RESOURCE", "Minibar", devCtx, authzentity.Parent{EntityType: "Room", EntityID: "ritz-suite"}),
			newEntity("concierge", "PRINCIPAL", "User", devCtx, ritz),
			newEntity("alice", "PRINCIPAL", "User", devCtx, agents),
			newEntity("alice-laptop", "RESOURCE", "Device", devCtx, alice),
		)
		q := newQueue()

		watch.New(ctx, cli, &v1alpha1.AuthzEntityList{}).WatchAuthzEntityParents().
			Delete(ctx, event.DeleteEvent{Object: parent}, q)

		Expect(queued(q)).To(ConsistOf(
			request("ritz"), request("ritz-suite"), request("ritz-suite-minibar"), request("concierge"), request("alice"),
		))
	})

	// The platform refuses a cycle, but the dry run that would say so is skipped while an AuthzPdp is
	// pending: the cluster can hold one.
	It("queues each entity once when RESOURCE parents loop", func() {
		parent := newEntity("agents", "PRINCIPAL", "Group", devCtx)
		cli := cluster(
			newEntity("ritz", "RESOURCE", "Hotel", devCtx, agents, authzentity.Parent{EntityType: "Room", EntityID: "ritz-suite"}),
			newEntity("ritz-suite", "RESOURCE", "Room", devCtx, authzentity.Parent{EntityType: "Hotel", EntityID: "ritz"}),
		)
		q := newQueue()
		done := make(chan struct{})

		go func() {
			defer GinkgoRecover()
			defer close(done)
			watch.New(ctx, cli, &v1alpha1.AuthzEntityList{}).WatchAuthzEntityParents().
				Delete(ctx, event.DeleteEvent{Object: parent}, q)
		}()

		Eventually(done).Should(BeClosed(), "the cascade search ends")
		Expect(queued(q)).To(ConsistOf(request("ritz"), request("ritz-suite")))
	})

	It("queues nothing on a create or an update", func() {
		parent := newEntity("agents", "PRINCIPAL", "Group", devCtx)
		cli := cluster(newEntity("alice", "PRINCIPAL", "User", devCtx, agents))
		q := newQueue()
		handler := watch.New(ctx, cli, &v1alpha1.AuthzEntityList{}).WatchAuthzEntityParents()

		handler.Create(ctx, event.CreateEvent{Object: parent}, q)
		handler.Update(ctx, event.UpdateEvent{ObjectOld: parent, ObjectNew: parent}, q)

		Expect(queued(q)).To(BeEmpty())
	})
})
