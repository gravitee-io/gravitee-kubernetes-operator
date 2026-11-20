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

package service

import (
	"strconv"

	sdk "github.com/gravitee-io/gravitee-automation-sdk/apim-sdk/pkg/sdk"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/env"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
	httputil "github.com/gravitee-io/gravitee-kubernetes-operator/internal/http"
)

const (
	crossIDParam     = "crossId"
	stateActionParam = "action"
	planParam        = "plan"
	applicationParam = "application"
)

var deleteParams = map[string]string{
	"closePlans": "true",
}

// APIs brings support for managing gravitee.io APIM APIs.
type APIs struct {
	*client.Client
}

func NewAPIs(client *client.Client) *APIs {
	return &APIs{Client: client}
}

// GetByCrossID for test purposes only.
func (svc *APIs) GetByCrossID(crossID string) (*model.ApiListItem, error) {
	url := svc.EnvV1Target("apis").WithQueryParam(crossIDParam, crossID)
	apis := new([]model.ApiListItem)

	if err := svc.HTTP.Get(url.String(), apis); err != nil {
		return nil, err
	}

	if len(*apis) == 0 {
		return nil, errors.NewNotFoundError()
	}

	return &(*apis)[0], nil
}

// GetByID For tests purposes only.
func (svc *APIs) GetByID(apiID string) (*model.ApiEntity, error) {
	url := svc.EnvV1Target("apis").WithPath(apiID)
	api := new(model.ApiEntity)

	if err := svc.HTTP.Get(url.String(), api); err != nil {
		return nil, err
	}

	return api, nil
}

// GetV4ByHRID fetches a V4 API from the Automation API by HRID. For test purposes only.
func (svc *APIs) GetV4ByHRID(hrid string) (*model.APIV4DTO, error) {
	return svc.getV4(svc.AutomationTarget("apis").WithPath(hrid))
}

// GetV4WithUUID fetches a pre-Automation V4 API from the Automation API by UUID.
// For test purposes only.
func (svc *APIs) GetV4WithUUID(apiID string) (*model.APIV4DTO, error) {
	url := svc.AutomationTarget("apis").WithPath(apiID).
		WithQueryParam("hridContainsUUID", strconv.FormatBool(true))
	return svc.getV4(url)
}

func (svc *APIs) getV4(url *httputil.URL) (*model.APIV4DTO, error) {
	resp := new(model.APIV4DTO)
	if err := svc.HTTP.Get(url.String(), resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (svc *APIs) ImportV2(spec *v2.Api) (*base.Status, error) {
	return svc.importV2(spec, false)
}

func (svc *APIs) DryRunImportV2(spec *v2.Api) (*base.Status, error) {
	return svc.importV2(spec, true)
}

func (svc *APIs) importV2(spec *v2.Api, dryRun bool) (*base.Status, error) {
	url := svc.EnvV1Target("apis/import-crd").WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	status := new(base.Status)
	if err := svc.HTTP.Put(url.String(), spec, status); err != nil {
		return nil, err
	}

	return status, nil
}

// UpsertV4Legacy PUTs a v4 API created before the Automation API (pre-4.12), whose UUID is sent as the hrid, with
// hridContainsUUID=true. It is the only v4 write left on GKO's HTTP client: APIM reads hridContainsUUID on
// PUT /apis but its Automation OAS does not declare it, so the generated SDK cannot send it (APIM-16599, Issue 01).
// Delete it once APIM declares the parameter and the SDK syncs it. Every other v4 call goes through
// internal/apim/apiv4.
func (svc *APIs) UpsertV4Legacy(payload sdk.ApiV4, dryRun, ignoreGroups bool) (*sdk.ApiV4State, error) {
	url := svc.AutomationTarget("apis").
		WithQueryParam("dryRun", strconv.FormatBool(dryRun)).
		WithQueryParam("hridContainsUUID", strconv.FormatBool(true))
	if ignoreGroups {
		url = url.WithQueryParam("ignoreGroups", env.TrueString)
	}

	state := new(sdk.ApiV4State)
	if err := svc.HTTP.Put(url.String(), payload, state); err != nil {
		return nil, err
	}
	return state, nil
}

func (svc *APIs) UpdateState(apiID string, action model.Action) error {
	url := svc.EnvV1Target("apis").WithPath(apiID).WithQueryParam(stateActionParam, string(action))
	return svc.HTTP.Post(url.String(), nil, nil)
}

func (svc *APIs) DeleteV2(apiID string) error {
	url := svc.EnvV1Target("apis").WithPath(apiID).WithQueryParams(deleteParams)
	return svc.HTTP.Delete(url.String(), nil)
}

func (svc *APIs) SetKubernetesContext(apiID string) error {
	url := svc.EnvV1Target("apis").WithPath(apiID).WithPath("definition-context")
	return svc.HTTP.Put(url.String(), model.NewKubernetesContext(), nil)
}

func (svc *APIs) Deploy(id string) error {
	url := svc.EnvV1Target("apis").WithPath(id).WithPath("deploy")
	return svc.HTTP.Post(url.String(), new(model.ApiDeployment), nil)
}
