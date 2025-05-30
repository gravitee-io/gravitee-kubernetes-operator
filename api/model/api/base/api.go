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

package base

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

type ApiBase struct {
	// The API ID. If empty, this field will take the value of the `metadata.uid`
	// field of the resource.
	ID string `json:"id,omitempty"`
	// When promoting an API from one environment to the other,
	// this ID identifies the API across those different environments.
	// Setting this ID also allows to take control over an existing API on an APIM instance
	// (by setting the same value as defined in APIM).
	// If empty, a UUID will be generated based on the namespace and name of the resource.
	CrossID string `json:"crossId,omitempty"`
	// API name
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	// API version
	Version string `json:"version"`
	// +kubebuilder:default:=`STARTED`
	// The state of API (setting the value to `STOPPED` will make the API un-reachable from the gateway)
	State ApiState `json:"state,omitempty"`
	// +kubebuilder:validation:Optional
	// List of Tags of the API
	Tags []string `json:"tags"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	// List of labels of the API
	Labels []string `json:"labels"`
	// +kubebuilder:default:=PRIVATE
	// Should the API be publicly available from the portal or not ?
	Visibility ApiVisibility `json:"visibility,omitempty"`
	// +kubebuilder:validation:Optional
	// List of Properties for the API
	// +kubebuilder:default:={}
	Properties []*Property `json:"properties"`
	// +kubebuilder:validation:Optional
	// List of API metadata entries
	// +kubebuilder:default:={}
	Metadata []*MetadataEntry `json:"metadata"`
	// +kubebuilder:validation:Optional
	// Resources can be either inlined or reference the namespace and name
	// of an <a href="#apiresource">existing API resource definition</a>.
	// +kubebuilder:default:={}
	Resources []*ResourceOrRef `json:"resources"`
	// List of groups associated with the API.
	// This groups are id or name references to existing groups in APIM.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Groups []string `json:"groups"`
	// List of group references associated with the API
	// These groups are references to Group custom resources created on the cluster.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	GroupRefs []refs.NamespacedName `json:"groupRefs"`
	// +kubebuilder:validation:Optional
	// The list of categories the API belongs to.
	// Categories are reflected in APIM portal so that consumers can easily find the APIs they need.
	// +kubebuilder:default:={}
	Categories []string `json:"categories"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	// If true, new members added to the API spec will
	// be notified when the API is synced with APIM.
	NotifyMembers bool `json:"notifyMembers"`
	// References to Notification custom resources to setup notifications.
	// For an API Notification CRD `eventType` field must be set to `api`
	// and only events set via `apiEvents` attributes are used.
	// Only one notification with `target` equals to `console` is admitted.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	NotificationsRefs []refs.NamespacedName `json:"notificationsRefs,omitempty"`
	// ConsoleNotification struct sent to the mAPI, not part of the CRD spec.
	// +kubebuilder:skipversion
	ConsoleNotification *ConsoleNotificationConfiguration `json:"consoleNotificationConfiguration,omitempty"`
}

// GetResources implements core.ApiDefinitionModel.
func (api *ApiBase) GetResources() []core.ObjectOrRef[core.ResourceModel] {
	refs := make([]core.ObjectOrRef[core.ResourceModel], len(api.Resources))
	for i := range api.Resources {
		refs[i] = api.Resources[i]
	}
	return refs
}

// +kubebuilder:validation:Enum=PUBLIC;PRIVATE;
type ApiVisibility string

type DefinitionVersion string

const (
	DefinitionVersionV1 DefinitionVersion = "1.0.0"
	DefinitionVersionV2 DefinitionVersion = "2.0.0"
	DefinitionVersionV4 DefinitionVersion = "V4"
	GatewayDefinitionV4 DefinitionVersion = "4.0.0"
)

// +kubebuilder:validation:Enum=CREATED;PUBLISHED;UNPUBLISHED;DEPRECATED;ARCHIVED;
type LifecycleState string

// +kubebuilder:validation:Enum=STARTED;STOPPED;
type ApiState string

const (
	StateStarted ApiState = "STARTED"
	StateStopped ApiState = "STOPPED"
)

type ResponseTemplate struct {
	// +kubebuilder:validation:Optional
	StatusCode *int `json:"status,omitempty"`
	// +kubebuilder:validation:Optional
	Headers *map[string]string `json:"headers,omitempty"`
	// +kubebuilder:validation:Optional
	Body *string `json:"body,omitempty"`
	// Propagate error key to logs
	PropagateErrorKeyToLogs *bool `json:"propagateErrorKeyToLogs,omitempty"`
}

func (api *ApiBase) GetName() string {
	return api.Name
}
