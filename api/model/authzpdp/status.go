// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authzpdp

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Status struct {
	// The ID of the PDP gateway in the Gravitee API Management instance
	// +kubebuilder:validation:Optional
	ID string `json:"id,omitempty"`
	// The organization ID defined in the management context
	// +kubebuilder:validation:Optional
	OrgID string `json:"organizationId,omitempty"`
	// The environment ID defined in the management context
	// +kubebuilder:validation:Optional
	EnvID string `json:"environmentId,omitempty"`
	// The human-readable ID the platform addresses this PDP gateway by, derived from the
	// resource namespace and name.
	// +kubebuilder:validation:Optional
	HRID string `json:"hrid,omitempty"`
	// targetPdpId@tag, or targetPdpId when untagged: what the pdpRef of a schema, an entity or a
	// policy resolves to. Empty until the platform has the PDP gateway.
	// +kubebuilder:validation:Optional
	RoutingScope string `json:"routingScope,omitempty"`
	// The AuthZEN discovery path of the endpoint. Absent without an endpoint, or when another API
	// already held that path when the endpoint was created.
	// +kubebuilder:validation:Optional
	WellKnownPath string `json:"wellKnownPath,omitempty"`
	// Conditions describe the current conditions of the AuthzPdp.
	//
	// Known condition types are:
	// * "Accepted"
	// * "ResolvedRefs"
	// * "AutomationAPIManaged"
	//
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:default={}
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// When the PDP gateway has been created regardless of errors, this field is
	// used to persist the error message encountered during admission
	Errors status.Errors `json:"errors,omitempty"`
}
