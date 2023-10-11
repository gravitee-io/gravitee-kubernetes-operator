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
	Policy        string                  `json:"policy,omitempty"`
	Resource      string                  `json:"resource,omitempty"`
	Configuration *utils.GenericStringMap `json:"configuration,omitempty"`
}

type PluginReference struct {
	Namespace string `json:"namespace,omitempty"`
	Resource  string `json:"resource,omitempty"`
	Name      string `json:"name,omitempty"`
}

type PluginRevision struct {
	PluginReference *PluginReference `json:"pluginReference,omitempty"`
	Generation      int64            `json:"generation,omitempty"`
	Plugin          *Plugin          `json:"plugin,omitempty"`
	HashCode        string           `json:"hashCode,omitempty"`
}
