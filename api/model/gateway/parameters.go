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

package gateway

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// The GatewayClassParameters custom resource is
// the Gravitee.io extension point that allows you to configure
// our implementation of the [Kubernetes Gateway API](https://gateway-api.sigs.k8s.io/).
// It defines a set of configuration options to control how
// Gravitee Gateways are deployed and behave when managed via the Gateway API,
// including licensing, Kafka support, and Kubernetes-specific deployment settings.
type GatewayClassParameters struct {
	// The gravitee section controls Gravitee specific features
	// and allows you to configure and customize our implementation
	// of the Kubernetes Gateway API.
	// +kubebuilder:validation:Optional
	Gravitee *GraviteeConfig `json:"gravitee"`
	// The kubernetes section of the GatewayClassParameters
	// spec lets you customize core Kubernetes resources
	// that are part of your Gateway deployments.
	// +kubebuilder:validation:Optional
	Kubernetes *KubernetesConfig `json:"kubernetes"`
}

type GraviteeConfig struct {
	// +kubebuilder:validation:Optional
	// A reference to a Kubernetes Secret that contains your Gravitee license key.
	// This license is required to unlock advanced capabilities like Kafka protocol support.
	LicenseRef *gwAPIv1.SecretObjectReference `json:"licenseRef"`
	// Use this field to enable Kafka support in the Gateway.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default={enabled: false}
	Kafka *GraviteeKafkaConfig `json:"kafka"`
	// +kubebuilder:validation:Optional
	// Use this field to provide custom gateway configuration,
	// giving you control over additional configuration blocks
	// available in the gateway
	// [settings](https://documentation.gravitee.io/apim/configure-apim/apim-components/gravitee-gateway).
	YAML *utils.GenericStringMap `json:"yaml,omitempty"`
}

type GraviteeKafkaConfig struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=true
	Enabled bool `json:"enabled"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	RoutingHostMode *GraviteeKafkaRoutingHostModeConfig `json:"routingHostMode"`
}

type GraviteeKafkaRoutingHostModeConfig struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="broker-{brokerId}-{apiHost}"
	BokerDomainPattern string `json:"brokerDomainPattern"`
	// You can find details about these configurations options in our
	// [documentation](https://documentation.gravitee.io/apim/kafka-gateway/configure-the-kafka-gateway-and-client).
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="{apiHost}"
	BootstrapDomainPattern string `json:"bootstrapDomainPattern"`
}

type KubernetesConfig struct {
	// Use this field to modify pod labels and annotations,
	// adjust the number of replicas to control scaling,
	// specify update strategies for rolling updates,
	// and override the pod template to customize container specs,
	// security settings, or environment variables.
	// +kubebuilder:validation:Optional
	Deployment *Deployment `json:"deployment"`
	// Use this field to customize the Kubernetes Service that exposes the Gateway
	// by adding labels and annotations, choosing the service type,
	// configuring the external traffic policy, and specifying the load balancer class.`
	// +kubebuilder:validation:Optional
	Service *Service `json:"service"`
}
