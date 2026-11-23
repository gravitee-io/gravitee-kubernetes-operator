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
	"context"
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimauthzpdp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/authzpdp"
	xErrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

func newAuthzPdp() *v1alpha1.AuthzPdp {
	return &v1alpha1.AuthzPdp{
		ObjectMeta: metav1.ObjectMeta{Name: "orders-eu", Namespace: "gravitee"},
		Spec: v1alpha1.AuthzPdpSpec{
			Context: &refs.NamespacedName{Name: "dev-ctx"},
			Type: authzpdp.Type{
				Name:        "Orders PDP (EU)",
				TargetPdpID: "orders",
				Tag:         "eu",
				Strict:      true,
				Endpoint:    &authzpdp.Endpoint{Path: "/authz/orders-eu"},
			},
		},
	}
}

func toAuthzPdpDTO(pdp *v1alpha1.AuthzPdp) sdk.AuthzPdp {
	GinkgoHelper()
	dto, err := apimauthzpdp.ToDTO(pdp)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

// authzPdpState is what the platform answers for orders-eu: the payload, the ids, the findings, the routing
// scope and the discovery path.
const authzPdpState = `{
	"hrid": "gravitee-orders-eu",
	"name": "Orders PDP (EU)",
	"targetPdpId": "orders",
	"tag": "eu",
	"strict": true,
	"endpoint": {"path": "/authz/orders-eu"},
	"id": "5d0c1e7a", "organizationId": "DEFAULT", "environmentId": "DEFAULT",
	"routingScope": "orders@eu",
	"wellKnownPath": "/authz/orders-eu/.well-known/authzen-configuration",
	"errors": {"warning": ["strict is shared by every PDP of targetPdpId orders"]}
}`

var _ = Describe("Authz PDP", func() {
	Describe("spec hash", func() {
		DescribeTable("changes with each field",
			func(change func(pdp *v1alpha1.AuthzPdp)) {
				changed := newAuthzPdp()
				change(changed)
				Expect(changed.Spec.Hash()).ToNot(Equal(newAuthzPdp().Spec.Hash()))
			},
			Entry("name", func(pdp *v1alpha1.AuthzPdp) { pdp.Spec.Name = "Orders" }),
			Entry("targetPdpId", func(pdp *v1alpha1.AuthzPdp) { pdp.Spec.TargetPdpID = "billing" }),
			Entry("tag", func(pdp *v1alpha1.AuthzPdp) { pdp.Spec.Tag = "us" }),
			Entry("strict", func(pdp *v1alpha1.AuthzPdp) { pdp.Spec.Strict = false }),
			Entry("endpoint path", func(pdp *v1alpha1.AuthzPdp) { pdp.Spec.Endpoint.Path = "/authz/orders" }),
			Entry("no endpoint", func(pdp *v1alpha1.AuthzPdp) { pdp.Spec.Endpoint = nil }),
			Entry("context", func(pdp *v1alpha1.AuthzPdp) { pdp.Spec.Context.Name = "prod-ctx" }),
		)

		It("is stable for the same spec", func() {
			Expect(newAuthzPdp().Spec.Hash()).To(Equal(newAuthzPdp().Spec.Hash()))
		})
	})

	Describe("wire mapping", func() {
		It("derives the hrid from the namespace and name", func() {
			Expect(toAuthzPdpDTO(newAuthzPdp()).Hrid).To(Equal("gravitee-orders-eu"))
		})

		It("serialises the payload with the wire field names", func() {
			Expect(marshalJSON(toAuthzPdpDTO(newAuthzPdp()))).To(MatchJSON(`{
				"hrid": "gravitee-orders-eu",
				"name": "Orders PDP (EU)",
				"targetPdpId": "orders",
				"tag": "eu",
				"strict": true,
				"endpoint": {"path": "/authz/orders-eu"}
			}`))
		})

		It("omits an empty tag and a missing endpoint, and sends strict false", func() {
			pdp := newAuthzPdp()
			pdp.Spec.Tag = ""
			pdp.Spec.Strict = false
			pdp.Spec.Endpoint = nil

			Expect(marshalJSON(toAuthzPdpDTO(pdp))).To(MatchJSON(`{
				"hrid": "gravitee-orders-eu",
				"name": "Orders PDP (EU)",
				"targetPdpId": "orders",
				"strict": false
			}`))
		})
	})

	Describe("status", func() {
		It("copies the ids, the findings, the routing scope and the discovery path", func() {
			var state sdk.AuthzPdpState
			Expect(json.Unmarshal([]byte(authzPdpState), &state)).To(Succeed())
			pdp := newAuthzPdp()

			Expect(apimauthzpdp.UpdateStatus(context.Background(), pdp, apimauthzpdp.Response{State: state})).To(Succeed())

			Expect(pdp.Status.ID).To(Equal("5d0c1e7a"))
			Expect(pdp.Status.OrgID).To(Equal("DEFAULT"))
			Expect(pdp.Status.EnvID).To(Equal("DEFAULT"))
			Expect(pdp.Status.HRID).To(Equal("gravitee-orders-eu"))
			Expect(pdp.Status.RoutingScope).To(Equal("orders@eu"))
			Expect(pdp.Status.WellKnownPath).To(Equal("/authz/orders-eu/.well-known/authzen-configuration"))
			Expect(pdp.Status.Errors.Warning).To(ConsistOf("strict is shared by every PDP of targetPdpId orders"))
			Expect(k8s.IsAutomationAPIManaged(pdp)).To(BeTrue())
		})
	})

	Describe("calls", func() {
		const collection = "/automation/organizations/org1/environments/env1/authz/pdps"

		var (
			ctx    context.Context
			fake   *fakeAutomationAPI
			client *apim.SDKClient
		)

		BeforeEach(func() {
			ctx = context.Background()
			fake = newFakeAutomationAPI()
			fake.response = authzPdpState
			DeferCleanup(fake.Close)
			var err error
			client, err = apim.NewSDKClient(ctx, managementContext(fake.URL, bearer("token")))
			Expect(err).ToNot(HaveOccurred())
		})

		It("puts the payload on the pdps collection and answers the state", func() {
			resp, err := apimauthzpdp.Upsert(ctx, client, toAuthzPdpDTO(newAuthzPdp()))
			Expect(err).ToNot(HaveOccurred())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodPut))
			Expect(req.Path).To(Equal(collection))
			Expect(req.Query.Get("dryRun")).To(Equal("false"))
			Expect(marshalJSON(req.Body)).To(MatchJSON(`{
				"hrid": "gravitee-orders-eu",
				"name": "Orders PDP (EU)",
				"targetPdpId": "orders",
				"tag": "eu",
				"strict": true,
				"endpoint": {"path": "/authz/orders-eu"}
			}`))
			Expect(resp.State.Id).To(HaveValue(Equal("5d0c1e7a")))
			Expect(resp.State.RoutingScope).To(HaveValue(Equal("orders@eu")))
			Expect(resp.GetOrgID()).To(Equal("DEFAULT"))
			Expect(resp.GetEnvID()).To(Equal("DEFAULT"))
		})

		It("dry runs and returns the findings of the state", func() {
			fake.response = `{"hrid": "gravitee-orders-eu", "name": "Orders PDP (EU)", "targetPdpId": "orders",
				"errors": {"severe": ["routing scope orders@eu is held by gravitee-orders"], "warning": ["shared strict"]}}`

			findings, err := apimauthzpdp.DryRun(ctx, client, toAuthzPdpDTO(newAuthzPdp()))
			Expect(err).ToNot(HaveOccurred())

			Expect(fake.lastRequest().Query.Get("dryRun")).To(Equal("true"))
			Expect(findings.Severe).To(HaveLen(1))
			Expect(findings.Severe[0].Message).To(Equal("routing scope orders@eu is held by gravitee-orders"))
			Expect(findings.Warning).To(HaveLen(1))
			Expect(findings.Warning[0].Message).To(Equal("shared strict"))
		})

		It("returns the findings of a refused dry run", func() {
			fake.status = http.StatusBadRequest
			fake.errorResponse = `{"hrid": "gravitee-orders-eu",
				"errors": {"severe": ["endpoint path /authz/orders-eu is used by another API"]}}`

			findings, err := apimauthzpdp.DryRun(ctx, client, toAuthzPdpDTO(newAuthzPdp()))
			Expect(err).ToNot(HaveOccurred())

			Expect(findings.Severe).To(HaveLen(1))
			Expect(findings.Severe[0].Message).To(Equal("endpoint path /authz/orders-eu is used by another API"))
		})

		It("returns a 400 that carries no findings as the error", func() {
			fake.status = http.StatusBadRequest
			fake.errorResponse = `{"message": "hrid is required", "http_status": 400}`

			_, err := apimauthzpdp.DryRun(ctx, client, toAuthzPdpDTO(newAuthzPdp()))

			Expect(xErrors.IsBadRequest(err)).To(BeTrue())
		})

		It("returns the 403 of a strict change without the permission as the error", func() {
			fake.status = http.StatusForbidden
			fake.errorResponse = `{"message": "You do not have sufficient rights", "http_status": 403}`

			_, err := apimauthzpdp.DryRun(ctx, client, toAuthzPdpDTO(newAuthzPdp()))

			Expect(err).To(MatchError(ContainSubstring("You do not have sufficient rights")))
		})

		It("gets the remote as the payload a PUT sends", func() {
			remote, err := apimauthzpdp.GetRemote(ctx, client, toAuthzPdpDTO(newAuthzPdp()))
			Expect(err).ToNot(HaveOccurred())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodGet))
			Expect(req.Path).To(Equal(collection + "/gravitee-orders-eu"))
			Expect(remote).To(Equal(toAuthzPdpDTO(newAuthzPdp())))
		})

		It("reports a remote that is gone as not found", func() {
			fake.status = http.StatusNotFound

			_, err := apimauthzpdp.GetRemote(ctx, client, toAuthzPdpDTO(newAuthzPdp()))

			Expect(xErrors.IsNotFound(err)).To(BeTrue())
		})

		It("deletes by hrid without dryRun and treats a PDP already gone as deleted", func() {
			dto := toAuthzPdpDTO(newAuthzPdp())

			Expect(apimauthzpdp.Delete(ctx, client, dto)).To(Succeed())
			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodDelete))
			Expect(req.Path).To(Equal(collection + "/gravitee-orders-eu"))
			// APIM 4.13 ignores dryRun on a DELETE and deletes: it is never sent.
			Expect(req.Query).ToNot(HaveKey("dryRun"))

			fake.status = http.StatusNotFound
			Expect(apimauthzpdp.Delete(ctx, client, dto)).To(Succeed())
		})

		It("returns the 409 of a PDP still needed as a target", func() {
			fake.status = http.StatusConflict
			fake.errorResponse = `{"message": "policy hotel-read would have no target", "http_status": 409}`

			err := apimauthzpdp.Delete(ctx, client, toAuthzPdpDTO(newAuthzPdp()))

			Expect(err).To(MatchError(ContainSubstring("policy hotel-read would have no target")))
		})
	})

	Describe("a refused apply", func() {
		refusedWith := func(statusCode int, body string) error {
			return xErrors.ServerError{StatusCode: statusCode, Body: body}
		}

		It("records the findings in the status, keeps the routing scope and stays recoverable", func() {
			pdp := newAuthzPdp()
			pdp.Status.RoutingScope = "orders@eu"

			err := apimauthzpdp.OnUpsertError(context.Background(), pdp,
				refusedWith(http.StatusBadRequest, `{"errors": {"severe": ["endpoint path taken"]}}`))

			Expect(err).To(MatchError(ContainSubstring("authz pdp [orders-eu] was refused by the automation api")))
			Expect(xErrors.IsRecoverable(err)).To(BeTrue())
			Expect(pdp.Status.Errors.Severe).To(ConsistOf("endpoint path taken"))
			Expect(pdp.Status.RoutingScope).To(Equal("orders@eu"))
		})

		It("stays recoverable when the 400 carries no findings", func() {
			pdp := newAuthzPdp()

			err := apimauthzpdp.OnUpsertError(context.Background(), pdp,
				refusedWith(http.StatusBadRequest, `{"message": "hrid is required", "http_status": 400}`))

			Expect(err).To(MatchError(ContainSubstring("authz pdp [orders-eu] was refused by the automation api")))
			Expect(xErrors.IsRecoverable(err)).To(BeTrue())
			Expect(pdp.Status.Errors.Severe).To(BeEmpty())
		})

		It("leaves any other failure to the lifecycle", func() {
			err := apimauthzpdp.OnUpsertError(context.Background(), newAuthzPdp(),
				refusedWith(http.StatusInternalServerError, `{"errors": {"severe": ["boom"]}}`))

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
