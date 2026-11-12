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

package apim

import (
	"strconv"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// GKO-3452: APIM stores each plan's order from its position in the list GKO sends,
// and the Automation API GET lists the plans in that stored order.
var _ = Describe("API v4 plan order", func() {
	DescribeTable("plans are sent sorted by order, unset orders last, then by key",
		func(crd string, expected []string) {
			Expect(sentPlanHRIDs(crd)).To(Equal(expected))
		},
		Entry("orders against the key order",
			`{"name": "api", "plans": {"foo": `+planWithOrder("foo", 1)+`, "bar": `+planWithOrder("bar", 2)+`}}`,
			[]string{"foo", "bar"},
		),
		Entry("orders 20 and 10",
			`{"name": "api", "plans": {"a": `+planWithOrder("a", 20)+`, "b": `+planWithOrder("b", 10)+`}}`,
			[]string{"b", "a"},
		),
		Entry("unset orders after the set ones, by key",
			`{"name": "api", "plans": {"c": `+plan("c")+`, "b": `+plan("b")+`, "z": `+planWithOrder("z", 5)+`}}`,
			[]string{"z", "b", "c"},
		),
		Entry("the same order, by key",
			`{"name": "api", "plans": {"b": `+planWithOrder("b", 1)+`, "a": `+planWithOrder("a", 1)+`}}`,
			[]string{"a", "b"},
		),
	)

	It("reads a remote plan's order from its position in the GET list", func() {
		remote := parseJSON[model.APIV4DTO](`{"name": "api", "plans": [` + plan("foo") + `, ` + plan("bar") + `]}`)

		Expect(remote.Plans).To(HaveLen(2))
		Expect(remote.Plans[0].HRID).To(Equal("foo"))
		Expect(remote.Plans[0].Order).To(HaveValue(Equal(1)))
		Expect(remote.Plans[1].HRID).To(Equal("bar"))
		Expect(remote.Plans[1].Order).To(HaveValue(Equal(2)))
	})

	DescribeTable("no drift when APIM keeps the order GKO sent",
		func(crd, remote string) {
			expectNoDrift(detectPlanOrderDrift(crd, remote))
		},
		Entry("orders against the key order",
			`{"name": "api", "plans": {"foo": `+planWithOrder("foo", 1)+`, "bar": `+planWithOrder("bar", 2)+`}}`,
			`{"name": "api", "plans": [`+plan("foo")+`, `+plan("bar")+`]}`,
		),
		Entry("orders 10 and 20 compared as their rank",
			`{"name": "api", "plans": {"foo": `+planWithOrder("foo", 10)+`, "bar": `+planWithOrder("bar", 20)+`}}`,
			`{"name": "api", "plans": [`+plan("foo")+`, `+plan("bar")+`]}`,
		),
		Entry("the same order, ranked by key",
			`{"name": "api", "plans": {"b": `+planWithOrder("b", 1)+`, "a": `+planWithOrder("a", 1)+`}}`,
			`{"name": "api", "plans": [`+plan("a")+`, `+plan("b")+`]}`,
		),
		Entry("plans without a status, published by default",
			`{"name": "api", "plans": {"foo": `+planWithOrderNoStatus("foo", 1)+`, "bar": `+
				planWithOrderNoStatus("bar", 2)+`}}`,
			`{"name": "api", "plans": [`+planNoStatus("foo")+`, `+planNoStatus("bar")+`]}`,
		),
	)

	DescribeTable("no drift on the plan order when a plan has no order, whatever order APIM lists",
		func(crd, remote string) {
			expectNoDrift(detectPlanOrderDrift(crd, remote))
		},
		Entry("no plan sets an order",
			`{"name": "api", "plans": {"foo": `+plan("foo")+`, "bar": `+plan("bar")+`}}`,
			`{"name": "api", "plans": [`+plan("foo")+`, `+plan("bar")+`]}`,
		),
		Entry("one plan without an order",
			`{"name": "api", "plans": {"foo": `+planWithOrder("foo", 1)+`, "bar": `+planWithOrder("bar", 2)+`, "zed": `+plan("zed")+`}}`,
			`{"name": "api", "plans": [`+plan("zed")+`, `+plan("bar")+`, `+plan("foo")+`]}`,
		),
	)

	// APIM renumbers published plans only, so a deprecated plan can tie with another and the GET then
	// lists it after the published plans.
	DescribeTable("no drift on the plan order when a plan is not published, whatever order APIM lists",
		func(crd, remote string) {
			expectNoDrift(detectPlanOrderDrift(crd, remote))
		},
		Entry("a deprecated plan in the middle",
			`{"name": "api", "plans": {"a": `+planWithOrder("a", 1)+`, "b": `+
				planWithOrderAndStatus("b", 2, "DEPRECATED")+`, "c": `+planWithOrder("c", 3)+`}}`,
			`{"name": "api", "plans": [`+plan("a")+`, `+plan("c")+`, `+planWithStatus("b", "DEPRECATED")+`]}`,
		),
	)

	It("drifts on the plan orders when the plans were reordered in APIM", func() {
		result := detectPlanOrderDrift(
			`{"name": "api", "plans": {"foo": `+planWithOrder("foo", 1)+`, "bar": `+planWithOrder("bar", 2)+`}}`,
			`{"name": "api", "plans": [`+plan("bar")+`, `+plan("foo")+`]}`,
		)

		Expect(result.DriftDetected()).To(BeTrue())
		Expect(result.String()).To(ContainSubstring("order: 1 != 2"))
		Expect(result.String()).To(ContainSubstring("order: 2 != 1"))
		Expect(result.String()).NotTo(ContainSubstring("name:"))
		Expect(result.String()).NotTo(ContainSubstring("hrid:"))
	})

	It("drifts on the plan orders of plans without a status, published by default", func() {
		result := detectPlanOrderDrift(
			`{"name": "api", "plans": {"foo": `+planWithOrderNoStatus("foo", 1)+`, "bar": `+
				planWithOrderNoStatus("bar", 2)+`}}`,
			`{"name": "api", "plans": [`+planNoStatus("bar")+`, `+planNoStatus("foo")+`]}`,
		)

		Expect(result.DriftDetected()).To(BeTrue())
		Expect(result.String()).To(ContainSubstring("order: 1 != 2"))
	})
})

func sentPlanHRIDs(crd string) []string {
	GinkgoHelper()
	dto := model.ToAPIV4DTO(parseJSON[*v4.Api](crd))
	hrids := make([]string, 0, len(dto.Plans))
	for _, p := range dto.Plans {
		hrids = append(hrids, p.HRID)
	}
	return hrids
}

func detectPlanOrderDrift(crd, remote string) drift.Result {
	GinkgoHelper()
	crdDTO := crdDriftAPIV4DTO(crd)
	remoteDTO := parseJSON[model.APIV4DTO](remote).Canonicalize()
	return drift.DetectWithNamespace(crdDTO, remoteDTO, "")
}

// crdDriftAPIV4DTO maps a CRD spec the way the admission drift check does.
func crdDriftAPIV4DTO(crd string) model.APIV4DTO {
	GinkgoHelper()
	return model.ToAPIV4DriftDTO(parseJSON[*v4.Api](crd))
}

func plan(hrid string) string {
	return planWithStatus(hrid, "PUBLISHED")
}

func planWithStatus(hrid, status string) string {
	return `{"hrid": "` + hrid + `", "name": "` + hrid + ` plan", "status": "` + status +
		`", "security": {"type": "KEY_LESS"}}`
}

func planWithOrder(hrid string, order int) string {
	return planWithOrderAndStatus(hrid, order, "PUBLISHED")
}

func planWithOrderAndStatus(hrid string, order int, status string) string {
	return `{"hrid": "` + hrid + `", "name": "` + hrid + ` plan", "status": "` + status + `", "order": ` +
		strconv.Itoa(order) + `, "security": {"type": "KEY_LESS"}}`
}

func planNoStatus(hrid string) string {
	return `{"hrid": "` + hrid + `", "name": "` + hrid + ` plan", "security": {"type": "KEY_LESS"}}`
}

func planWithOrderNoStatus(hrid string, order int) string {
	return `{"hrid": "` + hrid + `", "name": "` + hrid + ` plan", "order": ` + strconv.Itoa(order) +
		`, "security": {"type": "KEY_LESS"}}`
}
