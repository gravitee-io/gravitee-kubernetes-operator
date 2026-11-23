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
	"context"
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpolicy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimauthzpolicy "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/authzpolicy"
	xErrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
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

// resolvedAuthzPolicy is newAuthzPolicy with its pdpRef resolved, as dynamic.ResolveAuthzPolicyRefs leaves it.
func resolvedAuthzPolicy() *v1alpha1.AuthzPolicy {
	policy := newAuthzPolicy()
	policy.Spec.Targets = []authzpdp.Target{{Scope: "orders@eu"}, {Scope: "default"}}
	return policy
}

func toAuthzPolicyDTO(policy *v1alpha1.AuthzPolicy) sdk.AuthzPolicy {
	GinkgoHelper()
	dto, err := apimauthzpolicy.ToDTO(policy)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

// authzPolicyState is what the platform answers for hotel-booking-read: the payload, its targets sorted, the ids,
// the findings and the deployment.
const authzPolicyState = `{
	"hrid": "gravitee-hotel-booking-read",
	"entityId": "policy.hotel-booking-read",
	"name": "Hotel booking — read access",
	"policyType": "API",
	"policyText": "permit(principal in Group::\"agents\", action == Action::\"read\", resource == Hotel::\"grand-paris\");\n",
	"targets": ["default", "orders@eu"],
	"id": "5f1c", "organizationId": "DEFAULT", "environmentId": "DEFAULT",
	"errors": {"warning": ["the policy is DISABLED: the apply deploys it again, and the PDP gateways it targets enforce it"]},
	"policyStatus": "DEPLOYED"
}`

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
			Expect(toAuthzPolicyDTO(resolvedAuthzPolicy()).Hrid).To(Equal("gravitee-hotel-booking-read"))
		})

		It("sends the scope of each resolved target, in declaration order", func() {
			Expect(marshalJSON(toAuthzPolicyDTO(resolvedAuthzPolicy()))).To(MatchJSON(`{
				"hrid": "gravitee-hotel-booking-read",
				"entityId": "policy.hotel-booking-read",
				"name": "Hotel booking — read access",
				"policyType": "API",
				"policyText": ` + mustJSON(hotelBookingReadText) + `,
				"targets": ["orders@eu", "default"]
			}`))
		})

		It("omits absent targets: the platform deploys to the default engine on every gateway", func() {
			untargeted := newAuthzPolicy()
			untargeted.Spec.Targets = nil

			Expect(marshalJSON(toAuthzPolicyDTO(untargeted))).ToNot(ContainSubstring(`"targets"`))
		})

		It("sends the text verbatim, its trailing newline included", func() {
			Expect(toAuthzPolicyDTO(resolvedAuthzPolicy()).PolicyText).To(Equal(hotelBookingReadText))
		})

		It("never sends policyStatus, which the platform owns", func() {
			Expect(marshalJSON(toAuthzPolicyDTO(resolvedAuthzPolicy()))).ToNot(ContainSubstring(`policyStatus`))
		})

		// A pdpRef has no scope to send until its AuthzPdp's routing scope replaces it.
		It("refuses a pdpRef left unresolved, naming its index", func() {
			_, err := apimauthzpolicy.ToDTO(newAuthzPolicy())

			Expect(err).To(MatchError("targets[0].pdpRef [gravitee/orders-eu] is not resolved"))
		})
	})

	Describe("status", func() {
		It("copies the ids, the findings and the deployment", func() {
			var state sdk.AuthzPolicyState
			Expect(json.Unmarshal([]byte(authzPolicyState), &state)).To(Succeed())
			policy := newAuthzPolicy()

			Expect(apimauthzpolicy.UpdateStatus(context.Background(), policy,
				apimauthzpolicy.Response{State: state})).To(Succeed())

			Expect(policy.Status.ID).To(Equal("5f1c"))
			Expect(policy.Status.OrgID).To(Equal("DEFAULT"))
			Expect(policy.Status.EnvID).To(Equal("DEFAULT"))
			Expect(policy.Status.HRID).To(Equal("gravitee-hotel-booking-read"))
			Expect(policy.Status.PolicyStatus).To(Equal("DEPLOYED"))
			Expect(policy.Status.Errors.Warning).To(HaveLen(1))
			Expect(k8s.IsAutomationAPIManaged(policy)).To(BeTrue())
		})
	})

	Describe("calls", func() {
		const collection = "/automation/organizations/org1/environments/env1/authz/policies"

		var (
			ctx    context.Context
			fake   *fakeAutomationAPI
			client *apim.SDKClient
		)

		BeforeEach(func() {
			ctx = context.Background()
			fake = newFakeAutomationAPI()
			fake.response = authzPolicyState
			DeferCleanup(fake.Close)
			var err error
			client, err = apim.NewSDKClient(ctx, managementContext(fake.URL, bearer("token")))
			Expect(err).ToNot(HaveOccurred())
		})

		It("puts the payload on the policies collection and answers the state", func() {
			resp, err := apimauthzpolicy.Upsert(ctx, client, toAuthzPolicyDTO(resolvedAuthzPolicy()))
			Expect(err).ToNot(HaveOccurred())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodPut))
			Expect(req.Path).To(Equal(collection))
			Expect(req.Query.Get("dryRun")).To(Equal("false"))
			Expect(marshalJSON(req.Body)).To(MatchJSON(marshalJSON(toAuthzPolicyDTO(resolvedAuthzPolicy()))))
			Expect(resp.State.Id).To(HaveValue(Equal("5f1c")))
			Expect(resp.GetOrgID()).To(Equal("DEFAULT"))
			Expect(resp.GetEnvID()).To(Equal("DEFAULT"))
		})

		It("dry runs and returns the findings of the state", func() {
			fake.response = `{"hrid": "gravitee-hotel-booking-read", "entityId": "policy.hotel-booking-read",
				"name": "", "policyType": "API", "policyText": "",
				"errors": {"severe": [` + mustJSON(policyOwnershipFinding) + `]}}`

			findings, err := apimauthzpolicy.DryRun(ctx, client, toAuthzPolicyDTO(resolvedAuthzPolicy()))
			Expect(err).ToNot(HaveOccurred())

			Expect(fake.lastRequest().Query.Get("dryRun")).To(Equal("true"))
			Expect(findings.Severe).To(HaveLen(1))
			Expect(findings.Severe[0].Message).To(Equal(policyOwnershipFinding))
		})

		It("returns the findings of a refused dry run", func() {
			fake.status = http.StatusBadRequest
			fake.errorResponse = `{"hrid": "gravitee-hotel-booking-read",
				"errors": {"severe": ["policy text does not parse: unexpected token at 1:7"]}}`

			findings, err := apimauthzpolicy.DryRun(ctx, client, toAuthzPolicyDTO(resolvedAuthzPolicy()))
			Expect(err).ToNot(HaveOccurred())

			Expect(findings.Severe).To(HaveLen(1))
			Expect(findings.Severe[0].Message).To(Equal("policy text does not parse: unexpected token at 1:7"))
		})

		It("returns a 400 that carries no findings as the error", func() {
			fake.status = http.StatusBadRequest
			fake.errorResponse = `{"message": "hrid is required", "http_status": 400}`

			_, err := apimauthzpolicy.DryRun(ctx, client, toAuthzPolicyDTO(resolvedAuthzPolicy()))

			Expect(xErrors.IsBadRequest(err)).To(BeTrue())
		})

		It("gets the remote as the payload a PUT sends, without its deployment", func() {
			remote, err := apimauthzpolicy.GetRemote(ctx, client, toAuthzPolicyDTO(resolvedAuthzPolicy()))
			Expect(err).ToNot(HaveOccurred())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodGet))
			Expect(req.Path).To(Equal(collection + "/gravitee-hotel-booking-read"))
			Expect(marshalJSON(remote)).To(MatchJSON(`{
				"hrid": "gravitee-hotel-booking-read",
				"entityId": "policy.hotel-booking-read",
				"name": "Hotel booking — read access",
				"policyType": "API",
				"policyText": ` + mustJSON(hotelBookingReadText) + `,
				"targets": ["default", "orders@eu"]
			}`))
		})

		It("reports a remote that is gone as not found", func() {
			fake.status = http.StatusNotFound

			_, err := apimauthzpolicy.GetRemote(ctx, client, toAuthzPolicyDTO(resolvedAuthzPolicy()))

			Expect(xErrors.IsNotFound(err)).To(BeTrue())
		})

		It("deletes by hrid and treats a policy already gone as deleted", func() {
			dto := toAuthzPolicyDTO(resolvedAuthzPolicy())

			Expect(apimauthzpolicy.Delete(ctx, client, dto)).To(Succeed())
			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodDelete))
			Expect(req.Path).To(Equal(collection + "/gravitee-hotel-booking-read"))
			Expect(req.Query.Has("dryRun")).To(BeFalse())

			fake.status = http.StatusNotFound
			Expect(apimauthzpolicy.Delete(ctx, client, dto)).To(Succeed())
		})

		// The lifecycle retries a failed delete, keeping the finalizer until the token may withdraw the policy.
		It("returns the 403 of a token that may not withdraw a deployed policy", func() {
			fake.status = http.StatusForbidden
			fake.errorResponse = `{"message": "ENVIRONMENT_AUTHZ_PDP[UPDATE] is required", "http_status": 403}`

			err := apimauthzpolicy.Delete(ctx, client, toAuthzPolicyDTO(resolvedAuthzPolicy()))

			Expect(err).To(MatchError(ContainSubstring("ENVIRONMENT_AUTHZ_PDP[UPDATE] is required")))
		})
	})

	Describe("a refused apply", func() {
		refusedWith := func(statusCode int, body string) error {
			return xErrors.ServerError{StatusCode: statusCode, Body: body}
		}

		It("records the findings in the status, keeps the deployment and stays recoverable", func() {
			policy := newAuthzPolicy()
			policy.Status.PolicyStatus = "DEPLOYED"

			err := apimauthzpolicy.OnUpsertError(context.Background(), policy, refusedWith(http.StatusBadRequest,
				`{"errors": {"severe": [`+mustJSON(policyOwnershipFinding)+`]}}`))

			Expect(err).To(MatchError(ContainSubstring("authz policy [hotel-booking-read] was refused by the automation api")))
			Expect(xErrors.IsRecoverable(err)).To(BeTrue())
			Expect(policy.Status.Errors.Severe).To(ConsistOf(policyOwnershipFinding))
			Expect(policy.Status.PolicyStatus).To(Equal("DEPLOYED"))
		})

		It("stays recoverable when the 400 carries no findings", func() {
			policy := newAuthzPolicy()

			err := apimauthzpolicy.OnUpsertError(context.Background(), policy,
				refusedWith(http.StatusBadRequest, `{"message": "hrid is required", "http_status": 400}`))

			Expect(xErrors.IsRecoverable(err)).To(BeTrue())
			Expect(policy.Status.Errors.Severe).To(BeEmpty())
		})

		It("leaves any other failure to the lifecycle", func() {
			err := apimauthzpolicy.OnUpsertError(context.Background(), newAuthzPolicy(),
				refusedWith(http.StatusInternalServerError, `{"errors": {"severe": ["boom"]}}`))

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
