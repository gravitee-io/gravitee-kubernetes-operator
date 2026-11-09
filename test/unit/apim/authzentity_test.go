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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzentity"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzpdp"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
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
			Expect(model.ToAuthzEntityDTO(newAuthzEntity()).HRID).To(Equal("gravitee-agents-group"))
		})

		// dynamic.ResolveAuthzEntityRefs has replaced each pdpRef with its AuthzPdp's routing scope.
		It("sends each target's scope, in declaration order", func() {
			entity := newAuthzEntity()
			entity.Spec.Targets = []authzpdp.Target{{Scope: "orders@eu"}, {Scope: "default"}}
			raw, err := json.Marshal(model.ToAuthzEntityDTO(entity))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(raw)).To(MatchJSON(`{
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
			raw, err := json.Marshal(model.ToAuthzEntityDTO(entity))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(raw)).To(MatchJSON(`{
				"hrid": "gravitee-agents-group",
				"entityId": "agents",
				"kind": "PRINCIPAL",
				"entityType": "Group",
				"active": true
			}`))
		})

		It("sends an unset active as true", func() {
			entity := newAuthzEntity()
			entity.Spec.Active = nil
			Expect(model.ToAuthzEntityDTO(entity).Active).To(BeTrue())
		})

		It("sends active false", func() {
			entity := newAuthzEntity()
			entity.Spec.Active = new(false)
			raw, err := json.Marshal(model.ToAuthzEntityDTO(entity))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(raw)).To(ContainSubstring(`"active":false`))
		})

		It("does not share the attributes with the resource", func() {
			entity := newAuthzEntity()
			dto := model.ToAuthzEntityDTO(entity)
			dto.Attributes.Put("department", "sales")
			Expect(entity.Spec.Attributes.GetString("department")).To(Equal("travel"))
		})
	})

	Describe("state decoding", func() {
		It("decodes the payload, the ids and the findings", func() {
			body := `{
				"hrid": "gravitee-agents-group",
				"entityId": "agents",
				"kind": "PRINCIPAL",
				"entityType": "Group",
				"attributes": {"department": "travel"},
				"parents": [{"entityType": "Group", "entityId": "travel"}],
				"active": false,
				"targets": ["default", "orders@eu"],
				"id": "3b9d0c2b", "organizationId": "DEFAULT", "environmentId": "DEFAULT",
				"errors": {"severe": [` + mustJSON(ownershipFinding) + `]}
			}`
			state := new(model.AuthzEntityState)
			Expect(json.Unmarshal([]byte(body), state)).To(Succeed())

			Expect(state.HRID).To(Equal("gravitee-agents-group"))
			Expect(state.EntityID).To(Equal("agents"))
			Expect(state.Kind).To(Equal("PRINCIPAL"))
			Expect(state.EntityType).To(Equal("Group"))
			Expect(state.Attributes.GetString("department")).To(Equal("travel"))
			Expect(state.Parents).To(Equal([]model.AuthzEntityParentDTO{{EntityType: "Group", EntityID: "travel"}}))
			Expect(state.Active).To(BeFalse())
			Expect(state.Targets).To(Equal([]string{"default", "orders@eu"}))
			Expect(state.ID).To(Equal("3b9d0c2b"))
			Expect(state.OrgID).To(Equal("DEFAULT"))
			Expect(state.EnvID).To(Equal("DEFAULT"))
			Expect(state.Errors.Severe).To(ConsistOf(ownershipFinding))
		})
	})
})

func mustJSON(s string) string {
	raw, err := json.Marshal(s)
	Expect(err).ToNot(HaveOccurred())
	return string(raw)
}
