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

// +kubebuilder:object:generate=true
package gateway

import (
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type GatewayClassParameters struct {
	// +kubebuilder:validation:Optional
	Gravitee *GraviteeConfig `json:"gravitee"`
	// +kubebuilder:validation:Optional
	Kubernetes *KubernetesConfig `json:"kubernetes"`
}

type GraviteeConfig struct {
	// +kubebuilder:validation:Optional
	LicenseRef *gwAPIv1.SecretObjectReference `json:"licenseRef"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default={enabled: false}
	Kafka *GraviteeKafkaConfig `json:"kafka"`
}

type GraviteeKafkaConfig struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	Enabled bool `json:"enabled"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default={brokerPrefix: broker-, domainSeparator: -}
	RoutingHostMode *GraviteeKafkaRoutingHostModeConfig `json:"routingHostMode"`
}

type GraviteeKafkaRoutingHostModeConfig struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=broker-
	BrokerPrefix string `json:"brokerPrefix"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=-
	DomainSeparator string `json:"domainSeparator"`
}

type GraviteeKafkaListenerConfig struct {
	Name gwAPIv1.SectionName `json:"name"`
}

type KubernetesConfig struct {
	// +kubebuilder:validation:Optional
	Deployment *Deployment `json:"deployment"`
	// +kubebuilder:validation:Optional
	Service *Service `json:"service"`
}
