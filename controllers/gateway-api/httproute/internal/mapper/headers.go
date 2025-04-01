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
	"strings"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/utils"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var (
	transformHeadersPolicyName = "transform-headers"
)

const (
	removeHeadersKey      = "removeHeaders"
	addHeadersKey         = "addHeaders"
	appendHeadersKey      = "appendHeaders"
	headerAppenderPattern = "{#request.headers['%s'] == null ? '%s' : #request.headers['%s'] +','+'%s'}"
)

func buildHeaderTransformer(filter gwAPIv1.HTTPHeaderFilter) *v4.FlowStep {
	return v4.NewFlowStep(base.FlowStep{
		Policy:        &transformHeadersPolicyName,
		Enabled:       true,
		Configuration: buildHeaderTransformerConfig(filter),
	})
}

func buildHeaderTransformerConfig(filter gwAPIv1.HTTPHeaderFilter) *utils.GenericStringMap {
	config := utils.NewGenericStringMap()
	config.Put(addHeadersKey, mapSetHeader(filter.Set))
	config.Put(appendHeadersKey, mapAddHeader(filter.Add))
	config.Put(removeHeadersKey, filter.Remove)
	return config
}

func mapSetHeader(headers []gwAPIv1.HTTPHeader) []map[string]any {
	set := make([]map[string]any, len(headers))
	for i := range headers {
		set[i] = map[string]any{
			"name":  string(headers[i].Name),
			"value": headers[i].Value,
		}
	}
	return set
}

func mapAddHeader(headers []gwAPIv1.HTTPHeader) []map[string]any {
	add := []map[string]any{}
	for _, header := range headers {
		values := strings.Split(header.Value, ",")
		for _, val := range values {
			add = append(
				add,
				map[string]any{
					"name":  string(header.Name),
					"value": val,
				},
			)
		}
	}
	return add
}
