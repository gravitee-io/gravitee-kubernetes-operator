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

package idpgroupmapping

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	Groups []string `json:"groups"`
	// +kubebuilder:validation:Required
	IDPID string `json:"idpId"`
	// +kubebuilder:validation:Required
	Condition string `json:"condition"`
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
	// The group name in the Gravitee API Management instance
	// +kubebuilder:validation:Optional
	Groups []string `json:"group,omitempty"`
	// Conditions describe the current conditions of the Group.
	//
	// Known condition types are:
	// * "Accepted"
	// * "ResolvedRefs"
	//
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:default={}
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// The processing status of the Group. *** DEPRECATED ***
	ProcessingStatus core.ProcessingStatus `json:"processingStatus,omitempty"`
	// When group has been created regardless of errors, this field is
	// used to persist the error message encountered during admission
	Errors status.Errors `json:"errors,omitempty"`
}
