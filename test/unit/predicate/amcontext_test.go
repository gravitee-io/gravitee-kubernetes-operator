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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/am"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func amContext(url string) *v1alpha1.AMContext {
	return &v1alpha1.AMContext{
		ObjectMeta: metav1.ObjectMeta{Name: "ctx", Namespace: "default"},
		Spec: v1alpha1.AMContextSpec{
			Context: &am.Context{BaseUrl: url, OrgID: "DEFAULT", EnvID: "DEFAULT"},
		},
	}
}

var _ = Describe("LastSpecHashPredicate AMContext", func() {
	p := predicate.LastSpecHashPredicate{}

	It("reconciles a create with no last-spec-hash", func() {
		Expect(p.Create(event.CreateEvent{Object: amContext("http://am")})).To(BeTrue())
	})

	It("skips a create whose last-spec-hash already matches the spec", func() {
		obj := amContext("http://am")
		obj.Annotations = map[string]string{
			core.LastSpecHashAnnotation: hash.Calculate(&obj.Spec),
		}
		Expect(p.Create(event.CreateEvent{Object: obj})).To(BeFalse())
	})

	It("reconciles an update when the spec changes", func() {
		oldObj := amContext("http://am")
		newObj := amContext("http://am-other")
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeTrue())
	})

	It("skips an update when the spec is unchanged", func() {
		oldObj := amContext("http://am")
		newObj := amContext("http://am")
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeFalse())
	})
})
