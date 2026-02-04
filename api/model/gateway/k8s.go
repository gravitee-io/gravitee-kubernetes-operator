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
	autoscalingV2 "k8s.io/api/autoscaling/v2"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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

type Autoscaling struct {
	// Use this field to enable HorizontalPodAutoscaler reconciliation for the Gateway.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	Enabled bool `json:"enabled"`
	// The minimum number of replicas when autoscaling is enabled.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:default:=1
	MinReplicas *int32 `json:"minReplicas,omitempty"`
	// The maximum number of replicas when autoscaling is enabled.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:default:=10
	MaxReplicas *int32 `json:"maxReplicas,omitempty"`
	// The metrics used by the HorizontalPodAutoscaler to determine the desired replica count.
	// If empty and autoscaling is enabled, the operator will use a default CPU utilization metric
	// targeting 80% average utilization.
	// +kubebuilder:validation:Optional
	Metrics []autoscalingV2.MetricSpec `json:"metrics,omitempty"`
	// Behavior configures scaling behavior for the HorizontalPodAutoscaler.
	// +kubebuilder:validation:Optional
	Behavior *autoscalingV2.HorizontalPodAutoscalerBehavior `json:"behavior,omitempty"`
}

type PodDisruptionBudget struct {
	// Use this field to enable PodDisruptionBudget reconciliation for the Gateway.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	Enabled bool `json:"enabled"`
	// The minimum number of pods that must be available after an eviction.
	// +kubebuilder:validation:Optional
	MinAvailable *intstr.IntOrString `json:"minAvailable,omitempty"`
	// The maximum number of pods that can be unavailable after an eviction.
	// If neither minAvailable nor maxUnavailable is provided when enabled, the operator will
	// default maxUnavailable to 1.
	// +kubebuilder:validation:Optional
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`
}
