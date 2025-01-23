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

import "time"

type SharedPolicyGroup struct {
	ID                  string    `json:"id,omitempty"`
	CrossID             string    `json:"crossId,omitempty"`
	Name                string    `json:"name,omitempty"`
	Status              string    `json:"status,omitempty"`
	Description         string    `json:"description,omitempty"`
	PrerequisiteMessage string    `json:"prerequisiteMessage,omitempty"`
	Version             int       `json:"version,omitempty"`
	AppType             string    `json:"type,omitempty"`
	Steps               []any     `json:"steps,omitempty"` // it is not really important for us
	Phase               string    `json:"phase,omitempty"`
	DeployedAt          time.Time `json:"deployedAt,omitempty"`
	CreatedAt           time.Time `json:"createdAt,omitempty"`
	UpdatedAt           time.Time `json:"updatedAt,omitempty"`
	LifecycleState      string    `json:"lifecycleState,omitempty"`
}
