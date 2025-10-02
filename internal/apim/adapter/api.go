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

package adapter

import (
	"encoding/json"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"
	"k8s.io/utils/ptr"
)

type ApiV2 struct {
	v2.Api `json:",inline"`
}

func NewApiV2(a v2.Api) *ApiV2 {
	api := &ApiV2{Api: a}
	// Set default values
	if api.Flows == nil {
		api.Flows = []v2.Flow{}
	}

	api.populateDefaultFlowValues()

	if api.Plans == nil {
		api.Plans = []*v2.Plan{}
	}

	api.populateDefaultPlanValues()

	for _, p := range api.Plans {
		if p.Tags == nil {
			p.Tags = []string{}
		}
	}
	if api.Proxy != nil {
		api.populateDefaultApiProxyValues()
	}

	if api.Services != nil {
		api.populateDefaultApiServiceValues()
	}

	for _, r := range api.Resources {
		if r.Enabled == nil {
			r.Enabled = ptr.To(true)
		}
	}

	return api
}

func (api *ApiV2) MarshalJSON() ([]byte, error) {
	return json.Marshal(api.Api)
}

func populateDefaultHttpClientSslOptions(o *base.HttpClientSslOptions) {
	if o.TrustAll == nil {
		o.TrustAll = ptr.To(false)
	}
	if o.HostnameVerifier == nil {
		o.HostnameVerifier = ptr.To(true)
	}
}

func populateDefaultHttpClientOptions(o *base.HttpClientOptions) {
	if o.KeepAlive == nil {
		o.KeepAlive = ptr.To(true)
	}
	if o.KeepAliveTimeout == nil {
		o.KeepAliveTimeout = ptr.To(uint64(30000))
	}
	if o.Pipelining == nil {
		o.Pipelining = ptr.To(false)
	}
	if o.UseCompression == nil {
		o.UseCompression = ptr.To(false)
	}
	if o.PropagateClientAcceptEncoding == nil {
		o.PropagateClientAcceptEncoding = ptr.To(false)
	}
	if o.FollowRedirects == nil {
		o.FollowRedirects = ptr.To(false)
	}
	if o.ClearTextUpgrade == nil {
		o.ClearTextUpgrade = ptr.To(true)
	}
	if o.ProtocolVersion == nil {
		o.ProtocolVersion = ptr.To(base.ProtocolVersion("HTTP_1_1"))
	}
}

func (api *ApiV2) populateDefaultApiServiceValues() {
	if api.Services.EndpointDiscoveryService != nil {
		if api.Services.EndpointDiscoveryService.Tenants == nil {
			api.Services.EndpointDiscoveryService.Tenants = []string{}
		}
	}
	hcs := api.Services.HealthCheckService
	if hcs != nil {
		if hcs.Enabled == nil {
			hcs.Enabled = ptr.To(false)
		}

		if hcs.Steps == nil {
			hcs.Steps = []*v2.HealthCheckStep{}
		}

		for _, s := range hcs.Steps {
			if s.Request.Headers == nil {
				s.Request.Headers = []base.HttpHeader{}
			}
		}
	}
}

func (api *ApiV2) populateDefaultApiProxyValues() { //nolint:gocognit // normal complexity
	if api.Proxy.Groups == nil {
		api.Proxy.Groups = []*v2.EndpointGroup{}
	}

	for _, g := range api.Proxy.Groups {
		if g.HttpClientOptions != nil {
			populateDefaultHttpClientOptions(g.HttpClientOptions)
		}
		if g.HttpClientSslOptions != nil {
			populateDefaultHttpClientSslOptions(g.HttpClientSslOptions)
		}
		if g.HttpProxy != nil {
			if g.HttpProxy.Enabled == nil {
				g.HttpProxy.Enabled = ptr.To(false)
			}
			if g.HttpProxy.UseSystemProxy == nil {
				g.HttpProxy.UseSystemProxy = ptr.To(false)
			}
		}
		for _, e := range g.Endpoints {
			if e.Tenants == nil {
				e.Tenants = []string{}
			}
			if e.Headers == nil {
				e.Headers = []base.HttpHeader{}
			}
			if e.HttpClientOptions != nil {
				populateDefaultHttpClientOptions(e.HttpClientOptions)
			}
			if e.HttpClientSslOptions != nil {
				populateDefaultHttpClientSslOptions(e.HttpClientSslOptions)
			}
			if e.HttpProxy != nil {
				if e.HttpProxy.Enabled == nil {
					e.HttpProxy.Enabled = ptr.To(false)
				}
				if e.HttpProxy.UseSystemProxy == nil {
					e.HttpProxy.UseSystemProxy = ptr.To(false)
				}
			}
		}
		if g.Services != nil {
			if g.Services.EndpointDiscoveryService != nil && g.Services.EndpointDiscoveryService.Tenants == nil {
				g.Services.EndpointDiscoveryService.Tenants = []string{}
			}
			if g.Services.HealthCheckService != nil && g.Services.HealthCheckService.Steps == nil {
				g.Services.HealthCheckService.Steps = []*v2.HealthCheckStep{}
			}
		}
	}
}

func (api *ApiV2) populateDefaultPlanValues() { //nolint:gocognit // normal complexity
	for _, p := range api.Plans {
		if p.Validation == nil {
			p.Validation = ptr.To(base.PlanValidation("AUTO"))
		}
		if p.Status == nil {
			p.Status = ptr.To(base.PlanStatus("PUBLISHED"))
		}
		if p.Type == nil {
			p.Type = ptr.To(base.PlanType("API"))
		}
		if p.Flows == nil {
			p.Flows = []v2.Flow{}
		}

		for _, f := range p.Flows {
			if f.Pre == nil {
				f.Pre = []base.FlowStep{}
			}
			if f.Post == nil {
				f.Post = []base.FlowStep{}
			}
			if f.Enabled == nil {
				f.Enabled = ptr.To(true)
			}
			if f.Methods == nil {
				f.Methods = []base.HttpMethod{}
			}
			if f.Consumers == nil {
				f.Consumers = []v2.Consumer{}
			}
		}
		if p.ExcludedGroups == nil {
			p.ExcludedGroups = []string{}
		}
	}
}

func (api *ApiV2) populateDefaultFlowValues() {
	for _, f := range api.Flows {
		if f.Pre == nil {
			f.Pre = []base.FlowStep{}
		}
		if f.Post == nil {
			f.Post = []base.FlowStep{}
		}
		if f.Enabled == nil {
			f.Enabled = ptr.To(true)
		}
		if f.Methods == nil {
			f.Methods = []base.HttpMethod{}
		}
		if f.Consumers == nil {
			f.Consumers = []v2.Consumer{}
		}
	}
}

type ApiV4 struct {
	v4.Api `json:",inline"`
}

func NewApiV4(a v4.Api) *ApiV4 {
	api := &ApiV4{Api: a}
	for _, eg := range api.EndpointGroups {
		api.populateDefaultApiEndpointGroupValues(eg)
	}

	if api.Flows == nil {
		api.Flows = []*v4.Flow{}
	}

	populateDefaultFlows(api.Flows)

	if api.Plans == nil {
		api.Plans = &map[string]*v4.Plan{}
	}

	api.populateDefaultApiPlanValues()

	for _, r := range api.Resources {
		if r.Enabled == nil {
			r.Enabled = ptr.To(true)
		}
	}

	return api
}

func (api *ApiV4) MarshalJSON() ([]byte, error) {
	return json.Marshal(api.Api)
}

func (api *ApiV4) populateDefaultApiPlanValues() {
	for _, p := range *api.Plans {
		if p.DefinitionVersion == nil {
			p.DefinitionVersion = ptr.To(v4.PlanDefinitionVersion)
		}
		if p.Validation == nil {
			p.Validation = ptr.To(base.PlanValidation("AUTO"))
		}
		if p.Status == nil {
			p.Status = ptr.To(base.PlanStatus("PUBLISHED"))
		}
		if p.Type == nil {
			p.Type = ptr.To(base.PlanType("API"))
		}
		if p.Flows == nil {
			p.Flows = []*v4.Flow{}
		} else {
			populateDefaultFlows(p.Flows)
		}
		if p.ExcludedGroups == nil {
			p.ExcludedGroups = []string{}
		}
	}
}

func populateDefaultFlows(flows []*v4.Flow) { //nolint:gocognit // normal complexity
	for _, f := range flows {
		if f.Enabled == nil {
			f.Enabled = ptr.To(true)
		}
		for _, r := range f.Request {
			if r.Enabled == nil {
				r.Enabled = ptr.To(true)
			}
		}
		for _, r := range f.Response {
			if r.Enabled == nil {
				r.Enabled = ptr.To(true)
			}
		}
		for _, r := range f.Subscribe {
			if r.Enabled == nil {
				r.Enabled = ptr.To(true)
			}
		}
		for _, r := range f.Publish {
			if r.Enabled == nil {
				r.Enabled = ptr.To(true)
			}
		}
		for _, r := range f.Connect {
			if r.Enabled == nil {
				r.Enabled = ptr.To(true)
			}
		}
		for _, r := range f.Interact {
			if r.Enabled == nil {
				r.Enabled = ptr.To(true)
			}
		}
	}
}

func (api *ApiV4) populateDefaultApiEndpointGroupValues(eg *v4.EndpointGroup) {
	if eg.Endpoints == nil {
		eg.Endpoints = []*v4.Endpoint{}
	}
	for _, ep := range eg.Endpoints {
		if ep.Tenants == nil {
			ep.Tenants = []string{}
		}
	}

	if eg.HttpClientOptions != nil {
		populateDefaultHttpClientOptions(eg.HttpClientOptions)
	}
	if eg.HttpClientSslOptions != nil {
		populateDefaultHttpClientSslOptions(eg.HttpClientSslOptions)
	}
}
