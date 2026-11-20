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

package apim

import (
	"slices"
	"strings"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	apimmcpproxy "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/drift"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MCP proxy Drift detection", func() {

	DescribeTable("equivalent values",
		func(crd, remote any) {
			expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
		},
		Entry("empty struct",
			sdk.AimMcpProxy{},
			sdk.AimMcpProxy{},
		),
		Entry("equal struct",
			completeMcpProxyDTO(),
			completeMcpProxyDTO(),
		),
		Entry("PROXY: the platform never returns credentials",
			mcpProxyDTO(proxyMcpProxyCRD()),
			mcpProxyRemote(mcpProxyDTO(proxyMcpProxyCRD())),
		),
		Entry("STUDIO: the platform never returns credentials and adds the tool entity ids",
			mcpProxyDTO(studioMcpProxyCRD()),
			withToolEntityIDs(mcpProxyRemote(mcpProxyDTO(studioMcpProxyCRD()))),
		),
		Entry("PROXY: NONE upstream auth is stored as passthrough and omitted",
			mcpProxyDTO(withProxyAuth(proxyMcpProxyCRD(), &mcpproxy.UpstreamAuth{Type: mcpproxy.UpstreamAuthTypeNone})),
			withUpstreamAuth(mcpProxyRemote(mcpProxyDTO(proxyMcpProxyCRD())), nil),
		),
		Entry("default flow execution is omitted by the platform",
			mcpProxyDTO(withFlowExecution(proxyMcpProxyCRD(), &v4.FlowExecution{Mode: v4.FlowModeDefault})),
			withFlowExecutionDTO(mcpProxyRemote(mcpProxyDTO(proxyMcpProxyCRD())), nil),
		),
		Entry("empty flows are omitted by the platform",
			withFlows(mcpProxyDTO(proxyMcpProxyCRD()), []sdk.AimMcpProxyFlow{}),
			withFlows(mcpProxyRemote(mcpProxyDTO(proxyMcpProxyCRD())), nil),
		),
		Entry("STUDIO: upstream auth in another order once canonicalized",
			mcpProxyDTO(studioMcpProxyCRD()),
			apimmcpproxy.Canonical(reversedStudioAuth(mcpProxyRemote(mcpProxyDTO(studioMcpProxyCRD())))),
		),
		Entry("hrid is identity, not payload",
			withMcpProxyHRID(completeMcpProxyDTO(), "default-drift-full"),
			withMcpProxyHRID(completeMcpProxyDTO(), "something-else"),
		),
		Entry("STUDIO: tool-level authorization off is omitted by the platform",
			withEnableFGA(mcpProxyDTO(studioMcpProxyCRD()), new(false)),
			withEnableFGA(mcpProxyRemote(mcpProxyDTO(studioMcpProxyCRD())), nil),
		),
		Entry("a flow execution that does not require a match is omitted by the platform",
			mcpProxyDTO(withFlowExecution(proxyMcpProxyCRD(), &v4.FlowExecution{Mode: v4.FlowModeBestMatch})),
			withFlowExecutionDTO(mcpProxyRemote(mcpProxyDTO(proxyMcpProxyCRD())), &sdk.AimMcpProxyFlowExecution{
				Mode: new(sdk.AimMcpProxyFlowExecutionModeBESTMATCH),
			}),
		),
	)

	DescribeTable("drifted values",
		func(crd, remote any) {
			result := drift.DetectWithNamespace(crd, remote, "")
			Expect(result.DriftDetected()).To(BeTrue())
		},
		Entry("context path changed on the platform",
			mcpProxyDTO(proxyMcpProxyCRD()),
			withContextPath(mcpProxyRemote(mcpProxyDTO(proxyMcpProxyCRD())), "/mcp/elsewhere"),
		),
		Entry("plan added on the platform",
			mcpProxyDTO(proxyMcpProxyCRD()),
			withExtraPlan(mcpProxyRemote(mcpProxyDTO(proxyMcpProxyCRD())), "Console"),
		),
		Entry("tool-level authorization switched off on the platform",
			mcpProxyDTO(studioMcpProxyCRD()),
			withEnableFGA(mcpProxyRemote(mcpProxyDTO(studioMcpProxyCRD())), new(false)),
		),
		Entry("tool-level authorization omitted by the platform",
			mcpProxyDTO(studioMcpProxyCRD()),
			withEnableFGA(mcpProxyRemote(mcpProxyDTO(studioMcpProxyCRD())), nil),
		),
		Entry("proxy stopped on the platform",
			mcpProxyDTO(proxyMcpProxyCRD()),
			withState(mcpProxyRemote(mcpProxyDTO(proxyMcpProxyCRD())), sdk.AimMcpProxyLifecycleStateSTOPPED),
		),
		Entry("upstream auth type changed on the platform",
			mcpProxyDTO(proxyMcpProxyCRD()),
			withUpstreamAuth(mcpProxyRemote(mcpProxyDTO(proxyMcpProxyCRD())), &sdk.AimMcpProxyUpstreamAuth{
				Bearer: &sdk.AimMcpProxyBearerUpstreamAuth{Type: sdk.AimMcpProxyBearerUpstreamAuthTypeBEARER},
			}),
		),
	)

	It("names the upstream auth type the platform switched to", func() {
		crd := mcpProxyDTO(proxyMcpProxyCRD())
		remote := withUpstreamAuth(mcpProxyRemote(crd), &sdk.AimMcpProxyUpstreamAuth{
			Bearer: &sdk.AimMcpProxyBearerUpstreamAuth{Type: sdk.AimMcpProxyBearerUpstreamAuthTypeBEARER},
		})

		merged := admissionDrift(crd, remote)

		Expect(merged.DriftDetected()).To(BeTrue())
		Expect(merged.String()).To(Equal(`proxy:
  upstreamAuth:
    type: "BASIC" != "BEARER"
    username: "gravitee" != <nil>`))
	})

	It("admits the defaults APIM saves in a flow step configuration", func() {
		crd := mcpProxyWithRequestStep(`{"rate": {"limit": 100}}`)
		remote := mcpProxyWithRequestStep(`{"rate": {"limit": 100}, "addHeaders": false, "async": false}`)

		expectNoDrift(drift.DetectWithNamespace(crd, remote, ""))
	})

	DescribeTable("user info endpoint of an identity provider",
		func(crdEndpoint, remoteEndpoint *string, report string) {
			crd := mcpProxyDTO(withGenericIdentityProvider(proxyMcpProxyCRD(), crdEndpoint))
			remote := withUserInfoEndpoint(mcpProxyRemote(crd), remoteEndpoint)

			merged := admissionDrift(crd, remote)

			if report == "" {
				expectNoDrift(merged)
				return
			}
			Expect(merged.DriftDetected()).To(BeTrue())
			Expect(merged.String()).To(ContainSubstring(report))
		},
		Entry("left out, APIM saves /userinfo", nil, new("/userinfo"), ""),
		Entry("left out, APIM stores it empty", nil, new(""), ""),
		Entry("empty, APIM omits it", new(""), nil, ""),
		Entry("empty, APIM saves /userinfo", new(""), new("/userinfo"), ""),
		Entry("left out, changed to /me in APIM", nil, new("/me"), `userInfoEndpoint: <nil> != "/me"`),
		Entry("declared /me, APIM has /userinfo", new("/me"), new("/userinfo"), `userInfoEndpoint: "/me" != "/userinfo"`),
	)

	It("detects an identity provider only APIM has", func() {
		crd := mcpProxyWithIdentityProviders(keycloakIdentityProvider(""))
		remote := mcpProxyWithIdentityProviders(
			keycloakIdentityProvider(""),
			`{"name": "console", "type": "OAUTH2_GENERIC", "issuerUrl": "https://console.example.com"}`,
		)

		merged := admissionDrift(crd, remote)

		Expect(merged.DriftDetected()).To(BeTrue())
	})

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeMcpProxyDTO(), completeMcpProxyDTO())
		})
	})
})

// completeMcpProxyDTO sets every field, every union variant included. JSON carries one variant of a
// union, so the fixture holds one and the others are set here.
func completeMcpProxyDTO() sdk.AimMcpProxy {
	GinkgoHelper()
	dto := loadFixture[sdk.AimMcpProxy]("mcpproxy_full_dto.json")

	dto.Proxy.UpstreamAuth = new(withEveryUpstreamAuth(*dto.Proxy.UpstreamAuth))
	dto.Studio.UpstreamAuth[0].Auth = withEveryUpstreamAuth(dto.Studio.UpstreamAuth[0].Auth)

	dto.Flows[0].Selectors[0].Condition = &sdk.AimMcpProxyConditionFlowSelector{
		Type:      sdk.AimMcpProxyConditionFlowSelectorTypeCONDITION,
		Condition: "{#request.headers['X-Audit'] != null}",
	}
	dto.Plans[0].Flows[0].Selectors[0].Mcp = &sdk.AimMcpProxyMcpFlowSelector{
		Type:    sdk.AimMcpProxyMcpFlowSelectorTypeMCP,
		Methods: []string{"tools/list"},
	}

	dto.IdentityProviders[0].GraviteeAm = &sdk.AimMcpProxyGraviteeAmIdentityProvider{
		Type: sdk.AimMcpProxyGraviteeAmIdentityProviderTypeGRAVITEEAM,
		Name: "am",
	}
	dto.IdentityProviders[0].OAuth2Auth0 = &sdk.AimMcpProxyAuth0IdentityProvider{
		Type:     sdk.AimMcpProxyAuth0IdentityProviderTypeOAUTH2AUTH0,
		Name:     "auth0",
		Domain:   "gravitee.eu.auth0.com",
		Audience: "https://mcp.example.com",
	}

	dto.Plans[0].Security.KeyLess = &sdk.AimMcpProxyKeyLessPlanSecurity{
		Type: sdk.AimMcpProxyKeyLessPlanSecurityTypeKEYLESS,
	}
	dto.Plans[0].Security.OAuth2 = &sdk.AimMcpProxyOAuth2PlanSecurity{
		Type:     sdk.AimMcpProxyOAuth2PlanSecurityTypeOAUTH2,
		Provider: "keycloak",
		Scopes:   []string{"repo:read"},
	}
	return dto
}

// withEveryUpstreamAuth adds every upstream auth variant to the one the fixture holds.
func withEveryUpstreamAuth(auth sdk.AimMcpProxyUpstreamAuth) sdk.AimMcpProxyUpstreamAuth {
	auth.None = &sdk.AimMcpProxyNoUpstreamAuth{Type: sdk.AimMcpProxyNoUpstreamAuthTypeNONE}
	auth.Passthrough = &sdk.AimMcpProxyPassthroughUpstreamAuth{
		Type: sdk.AimMcpProxyPassthroughUpstreamAuthTypePASSTHROUGH,
	}
	auth.ApiKey = &sdk.AimMcpProxyApiKeyUpstreamAuth{
		Type:         sdk.AimMcpProxyApiKeyUpstreamAuthTypeAPIKEY,
		ApiKeyHeader: "X-Api-Key",
		ApiKey:       new("k"),
	}
	auth.Bearer = &sdk.AimMcpProxyBearerUpstreamAuth{
		Type:  sdk.AimMcpProxyBearerUpstreamAuthTypeBEARER,
		Token: new("t"),
	}
	auth.Basic = &sdk.AimMcpProxyBasicUpstreamAuth{
		Type:     sdk.AimMcpProxyBasicUpstreamAuthTypeBASIC,
		Username: "gravitee",
		Password: new("p"),
	}
	return auth
}

func mcpProxyDTO(proxy *v1alpha1.McpProxy) sdk.AimMcpProxy {
	GinkgoHelper()
	dto, err := apimmcpproxy.ToDTO(proxy)
	Expect(err).ToNot(HaveOccurred())
	return dto
}

func proxyMcpProxyCRD() *v1alpha1.McpProxy {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.McpProxy]("mcpproxy_proxy_crd.json")
	return &fixture
}

func studioMcpProxyCRD() *v1alpha1.McpProxy {
	GinkgoHelper()
	fixture := loadFixture[v1alpha1.McpProxy]("mcpproxy_studio_crd.json")
	return &fixture
}

// mcpProxyRemote is what the platform answers for a payload: every credential removed.
func mcpProxyRemote(dto sdk.AimMcpProxy) sdk.AimMcpProxy {
	if dto.Proxy != nil {
		proxy := *dto.Proxy
		if proxy.UpstreamAuth != nil {
			proxy.UpstreamAuth = new(withoutSecrets(*proxy.UpstreamAuth))
		}
		dto.Proxy = &proxy
	}
	if dto.Studio != nil {
		studio := *dto.Studio
		studio.UpstreamAuth = slices.Clone(studio.UpstreamAuth)
		for i := range studio.UpstreamAuth {
			studio.UpstreamAuth[i].Auth = withoutSecrets(studio.UpstreamAuth[i].Auth)
		}
		dto.Studio = &studio
	}
	dto.IdentityProviders = slices.Clone(dto.IdentityProviders)
	for i := range dto.IdentityProviders {
		if generic := dto.IdentityProviders[i].OAuth2Generic; generic != nil {
			stripped := *generic
			stripped.ClientSecret = nil
			dto.IdentityProviders[i].OAuth2Generic = &stripped
		}
	}
	return dto
}

func withoutSecrets(auth sdk.AimMcpProxyUpstreamAuth) sdk.AimMcpProxyUpstreamAuth {
	if auth.ApiKey != nil {
		apiKey := *auth.ApiKey
		apiKey.ApiKey = nil
		auth.ApiKey = &apiKey
	}
	if auth.Bearer != nil {
		bearer := *auth.Bearer
		bearer.Token = nil
		auth.Bearer = &bearer
	}
	if auth.Basic != nil {
		basic := *auth.Basic
		basic.Password = nil
		auth.Basic = &basic
	}
	if auth.OAuth2 != nil {
		oauth2 := *auth.OAuth2
		oauth2.ClientSecret = nil
		auth.OAuth2 = &oauth2
	}
	return auth
}

func withToolEntityIDs(dto sdk.AimMcpProxy) sdk.AimMcpProxy {
	studio := *dto.Studio
	studio.Tools = slices.Clone(studio.Tools)
	for i := range studio.Tools {
		studio.Tools[i].EntityId = new("mcp-tool." + studio.Tools[i].Tool)
	}
	dto.Studio = &studio
	return dto
}

func reversedStudioAuth(dto sdk.AimMcpProxy) sdk.AimMcpProxy {
	studio := *dto.Studio
	studio.UpstreamAuth = slices.Clone(studio.UpstreamAuth)
	slices.Reverse(studio.UpstreamAuth)
	dto.Studio = &studio
	return dto
}

func withProxyAuth(crd *v1alpha1.McpProxy, auth *mcpproxy.UpstreamAuth) *v1alpha1.McpProxy {
	crd.Spec.Proxy.UpstreamAuth = auth
	return crd
}

func withFlowExecution(crd *v1alpha1.McpProxy, execution *v4.FlowExecution) *v1alpha1.McpProxy {
	crd.Spec.FlowExecution = execution
	return crd
}

func withUpstreamAuth(dto sdk.AimMcpProxy, auth *sdk.AimMcpProxyUpstreamAuth) sdk.AimMcpProxy {
	proxy := *dto.Proxy
	proxy.UpstreamAuth = auth
	dto.Proxy = &proxy
	return dto
}

func withFlowExecutionDTO(dto sdk.AimMcpProxy, execution *sdk.AimMcpProxyFlowExecution) sdk.AimMcpProxy {
	dto.FlowExecution = execution
	return dto
}

func withFlows(dto sdk.AimMcpProxy, flows []sdk.AimMcpProxyFlow) sdk.AimMcpProxy {
	dto.Flows = flows
	return dto
}

func withContextPath(dto sdk.AimMcpProxy, contextPath string) sdk.AimMcpProxy {
	dto.ContextPath = contextPath
	return dto
}

func withExtraPlan(dto sdk.AimMcpProxy, name string) sdk.AimMcpProxy {
	dto.Plans = append(slices.Clone(dto.Plans), sdk.AimMcpProxyPlan{
		Name: name,
		Security: sdk.AimMcpProxyPlanSecurity{
			KeyLess: &sdk.AimMcpProxyKeyLessPlanSecurity{Type: sdk.AimMcpProxyKeyLessPlanSecurityTypeKEYLESS},
		},
	})
	return dto
}

func withEnableFGA(dto sdk.AimMcpProxy, enabled *bool) sdk.AimMcpProxy {
	studio := *dto.Studio
	studio.EnableFGA = enabled
	dto.Studio = &studio
	return dto
}

func withState(dto sdk.AimMcpProxy, state sdk.AimMcpProxyLifecycleState) sdk.AimMcpProxy {
	dto.State = &state
	return dto
}

func withMcpProxyHRID(dto sdk.AimMcpProxy, hrid string) sdk.AimMcpProxy {
	dto.Hrid = hrid
	return dto
}

func mcpProxyWithRequestStep(configuration string) sdk.AimMcpProxy {
	GinkgoHelper()
	return parseJSON[sdk.AimMcpProxy](`{"name": "mcp", "flows": [{"name": "Log tool calls", "enabled": true, "request": [` +
		`{"name": "Quota", "policy": "rate-limit", "enabled": true, "configuration": ` + configuration + `}]}]}`)
}

func mcpProxyWithIdentityProviders(providers ...string) sdk.AimMcpProxy {
	GinkgoHelper()
	return parseJSON[sdk.AimMcpProxy](`{"name": "mcp", "identityProviders": [` + strings.Join(providers, ", ") + `]}`)
}

// keycloakIdentityProvider is an OAUTH2_GENERIC identity provider; extra holds more fields, each
// with a leading comma.
func keycloakIdentityProvider(extra string) string {
	return `{"name": "keycloak", "type": "OAUTH2_GENERIC", "issuerUrl": "https://kc.example.com/realms/acme"` + extra + `}`
}

// withGenericIdentityProvider adds an OAUTH2_GENERIC identity provider with the given user info endpoint.
func withGenericIdentityProvider(crd *v1alpha1.McpProxy, userInfoEndpoint *string) *v1alpha1.McpProxy {
	crd.Spec.IdentityProviders = append(crd.Spec.IdentityProviders, mcpproxy.IdentityProvider{
		Name: "keycloak",
		Type: mcpproxy.IdentityProviderOAuth2Generic,
		OAuth2Generic: &mcpproxy.OAuth2GenericProvider{
			IssuerURL:                   "https://kc.example.com/realms/acme",
			IntrospectionEndpoint:       "/protocol/openid-connect/token/introspect",
			IntrospectionEndpointMethod: "POST",
			UserInfoEndpoint:            userInfoEndpoint,
		},
	})
	return crd
}

// withUserInfoEndpoint sets the user info endpoint of the keycloak identity provider.
func withUserInfoEndpoint(dto sdk.AimMcpProxy, userInfoEndpoint *string) sdk.AimMcpProxy {
	dto.IdentityProviders = slices.Clone(dto.IdentityProviders)
	for i := range dto.IdentityProviders {
		if generic := dto.IdentityProviders[i].OAuth2Generic; generic != nil && generic.Name == "keycloak" {
			keycloak := *generic
			keycloak.UserInfoEndpoint = userInfoEndpoint
			dto.IdentityProviders[i].OAuth2Generic = &keycloak
		}
	}
	return dto
}
