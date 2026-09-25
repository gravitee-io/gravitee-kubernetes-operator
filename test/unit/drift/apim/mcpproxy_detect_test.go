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

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
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
			model.McpProxyDTO{},
			model.McpProxyDTO{},
		),
		Entry("equal struct",
			completeMcpProxyDTO(),
			completeMcpProxyDTO(),
		),
		Entry("PROXY: the platform never returns credentials",
			model.ToMcpProxyDTO(proxyMcpProxyCRD()),
			mcpProxyRemote(model.ToMcpProxyDTO(proxyMcpProxyCRD())),
		),
		Entry("STUDIO: the platform never returns credentials and adds the tool entity ids",
			model.ToMcpProxyDTO(studioMcpProxyCRD()),
			withToolEntityIDs(mcpProxyRemote(model.ToMcpProxyDTO(studioMcpProxyCRD()))),
		),
		Entry("PROXY: NONE upstream auth is stored as passthrough and omitted",
			model.ToMcpProxyDTO(withProxyAuth(proxyMcpProxyCRD(), &mcpproxy.UpstreamAuth{Type: mcpproxy.UpstreamAuthTypeNone})),
			withUpstreamAuth(mcpProxyRemote(model.ToMcpProxyDTO(proxyMcpProxyCRD())), nil),
		),
		Entry("default flow execution is omitted by the platform",
			model.ToMcpProxyDTO(withFlowExecution(proxyMcpProxyCRD(), &v4.FlowExecution{Mode: v4.FlowModeDefault})),
			withFlowExecutionDTO(mcpProxyRemote(model.ToMcpProxyDTO(proxyMcpProxyCRD())), nil),
		),
		Entry("empty flows are omitted by the platform",
			withFlows(model.ToMcpProxyDTO(proxyMcpProxyCRD()), []model.McpProxyFlowDTO{}),
			withFlows(mcpProxyRemote(model.ToMcpProxyDTO(proxyMcpProxyCRD())), nil),
		),
		Entry("STUDIO: upstream auth in another order once canonicalized",
			model.ToMcpProxyDTO(studioMcpProxyCRD()),
			canonicalized(reversedStudioAuth(mcpProxyRemote(model.ToMcpProxyDTO(studioMcpProxyCRD())))),
		),
		Entry("hrid is identity, not payload",
			withMcpProxyHRID(completeMcpProxyDTO(), "default-drift-full"),
			withMcpProxyHRID(completeMcpProxyDTO(), "something-else"),
		),
	)

	DescribeTable("drifted values",
		func(crd, remote any) {
			result := drift.DetectWithNamespace(crd, remote, "")
			Expect(result.DriftDetected()).To(BeTrue())
		},
		Entry("context path changed on the platform",
			model.ToMcpProxyDTO(proxyMcpProxyCRD()),
			withContextPath(mcpProxyRemote(model.ToMcpProxyDTO(proxyMcpProxyCRD())), "/mcp/elsewhere"),
		),
		Entry("plan added on the platform",
			model.ToMcpProxyDTO(proxyMcpProxyCRD()),
			withExtraPlan(mcpProxyRemote(model.ToMcpProxyDTO(proxyMcpProxyCRD())), "Console"),
		),
		Entry("tool-level authorization switched off on the platform",
			model.ToMcpProxyDTO(studioMcpProxyCRD()),
			withEnableFGA(mcpProxyRemote(model.ToMcpProxyDTO(studioMcpProxyCRD())), false),
		),
		Entry("proxy stopped on the platform",
			model.ToMcpProxyDTO(proxyMcpProxyCRD()),
			withState(mcpProxyRemote(model.ToMcpProxyDTO(proxyMcpProxyCRD())), "STOPPED"),
		),
	)

	Describe("All properties regression test", func() {
		It("ensure no new property isn't tested are tested", func() {
			expectedEquivalentNotHavingAnyZeroValue(completeMcpProxyDTO(), completeMcpProxyDTO())
		})
	})
})

func completeMcpProxyDTO() model.McpProxyDTO {
	GinkgoHelper()
	return loadFixture[model.McpProxyDTO]("mcpproxy_full_dto.json")
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
func mcpProxyRemote(dto model.McpProxyDTO) model.McpProxyDTO {
	strip := func(auth model.McpProxyAuthDTO) model.McpProxyAuthDTO {
		auth.APIKey, auth.Token, auth.Password, auth.ClientSecret = nil, nil, nil, nil
		return auth
	}
	if dto.UpstreamAuth != nil {
		dto.UpstreamAuth = new(strip(*dto.UpstreamAuth))
	}
	if dto.Studio != nil {
		studio := *dto.Studio
		studio.UpstreamAuth = slices.Clone(studio.UpstreamAuth)
		for i := range studio.UpstreamAuth {
			studio.UpstreamAuth[i].Auth = strip(studio.UpstreamAuth[i].Auth)
		}
		dto.Studio = &studio
	}
	dto.IdentityProviders = slices.Clone(dto.IdentityProviders)
	for i := range dto.IdentityProviders {
		dto.IdentityProviders[i].ClientSecret = nil
	}
	return dto
}

func withToolEntityIDs(dto model.McpProxyDTO) model.McpProxyDTO {
	studio := *dto.Studio
	studio.Tools = slices.Clone(studio.Tools)
	for i := range studio.Tools {
		studio.Tools[i].EntityID = "mcp-tool." + studio.Tools[i].Tool
	}
	dto.Studio = &studio
	return dto
}

func reversedStudioAuth(dto model.McpProxyDTO) model.McpProxyDTO {
	studio := *dto.Studio
	studio.UpstreamAuth = slices.Clone(studio.UpstreamAuth)
	slices.Reverse(studio.UpstreamAuth)
	dto.Studio = &studio
	return dto
}

func canonicalized(dto model.McpProxyDTO) model.McpProxyDTO {
	dto.Canonicalize()
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

func withUpstreamAuth(dto model.McpProxyDTO, auth *model.McpProxyAuthDTO) model.McpProxyDTO {
	dto.UpstreamAuth = auth
	return dto
}

func withFlowExecutionDTO(dto model.McpProxyDTO, execution *v4.FlowExecution) model.McpProxyDTO {
	dto.FlowExecution = execution
	return dto
}

func withFlows(dto model.McpProxyDTO, flows []model.McpProxyFlowDTO) model.McpProxyDTO {
	dto.Flows = flows
	return dto
}

func withContextPath(dto model.McpProxyDTO, contextPath string) model.McpProxyDTO {
	dto.ContextPath = contextPath
	return dto
}

func withExtraPlan(dto model.McpProxyDTO, name string) model.McpProxyDTO {
	dto.Plans = append(slices.Clone(dto.Plans), model.McpProxyPlanDTO{
		Name:     name,
		Security: model.McpProxyPlanSecurityDTO{Type: "KEY_LESS"},
	})
	return dto
}

func withEnableFGA(dto model.McpProxyDTO, enabled bool) model.McpProxyDTO {
	studio := *dto.Studio
	studio.EnableFGA = enabled
	dto.Studio = &studio
	return dto
}

func withState(dto model.McpProxyDTO, state string) model.McpProxyDTO {
	dto.State = state
	return dto
}

func withMcpProxyHRID(dto model.McpProxyDTO, hrid string) model.McpProxyDTO {
	dto.HRID = hrid
	return dto
}
