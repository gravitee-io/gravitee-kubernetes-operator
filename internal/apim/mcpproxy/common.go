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

// Package mcpproxy maps the McpProxy CRD to the AIM MCP proxies of the APIM Automation API through the SDK.
package mcpproxy

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"
	"github.com/gravitee-io/gravitee-automation-sdk/common/pkg/unstructured"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/mcpproxy"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim"
)

// ToDTO maps the CRD onto the payload a PUT sends. It is the single mapping used by the sync path, the dry
// run and drift detection. SDK defaults are not applied. The CRD nests each union variant under its own key
// while the payload carries it flat with its type, so unions are mapped field by field, not through JSON.
// Studio references are resolved to the HRIDs of their CatalogMcpServers.
func ToDTO(obj *v1alpha1.McpProxy) (sdk.AimMcpProxy, error) {
	spec := obj.Spec.Type

	state := sdk.AimMcpProxyLifecycleState(spec.State)
	if state == "" {
		state = sdk.AimMcpProxyLifecycleStateSTARTED
	}

	flows, err := toFlows(spec.Flows)
	if err != nil {
		return sdk.AimMcpProxy{}, err
	}
	providers, err := toIdentityProviders(spec.IdentityProviders)
	if err != nil {
		return sdk.AimMcpProxy{}, err
	}
	plans, err := toPlans(spec.Plans)
	if err != nil {
		return sdk.AimMcpProxy{}, err
	}

	dto := sdk.AimMcpProxy{
		Hrid:              refs.NewNamespacedNameFromObject(obj).HRID(),
		EntityId:          spec.EntityID,
		Name:              spec.Name,
		Description:       spec.Description,
		ContextPath:       spec.ContextPath,
		Mode:              sdk.AimMcpProxyMode(spec.Mode),
		State:             &state,
		FlowExecution:     toFlowExecution(spec.FlowExecution),
		Flows:             flows,
		IdentityProviders: providers,
		Plans:             plans,
	}

	switch spec.Mode {
	case mcpproxy.ModeProxy:
		if spec.Proxy != nil {
			if dto.Proxy, err = toProxy(spec.Proxy); err != nil {
				return sdk.AimMcpProxy{}, err
			}
		}
	case mcpproxy.ModeStudio:
		if spec.Studio != nil {
			if dto.Studio, err = toStudio(spec.Studio, obj.GetNamespace()); err != nil {
				return sdk.AimMcpProxy{}, err
			}
		}
	}

	return Canonical(dto), nil
}

// Canonical sorts the collections the platform returns sorted, so that a mapped CRD and a remote state
// compare item by item. It sorts copies: the payload may share its slices with a state.
func Canonical(dto sdk.AimMcpProxy) sdk.AimMcpProxy {
	dto.Plans = sortedCopy(dto.Plans, func(a, b sdk.AimMcpProxyPlan) int {
		return cmp.Compare(a.Name, b.Name)
	})
	dto.IdentityProviders = sortedCopy(dto.IdentityProviders, func(a, b sdk.AimMcpProxyIdentityProvider) int {
		return cmp.Compare(identityProviderName(a), identityProviderName(b))
	})
	if dto.Studio != nil {
		studio := *dto.Studio
		studio.Tools = sortedCopy(studio.Tools, func(a, b sdk.AimMcpProxyStudioTool) int {
			return cmp.Or(cmp.Compare(a.Server, b.Server), cmp.Compare(a.Tool, b.Tool))
		})
		studio.UpstreamAuth = sortedCopy(studio.UpstreamAuth, func(a, b sdk.AimMcpProxyStudioUpstreamAuth) int {
			return cmp.Compare(a.Server, b.Server)
		})
		dto.Studio = &studio
	}
	return dto
}

// sortedCopy returns a sorted copy of s, or s itself when there is nothing to sort, so that a nil or an
// empty list keeps its JSON form.
func sortedCopy[S ~[]E, E any](s S, compare func(a, b E) int) S {
	if len(s) < 2 {
		return s
	}
	sorted := slices.Clone(s)
	slices.SortStableFunc(sorted, compare)
	return sorted
}

func identityProviderName(provider sdk.AimMcpProxyIdentityProvider) string {
	switch {
	case provider.GraviteeAm != nil:
		return provider.GraviteeAm.Name
	case provider.OAuth2Generic != nil:
		return provider.OAuth2Generic.Name
	case provider.OAuth2Auth0 != nil:
		return provider.OAuth2Auth0.Name
	}
	return ""
}

// toProxy leaves out an upstream auth of type NONE: the platform stores it as no authentication and omits it.
func toProxy(proxy *mcpproxy.Proxy) (*sdk.AimMcpProxyProxy, error) {
	dto := &sdk.AimMcpProxyProxy{ServerUrl: proxy.ServerURL}
	if auth := proxy.UpstreamAuth; auth != nil && auth.Type != mcpproxy.UpstreamAuthTypeNone {
		upstreamAuth, err := toUpstreamAuth(*auth)
		if err != nil {
			return nil, err
		}
		dto.UpstreamAuth = &upstreamAuth
	}
	return dto, nil
}

func toStudio(studio *mcpproxy.Studio, ns string) (*sdk.AimMcpProxyStudio, error) {
	dto := &sdk.AimMcpProxyStudio{
		Tools:        make([]sdk.AimMcpProxyStudioTool, 0, len(studio.Tools)),
		UpstreamAuth: make([]sdk.AimMcpProxyStudioUpstreamAuth, 0, len(studio.UpstreamAuth)),
		EnableFGA:    new(studio.EnableFGA),
	}
	for _, tool := range studio.Tools {
		dto.Tools = append(dto.Tools, sdk.AimMcpProxyStudioTool{
			Server: serverHRID(tool.ServerRef, ns),
			Tool:   tool.Tool,
			Alias:  tool.Alias,
		})
	}
	for _, entry := range studio.UpstreamAuth {
		auth, err := toUpstreamAuth(entry.Auth)
		if err != nil {
			return nil, err
		}
		dto.UpstreamAuth = append(dto.UpstreamAuth, sdk.AimMcpProxyStudioUpstreamAuth{
			Server: serverHRID(entry.ServerRef, ns),
			Auth:   auth,
		})
	}
	return dto, nil
}

// serverHRID is the HRID of a referenced CatalogMcpServer, derived like its own.
func serverHRID(ref refs.NamespacedName, ns string) string {
	server := mcpproxy.WithNamespace(ref, ns)
	return server.HRID()
}

// toUpstreamAuth maps an auth without a type to NONE.
func toUpstreamAuth(auth mcpproxy.UpstreamAuth) (sdk.AimMcpProxyUpstreamAuth, error) {
	switch auth.Type {
	case "", mcpproxy.UpstreamAuthTypeNone:
		return sdk.AimMcpProxyUpstreamAuth{
			None: &sdk.AimMcpProxyNoUpstreamAuth{Type: sdk.AimMcpProxyNoUpstreamAuthTypeNONE},
		}, nil
	case mcpproxy.UpstreamAuthTypeAPIKey:
		if auth.APIKey != nil {
			return sdk.AimMcpProxyUpstreamAuth{ApiKey: &sdk.AimMcpProxyApiKeyUpstreamAuth{
				Type:         sdk.AimMcpProxyApiKeyUpstreamAuthTypeAPIKEY,
				ApiKeyHeader: auth.APIKey.APIKeyHeader,
				ApiKey:       new(auth.APIKey.APIKey),
			}}, nil
		}
	case mcpproxy.UpstreamAuthTypeBearer:
		if auth.Bearer != nil {
			return sdk.AimMcpProxyUpstreamAuth{Bearer: &sdk.AimMcpProxyBearerUpstreamAuth{
				Type:  sdk.AimMcpProxyBearerUpstreamAuthTypeBEARER,
				Token: new(auth.Bearer.Token),
			}}, nil
		}
	case mcpproxy.UpstreamAuthTypeBasic:
		if auth.Basic != nil {
			return sdk.AimMcpProxyUpstreamAuth{Basic: &sdk.AimMcpProxyBasicUpstreamAuth{
				Type:     sdk.AimMcpProxyBasicUpstreamAuthTypeBASIC,
				Username: auth.Basic.Username,
				Password: new(auth.Basic.Password),
			}}, nil
		}
	case mcpproxy.UpstreamAuthTypeOAuth2:
		if auth.OAuth2 != nil {
			return sdk.AimMcpProxyUpstreamAuth{OAuth2: &sdk.AimMcpProxyOAuth2UpstreamAuth{
				Type:         sdk.AimMcpProxyOAuth2UpstreamAuthTypeOAUTH2,
				AuthorizeUrl: auth.OAuth2.AuthorizeURL,
				TokenUrl:     auth.OAuth2.TokenURL,
				ClientId:     auth.OAuth2.ClientID,
				ClientSecret: new(auth.OAuth2.ClientSecret),
				Scopes:       auth.OAuth2.Scopes,
			}}, nil
		}
	default:
		return sdk.AimMcpProxyUpstreamAuth{}, fmt.Errorf("unknown upstream auth type %s", auth.Type)
	}
	return sdk.AimMcpProxyUpstreamAuth{}, fmt.Errorf("upstream auth type %s has no block for it", auth.Type)
}

// toPlans always returns a list: plans are required.
func toPlans(plans []mcpproxy.Plan) ([]sdk.AimMcpProxyPlan, error) {
	dtos := make([]sdk.AimMcpProxyPlan, 0, len(plans))
	for _, plan := range plans {
		security, err := toPlanSecurity(plan.Security)
		if err != nil {
			return nil, fmt.Errorf("plan %s: %w", plan.Name, err)
		}
		flows, err := toFlows(plan.Flows)
		if err != nil {
			return nil, fmt.Errorf("plan %s: %w", plan.Name, err)
		}
		dtos = append(dtos, sdk.AimMcpProxyPlan{Name: plan.Name, Security: security, Flows: flows})
	}
	return dtos, nil
}

// toPlanSecurity sends propagateApiKey only when it is true, its default being false.
func toPlanSecurity(security mcpproxy.PlanSecurity) (sdk.AimMcpProxyPlanSecurity, error) {
	switch security.Type {
	case mcpproxy.PlanSecurityKeyLess:
		return sdk.AimMcpProxyPlanSecurity{
			KeyLess: &sdk.AimMcpProxyKeyLessPlanSecurity{Type: sdk.AimMcpProxyKeyLessPlanSecurityTypeKEYLESS},
		}, nil
	case mcpproxy.PlanSecurityAPIKey:
		if security.APIKey != nil {
			apiKey := &sdk.AimMcpProxyApiKeyPlanSecurity{
				Type:         sdk.AimMcpProxyApiKeyPlanSecurityTypeAPIKEY,
				Source:       sdk.AimMcpProxyApiKeyPlanSecuritySource(security.APIKey.Source),
				ApiKeyHeader: security.APIKey.APIKeyHeader,
			}
			if security.APIKey.PropagateAPIKey {
				apiKey.PropagateApiKey = new(true)
			}
			return sdk.AimMcpProxyPlanSecurity{ApiKey: apiKey}, nil
		}
	case mcpproxy.PlanSecurityOAuth2:
		if security.OAuth2 != nil {
			return sdk.AimMcpProxyPlanSecurity{OAuth2: &sdk.AimMcpProxyOAuth2PlanSecurity{
				Type:     sdk.AimMcpProxyOAuth2PlanSecurityTypeOAUTH2,
				Provider: security.OAuth2.Provider,
				Scopes:   security.OAuth2.Scopes,
			}}, nil
		}
	default:
		return sdk.AimMcpProxyPlanSecurity{}, fmt.Errorf("unknown security type %s", security.Type)
	}
	return sdk.AimMcpProxyPlanSecurity{}, fmt.Errorf("security type %s has no block for it", security.Type)
}

func toIdentityProviders(providers []mcpproxy.IdentityProvider) ([]sdk.AimMcpProxyIdentityProvider, error) {
	if len(providers) == 0 {
		return nil, nil
	}
	dtos := make([]sdk.AimMcpProxyIdentityProvider, 0, len(providers))
	for _, provider := range providers {
		dto, err := toIdentityProvider(provider)
		if err != nil {
			return nil, fmt.Errorf("identity provider %s: %w", provider.Name, err)
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

func toIdentityProvider(provider mcpproxy.IdentityProvider) (sdk.AimMcpProxyIdentityProvider, error) {
	switch provider.Type {
	case mcpproxy.IdentityProviderGraviteeAM:
		return sdk.AimMcpProxyIdentityProvider{GraviteeAm: &sdk.AimMcpProxyGraviteeAmIdentityProvider{
			Name: provider.Name,
			Type: sdk.AimMcpProxyGraviteeAmIdentityProviderTypeGRAVITEEAM,
		}}, nil
	case mcpproxy.IdentityProviderOAuth2Generic:
		if generic := provider.OAuth2Generic; generic != nil {
			return sdk.AimMcpProxyIdentityProvider{OAuth2Generic: &sdk.AimMcpProxyGenericIdentityProvider{
				Name:                  provider.Name,
				Type:                  sdk.AimMcpProxyGenericIdentityProviderTypeOAUTH2GENERIC,
				IssuerUrl:             generic.IssuerURL,
				IntrospectionEndpoint: generic.IntrospectionEndpoint,
				IntrospectionEndpointMethod: sdk.AimMcpProxyGenericIdentityProviderIntrospectionEndpointMethod(
					generic.IntrospectionEndpointMethod),
				ClientId:         generic.ClientID,
				UserInfoEndpoint: generic.UserInfoEndpoint,
				ClientSecret:     generic.ClientSecret,
			}}, nil
		}
	case mcpproxy.IdentityProviderOAuth2Auth0:
		if auth0 := provider.Auth0; auth0 != nil {
			return sdk.AimMcpProxyIdentityProvider{OAuth2Auth0: &sdk.AimMcpProxyAuth0IdentityProvider{
				Name:     provider.Name,
				Type:     sdk.AimMcpProxyAuth0IdentityProviderTypeOAUTH2AUTH0,
				Domain:   auth0.Domain,
				Audience: auth0.Audience,
			}}, nil
		}
	default:
		return sdk.AimMcpProxyIdentityProvider{}, fmt.Errorf("unknown type %s", provider.Type)
	}
	return sdk.AimMcpProxyIdentityProvider{}, fmt.Errorf("type %s has no block for it", provider.Type)
}

// toFlowExecution drops the platform default (DEFAULT, match not required), which the platform omits from
// its responses.
func toFlowExecution(execution *v4.FlowExecution) *sdk.AimMcpProxyFlowExecution {
	if execution == nil {
		return nil
	}
	if (execution.Mode == "" || execution.Mode == v4.FlowModeDefault) && !execution.MatchRequired {
		return nil
	}
	dto := &sdk.AimMcpProxyFlowExecution{MatchRequired: new(execution.MatchRequired)}
	if execution.Mode != "" {
		dto.Mode = new(sdk.AimMcpProxyFlowExecutionMode(execution.Mode))
	}
	return dto
}

func toFlows(flows []mcpproxy.Flow) ([]sdk.AimMcpProxyFlow, error) {
	if len(flows) == 0 {
		return nil, nil
	}
	dtos := make([]sdk.AimMcpProxyFlow, 0, len(flows))
	for _, flow := range flows {
		selectors, err := toFlowSelectors(flow.Selectors)
		if err != nil {
			return nil, fmt.Errorf("flow %s: %w", flow.Name, err)
		}
		dtos = append(dtos, sdk.AimMcpProxyFlow{
			Name:      flow.Name,
			Enabled:   new(flow.Enabled),
			Selectors: selectors,
			Request:   toFlowSteps(flow.Request),
			Response:  toFlowSteps(flow.Response),
			Tags:      flow.Tags,
		})
	}
	return dtos, nil
}

func toFlowSelectors(selectors []mcpproxy.FlowSelector) ([]sdk.AimMcpProxyFlowSelector, error) {
	if len(selectors) == 0 {
		return nil, nil
	}
	dtos := make([]sdk.AimMcpProxyFlowSelector, 0, len(selectors))
	for _, selector := range selectors {
		dto, err := toFlowSelector(selector)
		if err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

func toFlowSelector(selector mcpproxy.FlowSelector) (sdk.AimMcpProxyFlowSelector, error) {
	switch selector.Type {
	case mcpproxy.FlowSelectorMCP:
		if selector.MCP != nil {
			return sdk.AimMcpProxyFlowSelector{Mcp: &sdk.AimMcpProxyMcpFlowSelector{
				Type:    sdk.AimMcpProxyMcpFlowSelectorTypeMCP,
				Methods: selector.MCP.Methods,
			}}, nil
		}
	case mcpproxy.FlowSelectorCondition:
		if selector.Condition != nil {
			return sdk.AimMcpProxyFlowSelector{Condition: &sdk.AimMcpProxyConditionFlowSelector{
				Type:      sdk.AimMcpProxyConditionFlowSelectorTypeCONDITION,
				Condition: selector.Condition.Condition,
			}}, nil
		}
	default:
		return sdk.AimMcpProxyFlowSelector{}, fmt.Errorf("unknown selector type %s", selector.Type)
	}
	return sdk.AimMcpProxyFlowSelector{}, fmt.Errorf("selector type %s has no block for it", selector.Type)
}

func toFlowSteps(steps []mcpproxy.FlowStep) []sdk.AimMcpProxyFlowStep {
	if len(steps) == 0 {
		return nil
	}
	dtos := make([]sdk.AimMcpProxyFlowStep, 0, len(steps))
	for _, step := range steps {
		dtos = append(dtos, sdk.AimMcpProxyFlowStep{
			Name:          step.Name,
			Policy:        step.Policy,
			Description:   step.Description,
			Enabled:       new(step.Enabled),
			Condition:     step.Condition,
			Configuration: toConfiguration(step.Configuration),
		})
	}
	return dtos
}

func toConfiguration(configuration *utils.GenericStringMap) *unstructured.Unstructured {
	if configuration == nil {
		return nil
	}
	return unstructured.From(configuration.DeepCopy().Object)
}

// CreateAPIMClient builds the SDK client of the management context the proxy references.
func CreateAPIMClient(ctx context.Context, obj *v1alpha1.McpProxy) (*apim.SDKClient, error) {
	return apim.SDKFromContextRef(ctx, obj.ContextRef(), obj.GetNamespace())
}
