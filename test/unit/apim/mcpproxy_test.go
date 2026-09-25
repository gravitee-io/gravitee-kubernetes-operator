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
	"fmt"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	xErrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hrid"
)

func newMcpProxy(auth *mcpproxy.UpstreamAuth) *v1alpha1.McpProxy {
	return &v1alpha1.McpProxy{
		ObjectMeta: metav1.ObjectMeta{Name: "github-mcp", Namespace: "gravitee"},
		Spec: v1alpha1.McpProxySpec{
			Context: &refs.NamespacedName{Name: "dev-ctx"},
			Type: mcpproxy.Type{
				EntityID:        "mcp-proxy.github",
				Name:            "GitHub",
				ContextPath:     "/mcp/github",
				ProtocolVersion: "2025-03-26",
				Mode:            mcpproxy.ModeProxy,
				State:           mcpproxy.StateStarted,
				Proxy: &mcpproxy.Proxy{
					ServerURL:    "https://api.githubcopilot.com/mcp/",
					UpstreamAuth: auth,
				},
				Plans: []mcpproxy.Plan{{
					Name: "gold",
					Security: mcpproxy.PlanSecurity{
						Type:   mcpproxy.PlanSecurityAPIKey,
						APIKey: &mcpproxy.PlanAPIKey{Source: "HEADER"},
					},
				}},
			},
		},
	}
}

func newMcpStudio() *v1alpha1.McpProxy {
	return &v1alpha1.McpProxy{
		ObjectMeta: metav1.ObjectMeta{Name: "support-studio", Namespace: "gravitee"},
		Spec: v1alpha1.McpProxySpec{
			Context: &refs.NamespacedName{Name: "dev-ctx"},
			Type: mcpproxy.Type{
				EntityID:        "mcp-proxy.support",
				Name:            "Support Studio",
				ContextPath:     "/mcp/support",
				ProtocolVersion: "2025-03-26",
				Mode:            mcpproxy.ModeStudio,
				State:           mcpproxy.StateStarted,
				Studio: &mcpproxy.Studio{
					Tools: []mcpproxy.StudioTool{
						{ServerRef: refs.NamespacedName{Name: "jira-mcp", Namespace: "tools"}, Tool: "create_ticket"},
						{ServerRef: refs.NamespacedName{Name: "github-mcp"}, Tool: "search_issues", Alias: new("gh_search")},
						{ServerRef: refs.NamespacedName{Name: "github-mcp"}, Tool: "create_issue"},
					},
					UpstreamAuth: []mcpproxy.StudioUpstreamAuth{
						{
							ServerRef: refs.NamespacedName{Name: "jira-mcp", Namespace: "tools"},
							Auth:      mcpproxy.UpstreamAuth{Type: mcpproxy.UpstreamAuthTypeNone},
						},
						{
							ServerRef: refs.NamespacedName{Name: "github-mcp"},
							Auth: mcpproxy.UpstreamAuth{
								Type:   mcpproxy.UpstreamAuthTypeBearer,
								Bearer: &mcpproxy.BearerAuth{Token: "secret://kubernetes/github-mcp:token"},
							},
						},
					},
					EnableFGA: true,
				},
				Plans: []mcpproxy.Plan{{
					Name: "gold",
					Security: mcpproxy.PlanSecurity{
						Type:   mcpproxy.PlanSecurityAPIKey,
						APIKey: &mcpproxy.PlanAPIKey{Source: "HEADER"},
					},
				}},
			},
		},
	}
}

func marshal(v any) string {
	raw, err := json.Marshal(v)
	Expect(err).ToNot(HaveOccurred())
	return string(raw)
}

var _ = Describe("MCP proxy", func() {
	Describe("spec hash", func() {
		It("changes with the context path and the upstream", func() {
			base := newMcpProxy(nil)
			otherPath := newMcpProxy(nil)
			otherPath.Spec.ContextPath = "/mcp/other"
			otherURL := newMcpProxy(nil)
			otherURL.Spec.Proxy.ServerURL = "https://mcp.gitlab.example.com/"

			Expect(base.Spec.Hash()).To(Equal(newMcpProxy(nil).Spec.Hash()))
			Expect(base.Spec.Hash()).ToNot(Equal(otherPath.Spec.Hash()))
			Expect(base.Spec.Hash()).ToNot(Equal(otherURL.Spec.Hash()))
		})
	})

	Describe("PROXY wire mapping", func() {
		It("derives the hrid and moves the proxy block to the top level", func() {
			dto := model.ToMcpProxyDTO(newMcpProxy(nil))
			Expect(dto.HRID).To(Equal("gravitee-github-mcp"))
			Expect(*dto.ServerURL).To(Equal("https://api.githubcopilot.com/mcp/"))
			Expect(dto.Studio).To(BeNil())
		})

		It("defaults the mode and the lifecycle", func() {
			proxy := newMcpProxy(nil)
			proxy.Spec.Mode = ""
			proxy.Spec.State = ""
			dto := model.ToMcpProxyDTO(proxy)
			Expect(dto.Mode).To(Equal("PROXY"))
			Expect(dto.State).To(Equal("STARTED"))
		})

		DescribeTable("sends no upstream auth for a passthrough proxy",
			func(auth *mcpproxy.UpstreamAuth) {
				Expect(model.ToMcpProxyDTO(newMcpProxy(auth)).UpstreamAuth).To(BeNil())
			},
			Entry("when omitted", nil),
			Entry("when NONE", &mcpproxy.UpstreamAuth{Type: mcpproxy.UpstreamAuthTypeNone}),
		)

		DescribeTable("flattens each upstream auth variant",
			func(auth mcpproxy.UpstreamAuth, expected string) {
				Expect(marshal(model.ToMcpProxyDTO(newMcpProxy(&auth)).UpstreamAuth)).To(MatchJSON(expected))
			},
			Entry("API_KEY", mcpproxy.UpstreamAuth{
				Type:   mcpproxy.UpstreamAuthTypeAPIKey,
				APIKey: &mcpproxy.APIKeyAuth{Header: "X-Api-Key", Value: "k"},
			}, `{"type": "API_KEY", "apiKeyHeader": "X-Api-Key", "apiKey": "k"}`),
			Entry("BEARER", mcpproxy.UpstreamAuth{
				Type:   mcpproxy.UpstreamAuthTypeBearer,
				Bearer: &mcpproxy.BearerAuth{Token: "t"},
			}, `{"type": "BEARER", "token": "t"}`),
			Entry("BASIC", mcpproxy.UpstreamAuth{
				Type:  mcpproxy.UpstreamAuthTypeBasic,
				Basic: &mcpproxy.BasicAuth{Username: "u", Password: "p"},
			}, `{"type": "BASIC", "username": "u", "password": "p"}`),
			Entry("OAUTH2", mcpproxy.UpstreamAuth{
				Type: mcpproxy.UpstreamAuthTypeOAuth2,
				OAuth2: &mcpproxy.OAuth2Auth{
					AuthorizeURL: "https://auth.example.com/authorize",
					TokenURL:     "https://auth.example.com/token",
					ClientID:     "gravitee",
					ClientSecret: "s",
					Scopes:       []string{"repo"},
				},
			}, `{
				"type": "OAUTH2", "authorizeUrl": "https://auth.example.com/authorize",
				"tokenUrl": "https://auth.example.com/token", "clientId": "gravitee",
				"clientSecret": "s", "scopes": ["repo"]
			}`),
		)

		DescribeTable("flattens each plan security variant",
			func(security mcpproxy.PlanSecurity, expected string) {
				proxy := newMcpProxy(nil)
				proxy.Spec.Plans[0].Security = security
				Expect(marshal(model.ToMcpProxyDTO(proxy).Plans[0].Security)).To(MatchJSON(expected))
			},
			Entry("KEY_LESS", mcpproxy.PlanSecurity{Type: mcpproxy.PlanSecurityKeyLess}, `{"type": "KEY_LESS"}`),
			Entry("API_KEY without propagation", mcpproxy.PlanSecurity{
				Type:   mcpproxy.PlanSecurityAPIKey,
				APIKey: &mcpproxy.PlanAPIKey{Source: "HEADER", Header: new("X-Key")},
			}, `{"type": "API_KEY", "source": "HEADER", "apiKeyHeader": "X-Key"}`),
			Entry("API_KEY with propagation", mcpproxy.PlanSecurity{
				Type:   mcpproxy.PlanSecurityAPIKey,
				APIKey: &mcpproxy.PlanAPIKey{Source: "BEARER", Propagate: true},
			}, `{"type": "API_KEY", "source": "BEARER", "propagateApiKey": true}`),
			Entry("OAUTH2", mcpproxy.PlanSecurity{
				Type:   mcpproxy.PlanSecurityOAuth2,
				OAuth2: &mcpproxy.PlanOAuth2{Provider: "am", Scopes: []string{"repo:read"}},
			}, `{"type": "OAUTH2", "provider": "am", "scopes": ["repo:read"]}`),
		)

		It("drops the platform default flow execution", func() {
			proxy := newMcpProxy(nil)
			proxy.Spec.FlowExecution = &v4.FlowExecution{Mode: v4.FlowModeDefault}
			Expect(model.ToMcpProxyDTO(proxy).FlowExecution).To(BeNil())

			proxy.Spec.FlowExecution = &v4.FlowExecution{Mode: v4.FlowModeBestMatch}
			Expect(model.ToMcpProxyDTO(proxy).FlowExecution).To(Equal(&v4.FlowExecution{Mode: v4.FlowModeBestMatch}))
		})

		It("serialises the proxyBehindKeycloak shape of the automation OAS with the wire field names", func() {
			proxy := newMcpProxy(nil)
			proxy.Spec.IdentityProviders = []mcpproxy.IdentityProvider{
				{
					Name: "keycloak",
					Type: mcpproxy.IdentityProviderOAuth2Generic,
					OAuth2Generic: &mcpproxy.OAuth2GenericProvider{
						IssuerURL:                   "https://kc.example.com/realms/acme",
						IntrospectionEndpoint:       "/protocol/openid-connect/token/introspect",
						IntrospectionEndpointMethod: "POST",
						ClientID:                    new("gateway"),
						ClientSecret:                new("secret://kubernetes/keycloak-gateway:client-secret"),
					},
				},
				{Name: "am", Type: mcpproxy.IdentityProviderGraviteeAM},
			}
			proxy.Spec.Flows = []mcpproxy.Flow{{
				Name:    "Tools ACL",
				Enabled: true,
				Selectors: []mcpproxy.FlowSelector{{
					Type: mcpproxy.FlowSelectorMCP,
					MCP:  &mcpproxy.McpSelector{Methods: []string{"tools/list", "tools/call"}},
				}},
				Request: []mcpproxy.FlowStep{{
					Name:    "MCP ACL",
					Policy:  "mcp-acl",
					Enabled: true,
					Configuration: utils.NewGenericStringMap().Put("authorizations", []any{
						map[string]any{"patternType": "ANY"},
					}),
				}},
			}}
			proxy.Spec.Plans = []mcpproxy.Plan{
				{
					Name: "partners",
					Security: mcpproxy.PlanSecurity{
						Type:   mcpproxy.PlanSecurityOAuth2,
						OAuth2: &mcpproxy.PlanOAuth2{Provider: "keycloak"},
					},
				},
				{
					Name: "gold",
					Security: mcpproxy.PlanSecurity{
						Type:   mcpproxy.PlanSecurityOAuth2,
						OAuth2: &mcpproxy.PlanOAuth2{Provider: "am", Scopes: []string{"repo:read"}},
					},
					Flows: []mcpproxy.Flow{{
						Name:    "Gold rate limit",
						Enabled: true,
						Selectors: []mcpproxy.FlowSelector{{
							Type:      mcpproxy.FlowSelectorCondition,
							Condition: &mcpproxy.ConditionSelector{Condition: "{#request.headers['X-Gold'] != null}"},
						}},
					}},
				},
			}
			proxy.Spec.Proxy.UpstreamAuth = &mcpproxy.UpstreamAuth{
				Type:   mcpproxy.UpstreamAuthTypeBearer,
				Bearer: &mcpproxy.BearerAuth{Token: "secret://kubernetes/github-mcp:token"},
			}

			Expect(marshal(model.ToMcpProxyDTO(proxy))).To(MatchJSON(`{
				"hrid": "gravitee-github-mcp",
				"entityId": "mcp-proxy.github",
				"name": "GitHub",
				"contextPath": "/mcp/github",
				"mode": "PROXY",
				"protocolVersion": "2025-03-26",
				"state": "STARTED",
				"serverUrl": "https://api.githubcopilot.com/mcp/",
				"upstreamAuth": {"type": "BEARER", "token": "secret://kubernetes/github-mcp:token"},
				"identityProviders": [
					{"name": "am", "type": "GRAVITEE_AM"},
					{
						"name": "keycloak", "type": "OAUTH2_GENERIC",
						"issuerUrl": "https://kc.example.com/realms/acme",
						"introspectionEndpoint": "/protocol/openid-connect/token/introspect",
						"introspectionEndpointMethod": "POST",
						"clientId": "gateway",
						"clientSecret": "secret://kubernetes/keycloak-gateway:client-secret"
					}
				],
				"flows": [{
					"name": "Tools ACL", "enabled": true,
					"selectors": [{"type": "MCP", "methods": ["tools/list", "tools/call"]}],
					"request": [{
						"name": "MCP ACL", "policy": "mcp-acl", "enabled": true,
						"configuration": {"authorizations": [{"patternType": "ANY"}]}
					}]
				}],
				"plans": [
					{
						"name": "gold",
						"security": {"type": "OAUTH2", "provider": "am", "scopes": ["repo:read"]},
						"flows": [{
							"name": "Gold rate limit", "enabled": true,
							"selectors": [{"type": "CONDITION", "condition": "{#request.headers['X-Gold'] != null}"}]
						}]
					},
					{"name": "partners", "security": {"type": "OAUTH2", "provider": "keycloak"}}
				]
			}`))
		})
	})

	Describe("STUDIO wire mapping", func() {
		It("serialises the studio shape of the automation OAS with the wire field names", func() {
			Expect(marshal(model.ToMcpProxyDTO(newMcpStudio()))).To(MatchJSON(`{
				"hrid": "gravitee-support-studio",
				"entityId": "mcp-proxy.support",
				"name": "Support Studio",
				"contextPath": "/mcp/support",
				"protocolVersion": "2025-03-26",
				"mode": "STUDIO",
				"state": "STARTED",
				"studio": {
					"tools": [
						{"server": "gravitee-github-mcp", "tool": "create_issue"},
						{"server": "gravitee-github-mcp", "tool": "search_issues", "alias": "gh_search"},
						{"server": "tools-jira-mcp", "tool": "create_ticket"}
					],
					"upstreamAuth": [
						{"server": "gravitee-github-mcp", "auth": {"type": "BEARER", "token": "secret://kubernetes/github-mcp:token"}},
						{"server": "tools-jira-mcp", "auth": {"type": "NONE"}}
					],
					"enableFGA": true
				},
				"plans": [{"name": "gold", "security": {"type": "API_KEY", "source": "HEADER"}}]
			}`))
		})
	})

	Describe("studio references", func() {
		It("lists every referenced server once, with the namespace defaulted", func() {
			servers := newMcpStudio().Spec.ServerRefs("gravitee")
			Expect(servers).To(Equal([]refs.NamespacedName{
				{Namespace: "tools", Name: "jira-mcp"},
				{Namespace: "gravitee", Name: "github-mcp"},
			}))
		})

		It("finds no server without upstream auth when each has an entry", func() {
			Expect(newMcpStudio().Spec.ServersWithoutUpstreamAuth("gravitee")).To(BeEmpty())
		})

		It("names a server whose tools are selected without an upstream auth entry", func() {
			studio := newMcpStudio()
			studio.Spec.Studio.UpstreamAuth = studio.Spec.Studio.UpstreamAuth[:1]
			Expect(studio.Spec.ServersWithoutUpstreamAuth("gravitee")).To(Equal([]refs.NamespacedName{
				{Namespace: "gravitee", Name: "github-mcp"},
			}))
		})

		It("matches an entry naming the proxy namespace explicitly", func() {
			studio := newMcpStudio()
			studio.Spec.Studio.UpstreamAuth[1].ServerRef.Namespace = "gravitee"
			Expect(studio.Spec.ServersWithoutUpstreamAuth("gravitee")).To(BeEmpty())
		})

		It("does not match a server of the same name in another namespace", func() {
			studio := newMcpStudio()
			studio.Spec.Studio.UpstreamAuth[1].ServerRef.Namespace = "elsewhere"
			Expect(studio.Spec.ServersWithoutUpstreamAuth("gravitee")).To(Equal([]refs.NamespacedName{
				{Namespace: "gravitee", Name: "github-mcp"},
			}))
		})
	})

	Describe("state decoding", func() {
		It("decodes the ids, the observed lifecycle and the studio tool entity ids", func() {
			body := `{
				"hrid": "gravitee-support-studio",
				"entityId": "mcp-proxy.support",
				"name": "Support Studio",
				"contextPath": "/mcp/support",
				"mode": "STUDIO",
				"protocolVersion": "2025-03-26",
				"state": "STOPPED",
				"studio": {
					"tools": [{"server": "gravitee-github-mcp", "tool": "create_issue", "entityId": "mcp-tool.github.create_issue"}],
					"upstreamAuth": [{"server": "gravitee-github-mcp", "auth": {"type": "BEARER"}}],
					"enableFGA": true
				},
				"plans": [{"name": "gold", "security": {"type": "API_KEY", "source": "HEADER"}}],
				"id": "1f3c", "organizationId": "DEFAULT", "environmentId": "DEFAULT",
				"errors": {"warning": ["tool added after creation"]}
			}`
			state := new(model.McpProxyState)
			Expect(json.Unmarshal([]byte(body), state)).To(Succeed())

			Expect(state.HRID).To(Equal("gravitee-support-studio"))
			Expect(state.ID).To(Equal("1f3c"))
			Expect(state.EnvID).To(Equal("DEFAULT"))
			Expect(state.State).To(Equal("STOPPED"))
			Expect(state.Errors.Warning).To(ConsistOf("tool added after creation"))
			Expect(state.Studio.Tools[0].EntityID).To(Equal("mcp-tool.github.create_issue"))
			Expect(state.Studio.UpstreamAuth[0].Auth.Token).To(BeNil())
		})
	})

	Describe("refusal findings", func() {
		refusedWith := func(statusCode int, body string) error {
			return xErrors.ServerError{StatusCode: statusCode, Body: body}
		}

		It("reads the findings a refused apply carries in errors.severe", func() {
			findings, refused := model.AutomationRefusal(refusedWith(http.StatusBadRequest, `{
				"name": "GitHub",
				"errors": {"severe": ["plan [gold] cannot be changed in place"], "warning": ["slow"]}
			}`))

			Expect(refused).To(BeTrue())
			Expect(findings.Severe).To(ConsistOf("plan [gold] cannot be changed in place"))
			Expect(findings.Warning).To(ConsistOf("slow"))
		})

		DescribeTable("is not a refusal",
			func(err error) {
				_, refused := model.AutomationRefusal(err)
				Expect(refused).To(BeFalse())
			},
			Entry("a 400 without severe findings",
				refusedWith(http.StatusBadRequest, `{"errors": {"warning": ["slow"]}}`)),
			Entry("a lifecycle failure in the platform error shape",
				refusedWith(http.StatusInternalServerError,
					`{"httpStatus": 500, "technicalCode": "gamma.mcp.lifecycle.failed"}`)),
			Entry("an error that is not a server answer",
				fmt.Errorf("connection refused")),
		)
	})

	Describe("derived hrid", func() {
		DescribeTable("is refused when the automation API cannot address it",
			func(value string) {
				Expect(hrid.Validate(value)).To(MatchError(ContainSubstring(value)))
			},
			Entry("a dotted name", "gravitee-github.mcp"),
			Entry("longer than 256 characters", "gravitee-"+strings.Repeat("a", 248)),
		)

		It("is accepted when it matches the grammar", func() {
			Expect(hrid.Validate("gravitee-github_mcp")).To(Succeed())
		})
	})
})
