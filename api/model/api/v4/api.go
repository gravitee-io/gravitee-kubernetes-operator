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

package v4

import (
	"fmt"
	"reflect"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

// +kubebuilder:validation:Enum=PROXY;MESSAGE;NATIVE;
type ApiType string

// +kubebuilder:validation:Enum=PUBLISHED;UNPUBLISHED;
type ApiV4LifecycleState string

var _ core.ApiDefinitionModel = &Api{}

type Api struct {
	*base.ApiBase `json:",inline"`
	// API description
	// +kubebuilder:validation:Optional
	Description *string `json:"description,omitempty"`
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
	// A map of plan identifiers to plan
	// Keys uniquely identify plans and are used to keep them in sync
	// when using a management context.
	// +kubebuilder:validation:Optional
	Plans *map[string]*Plan `json:"plans,omitempty"`
	// API Flow Execution (Not applicable for Native API)
	FlowExecution *FlowExecution `json:"flowExecution,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	// List of flows for the API
	Flows []*Flow `json:"flows"`
	// API Analytics (Not applicable for Native API)
	Analytics *Analytics `json:"analytics,omitempty"`
	// API Services (Not applicable for Native API)
	Services *ApiServices `json:"services,omitempty"`
	// A list of Response Templates for the API (Not applicable for Native API)
	// +kubebuilder:validation:Optional
	ResponseTemplates *map[string]map[string]*base.ResponseTemplate `json:"responseTemplates,omitempty"`
	// List of members associated with the API
	// +kubebuilder:validation:Optional
	Members []*base.Member `json:"members,omitempty"`
	// +kubebuilder:validation:Optional
	// A map of pages objects.
	//
	// Keys uniquely identify pages and are used to keep them in sync
	// with APIM when using a management context.
	//
	// Renaming a key is the equivalent of deleting the page and recreating
	// it holding a new ID in APIM.
	Pages *map[string]*Page `json:"pages"`
	// API Failover
	Failover *Failover `json:"failover,omitempty"`
}

func (api *Api) GetType() string {
	return string(api.Type)
}

func (api *Api) GetGroupRefs() []core.ObjectRef {
	refs := make([]core.ObjectRef, 0)
	for i := range api.GroupRefs {
		refs = append(refs, &api.GroupRefs[i])
	}
	return refs
}

func (api *Api) GetGroups() []string {
	return api.Groups
}

func (api *Api) SetGroups(groups []string) {
	api.Groups = groups
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

var _ core.DefinitionContext = &DefinitionContext{}

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
	// +kubebuilder:default:=`MANAGEMENT`
	// +kubebuilder:validation:Enum=KUBERNETES;MANAGEMENT;
	SyncFrom DefinitionContextOrigin `json:"syncFrom,omitempty"`
}

type Failover struct {
	// API Failover is enabled?
	// +kubebuilder:default:=false
	Enabled *bool `json:"enabled,omitempty"`
	// API Failover max retires
	// +kubebuilder:default:=2
	MaxRetries *int `json:"maxRetries,omitempty"`
	// API Failover slow call duration
	// +kubebuilder:default:=2000
	SlowCallDuration *int64 `json:"slowCallDuration,omitempty"`
	// API Failover  open state duration
	// +kubebuilder:default:=10000
	OpenStateDuration *int64 `json:"openStateDuration,omitempty"`
	// API Failover max failures
	// +kubebuilder:default:=5
	MaxFailures *int `json:"maxFailures,omitempty"`
	// API Failover  per subscription
	// +kubebuilder:default:=true
	PerSubscription *bool `json:"perSubscription,omitempty"`
}

func NewDefaultKubernetesContext() *DefinitionContext {
	return &DefinitionContext{
		Origin:   OriginKubernetes,
		SyncFrom: OriginManagement,
	}
}

func (ctx *DefinitionContext) MergeWith(rhs core.DefinitionContext) *DefinitionContext {
	if reflect.ValueOf(rhs).IsNil() {
		return ctx
	}
	if impl, ok := rhs.(*DefinitionContext); ok {
		if ctx == nil {
			return impl
		}
		lhs := new(DefinitionContext)
		ctx.DeepCopyInto(lhs)
		if impl.Origin != "" {
			lhs.Origin = impl.Origin
		}
		if impl.SyncFrom != "" {
			lhs.SyncFrom = impl.SyncFrom
		}
		return lhs
	}
	return ctx
}

func (ctx *DefinitionContext) GetOrigin() string {
	if ctx == nil {
		return string(OriginKubernetes)
	}
	return string(ctx.Origin)
}

func (ctx *DefinitionContext) SetOrigin(origin string) {
	if ctx != nil {
		ctx.Origin = DefinitionContextOrigin(origin)
	}
}

func (api *Api) GetDefinitionContext() core.DefinitionContext {
	return api.DefinitionContext
}

func (api *Api) SetDefinitionContext(ctx core.DefinitionContext) {
	if impl, ok := ctx.(*DefinitionContext); ok {
		api.DefinitionContext = impl
	}
}

func (api *Api) GetState() string {
	return string(api.State)
}

func (api *Api) HasPlans() bool {
	return api.Plans != nil && len(*api.Plans) > 0
}

func (api *Api) GetPlan(name string) core.PlanModel {
	if api.Plans == nil {
		return nil
	}
	return (*api.Plans)[name]
}

func (api *Api) IsStopped() bool {
	return api.State == base.StateStopped
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
	if api.Plans != nil {
		for name, plan := range *api.Plans {
			plans = append(plans, plan.ToGatewayDefinition(name))
		}
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

func (api *Api) GetDefinitionVersion() core.ApiDefinitionVersion {
	return core.ApiV4
}

func (api *Api) GetContextPaths() []string {
	paths := make([]string, 0)
	for _, l := range api.Listeners {
		paths = append(paths, parseListener(l)...)
	}
	return paths
}

func parseListener(l Listener) []string {
	if l == nil {
		return []string{}
	}

	switch t := l.(type) {
	case *GenericListener:
		return parseListener(t.ToListener())
	case *HttpListener:
		{
			paths := make([]string, 0)
			for _, path := range t.Paths {
				if path.Host != "" {
					p := fmt.Sprintf("%s/%s", path.Host, path.Path)
					paths = append(paths, p)
				} else {
					paths = append(paths, path.Path)
				}
			}
			return paths
		}
	case *TCPListener:
		return t.Hosts
	}

	return []string{}
}

func (api *Api) GetAllSharedPolicyGroups() []*refs.NamespacedName {
	var results []*refs.NamespacedName

	if api.Flows != nil {
		results = append(results, getFLowSharedPolicyGroupsReferences(api.Flows)...)
	}

	if api.Plans != nil {
		for _, plan := range *api.Plans {
			if plan.Flows != nil {
				results = append(results, getFLowSharedPolicyGroupsReferences(plan.Flows)...)
			}
		}
	}

	return results
}

//nolint:gocognit // acceptable complexity
func getFLowSharedPolicyGroupsReferences(flows []*Flow) []*refs.NamespacedName {
	var results []*refs.NamespacedName

	for _, flow := range flows {
		for _, flowStep := range flow.Request {
			if flowStep.SharedPolicyGroup != nil {
				results = append(results, flowStep.SharedPolicyGroup)
			}
		}
		for _, flowStep := range flow.Response {
			if flowStep.SharedPolicyGroup != nil {
				results = append(results, flowStep.SharedPolicyGroup)
			}
		}
		for _, flowStep := range flow.Connect {
			if flowStep.SharedPolicyGroup != nil {
				results = append(results, flowStep.SharedPolicyGroup)
			}
		}
		for _, flowStep := range flow.Interact {
			if flowStep.SharedPolicyGroup != nil {
				results = append(results, flowStep.SharedPolicyGroup)
			}
		}
		for _, flowStep := range flow.Publish {
			if flowStep.SharedPolicyGroup != nil {
				results = append(results, flowStep.SharedPolicyGroup)
			}
		}
		for _, flowStep := range flow.Subscribe {
			if flowStep.SharedPolicyGroup != nil {
				results = append(results, flowStep.SharedPolicyGroup)
			}
		}
	}

	return results
}
