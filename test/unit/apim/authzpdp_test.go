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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
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
			Expect(model.ToAuthzPdpDTO(newAuthzPdp()).HRID).To(Equal("gravitee-orders-eu"))
		})

		It("serialises the payload with the wire field names", func() {
			raw, err := json.Marshal(model.ToAuthzPdpDTO(newAuthzPdp()))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(raw)).To(MatchJSON(`{
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

			raw, err := json.Marshal(model.ToAuthzPdpDTO(pdp))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(raw)).To(MatchJSON(`{
				"hrid": "gravitee-orders-eu",
				"name": "Orders PDP (EU)",
				"targetPdpId": "orders",
				"strict": false
			}`))
		})
	})

	Describe("state decoding", func() {
		It("decodes the payload, the ids, the routing scope and the discovery path", func() {
			body := `{
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
			state := new(model.AuthzPdpState)
			Expect(json.Unmarshal([]byte(body), state)).To(Succeed())

			Expect(state.HRID).To(Equal("gravitee-orders-eu"))
			Expect(state.TargetPdpID).To(Equal("orders"))
			Expect(state.Tag).To(Equal("eu"))
			Expect(state.Strict).To(BeTrue())
			Expect(state.Endpoint.Path).To(Equal("/authz/orders-eu"))
			Expect(state.ID).To(Equal("5d0c1e7a"))
			Expect(state.OrgID).To(Equal("DEFAULT"))
			Expect(state.EnvID).To(Equal("DEFAULT"))
			Expect(state.RoutingScope).To(Equal("orders@eu"))
			Expect(state.WellKnownPath).To(Equal("/authz/orders-eu/.well-known/authzen-configuration"))
			Expect(state.Errors.Warning).To(ConsistOf("strict is shared by every PDP of targetPdpId orders"))
		})
	})
})
