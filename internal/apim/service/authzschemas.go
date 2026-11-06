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

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
)

// authzSchemasPath is the authz module's schema documents collection under the automation
// environment root.
const authzSchemasPath = "authz/schemas"

type AuthzSchemas struct {
	*client.Client
}

func NewAuthzSchemas(client *client.Client) *AuthzSchemas {
	return &AuthzSchemas{Client: client}
}

func (svc *AuthzSchemas) CreateOrUpdate(dto model.AuthzSchemaDTO) (*model.AuthzSchemaState, error) {
	return svc.createOrUpdate(dto, false)
}

func (svc *AuthzSchemas) DryRunCreateOrUpdate(dto model.AuthzSchemaDTO) (*model.AuthzSchemaState, error) {
	return svc.createOrUpdate(dto, true)
}

// createOrUpdate PUTs the payload on the collection and answers the state: ids, findings and the
// target the document is published to. A dry run answers 200 with its findings in
// errors.severe; a real apply refused by a finding answers 400. On a non-2xx the client returns
// the error without decoding the body.
func (svc *AuthzSchemas) createOrUpdate(dto model.AuthzSchemaDTO, dryRun bool) (*model.AuthzSchemaState, error) {
	url := svc.AutomationTarget(authzSchemasPath).
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	state := new(model.AuthzSchemaState)
	if err := svc.HTTP.Put(url.String(), dto, state); err != nil {
		return nil, err
	}

	return state, nil
}

// Delete withdraws the schema document from the PDP gateway it targets and deletes it.
func (svc *AuthzSchemas) Delete(hrid string) error {
	url := svc.AutomationTarget(authzSchemasPath).WithPath(hrid)
	return svc.HTTP.Delete(url.String(), nil)
}

// GetByHRID reads the schema document's state. Used by drift detection and tests.
func (svc *AuthzSchemas) GetByHRID(hrid string) (*model.AuthzSchemaState, error) {
	url := svc.AutomationTarget(authzSchemasPath).WithPath(hrid)
	state := new(model.AuthzSchemaState)
	if err := svc.HTTP.Get(url.String(), state); err != nil {
		return nil, err
	}
	return state, nil
}
