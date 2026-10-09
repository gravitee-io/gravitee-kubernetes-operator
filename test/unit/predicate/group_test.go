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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func groupNamed(name string) *v1alpha1.Group {
	return &v1alpha1.Group{
		ObjectMeta: metav1.ObjectMeta{Name: "developers", Namespace: "default"},
		Spec:       v1alpha1.GroupSpec{Type: &group.Type{Name: name}},
	}
}

func synced(obj *v1alpha1.Group) *v1alpha1.Group {
	synced := obj.DeepCopy()
	synced.Annotations = map[string]string{core.LastSpecHashAnnotation: hash.Calculate(&synced.Spec)}
	return synced
}

var _ = Describe("GroupSyncedPredicate", func() {
	p := predicate.GroupSyncedPredicate{}

	It("skips a create, the group is not in APIM yet", func() {
		Expect(p.Create(event.CreateEvent{Object: groupNamed("developers")})).To(BeFalse())
	})

	It("passes the first sync, which adds last-spec-hash and leaves the spec unchanged", func() {
		created := groupNamed("developers")
		Expect(p.Update(event.UpdateEvent{ObjectOld: created, ObjectNew: synced(created)})).To(BeTrue())
	})

	It("passes a synced rename, which changes last-spec-hash", func() {
		oldObj := synced(groupNamed("developers"))
		newObj := synced(groupNamed("engineers"))
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeTrue())
	})

	It("skips a status-only update", func() {
		oldObj := synced(groupNamed("developers"))
		newObj := oldObj.DeepCopy()
		newObj.Status.ID = "group-id"
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeFalse())
	})

	It("skips a spec change not synced yet", func() {
		oldObj := synced(groupNamed("developers"))
		newObj := oldObj.DeepCopy()
		newObj.Spec.Name = "engineers"
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeFalse())
	})

	It("skips objects that are not groups", func() {
		oldObj := catalogMcpServer("http://mcp")
		newObj := catalogMcpServer("http://mcp")
		newObj.Annotations = map[string]string{core.LastSpecHashAnnotation: "hash"}
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeFalse())
	})
})
