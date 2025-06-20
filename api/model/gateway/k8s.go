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
	appV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
)

type Deployment struct {
	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas *int32 `json:"replicas,omitempty"`
	// +kubebuilder:validation:Optional
	Strategy *appV1.DeploymentStrategy `json:"strategy,omitempty"`
	// The template.spec field uses the standard Kubernetes Pod template specification,
	// and its contents will be merged using a
	// [strategic merge patch](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/)
	// with Gravitee's default deployment configuration.
	// +kubebuilder:validation:Optional
	Template *coreV1.PodTemplateSpec `json:"template,omitempty"`
}

type Service struct {
	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=LoadBalancer
	Type *coreV1.ServiceType `json:"type,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=Cluster
	ExternalTrafficPolicy coreV1.ServiceExternalTrafficPolicy `json:"externalTrafficPolicy,omitempty"`
	// +kubebuilder:validation:Optional
	LoadBalancerClass *string `json:"loadBalancerClass,omitempty"`
}
