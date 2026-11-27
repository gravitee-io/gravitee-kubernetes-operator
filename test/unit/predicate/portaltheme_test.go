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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portaltheme"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func travelTheme(name, id string) *v1alpha1.PortalTheme {
	thm := &v1alpha1.PortalTheme{
		ObjectMeta: metav1.ObjectMeta{Name: "travel-theme", Namespace: "default"},
		Spec:       v1alpha1.PortalThemeSpec{Type: portaltheme.Type{Name: name}},
	}
	thm.Status.ID = id
	return thm
}

// A portal applied with its theme can reach APIM before the theme does and be refused: the update that gives
// the theme its APIM id changes the status only, so the theme watch must let it through.
var _ = Describe("PortalTheme watch of a portal", func() {
	p := predicate.PortalThemeChildPredicate()

	It("wakes the portal when the theme gets its APIM id", func() {
		Expect(p.Update(event.UpdateEvent{
			ObjectOld: travelTheme("Travel", ""),
			ObjectNew: travelTheme("Travel", "ea48e54c-768c-3b0d-9638-406f9e69348d"),
		})).To(BeTrue())
	})

	It("skips any other status-only update of the theme", func() {
		oldThm := travelTheme("Travel", "")
		newThm := travelTheme("Travel", "")
		newThm.Status.Conditions = []metav1.Condition{{Type: "Accepted", Status: metav1.ConditionFalse}}
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldThm, ObjectNew: newThm})).To(BeFalse())
	})
})
