// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

import (
	"strconv"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"
)

// mcpProxiesPath is the aim module's MCP proxies collection under the
// automation environment root.
const mcpProxiesPath = "aim/mcp-proxies"

type McpProxies struct {
	*client.Client
}

func NewMcpProxies(client *client.Client) *McpProxies {
	return &McpProxies{Client: client}
}

func (svc *McpProxies) CreateOrUpdate(proxy *v1alpha1.McpProxy) (*model.McpProxyState, error) {
	return svc.createOrUpdate(proxy, false)
}

func (svc *McpProxies) DryRunCreateOrUpdate(
	proxy *v1alpha1.McpProxy,
) (*model.McpProxyState, error) {
	return svc.createOrUpdate(proxy, true)
}

// createOrUpdate PUTs the spec on the collection and answers the state: ids, findings, the
// observed lifecycle and, for a Studio, the tools with their entity ids. A dry run answers 200
// with its findings in errors.severe; a real apply refused by a finding answers 400. On a non-2xx
// the client returns the error without decoding the body.
func (svc *McpProxies) createOrUpdate(
	proxy *v1alpha1.McpProxy,
	dryRun bool,
) (*model.McpProxyState, error) {
	url := svc.AutomationTarget(mcpProxiesPath).
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	dto := model.ToMcpProxyDTO(proxy)
	state := new(model.McpProxyState)

	if err := svc.HTTP.Put(url.String(), dto, state); err != nil {
		return nil, err
	}

	k8s.AddAutomationAPIManagedCondition(proxy)

	return state, nil
}

func (svc *McpProxies) Delete(proxy *v1alpha1.McpProxy) error {
	hrid := refs.NewNamespacedNameFromObject(proxy).HRID()
	url := svc.AutomationTarget(mcpProxiesPath).WithPath(hrid)
	return svc.HTTP.Delete(url.String(), nil)
}

// GetByHRID reads the proxy's state. Used by drift detection and tests.
func (svc *McpProxies) GetByHRID(hrid string) (*model.McpProxyState, error) {
	url := svc.AutomationTarget(mcpProxiesPath).WithPath(hrid)
	state := new(model.McpProxyState)
	if err := svc.HTTP.Get(url.String(), state); err != nil {
		return nil, err
	}
	return state, nil
}
