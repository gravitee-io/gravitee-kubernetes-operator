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

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/refs"
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/core"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/k8s"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"

	v4 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v4"

	v2 "github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/v2"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/model"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/errors"
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

func (svc *APIs) GetV4ByID(apiID string) (*v4.Api, error) {
	url := svc.EnvV2Target("apis").WithPath(apiID)
	resp := new(v4.Api)

	if err := svc.HTTP.Get(url.String(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// GetV4ByHRID fetches a V4 API from the Automation API by HRID. For test purposes only.
func (svc *APIs) GetV4ByHRID(hrid string) (*v4.AutomationApi, error) {
	url := svc.AutomationTarget("apis").WithPath(hrid)

	resp := new(v4.AutomationApi)

	if err := svc.HTTP.Get(url.String(), &resp); err != nil {
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

func (svc *APIs) ImportV4(api *v1alpha1.ApiV4Definition) (*base.Status, error) {
	return svc.applyV4(api, false)
}

func (svc *APIs) DryRunImportV4(api *v1alpha1.ApiV4Definition) (*base.Status, error) {
	return svc.applyV4(api, true)
}

func (svc *APIs) applyV4(api *v1alpha1.ApiV4Definition, dryRun bool) (*base.Status, error) {
	url := svc.AutomationTarget("apis").WithQueryParam("dryRun", strconv.FormatBool(dryRun))

	// populateIDs only set this if it an upgraded API
	setHridWithUUID := api.Spec.ID != ""

	automation := api.Spec.ToAutomation()

	// If an ID is set, it can be:
	// 1. export/import case (user is trying to manage an existing API)
	// 2. CRD created before upgrade to AutomationAPI (prior to 4.12)
	// Anyhow, use HRID to store ID and pass setHridWithUUID query param
	// To tell AutomationAPI that this it is not an HRID-managed resource
	// UUIDs are already computed
	if setHridWithUUID {
		automation.HRID = api.Spec.ID
		url = url.WithQueryParam("hridContainsUUID", strconv.FormatBool(true))
	}

	status := new(v4.AutomationStatus)
	if err := svc.HTTP.Put(url.String(), automation, status); err != nil {
		return nil, err
	}

	if !setHridWithUUID {
		k8s.AddAutomationAPIManagedCondition(api)
	}
	return &base.Status{
		ApiStatus: status.ApiStatus,
		Plans:     status.PlansAsMap(),
	}, nil
}

func (svc *APIs) UpdateState(apiID string, action model.Action) error {
	url := svc.EnvV1Target("apis").WithPath(apiID).WithQueryParam(stateActionParam, string(action))
	return svc.HTTP.Post(url.String(), nil, nil)
}

func (svc *APIs) DeleteV2(apiID string) error {
	url := svc.EnvV1Target("apis").WithPath(apiID).WithQueryParams(deleteParams)
	return svc.HTTP.Delete(url.String(), nil)
}

func (svc *APIs) DeleteV4(api core.ApiDefinitionObject) error {
	apiID, hridContainsUUID := getApiID(api)
	url := svc.AutomationTarget("apis").WithPath(apiID).
		WithQueryParam("hridContainsUUID", strconv.FormatBool(hridContainsUUID))
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

func getApiID(api core.ApiDefinitionObject) (string, bool) {
	if k8s.IsAutomationAPIManaged(api) {
		return refs.NewNamespacedNameFromObject(api).HRID(), false
	}
	return api.GetID(), true
}
