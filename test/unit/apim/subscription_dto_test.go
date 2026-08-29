// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apim_test

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newSub(name, statusID string, automationManaged bool) *v1alpha1.Subscription {
	sub := &v1alpha1.Subscription{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default"},
		Status:     v1alpha1.SubscriptionStatus{ID: statusID},
	}
	sub.Spec.Plan = "jwt"
	if automationManaged {
		k8s.AddAutomationAPIManagedCondition(sub)
	}
	return sub
}

func newAPI(managed bool) *v1alpha1.ApiV4Definition {
	api := &v1alpha1.ApiV4Definition{
		ObjectMeta: metav1.ObjectMeta{Name: "e2e-v4-jwt-api", Namespace: "default"},
	}
	api.Spec.V4BaseApi = &v4.V4BaseApi{ApiBase: &base.ApiBase{}}
	api.Status.ID = "d1912aa8-3051-3e80-bbd4-c0097781bb97"
	if managed {
		k8s.AddAutomationAPIManagedCondition(api)
	}
	return api
}

func newApp(managed bool) *v1alpha1.Application {
	app := &v1alpha1.Application{
		ObjectMeta: metav1.ObjectMeta{Name: "e2e-app-simple", Namespace: "default"},
	}
	app.Status.ID = "f0848adf-3c29-300f-93bf-b8906295b0d2"
	if managed {
		k8s.AddAutomationAPIManagedCondition(app)
	}
	return app
}

var _ = Describe("ToSubscriptionDTO", func() {
	DescribeTable("SubscriptionUsesUUID",
		func(sub *v1alpha1.Subscription, want bool) {
			Expect(model.SubscriptionUsesUUID(sub)).To(Equal(want))
		},
		Entry("create: no condition, empty status ID → HRID",
			newSub("e2e-sub-jwt-v4", "", false), false),
		Entry("legacy: no condition, status UUID → UUID",
			newSub("e2e-sub-jwt-v4", "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", false), true),
		Entry("reconciled: condition set, empty status ID → HRID",
			newSub("e2e-sub-jwt-v4", "", true), false),
		Entry("reconciled: condition set, status UUID → HRID",
			newSub("e2e-sub-jwt-v4", "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", true), false),
	)

	It("uses HRIDs on create before the Automation condition is persisted", func() {
		dto := model.ToSubscriptionDTO(newSub("e2e-sub-jwt-v4", "", false), newAPI(true), newApp(true))

		Expect(dto.ID).To(Equal("default-e2e-sub-jwt-v4"))
		Expect(dto.ApiID).To(Equal("default-e2e-v4-jwt-api"))
		Expect(dto.AppID).To(Equal("default-e2e-app-simple"))
		Expect(dto.PlanID).To(Equal("jwt"))
	})

	It("keeps API and app as HRIDs when only the subscription is pre-Automation", func() {
		dto := model.ToSubscriptionDTO(
			newSub("e2e-sub-jwt-v4", "sub-uuid-1", false),
			newAPI(true),
			newApp(true),
		)

		Expect(dto.ID).To(Equal("sub-uuid-1"))
		Expect(dto.ApiID).To(Equal("default-e2e-v4-jwt-api"))
		Expect(dto.AppID).To(Equal("default-e2e-app-simple"))
		Expect(dto.PlanID).To(Equal("jwt"))
	})

	It("uses UUIDs per field when API and app are also pre-Automation", func() {
		api := newAPI(false)
		plan := &v4.Plan{Plan: &base.Plan{}}
		plan.ID = "plan-uuid"
		plans := map[string]*v4.Plan{"jwt": plan}
		api.Spec.Plans = &plans

		dto := model.ToSubscriptionDTO(
			newSub("e2e-sub-jwt-v4", "sub-uuid-1", false),
			api,
			newApp(false),
		)

		Expect(dto.ID).To(Equal("sub-uuid-1"))
		Expect(dto.ApiID).To(Equal(api.Status.ID))
		Expect(dto.AppID).To(Equal("f0848adf-3c29-300f-93bf-b8906295b0d2"))
		Expect(dto.PlanID).To(Equal("plan-uuid"))
	})
})
