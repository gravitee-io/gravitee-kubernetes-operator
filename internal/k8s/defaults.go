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
	"maps"

	appV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	DefaultCPURequest        = "200m"
	DefaultCPULimit          = "500m"
	DefaultMemRequest        = "256Mi"
	DefaultMemLimit          = "512Mi"
	DefaultConfigVolumeName  = "config"
	DefaultLicenseVolumeName = "license"
	DefaultConfigFileEntry   = "gravitee.yml"
	DefaultProbePort         = 18082

	GatewayConfigMapPrefix     = "gio-gw-config-"
	PEMRegistryConfigMapPrefix = "gio-pem-registry-"

	GatewayContainerName     = "gateway"
	DefaultGatewayImage      = "graviteeio/apim-gateway"
	DefaultGatewayConfigFile = "/opt/graviteeio-gateway/config/gravitee.yml"

	DefaultLicenseMountPath = "/opt/graviteeio-gateway/license"

	InstanceLabelKey  = "app.kubernetes.io/instance"
	ComponentLabelKey = "app.kubernetes.io/component"
	VersionLabelKey   = "app.kubernetes.io/version"
	NameLabelKey      = "app.kubernetes.io/name"
	PartOfLabelKey    = "app.kubernetes.io/part-of"
	ManagedByLabelKey = "app.kubernetes.io/managed-by"

	ManagedByLabelValue            = "gko.gravitee.io"
	PartOfLabelValue               = "apim.gravitee.io"
	GatewayComponentLabelValue     = "gateway"
	HTTPRouteComponentlabelValue   = "http-route"
	PEMRegistryComponentLabelValue = "kubernetes-pem-registry"
)

var (
	DefaultReplicas         int32 = 1
	DefaultVolumeSourceMode int32 = 420
)

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

var ConfigVolumeMount = coreV1.VolumeMount{
	Name:      DefaultConfigVolumeName,
	MountPath: DefaultGatewayConfigFile,
	SubPath:   DefaultConfigFileEntry,
	ReadOnly:  true,
}

var LicenseVolumeMount = coreV1.VolumeMount{
	Name:      DefaultLicenseVolumeName,
	MountPath: DefaultLicenseMountPath,
	ReadOnly:  true,
}

var DefaultGatewayContainer = coreV1.Container{
	Image:           DefaultGatewayImage,
	ImagePullPolicy: coreV1.PullIfNotPresent,
	Name:            GatewayContainerName,
	LivenessProbe:   DefaultLivenessProbe,
	ReadinessProbe:  DefaultReadinessProbe,
	StartupProbe:    DefaultStartupProbe,
	VolumeMounts:    []coreV1.VolumeMount{},
	Ports:           []coreV1.ContainerPort{},
}

var DefaultGatewayConfigVolume = coreV1.Volume{
	Name:         DefaultConfigVolumeName,
	VolumeSource: DefaultConfigVolumeSource,
}

var DefaultLicenseConfigVolume = coreV1.Volume{
	Name:         DefaultLicenseVolumeName,
	VolumeSource: DefaultLicenseVolumeSource,
}

var DefaultLicenseVolumeSource = coreV1.VolumeSource{
	Secret: &coreV1.SecretVolumeSource{
		DefaultMode: &DefaultVolumeSourceMode,
	},
}

var DefaultConfigVolumeSource = coreV1.VolumeSource{
	ConfigMap: &coreV1.ConfigMapVolumeSource{
		DefaultMode:          &DefaultVolumeSourceMode,
		LocalObjectReference: coreV1.LocalObjectReference{},
	},
}

var DefaultGatewayConfigMap = &coreV1.ConfigMap{
	ObjectMeta: metaV1.ObjectMeta{},
	Data:       map[string]string{},
}

var DefaultPEMRegistryConfigMap = &coreV1.ConfigMap{
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

var CommonLabels = map[string]string{
	PartOfLabelKey:    PartOfLabelValue,
	ManagedByLabelKey: ManagedByLabelValue,
}

func GIOPemRegistryLabels(gwName string) map[string]string {
	labels := map[string]string{
		ComponentLabelKey: PEMRegistryComponentLabelValue,
		InstanceLabelKey:  gwName,
		NameLabelKey:      gwName,
	}
	maps.Copy(labels, CommonLabels)
	return labels
}

func GwAPIv1GatewayLabels(gwName string) map[string]string {
	labels := map[string]string{
		ComponentLabelKey: GatewayComponentLabelValue,
		InstanceLabelKey:  gwName,
		NameLabelKey:      gwName,
	}
	maps.Copy(labels, CommonLabels)
	return labels
}

func GwAPIv1HTTPRouteLabels(routeNAme string) map[string]string {
	labels := map[string]string{
		ComponentLabelKey: HTTPRouteComponentlabelValue,
		InstanceLabelKey:  routeNAme,
		NameLabelKey:      routeNAme,
	}
	maps.Copy(labels, CommonLabels)
	return labels
}
