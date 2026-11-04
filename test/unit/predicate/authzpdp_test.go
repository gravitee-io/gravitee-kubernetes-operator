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

package predicate_test

import (
	"time"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func authzPdp(name, scope string) *v1alpha1.AuthzPdp {
	return &v1alpha1.AuthzPdp{
		ObjectMeta: metav1.ObjectMeta{Name: "orders-eu", Namespace: "default"},
		Spec: v1alpha1.AuthzPdpSpec{
			Type: authzpdp.Type{Name: name, TargetPdpID: "orders", Tag: "eu"},
		},
		Status: v1alpha1.AuthzPdpStatus{Status: authzpdp.Status{RoutingScope: scope}},
	}
}

var _ = Describe("LastSpecHashPredicate AuthzPdp", func() {
	p := predicate.LastSpecHashPredicate{}

	It("reconciles a create with no last-spec-hash", func() {
		Expect(p.Create(event.CreateEvent{Object: authzPdp("Orders", "")})).To(BeTrue())
	})

	It("skips a create whose last-spec-hash already matches the spec", func() {
		obj := authzPdp("Orders", "")
		obj.Annotations = map[string]string{
			core.LastSpecHashAnnotation: hash.Calculate(&obj.Spec),
		}
		Expect(p.Create(event.CreateEvent{Object: obj})).To(BeFalse())
	})

	It("reconciles an update when the spec changes", func() {
		oldObj := authzPdp("Orders", "")
		newObj := authzPdp("Orders (EU)", "")
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeTrue())
	})

	It("skips an update when only the status changes", func() {
		oldObj := authzPdp("Orders", "")
		newObj := authzPdp("Orders", "orders@eu")
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeFalse())
	})
})

var _ = Describe("AuthzPdpRoutedPredicate", func() {
	p := predicate.AuthzPdpRoutedPredicate{}

	It("drops a create: a new PDP is not synced yet", func() {
		Expect(p.Create(event.CreateEvent{Object: authzPdp("Orders", "orders@eu")})).To(BeFalse())
	})

	It("passes a first sync", func() {
		oldObj := authzPdp("Orders", "")
		newObj := authzPdp("Orders", "orders@eu")
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeTrue())
	})

	It("passes a routing scope change", func() {
		oldObj := authzPdp("Orders", "orders@eu")
		newObj := authzPdp("Orders", "orders@us")
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeTrue())
	})

	It("drops an update that leaves the routing scope as it was", func() {
		oldObj := authzPdp("Orders", "orders@eu")
		newObj := authzPdp("Orders (EU)", "orders@eu")
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeFalse())
	})

	It("passes the update that marks the PDP for deletion", func() {
		oldObj := authzPdp("Orders", "orders@eu")
		newObj := authzPdp("Orders", "orders@eu")
		newObj.Finalizers = []string{core.AuthzPdpFinalizer}
		newObj.DeletionTimestamp = &metav1.Time{Time: time.Now()}
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeTrue())
	})

	It("drops a later update of a PDP already being deleted", func() {
		deleting := func() *v1alpha1.AuthzPdp {
			obj := authzPdp("Orders", "orders@eu")
			obj.Finalizers = []string{core.AuthzPdpFinalizer}
			obj.DeletionTimestamp = &metav1.Time{Time: time.Unix(1793800000, 0)}
			return obj
		}
		Expect(p.Update(event.UpdateEvent{ObjectOld: deleting(), ObjectNew: deleting()})).To(BeFalse())
	})

	It("drops an update of another kind", func() {
		oldObj := catalogMcpServer("http://mcp")
		newObj := catalogMcpServer("http://mcp-other")
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeFalse())
	})

	It("passes a delete: the platform removes the scope from every target", func() {
		Expect(p.Delete(event.DeleteEvent{Object: authzPdp("Orders", "orders@eu")})).To(BeTrue())
	})

	It("drops a generic event", func() {
		Expect(p.Generic(event.GenericEvent{Object: authzPdp("Orders", "orders@eu")})).To(BeFalse())
	})
})
