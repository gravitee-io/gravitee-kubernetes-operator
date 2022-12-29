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
	"net/http"

	kModel "github.com/gravitee-io/gravitee-kubernetes-operator/api/model"

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

var importParams = map[string]string{
	"definitionVersion": "2.0.0",
}

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

func (svc *APIs) GetByCrossID(crossID string) (*model.ApiListItem, error) {
	url := svc.EnvTarget("apis").WithQueryParam(crossIDParam, crossID)
	apis := new([]model.ApiListItem)

	if err := svc.HTTP.Get(url.String(), apis); err != nil {
		return nil, err
	}

	if len(*apis) == 0 {
		return nil, errors.NewNotFoundError()
	}

	return &(*apis)[0], nil
}

func (svc *APIs) GetByID(apiID string) (*model.ApiEntity, error) {
	url := svc.EnvTarget("apis").WithPath(apiID)
	api := new(model.ApiEntity)

	if err := svc.HTTP.Get(url.String(), api); err != nil {
		return nil, err
	}

	return api, nil
}

func (svc *APIs) Import(method string, spec *kModel.Api) (*model.ApiEntity, error) {
	url := svc.EnvTarget("apis/import").WithQueryParams(importParams)
	api := new(model.ApiEntity)
	fun := svc.getImportFunc(method)

	if err := fun(url.String(), spec, api); err != nil {
		return nil, err
	}

	return api, nil
}

func (svc *APIs) getImportFunc(method string) func(string, any, any) error {
	if method == http.MethodPost {
		return svc.HTTP.Post
	}
	return svc.HTTP.Put
}

func (svc *APIs) UpdateState(apiID string, action model.Action) error {
	url := svc.EnvTarget("apis").WithPath(apiID).WithQueryParam(stateActionParam, string(action))
	return svc.HTTP.Post(url.String(), nil, nil)
}

func (svc *APIs) Delete(apiID string) error {
	url := svc.EnvTarget("apis").WithPath(apiID).WithQueryParams(deleteParams)
	return svc.HTTP.Delete(url.String(), nil)
}

func (svc *APIs) SetKubernetesContext(apiID string) error {
	url := svc.EnvTarget("apis").WithPath(apiID).WithPath("definition-context")
	return svc.HTTP.Post(url.String(), model.NewKubernetesContext(), nil)
}
