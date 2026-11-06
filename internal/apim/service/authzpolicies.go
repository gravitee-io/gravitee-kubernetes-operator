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

// authzPoliciesPath is the authz module's policies collection under the automation environment
// root.
const authzPoliciesPath = "authz/policies"

type AuthzPolicies struct {
	*client.Client
}

func NewAuthzPolicies(client *client.Client) *AuthzPolicies {
	return &AuthzPolicies{Client: client}
}

func (svc *AuthzPolicies) CreateOrUpdate(dto model.AuthzPolicyDTO) (*model.AuthzPolicyState, error) {
	return svc.createOrUpdate(dto, false)
}

func (svc *AuthzPolicies) DryRunCreateOrUpdate(dto model.AuthzPolicyDTO) (*model.AuthzPolicyState, error) {
	return svc.createOrUpdate(dto, true)
}

// createOrUpdate PUTs the payload on the collection and answers the state: ids, findings, the
// policy as the platform stores it and whether it is deployed. Every apply deploys the policy. A
// dry run answers 200 with its findings in errors.severe; a real apply refused by a finding
// answers 400. On a non-2xx the client returns the error without decoding the body.
func (svc *AuthzPolicies) createOrUpdate(dto model.AuthzPolicyDTO, dryRun bool) (*model.AuthzPolicyState, error) {
	url := svc.AutomationTarget(authzPoliciesPath).
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	state := new(model.AuthzPolicyState)
	if err := svc.HTTP.Put(url.String(), dto, state); err != nil {
		return nil, err
	}

	return state, nil
}

// Delete withdraws the policy from the PDP gateways it targets and deletes it; the entities it
// names are kept. No dry run: APIM 4.13 ignores dryRun on a delete and deletes.
func (svc *AuthzPolicies) Delete(hrid string) error {
	url := svc.AutomationTarget(authzPoliciesPath).WithPath(hrid)
	return svc.HTTP.Delete(url.String(), nil)
}

// GetByHRID reads the policy's state. Used by drift detection and tests.
func (svc *AuthzPolicies) GetByHRID(hrid string) (*model.AuthzPolicyState, error) {
	url := svc.AutomationTarget(authzPoliciesPath).WithPath(hrid)
	state := new(model.AuthzPolicyState)
	if err := svc.HTTP.Get(url.String(), state); err != nil {
		return nil, err
	}
	return state, nil
}
