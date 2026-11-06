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

package apim_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpolicy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
)

const policyOwnershipFinding = `entityId [policy.hotel-booking-read] is held by policy "Hotel booking" ` +
	`(id: 5f1c), which this manifest did not create. Declare another entityId, or delete that policy first`

const hotelBookingReadText = "permit(principal in Group::\"agents\", action == Action::\"read\", " +
	"resource == Hotel::\"grand-paris\");\n"

func newAuthzPolicy() *v1alpha1.AuthzPolicy {
	return &v1alpha1.AuthzPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "hotel-booking-read", Namespace: "gravitee"},
		Spec: v1alpha1.AuthzPolicySpec{
			Context: &refs.NamespacedName{Name: "dev-ctx"},
			Type: authzpolicy.Type{
				EntityID:   "policy.hotel-booking-read",
				Name:       "Hotel booking — read access",
				PolicyType: "API",
				PolicyText: hotelBookingReadText,
				Targets: []authzpdp.Target{
					{PdpRef: &refs.NamespacedName{Name: "orders-eu"}},
					{Scope: "default"},
				},
			},
		},
	}
}

var _ = Describe("Authz policy", func() {
	Describe("spec hash", func() {
		DescribeTable("changes with each field",
			func(change func(policy *v1alpha1.AuthzPolicy)) {
				changed := newAuthzPolicy()
				change(changed)
				Expect(changed.Spec.Hash()).ToNot(Equal(newAuthzPolicy().Spec.Hash()))
			},
			Entry("entityId", func(policy *v1alpha1.AuthzPolicy) { policy.Spec.EntityID = "policy.hotel-booking-write" }),
			Entry("name", func(policy *v1alpha1.AuthzPolicy) { policy.Spec.Name = "Hotel booking" }),
			Entry("policyType", func(policy *v1alpha1.AuthzPolicy) { policy.Spec.PolicyType = "MCP" }),
			Entry("policyText", func(policy *v1alpha1.AuthzPolicy) { policy.Spec.PolicyText = "forbid(principal, action, resource);" }),
			Entry("pdpRef", func(policy *v1alpha1.AuthzPolicy) { policy.Spec.Targets[0].PdpRef.Name = "orders-us" }),
			Entry("scope instead of pdpRef", func(policy *v1alpha1.AuthzPolicy) {
				policy.Spec.Targets[0] = authzpdp.Target{Scope: "orders@eu"}
			}),
			Entry("no targets", func(policy *v1alpha1.AuthzPolicy) { policy.Spec.Targets = nil }),
			Entry("context", func(policy *v1alpha1.AuthzPolicy) { policy.Spec.Context.Name = "prod-ctx" }),
		)

		It("is stable for the same spec", func() {
			Expect(newAuthzPolicy().Spec.Hash()).To(Equal(newAuthzPolicy().Spec.Hash()))
		})
	})

	Describe("wire mapping", func() {
		It("derives the hrid from the namespace and name", func() {
			Expect(model.ToAuthzPolicyDTO(newAuthzPolicy(), nil).HRID).To(Equal("gravitee-hotel-booking-read"))
		})

		It("sends the targets the caller resolved, in their order, not the pdpRefs", func() {
			raw, err := json.Marshal(model.ToAuthzPolicyDTO(newAuthzPolicy(), []string{"orders@eu", "default"}))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(raw)).To(MatchJSON(`{
				"hrid": "gravitee-hotel-booking-read",
				"entityId": "policy.hotel-booking-read",
				"name": "Hotel booking — read access",
				"policyType": "API",
				"policyText": ` + mustJSON(hotelBookingReadText) + `,
				"targets": ["orders@eu", "default"]
			}`))
		})

		It("omits absent targets: the platform deploys to the default engine on every gateway", func() {
			raw, err := json.Marshal(model.ToAuthzPolicyDTO(newAuthzPolicy(), nil))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(raw)).ToNot(ContainSubstring(`"targets"`))
		})

		It("sends the text verbatim, its trailing newline included", func() {
			Expect(model.ToAuthzPolicyDTO(newAuthzPolicy(), nil).PolicyText).To(Equal(hotelBookingReadText))
		})

		It("never sends policyStatus, which the platform owns", func() {
			raw, err := json.Marshal(model.ToAuthzPolicyDTO(newAuthzPolicy(), nil))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(raw)).ToNot(ContainSubstring(`policyStatus`))
		})
	})

	Describe("state decoding", func() {
		It("decodes the payload, the ids, the findings and the deployment", func() {
			body := `{
				"hrid": "gravitee-hotel-booking-read",
				"entityId": "policy.hotel-booking-read",
				"name": "Hotel booking — read access",
				"policyType": "API",
				"policyText": ` + mustJSON(hotelBookingReadText) + `,
				"targets": ["default", "orders@eu"],
				"id": "5f1c", "organizationId": "DEFAULT", "environmentId": "DEFAULT",
				"errors": {"warning": ["the policy is DISABLED: the apply deploys it again, and the PDP gateways it targets enforce it"]},
				"policyStatus": "DEPLOYED"
			}`
			state := new(model.AuthzPolicyState)
			Expect(json.Unmarshal([]byte(body), state)).To(Succeed())

			Expect(state.HRID).To(Equal("gravitee-hotel-booking-read"))
			Expect(state.EntityID).To(Equal("policy.hotel-booking-read"))
			Expect(state.Name).To(Equal("Hotel booking — read access"))
			Expect(state.PolicyType).To(Equal("API"))
			Expect(state.PolicyText).To(Equal(hotelBookingReadText))
			Expect(state.Targets).To(Equal([]string{"default", "orders@eu"}))
			Expect(state.ID).To(Equal("5f1c"))
			Expect(state.OrgID).To(Equal("DEFAULT"))
			Expect(state.EnvID).To(Equal("DEFAULT"))
			Expect(state.Errors.Warning).To(HaveLen(1))
			Expect(state.PolicyStatus).To(Equal("DEPLOYED"))
		})
	})
})
