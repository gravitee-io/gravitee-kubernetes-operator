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
	"fmt"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/group"
)

type Member struct {
	// Member source
	// +kubebuilder:validation:Required
	// +kubebuilder:example:=gravitee
	Source string `json:"source"`
	// Member source ID
	// +kubebuilder:validation:Required
	// +kubebuilder:example:=user@email.com
	SourceID string `json:"sourceId"`
	// The API role associated with this Member
	// +kubebuilder:default:=USER
	Role string `json:"role,omitempty"`
}

func (m *Member) String() string {
	return fmt.Sprintf(
		"{Source:%s,SourceID:%s,Role:%s}",
		m.Source, m.SourceID, m.Role,
	)
}

func NewGraviteeMember(username, role string) *Member {
	return &Member{
		Source:   "gravitee",
		SourceID: username,
		Role:     role,
	}
}

func NewMemoryMember(username, role string) *Member {
	return &Member{
		Source:   "memory",
		SourceID: username,
		Role:     role,
	}
}

func NewGraviteeGroupMember(username, role string) group.Member {
	return group.Member{
		Source:   "gravitee",
		SourceID: username,
		Roles: map[group.RoleScope]string{
			group.APIRoleScope:         role,
			group.ApplicationRoleScope: role,
			group.IntegrationRoleScope: role,
		},
	}
}
