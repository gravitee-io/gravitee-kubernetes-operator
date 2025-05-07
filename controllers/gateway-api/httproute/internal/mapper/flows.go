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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/el"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	contextPathEqualsCondition = el.Expression("#request.contextPath eq '%s'")
	pathInfoEqualsCondition    = el.Expression("#request.pathInfo eq '%s'")
	pathInfoMatchesCondition   = el.Expression("#request.pathInfo matches '%s'")
	headerEqualsCondition      = el.Expression("#request.headers['%s'][0] eq '%s'")
	headerMatchesCondition     = el.Expression("#request.headers['%s'][0] matches '%s'")
	paramEqualsCondition       = el.Expression("#request.params['%s'] eq '%s'")
	paramMatchesCondition      = el.Expression("#request.params['%s'] matches '%s'")

	routingPolicyName = "dynamic-routing"
	routingRulesKey   = "rules"
	routingPatternKey = "pattern"
	routingPattern    = "(.*)"
	routingURLKey     = "url"

	endpointMatcherPattern = "%s:{#group[0]}"
)

func buildFlows(route *gwAPIv1.HTTPRoute) []*v4.Flow {
	flows := make([]*v4.Flow, 0)
	for ruleIndex, rule := range route.Spec.Rules {
		for matchIndex, match := range rule.Matches {
			conditionsExpressions := buildFlowConditionExpressions(match)
			flows = append(
				flows,
				buildFlow(
					rule, ruleIndex, match, matchIndex, conditionsExpressions,
				),
			)
		}
	}
	return flows
}

func buildFlow(
	rule gwAPIv1.HTTPRouteRule,
	ruleIndex int,
	match gwAPIv1.HTTPRouteMatch,
	matchIndex int,
	conditionsExpressions []el.Expression,
) *v4.Flow {
	flowName := fmt.Sprintf("rule-%d-match%d", ruleIndex, matchIndex)
	return &v4.Flow{
		Name:      &flowName,
		Request:   buildRequestFlow(rule, ruleIndex),
		Response:  buildResponseFlow(rule),
		Enabled:   true,
		Selectors: buildFlowSelectors(match, conditionsExpressions),
	}
}

func buildFlowSelectors(
	match gwAPIv1.HTTPRouteMatch,
	conditionsExpressions []el.Expression,
) []*v4.FlowSelector {
	return []*v4.FlowSelector{
		buildHTTPSelector(match),
		v4.NewConditionSelector(buildCondition(conditionsExpressions)),
	}
}

func buildCondition(conditionsExpressions []el.Expression) string {
	el := el.Empty()
	for _, exp := range conditionsExpressions {
		el = el.And(exp)
	}
	return el.Closed().String()
}

func buildFlowConditionExpressions(match gwAPIv1.HTTPRouteMatch) []el.Expression {
	expressions := []el.Expression{buildPathCondition(match.Path)}
	for _, header := range match.Headers {
		expressions = append(expressions, buildHeaderCondition(header))
	}
	for _, param := range match.QueryParams {
		expressions = append(expressions, buildParamCondition(param))
	}
	return expressions
}

func buildParamCondition(match gwAPIv1.HTTPQueryParamMatch) el.Expression {
	if *match.Type == gwAPIv1.QueryParamMatchExact {
		return paramEqualsCondition.Format(match.Name, match.Value)
	}
	return paramMatchesCondition.Format(match.Name, match.Value)
}

func buildHeaderCondition(match gwAPIv1.HTTPHeaderMatch) el.Expression {
	if *match.Type == gwAPIv1.HeaderMatchExact {
		return headerEqualsCondition.Format(match.Name, match.Value)
	}
	return headerMatchesCondition.Format(match.Name, match.Value)
}

func buildPathCondition(match *gwAPIv1.HTTPPathMatch) el.Expression {
	switch *match.Type {
	case gwAPIv1.PathMatchPathPrefix:
		return contextPathEqualsCondition.Format(addTrailingSlash(*match.Value))
	case gwAPIv1.PathMatchRegularExpression:
		return contextPathEqualsCondition.Format(rootPath).
			And(pathInfoMatchesCondition.Format(*match.Value))
	case gwAPIv1.PathMatchExact:
		return contextPathEqualsCondition.Format(addTrailingSlash(*match.Value)).
			And(pathInfoEqualsCondition.Format(rootPath))
	default:
		panic(fmt.Sprintf("unsupported path match type: %s", *match.Type))
	}
}

func addTrailingSlash(s string) string {
	if s == rootPath {
		return s
	}
	return s + "/"
}

func buildHTTPSelector(match gwAPIv1.HTTPRouteMatch) *v4.FlowSelector {
	methods := []base.HttpMethod{}
	if match.Method != nil {
		methods = append(methods, base.HttpMethod(*match.Method))
	}
	return v4.NewHTTPSelector("/", "START_WITH", methods)
}

func buildRequestFlow(rule gwAPIv1.HTTPRouteRule, ruleIndex int) []*v4.FlowStep {
	return append(
		[]*v4.FlowStep{buildRoutingStep(ruleIndex)},
		buildRequestFilters(rule)...,
	)
}

func buildResponseFlow(rule gwAPIv1.HTTPRouteRule) []*v4.FlowStep {
	return buildResponseFilters(rule)
}

func buildRoutingStep(ruleIndex int) *v4.FlowStep {
	policyName := routingPolicyName
	return v4.NewFlowStep(base.FlowStep{
		Policy:  &policyName,
		Enabled: true,
		Configuration: utils.NewGenericStringMap().
			Put(routingRulesKey, buildRoutingRule(ruleIndex)),
	})
}

func buildRoutingRule(index int) []interface{} {
	return []interface{}{
		map[string]interface{}{
			routingPatternKey: routingPattern,
			routingURLKey:     buildRoutingTarget(index),
		},
	}
}

func buildRoutingTarget(index int) string {
	return fmt.Sprintf(endpointMatcherPattern, buildEndpointGroupName(index))
}
