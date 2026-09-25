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

package model

import (
	"cmp"
	"slices"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

// McpProxyDTO is the automation API wire representation of an McpProxy (PUT /aim/mcp-proxies).
// The platform refuses unknown properties, so this carries exactly the wire fields: the CRD's
// proxy block is flattened to serverUrl and upstreamAuth, and every union is flat, discriminated
// by type.
type McpProxyDTO struct {
	HRID            string  `json:"hrid,omitempty" drift:"ignore"`
	EntityID        string  `json:"entityId"`
	Name            string  `json:"name"`
	Description     *string `json:"description,omitempty" drift:"empty-is-nil"`
	ContextPath     string  `json:"contextPath"`
	Mode            string  `json:"mode"`
	ProtocolVersion string  `json:"protocolVersion"`
	// Declared lifecycle on a PUT; the lifecycle the platform observes in a response.
	State     string  `json:"state"`
	ServerURL *string `json:"serverUrl,omitempty" drift:"empty-is-nil"`
	// Nil for a passthrough proxy: the platform stores NONE as no authentication and omits it.
	UpstreamAuth      *McpProxyAuthDTO              `json:"upstreamAuth,omitempty" drift:"empty-is-nil"`
	Studio            *McpProxyStudioDTO            `json:"studio,omitempty" drift:"empty-is-nil"`
	FlowExecution     *v4.FlowExecution             `json:"flowExecution,omitempty" drift:"empty-is-nil"`
	Flows             []McpProxyFlowDTO             `json:"flows,omitempty" drift:"empty-is-nil"`
	IdentityProviders []McpProxyIdentityProviderDTO `json:"identityProviders,omitempty" drift:"empty-is-nil"`
	Plans             []McpProxyPlanDTO             `json:"plans"`
}

// McpProxyAuthDTO is the wire shape of an upstream credential. Secret values are never
// returned by the platform, so they are ignored for drift.
type McpProxyAuthDTO struct {
	Type         string   `json:"type"`
	APIKeyHeader *string  `json:"apiKeyHeader,omitempty"`
	APIKey       *string  `json:"apiKey,omitempty" drift:"ignore"`
	Token        *string  `json:"token,omitempty" drift:"ignore"`
	Username     *string  `json:"username,omitempty"`
	Password     *string  `json:"password,omitempty" drift:"ignore"`
	AuthorizeURL *string  `json:"authorizeUrl,omitempty"`
	TokenURL     *string  `json:"tokenUrl,omitempty"`
	ClientID     *string  `json:"clientId,omitempty"`
	ClientSecret *string  `json:"clientSecret,omitempty" drift:"ignore"`
	Scopes       []string `json:"scopes,omitempty" drift:"empty-is-nil"`
}

type McpProxyStudioDTO struct {
	Tools        []McpProxyStudioToolDTO `json:"tools"`
	UpstreamAuth []McpProxyStudioAuthDTO `json:"upstreamAuth"`
	EnableFGA    bool                    `json:"enableFGA"`
}

// McpProxyStudioToolDTO selects a tool by the HRID of its catalog server. The platform adds the
// composed entityId in responses; it is read-only and never sent.
type McpProxyStudioToolDTO struct {
	Server   string  `json:"server"`
	Tool     string  `json:"tool"`
	Alias    *string `json:"alias,omitempty" drift:"empty-is-nil"`
	EntityID string  `json:"entityId,omitempty" drift:"ignore"`
}

type McpProxyStudioAuthDTO struct {
	Server string          `json:"server"`
	Auth   McpProxyAuthDTO `json:"auth"`
}

type McpProxyPlanDTO struct {
	Name     string                  `json:"name"`
	Security McpProxyPlanSecurityDTO `json:"security"`
	Flows    []McpProxyFlowDTO       `json:"flows,omitempty" drift:"empty-is-nil"`
}

type McpProxyPlanSecurityDTO struct {
	Type            string   `json:"type"`
	Source          *string  `json:"source,omitempty"`
	APIKeyHeader    *string  `json:"apiKeyHeader,omitempty" drift:"empty-is-nil"`
	PropagateAPIKey bool     `json:"propagateApiKey,omitempty" drift:"empty-is-nil"`
	Provider        *string  `json:"provider,omitempty"`
	Scopes          []string `json:"scopes,omitempty" drift:"empty-is-nil"`
}

type McpProxyIdentityProviderDTO struct {
	Name                        string  `json:"name"`
	Type                        string  `json:"type"`
	IssuerURL                   *string `json:"issuerUrl,omitempty"`
	IntrospectionEndpoint       *string `json:"introspectionEndpoint,omitempty"`
	IntrospectionEndpointMethod *string `json:"introspectionEndpointMethod,omitempty"`
	ClientID                    *string `json:"clientId,omitempty" drift:"empty-is-nil"`
	UserInfoEndpoint            *string `json:"userInfoEndpoint,omitempty" drift:"empty-is-nil"`
	ClientSecret                *string `json:"clientSecret,omitempty" drift:"ignore"`
	Domain                      *string `json:"domain,omitempty"`
	Audience                    *string `json:"audience,omitempty"`
}

type McpProxyFlowDTO struct {
	Name      string                    `json:"name"`
	Enabled   bool                      `json:"enabled"`
	Selectors []McpProxyFlowSelectorDTO `json:"selectors,omitempty" drift:"empty-is-nil"`
	Request   []McpProxyFlowStepDTO     `json:"request,omitempty" drift:"empty-is-nil"`
	Response  []McpProxyFlowStepDTO     `json:"response,omitempty" drift:"empty-is-nil"`
	Tags      []string                  `json:"tags,omitempty" drift:"empty-is-nil"`
}

type McpProxyFlowSelectorDTO struct {
	Type      string   `json:"type"`
	Methods   []string `json:"methods,omitempty" drift:"empty-is-nil"`
	Condition *string  `json:"condition,omitempty" drift:"empty-is-nil"`
}

type McpProxyFlowStepDTO struct {
	Name          string                  `json:"name"`
	Policy        string                  `json:"policy"`
	Description   *string                 `json:"description,omitempty" drift:"empty-is-nil"`
	Enabled       bool                    `json:"enabled"`
	Condition     *string                 `json:"condition,omitempty" drift:"empty-is-nil"`
	Configuration *utils.GenericStringMap `json:"configuration,omitempty" drift:"unstructured"`
}

// McpProxyState is what a PUT or a GET answers: the payload plus the platform's identifiers
// and findings.
type McpProxyState struct {
	McpProxyDTO `json:",inline"`
	ID          string        `json:"id,omitempty"`
	OrgID       string        `json:"organizationId,omitempty"`
	EnvID       string        `json:"environmentId,omitempty"`
	Errors      status.Errors `json:"errors,omitempty"`
}

// ToMcpProxyDTO maps the CRD onto the wire payload. It is the single mapping used by the sync
// path and by drift detection.
func ToMcpProxyDTO(crd *v1alpha1.McpProxy) McpProxyDTO {
	spec := crd.Spec.Type

	mode := spec.Mode
	if mode == "" {
		mode = mcpproxy.ModeProxy
	}
	state := spec.State
	if state == "" {
		state = mcpproxy.StateStarted
	}

	dto := McpProxyDTO{
		HRID:              refs.NewNamespacedNameFromObject(crd).HRID(),
		EntityID:          spec.EntityID,
		Name:              spec.Name,
		Description:       spec.Description,
		ContextPath:       spec.ContextPath,
		Mode:              string(mode),
		ProtocolVersion:   spec.ProtocolVersion,
		State:             string(state),
		FlowExecution:     toMcpProxyFlowExecution(spec.FlowExecution),
		Flows:             toMcpProxyFlowDTOs(spec.Flows),
		IdentityProviders: toMcpProxyIdentityProviderDTOs(spec.IdentityProviders),
		Plans:             toMcpProxyPlanDTOs(spec.Plans),
	}

	switch mode {
	case mcpproxy.ModeProxy:
		if spec.Proxy != nil {
			dto.ServerURL = new(spec.Proxy.ServerURL)
			if auth := spec.Proxy.UpstreamAuth; auth != nil && auth.Type != mcpproxy.UpstreamAuthTypeNone {
				dto.UpstreamAuth = new(toMcpProxyAuthDTO(*auth))
			}
		}
	case mcpproxy.ModeStudio:
		if spec.Studio != nil {
			dto.Studio = toMcpProxyStudioDTO(spec.Studio, crd.GetNamespace())
		}
	}

	dto.Canonicalize()

	return dto
}

// Canonicalize sorts the collections the platform returns sorted, so that a mapped CRD and a
// remote state compare item by item.
func (dto *McpProxyDTO) Canonicalize() {
	slices.SortStableFunc(dto.Plans, func(a, b McpProxyPlanDTO) int { return cmp.Compare(a.Name, b.Name) })
	slices.SortStableFunc(dto.IdentityProviders, func(a, b McpProxyIdentityProviderDTO) int {
		return cmp.Compare(a.Name, b.Name)
	})
	if dto.Studio != nil {
		slices.SortStableFunc(dto.Studio.Tools, func(a, b McpProxyStudioToolDTO) int {
			return cmp.Or(cmp.Compare(a.Server, b.Server), cmp.Compare(a.Tool, b.Tool))
		})
		slices.SortStableFunc(dto.Studio.UpstreamAuth, func(a, b McpProxyStudioAuthDTO) int {
			return cmp.Compare(a.Server, b.Server)
		})
	}
}

func toMcpProxyStudioDTO(studio *mcpproxy.Studio, ns string) *McpProxyStudioDTO {
	dto := &McpProxyStudioDTO{
		Tools:        make([]McpProxyStudioToolDTO, 0, len(studio.Tools)),
		UpstreamAuth: make([]McpProxyStudioAuthDTO, 0, len(studio.UpstreamAuth)),
		EnableFGA:    studio.EnableFGA,
	}
	for _, tool := range studio.Tools {
		dto.Tools = append(dto.Tools, McpProxyStudioToolDTO{
			Server: serverHRID(tool.ServerRef, ns),
			Tool:   tool.Tool,
			Alias:  tool.Alias,
		})
	}
	for _, entry := range studio.UpstreamAuth {
		dto.UpstreamAuth = append(dto.UpstreamAuth, McpProxyStudioAuthDTO{
			Server: serverHRID(entry.ServerRef, ns),
			Auth:   toMcpProxyAuthDTO(entry.Auth),
		})
	}
	return dto
}

// serverHRID is the HRID of a referenced CatalogMcpServer, derived like its own.
func serverHRID(ref refs.NamespacedName, ns string) string {
	server := mcpproxy.WithNamespace(ref, ns)
	return server.HRID()
}

func toMcpProxyAuthDTO(auth mcpproxy.UpstreamAuth) McpProxyAuthDTO {
	dto := McpProxyAuthDTO{Type: string(auth.Type)}
	if dto.Type == "" {
		dto.Type = string(mcpproxy.UpstreamAuthTypeNone)
	}

	switch auth.Type {
	case mcpproxy.UpstreamAuthTypeAPIKey:
		if auth.APIKey != nil {
			dto.APIKeyHeader = new(auth.APIKey.Header)
			dto.APIKey = new(auth.APIKey.Value)
		}
	case mcpproxy.UpstreamAuthTypeBearer:
		if auth.Bearer != nil {
			dto.Token = new(auth.Bearer.Token)
		}
	case mcpproxy.UpstreamAuthTypeBasic:
		if auth.Basic != nil {
			dto.Username = new(auth.Basic.Username)
			dto.Password = new(auth.Basic.Password)
		}
	case mcpproxy.UpstreamAuthTypeOAuth2:
		if auth.OAuth2 != nil {
			dto.AuthorizeURL = new(auth.OAuth2.AuthorizeURL)
			dto.TokenURL = new(auth.OAuth2.TokenURL)
			dto.ClientID = new(auth.OAuth2.ClientID)
			dto.ClientSecret = new(auth.OAuth2.ClientSecret)
			dto.Scopes = auth.OAuth2.Scopes
		}
	case mcpproxy.UpstreamAuthTypeNone:
	}

	return dto
}

func toMcpProxyPlanDTOs(plans []mcpproxy.Plan) []McpProxyPlanDTO {
	dtos := make([]McpProxyPlanDTO, 0, len(plans))
	for _, plan := range plans {
		dtos = append(dtos, McpProxyPlanDTO{
			Name:     plan.Name,
			Security: toMcpProxyPlanSecurityDTO(plan.Security),
			Flows:    toMcpProxyFlowDTOs(plan.Flows),
		})
	}
	return dtos
}

func toMcpProxyPlanSecurityDTO(security mcpproxy.PlanSecurity) McpProxyPlanSecurityDTO {
	dto := McpProxyPlanSecurityDTO{Type: string(security.Type)}

	switch security.Type {
	case mcpproxy.PlanSecurityAPIKey:
		if security.APIKey != nil {
			dto.Source = new(string(security.APIKey.Source))
			dto.APIKeyHeader = security.APIKey.Header
			dto.PropagateAPIKey = security.APIKey.Propagate
		}
	case mcpproxy.PlanSecurityOAuth2:
		if security.OAuth2 != nil {
			dto.Provider = new(security.OAuth2.Provider)
			dto.Scopes = security.OAuth2.Scopes
		}
	case mcpproxy.PlanSecurityKeyLess:
	}

	return dto
}

func toMcpProxyIdentityProviderDTOs(providers []mcpproxy.IdentityProvider) []McpProxyIdentityProviderDTO {
	if len(providers) == 0 {
		return nil
	}
	dtos := make([]McpProxyIdentityProviderDTO, 0, len(providers))
	for _, provider := range providers {
		dto := McpProxyIdentityProviderDTO{Name: provider.Name, Type: string(provider.Type)}
		switch provider.Type {
		case mcpproxy.IdentityProviderOAuth2Generic:
			if generic := provider.OAuth2Generic; generic != nil {
				dto.IssuerURL = new(generic.IssuerURL)
				dto.IntrospectionEndpoint = new(generic.IntrospectionEndpoint)
				dto.IntrospectionEndpointMethod = new(string(generic.IntrospectionEndpointMethod))
				dto.ClientID = generic.ClientID
				dto.UserInfoEndpoint = generic.UserInfoEndpoint
				dto.ClientSecret = generic.ClientSecret
			}
		case mcpproxy.IdentityProviderOAuth2Auth0:
			if auth0 := provider.Auth0; auth0 != nil {
				dto.Domain = new(auth0.Domain)
				dto.Audience = new(auth0.Audience)
			}
		case mcpproxy.IdentityProviderGraviteeAM:
		}
		dtos = append(dtos, dto)
	}
	return dtos
}

// toMcpProxyFlowExecution drops the platform default (DEFAULT, match not required), which the
// platform omits from its responses.
func toMcpProxyFlowExecution(execution *v4.FlowExecution) *v4.FlowExecution {
	if execution == nil {
		return nil
	}
	if (execution.Mode == "" || execution.Mode == v4.FlowModeDefault) && !execution.MatchRequired {
		return nil
	}
	return new(*execution)
}

func toMcpProxyFlowDTOs(flows []mcpproxy.Flow) []McpProxyFlowDTO {
	if len(flows) == 0 {
		return nil
	}
	dtos := make([]McpProxyFlowDTO, 0, len(flows))
	for _, flow := range flows {
		dtos = append(dtos, McpProxyFlowDTO{
			Name:      flow.Name,
			Enabled:   flow.Enabled,
			Selectors: toMcpProxyFlowSelectorDTOs(flow.Selectors),
			Request:   toMcpProxyFlowStepDTOs(flow.Request),
			Response:  toMcpProxyFlowStepDTOs(flow.Response),
			Tags:      flow.Tags,
		})
	}
	return dtos
}

func toMcpProxyFlowSelectorDTOs(selectors []mcpproxy.FlowSelector) []McpProxyFlowSelectorDTO {
	if len(selectors) == 0 {
		return nil
	}
	dtos := make([]McpProxyFlowSelectorDTO, 0, len(selectors))
	for _, selector := range selectors {
		dto := McpProxyFlowSelectorDTO{Type: string(selector.Type)}
		switch selector.Type {
		case mcpproxy.FlowSelectorMCP:
			if selector.MCP != nil {
				dto.Methods = selector.MCP.Methods
			}
		case mcpproxy.FlowSelectorCondition:
			if selector.Condition != nil {
				dto.Condition = new(selector.Condition.Condition)
			}
		}
		dtos = append(dtos, dto)
	}
	return dtos
}

func toMcpProxyFlowStepDTOs(steps []mcpproxy.FlowStep) []McpProxyFlowStepDTO {
	if len(steps) == 0 {
		return nil
	}
	dtos := make([]McpProxyFlowStepDTO, 0, len(steps))
	for _, step := range steps {
		dtos = append(dtos, McpProxyFlowStepDTO{
			Name:          step.Name,
			Policy:        step.Policy,
			Description:   step.Description,
			Enabled:       step.Enabled,
			Condition:     step.Condition,
			Configuration: step.Configuration,
		})
	}
	return dtos
}
