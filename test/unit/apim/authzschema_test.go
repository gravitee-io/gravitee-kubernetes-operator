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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzschema"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimauthzschema "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/authzschema"
	xErrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

const travelSchemaText = "entity User in [Group];\nentity Group;\nentity Hotel;\n"

func newAuthzSchema() *v1alpha1.AuthzSchema {
	return &v1alpha1.AuthzSchema{
		ObjectMeta: metav1.ObjectMeta{Name: "travel-schema", Namespace: "gravitee"},
		Spec: v1alpha1.AuthzSchemaSpec{
			Context: &refs.NamespacedName{Name: "dev-ctx"},
			Type: authzschema.Type{
				Name:   "Travel schema",
				Target: &authzpdp.Target{PdpRef: &refs.NamespacedName{Name: "orders-eu"}},
				Schema: travelSchemaText,
			},
		},
	}
}

// resolvedAuthzSchema is newAuthzSchema with its pdpRef resolved, as dynamic.ResolveAuthzSchemaRefs leaves it.
func resolvedAuthzSchema() *v1alpha1.AuthzSchema {
	schema := newAuthzSchema()
	schema.Spec.Target = &authzpdp.Target{Scope: "orders@eu"}
	return schema
}

func toAuthzSchemaDTO(schema *v1alpha1.AuthzSchema) sdk.AuthzSchema {
	GinkgoHelper()
	dto, err := apimauthzschema.ToDTO(schema)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

// authzSchemaState is what the platform answers for travel-schema: the payload, the ids and the findings.
const authzSchemaState = `{
	"hrid": "gravitee-travel-schema",
	"name": "Travel schema",
	"target": "orders@eu",
	"schema": "entity User in [Group];\nentity Group;\nentity Hotel;\n",
	"id": "7a1e0c2b", "organizationId": "DEFAULT", "environmentId": "DEFAULT",
	"errors": {"warning": ["entity Hotel has no attributes"]}
}`

var _ = Describe("Authz schema", func() {
	Describe("spec hash", func() {
		DescribeTable("changes with each field",
			func(change func(schema *v1alpha1.AuthzSchema)) {
				changed := newAuthzSchema()
				change(changed)
				Expect(changed.Spec.Hash()).ToNot(Equal(newAuthzSchema().Spec.Hash()))
			},
			Entry("name", func(schema *v1alpha1.AuthzSchema) { schema.Spec.Name = "Travel" }),
			Entry("schema", func(schema *v1alpha1.AuthzSchema) { schema.Spec.Schema = "entity User;\n" }),
			Entry("pdpRef", func(schema *v1alpha1.AuthzSchema) { schema.Spec.Target.PdpRef.Name = "orders-us" }),
			Entry("scope instead of pdpRef", func(schema *v1alpha1.AuthzSchema) {
				schema.Spec.Target = &authzpdp.Target{Scope: "orders@eu"}
			}),
			Entry("no target", func(schema *v1alpha1.AuthzSchema) { schema.Spec.Target = nil }),
			Entry("context", func(schema *v1alpha1.AuthzSchema) { schema.Spec.Context.Name = "prod-ctx" }),
		)

		It("is stable for the same spec", func() {
			Expect(newAuthzSchema().Spec.Hash()).To(Equal(newAuthzSchema().Spec.Hash()))
		})
	})

	Describe("wire mapping", func() {
		It("derives the hrid from the namespace and name", func() {
			Expect(toAuthzSchemaDTO(resolvedAuthzSchema()).Hrid).To(Equal("gravitee-travel-schema"))
		})

		It("sends the scope of the resolved target", func() {
			Expect(marshalJSON(toAuthzSchemaDTO(resolvedAuthzSchema()))).To(MatchJSON(`{
				"hrid": "gravitee-travel-schema",
				"name": "Travel schema",
				"target": "orders@eu",
				"schema": "entity User in [Group];\nentity Group;\nentity Hotel;\n"
			}`))
		})

		It("omits an absent target", func() {
			untargeted := newAuthzSchema()
			untargeted.Spec.Target = nil

			Expect(marshalJSON(toAuthzSchemaDTO(untargeted))).To(MatchJSON(`{
				"hrid": "gravitee-travel-schema",
				"name": "Travel schema",
				"schema": "entity User in [Group];\nentity Group;\nentity Hotel;\n"
			}`))
		})

		// Sent as no target, the document would go to the default engine.
		It("refuses a pdpRef left unresolved", func() {
			_, err := apimauthzschema.ToDTO(newAuthzSchema())

			Expect(err).To(MatchError("target.pdpRef [gravitee/orders-eu] is not resolved"))
		})
	})

	Describe("status", func() {
		It("copies the ids, the findings and the target", func() {
			var state sdk.AuthzSchemaState
			Expect(json.Unmarshal([]byte(authzSchemaState), &state)).To(Succeed())
			schema := newAuthzSchema()

			Expect(apimauthzschema.UpdateStatus(context.Background(), schema,
				apimauthzschema.Response{State: state})).To(Succeed())

			Expect(schema.Status.ID).To(Equal("7a1e0c2b"))
			Expect(schema.Status.OrgID).To(Equal("DEFAULT"))
			Expect(schema.Status.EnvID).To(Equal("DEFAULT"))
			Expect(schema.Status.HRID).To(Equal("gravitee-travel-schema"))
			Expect(schema.Status.Target).To(Equal("orders@eu"))
			Expect(schema.Status.Errors.Warning).To(ConsistOf("entity Hotel has no attributes"))
			Expect(k8s.IsAutomationAPIManaged(schema)).To(BeTrue())
		})
	})

	Describe("calls", func() {
		const collection = "/automation/organizations/org1/environments/env1/authz/schemas"

		var (
			ctx    context.Context
			fake   *fakeAutomationAPI
			client *apim.SDKClient
		)

		BeforeEach(func() {
			ctx = context.Background()
			fake = newFakeAutomationAPI()
			fake.response = authzSchemaState
			DeferCleanup(fake.Close)
			var err error
			client, err = apim.NewSDKClient(ctx, managementContext(fake.URL, bearer("token")))
			Expect(err).ToNot(HaveOccurred())
		})

		It("puts the payload on the schemas collection and answers the state", func() {
			resp, err := apimauthzschema.Upsert(ctx, client, toAuthzSchemaDTO(resolvedAuthzSchema()))
			Expect(err).ToNot(HaveOccurred())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodPut))
			Expect(req.Path).To(Equal(collection))
			Expect(req.Query.Get("dryRun")).To(Equal("false"))
			Expect(marshalJSON(req.Body)).To(MatchJSON(`{
				"hrid": "gravitee-travel-schema",
				"name": "Travel schema",
				"target": "orders@eu",
				"schema": "entity User in [Group];\nentity Group;\nentity Hotel;\n"
			}`))
			Expect(*resp.State.Id).To(Equal("7a1e0c2b"))
			Expect(resp.GetOrgID()).To(Equal("DEFAULT"))
			Expect(resp.GetEnvID()).To(Equal("DEFAULT"))
		})

		It("dry runs and returns the findings of the state", func() {
			fake.response = `{"hrid": "gravitee-travel-schema", "name": "Travel schema", "schema": "",
				"errors": {"severe": ["target [orders@eu] already holds schema document 'Orders schema'"]}}`

			findings, err := apimauthzschema.DryRun(ctx, client, toAuthzSchemaDTO(resolvedAuthzSchema()))
			Expect(err).ToNot(HaveOccurred())

			Expect(fake.lastRequest().Query.Get("dryRun")).To(Equal("true"))
			Expect(findings.Severe).To(HaveLen(1))
			Expect(findings.Severe[0].Message).To(Equal("target [orders@eu] already holds schema document 'Orders schema'"))
		})

		It("returns the findings of a refused dry run", func() {
			fake.status = http.StatusBadRequest
			fake.errorResponse = `{"hrid": "gravitee-travel-schema",
				"errors": {"severe": ["schema does not parse: unexpected token at 1:7"]}}`

			findings, err := apimauthzschema.DryRun(ctx, client, toAuthzSchemaDTO(resolvedAuthzSchema()))
			Expect(err).ToNot(HaveOccurred())

			Expect(findings.Severe).To(HaveLen(1))
			Expect(findings.Severe[0].Message).To(Equal("schema does not parse: unexpected token at 1:7"))
		})

		It("returns a 400 that carries no findings as the error", func() {
			fake.status = http.StatusBadRequest
			fake.errorResponse = `{"message": "hrid is required", "http_status": 400}`

			_, err := apimauthzschema.DryRun(ctx, client, toAuthzSchemaDTO(resolvedAuthzSchema()))

			Expect(xErrors.IsBadRequest(err)).To(BeTrue())
		})

		It("gets the remote as the payload a PUT sends", func() {
			remote, err := apimauthzschema.GetRemote(ctx, client, toAuthzSchemaDTO(resolvedAuthzSchema()))
			Expect(err).ToNot(HaveOccurred())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodGet))
			Expect(req.Path).To(Equal(collection + "/gravitee-travel-schema"))
			Expect(remote).To(Equal(toAuthzSchemaDTO(resolvedAuthzSchema())))
		})

		It("reports a remote that is gone as not found", func() {
			fake.status = http.StatusNotFound

			_, err := apimauthzschema.GetRemote(ctx, client, toAuthzSchemaDTO(resolvedAuthzSchema()))

			Expect(xErrors.IsNotFound(err)).To(BeTrue())
		})

		It("deletes by hrid and treats a schema already gone as deleted", func() {
			dto := toAuthzSchemaDTO(resolvedAuthzSchema())

			Expect(apimauthzschema.Delete(ctx, client, dto)).To(Succeed())
			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodDelete))
			Expect(req.Path).To(Equal(collection + "/gravitee-travel-schema"))

			fake.status = http.StatusNotFound
			Expect(apimauthzschema.Delete(ctx, client, dto)).To(Succeed())
		})
	})

	Describe("a refused apply", func() {
		refusedWith := func(statusCode int, body string) error {
			return xErrors.ServerError{StatusCode: statusCode, Body: body}
		}

		It("records the findings in the status, keeps the target and stays recoverable", func() {
			schema := newAuthzSchema()
			schema.Status.Target = "orders@eu"

			err := apimauthzschema.OnUpsertError(context.Background(), schema, refusedWith(http.StatusBadRequest,
				`{"errors": {"severe": ["target [orders@eu] already holds schema document 'Orders schema'"]}}`))

			Expect(err).To(MatchError(ContainSubstring("authz schema [travel-schema] was refused by the automation api")))
			Expect(xErrors.IsRecoverable(err)).To(BeTrue())
			Expect(schema.Status.Errors.Severe).To(ConsistOf("target [orders@eu] already holds schema document 'Orders schema'"))
			Expect(schema.Status.Target).To(Equal("orders@eu"))
		})

		It("stays recoverable when the 400 carries no findings", func() {
			schema := newAuthzSchema()

			err := apimauthzschema.OnUpsertError(context.Background(), schema,
				refusedWith(http.StatusBadRequest, `{"message": "hrid is required", "http_status": 400}`))

			Expect(xErrors.IsRecoverable(err)).To(BeTrue())
			Expect(schema.Status.Errors.Severe).To(BeEmpty())
		})

		It("leaves any other failure to the lifecycle", func() {
			err := apimauthzschema.OnUpsertError(context.Background(), newAuthzSchema(),
				refusedWith(http.StatusInternalServerError, `{"errors": {"severe": ["boom"]}}`))

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
