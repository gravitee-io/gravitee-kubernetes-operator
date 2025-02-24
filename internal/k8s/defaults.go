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

package k8s

import (
	appV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	DefaultCPURequest       = "200m"
	DefaultCPULimit         = "500m"
	DefaultMemRequest       = "256Mi"
	DefaultMemLimit         = "512Mi"
	DefaultConfigVolumeName = "config"
	DefaultConfigFileEntry  = "gravitee.yml"
	DefaultProbePort        = 18082

	GatewayContainerName     = "gateway"
	DefaultGatewayImage      = "graviteeio/apim-gateway"
	DefaultGatewayConfigFile = "/opt/graviteeio-gateway/config/gravitee.yml"

	InstanceLabelKey  = "app.kubernetes.io/instance"
	ComponentLabelKey = "app.kubernetes.io/component"
	VersionLabelKey   = "app.kubernetes.io/version"
	NameLabelKey      = "app.kubernetes.io/name"
	PartOfLabelKey    = "app.kubernetes.io/part-of"
	ManagedByLabelKey = "app.kubernetes.io/managed-by"

	ManagedByLabelValue        = "gko.gravitee.io"
	PartOfLabelValue           = "apim.gravitee.io"
	GatewayComponentLabelValue = "gateway"
)

var DefaultReplicas int32 = 1

var DefaultLivenessProbe = &coreV1.Probe{
	FailureThreshold: 3,
	PeriodSeconds:    5,
	SuccessThreshold: 1,
	TimeoutSeconds:   2,
	ProbeHandler: coreV1.ProbeHandler{
		HTTPGet: &coreV1.HTTPGetAction{
			HTTPHeaders: []coreV1.HTTPHeader{
				{Name: "Authorization", Value: "Basic YWRtaW46YWRtaW4="},
			},
			Path:   "/_node/health?probes=http-server",
			Port:   intstr.FromInt32(DefaultProbePort),
			Scheme: "HTTP",
		},
	},
}

var DefaultReadinessProbe = &coreV1.Probe{
	FailureThreshold: 2,
	PeriodSeconds:    10,
	SuccessThreshold: 1,
	TimeoutSeconds:   2,
	ProbeHandler: coreV1.ProbeHandler{
		HTTPGet: &coreV1.HTTPGetAction{
			HTTPHeaders: []coreV1.HTTPHeader{
				{Name: "Authorization", Value: "Basic YWRtaW46YWRtaW4="},
			},
			Path:   "/_node/health?probes=http-server",
			Port:   intstr.FromInt32(DefaultProbePort),
			Scheme: "HTTP",
		},
	},
}

var DefaultStartupProbe = &coreV1.Probe{
	FailureThreshold:    100,
	InitialDelaySeconds: 5,
	PeriodSeconds:       2,
	SuccessThreshold:    1,
	TimeoutSeconds:      2,
	ProbeHandler: coreV1.ProbeHandler{
		HTTPGet: &coreV1.HTTPGetAction{
			HTTPHeaders: []coreV1.HTTPHeader{
				{Name: "Authorization", Value: "Basic YWRtaW46YWRtaW4="},
			},
			Path:   "/_node/health?probes=http-server",
			Port:   intstr.FromInt32(DefaultProbePort),
			Scheme: "HTTP",
		},
	},
}

var DefaultResources = &coreV1.ResourceRequirements{
	Requests: coreV1.ResourceList{
		coreV1.ResourceCPU:    resource.MustParse(DefaultCPURequest),
		coreV1.ResourceMemory: resource.MustParse(DefaultMemRequest),
	},
	Limits: coreV1.ResourceList{
		coreV1.ResourceCPU:    resource.MustParse(DefaultCPULimit),
		coreV1.ResourceMemory: resource.MustParse(DefaultMemLimit),
	},
}

var DefaultGatewayVolumeMounts = []coreV1.VolumeMount{
	{
		Name:      DefaultConfigVolumeName,
		MountPath: DefaultGatewayConfigFile,
		SubPath:   DefaultConfigFileEntry,
	},
}

var DefaultGatewayContainer = coreV1.Container{
	Image:          DefaultGatewayImage,
	Name:           GatewayContainerName,
	LivenessProbe:  DefaultLivenessProbe,
	ReadinessProbe: DefaultReadinessProbe,
	StartupProbe:   DefaultStartupProbe,
	VolumeMounts:   DefaultGatewayVolumeMounts,
	Ports:          []coreV1.ContainerPort{},
}

var DefaultGatewayVolume = coreV1.Volume{
	Name:         DefaultConfigVolumeName,
	VolumeSource: DefaultVolumeSource,
}

var DefaultVolumeSource = coreV1.VolumeSource{
	ConfigMap: &coreV1.ConfigMapVolumeSource{
		LocalObjectReference: coreV1.LocalObjectReference{},
	},
}

var DefaultGatewayConfigMap = &coreV1.ConfigMap{
	ObjectMeta: metaV1.ObjectMeta{},
	Data:       map[string]string{},
}

var DefaultGatewayPodSpec = &coreV1.PodSpec{
	Containers: []coreV1.Container{
		DefaultGatewayContainer,
	},
	Resources: DefaultResources,
	Volumes:   []coreV1.Volume{},
}

var DefaultGatewayPodTemplate = &coreV1.PodTemplateSpec{
	ObjectMeta: metaV1.ObjectMeta{},
	Spec:       *DefaultGatewayPodSpec,
}

var DefaultGatewayDeploymentSpec = &appV1.DeploymentSpec{
	Replicas: &DefaultReplicas,
	Selector: &metaV1.LabelSelector{},
	Template: coreV1.PodTemplateSpec{
		Spec: *DefaultGatewayPodSpec,
	},
}

var DefaultGatewayDeployment = &appV1.Deployment{
	ObjectMeta: metaV1.ObjectMeta{},
	Spec:       *DefaultGatewayDeploymentSpec,
}

var DefaultServiceSpec = &coreV1.ServiceSpec{
	Ports: []coreV1.ServicePort{},
	Type:  coreV1.ServiceTypeLoadBalancer,
}

var DefaultService = &coreV1.Service{
	ObjectMeta: metaV1.ObjectMeta{},
	Spec:       *DefaultServiceSpec,
}

func DefaultLabels(appName string) map[string]string {
	return map[string]string{
		ComponentLabelKey: GatewayComponentLabelValue,
		InstanceLabelKey:  GatewayComponentLabelValue,
		PartOfLabelKey:    PartOfLabelValue,
		ManagedByLabelKey: ManagedByLabelValue,
		NameLabelKey:      appName,
	}
}
