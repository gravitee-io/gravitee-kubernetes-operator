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
	// Hostname is used to uniquely route clients to this API.
	// Your client must trust the certificate provided by the gateway,
	// and as there is a variable host in the proxy bootstrap server URL,
	// you likely need to request a wildcard SAN for the certificate presented by the gateway.
	// If empty, the hostname defined in the Kafka listener of the parent will be used.
	// +optional
	Hostname *gwAPIv1.Hostname `json:"hostname,omitempty"`
	// BackendRefs defines the backend(s) where matching requests should be sent.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	BackendRefs []KafkaBackendRef `json:"backendRefs,omitempty"`
	// Filters define the filters that are applied to Kafka trafic matching this route.
	// +kubebuilder:validation:MaxItems=16
	// +optional
	Filters []KafkaRouteFilter `json:"filters,omitempty"`
	// +optional
	// +kubebuilder:validation:MaxProperties=16
	// Options are a list of key/value pairs to enable extended configuration specific
	// to an
	Options map[gwAPIv1.AnnotationKey]gwAPIv1.AnnotationValue `json:"options,omitempty"`
}

// This currently wraps the code gateway API BackendObjectReference type,
// leaving room for e.g. backend security configuration.
type KafkaBackendRef struct {
	// BackendObjectReference defines how an ObjectReference that is
	// specific to BackendRef. It includes a few additional fields and features
	// than a regular ObjectReference.
	gwAPIv1.BackendObjectReference `json:",inline"`
}
