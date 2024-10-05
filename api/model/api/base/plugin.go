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

import "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"

type Plugin struct {
	// Plugin Policy
	// +kubebuilder:validation:Optional
	Policy string `json:"policy,omitempty"`
	// Plugin Resource
	// +kubebuilder:validation:Optional
	Resource string `json:"resource,omitempty"`
	// Plugin Configuration, a map of arbitrary key-values
	// +kubebuilder:validation:Optional
	Configuration *utils.GenericStringMap `json:"configuration,omitempty"`
}

type PluginReference struct {
	// Plugin Reference Namespace
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
	// Plugin Reference Resource
	// +kubebuilder:validation:Optional
	Resource string `json:"resource,omitempty"`
	// Plugin Reference Name
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
}

type PluginRevision struct {
	// Plugin reference
	PluginReference *PluginReference `json:"pluginReference,omitempty"`
	// Plugin Generation
	// +kubebuilder:validation:Optional
	Generation int64 `json:"generation,omitempty"`
	// Plugin
	Plugin *Plugin `json:"plugin,omitempty"`
	// Plugin Revision Hash code
	// +kubebuilder:validation:Optional
	HashCode string `json:"hashCode,omitempty"`
}
