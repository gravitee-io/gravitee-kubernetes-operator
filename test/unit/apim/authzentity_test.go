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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimauthzentity "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/authzentity"
	xErrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

const ownershipFinding = `entityId [agents] is held by Group::"agents", which this manifest did not create ` +
	`(source: gravitee_am). Declare another entityId, or delete that entity first`

func newAuthzEntity() *v1alpha1.AuthzEntity {
	return &v1alpha1.AuthzEntity{
		ObjectMeta: metav1.ObjectMeta{Name: "agents-group", Namespace: "gravitee"},
		Spec: v1alpha1.AuthzEntitySpec{
			Context: &refs.NamespacedName{Name: "dev-ctx"},
			Type: authzentity.Type{
				EntityID:   "agents",
				Kind:       "PRINCIPAL",
				EntityType: "Group",
				Attributes: utils.ToGenericStringMap(map[string]any{"department": "travel"}),
				Parents:    []authzentity.Parent{{EntityType: "Group", EntityID: "travel"}},
				Active:     new(true),
				Targets: []authzpdp.Target{
					{PdpRef: &refs.NamespacedName{Name: "orders-eu"}},
					{Scope: "default"},
				},
			},
		},
	}
}

// resolvedAuthzEntity is newAuthzEntity with its pdpRef resolved, as dynamic.ResolveAuthzEntityRefs leaves it.
func resolvedAuthzEntity() *v1alpha1.AuthzEntity {
	entity := newAuthzEntity()
	entity.Spec.Targets = []authzpdp.Target{{Scope: "orders@eu"}, {Scope: "default"}}
	return entity
}

func toAuthzEntityDTO(entity *v1alpha1.AuthzEntity) sdk.AuthzEntity {
	GinkgoHelper()
	dto, err := apimauthzentity.ToDTO(entity)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

// authzEntityState is what the platform answers for agents-group: the payload, its lists sorted, the ids and the
// findings.
const authzEntityState = `{
	"hrid": "gravitee-agents-group",
	"entityId": "agents",
	"kind": "PRINCIPAL",
	"entityType": "Group",
	"attributes": {"department": "travel"},
	"parents": [{"entityType": "Group", "entityId": "travel"}],
	"active": true,
	"targets": ["default", "orders@eu"],
	"id": "3b9d0c2b", "organizationId": "DEFAULT", "environmentId": "DEFAULT",
	"errors": {"warning": ["attribute [department] is not declared by the schema of Group"]}
}`

var _ = Describe("Authz entity", func() {
	Describe("spec hash", func() {
		DescribeTable("changes with each field",
			func(change func(entity *v1alpha1.AuthzEntity)) {
				changed := newAuthzEntity()
				change(changed)
				Expect(changed.Spec.Hash()).ToNot(Equal(newAuthzEntity().Spec.Hash()))
			},
			Entry("entityId", func(entity *v1alpha1.AuthzEntity) { entity.Spec.EntityID = "travel" }),
			Entry("kind", func(entity *v1alpha1.AuthzEntity) { entity.Spec.Kind = "RESOURCE" }),
			Entry("entityType", func(entity *v1alpha1.AuthzEntity) { entity.Spec.EntityType = "Team" }),
			Entry("attributes", func(entity *v1alpha1.AuthzEntity) { entity.Spec.Attributes.Put("department", "sales") }),
			Entry("parents", func(entity *v1alpha1.AuthzEntity) { entity.Spec.Parents[0].EntityID = "sales" }),
			Entry("active", func(entity *v1alpha1.AuthzEntity) { entity.Spec.Active = new(false) }),
			Entry("pdpRef", func(entity *v1alpha1.AuthzEntity) { entity.Spec.Targets[0].PdpRef.Name = "orders-us" }),
			Entry("scope instead of pdpRef", func(entity *v1alpha1.AuthzEntity) {
				entity.Spec.Targets[0] = authzpdp.Target{Scope: "orders@eu"}
			}),
			Entry("no targets", func(entity *v1alpha1.AuthzEntity) { entity.Spec.Targets = nil }),
			Entry("context", func(entity *v1alpha1.AuthzEntity) { entity.Spec.Context.Name = "prod-ctx" }),
		)

		It("is stable for the same spec", func() {
			Expect(newAuthzEntity().Spec.Hash()).To(Equal(newAuthzEntity().Spec.Hash()))
		})
	})

	Describe("wire mapping", func() {
		It("derives the hrid from the namespace and name", func() {
			Expect(toAuthzEntityDTO(resolvedAuthzEntity()).Hrid).To(Equal("gravitee-agents-group"))
		})

		It("sends each target's scope, in declaration order", func() {
			Expect(marshalJSON(toAuthzEntityDTO(resolvedAuthzEntity()))).To(MatchJSON(`{
				"hrid": "gravitee-agents-group",
				"entityId": "agents",
				"kind": "PRINCIPAL",
				"entityType": "Group",
				"attributes": {"department": "travel"},
				"parents": [{"entityType": "Group", "entityId": "travel"}],
				"active": true,
				"targets": ["orders@eu", "default"]
			}`))
		})

		It("omits absent attributes, parents and targets", func() {
			entity := newAuthzEntity()
			entity.Spec.Attributes = nil
			entity.Spec.Parents = nil
			entity.Spec.Targets = nil

			Expect(marshalJSON(toAuthzEntityDTO(entity))).To(MatchJSON(`{
				"hrid": "gravitee-agents-group",
				"entityId": "agents",
				"kind": "PRINCIPAL",
				"entityType": "Group",
				"active": true
			}`))
		})

		It("sends an unset active as true", func() {
			entity := resolvedAuthzEntity()
			entity.Spec.Active = nil

			Expect(toAuthzEntityDTO(entity).Active).To(HaveValue(BeTrue()))
		})

		It("sends active false", func() {
			entity := resolvedAuthzEntity()
			entity.Spec.Active = new(false)

			Expect(marshalJSON(toAuthzEntityDTO(entity))).To(ContainSubstring(`"active":false`))
		})

		It("does not share the attributes with the resource", func() {
			entity := resolvedAuthzEntity()
			dto := toAuthzEntityDTO(entity)
			dto.Attributes.Put("department", "sales")

			Expect(entity.Spec.Attributes.GetString("department")).To(Equal("travel"))
		})

		// A pdpRef has no scope to send until its AuthzPdp's routing scope replaces it.
		It("refuses a pdpRef left unresolved, naming its index", func() {
			_, err := apimauthzentity.ToDTO(newAuthzEntity())

			Expect(err).To(MatchError("targets[0].pdpRef [gravitee/orders-eu] is not resolved"))
		})
	})

	Describe("status", func() {
		It("copies the ids and the findings", func() {
			var state sdk.AuthzEntityState
			Expect(json.Unmarshal([]byte(authzEntityState), &state)).To(Succeed())
			entity := newAuthzEntity()

			Expect(apimauthzentity.UpdateStatus(context.Background(), entity,
				apimauthzentity.Response{State: state})).To(Succeed())

			Expect(entity.Status.ID).To(Equal("3b9d0c2b"))
			Expect(entity.Status.OrgID).To(Equal("DEFAULT"))
			Expect(entity.Status.EnvID).To(Equal("DEFAULT"))
			Expect(entity.Status.HRID).To(Equal("gravitee-agents-group"))
			Expect(entity.Status.Errors.Warning).To(ConsistOf("attribute [department] is not declared by the schema of Group"))
			Expect(k8s.IsAutomationAPIManaged(entity)).To(BeTrue())
		})
	})

	Describe("calls", func() {
		const collection = "/automation/organizations/org1/environments/env1/authz/entities"

		var (
			ctx    context.Context
			fake   *fakeAutomationAPI
			client *apim.SDKClient
		)

		BeforeEach(func() {
			ctx = context.Background()
			fake = newFakeAutomationAPI()
			fake.response = authzEntityState
			DeferCleanup(fake.Close)
			var err error
			client, err = apim.NewSDKClient(ctx, managementContext(fake.URL, bearer("token")))
			Expect(err).ToNot(HaveOccurred())
		})

		It("puts the payload on the entities collection and answers the state", func() {
			resp, err := apimauthzentity.Upsert(ctx, client, toAuthzEntityDTO(resolvedAuthzEntity()))
			Expect(err).ToNot(HaveOccurred())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodPut))
			Expect(req.Path).To(Equal(collection))
			Expect(req.Query.Get("dryRun")).To(Equal("false"))
			Expect(marshalJSON(req.Body)).To(MatchJSON(marshalJSON(toAuthzEntityDTO(resolvedAuthzEntity()))))
			Expect(*resp.State.Id).To(Equal("3b9d0c2b"))
			Expect(resp.GetOrgID()).To(Equal("DEFAULT"))
			Expect(resp.GetEnvID()).To(Equal("DEFAULT"))
		})

		It("dry runs and returns the findings of the state", func() {
			fake.response = `{"hrid": "gravitee-agents-group", "entityId": "agents", "kind": "PRINCIPAL",
				"entityType": "Group", "errors": {"severe": [` + mustJSON(ownershipFinding) + `]}}`

			findings, err := apimauthzentity.DryRun(ctx, client, toAuthzEntityDTO(resolvedAuthzEntity()))
			Expect(err).ToNot(HaveOccurred())

			Expect(fake.lastRequest().Query.Get("dryRun")).To(Equal("true"))
			Expect(findings.Severe).To(HaveLen(1))
			Expect(findings.Severe[0].Message).To(Equal(ownershipFinding))
		})

		It("returns the findings of a refused dry run", func() {
			fake.status = http.StatusBadRequest
			fake.errorResponse = `{"hrid": "gravitee-agents-group",
				"errors": {"severe": ["attribute key [_owner] is reserved"]}}`

			findings, err := apimauthzentity.DryRun(ctx, client, toAuthzEntityDTO(resolvedAuthzEntity()))
			Expect(err).ToNot(HaveOccurred())

			Expect(findings.Severe).To(HaveLen(1))
			Expect(findings.Severe[0].Message).To(Equal("attribute key [_owner] is reserved"))
		})

		It("returns a 400 that carries no findings as the error", func() {
			fake.status = http.StatusBadRequest
			fake.errorResponse = `{"message": "hrid is required", "http_status": 400}`

			_, err := apimauthzentity.DryRun(ctx, client, toAuthzEntityDTO(resolvedAuthzEntity()))

			Expect(xErrors.IsBadRequest(err)).To(BeTrue())
		})

		It("gets the remote as the payload a PUT sends, its lists as the platform sorts them", func() {
			remote, err := apimauthzentity.GetRemote(ctx, client, toAuthzEntityDTO(resolvedAuthzEntity()))
			Expect(err).ToNot(HaveOccurred())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodGet))
			Expect(req.Path).To(Equal(collection + "/gravitee-agents-group"))
			Expect(marshalJSON(remote)).To(MatchJSON(`{
				"hrid": "gravitee-agents-group",
				"entityId": "agents",
				"kind": "PRINCIPAL",
				"entityType": "Group",
				"attributes": {"department": "travel"},
				"parents": [{"entityType": "Group", "entityId": "travel"}],
				"active": true,
				"targets": ["default", "orders@eu"]
			}`))
		})

		It("reports a remote that is gone as not found", func() {
			fake.status = http.StatusNotFound

			_, err := apimauthzentity.GetRemote(ctx, client, toAuthzEntityDTO(resolvedAuthzEntity()))

			Expect(xErrors.IsNotFound(err)).To(BeTrue())
		})

		It("deletes by hrid and treats an entity already gone as deleted", func() {
			dto := toAuthzEntityDTO(resolvedAuthzEntity())

			Expect(apimauthzentity.Delete(ctx, client, dto)).To(Succeed())
			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodDelete))
			Expect(req.Path).To(Equal(collection + "/gravitee-agents-group"))
			Expect(req.Query.Has("dryRun")).To(BeFalse())

			fake.status = http.StatusNotFound
			Expect(apimauthzentity.Delete(ctx, client, dto)).To(Succeed())
		})

		// The lifecycle retries a failed delete, keeping the finalizer until the cascade fits.
		It("returns the 413 of a cascade over the platform's limit", func() {
			fake.status = http.StatusRequestEntityTooLarge
			fake.errorResponse = `{"message": "deleting agents would delete 1200 entities, over the limit of 1000",
				"http_status": 413}`

			err := apimauthzentity.Delete(ctx, client, toAuthzEntityDTO(resolvedAuthzEntity()))

			Expect(err).To(MatchError(ContainSubstring("over the limit of 1000")))
		})
	})

	Describe("a refused apply", func() {
		refusedWith := func(statusCode int, body string) error {
			return xErrors.ServerError{StatusCode: statusCode, Body: body}
		}

		It("records the findings in the status and stays recoverable", func() {
			entity := newAuthzEntity()

			err := apimauthzentity.OnUpsertError(context.Background(), entity, refusedWith(http.StatusBadRequest,
				`{"errors": {"severe": [`+mustJSON(ownershipFinding)+`]}}`))

			Expect(err).To(MatchError(ContainSubstring("authz entity [agents-group] was refused by the automation api")))
			Expect(xErrors.IsRecoverable(err)).To(BeTrue())
			Expect(entity.Status.Errors.Severe).To(ConsistOf(ownershipFinding))
		})

		It("stays recoverable when the 400 carries no findings", func() {
			entity := newAuthzEntity()

			err := apimauthzentity.OnUpsertError(context.Background(), entity,
				refusedWith(http.StatusBadRequest, `{"message": "hrid is required", "http_status": 400}`))

			Expect(xErrors.IsRecoverable(err)).To(BeTrue())
			Expect(entity.Status.Errors.Severe).To(BeEmpty())
		})

		It("leaves any other failure to the lifecycle", func() {
			err := apimauthzentity.OnUpsertError(context.Background(), newAuthzEntity(),
				refusedWith(http.StatusInternalServerError, `{"errors": {"severe": ["boom"]}}`))

			Expect(err).ToNot(HaveOccurred())
		})
	})
})

func mustJSON(s string) string {
	raw, err := json.Marshal(s)
	Expect(err).ToNot(HaveOccurred())
	return string(raw)
}
