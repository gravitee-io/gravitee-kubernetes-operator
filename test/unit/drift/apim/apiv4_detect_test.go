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
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("API v4 Drift detection", func() {
	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, "gravitee"))
		},
		Entry("empty struct",
			model.APIV4DTO{},
			model.APIV4DTO{},
		),
		Entry("equal struct",
			completeAPIV4DTO(),
			completeAPIV4DTO(),
		),
		Entry("equal struct from CRD mapping",
			completeAPIV4DTO(),
			model.ToAPIV4DTO(completeAPIV4Spec()),
		),
		Entry("empty collections equivalent to nil",
			model.APIV4DTO{
				FlowExecution: &model.APIV4FlowExecutionDTO{},
				Plans: []*model.APIV4PlanDTO{{
					Security: &model.APIV4PlanSecurityDTO{
						Type: "foo",
					},
				}},
			},
			model.APIV4DTO{
				Tags:       []string{},
				Labels:     []string{},
				Properties: []*model.APIV4PropertyDTO{},
				Resources:  []*model.APIV4ResourceDTO{},
				Metadata:   []*model.APIV4MetadataEntryDTO{},
				FlowExecution: &model.APIV4FlowExecutionDTO{
					Mode: "DEFAULT",
				},
				Groups:           []string{},
				Categories:       []string{},
				Flows:            []*model.APIV4FlowDTO{},
				Members:          []*model.APIV4MemberDTO{},
				PortalNavigation: []*model.APIV4NavigationPathDTO{},
				Pages:            []*model.APIV4PageDTO{},
				Plans: []*model.APIV4PlanDTO{
					{
						Characteristics: make([]string, 0),
						Security: &model.APIV4PlanSecurityDTO{
							Type: "FOO",
						},
					},
				},
				ConsoleNotification: &model.APIV4ConsoleNotificationDTO{
					Events: []string{},
					Groups: []string{},
				},
			},
		),
		Entry("differing collections equivalences",
			model.APIV4DTO{
				Groups: []string{
					"other",
					"developers",
				},
				ConsoleNotification: &model.APIV4ConsoleNotificationDTO{
					Groups: []string{
						"other",
						"developers",
					},
				},
			},
			model.APIV4DTO{
				Metadata: []*model.APIV4MetadataEntryDTO{
					{
						BaseMetadata: model.BaseMetadata{Name: "foo", Value: nil, DefaultValue: new("bar")},
						Key:          "foo",
						Format:       "STRING",
					},
				},
				Groups: []string{
					"gravitee-developers",
				},
				ConsoleNotification: &model.APIV4ConsoleNotificationDTO{
					Groups: []string{
						"gravitee-developers",
					},
				},
			},
		),
		Entry("empty reporters equivalent to true",
			model.APIV4DTO{
				Analytics: nil,
			},
			model.APIV4DTO{
				Analytics: &model.APIV4AnalyticsDTO{
					Enabled:                true,
					ReporterMetricsEnabled: new(true),
				},
			},
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeAPIV4DTO(), completeAPIV4DTO())
		})
	})
})

func completeAPIV4Spec() *v4.Api {
	GinkgoHelper()
	fixture := loadFixture[v4.Api]("api_v4_spec.json")
	return &fixture
}

func completeAPIV4DTO() model.APIV4DTO {
	GinkgoHelper()
	return loadFixture[model.APIV4DTO]("api_v4_dto.json")
}
