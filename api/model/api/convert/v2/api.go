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

package v2

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
)

func FromV4(v4Api *v4.Api) *v2.Api {
	return &v2.Api{
		ApiBase:           v4Api.ApiBase,
		Version:           v4Api.ApiVersion,
		DefinitionVersion: base.DefinitionVersionV2,
		FlowMode:          v2.DefaultFlowMode,
		Proxy:             toProxy(v4Api.Listeners, v4Api.EndpointGroups, v4Api.Analytics),
		Plans:             toPlans(v4Api.Plans),
		Flows:             toFlows(v4Api.Flows),
		Services:          toApiServices(v4Api.Services),
	}
}
