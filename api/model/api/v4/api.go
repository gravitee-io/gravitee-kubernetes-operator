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

// +kubebuilder:object:generate=true
package v4

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
)

// +kubebuilder:validation:Enum=PROXY;MESSAGE;
type ApiType string

// +kubebuilder:validation:Enum=PUBLISHED;UNPUBLISHED;
type ApiV4LifecycleState string

type Api struct {
	*base.ApiBase `json:",inline"`
	// +kubebuilder:default:=`V4`
	// +kubebuilder:validation:Enum=`V4`;
	// The definition version of the API.
	DefinitionVersion base.DefinitionVersion `json:"definitionVersion,omitempty"`
	// The API Definition context is used to identify the Kubernetes origin of the API,
	// and define whether the API definition should be synchronized
	// from an API instance or from a config map created in the cluster (which is the default)
	DefinitionContext *DefinitionContext `json:"definitionContext,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=`UNPUBLISHED`
	// API life cycle state can be one of the values PUBLISHED, UNPUBLISHED
	LifecycleState ApiV4LifecycleState `json:"lifecycleState,omitempty"`
	// +kubebuilder:validation:Required
	// Api Type (proxy or message)
	Type ApiType `json:"type"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems:=1
	// List of listeners for this API
	Listeners []*GenericListener `json:"listeners"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems:=1
	// List of Endpoint groups
	EndpointGroups []*EndpointGroup `json:"endpointGroups"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinProperties:=1
	// A map of plan identifiers to plan
	// Keys uniquely identify plans and are used to keep them in sync
	// when using a management context.
	Plans map[string]*Plan `json:"plans"`
	// API Flow Execution
	FlowExecution *FlowExecution `json:"flowExecution,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	// List of flows for the API
	Flows []*Flow `json:"flows"`
	// API Analytics
	Analytics *Analytics `json:"analytics,omitempty"`
	// API Services
	Services *ApiServices `json:"services,omitempty"`
	// A list of Response Templates for the API
	ResponseTemplates map[string]map[string]*base.ResponseTemplate `json:"responseTemplates,omitempty"`
	// List of members associated with the API
	Members []*base.Member `json:"members,omitempty"`
	// +kubebuilder:validation:Optional
	// A map of pages objects.
	//
	// Keys uniquely identify pages and are used to keep them in sync
	// with APIM when using a management context.
	//
	// Renaming a key is the equivalent of deleting the page and recreating
	// it holding a new ID in APIM.
	Pages map[string]*Page `json:"pages"`
}

type GatewayDefinitionApi struct {
	*Api    `json:",inline"`
	Version string                   `json:"apiVersion"`
	Plans   []*GatewayDefinitionPlan `json:"plans"`
}

// +kubebuilder:validation:Enum=FULLY_MANAGED;
type DefinitionContextMode string

type DefinitionContextOrigin string

const (
	ModeFullyManaged DefinitionContextOrigin = "FULLY_MANAGED"
	OriginKubernetes DefinitionContextOrigin = "KUBERNETES"
	OriginManagement DefinitionContextOrigin = "MANAGEMENT"
)

type DefinitionContext struct {
	// The definition context origin where the API definition is managed.
	// The value is always `KUBERNETES` for an API managed by the operator.
	// +kubebuilder:validation:Enum=KUBERNETES;
	// +kubebuilder:default:=`KUBERNETES`
	Origin DefinitionContextOrigin `json:"origin,omitempty"`
	// The syncFrom field defines where the gateways should source the API definition from.
	// If the value is `MANAGEMENT`, the API definition will be sourced from an APIM instance.
	// This means that the API definition *must* hold a context reference in that case.
	// Setting the value to `MANAGEMENT` allows to make an API definition available on
	// gateways deployed across multiple clusters / regions.
	// If the value is `KUBERNETES`, the API definition will be sourced from a config map.
	// This means that only gateways deployed in the same cluster will be able to sync the API definition.
	// +kubebuilder:default:=`KUBERNETES`
	// +kubebuilder:validation:Enum=KUBERNETES;MANAGEMENT;
	SyncFrom DefinitionContextOrigin `json:"syncFrom,omitempty"`
}

func NewDefaultKubernetesContext() *DefinitionContext {
	return &DefinitionContext{
		Origin:   OriginKubernetes,
		SyncFrom: OriginKubernetes,
	}
}

func (ctx DefinitionContext) MergeWith(rhs *DefinitionContext) *DefinitionContext {
	lhs := new(DefinitionContext)
	ctx.DeepCopyInto(lhs)
	if rhs == nil {
		return lhs
	}
	if rhs.Origin != "" {
		lhs.Origin = rhs.Origin
	}
	if rhs.SyncFrom != "" {
		lhs.SyncFrom = rhs.SyncFrom
	}
	return lhs
}

// Converts the API to its gateway definition equivalent.
func (api *Api) ToGatewayDefinition() GatewayDefinitionApi {
	def := GatewayDefinitionApi{Api: api}
	def.Version = api.Version
	def.DefinitionVersion = base.GatewayDefinitionV4
	def.Type = ApiType(Enum(api.Type).ToGatewayDefinition())
	def.Listeners = api.getGatewayDefinitionListener()
	def.EndpointGroups = api.getGatewayDefinitionEndpointGroups()
	def.Plans = api.getGatewayDefinitionPlans()
	def.Flows = api.getApiDefinitionFlows()
	if api.FlowExecution != nil {
		api.FlowExecution.Mode = FlowMode(Enum(api.FlowExecution.Mode).ToGatewayDefinition())
	}
	return def
}

func (api *Api) getGatewayDefinitionPlans() []*GatewayDefinitionPlan {
	plans := make([]*GatewayDefinitionPlan, 0)
	for name, plan := range api.Plans {
		plans = append(plans, plan.ToGatewayDefinition(name))
	}
	return plans
}

func (api *Api) getGatewayDefinitionListener() []*GenericListener {
	listeners := make([]*GenericListener, len(api.Listeners))
	for i, listener := range api.Listeners {
		listeners[i] = ToListenerGatewayDefinition(listener)
	}
	return listeners
}

func (api *Api) getApiDefinitionFlows() []*Flow {
	flows := make([]*Flow, len(api.Flows))
	for i, flow := range api.Flows {
		flows[i] = flow.ToGatewayDefinition()
	}
	return flows
}

func (api *Api) getGatewayDefinitionEndpointGroups() []*EndpointGroup {
	endpointGroups := make([]*EndpointGroup, len(api.EndpointGroups))
	for i, endpointGroup := range api.EndpointGroups {
		endpointGroups[i] = endpointGroup.ToGatewayDefinition()
	}
	return endpointGroups
}
