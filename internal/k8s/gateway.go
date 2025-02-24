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
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"slices"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/gateway/yaml"
	appV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/sets"
	patchUtil "k8s.io/apimachinery/pkg/util/strategicpatch"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	GwAPIv1HTTPRouteKind = "HTTPRoute"
)

var GwAPIv1Group = gwAPIv1.Group(gwAPIv1.GroupVersion.Group)

var SupportedGwAPIProtocols = sets.New(
	gwAPIv1.HTTPProtocolType,
	gwAPIv1.HTTPSProtocolType,
)

var ProtocolToRouteKinds = map[gwAPIv1.ProtocolType][]gwAPIv1.RouteGroupKind{
	gwAPIv1.HTTPProtocolType: {
		{
			Group: &GwAPIv1Group,
			Kind:  GwAPIv1HTTPRouteKind,
		},
	},
	gwAPIv1.HTTPSProtocolType: {
		{
			Group: &GwAPIv1Group,
			Kind:  GwAPIv1HTTPRouteKind,
		},
	},
}

var ProtocolToServerType = map[gwAPIv1.ProtocolType]string{
	gwAPIv1.HTTPProtocolType:  "http",
	gwAPIv1.HTTPSProtocolType: "http",
}

func GetSupportedRouteKinds(listener gwAPIv1.Listener) []gwAPIv1.RouteGroupKind {
	if kinds, ok := ProtocolToRouteKinds[listener.Protocol]; ok {
		return kinds
	}
	return []gwAPIv1.RouteGroupKind{}
}

func IsListenerRef(
	gw *gwAPIv1.Gateway,
	listener gwAPIv1.Listener,
	ref gwAPIv1.ParentReference,
) bool {
	if !IsGatewayRef(gw, ref) {
		return false
	}
	if ref.SectionName == nil {
		return true
	}
	if *ref.SectionName != listener.Name {
		return false
	}
	if ref.Port == nil {
		return true
	}
	return *ref.Port == listener.Port
}

func IsGatewayRef(gw *gwAPIv1.Gateway, ref gwAPIv1.ParentReference) bool {
	if ref.Group != nil && *ref.Group != GwAPIv1Group {
		return false
	}
	if ref.Kind != nil && string(*ref.Kind) != gw.Kind {
		return false
	}
	if string(ref.Name) != gw.Name {
		return false
	}
	return true
}

func HasHTTPSupport(listener gwAPIv1.Listener) bool {
	if kinds, ok := ProtocolToRouteKinds[listener.Protocol]; !ok {
		return false
	} else {
		for _, k := range kinds {
			if k.Kind == GwAPIv1HTTPRouteKind {
				return true
			}
		}
	}
	return false
}

func DeployGateway(
	ctx context.Context,
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
) error {
	if configMap, err := getConfigMap(gw); err != nil {
		return err
	} else if err = CreateOrUpdate(ctx, configMap); err != nil {
		return err
	}

	if deployment, err := getDeployment(gw, params); err != nil {
		return err
	} else if err := CreateOrUpdate(ctx, deployment); err != nil {
		return err
	}

	service := getService(gw, params)
	return CreateOrUpdate(ctx, service)
}

func getConfigMap(gw *gateway.Gateway) (*coreV1.ConfigMap, error) {
	labels := DefaultLabels(gw.Object.Name)
	configMap := DefaultGatewayConfigMap.DeepCopy()

	servers, err := getYAMLServers(gw.Object.Spec.Listeners)
	if err != nil {
		return nil, err
	}

	yaml := yaml.DBLess.DeepCopy()
	yaml.Put("servers", servers)

	if httpPort := getHTTPPort(gw.Object.Spec.Listeners); httpPort != nil {
		yaml.Put("http", map[string]int32{"port": *httpPort})
	}

	yamlData, err := yaml.MarshalYAML()
	if err != nil {
		return nil, err
	}

	configMap.Name = gw.Object.Name
	configMap.Namespace = gw.Object.Namespace
	configMap.Data[DefaultConfigFileEntry] = string(yamlData)
	configMap.Labels = labels

	setOwnerReference(gw.Object, configMap)

	return configMap, nil
}

func getDeployment(gw *gateway.Gateway, params *v1alpha1.GatewayClassParameters) (*appV1.Deployment, error) {
	labels := DefaultLabels(gw.Object.Name)
	deployment := DefaultGatewayDeployment.DeepCopy()

	deployment.Name = gw.Object.Name
	deployment.Namespace = gw.Object.Namespace
	deployment.Labels = labels
	deployment.Spec.Template.Labels = labels
	deployment.Spec.Selector.MatchLabels = labels

	template, err := getPodTemplateSpec(params)
	if err != nil {
		return nil, err
	}

	volume := DefaultGatewayVolume.DeepCopy()
	volume.ConfigMap.LocalObjectReference.Name = gw.Object.Name

	template.Labels = labels
	template.Spec.Volumes = []coreV1.Volume{*volume}

	prepareContainer(template, gw.Object)

	deployment.Spec.Template = *template

	setOwnerReference(gw.Object, deployment)

	return deployment, nil
}

func getService(gw *gateway.Gateway, params *v1alpha1.GatewayClassParameters) *coreV1.Service {
	labels := DefaultLabels(gw.Object.Name)
	svc := DefaultService.DeepCopy()

	svc.Name = gw.Object.Name
	svc.Namespace = gw.Object.Namespace
	svc.Labels = labels
	svc.Spec.Selector = labels

	if params.Spec.Kubernetes == nil {
		return svc
	}
	if params.Spec.Kubernetes.Service == nil {
		return svc
	}

	svc.Spec.Type = *params.Spec.Kubernetes.Service.Type

	setServicePorts(svc, gw.Object.Spec.Listeners)

	setOwnerReference(gw.Object, svc)

	return svc
}

func getPodTemplateSpec(parameters *v1alpha1.GatewayClassParameters) (*coreV1.PodTemplateSpec, error) {
	base := DefaultGatewayPodTemplate.DeepCopy()
	if parameters.Spec.Kubernetes == nil {
		return base, nil
	}
	if parameters.Spec.Kubernetes.Deployment == nil {
		return base, nil
	}
	if parameters.Spec.Kubernetes.Deployment.Template == nil {
		return base, nil
	}
	patch := parameters.Spec.Kubernetes.Deployment.Template
	if template, err := mergePodTemplates(base, patch); err != nil {
		return nil, err
	} else {
		return template, nil
	}
}

func prepareContainer(template *coreV1.PodTemplateSpec, gw *gwAPIv1.Gateway) {
	container := getGatewayContainer(template.Spec.Containers)
	container.VolumeMounts = DefaultGatewayVolumeMounts
	setContainerPorts(container, gw.Spec.Listeners)
	for i := range template.Spec.Containers {
		if template.Spec.Containers[i].Name == GatewayContainerName {
			template.Spec.Containers[i] = *container
		}
	}
}

func getGatewayContainer(containers []coreV1.Container) *coreV1.Container {
	for _, container := range containers {
		if container.Name == GatewayContainerName {
			return &container
		}
	}
	return nil
}

func setContainerPorts(container *coreV1.Container, listeners []gwAPIv1.Listener) {
	ports := make([]coreV1.ContainerPort, 0)
	for _, listener := range listeners {
		port := coreV1.ContainerPort{
			ContainerPort: int32(listener.Port),
			Protocol:      coreV1.ProtocolTCP,
		}
		ports = append(ports, port)
	}
	probePorts := make(map[int32]coreV1.ContainerPort)
	if rp := getProbePort(container.ReadinessProbe); rp != nil {
		probePorts[*rp] = coreV1.ContainerPort{
			ContainerPort: *rp,
			Protocol:      coreV1.ProtocolTCP,
		}
	}
	if lp := getProbePort(container.LivenessProbe); lp != nil {
		probePorts[*lp] = coreV1.ContainerPort{
			ContainerPort: *lp,
			Protocol:      coreV1.ProtocolTCP,
		}
	}
	if sp := getProbePort(container.LivenessProbe); sp != nil {
		probePorts[*sp] = coreV1.ContainerPort{
			ContainerPort: *sp,
			Protocol:      coreV1.ProtocolTCP,
		}
	}
	ports = append(ports, slices.Collect(maps.Values(probePorts))...)
	container.Ports = ports
}

func getProbePort(probe *coreV1.Probe) *int32 {
	if probe == nil {
		return nil
	}
	if probe.TCPSocket != nil {
		port := probe.TCPSocket.Port
		return &port.IntVal
	}
	if probe.HTTPGet != nil {
		port := probe.HTTPGet.Port
		return &port.IntVal
	}
	return nil
}

func setServicePorts(svc *coreV1.Service, listeners []gwAPIv1.Listener) {
	servicePorts := make([]coreV1.ServicePort, 0)
	knownPortNumbers := sets.New[int32]()
	for _, listener := range listeners {
		portNumber := int32(listener.Port)
		if knownPortNumbers.Has(portNumber) {
			continue
		}
		knownPortNumbers.Insert(portNumber)
		servicePort := coreV1.ServicePort{
			Name:       string(listener.Name),
			Port:       portNumber,
			TargetPort: intstr.FromInt32(int32(listener.Port)),
			Protocol:   coreV1.ProtocolTCP,
		}
		if svc.Spec.Type == "NodePort" && portNumber >= 30000 {
			servicePort.NodePort = portNumber
		}
		servicePorts = append(servicePorts, servicePort)
	}
	svc.Spec.Ports = servicePorts
}

func setOwnerReference(gw *gwAPIv1.Gateway, obj client.Object) {
	kind := gw.GetObjectKind().GroupVersionKind().Kind
	version := gw.GetObjectKind().GroupVersionKind().GroupVersion().String()
	obj.SetOwnerReferences(
		[]metaV1.OwnerReference{
			{
				Kind:       kind,
				APIVersion: version,
				Name:       gw.GetName(),
				UID:        gw.GetUID(),
			},
		},
	)
}

func getHTTPPort(listeners []gwAPIv1.Listener) *int32 {
	port := int32(0)
	for _, listener := range listeners {
		if listener.Protocol == gwAPIv1.HTTPProtocolType {
			port = int32(listener.Port)
		}
		if port == 0 && listener.Protocol == gwAPIv1.HTTPSProtocolType {
			port = int32(listener.Port)
		}
	}
	if port > 0 {
		return &port
	}
	return nil
}

func getYAMLServers(listeners []gwAPIv1.Listener) ([]map[string]interface{}, error) {
	servers := make([]map[string]interface{}, len(listeners))
	for i, listener := range listeners {
		server := make(map[string]interface{})
		serverType, ok := ProtocolToServerType[listener.Protocol]
		if !ok {
			return nil, fmt.Errorf("unknown server protocol")
		}
		server["type"] = serverType
		server["port"] = listener.Port
		server["hostname"] = "0.0.0.0"
		if listener.Hostname != nil {
			server["hostname"] = listener.Hostname
		}
		servers[i] = server
	}
	return servers, nil
}

func mergePodTemplates(base, patch *coreV1.PodTemplateSpec) (*coreV1.PodTemplateSpec, error) {
	if patch == nil {
		return base, nil
	}

	baseBytes, err := json.Marshal(base)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON for base %s: %w", base.Name, err)
	}

	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON for patch %s: %w", patch.Name, err)
	}

	jsonResultBytes, err := patchUtil.StrategicMergePatch(baseBytes, patchBytes, &coreV1.PodTemplateSpec{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate merge patch for %s: %w", base.Name, err)
	}

	patchResult := base.DeepCopy()
	if err := json.Unmarshal(jsonResultBytes, patchResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal merged %s: %w", base.Name, err)
	}

	return patchResult, nil
}
