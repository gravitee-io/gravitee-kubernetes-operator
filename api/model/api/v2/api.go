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

// +kubebuilder:object:generate=true
package v2

import (
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

var _ core.ApiDefinitionModel = &Api{}
var _ core.DefinitionContext = &DefinitionContext{}

type Api struct {
	*base.ApiBase `json:",inline"`
	// API description
	Description string `json:"description"`
	// +kubebuilder:validation:Required
	// The definition context is used to inform a management API instance that this API definition
	// is managed using a kubernetes operator
	DefinitionContext *DefinitionContext `json:"definition_context,omitempty"`
	// +kubebuilder:default:=`CREATED`
	// API life cycle state can be one of the values CREATED, PUBLISHED, UNPUBLISHED, DEPRECATED, ARCHIVED
	LifecycleState base.LifecycleState `json:"lifecycle_state,omitempty"`
	// Shows the time that the API is deployed
	// +kubebuilder:validation:Optional
	DeployedAt *uint64 `json:"deployedAt,omitempty"`
	// +kubebuilder:default:=`2.0.0`
	// The definition version of the API. For v1alpha1 resources, this field should always set to `2.0.0`.
	DefinitionVersion base.DefinitionVersion `json:"gravitee,omitempty"`
	// +kubebuilder:default:=DEFAULT
	// The flow mode of the API. The value is either `DEFAULT` or `BEST_MATCH`.
	FlowMode FlowMode `json:"flow_mode,omitempty"`
	// The proxy of the API that specifies its VirtualHosts and Groups.
	Proxy *Proxy `json:"proxy,omitempty"`
	// Contains different services for the API (EndpointDiscovery, HealthCheck ...)
	Services *Services `json:"services,omitempty"`
	// +kubebuilder:validation:Optional
	// The flow of the API
	Flows []Flow `json:"flows"`
	// +kubebuilder:validation:Optional
	// API Path mapping
	PathMappings []string `json:"path_mappings"`
	// +kubebuilder:validation:Optional
	// API plans
	Plans []*Plan `json:"plans"`
	// A list of Response Templates for the API
	// +kubebuilder:validation:Optional
	ResponseTemplates *map[string]map[string]*base.ResponseTemplate `json:"response_templates,omitempty"`
	// List of members associated with the API
	// +kubebuilder:validation:Optional
	Members []*base.Member `json:"members"`
	// +kubebuilder:validation:Optional
	// A map of pages objects.
	//
	// Keys uniquely identify pages and are used to keep them in sync
	// with APIM when using a management context.
	//
	// Renaming a key is the equivalent of deleting the page and recreating
	// it holding a new ID in APIM.
	Pages map[string]*Page `json:"pages"`

	// Execution mode that eventually runs the API in the gateway
	// +kubebuilder:default:=v4-emulation-engine
	// +kubebuilder:validation:Enum=v3;v4-emulation-engine;
	ExecutionMode string `json:"execution_mode,omitempty"`
}

func (api *Api) GetDefinitionVersion() core.ApiDefinitionVersion {
	return core.ApiV2
}

func (api *Api) GetState() string {
	return string(api.State)
}

func (api *Api) HasPlans() bool {
	return len(api.Plans) > 0
}

func (api *Api) GetContextPaths() []string {
	paths := make([]string, 0)
	proxy := api.Proxy
	if proxy == nil {
		return paths
	}
	for _, vh := range proxy.VirtualHosts {
		host, path := vh.Host, vh.Path
		if host != "" {
			paths = append(paths, fmt.Sprintf("%s/%s", host, path))
		} else {
			paths = append(paths, path)
		}
	}

	return paths
}

const (
	ModeFullyManaged = "fully_managed"
	OriginKubernetes = "kubernetes"
	OriginManagement = "management"
)

type DefinitionContext struct {
	// +kubebuilder:default:=kubernetes
	Origin string `json:"origin,omitempty"`
	// +kubebuilder:default:=fully_managed
	Mode string `json:"mode,omitempty"`
	// +kubebuilder:default:=kubernetes
	SyncFrom string `json:"syncFrom,omitempty"`
}

func (api *Api) GetDefinitionContext() core.DefinitionContext {
	return api.DefinitionContext
}

func (api *Api) SetDefinitionContext(ctx core.DefinitionContext) {
	if impl, ok := ctx.(*DefinitionContext); ok {
		api.DefinitionContext = impl
	}
}

func (ctx *DefinitionContext) GetOrigin() string {
	if ctx == nil {
		return OriginKubernetes
	}
	return ctx.Origin
}

func (ctx *DefinitionContext) SetOrigin(origin string) {
	if ctx != nil {
		ctx.Origin = origin
	}
}
