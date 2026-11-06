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

// authzEntitiesPath is the authz module's entities collection under the automation environment
// root.
const authzEntitiesPath = "authz/entities"

type AuthzEntities struct {
	*client.Client
}

func NewAuthzEntities(client *client.Client) *AuthzEntities {
	return &AuthzEntities{Client: client}
}

func (svc *AuthzEntities) CreateOrUpdate(dto model.AuthzEntityDTO) (*model.AuthzEntityState, error) {
	return svc.createOrUpdate(dto, false)
}

func (svc *AuthzEntities) DryRunCreateOrUpdate(dto model.AuthzEntityDTO) (*model.AuthzEntityState, error) {
	return svc.createOrUpdate(dto, true)
}

// createOrUpdate PUTs the payload on the collection and answers the state: ids, findings and the
// entity as the platform stores it. A dry run answers 200 with its findings in errors.severe; a
// real apply refused by a finding answers 400. On a non-2xx the client returns
// the error without decoding the body.
func (svc *AuthzEntities) createOrUpdate(dto model.AuthzEntityDTO, dryRun bool) (*model.AuthzEntityState, error) {
	url := svc.AutomationTarget(authzEntitiesPath).
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	state := new(model.AuthzEntityState)
	if err := svc.HTTP.Put(url.String(), dto, state); err != nil {
		return nil, err
	}

	return state, nil
}

// Delete deletes the entity, its RESOURCE children with it, and detaches its PRINCIPAL children.
// Over the platform's cascade limit, it answers 413 and deletes nothing. No dry run: APIM 4.13
// ignores dryRun on a delete and deletes.
func (svc *AuthzEntities) Delete(hrid string) error {
	url := svc.AutomationTarget(authzEntitiesPath).WithPath(hrid)
	return svc.HTTP.Delete(url.String(), nil)
}

// GetByHRID reads the entity's state. Used by drift detection and tests.
func (svc *AuthzEntities) GetByHRID(hrid string) (*model.AuthzEntityState, error) {
	url := svc.AutomationTarget(authzEntitiesPath).WithPath(hrid)
	state := new(model.AuthzEntityState)
	if err := svc.HTTP.Get(url.String(), state); err != nil {
		return nil, err
	}
	return state, nil
}
