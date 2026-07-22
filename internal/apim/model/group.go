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

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
)

// GroupState is the automation API representation of a group returned by GET.
type GroupState struct {
	GroupDTO    `json:",inline"`
	MemberCount uint          `json:"memberCount,omitempty"`
	Errors      status.Errors `json:"errors,omitempty"`
}

type GroupDTO struct {
	ID            string   `json:"id,omitempty" drift:"ignore"`
	HRID          string   `json:"hrid,omitempty" drift:"ignore"`
	Name          string   `json:"name"`
	NotifyMembers bool     `json:"notifyMembers" drift:"ignore"` // send empty returns true, so need to ignore
	Members       []Member `json:"members" drift:"empty-is-nil"`
}

type Member struct {
	Source   string                     `json:"source"`
	SourceID string                     `json:"sourceId"`
	Roles    map[group.RoleScope]string `json:"roles"`
}

func ToGroupDTO(grp group.Type) GroupDTO {
	return mapViaJSON[GroupDTO](grp)
}
