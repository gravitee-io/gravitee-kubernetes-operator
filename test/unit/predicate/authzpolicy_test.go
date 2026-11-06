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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpolicy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func authzPolicy(pdpName, id string) *v1alpha1.AuthzPolicy {
	return &v1alpha1.AuthzPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "hotel-booking-read", Namespace: "default"},
		Spec: v1alpha1.AuthzPolicySpec{
			Type: authzpolicy.Type{
				EntityID:   "policy.hotel-booking-read",
				Name:       "Hotel booking",
				PolicyType: "API",
				PolicyText: "permit(principal, action, resource);",
				Targets:    []authzpdp.Target{{PdpRef: &refs.NamespacedName{Name: pdpName}}},
			},
		},
		Status: v1alpha1.AuthzPolicyStatus{Status: authzpolicy.Status{ID: id}},
	}
}

var _ = Describe("LastSpecHashPredicate AuthzPolicy", func() {
	p := predicate.LastSpecHashPredicate{}

	It("reconciles a create with no last-spec-hash", func() {
		Expect(p.Create(event.CreateEvent{Object: authzPolicy("orders-eu", "")})).To(BeTrue())
	})

	It("skips a create whose last-spec-hash already matches the spec", func() {
		obj := authzPolicy("orders-eu", "")
		obj.Annotations = map[string]string{
			core.LastSpecHashAnnotation: hash.Calculate(&obj.Spec),
		}
		Expect(p.Create(event.CreateEvent{Object: obj})).To(BeFalse())
	})

	It("reconciles an update that retargets the policy", func() {
		oldObj := authzPolicy("orders-eu", "")
		newObj := authzPolicy("orders-us", "")
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeTrue())
	})

	It("skips an update when only the deployment status changes", func() {
		oldObj := authzPolicy("orders-eu", "5f1c")
		newObj := authzPolicy("orders-eu", "5f1c")
		newObj.Status.PolicyStatus = "DISABLED"
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeFalse())
	})

	It("skips an update when only the status changes", func() {
		oldObj := authzPolicy("orders-eu", "")
		newObj := authzPolicy("orders-eu", "5f1c")
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj})).To(BeFalse())
	})
})
