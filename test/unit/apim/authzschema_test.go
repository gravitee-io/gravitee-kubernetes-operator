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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/authzschema"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
)

const travelSchemaText = "entity User in Group;\nentity Hotel;\n"

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
			Expect(model.ToAuthzSchemaDTO(newAuthzSchema(), "orders@eu").HRID).To(Equal("gravitee-travel-schema"))
		})

		It("sends the target the caller resolved, not the pdpRef", func() {
			raw, err := json.Marshal(model.ToAuthzSchemaDTO(newAuthzSchema(), "orders@eu"))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(raw)).To(MatchJSON(`{
				"hrid": "gravitee-travel-schema",
				"name": "Travel schema",
				"target": "orders@eu",
				"schema": "entity User in Group;\nentity Hotel;\n"
			}`))
		})

		It("omits an empty target", func() {
			raw, err := json.Marshal(model.ToAuthzSchemaDTO(newAuthzSchema(), ""))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(raw)).To(MatchJSON(`{
				"hrid": "gravitee-travel-schema",
				"name": "Travel schema",
				"schema": "entity User in Group;\nentity Hotel;\n"
			}`))
		})
	})

	Describe("state decoding", func() {
		It("decodes the payload, the ids, the target and the findings", func() {
			body := `{
				"hrid": "gravitee-travel-schema",
				"name": "Travel schema",
				"target": "orders@eu",
				"schema": "entity User in Group;\nentity Hotel;\n",
				"id": "7a1e0c2b", "organizationId": "DEFAULT", "environmentId": "DEFAULT",
				"errors": {"severe": ["target [orders@eu] already holds schema document 'Orders schema'"]}
			}`
			state := new(model.AuthzSchemaState)
			Expect(json.Unmarshal([]byte(body), state)).To(Succeed())

			Expect(state.HRID).To(Equal("gravitee-travel-schema"))
			Expect(state.Name).To(Equal("Travel schema"))
			Expect(state.Target).To(Equal("orders@eu"))
			Expect(state.Schema).To(Equal(travelSchemaText))
			Expect(state.ID).To(Equal("7a1e0c2b"))
			Expect(state.OrgID).To(Equal("DEFAULT"))
			Expect(state.EnvID).To(Equal("DEFAULT"))
			Expect(state.Errors.Severe).To(ConsistOf("target [orders@eu] already holds schema document 'Orders schema'"))
		})
	})
})
