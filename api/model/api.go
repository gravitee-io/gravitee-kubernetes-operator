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
package model

type Api struct {
	Description string `json:"description,omitempty"`
	DeployedAt  uint64 `json:"deployedAt,omitempty"`
	// The definition context is used to inform a management API instance that this API definition
	// is managed using a kubernetes operator
	DefinitionContext *DefinitionContext `json:"definition_context,omitempty"`

	// io.gravitee.definition.model.Api
	ID      string `json:"id,omitempty"`
	CrossID string `json:"crossId,omitempty"`
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Required
	Version string `json:"version,omitempty"`
	// +kubebuilder:default:=DEFAULT
	FlowMode FlowMode `json:"flow_mode,omitempty"`
	// +kubebuilder:default:=`2.0.0`
	DefinitionVersion DefinitionVersion `json:"gravitee,omitempty"`
	// +kubebuilder:default:=`STARTED`
	// +kubebuilder:validation:Enum=STARTED;STOPPED;
	State string `json:"state,omitempty"`
	// +kubebuilder:default:=`CREATED`
	LifecycleState LifecycleState `json:"lifecycle_state,omitempty"`
	// +kubebuilder:validation:Required
	Proxy             *Proxy                                  `json:"proxy,omitempty"`
	Services          *Services                               `json:"services,omitempty"`
	Resources         []*ResourceOrRef                        `json:"resources,omitempty"`
	Flows             []Flow                                  `json:"flows,omitempty"`
	Properties        []*Property                             `json:"properties,omitempty"`
	Tags              []string                                `json:"tags,omitempty"`
	Labels            []string                                `json:"labels,omitempty"`
	PathMappings      []string                                `json:"path_mappings,omitempty"`
	ResponseTemplates map[string]map[string]*ResponseTemplate `json:"response_templates,omitempty"`
	Plans             []*Plan                                 `json:"plans,omitempty"`
	// +kubebuilder:default:=PRIVATE
	Visibility   ApiVisibility `json:"visibility,omitempty"`
	Metadata     []*Metadata   `json:"metadata,omitempty"`
	PrimaryOwner *Member       `json:"primaryOwner,omitempty"`
	// +kubebuilder:default:=v4-emulation-engine
	ExecutionMode string `json:"execution_mode,omitempty"`
	// local defines if the api is local or not.
	//
	// If true, the Operator will create the ConfigMaps for the Gateway and pushes the API to the Management API
	// but without setting the update flag in the datastore.
	//
	// If false, the Operator will not create the ConfigMaps for the Gateway.
	// Instead, it pushes the API to the Management API and forces it to update the event in the datastore.
	// This will cause Gateways to fetch the APIs from the datastore
	//
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	IsLocal bool `json:"local"`
}

func NewApi() *Api {
	return &Api{
		FlowMode: "DEFAULT",
	}
}

// +kubebuilder:validation:Enum=PUBLIC;PRIVATE;
type ApiVisibility string

// +kubebuilder:validation:Enum=DEFAULT;BEST_MATCH;
type FlowMode string

const (
	BestMatchFlowMode = FlowMode("BEST_MATCH")
	DefaultFlowMode   = FlowMode("DEFAULT")
)

type DefinitionVersion string

const (
	V1 DefinitionVersion = "1.0.0"
	V2 DefinitionVersion = "2.0.0"
)

// +kubebuilder:validation:Enum=CREATED;PUBLISHED;UNPUBLISHED;DEPRECATED;ARCHIVED;
type LifecycleState string

const (
	StateStarted string = "STARTED"
	StateStopped string = "STOPPED"
)

type Resource struct {
	// +kubebuilder:validation:Optional
	Enabled       bool              `json:"enabled"`
	Name          string            `json:"name,omitempty"`
	ResourceType  string            `json:"type,omitempty"`
	Configuration *GenericStringMap `json:"configuration,omitempty"`
}

type ResourceOrRef struct {
	*Resource `json:",omitempty,inline"`
	Ref       *NamespacedName `json:"ref,omitempty"`
}

func (r *ResourceOrRef) IsRef() bool {
	return r.Ref != nil
}

func (r *ResourceOrRef) IsMatchingRef(name, namespace string) bool {
	return r.IsRef() && r.Ref.Name == name && r.Ref.Namespace == namespace
}

func NewResource() *Resource {
	return &Resource{
		Enabled: true,
	}
}

const (
	ModeFullyManaged = "fully_managed"
	OriginKubernetes = "kubernetes"
)

type DefinitionContext struct {
	// +kubebuilder:default:=kubernetes
	Origin string `json:"origin,omitempty"`
	// +kubebuilder:default:=fully_managed
	Mode string `json:"mode,omitempty"`
}

// +kubebuilder:validation:Enum=STRING;NUMERIC;BOOLEAN;DATE;MAIL;URL;
type MetadataFormat string
type Metadata struct {
	Key          string         `json:"key"`
	Name         string         `json:"name"`
	Format       MetadataFormat `json:"format"`
	Value        string         `json:"value,omitempty"`
	DefaultValue string         `json:"defaultValue,omitempty"`
}

func NewDefinitionContext() *DefinitionContext {
	return &DefinitionContext{
		Origin: "kubernetes",
	}
}

type Member struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
}
