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
	Description       string             `json:"description,omitempty"`
	DeployedAt        uint64             `json:"deployedAt,omitempty"`
	DefinitionContext *DefinitionContext `json:"definition_context,omitempty"`

	// io.gravitee.definition.model.Api
	Id      string `json:"id,omitempty"`
	CrossId string `json:"crossId,omitempty"`
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Required
	Version string `json:"version,omitempty"`
	// +kubebuilder:default:=DEFAULT
	FlowMode FlowMode `json:"flow_mode,omitempty"`
	// +kubebuilder:default:=`2.0.0`
	DefinitionVersion DefinitionVersion `json:"gravitee,omitempty"`
	// +kubebuilder:default:=`STARTED`
	State State `json:"state,omitempty"`
	// +kubebuilder:default:=`CREATED`
	LifecycleState LifecycleState `json:"lifecycle_state,omitempty"`
	// +kubebuilder:validation:Required
	Proxy     *Proxy      `json:"proxy,omitempty"`
	Services  *Services   `json:"services,omitempty"`
	Resources []*Resource `json:"resources,omitempty"`
	//	Paths             map[string][]interface{}                `json:"paths,omitempty"` // Different from Java
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

type DefinitionVersion string

const (
	V1 DefinitionVersion = "1.0.0"
	V2 DefinitionVersion = "2.0.0"
)

// +kubebuilder:validation:Enum=CREATED;PUBLISHED;UNPUBLISHED;DEPRECATED;ARCHIVED;
type LifecycleState string

// +kubebuilder:validation:Enum=STARTED;STOPPED;
type State string

const (
	StateStarted State = "STARTED"
	StateStopped State = "STOPPED"
)

type Resource struct {
	Enabled       bool              `json:"enabled,omitempty"`
	Name          string            `json:"name,omitempty"`
	ResourceType  string            `json:"type,omitempty"`
	Configuration *GenericStringMap `json:"configuration,omitempty"`
}

func NewResource() *Resource {
	return &Resource{
		Enabled: true,
	}
}

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
	Value        string         `json:"value"`
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
