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

// authzPdpsPath is the authz module's PDP gateways collection under the automation
// environment root.
const authzPdpsPath = "authz/pdps"

type AuthzPdps struct {
	*client.Client
}

func NewAuthzPdps(client *client.Client) *AuthzPdps {
	return &AuthzPdps{Client: client}
}

func (svc *AuthzPdps) CreateOrUpdate(pdp *v1alpha1.AuthzPdp) (*model.AuthzPdpState, error) {
	return svc.createOrUpdate(pdp, false)
}

func (svc *AuthzPdps) DryRunCreateOrUpdate(pdp *v1alpha1.AuthzPdp) (*model.AuthzPdpState, error) {
	return svc.createOrUpdate(pdp, true)
}

// createOrUpdate PUTs the spec on the collection and answers the state: ids, findings, the
// routing scope and the AuthZEN discovery path. A dry run answers 200 with its findings in
// errors.severe; a real apply refused by a finding answers 400. On a non-2xx the client returns
// the error without decoding the body.
func (svc *AuthzPdps) createOrUpdate(pdp *v1alpha1.AuthzPdp, dryRun bool) (*model.AuthzPdpState, error) {
	url := svc.AutomationTarget(authzPdpsPath).
		WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	dto := model.ToAuthzPdpDTO(pdp)
	state := new(model.AuthzPdpState)

	if err := svc.HTTP.Put(url.String(), dto, state); err != nil {
		return nil, err
	}

	k8s.AddAutomationAPIManagedCondition(pdp)

	return state, nil
}

// Delete removes the PDP gateway and its endpoint. The platform answers 409, deleting nothing,
// when the delete would leave a policy or an entity with no target.
func (svc *AuthzPdps) Delete(pdp *v1alpha1.AuthzPdp) error {
	hrid := refs.NewNamespacedNameFromObject(pdp).HRID()
	url := svc.AutomationTarget(authzPdpsPath).WithPath(hrid)
	return svc.HTTP.Delete(url.String(), nil)
}

// GetByHRID reads the PDP gateway's state. Used by drift detection and tests.
func (svc *AuthzPdps) GetByHRID(hrid string) (*model.AuthzPdpState, error) {
	url := svc.AutomationTarget(authzPdpsPath).WithPath(hrid)
	state := new(model.AuthzPdpState)
	if err := svc.HTTP.Get(url.String(), state); err != nil {
		return nil, err
	}
	return state, nil
}
