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

package v4

import "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"

type Service struct {
	// Is the service enabled or not ?
	Enabled bool `json:"enabled"`

	// Service Type
	// +kubebuilder:validation:Optional
	Type *string `json:"type,omitempty"`

	// Service Override Configuration or not?
	OverrideConfig bool `json:"overrideConfiguration"`

	// Service Configuration, a map of arbitrary key-values
	// +kubebuilder:validation:Optional
	Config *utils.GenericStringMap `json:"configuration,omitempty"`
}

func NewService(kind string, enabled bool) *Service {
	return &Service{
		Enabled: enabled,
		Type:    &kind,
	}
}

type EndpointServices struct {
	// Health check service
	HealthCheck *Service `json:"healthCheck,omitempty"`
}

type EndpointGroupServices struct {
	// Endpoint group discovery service
	Discovery *Service `json:"discovery,omitempty"`

	// Endpoint group health check service
	HealthCheck *Service `json:"healthCheck,omitempty"`
}

type ApiServices struct {
	// API dynamic property service
	DynamicProperty *Service `json:"dynamicProperty,omitempty"`
}
