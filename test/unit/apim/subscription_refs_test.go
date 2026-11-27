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

package apim_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	apimsubscription "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/subscription"
)

// The dynamic resolvers write the namespace into the reference they get. Admission compares the references of
// the previous and the new Subscription, so resolving must never change the Subscription's own.
var _ = Describe("Subscription references", func() {
	It("resolves copies in the subscription's namespace, leaving references without a namespace as written", func() {
		sub := newSub("sub", "", true)
		sub.Spec.API = refs.NamespacedName{Name: "api"}
		sub.Spec.App = refs.NamespacedName{Name: "app"}

		apiRef := apimsubscription.APIRef(sub)
		appRef := apimsubscription.AppRef(sub)
		apiRef.Name, appRef.Name = "changed", "changed"

		Expect(apiRef.Namespace).To(Equal("default"))
		Expect(appRef.Namespace).To(Equal("default"))
		Expect(sub.Spec.API).To(Equal(refs.NamespacedName{Name: "api"}))
		Expect(sub.Spec.App).To(Equal(refs.NamespacedName{Name: "app"}))
	})

	It("keeps the namespace a reference names", func() {
		sub := newSub("sub", "", true)
		sub.Spec.API = refs.NamespacedName{Name: "api", Namespace: "apis"}

		Expect(apimsubscription.APIRef(sub).Namespace).To(Equal("apis"))
	})
})


var _ = Describe("Subscription RecordRefs", func() {
	It("addresses an API, its plan and an application created before the Automation API by their UUIDs", func() {
		sub := newSub("sub", "", false)
		api, app := newAPIWithPlan(false, "plan-uuid"), newApp(false)

		apimsubscription.RecordRefs(sub, api, app)
		dto, err := apimsubscription.ToDTO(sub)

		Expect(err).ToNot(HaveOccurred())
		Expect(dto.Hrid).To(Equal("default-sub"))
		Expect(dto.ApiHrid).To(Equal(api.GetID()))
		Expect(dto.PlanHrid).To(Equal("plan-uuid"))
		Expect(dto.ApplicationHrid).To(Equal(app.GetID()))
		Expect(dto.Legacy).To(Equal(apimsubscription.Legacy{HridContainsApiUUID: true, HridContainsAppUUID: true}))
	})

	It("leaves the application out when it is gone", func() {
		sub := newSub("sub", "", false)

		apimsubscription.RecordRefs(sub, newAPIWithPlan(true, ""), nil)
		dto, err := apimsubscription.ToDTO(sub)

		Expect(err).ToNot(HaveOccurred())
		Expect(dto.ApplicationHrid).To(BeEmpty())
		Expect(dto.Legacy.HridContainsAppUUID).To(BeFalse())
	})
})
