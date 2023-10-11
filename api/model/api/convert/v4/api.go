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

package v4

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
)

func FromV2(api *v2.Api) *v4.Api {
	return &v4.Api{
		ApiBase:           api.ApiBase,
		ApiVersion:        api.Version,
		DefinitionVersion: base.DefinitionVersionV4,
		Type:              v4.ProxyType,
		Listeners:         toListeners(api.Proxy),
		EndpointGroups:    toEndpointGroups(api.Proxy),
		Plans:             toV4Plans(api.Plans),
		FlowExecution:     toFlowExecution(api.FlowMode),
		Flows:             toFlows(api.Flows),
		Analytics:         toAnalytics(api.Proxy),
		Services:          toApiServices(api.Services),
	}
}
