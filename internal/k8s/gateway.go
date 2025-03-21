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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/hash"
	appV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/sets"
	patchUtil "k8s.io/apimachinery/pkg/util/strategicpatch"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func DeployGateway(
	ctx context.Context,
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
) error {
	if err := CreateOrUpdate(ctx, getServiceAccount(gw)); err != nil {
		return err
	}
	if err := CreateOrUpdate(ctx, getRole(gw)); err != nil {
		return err
	}
	if err := CreateOrUpdate(ctx, getRoleBinding(gw)); err != nil {
		return err
	}

	if pemRegistry, err := getPEMRegistryConfigMap(gw); err != nil {
		return err
	} else if err := CreateOrUpdate(ctx, pemRegistry, func() error {
		return buildPEMRegistryData(gw, pemRegistry)
	}); err != nil {
		return err
	}

	gatewayConfig, err := getGatewayConfigMap(gw, params)
	if err != nil {
		return err
	} else if err := CreateOrUpdate(ctx, gatewayConfig, func() error {
		return buildGatewayConfigData(gw, params, gatewayConfig)
	}); err != nil {
		return err
	}

	if deployment, err := getDeployment(gw, params); err != nil {
		return err
	} else if err := CreateOrUpdate(ctx, deployment, func() error {
		AddAnnotation(deployment, "gravitee.io/config", hash.Calculate(gatewayConfig))
		return buildDeployment(gw, params, deployment)
	}); err != nil {
		return err
	}

	svc := getService(gw, params)
	return CreateOrUpdate(ctx, svc, func() error {
		buildServiceSpec(gw, params, svc)
		return nil
	})
}

func getPEMRegistryConfigMap(gw *gateway.Gateway) (*coreV1.ConfigMap, error) {
	configMap := DefaultPEMRegistryConfigMap.DeepCopy()
	configMap.Name = getPEMRegistryConfigMapName(gw.Object.Name)
	configMap.Namespace = gw.Object.Namespace
	configMap.Labels = GIOPemRegistryLabels(gw.Object.Name)
	if err := buildPEMRegistryData(gw, configMap); err != nil {
		return nil, err
	}
	return configMap, nil
}

func buildPEMRegistryData(gw *gateway.Gateway, configMap *coreV1.ConfigMap) error {
	data := make(map[string]string)
	for i := range gw.Object.Spec.Listeners {
		l := gw.Object.Spec.Listeners[i]
		if l.TLS != nil {
			ns := gw.Object.Namespace
			secretName := string(l.TLS.CertificateRefs[0].Name)
			v, err := buildPEMRegistryValue(ns, secretName)
			if err != nil {
				return err
			}
			k := string(l.Name)
			data[k] = v
		}
	}
	configMap.Data = data
	return nil
}

func buildPEMRegistryValue(ns, refName string) (string, error) {
	registry := []string{ns + "/" + refName}
	b, err := json.Marshal(registry)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func getGatewayConfigMap(
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
) (*coreV1.ConfigMap, error) {
	labels := GwAPIv1GatewayLabels(gw.Object.Name)
	configMap := DefaultGatewayConfigMap.DeepCopy()

	configMap.Name = getGatewayConfigMapName(gw.Object.Name)
	configMap.Namespace = gw.Object.Namespace
	configMap.Labels = labels

	setOwnerReference(gw.Object, configMap)

	if err := buildGatewayConfigData(gw, params, configMap); err != nil {
		return nil, err
	}

	return configMap, nil
}

func buildGatewayConfigData(
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
	configMap *coreV1.ConfigMap,
) error {
	servers, err := getServers(gw)
	if err != nil {
		return err
	}

	yaml := yaml.DBLess.DeepCopy()
	yaml.Put("servers", servers)
	yaml.Put("kafka", getKafkaServer(gw, params))
	yaml.Put("tags", BuildGatewayTag(gw))

	if httpPort := getHTTPPort(gw.Object.Spec.Listeners); httpPort != nil {
		yaml.Put("http", map[string]int32{"port": *httpPort})
	}

	yamlData, err := yaml.MarshalYAML()
	if err != nil {
		return err
	}
	configMap.Data[DefaultConfigFileEntry] = string(yamlData)
	return nil
}

func BuildGatewayTag(gw *gateway.Gateway) string {
	return BuildTag(gw.Object.Namespace, gw.Object.Name)
}

func BuildTag(namespace, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}

func getDeployment(
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
) (*appV1.Deployment, error) {
	labels := GwAPIv1GatewayLabels(gw.Object.Name)
	deployment := DefaultGatewayDeployment.DeepCopy()

	deployment.Name = gw.Object.Name
	deployment.Namespace = gw.Object.Namespace
	deployment.Labels = labels

	setOwnerReference(gw.Object, deployment)

	if err := buildDeployment(gw, params, deployment); err != nil {
		return nil, err
	}

	return deployment, nil
}

func buildDeployment(
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
	deployment *appV1.Deployment,
) error {
	labels := GwAPIv1GatewayLabels(gw.Object.Name)

	deployment.Spec.Selector.MatchLabels = labels
	deployment.Spec.Template.Labels = labels

	template, err := getPodTemplateSpec(gw, params)
	if err != nil {
		return err
	}

	prepareContainer(template, gw, params)

	template.Spec.Volumes = getVolumes(gw, params)

	template.Labels = labels

	deployment.Spec.Template = *template

	return nil
}

func getService(gw *gateway.Gateway, params *v1alpha1.GatewayClassParameters) *coreV1.Service {
	labels := GwAPIv1GatewayLabels(gw.Object.Name)
	svc := DefaultService.DeepCopy()

	svc.Name = gw.Object.Name
	svc.Namespace = gw.Object.Namespace
	svc.Labels = labels
	svc.Spec.Selector = labels

	buildServiceSpec(gw, params, svc)

	return svc
}

func buildServiceSpec(
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
	svc *coreV1.Service,
) {
	setServicePorts(svc, gw.Object.Spec.Listeners)

	setOwnerReference(gw.Object, svc)

	if params.Spec.Kubernetes == nil {
		return
	}

	if params.Spec.Kubernetes.Service == nil {
		return
	}

	svc.Spec.Type = *params.Spec.Kubernetes.Service.Type
}

func getServiceAccount(gw *gateway.Gateway) *coreV1.ServiceAccount {
	sa := &coreV1.ServiceAccount{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      gw.Object.Name,
			Namespace: gw.Object.Namespace,
		},
	}
	setOwnerReference(gw.Object, sa)
	return sa
}

func getRoleBinding(gw *gateway.Gateway) *rbacV1.RoleBinding {
	binding := &rbacV1.RoleBinding{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      gw.Object.Name,
			Namespace: gw.Object.Namespace,
		},
		Subjects: []rbacV1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      gw.Object.Name,
				Namespace: gw.Object.Namespace,
			},
		},
		RoleRef: rbacV1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     gw.Object.Name,
		},
	}
	setOwnerReference(gw.Object, binding)
	return binding
}

func getRole(gw *gateway.Gateway) *rbacV1.Role {
	role := &rbacV1.Role{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      gw.Object.Name,
			Namespace: gw.Object.Namespace,
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Verbs:     []string{"get", "list", "watch"},
				Resources: []string{"configmaps", "secrets"},
			},
		},
	}
	setOwnerReference(gw.Object, role)
	return role
}

func getPodTemplateSpec(
	gw *gateway.Gateway,
	parameters *v1alpha1.GatewayClassParameters,
) (*coreV1.PodTemplateSpec, error) {
	base := DefaultGatewayPodTemplate.DeepCopy()
	base.Spec.ServiceAccountName = gw.Object.Name

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
func getVolumes(
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
) []coreV1.Volume {
	configVolume := DefaultGatewayConfigVolume.DeepCopy()
	configVolume.ConfigMap.LocalObjectReference.Name = getGatewayConfigMapName(gw.Object.Name)
	volumes := []coreV1.Volume{*configVolume}

	graviteeParams := params.Spec.Gravitee
	if graviteeParams != nil && graviteeParams.LicenseRef != nil {
		licenseRef := graviteeParams.LicenseRef
		licenseVolume := DefaultLicenseConfigVolume.DeepCopy()
		licenseVolume.Secret.SecretName = string(licenseRef.Name)
		return append(volumes, *licenseVolume)
	}

	return volumes
}

func getVolumeMounts(
	params *v1alpha1.GatewayClassParameters,
) []coreV1.VolumeMount {
	vm := []coreV1.VolumeMount{ConfigVolumeMount}

	graviteeParams := params.Spec.Gravitee
	if graviteeParams != nil && graviteeParams.LicenseRef != nil {
		return append(vm, LicenseVolumeMount)
	}

	return vm
}

func prepareContainer(
	template *coreV1.PodTemplateSpec,
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
) {
	container := getGatewayContainer(template.Spec.Containers)
	container.VolumeMounts = getVolumeMounts(params)
	setContainerPorts(container, gw.Object.Spec.Listeners)
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
	blockOwnerDeletion := true
	obj.SetOwnerReferences(
		[]metaV1.OwnerReference{
			{
				Kind:               kind,
				APIVersion:         version,
				Name:               gw.GetName(),
				UID:                gw.GetUID(),
				BlockOwnerDeletion: &blockOwnerDeletion,
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

func getServers(gw *gateway.Gateway) ([]map[string]any, error) {
	listeners := gw.Object.Spec.Listeners
	statuses := gw.Object.Status.Listeners
	servers := make([]map[string]any, len(gw.Object.Spec.Listeners))
	for i, listener := range listeners {
		if IsKafkaListener(listener) {
			continue
		}
		status := gateway.WrapListenerStatus(&statuses[i])
		if !IsAccepted(status) {
			continue
		}
		server := make(map[string]any)
		serverType, ok := ProtocolToServerType[listener.Protocol]
		if !ok {
			return nil, fmt.Errorf("unknown server protocol")
		}
		server["type"] = serverType
		server["port"] = listener.Port
		if listener.Hostname != nil {
			server["hostname"] = *listener.Hostname
		} else {
			server["hostname"] = "0.0.0.0"
		}
		if listener.TLS != nil {
			server["secured"] = true
			server["ssl"] = buildTLS(listener)
		}
		servers[i] = server
	}
	return servers, nil
}

func getKafkaServer(
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
) map[string]any {
	kafkaParams := params.Spec.Gravitee.Kafka
	if !kafkaParams.Enabled {
		return yaml.Kafka.Object
	}

	kafkaListener := GetKafkaListener(gw)
	if kafkaListener == nil {
		return yaml.Kafka.Object
	}

	kafkaStatus := gateway.WrapListenerStatus(GetKafkaListenerStatus(gw))
	if !IsAccepted(kafkaStatus) {
		return yaml.Kafka.Object
	}

	kafka := yaml.Kafka.DeepCopy()

	kafka.Put("enabled", true)

	routingHostModeParams := kafkaParams.RoutingHostMode
	kafka.Put(
		"routingHostMode",
		map[string]any{
			"defaultDomain":   *kafkaListener.Hostname,
			"defaultPort":     kafkaListener.Port,
			"brokerPrefix":    routingHostModeParams.BrokerPrefix,
			"domainSeparator": routingHostModeParams.DomainSeparator,
		},
	)

	kafka.Put("port", kafkaListener.Port)

	kafka.Put(
		"ssl",
		buildTLS(*kafkaListener),
	)

	return kafka.Object
}

func buildTLS(listener gwAPIv1.Listener) map[string]any {
	tls := listener.TLS
	ssl := make(map[string]any)
	keystore := make(map[string]any)
	keystore["type"] = "pem"
	keystore["secret"] = buildKeystoreSecret(tls.CertificateRefs)
	ssl["keystore"] = keystore
	return ssl
}

func buildKeystoreSecret(certificateRefs []gwAPIv1.SecretObjectReference) string {
	ref := certificateRefs[0]
	return "secret://kubernetes/" + string(ref.Name)
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

func getGatewayConfigMapName(gwName string) string {
	return GatewayConfigMapPrefix + gwName
}

func getPEMRegistryConfigMapName(gwName string) string {
	return PEMRegistryConfigMapPrefix + gwName
}
