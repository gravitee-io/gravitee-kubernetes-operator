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
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpolicy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/watch"
)

// A policy that has moved off a PDP frees it as an entity does: once the platform holds its new
// targets, its status write queues every PDP being deleted, and the refused delete is retried.
var _ = Describe("WatchAuthzPdpReleases for an AuthzPolicy", func() {
	ctx := context.Background()

	retargeted := &v1alpha1.AuthzPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "hotel-booking-read", Namespace: "orders"},
		Spec: v1alpha1.AuthzPolicySpec{Type: authzpolicy.Type{
			EntityID: "policy.hotel-booking-read", PolicyType: "API",
			Targets: []authzpdp.Target{{Scope: "default"}},
		}},
	}
	releases := func() *watch.Type {
		deleting := &v1alpha1.AuthzPdp{ObjectMeta: metav1.ObjectMeta{
			Name: "orders-eu", Namespace: "orders",
			Finalizers:        []string{core.AuthzPdpFinalizer},
			DeletionTimestamp: &metav1.Time{Time: time.Now()},
		}}
		live := &v1alpha1.AuthzPdp{ObjectMeta: metav1.ObjectMeta{Name: "orders-us", Namespace: "orders"}}
		return watch.New(ctx, cluster(deleting, live), &v1alpha1.AuthzPdpList{})
	}
	freed := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "orders", Name: "orders-eu"}}

	It("retries the PDP being deleted when a policy is updated", func() {
		q := newQueue()

		releases().WatchAuthzPdpReleases().Update(ctx, event.UpdateEvent{ObjectOld: retargeted, ObjectNew: retargeted}, q)

		Expect(queued(q)).To(ConsistOf(freed))
	})

	It("retries the PDP being deleted when a policy is deleted", func() {
		q := newQueue()

		releases().WatchAuthzPdpReleases().Delete(ctx, event.DeleteEvent{Object: retargeted}, q)

		Expect(queued(q)).To(ConsistOf(freed))
	})
})
