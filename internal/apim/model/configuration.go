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

package model

type IDPConfiguration struct {
	ID                 string                 `json:"id,omitempty"`
	Name               string                 `json:"name,omitempty"`
	Description        string                 `json:"description,omitempty"`
	Type               string                 `json:"type,omitempty"` // GOOGLE, GITHUB, GRAVITEEIO_AM, OIDC
	Enabled            bool                   `json:"enabled,omitempty"`
	Configuration      map[string]interface{} `json:"configuration,omitempty"`
	GroupMappings      []GroupMapping         `json:"groupMappings,omitempty"`
	RoleMappings       []RoleMapping          `json:"roleMappings,omitempty"`
	UserProfileMapping map[string]string      `json:"userProfileMapping,omitempty"`
	EmailRequired      bool                   `json:"emailRequired,omitempty"`
	SyncMappings       bool                   `json:"syncMappings,omitempty"`
	Organization       string                 `json:"organization,omitempty"`
}

type GroupMapping struct {
	Condition string   `json:"condition,omitempty"`
	Groups    []string `json:"groups"`
}

type RoleMapping struct {
	Condition     string              `json:"condition,omitempty"`
	Organizations []string            `json:"organizations"`
	Environments  map[string][]string `json:"environments"`
}
