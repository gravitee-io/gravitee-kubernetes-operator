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
package k8s

import (
	appV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
)

type Deployment struct {
	Labels      map[string]string         `json:"labels,omitempty"`
	Annotations map[string]string         `json:"annotations,omitempty"`
	Replicas    *int32                    `json:"replicas,omitempty"`
	Strategy    *appV1.DeploymentStrategy `json:"strategy,omitempty"`
	Template    *coreV1.PodTemplateSpec   `json:"template,omitempty"`
}

type Service struct {
	Labels                map[string]string                   `json:"labels,omitempty"`
	Annotations           map[string]string                   `json:"annotations,omitempty"`
	Type                  *coreV1.ServiceType                 `json:"type,omitempty"`
	ExternalTrafficPolicy coreV1.ServiceExternalTrafficPolicy `json:"externalTrafficPolicy,omitempty"`
	LoadBalancerClass     *string                             `json:"loadBalancerClass,omitempty"`
}
