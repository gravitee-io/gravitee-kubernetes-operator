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
	"fmt"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
	apimmcpproxy "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	xErrors "github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hrid"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

func newMcpProxy(auth *mcpproxy.UpstreamAuth) *v1alpha1.McpProxy {
	return &v1alpha1.McpProxy{
		ObjectMeta: metav1.ObjectMeta{Name: "github-mcp", Namespace: "gravitee"},
		Spec: v1alpha1.McpProxySpec{
			Context: &refs.NamespacedName{Name: "dev-ctx"},
			Type: mcpproxy.Type{
				EntityID:    "mcp-proxy.github",
				Name:        "GitHub",
				ContextPath: "/mcp/github",
				Mode:        mcpproxy.ModeProxy,
				State:       mcpproxy.StateStarted,
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
				EntityID:    "mcp-proxy.support",
				Name:        "Support Studio",
				ContextPath: "/mcp/support",
				Mode:        mcpproxy.ModeStudio,
				State:       mcpproxy.StateStarted,
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

func toMcpProxyDTO(proxy *v1alpha1.McpProxy) sdk.AimMcpProxy {
	GinkgoHelper()
	dto, err := apimmcpproxy.ToDTO(proxy)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func decodeMcpProxyState(body string) sdk.AimMcpProxyState {
	GinkgoHelper()
	state := sdk.AimMcpProxyState{}
	Expect(json.Unmarshal([]byte(body), &state)).To(Succeed())
	return state
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
		It("derives the hrid and sends the proxy block as declared", func() {
			dto := toMcpProxyDTO(newMcpProxy(nil))
			Expect(dto.Hrid).To(Equal("gravitee-github-mcp"))
			Expect(dto.Proxy.ServerUrl).To(Equal("https://api.githubcopilot.com/mcp/"))
			Expect(dto.Studio).To(BeNil())
		})

		It("defaults the lifecycle", func() {
			proxy := newMcpProxy(nil)
			proxy.Spec.State = ""
			Expect(*toMcpProxyDTO(proxy).State).To(Equal(sdk.AimMcpProxyLifecycleStateSTARTED))
		})

		DescribeTable("sends no upstream auth for a passthrough proxy",
			func(auth *mcpproxy.UpstreamAuth) {
				Expect(toMcpProxyDTO(newMcpProxy(auth)).Proxy.UpstreamAuth).To(BeNil())
			},
			Entry("when omitted", nil),
			Entry("when NONE", &mcpproxy.UpstreamAuth{Type: mcpproxy.UpstreamAuthTypeNone}),
		)

		DescribeTable("flattens each upstream auth variant",
			func(auth mcpproxy.UpstreamAuth, expected string) {
				Expect(marshal(toMcpProxyDTO(newMcpProxy(&auth)).Proxy.UpstreamAuth)).To(MatchJSON(expected))
			},
			Entry("API_KEY", mcpproxy.UpstreamAuth{
				Type:   mcpproxy.UpstreamAuthTypeAPIKey,
				APIKey: &mcpproxy.APIKeyAuth{APIKeyHeader: "X-Api-Key", APIKey: "k"},
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
				Expect(marshal(toMcpProxyDTO(proxy).Plans[0].Security)).To(MatchJSON(expected))
			},
			Entry("KEY_LESS", mcpproxy.PlanSecurity{Type: mcpproxy.PlanSecurityKeyLess}, `{"type": "KEY_LESS"}`),
			Entry("API_KEY without propagation", mcpproxy.PlanSecurity{
				Type:   mcpproxy.PlanSecurityAPIKey,
				APIKey: &mcpproxy.PlanAPIKey{Source: "HEADER", APIKeyHeader: new("X-Key")},
			}, `{"type": "API_KEY", "source": "HEADER", "apiKeyHeader": "X-Key"}`),
			Entry("API_KEY with propagation", mcpproxy.PlanSecurity{
				Type:   mcpproxy.PlanSecurityAPIKey,
				APIKey: &mcpproxy.PlanAPIKey{Source: "BEARER", PropagateAPIKey: true},
			}, `{"type": "API_KEY", "source": "BEARER", "propagateApiKey": true}`),
			Entry("OAUTH2", mcpproxy.PlanSecurity{
				Type:   mcpproxy.PlanSecurityOAuth2,
				OAuth2: &mcpproxy.PlanOAuth2{Provider: "am", Scopes: []string{"repo:read"}},
			}, `{"type": "OAUTH2", "provider": "am", "scopes": ["repo:read"]}`),
		)

		It("drops the platform default flow execution", func() {
			proxy := newMcpProxy(nil)
			proxy.Spec.FlowExecution = &v4.FlowExecution{Mode: v4.FlowModeDefault}
			Expect(toMcpProxyDTO(proxy).FlowExecution).To(BeNil())

			proxy.Spec.FlowExecution = &v4.FlowExecution{Mode: v4.FlowModeBestMatch}
			Expect(marshal(toMcpProxyDTO(proxy).FlowExecution)).To(MatchJSON(`{"mode": "BEST_MATCH", "matchRequired": false}`))

			proxy.Spec.FlowExecution = &v4.FlowExecution{MatchRequired: true}
			Expect(marshal(toMcpProxyDTO(proxy).FlowExecution)).To(MatchJSON(`{"matchRequired": true}`))
		})

		DescribeTable("refuses a union it cannot send",
			func(mutate func(*v1alpha1.McpProxy)) {
				proxy := newMcpProxy(nil)
				mutate(proxy)
				_, err := apimmcpproxy.ToDTO(proxy)
				Expect(err).To(HaveOccurred())
			},
			Entry("an upstream auth without its block", func(proxy *v1alpha1.McpProxy) {
				proxy.Spec.Proxy.UpstreamAuth = &mcpproxy.UpstreamAuth{Type: mcpproxy.UpstreamAuthTypeBearer}
			}),
			Entry("an unknown upstream auth type", func(proxy *v1alpha1.McpProxy) {
				proxy.Spec.Proxy.UpstreamAuth = &mcpproxy.UpstreamAuth{Type: "DIGEST"}
			}),
			Entry("a plan security without its block", func(proxy *v1alpha1.McpProxy) {
				proxy.Spec.Plans[0].Security = mcpproxy.PlanSecurity{Type: mcpproxy.PlanSecurityOAuth2}
			}),
			Entry("an identity provider without its block", func(proxy *v1alpha1.McpProxy) {
				proxy.Spec.IdentityProviders = []mcpproxy.IdentityProvider{
					{Name: "auth0", Type: mcpproxy.IdentityProviderOAuth2Auth0},
				}
			}),
			Entry("a flow selector without its block", func(proxy *v1alpha1.McpProxy) {
				proxy.Spec.Flows = []mcpproxy.Flow{{
					Name:      "acl",
					Selectors: []mcpproxy.FlowSelector{{Type: mcpproxy.FlowSelectorMCP}},
				}}
			}),
		)

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

			Expect(marshal(toMcpProxyDTO(proxy))).To(MatchJSON(`{
				"hrid": "gravitee-github-mcp",
				"entityId": "mcp-proxy.github",
				"name": "GitHub",
				"contextPath": "/mcp/github",
				"mode": "PROXY",
				"state": "STARTED",
				"proxy": {
					"serverUrl": "https://api.githubcopilot.com/mcp/",
					"upstreamAuth": {"type": "BEARER", "token": "secret://kubernetes/github-mcp:token"}
				},
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
			Expect(marshal(toMcpProxyDTO(newMcpStudio()))).To(MatchJSON(`{
				"hrid": "gravitee-support-studio",
				"entityId": "mcp-proxy.support",
				"name": "Support Studio",
				"contextPath": "/mcp/support",
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
			state := decodeMcpProxyState(`{
				"hrid": "gravitee-support-studio",
				"entityId": "mcp-proxy.support",
				"name": "Support Studio",
				"contextPath": "/mcp/support",
				"mode": "STUDIO",
				"state": "STOPPED",
				"studio": {
					"tools": [{"server": "gravitee-github-mcp", "tool": "create_issue", "entityId": "mcp-tool.github.create_issue"}],
					"upstreamAuth": [{"server": "gravitee-github-mcp", "auth": {"type": "BEARER"}}],
					"enableFGA": true
				},
				"plans": [{"name": "gold", "security": {"type": "API_KEY", "source": "HEADER"}}],
				"id": "1f3c", "organizationId": "DEFAULT", "environmentId": "DEFAULT",
				"errors": {"warning": ["tool added after creation"]},
				"tools": [{"server": "gravitee-github-mcp", "tool": "create_issue", "alias": "open_ticket",
					"entityId": "mcp-tool.github.create_issue", "authPolicyRef": "MCPTool::\"github.create_issue\""}]
			}`)

			Expect(state.Hrid).To(Equal("gravitee-support-studio"))
			Expect(*state.Id).To(Equal("1f3c"))
			Expect(*state.EnvironmentId).To(Equal("DEFAULT"))
			Expect(*state.State).To(Equal(sdk.AimMcpProxyLifecycleStateSTOPPED))
			Expect(state.Errors.Warning).To(ConsistOf("tool added after creation"))
			Expect(*state.Studio.Tools[0].EntityId).To(Equal("mcp-tool.github.create_issue"))
			Expect(state.Studio.UpstreamAuth[0].Auth.Bearer.Token).To(BeNil())
		})

		// The status list, not studio.tools[], carries the policy reference: studio.tools[] is the
		// payload drift compares.
		It("leaves the tools status list out of the payload drift compares", func() {
			state := decodeMcpProxyState(`{
				"hrid": "gravitee-support-studio",
				"studio": {"tools": [{"server": "gravitee-github-mcp", "tool": "create_issue", "alias": "open_ticket"}]},
				"tools": [{"server": "gravitee-github-mcp", "tool": "create_issue", "alias": "open_ticket",
					"entityId": "mcp-tool.github.create_issue", "authPolicyRef": "MCPTool::\"github.create_issue\""}]
			}`)

			Expect(state.Tools).To(HaveLen(1))
			Expect(*state.Tools[0].AuthPolicyRef).To(Equal(`MCPTool::"github.create_issue"`))
			Expect(marshal(state.Spec())).ToNot(ContainSubstring("authPolicyRef"))
		})
	})

	// What the controller writes in status.tools, from the state a PUT answers.
	Describe("tools status", func() {
		studioTools := `"studio": {"tools": [{"server": "gravitee-github-mcp", "tool": "search_issues",
			"alias": "gh_search", "entityId": "mcp-tool.github.search_issues"}]}`
		fromStudio := []mcpproxy.StudioToolStatus{{
			Server:   "gravitee-github-mcp",
			Tool:     "search_issues",
			Alias:    new("gh_search"),
			EntityID: "mcp-tool.github.search_issues",
		}}

		DescribeTable("maps the state onto status.tools",
			func(body string, expected []mcpproxy.StudioToolStatus) {
				Expect(apimmcpproxy.ToolsStatus(decodeMcpProxyState(body))).To(Equal(expected))
			},
			Entry("from the tools list when the state answers both lists, references included",
				`{"mode": "STUDIO", `+studioTools+`, "tools": [{"server": "gravitee-github-mcp",
					"tool": "search_issues", "alias": "gh_search", "entityId": "mcp-tool.github.search_issues",
					"authPolicyRef": "MCPTool::\"github.search_issues\""}]}`,
				[]mcpproxy.StudioToolStatus{{
					Server:        "gravitee-github-mcp",
					Tool:          "search_issues",
					Alias:         new("gh_search"),
					EntityID:      "mcp-tool.github.search_issues",
					AuthPolicyRef: `MCPTool::"github.search_issues"`,
				}}),
			Entry("from studio.tools, with no reference, when the state has no tools list",
				`{"mode": "STUDIO", `+studioTools+`}`,
				fromStudio),
			Entry("as none when the state answers an empty tools list",
				`{"mode": "STUDIO", `+studioTools+`, "tools": []}`,
				nil),
			Entry("as none for a PROXY",
				`{"mode": "PROXY", "proxy": {"serverUrl": "https://api.githubcopilot.com/mcp/"}}`,
				nil),
		)
	})

	Describe("status", func() {
		It("copies the ids, the findings, the observed lifecycle and the tools", func() {
			state := decodeMcpProxyState(`{
				"hrid": "gravitee-support-studio", "mode": "STUDIO", "state": "STOPPED",
				"id": "1f3c", "organizationId": "DEFAULT", "environmentId": "DEFAULT",
				"errors": {"warning": ["tool added after creation"]},
				"tools": [{"server": "gravitee-github-mcp", "tool": "create_issue",
					"entityId": "mcp-tool.github.create_issue", "authPolicyRef": "MCPTool::\"github.create_issue\""}]
			}`)
			proxy := newMcpStudio()

			Expect(apimmcpproxy.UpdateStatus(context.Background(), proxy, apimmcpproxy.Response{State: state})).To(Succeed())

			Expect(proxy.Status.ID).To(Equal("1f3c"))
			Expect(proxy.Status.OrgID).To(Equal("DEFAULT"))
			Expect(proxy.Status.EnvID).To(Equal("DEFAULT"))
			Expect(proxy.Status.HRID).To(Equal("gravitee-support-studio"))
			Expect(proxy.Status.State).To(Equal("STOPPED"))
			Expect(proxy.Status.Errors.Warning).To(ConsistOf("tool added after creation"))
			Expect(proxy.Status.Tools).To(Equal([]mcpproxy.StudioToolStatus{{
				Server:        "gravitee-github-mcp",
				Tool:          "create_issue",
				EntityID:      "mcp-tool.github.create_issue",
				AuthPolicyRef: `MCPTool::"github.create_issue"`,
			}}))
			Expect(k8s.IsAutomationAPIManaged(proxy)).To(BeTrue())
		})

		It("leaves out what the state leaves out", func() {
			proxy := newMcpProxy(nil)
			proxy.Status.State = "STARTED"

			resp := apimmcpproxy.Response{State: sdk.AimMcpProxyState{Id: new("1f3c")}}
			Expect(apimmcpproxy.UpdateStatus(context.Background(), proxy, resp)).To(Succeed())

			Expect(proxy.Status.ID).To(Equal("1f3c"))
			Expect(proxy.Status.State).To(BeEmpty())
			Expect(proxy.Status.Tools).To(BeNil())
			Expect(proxy.Status.Errors.Severe).To(BeNil())
			Expect(proxy.Status.Errors.Warning).To(BeNil())
		})
	})

	Describe("calls", func() {
		const collection = "/automation/organizations/org1/environments/env1/aim/mcp-proxies"

		var (
			ctx    context.Context
			fake   *fakeAutomationAPI
			client *apim.SDKClient
		)

		BeforeEach(func() {
			ctx = context.Background()
			fake = newFakeAutomationAPI()
			fake.response = `{"hrid": "gravitee-github-mcp", "entityId": "mcp-proxy.github", "name": "GitHub",
				"contextPath": "/mcp/github", "mode": "PROXY", "state": "STARTED",
				"proxy": {"serverUrl": "https://api.githubcopilot.com/mcp/"},
				"plans": [{"name": "gold", "security": {"type": "API_KEY", "source": "HEADER"}}],
				"id": "1f3c", "organizationId": "DEFAULT", "environmentId": "DEFAULT"}`
			DeferCleanup(fake.Close)
			var err error
			client, err = apim.NewSDKClient(ctx, managementContext(fake.URL, bearer("token")))
			Expect(err).ToNot(HaveOccurred())
		})

		It("puts the payload on the proxies collection and answers the state", func() {
			resp, err := apimmcpproxy.Upsert(ctx, client, toMcpProxyDTO(newMcpProxy(nil)))
			Expect(err).ToNot(HaveOccurred())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodPut))
			Expect(req.Path).To(Equal(collection))
			Expect(req.Query.Get("dryRun")).To(Equal("false"))
			Expect(marshal(req.Body)).To(MatchJSON(`{
				"hrid": "gravitee-github-mcp",
				"entityId": "mcp-proxy.github",
				"name": "GitHub",
				"contextPath": "/mcp/github",
				"mode": "PROXY",
				"state": "STARTED",
				"proxy": {"serverUrl": "https://api.githubcopilot.com/mcp/"},
				"plans": [{"name": "gold", "security": {"type": "API_KEY", "source": "HEADER"}}]
			}`))
			Expect(*resp.State.Id).To(Equal("1f3c"))
			Expect(resp.GetOrgID()).To(Equal("DEFAULT"))
			Expect(resp.GetEnvID()).To(Equal("DEFAULT"))
		})

		It("dry runs and returns the findings of the state", func() {
			fake.response = `{"hrid": "gravitee-github-mcp", "mode": "PROXY", "plans": [],
				"errors": {"severe": ["context path [/mcp/github] is already used"], "warning": ["slow"]}}`

			findings, err := apimmcpproxy.DryRun(ctx, client, toMcpProxyDTO(newMcpProxy(nil)))
			Expect(err).ToNot(HaveOccurred())

			Expect(fake.lastRequest().Query.Get("dryRun")).To(Equal("true"))
			Expect(findings.Severe).To(HaveLen(1))
			Expect(findings.Severe[0].Message).To(Equal("context path [/mcp/github] is already used"))
			Expect(findings.Warning).To(HaveLen(1))
			Expect(findings.Warning[0].Message).To(Equal("slow"))
		})

		It("returns the findings of a refused dry run", func() {
			fake.status = http.StatusBadRequest
			fake.errorResponse = `{"name": "GitHub", "errors": {"severe": ["plan [gold] cannot be changed in place"]}}`

			findings, err := apimmcpproxy.DryRun(ctx, client, toMcpProxyDTO(newMcpProxy(nil)))
			Expect(err).ToNot(HaveOccurred())

			Expect(findings.Severe).To(HaveLen(1))
			Expect(findings.Severe[0].Message).To(Equal("plan [gold] cannot be changed in place"))
		})

		It("returns a failure that carries no findings as the error", func() {
			fake.status = http.StatusInternalServerError
			fake.errorResponse = `{"httpStatus": 500, "technicalCode": "gamma.mcp.lifecycle.failed"}`

			_, err := apimmcpproxy.DryRun(ctx, client, toMcpProxyDTO(newMcpProxy(nil)))

			Expect(err).To(MatchError(ContainSubstring("gamma.mcp.lifecycle.failed")))
		})

		It("gets the remote as the payload a PUT sends, sorted like the mapper sorts", func() {
			fake.response = `{"hrid": "gravitee-support-studio", "mode": "STUDIO", "state": "STARTED",
				"studio": {
					"tools": [{"server": "tools-jira-mcp", "tool": "create_ticket"},
						{"server": "gravitee-github-mcp", "tool": "create_issue", "entityId": "mcp-tool.github.create_issue"}],
					"upstreamAuth": [{"server": "tools-jira-mcp", "auth": {"type": "NONE"}},
						{"server": "gravitee-github-mcp", "auth": {"type": "BEARER"}}],
					"enableFGA": true
				},
				"plans": [{"name": "silver", "security": {"type": "KEY_LESS"}},
					{"name": "gold", "security": {"type": "KEY_LESS"}}],
				"tools": [{"server": "gravitee-github-mcp", "tool": "create_issue",
					"authPolicyRef": "MCPTool::\"github.create_issue\""}]}`

			remote, err := apimmcpproxy.GetRemote(ctx, client, toMcpProxyDTO(newMcpStudio()))
			Expect(err).ToNot(HaveOccurred())

			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodGet))
			Expect(req.Path).To(Equal(collection + "/gravitee-support-studio"))
			Expect(marshal(remote)).To(MatchJSON(`{
				"hrid": "gravitee-support-studio", "entityId": "", "name": "", "contextPath": "",
				"mode": "STUDIO", "state": "STARTED",
				"studio": {
					"tools": [{"server": "gravitee-github-mcp", "tool": "create_issue", "entityId": "mcp-tool.github.create_issue"},
						{"server": "tools-jira-mcp", "tool": "create_ticket"}],
					"upstreamAuth": [{"server": "gravitee-github-mcp", "auth": {"type": "BEARER"}},
						{"server": "tools-jira-mcp", "auth": {"type": "NONE"}}],
					"enableFGA": true
				},
				"plans": [{"name": "gold", "security": {"type": "KEY_LESS"}},
					{"name": "silver", "security": {"type": "KEY_LESS"}}]
			}`))
		})

		It("reports a remote that is gone as not found", func() {
			fake.status = http.StatusNotFound

			_, err := apimmcpproxy.GetRemote(ctx, client, toMcpProxyDTO(newMcpProxy(nil)))

			Expect(xErrors.IsNotFound(err)).To(BeTrue())
		})

		It("deletes by hrid and treats a proxy already gone as deleted", func() {
			dto := toMcpProxyDTO(newMcpProxy(nil))

			fake.status = http.StatusNoContent
			Expect(apimmcpproxy.Delete(ctx, client, dto)).To(Succeed())
			req := fake.lastRequest()
			Expect(req.Method).To(Equal(http.MethodDelete))
			Expect(req.Path).To(Equal(collection + "/gravitee-github-mcp"))

			fake.status = http.StatusNotFound
			Expect(apimmcpproxy.Delete(ctx, client, dto)).To(Succeed())
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
