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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
)

// syncedAPIV4 is an API whose last-spec-hash matches its spec, with Accepted set to accepted.
func syncedAPIV4(accepted metav1.ConditionStatus) *v1alpha1.ApiV4Definition {
	api := &v1alpha1.ApiV4Definition{ObjectMeta: metav1.ObjectMeta{Name: "api", Namespace: "default"}}
	api.Annotations = map[string]string{core.LastSpecHashAnnotation: hash.Calculate(&api.Spec)}
	api.Status.Conditions = []metav1.Condition{{Type: k8s.ConditionAccepted, Status: accepted}}
	return api
}

var _ = Describe("LastSpecHashPredicate Create on an ApiV4Definition", func() {
	p := predicate.LastSpecHashPredicate{}

	It("skips an API synced and accepted, replayed when the operator starts", func() {
		Expect(p.Create(event.CreateEvent{Object: syncedAPIV4(metav1.ConditionTrue)})).To(BeFalse())
	})

	It("replays an API whose last sync was not accepted", func() {
		Expect(p.Create(event.CreateEvent{Object: syncedAPIV4(metav1.ConditionFalse)})).To(BeTrue())
	})

	It("replays an API with no Accepted condition yet", func() {
		api := syncedAPIV4(metav1.ConditionTrue)
		api.Status.Conditions = nil
		Expect(p.Create(event.CreateEvent{Object: api})).To(BeTrue())
	})

	It("ignores a stale processingStatus", func() {
		api := syncedAPIV4(metav1.ConditionTrue)
		api.Status.ProcessingStatus = core.ProcessingStatusFailed
		Expect(p.Create(event.CreateEvent{Object: api})).To(BeFalse())
	})
})

var _ = Describe("ApiV4DefinitionStatus IsFailed", func() {
	It("is failed when Accepted is False", func() {
		Expect(syncedAPIV4(metav1.ConditionFalse).Status.IsFailed()).To(BeTrue())
	})

	It("is not failed when only ResolvedRefs is False (a missing group ref)", func() {
		api := syncedAPIV4(metav1.ConditionTrue)
		api.Status.Conditions = append(api.Status.Conditions,
			metav1.Condition{Type: k8s.ConditionResolvedRefs, Status: metav1.ConditionFalse})
		Expect(api.Status.IsFailed()).To(BeFalse())
	})
})
