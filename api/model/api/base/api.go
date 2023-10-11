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
package base

type ApiBase struct {
	ID          string `json:"id,omitempty"`
	CrossID     string `json:"crossId,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	DeployedAt  uint64 `json:"deployedAt,omitempty"`
	// +kubebuilder:validation:Required
	// The definition context is used to inform a management API instance that this API definition
	// is managed using a kubernetes operator
	DefinitionContext *DefinitionContext `json:"definition_context,omitempty"`
	// +kubebuilder:default:=`STARTED`
	// +kubebuilder:validation:Enum=STARTED;STOPPED;
	State string `json:"state,omitempty"`
	// +kubebuilder:default:=`CREATED`
	LifecycleState LifecycleState `json:"lifecycle_state,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	Labels         []string       `json:"labels,omitempty"`
	// +kubebuilder:default:=PRIVATE
	Visibility        ApiVisibility                           `json:"visibility,omitempty"`
	PrimaryOwner      *Member                                 `json:"primaryOwner,omitempty"`
	Properties        []*Property                             `json:"properties,omitempty"`
	Metadata          []*MetadataEntry                        `json:"metadata,omitempty"`
	ResponseTemplates map[string]map[string]*ResponseTemplate `json:"response_templates,omitempty"`
	Resources         []*ResourceOrRef                        `json:"resources,omitempty"`
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

// +kubebuilder:validation:Enum=PUBLIC;PRIVATE;
type ApiVisibility string

type DefinitionVersion string

const (
	DefinitionVersionV1 DefinitionVersion = "1.0.0"
	DefinitionVersionV2 DefinitionVersion = "2.0.0"
	DefinitionVersionV4 DefinitionVersion = "4.0.0"
)

// +kubebuilder:validation:Enum=CREATED;PUBLISHED;UNPUBLISHED;DEPRECATED;ARCHIVED;
type LifecycleState string

const (
	StateStarted string = "STARTED"
	StateStopped string = "STOPPED"
)

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

type ResponseTemplate struct {
	StatusCode int               `json:"status,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}
