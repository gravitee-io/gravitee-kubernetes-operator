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
	"encoding/json"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	nav "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/navigation"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
				Groups:           []model.APIGroup{},
				Categories:       []model.APICategory{},
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
					Groups: []model.APIGroup{},
				},
			},
		),
		Entry("differing collections equivalences",
			model.APIV4DTO{
				Groups: []model.APIGroup{
					"other",
					"developers",
				},
				ConsoleNotification: &model.APIV4ConsoleNotificationDTO{
					Groups: []model.APIGroup{
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
				Groups: []model.APIGroup{
					"gravitee-developers",
				},
				ConsoleNotification: &model.APIV4ConsoleNotificationDTO{
					Groups: []model.APIGroup{
						"gravitee-developers",
					},
				},
			},
		),
		Entry("crd-only categories ignored",
			model.APIV4DTO{
				Categories: []model.APICategory{"shared", "crd-only"},
			},
			model.APIV4DTO{
				Categories: []model.APICategory{"shared"},
			},
		),
		Entry("crd-only members ignored",
			model.APIV4DTO{
				Members: []*model.APIV4MemberDTO{
					{Source: "memory", SourceID: "admin", Role: "OWNER"},
					{Source: "memory", SourceID: "local-only", Role: "USER"},
				},
			},
			model.APIV4DTO{
				Members: []*model.APIV4MemberDTO{
					{Source: "memory", SourceID: "admin", Role: "OWNER"},
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
		Entry("unset portal navigation visibility resolved to PUBLIC matches the remote",
			model.APIV4DTO{PortalNavigation: []*model.APIV4NavigationPathDTO{
				{Path: "/alpha"},
			}}.WithResolvedVisibility(),
			model.APIV4DTO{PortalNavigation: []*model.APIV4NavigationPathDTO{
				{Path: "/alpha", Visibility: nav.Public},
			}}.WithResolvedVisibility(),
		),
		Entry("unset portal navigation visibility under a PRIVATE parent matches the inherited remote",
			model.APIV4DTO{PortalNavigation: []*model.APIV4NavigationPathDTO{
				{Path: "/alpha", Visibility: nav.Private},
				{Path: "/alpha/docs"},
			}}.WithResolvedVisibility(),
			model.APIV4DTO{PortalNavigation: []*model.APIV4NavigationPathDTO{
				{Path: "/alpha", Visibility: nav.Private},
				{Path: "/alpha/docs", Visibility: nav.Private},
			}}.WithResolvedVisibility(),
		),
		Entry("unset flow mode equivalent to the DEFAULT the remote resolves",
			model.APIV4DTO{FlowExecution: &model.APIV4FlowExecutionDTO{}},
			model.APIV4DTO{FlowExecution: &model.APIV4FlowExecutionDTO{Mode: v4.FlowModeDefault}},
		),
	)

	It("compares a flow mode the CRD sets against the remote default", func() {
		crd := model.APIV4DTO{FlowExecution: &model.APIV4FlowExecutionDTO{Mode: v4.FlowModeBestMatch}}
		remote := model.APIV4DTO{FlowExecution: &model.APIV4FlowExecutionDTO{Mode: v4.FlowModeDefault}}
		result := drift.DetectWithNamespace(crd, remote, "gravitee")
		Expect(result.DriftDetected()).To(BeTrue())
		Expect(result.String()).To(ContainSubstring(`mode: BEST_MATCH != DEFAULT`))
	})

	It("detects a console-side flip of a portal navigation entry the CRD leaves unset", func() {
		crd := model.APIV4DTO{PortalNavigation: []*model.APIV4NavigationPathDTO{
			{Path: "/alpha"},
		}}.WithResolvedVisibility()
		remote := model.APIV4DTO{PortalNavigation: []*model.APIV4NavigationPathDTO{
			{Path: "/alpha", Visibility: nav.Private},
		}}
		result := drift.DetectWithNamespace(crd, remote, "")
		Expect(result.DriftDetected()).To(BeTrue())
		Expect(result.String()).To(ContainSubstring(`visibility: PUBLIC != PRIVATE`))
	})

	It("compares native connection events as a set", func() {
		events := []v4.ConnectionEvent{v4.ConnectionEventConnected, v4.ConnectionEventDisconnected}
		crd := model.ToAPIV4DTO(&v4.Api{V4BaseApi: &v4.V4BaseApi{
			ApiBase:   &base.ApiBase{},
			Analytics: &v4.Analytics{Enabled: true, ConnectionEvents: &events},
		}})
		var remote model.APIV4DTO
		Expect(json.Unmarshal([]byte(`{"analytics": {"enabled": true, "connectionEvents": ["DISCONNECTED", "CONNECTED"]}}`), &remote)).To(Succeed())

		expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
	})

	DescribeTable("listener path overrideAccess is ignored",
		func(crdPaths, remotePaths string) {
			crd := apiV4WithListenerPaths(crdPaths).WithoutOverrideAccess()
			remote := apiV4WithListenerPaths(remotePaths).WithoutOverrideAccess()
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("unset against false",
			`[{"path": "/a"}]`,
			`[{"path": "/a", "overrideAccess": false}]`,
		),
		Entry("unset against the true a multi-tenant APIM forces on the first path",
			`[{"path": "/a"}, {"path": "/b"}]`,
			`[{"path": "/a", "overrideAccess": true}, {"path": "/b", "overrideAccess": false}]`,
		),
		Entry("unset against true on a path other than the first",
			`[{"path": "/a"}, {"path": "/b"}]`,
			`[{"path": "/a", "overrideAccess": false}, {"path": "/b", "overrideAccess": true}]`,
		),
		Entry("false against true",
			`[{"path": "/a", "overrideAccess": false}]`,
			`[{"path": "/a", "overrideAccess": true}]`,
		),
		Entry("true against false",
			`[{"path": "/a", "overrideAccess": true}]`,
			`[{"path": "/a", "overrideAccess": false}]`,
		),
	)

	It("still detects drift on the other keys of a listener path", func() {
		crd := apiV4WithListenerPaths(`[{"path": "/a", "overrideAccess": true}]`).WithoutOverrideAccess()
		remote := apiV4WithListenerPaths(`[{"path": "/b", "overrideAccess": true}]`).WithoutOverrideAccess()
		result := drift.DetectWithNamespace(crd, remote, "")
		Expect(result.DriftDetected()).To(BeTrue())
		Expect(result.String()).To(ContainSubstring(`path: "/a" != "/b"`))
	})

	DescribeTable("plans and pages are compared by HRID",
		func(crd, remote string) {
			crdDTO := model.ToAPIV4DTO(parseJSON[*v4.Api](crd))
			remoteDTO := parseJSON[model.APIV4DTO](remote).Canonicalize()
			expectNoDrift(drift.DetectWithNamespace(crdDTO, remoteDTO, ""))
		},
		Entry("plans returned out of HRID order",
			`{"name": "api", "plans": {"JWTPlan": `+apiV4JWTPlan+`, "OAuth2Plan": `+apiV4OAuth2Plan+`}}`,
			`{"name": "api", "plans": [`+apiV4OAuth2Plan+`, `+apiV4JWTPlan+`]}`,
		),
		Entry("pages returned out of HRID order",
			`{"name": "api", "pages": {"a-page": `+apiV4PageA+`, "b-page": `+apiV4PageB+`}}`,
			`{"name": "api", "pages": [`+apiV4PageB+`, `+apiV4PageA+`]}`,
		),
	)

	It("still detects a change on a plan returned out of order", func() {
		crd := model.ToAPIV4DTO(parseJSON[*v4.Api](
			`{"name": "api", "plans": {"JWTPlan": ` + apiV4JWTPlan + `, "OAuth2Plan": ` + apiV4OAuth2Plan + `}}`,
		))
		renamed := strings.Replace(apiV4OAuth2Plan, `"OAuth2 plan"`, `"OAuth2 plan v2"`, 1)
		remote := parseJSON[model.APIV4DTO](
			`{"name": "api", "plans": [` + renamed + `, ` + apiV4JWTPlan + `]}`,
		).Canonicalize()

		result := drift.DetectWithNamespace(crd, remote, "")
		Expect(result.DriftDetected()).To(BeTrue())
		Expect(result.String()).To(ContainSubstring(`plans[1]`))
		Expect(result.String()).To(ContainSubstring(`name: "OAuth2 plan" != "OAuth2 plan v2"`))
		Expect(result.String()).NotTo(ContainSubstring(`plans[0]`))
	})

	DescribeTable("still detects drift on plans and pages returned out of order",
		func(crd, remote string) {
			crdDTO := model.ToAPIV4DTO(parseJSON[*v4.Api](crd))
			remoteDTO := parseJSON[model.APIV4DTO](remote).Canonicalize()
			result := drift.DetectWithNamespace(crdDTO, remoteDTO, "")
			Expect(result.DriftDetected()).To(BeTrue())
		},
		Entry("page content changed",
			`{"name": "api", "pages": {"a-page": `+apiV4PageA+`, "b-page": `+apiV4PageB+`}}`,
			`{"name": "api", "pages": [`+apiV4PageB+`, `+strings.Replace(apiV4PageA, `"content": "a"`, `"content": "a2"`, 1)+`]}`,
		),
		Entry("plan added in APIM",
			`{"name": "api", "plans": {"JWTPlan": `+apiV4JWTPlan+`, "OAuth2Plan": `+apiV4OAuth2Plan+`}}`,
			`{"name": "api", "plans": [`+apiV4OAuth2Plan+`, `+apiV4KeylessPlan+`, `+apiV4JWTPlan+`]}`,
		),
		Entry("plan removed in APIM",
			`{"name": "api", "plans": {"JWTPlan": `+apiV4JWTPlan+`, "OAuth2Plan": `+apiV4OAuth2Plan+`}}`,
			`{"name": "api", "plans": [`+apiV4OAuth2Plan+`]}`,
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

func apiV4WithListenerPaths(paths string) model.APIV4DTO {
	GinkgoHelper()
	return parseJSON[model.APIV4DTO](`{"listeners": [{"type": "HTTP", "paths": ` + paths + `}]}`)
}

func completeAPIV4DTO() model.APIV4DTO {
	GinkgoHelper()
	return loadFixture[model.APIV4DTO]("api_v4_dto.json")
}

const (
	apiV4JWTPlan     = `{"hrid": "JWTPlan", "name": "JWT plan", "status": "PUBLISHED", "security": {"type": "JWT"}}`
	apiV4OAuth2Plan  = `{"hrid": "OAuth2Plan", "name": "OAuth2 plan", "status": "PUBLISHED", "security": {"type": "OAUTH2"}}`
	apiV4KeylessPlan = `{"hrid": "KeylessPlan", "name": "Keyless plan", "status": "PUBLISHED", "security": {"type": "KEY_LESS"}}`
	apiV4PageA       = `{"hrid": "a-page", "name": "A", "type": "MARKDOWN", "content": "a"}`
	apiV4PageB       = `{"hrid": "b-page", "name": "B", "type": "MARKDOWN", "content": "b"}`
)
