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
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var (
	httpRedirectPolicyName = "http-redirect"

	rulesKey        = "rules"
	rulePathKey     = "path"
	ruleLocationKey = "location"
	ruleStatusKey   = "status"

	rulePath = "/(.*)"

	locationPathWhenNoGivenPath   = "{#request.contextPath}/{#group[0]}"
	locationPathWhenReplacePrefix = "/%s/{#group[0]}"

	locationSchemeDefault = "{#request.scheme}"
	locationHostDefault   = "{#request.host}"

	statusCodeDefault = 302
)

func buildHTTPRedirect(filter gwAPIv1.HTTPRequestRedirectFilter) *v4.FlowStep {
	return v4.NewFlowStep(base.FlowStep{
		Policy:        &httpRedirectPolicyName,
		Enabled:       true,
		Configuration: buildHTTPRedirectConfig(filter),
	})
}

func buildHTTPRedirectConfig(filter gwAPIv1.HTTPRequestRedirectFilter) *utils.GenericStringMap {
	config := utils.NewGenericStringMap()
	rules := []any{buildRedirectRule(filter)}
	config.Put(rulesKey, rules)
	return config
}

func buildRedirectRule(filter gwAPIv1.HTTPRequestRedirectFilter) map[string]any {
	rule := make(map[string]any)
	rule[rulePathKey] = rulePath
	rule[ruleLocationKey] = buildRedirectLocation(filter)
	rule[ruleStatusKey] = getStatusCode(filter)
	return rule
}

func buildRedirectLocation(filter gwAPIv1.HTTPRequestRedirectFilter) string {
	scheme := getLocationScheme(filter)
	host := getLocationHost(filter)
	path := getLocationPath(filter)
	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}

func getLocationScheme(filter gwAPIv1.HTTPRequestRedirectFilter) string {
	if filter.Scheme != nil {
		return *filter.Scheme
	}
	return locationSchemeDefault
}

func getLocationHost(filter gwAPIv1.HTTPRequestRedirectFilter) string {
	if filter.Hostname != nil {
		host := string(*filter.Hostname)
		if filter.Port != nil {
			return fmt.Sprintf("%s:%d", host, int32(*filter.Port))
		}
		return host
	}
	return locationHostDefault
}

func getLocationPath(filter gwAPIv1.HTTPRequestRedirectFilter) string {
	if filter.Path == nil {
		return locationPathWhenNoGivenPath
	}
	if filter.Path.ReplaceFullPath != nil {
		return *filter.Path.ReplaceFullPath
	}
	return fmt.Sprintf(locationPathWhenReplacePrefix, *filter.Path.ReplacePrefixMatch)
}

func getStatusCode(filter gwAPIv1.HTTPRequestRedirectFilter) int {
	if filter.StatusCode != nil {
		return *filter.StatusCode
	}
	return statusCodeDefault
}
