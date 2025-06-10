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
	// hostHeaderWithoutPortEqualsCondition = el.Expression("#request.host.replaceAll(':(.*)$', '') eq '%s'").
	pathInfoEqualsCondition  = el.Expression("#request.pathInfo eq '%s'")
	pathInfoMatchesCondition = el.Expression("#request.pathInfo matches '%s'")
	headerEqualsCondition    = el.Expression("#request.headers['%s'][0] eq '%s'")
	headerMatchesCondition   = el.Expression("#request.headers['%s'][0] matches '%s'")
	paramEqualsCondition     = el.Expression("#request.params['%s'] eq '%s'")
	paramMatchesCondition    = el.Expression("#request.params['%s'] matches '%s'")

	routingPolicyName = "dynamic-routing"
	routingRulesKey   = "rules"
	routingPatternKey = "pattern"
	routingPattern    = `(.*)`
	routingURLKey     = "url"

	endpointMatcherPattern = "%s:{#group[0]}"
)

type weightedFlow struct {
	Flow   *v4.Flow
	Weight int
}

func buildFlows(route *gwAPIv1.HTTPRoute) []*v4.Flow {
	weightedFlows := make([]weightedFlow, 0)

	for ruleIndex, rule := range route.Spec.Rules {
		for matchIndex, match := range rule.Matches {
			conditionsExpressions := buildFlowConditionExpressions(match)
			weightedFlows = append(
				weightedFlows,
				weightedFlow{
					Flow: buildFlow(
						rule, ruleIndex, match, matchIndex, conditionsExpressions,
					),
					Weight: len(conditionsExpressions),
				},
			)
		}
	}
	return sortFlows(weightedFlows)
}

// must appear first in the list and in the same order as defined.
func sortFlows(weightedFlows []weightedFlow) []*v4.Flow {
	flows := []*v4.Flow{}
	maxWeight := getMaxWeight(weightedFlows)
	for len(flows) != len(weightedFlows) {
		for _, wf := range weightedFlows {
			if wf.Weight == maxWeight {
				flows = append(flows, wf.Flow)
			}
		}
		maxWeight--
	}
	return flows
}

func getMaxWeight(weightedFlows []weightedFlow) int {
	maxWeight := 0
	for i := range weightedFlows {
		if weightedFlows[i].Weight > maxWeight {
			maxWeight = weightedFlows[i].Weight
		}
	}
	return maxWeight
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
		Request:   buildRequestFlow(rule, ruleIndex, matchIndex),
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
			And(pathInfoEqualsCondition.Format(getExpectedPathInfo(*match.Value)))
	default:
		panic(fmt.Sprintf("unsupported path match type: %s", *match.Type))
	}
}

func buildHTTPSelector(match gwAPIv1.HTTPRouteMatch) *v4.FlowSelector {
	methods := []base.HttpMethod{}
	if match.Method != nil {
		methods = append(methods, base.HttpMethod(*match.Method))
	}
	return v4.NewHTTPSelector("/", "STARTS_WITH", methods)
}

func buildRequestFlow(rule gwAPIv1.HTTPRouteRule, ruleIndex, matchIndex int) []*v4.FlowStep {
	steps := []*v4.FlowStep{}
	if len(rule.BackendRefs) > 0 {
		steps = append(steps, buildRoutingStep(ruleIndex, matchIndex))
	}
	return append(steps, buildRequestFilters(rule)...)
}

func buildResponseFlow(rule gwAPIv1.HTTPRouteRule) []*v4.FlowStep {
	return buildResponseFilters(rule)
}

func buildRoutingStep(ruleIndex, matchIndex int) *v4.FlowStep {
	policyName := routingPolicyName
	return v4.NewFlowStep(base.FlowStep{
		Policy:  &policyName,
		Enabled: true,
		Configuration: utils.NewGenericStringMap().
			Put(routingRulesKey, buildRoutingRule(ruleIndex, matchIndex)),
	})
}

func buildRoutingRule(ruleIndex, matchIndex int) []interface{} {
	return []interface{}{
		map[string]interface{}{
			routingPatternKey: routingPattern,
			routingURLKey:     buildRoutingTarget(ruleIndex, matchIndex),
		},
	}
}

func buildRoutingTarget(ruleIndex, matchIndex int) string {
	return fmt.Sprintf(
		endpointMatcherPattern, buildEndpointGroupName(ruleIndex, matchIndex),
	)
}
