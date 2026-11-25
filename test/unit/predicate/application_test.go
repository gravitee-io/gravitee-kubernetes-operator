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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

// syncedApplication is an application whose last-spec-hash matches its spec, with the given conditions.
func syncedApplication(conditions ...metav1.Condition) *v1alpha1.Application {
	app := &v1alpha1.Application{ObjectMeta: metav1.ObjectMeta{Name: "app", Namespace: "default"}}
	app.Spec.Name = "app"
	app.Annotations = map[string]string{core.LastSpecHashAnnotation: hash.Calculate(&app.Spec)}
	app.Status.Conditions = conditions
	return app
}

var _ = Describe("LastSpecHashPredicate on an Application create, replayed when the operator starts", func() {
	p := predicate.LastSpecHashPredicate{}

	accepted := metav1.Condition{Type: "Accepted", Status: metav1.ConditionTrue}
	refused := metav1.Condition{Type: "Accepted", Status: metav1.ConditionFalse}

	DescribeTable("reconciles an application unless its last sync succeeded",
		func(app *v1alpha1.Application, reconcile bool) {
			Expect(p.Create(event.CreateEvent{Object: app})).To(Equal(reconcile))
		},
		Entry("synced and accepted", syncedApplication(accepted), false),
		Entry("synced but refused", syncedApplication(refused), true),
		Entry("synced, with no condition written yet", syncedApplication(), true),
		Entry("spec changed since the last sync", func() *v1alpha1.Application {
			app := syncedApplication(accepted)
			app.Spec.Name = "renamed"
			return app
		}(), true),
	)
})
