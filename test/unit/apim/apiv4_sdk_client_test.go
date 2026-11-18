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
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/apiv4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

const apiV4State = `{
	"id": "6b1d5f0a-0000-0000-0000-000000000001",
	"crossId": "6b1d5f0a-0000-0000-0000-000000000002",
	"hrid": "orders",
	"name": "orders",
	"version": "1",
	"type": "PROXY",
	"organizationId": "org1",
	"environmentId": "env1",
	"groups": ["remote-group"],
	"plans": [{"hrid": "KeyLess", "id": "6b1d5f0a-0000-0000-0000-000000000003", "name": "KeyLess",
		"security": {"type": "KEY_LESS"}, "status": "PUBLISHED", "mode": "STANDARD", "type": "API", "validation": "AUTO"}],
	"errors": {"severe": [], "warning": ["warning one"]}
}`

const apisPath = "/automation/organizations/org1/environments/env1/apis"

// apiV4ForSDK is an API with a management context that syncs from management (no gateway ConfigMap), with an HTTP
// flow selector that leaves path and pathOperator to APIM.
func apiV4ForSDK() *v1alpha1.ApiV4Definition {
	GinkgoHelper()
	api := apiV4WithContext()
	api.Spec.DefinitionContext = nil
	api.Spec.Groups = []string{"crd-group"}
	withSelector := apiV4(`{"flows": [{"name": "all", "enabled": true, "selectors": [{"type": "HTTP"}]}]}`)
	api.Spec.Flows = withSelector.Spec.Flows
	return api
}

var _ = Describe("ApiV4Definition through the APIM SDK", func() {
	var (
		ctx    context.Context
		fake   *fakeAutomationAPI
		client *apiv4.Client
	)

	BeforeEach(func() {
		ctx = context.Background()
		fake = newFakeAutomationAPI()
		fake.response = apiV4State
		DeferCleanup(fake.Close)
		registerFakeK8s()

		mctx := managementContext(fake.URL, bearer("token"))
		sdkClient, err := apim.NewSDKClient(ctx, mctx)
		Expect(err).ToNot(HaveOccurred())
		facade, err := apim.FromContext(ctx, mctx, "default")
		Expect(err).ToNot(HaveOccurred())
		client = &apiv4.Client{SDKClient: sdkClient, APIM: facade}
	})

	legacy := func() *v1alpha1.ApiV4Definition {
		api := apiV4ForSDK()
		api.Spec.ID = legacyAPIID
		return api
	}

	withoutHrid := func(body map[string]any) map[string]any {
		cp := make(map[string]any, len(body))
		for key, value := range body {
			if key != "hrid" {
				cp[key] = value
			}
		}
		return cp
	}

	It("upserts an API by hrid through the SDK", func() {
		resp, err := apiv4.Upsert(ctx, client, toAPIV4DTO(apiV4ForSDK()))
		Expect(err).ToNot(HaveOccurred())

		req := fake.lastRequest()
		Expect(req.Method).To(Equal(http.MethodPut))
		Expect(req.Path).To(Equal(apisPath))
		Expect(req.Query.Get("dryRun")).To(Equal("false"))
		Expect(req.Query.Has("hridContainsUUID")).To(BeFalse())
		Expect(req.Query.Has("ignoreGroups")).To(BeFalse())
		Expect(req.Body).To(HaveKeyWithValue("hrid", ordersHRID()))

		Expect(*resp.State.Id).To(Equal("6b1d5f0a-0000-0000-0000-000000000001"))
		Expect(resp.GetOrgID()).To(Equal("org1"))
		Expect(resp.GetEnvID()).To(Equal("env1"))
		Expect(resp.HridContainsUUID).To(BeFalse())
	})

	It("leaves the selector fields the CRD does not set out of the body", func() {
		_, err := apiv4.Upsert(ctx, client, toAPIV4DTO(apiV4ForSDK()))
		Expect(err).ToNot(HaveOccurred())

		flows, ok := fake.lastRequest().Body["flows"].([]any)
		Expect(ok).To(BeTrue())
		selector := flows[0].(map[string]any)["selectors"].([]any)[0]
		Expect(selector).To(Equal(map[string]any{"type": "HTTP"}))
	})

	It("upserts and dry-runs a legacy API through UpsertV4Legacy, with the same body as the SDK path", func() {
		_, err := apiv4.Upsert(ctx, client, toAPIV4DTO(apiV4ForSDK()))
		Expect(err).ToNot(HaveOccurred())
		sdkBody := fake.lastRequest().Body

		_, err = apiv4.Upsert(ctx, client, toAPIV4DTO(legacy()))
		Expect(err).ToNot(HaveOccurred())
		req := fake.lastRequest()
		Expect(req.Method).To(Equal(http.MethodPut))
		Expect(req.Path).To(Equal(apisPath))
		Expect(req.Query.Get("dryRun")).To(Equal("false"))
		Expect(req.Query.Get("hridContainsUUID")).To(Equal("true"))
		Expect(req.Body).To(HaveKeyWithValue("hrid", legacyAPIID))
		// the legacy API sends its spec.id, so its plan IDs differ: compare everything else
		Expect(withoutIDs(withoutHrid(req.Body))).To(Equal(withoutIDs(withoutHrid(sdkBody))))

		_, err = apiv4.DryRun(ctx, client, toAPIV4DTO(legacy()))
		Expect(err).ToNot(HaveOccurred())
		req = fake.lastRequest()
		Expect(req.Query.Get("dryRun")).To(Equal("true"))
		Expect(req.Query.Get("hridContainsUUID")).To(Equal("true"))
	})

	It("does not become Automation API managed when upserted as a legacy API", func() {
		api := legacy()
		resp, err := apiv4.Upsert(ctx, client, toAPIV4DTO(api))
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.HridContainsUUID).To(BeTrue())

		Expect(apiv4.UpdateStatus(ctx, api, resp)).To(Succeed())
		Expect(k8s.IsAutomationAPIManaged(api)).To(BeFalse())
		Expect(api.Status.Plans).To(HaveKeyWithValue("KeyLess", "6b1d5f0a-0000-0000-0000-000000000003"))
	})

	DescribeTable("asks APIM to leave groups alone under the ignore-groups annotation",
		func(api func() *v1alpha1.ApiV4Definition) {
			obj := api()
			obj.Annotations = map[string]string{core.IgnoreGroupsAnnotation: "true"}
			_, err := apiv4.Upsert(ctx, client, toAPIV4DTO(obj))
			Expect(err).ToNot(HaveOccurred())
			Expect(fake.lastRequest().Query.Get("ignoreGroups")).To(Equal("true"))
		},
		Entry("through the SDK", apiV4ForSDK),
		Entry("through UpsertV4Legacy", legacy),
	)

	It("returns the dry run warnings APIM reports when there is nothing severe", func() {
		errs, err := apiv4.DryRun(ctx, client, toAPIV4DTO(apiV4ForSDK()))
		Expect(err).ToNot(HaveOccurred())
		Expect(fake.lastRequest().Query.Get("dryRun")).To(Equal("true"))
		Expect(errs.IsSevere()).To(BeFalse())
		Expect(errs.Warning).To(ConsistOf(MatchError("warning one")))
	})

	It("returns only the severe dry run errors when there are any", func() {
		fake.response = `{"hrid": "orders", "name": "orders", "errors": {"severe": ["severe one"], "warning": ["warning one"]}}`
		errs, err := apiv4.DryRun(ctx, client, toAPIV4DTO(apiV4ForSDK()))
		Expect(err).ToNot(HaveOccurred())
		Expect(errs.Severe).To(ConsistOf(MatchError("severe one")))
		Expect(errs.Warning).To(BeEmpty())
	})

	It("deletes a legacy API by the UUID in its status", func() {
		api := apiV4ForSDK()
		api.Status.ID = legacyAPIID
		fake.status = http.StatusNoContent

		Expect(apiv4.Delete(ctx, client, toAPIV4DTO(api))).To(Succeed())

		req := fake.lastRequest()
		Expect(req.Method).To(Equal(http.MethodDelete))
		Expect(req.Path).To(Equal(apisPath + "/" + legacyAPIID))
		Expect(req.Query.Get("hridContainsUUID")).To(Equal("true"))
	})

	It("treats an API already gone as deleted", func() {
		fake.status = http.StatusNotFound

		Expect(apiv4.Delete(ctx, client, toAPIV4DTO(apiV4ForSDK()))).To(Succeed())

		req := fake.lastRequest()
		Expect(req.Path).To(Equal(apisPath + "/" + ordersHRID()))
		Expect(req.Query.Get("hridContainsUUID")).To(Equal("false"))
	})

	It("gets an API by hrid, or a legacy one by UUID", func() {
		synced := apiV4ForSDK()
		synced.Status.ID = legacyAPIID
		k8s.AddAutomationAPIManagedCondition(synced)

		remote, err := apiv4.GetRemote(ctx, client, toAPIV4DTO(synced))
		Expect(err).ToNot(HaveOccurred())
		Expect(remote.Name).To(Equal("orders"))
		req := fake.lastRequest()
		Expect(req.Method).To(Equal(http.MethodGet))
		Expect(req.Path).To(Equal(apisPath + "/" + ordersHRID()))
		Expect(req.Query.Get("hridContainsUUID")).To(Equal("false"))

		legacyAPI := legacy()
		legacyAPI.Status.ID = legacyAPIID
		_, err = apiv4.GetRemote(ctx, client, toAPIV4DTO(legacyAPI))
		Expect(err).ToNot(HaveOccurred())
		req = fake.lastRequest()
		Expect(req.Path).To(Equal(apisPath + "/" + legacyAPIID))
		Expect(req.Query.Get("hridContainsUUID")).To(Equal("true"))
	})

	It("fetches nothing for an API never synced", func() {
		dto := toAPIV4DTO(apiV4ForSDK())
		remote, err := apiv4.GetRemote(ctx, client, dto)
		Expect(err).ToNot(HaveOccurred())
		Expect(remote.Groups).To(Equal(dto.Groups))
		Expect(fake.lastRequest().Method).To(BeEmpty())
	})

	It("compares the CR's groups under ignore-groups, whatever APIM holds", func() {
		api := apiV4ForSDK()
		api.Status.ID = legacyAPIID
		k8s.AddAutomationAPIManagedCondition(api)

		remote, err := apiv4.GetRemote(ctx, client, toAPIV4DTO(api))
		Expect(err).ToNot(HaveOccurred())
		Expect(remote.Groups).To(ConsistOf("remote-group"))

		api.Annotations = map[string]string{core.IgnoreGroupsAnnotation: "true"}
		remote, err = apiv4.GetRemote(ctx, client, toAPIV4DTO(api))
		Expect(err).ToNot(HaveOccurred())
		Expect(remote.Groups).To(ConsistOf("crd-group"))
	})
})

// withoutIDs drops the API and plan ids from a request body: a legacy API sends its spec.id, and the plan ids
// are derived from it.
func withoutIDs(body map[string]any) map[string]any {
	cp := make(map[string]any, len(body))
	for key, value := range body {
		if key != "id" && key != "crossId" {
			cp[key] = value
		}
	}
	plans, ok := body["plans"].([]any)
	if !ok {
		return cp
	}
	stripped := make([]any, len(plans))
	for i, plan := range plans {
		fields := map[string]any{}
		for key, value := range plan.(map[string]any) {
			if key != "id" && key != "crossId" {
				fields[key] = value
			}
		}
		stripped[i] = fields
	}
	cp["plans"] = stripped
	return cp
}
