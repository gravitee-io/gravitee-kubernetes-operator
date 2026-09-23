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

// catalogMcpServersPath is the aim module's catalog MCP servers collection under the
// automation environment root.
const catalogMcpServersPath = "aim/catalog/mcp-servers"

type CatalogMcpServers struct {
	*client.Client
}

func NewCatalogMcpServers(client *client.Client) *CatalogMcpServers {
	return &CatalogMcpServers{Client: client}
}

func (svc *CatalogMcpServers) CreateOrUpdate(srv *v1alpha1.CatalogMcpServer) (*model.CatalogMcpServerState, error) {
	return svc.createOrUpdate(srv, false)
}

func (svc *CatalogMcpServers) DryRunCreateOrUpdate(
	srv *v1alpha1.CatalogMcpServer,
) (*model.CatalogMcpServerState, error) {
	return svc.createOrUpdate(srv, true)
}

// createOrUpdate PUTs the spec on the collection. The platform discovers the upstream before
// persisting anything (dry run included) and answers the state: ids, findings and what it
// discovered. On a non-2xx the client returns the error without decoding the body.
func (svc *CatalogMcpServers) createOrUpdate(
	srv *v1alpha1.CatalogMcpServer,
	dryRun bool,
) (*model.CatalogMcpServerState, error) {
	url := svc.AutomationTarget(catalogMcpServersPath).
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	dto := model.ToCatalogMcpServerDTO(srv)
	state := new(model.CatalogMcpServerState)

	if err := svc.HTTP.Put(url.String(), dto, state); err != nil {
		return nil, err
	}

	k8s.AddAutomationAPIManagedCondition(srv)

	return state, nil
}

func (svc *CatalogMcpServers) Delete(srv *v1alpha1.CatalogMcpServer) error {
	hrid := refs.NewNamespacedNameFromObject(srv).HRID()
	url := svc.AutomationTarget(catalogMcpServersPath).WithPath(hrid)
	return svc.HTTP.Delete(url.String(), nil)
}

// GetByHRID reads the server's state. Used by drift detection and tests.
func (svc *CatalogMcpServers) GetByHRID(hrid string) (*model.CatalogMcpServerState, error) {
	url := svc.AutomationTarget(catalogMcpServersPath).WithPath(hrid)
	state := new(model.CatalogMcpServerState)
	if err := svc.HTTP.Get(url.String(), state); err != nil {
		return nil, err
	}
	return state, nil
}
