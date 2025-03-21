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

package kafka

import gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"

type KafkaRouteFilterType string

const (
	KafkaRouteFilterTypeAccessControlList = KafkaRouteFilterType("ACL")
	KafkaRouteFilterTypeExtensionRef      = KafkaRouteFilterType("ExtensionRef")
)

type KafkaRouteFilter struct {
	// +unionDiscriminator
	// +kubebuilder:validation:Enum=ACL;ExtensionRef
	Type KafkaRouteFilterType `json:"type"`
	// +optional
	ACL []KafkaAccessControl `json:"acl"`
	// +optional
	ExtensionRef *gwAPIv1.LocalObjectReference `json:"extensionRef,omitempty"`
}
