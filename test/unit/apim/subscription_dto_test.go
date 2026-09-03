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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
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

func newAPIWithPlan(managed bool, planID string) *v1alpha1.ApiV4Definition {
	api := newAPI(managed)
	plan := &v4.Plan{Plan: &base.Plan{}}
	plan.ID = planID
	plans := map[string]*v4.Plan{"jwt": plan}
	api.Spec.Plans = &plans
	return api
}

func newV2API() *v1alpha1.ApiDefinition {
	api := &v1alpha1.ApiDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: "e2e-v2-jwt-api", Namespace: "default"},
	}
	api.Status.ID = "v2-api-uuid"
	plan := &v2.Plan{Plan: &base.Plan{ID: "v2-plan-uuid"}, Name: "jwt"}
	api.Spec.Plans = []*v2.Plan{plan}
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

	DescribeTable("APIUsesUUID",
		func(api core.ApiDefinitionObject, want bool) {
			Expect(model.APIUsesUUID(api)).To(Equal(want))
		},
		Entry("V4 create: no condition, empty status ID → HRID", emptyV4API(), false),
		Entry("V4 managed → HRID", newAPI(true), false),
		Entry("V4 pre-Automation with status ID → UUID", newAPI(false), true),
		Entry("V2 pre-Automation with status ID → UUID", newV2API(), true),
	)

	DescribeTable("ApplicationUsesUUID",
		func(app *v1alpha1.Application, want bool) {
			Expect(model.ApplicationUsesUUID(app)).To(Equal(want))
		},
		Entry("create: no condition, empty status ID → HRID", emptyApp(), false),
		Entry("managed → HRID", newApp(true), false),
		Entry("pre-Automation with status ID → UUID", newApp(false), true),
	)

	DescribeTable("SubscriptionRemoteUsesLegacyAPI",
		func(sub *v1alpha1.Subscription, api core.ApiDefinitionObject, want bool) {
			Expect(model.SubscriptionRemoteUsesLegacyAPI(sub, api)).To(Equal(want))
		},
		Entry("HRID sub + V4 managed API → GetByHRID",
			newSub("e2e-sub-jwt-v4", "", false), newAPI(true), false),
		Entry("HRID sub + V2 UUID API → GetByHRIDWithAPIUUID",
			newSub("e2e-sub-jwt-v2", "", false), newV2API(), true),
		Entry("legacy UUID sub + V2 API → GetWithUUID",
			newSub("e2e-sub-jwt-v2", "sub-uuid-1", false), newV2API(), false),
	)

	It("uses HRIDs on create before the Automation condition is persisted", func() {
		dto := model.ToSubscriptionDTO(newSub("e2e-sub-jwt-v4", "", false), newAPI(true), newApp(true))

		Expect(dto.ID).To(Equal("default-e2e-sub-jwt-v4"))
		Expect(dto.ApiID).To(Equal("default-e2e-v4-jwt-api"))
		Expect(dto.AppID).To(Equal("default-e2e-app-simple"))
		Expect(dto.PlanID).To(Equal("jwt"))
	})

	It("does not copy StartingAt onto the Import payload", func() {
		sub := newSub("e2e-sub-jwt-v4", "", true)
		sub.Status.StartedAt = "2023-07-24T20:43:16Z"

		dto := model.ToSubscriptionDTO(sub, newAPI(true), newApp(true))

		Expect(dto.StartingAt).To(BeEmpty())
	})

	It("maps spec fields used by Import and drift", func() {
		ending := "2024-08-25T23:43:16Z"
		sub := newSub("e2e-sub-jwt-v4", "", true)
		sub.Spec.EndingAt = &ending
		sub.Spec.Metadata = map[string]string{"team": "platform"}
		sub.Spec.ApiKeys = []subscription.ApiKeySpec{{Key: "key1", ExpireAt: &ending}}
		sub.Spec.ConsumerConfiguration = &subscription.ConsumerConfiguration{EntrypointID: "webhook"}

		dto := model.ToSubscriptionDTO(sub, newAPI(true), newApp(true))

		Expect(dto.EndingAt).To(Equal(ending))
		Expect(dto.Metadata).To(Equal(map[string]string{"team": "platform"}))
		Expect(dto.ApiKeys).To(HaveLen(1))
		Expect(dto.ApiKeys[0].Key).To(Equal("key1"))
		Expect(dto.ConsumerConfiguration).NotTo(BeNil())
		Expect(dto.ConsumerConfiguration.EntrypointID).To(Equal("webhook"))
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
		dto := model.ToSubscriptionDTO(
			newSub("e2e-sub-jwt-v4", "sub-uuid-1", false),
			newAPIWithPlan(false, "plan-uuid"),
			newApp(false),
		)

		Expect(dto.ID).To(Equal("sub-uuid-1"))
		Expect(dto.ApiID).To(Equal("d1912aa8-3051-3e80-bbd4-c0097781bb97"))
		Expect(dto.AppID).To(Equal("f0848adf-3c29-300f-93bf-b8906295b0d2"))
		Expect(dto.PlanID).To(Equal("plan-uuid"))
	})

	It("uses API UUID and plan UUID for a new HRID subscription on a V2 API", func() {
		dto := model.ToSubscriptionDTO(newSub("e2e-sub-jwt-v2", "", false), newV2API(), newApp(true))

		Expect(dto.ID).To(Equal("default-e2e-sub-jwt-v2"))
		Expect(dto.ApiID).To(Equal("v2-api-uuid"))
		Expect(dto.AppID).To(Equal("default-e2e-app-simple"))
		Expect(dto.PlanID).To(Equal("v2-plan-uuid"))
	})

	It("falls back to the plan HRID when a UUID API has no matching plan", func() {
		dto := model.ToSubscriptionDTO(newSub("e2e-sub-jwt-v4", "", false), newAPI(false), newApp(true))

		Expect(dto.PlanID).To(Equal("jwt"))
	})
})

var _ = Describe("ToSubscriptionDTOForDrift", func() {
	It("uses all UUIDs when the remote GET is Management v2", func() {
		dto := model.ToSubscriptionDTOForDrift(
			newSub("e2e-sub-jwt-v4", "sub-uuid-1", false),
			newAPIWithPlan(true, "plan-uuid"),
			newApp(true),
		)

		Expect(dto.ID).To(Equal("sub-uuid-1"))
		Expect(dto.ApiID).To(Equal("d1912aa8-3051-3e80-bbd4-c0097781bb97"))
		Expect(dto.AppID).To(Equal("f0848adf-3c29-300f-93bf-b8906295b0d2"))
		Expect(dto.PlanID).To(Equal("plan-uuid"))
	})

	It("keeps the mixed Import shape when the remote GET is Automation", func() {
		write := model.ToSubscriptionDTO(newSub("e2e-sub-jwt-v2", "", false), newV2API(), newApp(true))
		drift := model.ToSubscriptionDTOForDrift(newSub("e2e-sub-jwt-v2", "", false), newV2API(), newApp(true))

		Expect(drift).To(Equal(write))
		Expect(drift.ApiID).To(Equal("v2-api-uuid"))
		Expect(drift.ID).To(Equal("default-e2e-sub-jwt-v2"))
	})
})

func emptyV4API() *v1alpha1.ApiV4Definition {
	api := &v1alpha1.ApiV4Definition{
		ObjectMeta: metav1.ObjectMeta{Name: "e2e-v4-jwt-api", Namespace: "default"},
	}
	api.Spec.V4BaseApi = &v4.V4BaseApi{ApiBase: &base.ApiBase{}}
	return api
}

func emptyApp() *v1alpha1.Application {
	return &v1alpha1.Application{
		ObjectMeta: metav1.ObjectMeta{Name: "e2e-app-simple", Namespace: "default"},
	}
}
