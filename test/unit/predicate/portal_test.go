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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/portal"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/predicate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func travelPortal(name, id string) *v1alpha1.Portal {
	prtl := &v1alpha1.Portal{
		ObjectMeta: metav1.ObjectMeta{Name: "travel", Namespace: "default"},
		Spec:       v1alpha1.PortalSpec{Type: portal.Type{Name: name}},
	}
	prtl.Status.ID = id
	return prtl
}

// A listing, page or link waits for its portal to be synced (GKO-3480): the update that gives the portal
// its APIM id changes the status only, so the Portal watch must let it through.
var _ = Describe("Portal watch of a listing, page or link", func() {
	p := predicate.PortalChildPredicate()

	It("wakes the children when the portal gets its APIM id", func() {
		Expect(p.Update(event.UpdateEvent{
			ObjectOld: travelPortal("Travel", ""),
			ObjectNew: travelPortal("Travel", "f9326ba3-0da0-3d4f-a7d5-f4a729b9e3fc"),
		})).To(BeTrue())
	})

	It("wakes the children when the portal spec changes", func() {
		Expect(p.Update(event.UpdateEvent{
			ObjectOld: travelPortal("Travel", "f9326ba3-0da0-3d4f-a7d5-f4a729b9e3fc"),
			ObjectNew: travelPortal("Travel and leisure", "f9326ba3-0da0-3d4f-a7d5-f4a729b9e3fc"),
		})).To(BeTrue())
	})

	It("skips any other status-only update of the portal", func() {
		oldPrtl := travelPortal("Travel", "")
		newPrtl := travelPortal("Travel", "")
		newPrtl.Status.Conditions = []metav1.Condition{{Type: "Accepted", Status: metav1.ConditionFalse}}
		Expect(p.Update(event.UpdateEvent{ObjectOld: oldPrtl, ObjectNew: newPrtl})).To(BeFalse())
	})

	It("lets the synced check alone pass no create, delete or generic event", func() {
		synced := predicate.PortalSyncedPredicate{}
		prtl := travelPortal("Travel", "f9326ba3-0da0-3d4f-a7d5-f4a729b9e3fc")
		Expect(synced.Create(event.CreateEvent{Object: prtl})).To(BeFalse())
		Expect(synced.Delete(event.DeleteEvent{Object: prtl})).To(BeFalse())
		Expect(synced.Generic(event.GenericEvent{Object: prtl})).To(BeFalse())
	})
})
