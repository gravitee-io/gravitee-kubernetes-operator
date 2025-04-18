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
	"net/http"
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/el"
	xhttp "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	proxyName              = "default"
	serviceURIPattern      = "http://%s.%s.svc.cluster.local:%d"
	routingPattern         = "(.*)"
	routingPolicyName      = "dynamic-routing"
	routingStepName        = "Ingress Routing"
	routingRulesKey        = "rules"
	routingPatternKey      = "pattern"
	routingUrlKey          = "url"
	mockPolicyName         = "mock"
	mockStepName           = "No Route Found"
	mockContentKey         = "content"
	mockStatusKey          = "status"
	mockHeadersKey         = "headers"
	endpointNamePattern    = "rule%02d-path%02d"
	endpointMatcherPattern = "%s:{#group[0]}"
	rootPath               = "/"
)

var hostCondition = el.Expression("#request.host == '%s'")
var noHostCondition = el.Expression("#request.host != '%s'")
var pathCondition = el.Expression("#request.contextPath.startsWith('%s')")

type Mapper struct {
	opts       Opts
	hosts      map[string]bool
	conditions []el.Expression
}

func New(opts Opts) *Mapper {
	return &Mapper{
		opts:       mergeOpts(opts),
		hosts:      make(map[string]bool),
		conditions: make([]el.Expression, 0),
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
// to the backend service, identified by the endpoint name. Is no rule matches,
// a 404 response is returned by a flow that negates all the previous conditions.
func (m *Mapper) Map(apiDefinition *v1alpha1.ApiDefinition, ingress *v1.Ingress) *v1alpha1.ApiDefinition {
	m.hosts = getHosts(ingress)
	cp := buildApiCopy(apiDefinition, ingress)
	cp.Spec.Proxy = buildProxy(ingress)
	cp.Spec.Flows = m.buildFlows(ingress.Spec.Rules)
	if apiDefinition.Spec.Flows != nil {
		cp.Spec.FlowMode = v2.DefaultFlowMode
		cp.Spec.Flows = append(cp.Spec.Flows, apiDefinition.Spec.Flows...)
	}
	return cp
}

// Get all the host names defined in the ingress rules,
// in order to compute the condition for rules with no host,
// checking that none of the hosts we have processed matches the
// host header of the incoming request.
func getHosts(ingress *v1.Ingress) map[string]bool {
	hosts := make(map[string]bool)
	for _, rule := range ingress.Spec.Rules {
		if rule.Host != "" {
			hosts[rule.Host] = true
		}
	}
	return hosts
}

func buildApiCopy(apiDefinition *v1alpha1.ApiDefinition, ingress *v1.Ingress) *v1alpha1.ApiDefinition {
	spec := *apiDefinition.Spec.DeepCopy()
	spec.Name = ingress.Name
	spec.Description = core.IngressLabel
	spec.Version = v1alpha1.GroupVersion.Version

	return &v1alpha1.ApiDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingress.Name,
			Namespace: ingress.Namespace,
			Annotations: map[string]string{
				core.Extends: core.IngressLabel,
			},
		},
		Spec: spec,
	}
}

func (m *Mapper) buildFlows(rules []v1.IngressRule) []v2.Flow {
	flows := make([]v2.Flow, 0)
	for ruleIndex, rule := range rules {
		flows = append(flows, m.buildPathFlows(rule, ruleIndex)...)
	}
	return append(flows, m.buildNotFoundFlow())
}

func (m *Mapper) buildPathFlows(rule v1.IngressRule, ruleIndex int) []v2.Flow {
	flows := make([]v2.Flow, 0)
	for i := range rule.HTTP.Paths {
		path := rule.HTTP.Paths[i]
		flows = append(flows, m.buildRoutingFlow(rule, newIndexedPath(&path, ruleIndex, i)))
	}
	return flows
}

// Init a conditional flow matching a given HTTP path of a given ingress rule.
// The flow will match the ingress path as a path operator and define a condition
// based on the host of the rule. If no host is defined for the rule, then
// the condition will check that none of the host we have processed matches the Host header
// of the incoming request.
func (m *Mapper) buildRoutingFlow(rule v1.IngressRule, path *indexedPath) v2.Flow {
	flow := v2.Flow{Enabled: true}
	name := fmt.Sprintf("%s%s", rule.Host, path.Path)
	flow.Name = &name
	flow.PathOperator = buildPathOperator(path)
	flow.Pre = buildRouting(path)

	if rule.Host == "" {
		flow.Condition = m.buildNoHostCondition(path)
		return flow
	}

	flow.Condition = m.buildHostCondition(rule, path)
	return flow
}

func buildPathOperator(path *indexedPath) *v2.PathOperator {
	if *path.PathType == v1.PathTypeExact {
		return &v2.PathOperator{
			Operator: base.EqualsOperator,
			Path:     toPointer(rootPath),
		}
	}
	return &v2.PathOperator{
		Operator: base.StartWithOperator,
		Path:     toPointer(rootPath),
	}
}

func buildRouting(path *indexedPath) []base.FlowStep {
	return append([]base.FlowStep{}, buildRoutingStep(path))
}

func (m *Mapper) buildNoHostCondition(path *indexedPath) *string {
	condition := el.Empty()
	for host := range m.hosts {
		condition = condition.And(noHostCondition.Format(host))
	}
	contextPath := strings.TrimSuffix(path.Path, rootPath)
	condition = condition.And(pathCondition.Format(contextPath))
	return m.storeCondition(condition)
}

func (m *Mapper) buildHostCondition(rule v1.IngressRule, path *indexedPath) *string {
	condition := hostCondition.Format(rule.Host)
	contextPath := strings.TrimSuffix(path.Path, rootPath)
	condition = condition.And(pathCondition.Format(contextPath))
	return m.storeCondition(condition)
}

func (m *Mapper) storeCondition(condition el.Expression) *string {
	m.conditions = append(m.conditions, condition.Parenthesized())
	return toPointer(condition.Closed().String())
}

func buildRoutingStep(path *indexedPath) base.FlowStep {
	return base.FlowStep{
		Name:    toPointer(routingStepName),
		Policy:  toPointer(routingPolicyName),
		Enabled: true,
		Configuration: utils.NewGenericStringMap().
			Put(routingRulesKey, buildRoutingRules(path)),
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

// This flow is used to return a 404 HTTP response when no route is found.
func (m *Mapper) buildNotFoundFlow() v2.Flow {
	flow := v2.Flow{
		Name:    toPointer(mockStepName),
		Pre:     []base.FlowStep{m.buildNotFoundStep()},
		Enabled: true,
		PathOperator: &v2.PathOperator{
			Operator: base.StartWithOperator,
			Path:     toPointer(rootPath),
		},
	}

	condition := el.Empty()

	for _, c := range m.conditions {
		condition = condition.Or(c)
	}

	flow.Condition = toPointer(condition.Parenthesized().Negated().Closed().String())

	return flow
}

func (m *Mapper) buildNotFoundStep() base.FlowStep {
	template := m.opts.Templates[http.StatusNotFound]

	return base.FlowStep{
		Name:    toPointer(mockStepName),
		Policy:  toPointer(mockPolicyName),
		Enabled: true,
		Configuration: &utils.GenericStringMap{
			Unstructured: unstructured.Unstructured{
				Object: map[string]interface{}{
					mockContentKey: template.Content,
					mockStatusKey:  fmt.Sprint(http.StatusNotFound),
					mockHeadersKey: []interface{}{
						map[string]interface{}{
							"name":  xhttp.ContentTypeHeader,
							"value": template.ContentType,
						},
					},
				},
			},
		},
	}
}

func buildProxy(ingress *v1.Ingress) *v2.Proxy {
	return &v2.Proxy{
		VirtualHosts: buildVirtualHosts(ingress),
		Groups: []*v2.EndpointGroup{
			{
				Name:      toPointer(proxyName),
				Endpoints: buildEndpoints(ingress),
			},
		},
	}
}

func buildEndpoints(ingress *v1.Ingress) []*v2.Endpoint {
	eps := make([]*v2.Endpoint, 0)
	for ruleIndex, rule := range ingress.Spec.Rules {
		for pathIndex := range rule.HTTP.Paths {
			path := &rule.HTTP.Paths[pathIndex]
			eps = append(eps, buildEndpoint(ingress, newIndexedPath(path, ruleIndex, pathIndex)))
		}
	}
	return eps
}

// For each rule and path of an ingress, build an endpoint identified by the position of the path in the rule,
// in order to be able to match it in the routing step when handling an incoming request for routing.
func buildEndpoint(ingress *v1.Ingress, path *indexedPath) *v2.Endpoint {
	return &v2.Endpoint{
		Name:   toPointer(path.String()),
		Target: buildEndpointTarget(ingress, path),
	}
}

func buildEndpointTarget(ingress *v1.Ingress, path *indexedPath) *string {
	svc := path.Backend.Service
	et := fmt.Sprintf(serviceURIPattern, svc.Name, ingress.Namespace, svc.Port.Number)
	return &et
}

// For each ingress host and path, build a virtual host.
func buildVirtualHosts(ingress *v1.Ingress) []*v2.VirtualHost {
	vhs := make([]*v2.VirtualHost, 0)
	for _, rule := range ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			vhs = append(vhs, buildVirtualHost(rule, path))
		}
	}
	return vhs
}

func buildVirtualHost(rule v1.IngressRule, path v1.HTTPIngressPath) *v2.VirtualHost {
	return &v2.VirtualHost{Host: rule.Host, Path: path.Path}
}

func toPointer(s string) *string {
	return &s
}
