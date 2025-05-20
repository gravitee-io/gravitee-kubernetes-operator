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

package group

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
)

// +kubebuilder:validation:Enum=API;APPLICATION;INTEGRATION;
type RoleScope string

const (
	APIRoleScope         = RoleScope("API")
	ApplicationRoleScope = RoleScope("APPLICATION")
	IntegrationRoleScope = RoleScope("INTEGRATION")
)

type Type struct {
	// +kubebuilder:validation:Optional
	ID string `json:"id,omitempty"`
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	// If true, new members added to the API spec will
	// be notified when the API is synced with APIM.
	NotifyMembers bool     `json:"notifyMembers"`
	Members       []Member `json:"members"`
}

type Member struct {
	// Member source
	// +kubebuilder:validation:Required
	// +kubebuilder:example:=gravitee
	Source string `json:"source"`
	// Member source ID
	// +kubebuilder:validation:Required
	// +kubebuilder:example:=user@email.com
	SourceID string `json:"sourceId"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:={}
	Roles map[RoleScope]string `json:"roles"`
}

type Status struct {
	// The ID of the Group in the Gravitee API Management instance
	// +kubebuilder:validation:Optional
	ID string `json:"id,omitempty"`
	// The organization ID defined in the management context
	// +kubebuilder:validation:Optional
	OrgID string `json:"organizationId,omitempty"`
	// The environment ID defined in the management context
	// +kubebuilder:validation:Optional
	EnvID string `json:"environmentId,omitempty"`
	// The processing status of the Group.
	ProcessingStatus core.ProcessingStatus `json:"processingStatus,omitempty"`
	// The number of members added to this group
	Members uint `json:"members"`
	// When group has been created regardless of errors, this field is
	// used to persist the error message encountered during admission
	Errors status.Errors `json:"errors,omitempty"`
}
