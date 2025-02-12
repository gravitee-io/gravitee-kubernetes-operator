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
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/k8s"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	gAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type GraviteeGateway struct {
	// +kubebuilder:validation:MaxItems:=64
	Listeners  []GraviteeListener `json:"listeners"`
	Gravitee   *GraviteeConfig    `json:"gravitee"` // extended with listener extension
	Kubernetes *KubernetesConfig  `json:"kubernetes"`
}

type GraviteeConfig struct {
	// +kubebuilder:validation:Optional
	YAML   *utils.GenericStringMap `json:"yaml"`
	DBLess bool                    `json:"dbLess"`
}

type GraviteeListener struct {
	gAPIv1.Listener `json:",inline"`

	// +kubebuilder:validation:Optional
	Config *GraviteeListenerConfig `json:"config"`
}

type GraviteeListenerConfig struct {
	IdleTimeout  uint `json:"idleTimeout"`
	TCPKeepAlive bool `json:"tcpKeepAlive"`
}

type GraviteeListenerTLSConfig struct {
	Protocols  string `json:"tlsProtocols"`
	ClientAuth string `json:"clientAuth"`
}

type KubernetesConfig struct {
	Deployment k8s.Deployment `json:"deployment"`
	Service    k8s.Service    `json:"service"`
}
