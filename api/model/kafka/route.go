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

// +kubebuilder:object:generate=true
package kafka

import (
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type KafKaRoute struct {
	gwAPIv1.CommonRouteSpec `json:",inline"`
	// +optional
	// We do not accept the route in our case if hostname is not set (for now)
	Hostname *gwAPIv1.Hostname `json:"hostname,omitempty"`
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	BackendRefs []KafkaBackendRef `json:"backendRefs,omitempty"`
	// +kubebuilder:validation:MaxItems=16
	Filters []KafkaRouteFilter `json:"filters,omitempty"`
	// +optional
	// +kubebuilder:validation:MaxProperties=16
	Options map[gwAPIv1.AnnotationKey]gwAPIv1.AnnotationValue `json:"options,omitempty"`
}

// Leave room for backend security configuration.
type KafkaBackendRef struct {
	gwAPIv1.BackendObjectReference `json:",inline"`
}
