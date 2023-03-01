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
	endpointNamePattern    = "rule%02d-path%02d"
	endpointMatcherPattern = "%s:{#group[0]}"
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

// This wrapper is used to compute the endpoint name
// used for target selection in the routing policy.
type indexedPath struct {
	*v1.HTTPIngressPath
	ruleIndex int
	index     int
}

func (p indexedPath) String() string {
	ruleIndex := p.ruleIndex + 1
	pathIndex := p.index + 1
	return fmt.Sprintf(endpointNamePattern, ruleIndex, pathIndex)
}

func newIndexedPath(path *v1.HTTPIngressPath, ruleIndex, index int) *indexedPath {
	return &indexedPath{
		HTTPIngressPath: path,
		ruleIndex:       ruleIndex,
		index:           index,
	}
}

// Map maps an ingress to a graviteeio API definition, adding one virtual host per ingress rule,
// one endpoint per backend service, and one conditional flow per host and path of the rule.
// The host header is used to select the flow, and a dynamic routing policy routes the request
// to the backend service, identified by the host and path of the rule.
func (m *Mapper) Map(apiDefinition *gio.ApiDefinition, ingress *v1.Ingress) *gio.ApiDefinition {
	m.hosts = getHosts(ingress)
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
	flows := make([]model.Flow, 0)
	for ruleIndex, rule := range rules {
		flows = append(flows, m.buildPathFlows(rule, ruleIndex)...)
	}
	return flows
}

func (m *Mapper) buildPathFlows(rule v1.IngressRule, ruleIndex int) []model.Flow {
	flows := make([]model.Flow, 0)
	for i := range rule.HTTP.Paths {
		path := rule.HTTP.Paths[i]
		flows = append(flows, m.buildConditionalFlow(rule, newIndexedPath(&path, ruleIndex, i)))
	}
	return flows
}

// Init a conditional flow matching a given HTTP path of a given ingress rule.
// The flow will match the ingress path as a path operator and define a condition
// based on the host of the rule. If no host is defined for the rule, then
// the condition will check that none of the host we have processed matches the Host header
// of the incoming request.
func (m *Mapper) buildConditionalFlow(rule v1.IngressRule, path *indexedPath) model.Flow {
	flow := model.Flow{}
	flow.Name = rule.Host + path.Path
	flow.PathOperator = buildPathOperator(path)
	flow.Pre = buildRouting(path)

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

func buildRoutingStep(path *indexedPath) model.FlowStep {
	return model.FlowStep{
		Name:    routingStepName,
		Policy:  routingPolicyName,
		Enabled: true,
		Configuration: &model.GenericStringMap{
			Unstructured: unstructured.Unstructured{
				Object: map[string]interface{}{
					routingRulesKey: buildRoutingRules(path),
				},
			},
		},
	}
}

func buildRoutingRules(path *indexedPath) []interface{} {
	return []interface{}{
		map[string]interface{}{
			routingPatternKey: routingPattern,
			routingUrlKey:     buildRoutingTarget(path),
		},
	}
}

func buildRoutingTarget(path *indexedPath) string {
	return fmt.Sprintf(endpointMatcherPattern, path.String())
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
	for ruleIndex, rule := range ingress.Spec.Rules {
		for pathIndex := range rule.HTTP.Paths {
			path := &rule.HTTP.Paths[pathIndex]
			eps = append(eps, buildEndpoint(ingress, newIndexedPath(path, ruleIndex, pathIndex)))
		}
	}
	return eps
}

// For each rule and path of an ingress, build an endpoint identified by the host and path,
// in order to be able to match it in the routing step when handling an incoming request for routing.
func buildEndpoint(ingress *v1.Ingress, path *indexedPath) *model.HttpEndpoint {
	return &model.HttpEndpoint{
		Name:   path.String(),
		Target: buildEndpointTarget(ingress, path),
	}
}

func buildEndpointTarget(ingress *v1.Ingress, path *indexedPath) string {
	svc := path.Backend.Service
	return fmt.Sprintf(serviceURIPattern, svc.Name, ingress.Namespace, svc.Port.Number)
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
func buildRouting(path *indexedPath) []model.FlowStep {
	return append([]model.FlowStep{}, buildRoutingStep(path))
}

func buildPathOperator(path *indexedPath) *model.PathOperator {
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

// Get all the host names defined in the ingress rules,
// in order to compute the condition for rules with no host,
// checking that none of the hosts we have processed matches the
// host header of the incoming request.
func getHosts(ingress *v1.Ingress) []string {
	hosts := make([]string, 0)
	for _, rule := range ingress.Spec.Rules {
		if rule.Host != "" {
			hosts = append(hosts, rule.Host)
		}
	}
	return hosts
}
