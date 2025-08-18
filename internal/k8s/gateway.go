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
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/gateway/logback"
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

const portMin = int32(1024)

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

	portMapping := mapPorts(gw.Object)

	gatewayConfig, err := getOverridingGatewayConfigMap(ctx, gw, params, portMapping)
	if err != nil {
		return err
	} else if err := CreateOrUpdate(ctx, gatewayConfig, func() error {
		return buildOverridingGatewayConfigData(ctx, gw, params, gatewayConfig, portMapping)
	}); err != nil {
		return err
	}

	configData := map[string]string{}

	if HasGraviteeYAML(params) {
		userConfig, err := getUserConfigMap(gw, params)
		if err != nil {
			return err
		}
		if err := CreateOrUpdate(ctx, userConfig, func() error {
			return buildUserGatewayConfigData(params, userConfig)
		}); err != nil {
			return err
		}
		maps.Copy(configData, userConfig.Data)
	}

	maps.Copy(configData, gatewayConfig.Data)

	if deployment, err := getDeployment(gw, params, portMapping, configData); err != nil {
		return err
	} else if err := CreateOrUpdate(ctx, deployment, func() error {
		return buildDeployment(gw, params, deployment, portMapping, configData)
	}); err != nil {
		return err
	}

	svc := getService(gw, params, portMapping)
	return CreateOrUpdate(ctx, svc, func() error {
		buildServiceSpec(gw, params, svc, portMapping)
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

func getUserConfigMap(gw *gateway.Gateway, params *v1alpha1.GatewayClassParameters) (*coreV1.ConfigMap, error) {
	labels := GwAPIv1GatewayLabels(gw.Object.Name)
	configMap := DefaultGatewayConfigMap.DeepCopy()
	configMap.Name = getUserGatewayConfigMapName(gw.Object.Name)
	configMap.Namespace = gw.Object.Namespace
	configMap.Labels = labels

	setOwnerReference(gw.Object, configMap)

	if err := buildUserGatewayConfigData(params, configMap); err != nil {
		return nil, err
	}

	return configMap, nil
}

func getOverridingGatewayConfigMap(
	ctx context.Context,
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
	portMapping map[gwAPIv1.PortNumber]int32,
) (*coreV1.ConfigMap, error) {
	labels := GwAPIv1GatewayLabels(gw.Object.Name)
	configMap := DefaultGatewayConfigMap.DeepCopy()

	configMap.Name = getGatewayConfigMapName(gw.Object.Name)
	configMap.Namespace = gw.Object.Namespace
	configMap.Labels = labels

	setOwnerReference(gw.Object, configMap)

	if err := buildOverridingGatewayConfigData(ctx, gw, params, configMap, portMapping); err != nil {
		return nil, err
	}

	return configMap, nil
}

func buildUserGatewayConfigData(
	params *v1alpha1.GatewayClassParameters,
	configMap *coreV1.ConfigMap,
) error {
	yamlObj := params.Spec.Gravitee.YAML
	yamlData, err := yamlObj.MarshalYAML()
	if err != nil {
		return err
	}
	configMap.Data[UserConfigFileEntry] = string(yamlData)
	return nil
}

func buildOverridingGatewayConfigData(
	ctx context.Context,
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
	configMap *coreV1.ConfigMap,
	portMapping map[gwAPIv1.PortNumber]int32,
) error {
	servers, err := getServers(ctx, gw, portMapping)
	if err != nil {
		return err
	}

	yaml := yaml.DBLess.DeepCopy()
	yaml.Put("servers", servers)
	yaml.Put("kafka", getKafkaServer(gw, params, portMapping))
	yaml.Put("tags", BuildGatewayTag(gw))

	if httpPort := getHTTPPort(gw.Object.Spec.Listeners, portMapping); httpPort != nil {
		yaml.Put("http", map[string]int32{"port": *httpPort})
	}

	yamlData, err := yaml.MarshalYAML()
	if err != nil {
		return err
	}
	configMap.Data[DefaultConfigFileEntry] = string(yamlData)

	configMap.Data[DefaultLogConfigFileEntry] = logback.Config

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
	portMapping map[gwAPIv1.PortNumber]int32,
	config map[string]string,
) (*appV1.Deployment, error) {
	labels := GwAPIv1GatewayLabels(gw.Object.Name)
	deployment := DefaultGatewayDeployment.DeepCopy()

	deployment.Name = gw.Object.Name
	deployment.Namespace = gw.Object.Namespace
	deployment.Labels = labels

	setOwnerReference(gw.Object, deployment)

	if err := buildDeployment(gw, params, deployment, portMapping, config); err != nil {
		return nil, err
	}

	return deployment, nil
}

func buildDeployment(
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
	deployment *appV1.Deployment,
	portMapping map[gwAPIv1.PortNumber]int32,
	config map[string]string,
) error {
	labels := GwAPIv1GatewayLabels(gw.Object.Name)

	deployment.Spec.Selector.MatchLabels = labels
	deployment.Spec.Template.Labels = labels

	template, err := getPodTemplateSpec(gw, params)
	if err != nil {
		return err
	}

	prepareContainer(template, gw, params, portMapping)

	template.Spec.Volumes = getVolumes(gw, params)

	template.Labels = labels

	deployment.Spec.Template = *template

	if params.Spec.Kubernetes == nil {
		return nil
	}
	if params.Spec.Kubernetes.Deployment == nil {
		return nil
	}

	if params.Spec.Kubernetes.Deployment.Replicas != nil {
		deployment.Spec.Replicas = params.Spec.Kubernetes.Deployment.Replicas
	}

	if params.Spec.Kubernetes.Deployment.Strategy != nil {
		deployment.Spec.Strategy = *params.Spec.Kubernetes.Deployment.Strategy
	}

	if deployment.Annotations == nil {
		deployment.Annotations = make(map[string]string)
	}

	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}

	if deployment.Labels == nil {
		deployment.Labels = make(map[string]string)
	}

	if deployment.Spec.Template.Labels == nil {
		deployment.Spec.Template.Labels = make(map[string]string)
	}

	for k, v := range params.Spec.Kubernetes.Deployment.Annotations {
		deployment.Annotations[k] = v
		deployment.Spec.Template.Annotations[k] = v
	}

	for k, v := range params.Spec.Kubernetes.Deployment.Labels {
		deployment.Labels[k] = v
		deployment.Spec.Template.Labels[k] = v
	}

	deployment.Spec.Template.Annotations["gravitee.io/config"] = hash.Calculate(config)

	return nil
}

func getService(
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
	portMapping map[gwAPIv1.PortNumber]int32,
) *coreV1.Service {
	labels := GwAPIv1GatewayLabels(gw.Object.Name)
	svc := DefaultService.DeepCopy()

	svc.Name = gw.Object.Name
	svc.Namespace = gw.Object.Namespace
	svc.Labels = labels
	svc.Spec.Selector = labels

	buildServiceSpec(gw, params, svc, portMapping)

	return svc
}

func buildServiceSpec(
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
	svc *coreV1.Service,
	portMapping map[gwAPIv1.PortNumber]int32,
) {
	setServicePorts(svc, gw.Object.Spec.Listeners, portMapping)

	setOwnerReference(gw.Object, svc)

	if params.Spec.Kubernetes == nil {
		return
	}

	if params.Spec.Kubernetes.Service == nil {
		return
	}

	svc.Spec.Type = *params.Spec.Kubernetes.Service.Type
	svc.Spec.ExternalTrafficPolicy = params.Spec.Kubernetes.Service.ExternalTrafficPolicy
	svc.Spec.LoadBalancerClass = params.Spec.Kubernetes.Service.LoadBalancerClass

	for k, v := range params.Spec.Kubernetes.Service.Annotations {
		svc.Annotations[k] = v
	}

	for k, v := range params.Spec.Kubernetes.Deployment.Labels {
		svc.Labels[k] = v
	}
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

func getRoleBinding(gw *gateway.Gateway) *rbacV1.ClusterRoleBinding {
	binding := &rbacV1.ClusterRoleBinding{
		ObjectMeta: metaV1.ObjectMeta{
			Name: gw.Object.Name,
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
			Kind:     "ClusterRole",
			Name:     gw.Object.Name,
		},
	}
	setOwnerReference(gw.Object, binding)
	return binding
}

func getRole(gw *gateway.Gateway) *rbacV1.ClusterRole {
	role := &rbacV1.ClusterRole{
		ObjectMeta: metaV1.ObjectMeta{
			Name: gw.Object.Name,
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

	if HasGraviteeYAML(parameters) {
		setGraviteeConf(base)
	}

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

func setGraviteeConf(podTemplate *coreV1.PodTemplateSpec) {
	containers := podTemplate.Spec.Containers
	for i, container := range containers {
		if container.Name == GatewayContainerName {
			containers[i].Env = append(containers[i].Env, coreV1.EnvVar{
				Name: "JAVA_OPTS",
				Value: fmt.Sprintf(
					"-Dgravitee.conf=%s,%s",
					UserGatewayConfigFile,
					DefaultGatewayConfigFile,
				),
			})
		}
	}
}

func getVolumes(
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
) []coreV1.Volume {
	configVolume := DefaultGatewayConfigVolume.DeepCopy()
	configVolume.ConfigMap.LocalObjectReference.Name = getGatewayConfigMapName(gw.Object.Name)
	volumes := []coreV1.Volume{*configVolume}

	if HasGraviteeLicense(params) {
		licenseRef := params.Spec.Gravitee.LicenseRef
		licenseVolume := DefaultLicenseConfigVolume.DeepCopy()
		licenseVolume.Secret.SecretName = string(licenseRef.Name)
		volumes = append(volumes, *licenseVolume)
	}

	if HasGraviteeYAML(params) {
		userVolume := UserGatewayConfigVolume.DeepCopy()
		userVolume.ConfigMap.LocalObjectReference.Name = getUserGatewayConfigMapName(gw.Object.Name)
		volumes = append(volumes, *userVolume)
	}

	return volumes
}

func getVolumeMounts(
	params *v1alpha1.GatewayClassParameters,
) []coreV1.VolumeMount {
	vm := []coreV1.VolumeMount{ConfigVolumeMount, LogConfigVolumeMount}

	if HasGraviteeLicense(params) {
		vm = append(vm, LicenseVolumeMount)
	}

	if HasGraviteeYAML(params) {
		vm = append(vm, UserConfigVolumeMount)
	}

	return vm
}

func prepareContainer(
	template *coreV1.PodTemplateSpec,
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
	portMapping map[gwAPIv1.PortNumber]int32,
) {
	container := getGatewayContainer(template.Spec.Containers)
	container.VolumeMounts = getVolumeMounts(params)
	setContainerPorts(container, gw.Object.Spec.Listeners, portMapping)
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

func setContainerPorts(
	container *coreV1.Container,
	listeners []gwAPIv1.Listener,
	portMapping map[gwAPIv1.PortNumber]int32,
) {
	ports := make([]coreV1.ContainerPort, 0)
	knownPorts := sets.New[int32]()
	for _, listener := range listeners {
		mappedPort := portMapping[listener.Port]
		if !knownPorts.Has(mappedPort) {
			port := coreV1.ContainerPort{
				ContainerPort: mappedPort,
				Protocol:      coreV1.ProtocolTCP,
			}
			ports = append(ports, port)
		}
		knownPorts.Insert(mappedPort)
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

func setServicePorts(
	svc *coreV1.Service,
	listeners []gwAPIv1.Listener,
	portMapping map[gwAPIv1.PortNumber]int32,
) {
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
			TargetPort: intstr.FromInt32(portMapping[listener.Port]),
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

func getHTTPPort(listeners []gwAPIv1.Listener, portMapping map[gwAPIv1.PortNumber]int32) *int32 {
	port := int32(0)
	for _, listener := range listeners {
		if listener.Protocol == gwAPIv1.HTTPProtocolType {
			port = portMapping[listener.Port]
		}
		if port == 0 && listener.Protocol == gwAPIv1.HTTPSProtocolType {
			port = portMapping[listener.Port]
		}
	}
	if port > 0 {
		return &port
	}
	return nil
}

func getServers(
	_ context.Context,
	gw *gateway.Gateway,
	portMapping map[gwAPIv1.PortNumber]int32,
) ([]map[string]any, error) {
	listeners := gw.Object.Spec.Listeners
	statuses := gw.Object.Status.Listeners
	servers := make([]map[string]any, 0)
	knownPorts := sets.New[int32]()
	for i, listener := range listeners {
		if IsKafkaListener(listener) {
			continue
		}

		status := gateway.WrapListenerStatus(&statuses[i])
		if !IsAccepted(status) {
			continue
		}

		serverPort := portMapping[listener.Port]

		if knownPorts.Has(serverPort) {
			continue
		}

		server := make(map[string]any)
		serverType, ok := ProtocolToServerType[listener.Protocol]
		if !ok {
			return nil, fmt.Errorf("unknown server protocol")
		}
		server["type"] = serverType
		server["port"] = serverPort
		server["host"] = "0.0.0.0"

		if listener.TLS != nil {
			server["secured"] = true
			server["ssl"] = buildTLS(listener)
		}
		servers = append(servers, server)
		knownPorts.Insert(serverPort)
	}
	return servers, nil
}

func getKafkaServer(
	gw *gateway.Gateway,
	params *v1alpha1.GatewayClassParameters,
	portMapping map[gwAPIv1.PortNumber]int32,
) map[string]any {
	if !HasKafkaEnabled(params) {
		return yaml.Kafka.Object
	}

	kafkaListener := GetKafkaListener(gw.Object)
	if kafkaListener == nil {
		return yaml.Kafka.Object
	}

	kafkaStatus := gateway.WrapListenerStatus(GetKafkaListenerStatus(gw))
	if !IsAccepted(kafkaStatus) {
		return yaml.Kafka.Object
	}

	kafka := yaml.Kafka.DeepCopy()

	kafka.Put("enabled", true)

	kafkaParams := params.Spec.Gravitee.Kafka

	routingHostModeParams := kafkaParams.RoutingHostMode
	kafka.Put(
		"routingHostMode",
		map[string]any{
			"defaultPort":            portMapping[kafkaListener.Port],
			"brokerDomainPattern":    routingHostModeParams.BokerDomainPattern,
			"bootstrapDomainPattern": routingHostModeParams.BootstrapDomainPattern,
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

	acceptUserConflictingChanges(base, patch)

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

func acceptUserConflictingChanges(tmpl *coreV1.PodTemplateSpec, userTmpl *coreV1.PodTemplateSpec) {
	if userTmpl == nil {
		return
	}

	for i := range userTmpl.Spec.Containers {
		uC := &userTmpl.Spec.Containers[i]
		for j := range tmpl.Spec.Containers {
			c := &tmpl.Spec.Containers[j]
			if c.Name == uC.Name {
				givePrecedenceToUserProbes(c, uC)
			}
		}
	}
}

func givePrecedenceToUserProbes(container *coreV1.Container, userContainer *coreV1.Container) {
	givePrecedenceToUserProbe(container.StartupProbe, userContainer.StartupProbe)
	givePrecedenceToUserProbe(container.ReadinessProbe, userContainer.ReadinessProbe)
	givePrecedenceToUserProbe(container.LivenessProbe, userContainer.LivenessProbe)
}

func givePrecedenceToUserProbe(probe *coreV1.Probe, userProbe *coreV1.Probe) {
	switch {
	case userProbe == nil || probe == nil:
		return
	case userProbe.Exec != nil:
		probe.Exec = userProbe.Exec
		probe.GRPC = nil
		probe.HTTPGet = nil
		probe.TCPSocket = nil
	case userProbe.GRPC != nil:
		probe.GRPC = userProbe.GRPC
		probe.Exec = nil
		probe.HTTPGet = nil
		probe.TCPSocket = nil
	case userProbe.HTTPGet != nil:
		probe.HTTPGet = userProbe.HTTPGet
		probe.Exec = nil
		probe.GRPC = nil
		probe.TCPSocket = nil
	case userProbe.TCPSocket != nil:
		probe.TCPSocket = userProbe.TCPSocket
		probe.Exec = nil
		probe.GRPC = nil
		probe.HTTPGet = nil
	}
}

func getGatewayConfigMapName(gwName string) string {
	return GatewayConfigMapPrefix + gwName
}

func getUserGatewayConfigMapName(gwName string) string {
	return "user-" + GatewayConfigMapPrefix + gwName
}

func getPEMRegistryConfigMapName(gwName string) string {
	return PEMRegistryConfigMapPrefix + gwName
}

func mapPorts(gw *gwAPIv1.Gateway) map[gwAPIv1.PortNumber]int32 {
	listenerToInternal := make(map[gwAPIv1.PortNumber]int32)
	for _, l := range gw.Spec.Listeners {
		listenerToInternal[l.Port] = ensureBindablePort(gw, l.Port)
	}
	return listenerToInternal
}

func getListenerPorts(gw *gwAPIv1.Gateway) sets.Set[int32] {
	ports := sets.New[int32]()
	for _, l := range gw.Spec.Listeners {
		ports.Insert(int32(l.Port))
	}
	return ports
}

func ensureBindablePort(
	gw *gwAPIv1.Gateway, port gwAPIv1.PortNumber,
) int32 {
	val := int32(port)
	if val > portMin {
		return val
	}
	return ensureBindableWithNoConflict(val, getListenerPorts(gw))
}

func ensureBindableWithNoConflict(
	port int32,
	listenerPorts sets.Set[int32],
) int32 {
	bindable := port + portMin
	if listenerPorts.Has(bindable) {
		return ensureBindableWithNoConflict(bindable+1, listenerPorts)
	}
	return bindable
}
