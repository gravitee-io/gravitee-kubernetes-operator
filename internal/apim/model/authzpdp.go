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

package model

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/status"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
)

// AuthzPdpDTO is the automation API wire representation of an AuthzPdp (PUT /authz/pdps).
type AuthzPdpDTO struct {
	HRID        string `json:"hrid,omitempty" drift:"ignore"`
	Name        string `json:"name"`
	TargetPdpID string `json:"targetPdpId"`
	Tag         string `json:"tag,omitempty"`
	// Strict is always sent, so the wire and drift agree on false.
	Strict   bool                 `json:"strict"`
	Endpoint *AuthzPdpEndpointDTO `json:"endpoint,omitempty" drift:"empty-is-nil"`
}

// AuthzPdpEndpointDTO is the AuthZEN endpoint of the PDP gateway.
type AuthzPdpEndpointDTO struct {
	Path string `json:"path"`
}

// AuthzPdpState is what a PUT or a GET answers: the payload, the envelope and what the platform
// derives. It deliberately does not embed authzpdp.Status: both halves would promote a hrid field
// and encoding/json silently drops a name promoted twice at the same depth.
type AuthzPdpState struct {
	AuthzPdpDTO   `json:",inline"`
	ID            string        `json:"id,omitempty"`
	OrgID         string        `json:"organizationId,omitempty"`
	EnvID         string        `json:"environmentId,omitempty"`
	Errors        status.Errors `json:"errors,omitempty"`
	RoutingScope  string        `json:"routingScope,omitempty"`
	WellKnownPath string        `json:"wellKnownPath,omitempty"`
}

// ToAuthzPdpDTO maps the CRD onto the wire payload. It is the single mapping used by the sync
// path and by drift detection.
func ToAuthzPdpDTO(crd *v1alpha1.AuthzPdp) AuthzPdpDTO {
	spec := crd.Spec.Type

	dto := AuthzPdpDTO{
		HRID:        refs.NewNamespacedNameFromObject(crd).HRID(),
		Name:        spec.Name,
		TargetPdpID: spec.TargetPdpID,
		Tag:         spec.Tag,
		Strict:      spec.Strict,
	}
	if spec.Endpoint != nil {
		dto.Endpoint = &AuthzPdpEndpointDTO{Path: spec.Endpoint.Path}
	}

	return dto
}
