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

package predicate_test

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func withAnnotation(obj client.Object, key, value string) client.Object {
	if value != "" {
		obj.SetAnnotations(map[string]string{key: value})
	}
	return obj
}

func apiV4(key, value string) client.Object {
	return withAnnotation(&v1alpha1.ApiV4Definition{ObjectMeta: metav1.ObjectMeta{Name: "api", Namespace: "default"}}, key, value)
}

func application(key, value string) client.Object {
	return withAnnotation(&v1alpha1.Application{ObjectMeta: metav1.ObjectMeta{Name: "app", Namespace: "default"}}, key, value)
}

func grp(key, value string) client.Object {
	return withAnnotation(&v1alpha1.Group{
		ObjectMeta: metav1.ObjectMeta{Name: "grp", Namespace: "default"},
		Spec:       v1alpha1.GroupSpec{Type: &group.Type{Name: "developers"}},
	}, key, value)
}

var _ = Describe("LastSpecHashPredicate ignore annotations", func() {
	p := predicate.LastSpecHashPredicate{}

	DescribeTable("reconciles an update that only changes the annotation",
		func(build func(string, string) client.Object, key, oldValue, newValue string) {
			Expect(p.Update(event.UpdateEvent{ObjectOld: build(key, oldValue), ObjectNew: build(key, newValue)})).To(BeTrue())
		},
		Entry("ApiV4Definition ignore-groups added", apiV4, core.IgnoreGroupsAnnotation, "", "true"),
		Entry("ApiV4Definition ignore-groups removed", apiV4, core.IgnoreGroupsAnnotation, "true", ""),
		Entry("Application ignore-groups added", application, core.IgnoreGroupsAnnotation, "", "true"),
		Entry("Group ignore-members flipped", grp, core.IgnoreMembersAnnotation, "true", "false"),
	)

	DescribeTable("skips an update when neither spec nor annotation changes",
		func(build func(string, string) client.Object, key string) {
			Expect(p.Update(event.UpdateEvent{ObjectOld: build(key, "true"), ObjectNew: build(key, "true")})).To(BeFalse())
		},
		Entry("ApiV4Definition", apiV4, core.IgnoreGroupsAnnotation),
		Entry("Application", application, core.IgnoreGroupsAnnotation),
		Entry("Group", grp, core.IgnoreMembersAnnotation),
	)

	It("ignores an unrelated annotation", func() {
		Expect(p.Update(event.UpdateEvent{ObjectOld: apiV4("other", ""), ObjectNew: apiV4("other", "x")})).To(BeFalse())
	})
})

var _ = DescribeTable("k8s.HasTrueAnnotation",
	func(value string, expected bool) {
		Expect(k8s.HasTrueAnnotation(apiV4(core.IgnoreGroupsAnnotation, value), core.IgnoreGroupsAnnotation)).To(Equal(expected))
	},
	Entry("true", "true", true),
	Entry("false", "false", false),
	Entry("absent", "", false),
	Entry("not a boolean", "yes", false),
)
