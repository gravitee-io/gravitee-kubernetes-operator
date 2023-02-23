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

package mapper

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	proxyName              = "default"
	serviceURIPattern      = "http://%s.%s.svc.cluster.local:%d"
	flowNamePattern        = "%s%s"
	routingPattern         = "(.*)"
	routingPolicyName      = "dynamic-routing"
	routingStepName        = "Ingress Routing"
	routingRulesKey        = "rules"
	routingPatternKey      = "pattern"
	routingUrlKey          = "url"
	endpointMatcherPattern = "{#endpoints['%s']}{#group[0]}"
	hostConditionPattern   = "{#request.headers['Host'][0] == '%s'}"
	noHostConditionPattern = "#request.headers['Host'][0] != '%s' && "
	rootPath               = "/"
)

type Mapper struct {
	hosts []string
}

func New() *Mapper {
	return &Mapper{
		hosts: make([]string, 0),
	}
}

// Map maps an ingress to a graviteeio API definition, adding one virtual host per ingress rule,
// one endpoint per backend service, and one conditional flow per host and path of the rule.
// The host header is used to select the flow, and a dynamic routing policy routes the request
// to the backend service, identified by the host and path of the rule.
func (m *Mapper) Map(apiDefinition *gio.ApiDefinition, ingress *v1.Ingress) *gio.ApiDefinition {
	cp := buildApiCopy(apiDefinition, ingress)
	cp.Spec.Proxy = buildProxy(ingress)
	cp.Spec.Flows = m.buildFlows(ingress.Spec.Rules)
	return cp
}

func buildApiCopy(apiDefinition *gio.ApiDefinition, ingress *v1.Ingress) *gio.ApiDefinition {
	return &gio.ApiDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingress.Name,
			Namespace: ingress.Namespace,
			Annotations: map[string]string{
				keys.Extends: keys.IngressLabel,
			},
		},
		Spec: *apiDefinition.Spec.DeepCopy(),
	}
}

func (m *Mapper) buildFlows(rules []v1.IngressRule) []model.Flow {
	cp := sortHostsFirst(rules)
	flows := make([]model.Flow, 0)
	for _, rule := range cp {
		flows = append(flows, m.buildPathFlows(rule)...)
	}
	return flows
}

func (m *Mapper) buildPathFlows(rule v1.IngressRule) []model.Flow {
	flows := make([]model.Flow, 0)
	for _, path := range rule.HTTP.Paths {
		flows = append(flows, m.buildConditionalFlow(rule, path))
	}
	return flows
}

// Init a conditional flow matching a given HTTP path of a given ingress rule.
// The flow will match the ingress path as a path operator and define a condition
// based on the host of the rule. If no host is defined for the rule, then
// the condition will check that none of the host we have processed matches the Host header
// of the incoming request.
func (m *Mapper) buildConditionalFlow(rule v1.IngressRule, path v1.HTTPIngressPath) model.Flow {
	flow := model.Flow{}
	flow.Name = buildEndpointName(rule.Host, path.Path)
	flow.PathOperator = buildPathOperator(path)
	flow.Pre = buildRouting(rule, path)

	if rule.Host == "" {
		flow.Condition = m.buildNoHostCondition()
		return flow
	}

	flow.Condition = buildHostCondition(rule)
	m.hosts = append(m.hosts, rule.Host)
	return flow
}

func (m *Mapper) buildNoHostCondition() string {
	condition := "{"
	for _, host := range m.hosts {
		condition += fmt.Sprintf(noHostConditionPattern, host)
	}
	return condition[:len(condition)-4] + "}"
}

func buildHostCondition(rule v1.IngressRule) string {
	return fmt.Sprintf(hostConditionPattern, rule.Host)
}

func buildRoutingStep(rule v1.IngressRule, path v1.HTTPIngressPath) model.FlowStep {
	return model.FlowStep{
		Name:    routingStepName,
		Policy:  routingPolicyName,
		Enabled: true,
		Configuration: &model.GenericStringMap{
			Unstructured: unstructured.Unstructured{
				Object: map[string]interface{}{
					routingRulesKey: buildRoutingRules(rule, path),
				},
			},
		},
	}
}

func buildRoutingRules(rule v1.IngressRule, path v1.HTTPIngressPath) []interface{} {
	return []interface{}{
		map[string]interface{}{
			routingPatternKey: routingPattern,
			routingUrlKey:     buildRoutingTarget(rule.Host, path.Path),
		},
	}
}

func buildRoutingTarget(host, path string) string {
	return fmt.Sprintf(endpointMatcherPattern, buildEndpointName(host, path))
}

func buildProxy(ingress *v1.Ingress) *model.Proxy {
	return &model.Proxy{
		VirtualHosts: buildVirtualHosts(ingress),
		Groups: []*model.EndpointGroup{
			{
				Name:      proxyName,
				Endpoints: buildEndpoints(ingress),
			},
		},
	}
}

func buildEndpoints(ingress *v1.Ingress) []*model.HttpEndpoint {
	eps := make([]*model.HttpEndpoint, 0)
	for _, rule := range ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			eps = append(eps, buildEndpoint(ingress, rule, path))
		}
	}
	return eps
}

// For each rule and path of an ingress, build an endpoint identified by the host and path,
// in order to be able to match it in the routing step when handling an incoming request for routing.
func buildEndpoint(ingress *v1.Ingress, rule v1.IngressRule, path v1.HTTPIngressPath) *model.HttpEndpoint {
	return &model.HttpEndpoint{
		Name:   buildEndpointName(rule.Host, path.Path),
		Target: buildEndpointTarget(ingress, path),
	}
}

func buildEndpointTarget(ingress *v1.Ingress, path v1.HTTPIngressPath) string {
	svc := path.Backend.Service
	return fmt.Sprintf(serviceURIPattern, svc.Name, ingress.Namespace, svc.Port.Number)
}

func buildEndpointName(host, path string) string {
	rawName := host + path
	return strings.TrimPrefix(rawName, "/")
}

// For each ingress host and path, build a virtual host.
func buildVirtualHosts(ingress *v1.Ingress) []*model.VirtualHost {
	vhs := make([]*model.VirtualHost, 0)
	for _, rule := range ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			vhs = append(vhs, buildVirtualHost(rule, path))
		}
	}
	return vhs
}

func buildVirtualHost(rule v1.IngressRule, path v1.HTTPIngressPath) *model.VirtualHost {
	return &model.VirtualHost{Host: rule.Host, Path: path.Path}
}

// Builds a flow step with a single routing policy that will forward the request to
// its matching endpoint using the request ingress host and path as a endpoint name.
func buildRouting(rule v1.IngressRule, path v1.HTTPIngressPath) []model.FlowStep {
	return append([]model.FlowStep{}, buildRoutingStep(rule, path))
}

func buildPathOperator(path v1.HTTPIngressPath) *model.PathOperator {
	if *path.PathType == v1.PathTypeExact {
		return &model.PathOperator{
			Operator: model.EqualsOperator,
			Path:     rootPath,
		}
	}
	return &model.PathOperator{
		Operator: model.StartWithOperator,
		Path:     rootPath,
	}
}

// Sort the rules so that the ones with a host appear first,
// in order to compute the condition for rules with no host,
// checking that none of the hosts we have processed matches the
// host header of the incoming request.
func sortHostsFirst(rules []v1.IngressRule) []v1.IngressRule {
	cp := make([]v1.IngressRule, len(rules))
	copy(cp, rules)
	sort.Slice(cp, func(i, j int) bool {
		return cp[i].Host != ""
	})
	return cp
}
