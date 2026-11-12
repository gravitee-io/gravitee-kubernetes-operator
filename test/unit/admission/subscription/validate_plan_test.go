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

package subscription

import (
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	submodel "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/subscription"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/admission/subscription"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const pushPlanKey = "Default-PUSH-plan"

func pushPlan() *v4.Plan {
	plan := v4.NewPlan()
	plan.Name = "Default PUSH plan"
	plan.Mode = v4.PlanModePush
	return plan
}

func standardPlan(securityType string) *v4.Plan {
	plan := v4.NewPlan()
	plan.Name = "Standard plan"
	if securityType != "" {
		plan.WithSecurity(&v4.PlanSecurity{Type: securityType})
	}
	return plan
}

func apiWithPlan(key string, plan *v4.Plan) *v1alpha1.ApiV4Definition {
	plans := map[string]*v4.Plan{key: plan}
	return &v1alpha1.ApiV4Definition{
		Spec: v1alpha1.ApiV4DefinitionSpec{Api: v4.Api{Plans: &plans}},
	}
}

func subscriptionTo(planKey string, consumer *submodel.ConsumerConfiguration) *v1alpha1.Subscription {
	return &v1alpha1.Subscription{
		Spec: v1alpha1.SubscriptionSpec{Type: submodel.Type{
			API:                   refs.NamespacedName{Namespace: "default", Name: "webhook-api"},
			App:                   refs.NamespacedName{Namespace: "default", Name: "webhook-application"},
			Plan:                  planKey,
			ConsumerConfiguration: consumer,
		}},
	}
}

func webhookConsumer(entrypointID string) *submodel.ConsumerConfiguration {
	return &submodel.ConsumerConfiguration{EntrypointID: entrypointID}
}

// PUSH plans have no security block. Admission must accept them on their consumer
// configuration instead of reading a security type (GKO-3237).
var _ = Describe("Validating the subscribed plan", func() {
	DescribeTable("subscribing to a v4 plan",
		func(plan *v4.Plan, subscribedKey string, consumer *submodel.ConsumerConfiguration, expected string) {
			api := apiWithPlan(pushPlanKey, plan)
			sub := subscriptionTo(subscribedKey, consumer)

			err := subscription.ValidatePlan(sub, api)

			if expected == "" {
				Expect(err).To(BeNil())
				return
			}
			Expect(err).ToNot(BeNil())
			Expect(err.Message).To(ContainSubstring(expected))
		},
		Entry("is admitted on a PUSH plan with a webhook entrypoint",
			pushPlan(), pushPlanKey, webhookConsumer("webhook"), ""),
		Entry("is rejected on a PUSH plan without consumerConfiguration",
			pushPlan(), pushPlanKey, nil, "consumerConfiguration.entrypointId is required"),
		Entry("is rejected on a PUSH plan with an empty entrypointId",
			pushPlan(), pushPlanKey, webhookConsumer(""), "consumerConfiguration.entrypointId is required"),
		Entry("is rejected on a STANDARD plan without security",
			standardPlan(""), pushPlanKey, nil, "security type is not one of"),
		Entry("is rejected on a KEY_LESS plan",
			standardPlan("KEY_LESS"), pushPlanKey, nil, "security type is not one of"),
		Entry("is admitted on an API_KEY plan",
			standardPlan("API_KEY"), pushPlanKey, nil, ""),
		Entry("is rejected when the plan does not exist",
			pushPlan(), "unknown-plan", webhookConsumer("webhook"), "cannot be found"),
	)

	It("reads an empty security type from a plan without security", func() {
		Expect(pushPlan().GetSecurityType()).To(BeEmpty())
	})

	It("reads an empty security type and no PUSH mode from a nil plan", func() {
		var plan *v4.Plan
		Expect(plan.GetSecurityType()).To(BeEmpty())
		Expect(plan.IsPush()).To(BeFalse())
	})
})
